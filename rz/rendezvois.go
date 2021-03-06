// Package rz let one use pure procedural (rendez-vois-like)
// paradigm in communicating between asynchronous processes
// (goroutines).
package rz

import (
	"fmt"
	"reflect"
	"sync"
)

type syncstruct struct {
	vars []interface{}
}

type synchan chan syncstruct

// SyncPoint is the key concept. It provides a target for synchronous
// communication between processes through its Send() method.
type SyncPoint struct {
	fv      reflect.Value
	in, out synchan
}

// NewSyncPoint creates and initializes new syncpoint.
//    server.syncpoint = NewSyncPoint(server.method)
// fn specifies a handler function; actual parameters
// of Send() request MUST correspond to handler's formal
// parameters list. If fn is nil, no special handling
// is performing and empty result list is returned
// (it's fast and could be used just for syncronization
// purposes)
//    var syncpoint = NewSyncPoint(nil)
func NewSyncPoint(fn interface{}) *SyncPoint {
	var fv = reflect.ValueOf(fn)
	if fv.IsValid() && fv.Kind() != reflect.Func {
		panic(fmt.Sprintf("NewSyncPoint %T is not a valid function", fn))
	}
	return &SyncPoint{fv, make(synchan), make(synchan)}
}

// Send transfers arguments and callback function to syncpoint handler. Say:
//    result = server.syncpoint.Send(client.method, arg1, arg2)
// Then it waits for communication and finally retuns a slice of the callback's results.
func (self *SyncPoint) Send(fn interface{}, args ...interface{}) []interface{} {
	self.in <- syncstruct{args}
	var out syncstruct
	out = <-self.out
	return _call(nil, reflect.ValueOf(fn), out.vars)
}

// Accept waits for the syncpoint's Send, passes transferred args
// to handler function and finally transfers its result slice back
// to sender.
//    server.syncpoint.Accept()
// It always suspends the calling process until a Send()
// request, thus, providing purely synchronous communication.
func (self *SyncPoint) Accept() {
	var ss syncstruct
	ss = <-self.in
	ss.vars = _call(nil, self.fv, ss.vars)
	self.out <- ss
}

// Accept function accepts requests to multiple syncpoints. Thus it
// provides an asynchronous communication (i.e. the usual case).
//    Accept(blocking, &lock, server.syncpoint1, server.syncpoint2)
// If blocking is true, it suspends the calling process until a Send()
// request is issued to any of syncpoint from the list. Otherwise
// (not blocking, just like go's select with default clause) it lets the
// calling process continue.
// If any of syncpoints should access single shared resource, lock must
// be specified. If each syncpoint guards its own resource, lock could be
// nil.
func Accept(blocking bool, lock sync.Locker, spoints ...*SyncPoint) {
	var cases []reflect.SelectCase
	if blocking {
		cases = make([]reflect.SelectCase, len(spoints))
	} else {
		cases = make([]reflect.SelectCase, len(spoints)+1)
		cases[len(cases)-1] = reflect.SelectCase{
			Dir: reflect.SelectDefault,
		}
	}
	for i, sp := range spoints {
		cases[i] = reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(sp.in),
		}
	}
	if chosen, recv, recvOK := reflect.Select(cases); recvOK {
		var ss = recv.Interface().(syncstruct)
		var sp = spoints[chosen]
		ss.vars = _call(lock, sp.fv, ss.vars)
		sp.out <- ss
	}
}

// CallPanic is issued when panic raised in calling a handler.
type CallPanic struct {
	Func reflect.Value
	Err  interface{}
}

func (self CallPanic) Error() string {
	return fmt.Sprintf("CallPanic %s %s", self.Func.Kind(), self.Err)
}

func _call(lock sync.Locker, fv reflect.Value, args []interface{}) []interface{} {
	if !fv.IsValid() {
		return nil
	}
	var a = make([]reflect.Value, len(args))
	for i, arg := range args {
		a[i] = reflect.ValueOf(arg)
	}
	defer func() {
		if r := recover(); r != nil {
			panic(CallPanic{fv, r})
		}
	}()
	if lock != nil {
		a = _lcall(lock, fv, a)
	} else {
		a = fv.Call(a)
	}
	var res = make([]interface{}, len(a))
	for i, r := range a {
		res[i] = r.Interface()
	}
	return res
}

func _lcall(lock sync.Locker, fv reflect.Value, args []reflect.Value) []reflect.Value {
	if !fv.IsValid() {
		return nil
	}
	defer func() {
		if lock != nil {
			lock.Unlock()
		}
	}()
	if lock != nil {
		lock.Lock()
	}
	return fv.Call(args)
}

// rendevois1_test.go
package rz

import (
	"fmt"
	"testing"
	"time"
)

//-----[ SERVER ]-------------------------------------------------------------//

type ProducerTask1 struct {
	syncpoint, termination *SyncPoint
	// product number is the shared resource
	product int
}

func NewProducerTask1() (prod *ProducerTask1) {
	prod = &ProducerTask1{product: 0}
	prod.syncpoint = NewSyncPoint(prod.answer)
    // termination syncpoint is used for syncronization only
    // no need in handler function
	prod.termination = NewSyncPoint(nil)
    return
}

// the main server process
func (self *ProducerTask1) run(max int, q quiet) {
	q.println("*** Producer started")
    // just wait for and accept requests
    // for product numbers until max is reached
	for self.product < max {
		self.syncpoint.Accept()
	}
    // wait for and accept request for termination
	self.termination.Accept()
    // once request is accepted, quit
	q.println("Producer ended")
}

// the product request handler
// filter specifies what kind of number is requested: odd, even or any
func (self *ProducerTask1) answer(filter string) (out int, ok bool) {
	if product := self.product; passed(product, filter) {
        // the request can be satisfied
        // produce new number
		self.product++
        // return success
		return product, true
	}
    // otherwise refuse
	return -1, false
}

//-----[ CLIENT ]-------------------------------------------------------------//

type ConsumerTask1 struct {
	producer *ProducerTask1
	filter   string
}

// the main client process
func (self *ConsumerTask1) run(pace time.Duration, q quiet, tb testing.TB) {
	q.println("*** Consumer started", self.filter)

	for {
        // send a request and handle the answer
		var res = self.producer.syncpoint.Send(func(product int, ok bool) {
			if !ok {
                // request refused -- no approriate product number at the moment
				//q.println("Consumer", self.filter, "got NO product")
                // have to wait
				time.Sleep(10 * pace * time.Millisecond)
			} else if passed(product, self.filter) {
                // request satisfied and checked
				q.println("Consumer got", self.filter, product)
				time.Sleep(pace * time.Millisecond)
			} else {
                // something went wrong -- producer answerd with
                // inapproriate number (should never happen)
				tb.Error("Consumer", self.filter, "got WRONG", product)
				time.Sleep(5 * pace * time.Millisecond)
			}
		},
        // request parameters MUST correspond
        // the handler's formal argument list
        self.filter)
        // we don't use the result now
        // it's here just for demo purposes
		_ = res
	}
    // not reach here
}

func TestMain1(t *testing.T) {
    // create tasks
	var prod = NewProducerTask1()
	var cons1 = ConsumerTask1{prod, "odd"}
	var cons2 = ConsumerTask1{prod, "even"}
    // start tasks
    // consumer tasks request the single syncpoint
    // so they are mutally synchronized (may access
    // shared resources safely)
	go cons1.run(10, false, t)
	go prod.run(100, false)
    // another consumer started with delay to
    // demonstrate that it in sync anyway
	time.Sleep(1000 * time.Millisecond)
	go cons2.run(1, false, t)
    // now send a request *right from the main thread*
    // and wait for termination (see ProducerTask1.run)
	prod.termination.Send(func() {
		fmt.Println("Terminated")
	})
	// just let things get done
	time.Sleep(1000 * time.Millisecond)
}

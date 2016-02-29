// rendevois1_test.go
package rz

import (
	"testing"
	"time"
	"fmt"
)

type ProducerTask1 struct {
	syncpoint, termination *SyncPoint
	product                int
}

func NewProducerTask1() (prod *ProducerTask1) {
	prod = &ProducerTask1{product: 0}
	prod.syncpoint = NewSyncPoint(prod.answer)
	prod.termination = NewSyncPoint(nil)
	return
}

func (self *ProducerTask1) run(max int, q quiet) {
	q.println("*** Producer started")

	for self.product < max {
		self.syncpoint.Accept()
	}

	self.termination.Accept()

	q.println("Producer ended")
}

func (self *ProducerTask1) answer(filter string) (out int, ok bool) {
	if product := self.product; passed(product, filter) {
		self.product++
		return product, true
	}
	return -1, false
}

//----------------------------------------------------------------------------//

type ConsumerTask1 struct {
	producer *ProducerTask1
	filter   string
}

func (self *ConsumerTask1) run(pace time.Duration, q quiet, tb testing.TB) {
	q.println("*** Consumer started", self.filter)

	for {
		var res = self.producer.syncpoint.Ask(func(product int, ok bool) {
			if !ok {
				//q.println("Consumer", self.filter, "got NO product")
				time.Sleep(10 * pace * time.Millisecond)
			} else if passed(product, self.filter) {
				q.println("Consumer got", self.filter, product)
				time.Sleep(pace * time.Millisecond)
			} else {
				tb.Error("Consumer", self.filter, "got WRONG", product)
				time.Sleep(5 * pace * time.Millisecond)
			}
		}, self.filter)
		_ = res
	}

	q.println("Consumer ended", self.filter)
}

func TestMain1(t *testing.T) {
	var prod = NewProducerTask1()
	var cons1 = ConsumerTask1{prod, "odd"}
	var cons2 = ConsumerTask1{prod, "even"}
	go cons1.run(10, false, t)
	go prod.run(100, false)
	time.Sleep(1000 * time.Millisecond)
	go cons2.run(1, false, t)
	prod.termination.Ask(func() {
		fmt.Println("Terminated")
	})
}

package rz

import (
	"fmt"
	"sync"
	"time"
	"testing"
)

//----------------------------------------------------------------------------//

type ProducerTask2 struct {
	syncpoint1, syncpoint2 *SyncPoint
	product                int
}

func NewProducerTask2() (prod *ProducerTask2) {
	prod = &ProducerTask2{product: 0}
	prod.syncpoint1 = NewSyncPoint(prod.answer1)
	prod.syncpoint2 = NewSyncPoint(prod.answer2)
	return
}

func (self *ProducerTask2) run() {
	fmt.Println("*** Producer started")

	var mx = &sync.Mutex{}
	for {
		Accept(true, mx, self.syncpoint1, self.syncpoint2)
	}

	fmt.Println("Producer ended")

}

func (self *ProducerTask2) answer1(filter string) (out int, ok bool) {
	var product = self.product
	defer func() {
		//fmt.Println("Producer gives", product)
	}()
	if passed(product, filter) {
		self.product++
		return product, true
	}
	return -1, false
}

func (self *ProducerTask2) answer2(filter string) (out int, ok bool) {
	var product = self.product
	defer func() {
		//fmt.Println("Producer gives", product)
	}()
	if passed(product, filter) {
		self.product++
		return product + 10000, true
	}
	return -2, false
}

//----------------------------------------------------------------------------//

type ConsumerTask21 struct {
	producer *ProducerTask2
	filter   string
}

func (self *ConsumerTask21) run(max int, pace time.Duration) {
	fmt.Println("*** Consumer started", self.filter)

	for done := false; !done; {
		self.producer.syncpoint1.Send(func(product int, ok bool) {
			if !ok {
				//fmt.Println("Consumer", self.filter, "got NO product")
				time.Sleep(10 * pace * time.Millisecond)
			} else if passed(product, self.filter) {
				fmt.Println("Consumer got", self.filter, product)
				if product < max {
					time.Sleep(pace * time.Millisecond)
				} else {
					done = true
				}
			} else {
				fmt.Println("Consumer", self.filter, "got WRONG", product)
				time.Sleep(5 * pace * time.Millisecond)
			}
		}, self.filter)
	}
	termination.Send(func(msg string) {
		fmt.Println(msg, "for", self.filter)
	})
	fmt.Println("Consumer ended", self.filter)
}

type ConsumerTask22 struct {
	producer    *ProducerTask2
	filter      string
	termination *SyncPoint
}

func (self *ConsumerTask22) run(pace time.Duration) {
	self.termination = NewSyncPoint(func() string {
		return fmt.Sprintf("Termination 2 confirmed")
	})

	fmt.Println("*** Consumer started", self.filter)

	for done := false; !done; {
		self.producer.syncpoint2.Send(func(product int, ok bool) {
			if !ok {
				//fmt.Println("Consumer", self.filter, "got NO product")
				time.Sleep(10 * pace * time.Millisecond)
			} else if passed(product, self.filter) {
				fmt.Println("Consumer got", self.filter, product)
			} else {
				fmt.Println("Consumer", self.filter, "got WRONG", product)
				time.Sleep(5 * pace * time.Millisecond)
			}
		}, self.filter)
		Accept(false, nil, self.termination)
	}
	fmt.Println("Consumer ended", self.filter)
}

//----------------------------------------------------------------------------//

var termination *SyncPoint

func TestMain2(*testing.T) {
	var prod = NewProducerTask2()
	var cons1 = &ConsumerTask21{prod, "odd"}
	var cons2 = &ConsumerTask22{prod, "even", nil}
	termination = NewSyncPoint(func() string {
		return fmt.Sprintf("Termination 1 confirmed")
	})

	go cons1.run(100, 10)
	go prod.run()
	time.Sleep(1000 * time.Millisecond)
	go cons2.run(1)

	termination.Accept() // cons1

	cons2.termination.Send(func(msg string) {
		fmt.Println(msg, "for main()")
	})
}


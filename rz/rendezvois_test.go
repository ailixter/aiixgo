package rz

import (
	"fmt"
    "time"
	//"testing"
)

func ExampleSyncPoint/*_Send*/() {
    //  the shared resource
	var resource int = 0
    //  the syncpoint to access the resource
    //  the passed func is a syncpoint handler
    //  (it could be nil, then it returns
    //  a nil result slice)
	var syncpoint = NewSyncPoint(func(arg1 int) (arg2 int) {
        //  the result list MUST correspond to the
        //  formal argument list of sender's handler
        //  this one returns the resource value
        //  plus the parameter passed by a request
		arg2 = resource + arg1
		return
	})
    //  start the 'consumer' (client) process
	go func() {
		for {
            //  send the request for communication
            //  with its handler and parameter(s)
            //  (it could be nil, just for the purpose
            //  of synchronization)
			syncpoint.Send(func(arg2 int) {
                // just print the result of syncpoint's handler
				fmt.Println(arg2)
			},
            //  passed parameters MUST correspond
            //  the formal parameter list of
            //  syncpoint's handler
            /*arg1*/ 1000)
		}
	}()
    //  start the 'producer' (server) process
	for resource < 10 {
        //  accept requests to this syncpoint
		syncpoint.Accept()
        //  modify the resource
		resource++
        //  simulate some load
        time.Sleep(10*time.Microsecond)
	}

	// Output:
    // 1000
    // 1001
    // 1002
    // 1003
    // 1004
    // 1005
    // 1006
    // 1007
    // 1008
    // 1009
}

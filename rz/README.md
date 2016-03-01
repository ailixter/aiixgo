#Asymmetric Rendez-Vois for Concurrent Processes

It makes it possible to inter-communicate the processes in pure procudural manner.

## Trivial Example
``` go
//  the shared resource
var resource int = 0
//  the syncpoint to access the resource
//  the passed func is a syncpoint handler
var syncpoint = NewSyncPoint(func(arg1 int) (arg2 int) {
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
}

```

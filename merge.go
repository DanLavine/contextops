package contextops

import (
	"context"
	"reflect"
)

//	PARAMETERS:
//	- contexts - any number of contexts to merge so that any closed, will trigger a closure of the returned context
//
//	RETURNS:
//	- context.Context - context that will be closed if any parameter contexts are closed.
//
// MergeForDone can be used when wanting to monitor for any number of contexts that can cancel a long
// running operation into a single context. It is required that one of the original contexts passed in
// as a paremeter closes to ensure the async resources setup from this function are cleaned up properly.
//
// Example: when wanting to obtain a lock from a service, the service will need to respond to clients in a number of ways, including:
//  1. On a server shutdown, tell the clients to reconnect
//  2. When a client disconnects
//  3. Context to cleanup async resources setup within this function
func MergeForDone(contexts ...context.Context) context.Context {
	oneCtx, oneCancel := context.WithCancel(context.Background())

	go func() {
		cases := []reflect.SelectCase{}

		if len(contexts) == 0 {
			oneCancel()
		} else {
			for _, ctx := range contexts {
				cases = append(cases, reflect.SelectCase{
					Dir:  reflect.SelectRecv,
					Chan: reflect.ValueOf(ctx.Done()),
				})
			}

			// don't care about any of the return values here. This just needs to be triggered
			_, _, _ = reflect.Select(cases)
			oneCancel()
		}
	}()

	return oneCtx
}

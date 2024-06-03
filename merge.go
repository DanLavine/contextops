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

//	PARAMETERS:
//	- mainCtx - the context that all other contexts will be merged into. The Values of this context are preserved
//	- contexts - all other contexts whos closeure of the `Done()` channel will trigger a closure of the returned contex
//
//	RETURNS:
//	- context - the merged context that is canceled when any of the other contexts cancel
//	- func() - cancel function that must be called to not leak goroutines
//
// MergDone will merge the 'Done()' opertation of any number of contexts into the mainCtx. The
// values of mainCtx are preserved, but no values from any of the contexts will be preserved
//
// Example for when to use this operation:
//  1. An API's http.Request.Context() has a number of `context.Value(...)` set for Open Telemitry, Logging or other metadata assigned (mainCtx)
//  2. The same API can be canceled via a server shutdown request (...context), or some other operation that has no data assigned, just a shutdown trigger
//  3. Now rather than passing all possible context operations around into the possible shutdown cases, they can all be merged into 1 context and handled as a single unit
func MergeDone(mainCtx context.Context, contexts ...context.Context) (context.Context, func()) {
	mergeCtx, oneCancel := context.WithCancel(mainCtx)

	cases := []reflect.SelectCase{
		{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(mainCtx.Done()),
		},
	}

	for _, ctx := range contexts {
		cases = append(cases, reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(ctx.Done()),
		})
	}

	go func() {
		// don't care about any of the return values here. This just needs to be triggered
		_, _, _ = reflect.Select(cases)
		oneCancel()
	}()

	return mergeCtx, oneCancel
}

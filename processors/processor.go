// Package processors specifies the Processor functions.
//
// The http.ResponseWriter for "Processor" implementations will be a proxy object that tracks
// if response sending has been initiated -- call to Write(...).
package processors

import "net/http"

// A Processor is expected to do everything needed to process the Request into the ResponseWriter.
// If the processors is part of a CRUD app, then it should delegate to code to do the business logic.
// If the processors is part of a Proxy app, then it should delegate to the "source" app with
// "pre" and/or "post" processing delegated to the appropriate business logic.
//
// If Write(...) hasn't been called, then the "Processor" function's error will trigger a
// http.Error call with the results (note: any invalid status code will turn into a 500!).
// If Write(...) has been called, then resulting error text will simply be logged!
type Processor func(http.ResponseWriter, *http.Request) (statusCode int, err error)

// Return404 is an implementation of the Processor function, to be used if no "Processor" is
// passed to the Builder's Build method.
//
//goland:noinspection GoUnusedExportedFunction
func Return404(w http.ResponseWriter, r *http.Request) (statusCode int, err error) {
	http.Error(w, "path '"+r.URL.Path+"' Not Found", http.StatusNotFound)
	return // zero values!
}

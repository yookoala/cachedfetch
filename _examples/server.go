package main

import (
	"fmt"
	"net/http"
)

func ExampleServer() (mux *http.ServeMux) {

	mux = http.NewServeMux()

	// produce count with a channel
	counts := make(chan int64)
	go func() {
		for i := int64(1); ; i++ {
			counts <- i
		}
	}()

	// returns a handler that displays a simple notice
	getNoticePage := func(notice string) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, notice)
		})
	}

	// a simple page the return ever changing content
	getCounterPage := func(name string) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			count := <-counts
			fmt.Fprintf(w, "%s - Counter: %03d", name, count)
		})
	}

	// bind example paths
	mux.Handle("/example/1", getNoticePage("Hello, example 1"))
	mux.Handle("/example/2", getCounterPage("Example 2"))
	mux.Handle("/example/3", getCounterPage("Example 3"))
	for i := 1; i <= 10; i++ {
		mux.Handle(fmt.Sprintf("/example/4/%d", i),
			getCounterPage(fmt.Sprintf("Example 4.%d", i)))
	}
	mux.Handle("/example/5", getCounterPage("Example 5"))
	mux.Handle("/example/6", getCounterPage("Example 6"))
	return
}

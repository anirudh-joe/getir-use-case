package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/getircase/controller"
)

func main() {
	var memServer, mongoServer http.Handler
	memServer = new(controller.MemDbGate)
	mongoServer = new(controller.MongoDbGate)
	// Add request handlers for the given url path
	http.Handle("/in-memory", memServer)
	http.Handle("/mongo", mongoServer)
	// serve http requests on 8080
	log.Fatal(http.ListenAndServe(":8080", nil))
	// Run the server indefinitely
	// till os kills or user interrupt
	done := make(chan os.Signal, 2)
	signal.Notify(done, os.Interrupt)
	signal.Notify(done, os.Kill)
	_ = <-done
	fmt.Printf("\nThats all folks.\n")
}

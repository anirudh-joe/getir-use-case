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
	http.Handle("/in-memory", memServer)
	http.Handle("/mongo", mongoServer)

	log.Fatal(http.ListenAndServe(":8080", nil))

	done := make(chan os.Signal, 2)
	signal.Notify(done, os.Interrupt)
	signal.Notify(done, os.Kill)
	_ = <-done
	fmt.Printf("\nThats all folks.\n")
}

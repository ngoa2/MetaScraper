package main

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/UW-Info-441-Winter-Quarter-2020/homework-ngoa2/servers/gateway/handlers"
)

//main is the main entry point for the server
func main() {
	/* TODO: add code to do the following
	- Read the ADDR environment variable to get the address
	  the server should listen on. If empty, default to ":80"
	- Create a new mux for the web server.
	- Tell the mux to call your handlers.SummaryHandler function
	  when the "/v1/summary" URL path is requested.
	- Start a web server listening on the address you read from
	  the environment variable, using the mux you created as
	  the root handler. Use log.Fatal() to report any errors
	  that occur when trying to start the web server.
	*/

	// Reads the ADDR environment variabl to get the address the
	// server should listen on. Defaults to ":80" wheny empty
	addr := os.Getenv("ADDR")
	tlscert := os.Getenv("TLSCERT")
	tlskey := os.Getenv("TLSKEY")

	if len(tlscert) == 0 || len(tlskey) == 0 {
		err := errors.New("TLSCERT and TLSKEY env variables have not been set")
		log.Fatal(err)
	}
	if len(addr) == 0 {
		// changeed to :443 for standard HTTPS port
		addr = ":443"
	}

	// Creates a new mux for the web server and tells it to call
	// handlers.SummaryHandler when "/v1/summary" URL path is
	// requested
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/summary", handlers.SummaryHandler)

	log.Printf("server is listening at %s...", addr)
	log.Fatal(http.ListenAndServeTLS(addr, tlscert, tlskey, mux))

}

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"server/pages"
	"server/vrahos"

	"golang.org/x/net/http2"
)

func main() {

	components := []vrahos.Component{
		pages.Document{},
		pages.IndexPage{},
	}

	mux := http.NewServeMux()

	vrahos.Vrahos(mux, components)

	port := os.Getenv("PORT")
	fmt.Println("Listening " + port)
	server := &http.Server{
		Addr:    "0.0.0.0:" + port,
		Handler: mux,
	}
	http2.ConfigureServer(server, nil)

	log.Fatal(server.ListenAndServe())
}
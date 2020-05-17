package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type MyServer struct {
	httpServer *http.Server
}

func (server *MyServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

	log.Print(request.Method, " ", request.URL.RequestURI())
	dummy := make(map[string]interface{})
	dummy["ala"] = "bala"
	dummy["integer"] = 13
	js, err := json.Marshal(dummy)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.Write(js)
}

func (server *MyServer) setup() {
	log.Print("setup")
	defer log.Print("Finish setup")

}

func (server *MyServer) doRun(addr string) {
	defer log.Print("Funished running server")

	server.httpServer = &http.Server{Addr: addr, Handler: server}
	log.Print("Running server on address:", addr)
	server.httpServer.ListenAndServe()
}

func (server *MyServer) RunAsync(addr string) {
	go server.doRun(addr)
}

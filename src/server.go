package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
)

type MyServer struct {
	httpServer *http.Server
	muxHandler *http.ServeMux

	stopSignal     chan int
	serverFinished chan int
}

func (server *MyServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

	log.Print(request.Method, " ", request.URL.RequestURI())
	dummy := make(map[string]interface{})
	dummy["title"] = "This is my Golang demo"
	dummy["description"] = "For sparta!"
	js, err := json.Marshal(dummy)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.Write(js)
}

func (server *MyServer) Init() {
	log.Print("setup")
	defer log.Print("Finish setup")

	server.stopSignal = make(chan int)
	server.serverFinished = make(chan int)

	server.muxHandler = http.NewServeMux()
}

func (server *MyServer) waitForStopWorker() {
	log.Print("Waiting for stop event")
	signal := <-server.stopSignal
	log.Printf("Stop signal [%v] received", signal)

	if signal != 0 {
		return
	}

	if err := server.httpServer.Shutdown(context.Background()); err != nil {
		log.Printf("Error [%v] shuting down server", err)
		return
	}

	log.Print("Server was shut down successfully")
}

func (server *MyServer) onFinishServer() {
	server.serverFinished <- 0
}

func (server *MyServer) doRun(addr string) {

	defer server.onFinishServer()
	defer log.Print("Finished running server")

	server.muxHandler.Handle("/", server)

	server.httpServer = &http.Server{Addr: addr, Handler: server.muxHandler}
	log.Print("Running server on address:", addr)

	go server.waitForStopWorker()
	if err := server.httpServer.ListenAndServe(); err != http.ErrServerClosed {
		log.Printf("Error [%v] starting server", err)
		server.stopSignal <- -1
	}
}

func (server *MyServer) RunAsync(addr string) {
	go server.doRun(addr)
}

func (server *MyServer) Stop() {

	log.Printf("Stoping server")
	server.stopSignal <- 0

	log.Print("Wait for finish")
	<-server.serverFinished

	log.Print("Finish signal recerived")
}

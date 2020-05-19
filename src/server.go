package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type MyServer struct {
	httpServer *http.Server
	muxHandler *http.ServeMux

	stopSignal     chan int
	serverFinished chan int
}

type PersonRESTApi struct {
	PersonDB IPersonDB
}

func (api *PersonRESTApi) Init() {
	api.PersonDB, _ = NewPersonDB()
}

func (api *PersonRESTApi) onPOST(writer http.ResponseWriter, request *http.Request) {

	url := request.URL.RequestURI()
	if url != "/services/restapi/1.0/Persons/" {
		http.Error(writer, "Invalid service request", http.StatusBadRequest)
		return
	}

	var person Person
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&person)
	if err != nil {
		http.Error(writer, "Invalid service request", http.StatusBadRequest)
		return
	}

	err = api.PersonDB.Insert(&person)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	result := make(map[string]interface{})
	result["error"] = 0
	result["message"] = "all good"

	js, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.Write(js)
}

func (api *PersonRESTApi) getPersonCnpFromURL(url string) string {
	var paths []string
	paths = strings.Split(url, "/")
	return paths[len(paths)-1]
}

func (api *PersonRESTApi) onGET(writer http.ResponseWriter, request *http.Request) {

	url := request.URL.RequestURI()
	if !strings.HasPrefix(url, "/services/restapi/1.0/Persons/") {
		http.Error(writer, "Invalid service request", http.StatusBadRequest)
		return
	}

	cnp_str := api.getPersonCnpFromURL(url)
	cnp, _ := strconv.Atoi(cnp_str)

	person, _ := api.PersonDB.Get(cnp)

	js, err := json.MarshalIndent(person, "", "    ")
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.Write(js)
}

func (api *PersonRESTApi) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		api.onPOST(writer, request)
		return
	}

	if request.Method == "GET" {
		api.onGET(writer, request)
		return
	}

	http.Error(writer, "Invalid method for this url", http.StatusBadRequest)
}

type DummyHandler struct{}

func (handler *DummyHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
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

func (server *MyServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	log.Print(request.Method, " ", request.URL.RequestURI())
	server.muxHandler.ServeHTTP(writer, request)
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

func (server *MyServer) registerHandlers() {
	api := &PersonRESTApi{}
	api.Init()
	server.muxHandler.Handle("/", api)
}

func (server *MyServer) doRun(addr string) {

	defer server.onFinishServer()
	defer log.Print("Finished running server")

	server.registerHandlers()

	server.httpServer = &http.Server{Addr: addr, Handler: server}
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

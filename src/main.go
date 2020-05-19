package main

import (
	"log"
	"time"
)

func setupLogger() {
	var flags int
	flags |= log.Ltime
	flags |= log.Ldate
	flags |= log.Lmicroseconds
	//flags |= log.LUTC
	flags |= log.Lshortfile
	log.SetFlags(flags)
}

func main() {
	setupLogger()
	log.Print("Begin POC")
	defer log.Print("Exit main")

	db, _ := NewPersonDB()
	db.Insert(&Person{Name: "TestName", LastName: "TestLastName", CNP: 123})

	p, _ := db.Get(123)
	log.Print(p)

	server := &MyServer{}
	server.Init()

	server.RunAsync("192.168.1.192:80")

	time.Sleep(60 * time.Minute)

	server.Stop()
}

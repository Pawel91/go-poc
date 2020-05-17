package main

import (
	"log"
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
}

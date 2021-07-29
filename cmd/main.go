package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dkucheru/Calendar/app"
	"github.com/dkucheru/Calendar/structs"
)

var CheckSignals chan os.Signal

func main() {
	migrateDown := flag.Bool("down", false, "call to sql-migrate down")
	flag.Parse()
	structs.GlobalId = 1
	rand.Seed(time.Now().Unix())
	appNew, err := app.New(*migrateDown)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		err = appNew.Run()
		if err != nil {
			log.Fatal(err)
		}
	}()
	defer appNew.Stop()
	log.Printf("Started server")
	CheckSignals = make(chan os.Signal, 1)
	signal.Notify(CheckSignals, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Println(fmt.Sprint(<-CheckSignals))
	log.Println("Stopping API server.")
}

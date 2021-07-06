package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/dkucheru/Calendar/app"
	"github.com/dkucheru/Calendar/structs"
)

func main() {
	structs.GlobalId = 1

	rand.Seed(time.Now().Unix())
	appNew, err := app.New()
	if err != nil {
		log.Fatal(err)
	}

	err = appNew.Run()
	if err != nil {
		log.Fatal(err)
	}
}

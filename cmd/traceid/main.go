package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Println(uuid.Must(uuid.NewV7()))
		os.Exit(0)
	}

	uuid, err := uuid.Parse(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	t := time.Unix(uuid.Time().UnixTime())
	fmt.Println(t)
}

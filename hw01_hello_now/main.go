package main

import (
	"fmt"
	"log"
	"time"

	"github.com/beevik/ntp"
)

func main() {
	currTime := time.Now()
	ntpTime, err := ntp.Time("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		log.Fatal("Error:\n", err)
	}

	fmt.Printf("current time: %v\nexact time: %v\n", currTime.Round(0), ntpTime.Round(0))
}

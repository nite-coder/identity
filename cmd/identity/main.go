package main

import (
	"fmt"

	"github.com/jasonsoft/log"
)

func main() {
	log.SetAppID("identity") // unique id for the app
	fmt.Println("done")
}

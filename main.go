package main

import (
	"fmt"
	"net/http/httputil"

	"github.com/jasonsoft/log"
	"github.com/jasonsoft/log/handlers/console"

	"github.com/jasonsoft/napnap"
)

func main() {
	log.SetAppID("identity") // unique id for the app

	clog := console.New()
	log.RegisterHandler(clog, log.AllLevels...)

	nap := napnap.New()
	nap.Use(napnap.NewHealth())
	router := napnap.NewRouter()
	router.All("/auth", homeEndpoint)
	nap.Use(router)
	server := napnap.NewHttpEngine(":16000")
	nap.Run(server)
}

func homeEndpoint(c *napnap.Context) {

	requestDump, err := httputil.DumpRequest(c.Request, true)
	if err != nil {
		fmt.Println(err)
	}

	c.RespHeader("X-Auth-JWT", "jwt-1234567")

	log.Debugf("main: request: %s", string(requestDump))
	c.SetStatus(200)

}

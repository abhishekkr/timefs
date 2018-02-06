package main

import (
	"fmt"
	"log"
	"os"

	"github.com/abhishekkr/gol/golenv"
	"gopkg.in/alecthomas/kingpin.v2"

	timefsClient "github.com/abhishekkr/timefs/client/timefsClient"
	timedot "github.com/abhishekkr/timefs/timedot"
)

var (
	TIMEFS_AT = golenv.OverrideIfEnv("TIMEFS_AT", "127.0.0.1:7999")

	app = kingpin.New("timefs-client", "a client to use (or test) timefs-server")

	serverIP = app.Flag("server", "Server address.").Default(TIMEFS_AT).String()

	dummy     = app.Command("dummy", "to test timefs-server")
	dummyMode = dummy.Arg("mode", "create|read").Required().String()
)

func dummyCli(client *timedot.TimeFSClient, mode string) {
	if mode == "create" {
		timefsClient.DummyCreate(client)
	} else if mode == "read" {
		timefsClient.DummyRead(client)
	} else {
		fmt.Printf(`wrong dummy mode: %s
  available options are {create,read}
  example: tfclient --server='127.0.0.1:7999' dummy create%s`, mode, "\n")
	}
}

func main() {
	log.Println("starting client...")
	conn := timefsClient.LinkOpen(TIMEFS_AT)
	defer timefsClient.LinkClose(conn)
	client := timedot.NewTimeFSClient(conn)

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case dummy.FullCommand():
		dummyCli(&client, *dummyMode)

	default:
		log.Println("wrong usage, try running app with 'help'")
	}
}

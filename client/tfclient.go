package main

import (
	"fmt"
	"log"

	"github.com/abhishekkr/gol/golenv"
	"gopkg.in/alecthomas/kingpin.v2"

	timefsClient "github.com/abhishekkr/timefs/client/timefsClient"
	timedot "github.com/abhishekkr/timefs/timedot"
)

var (
	serverIP = kingpin.Flag("server", "Server address.").Required().String()

	flagMode = kingpin.Arg("mode", "mode to run in").Required().String()
	flagAxn  = kingpin.Arg("axn", "create|read").Required().String()
)

func dummyCli(client *timedot.TimeFSClient, axn string) {
	if axn == "create" {
		timefsClient.DummyCreate(client)
	} else if axn == "read" {
		timefsClient.DummyRead(client)
	} else {
		fmt.Printf(`wrong dummy axn: %s
  available options are {create,read}
  example: tfclient --server='127.0.0.1:7999' dummy create%s`, axn, "\n")
	}
}

func main() {
	kingpin.Version("0.1.0")
	kingpin.Parse()
	if *serverIP == "" {
		*serverIP = golenv.OverrideIfEnv("TIMEFS_AT", "127.0.0.1:7999")
	}

	log.Println("starting client...", *serverIP)
	conn := timefsClient.LinkOpen(*serverIP)
	defer timefsClient.LinkClose(conn)
	client := timedot.NewTimeFSClient(conn)

	switch *flagMode {
	case "dummy":
		dummyCli(&client, *flagAxn)

	default:
		log.Println("wrong usage, try running app with 'help'")
	}
}

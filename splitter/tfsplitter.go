package main

import (
	"context"
	"log"
	"net"
	"strings"

	"google.golang.org/grpc"

	"github.com/abhishekkr/gol/golenv"

	timefsClient "github.com/abhishekkr/timefs/client/timefsClient"
	timefsSplit "github.com/abhishekkr/timefs/splitter/timefsSplit"
	timedot "github.com/abhishekkr/timefs/timedot"
)

var (
	TIMEFS_PROXY_PORT = golenv.OverrideIfEnv("TIMEFS_PROXY_PORT", ":7799")
	TIMEFS_BACKENDS   = golenv.OverrideIfEnv("TIMEFS_BACKENDS", "127.0.0.1:7999")

	Clients []timedot.TimeFSClient
)

type Timedots struct {
	savedTimedots []*timedot.Record
}

func main() {
	connectBackends(TIMEFS_BACKENDS)

	conn, err := net.Listen("tcp", TIMEFS_PROXY_PORT)
	if err != nil {
		log.Fatalln("failed to bind at", TIMEFS_PROXY_PORT)
		return
	}

	log.Println("starting server... @", TIMEFS_PROXY_PORT)
	svr := grpc.NewServer()
	timedot.RegisterTimeFSServer(svr, &Timedots{})
	svr.Serve(conn)
}

func connectBackends(backendCSV string) {
	links := strings.Split(backendCSV, ",")
	Clients = make([]timedot.TimeFSClient, len(links))

	for idx, link := range links {
		clientConn := timefsClient.LinkOpen(link)
		defer timefsClient.LinkClose(clientConn)
		Clients[idx] = timedot.NewTimeFSClient(clientConn)
	}
}

func (tym *Timedots) CreateTimedot(c context.Context, input *timedot.Record) (*timedot.TimedotSave, error) {
	go timefsSplit.CreateTimeFS()
	return &timedot.TimedotSave{
		Success: true,
	}, nil
}

func (tym *Timedots) ReadTimedot(filtr *timedot.Record, stream timedot.TimeFS_ReadTimedotServer) error {
	recordChan := make(chan timedot.Record)

	go timefsSplit.GetTimeFS()
	for record := range recordChan {
		err := stream.Send(&record)
		if err != nil {
			return err
		}
	}
	return nil
}

func (tym *Timedots) DeleteTimedot(c context.Context, input *timedot.Record) (*timedot.TimedotDelete, error) {
	panic("WIP")
	return &timedot.TimedotDelete{
		Success: true,
		Count:   0,
	}, nil
}

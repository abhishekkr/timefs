package main

import (
	"context"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/abhishekkr/gol/golenv"

	timefsSplitter "github.com/abhishekkr/timefs/splitter/timefsSplitter"
	timedot "github.com/abhishekkr/timefs/timedot"
)

var (
	TIMEFS_PROXY_PORT = golenv.OverrideIfEnv("TIMEFS_PROXY_PORT", ":7799")
	TIMEFS_CREATE_ID  = golenv.OverrideIfEnv("TIMEFS_CREATE_ID", "recursion")

	createEngine timefsSplitter.CreateEngine
)

type Timedots struct {
	savedTimedots []*timedot.Record
}

func main() {
	defer timefsSplitter.CloseBackends()
	createEngine = timefsSplitter.GetCreateEngine(TIMEFS_CREATE_ID)

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

func (tym *Timedots) CreateTimedot(c context.Context, input *timedot.Record) (*timedot.TimedotSave, error) {
	go createEngine.CreateTimeFS(input)
	return &timedot.TimedotSave{
		Success: true,
	}, nil
}

func (tym *Timedots) ReadTimedot(filtr *timedot.Record, stream timedot.TimeFS_ReadTimedotServer) error {
	recordChan := make(chan timedot.Record)
	go timefsSplitter.GetTimeFS(recordChan, filtr)

	var err error
	for record := range recordChan {
		err = stream.Send(&record)
		if err == io.EOF {
			log.Println("[+] read stream ended by EOF;", err)
			break
		}
		if err != nil {
			log.Println("read stream send failed;", err)
			break
		}
	}
	log.Println("[+] read stream ended")
	return err
}

func (tym *Timedots) DeleteTimedot(c context.Context, input *timedot.Record) (*timedot.TimedotDelete, error) {
	panic("WIP")
	return &timedot.TimedotDelete{
		Success: true,
		Count:   0,
	}, nil
}

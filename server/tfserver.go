package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/abhishekkr/gol/golenv"

	timefs "github.com/abhishekkr/timefs/fs"
	timedot "github.com/abhishekkr/timefs/timedot"
)

var (
	TIMEFS_PORT = golenv.OverrideIfEnv("TIMEFS_PORT", ":7999")
)

type Timedots struct {
	savedTimedots []*timedot.Record
}

func main() {
	conn, err := net.Listen("tcp", TIMEFS_PORT)
	if err != nil {
		log.Fatalln("failed to bind at", TIMEFS_PORT)
		return
	}

	log.Println("starting server... @", TIMEFS_PORT)
	svr := grpc.NewServer()
	timedot.RegisterTimeFSServer(svr, &Timedots{})
	svr.Serve(conn)
}

func (tym *Timedots) CreateTimedot(c context.Context, input *timedot.Record) (*timedot.TimedotSave, error) {
	tym.savedTimedots = append(tym.savedTimedots, input)
	go timefs.CreateRecord(input)
	return &timedot.TimedotSave{
		Success: true,
	}, nil
}

func (tym *Timedots) ReadTimedot(filtr *timedot.Record, stream timedot.TimeFS_ReadTimedotServer) error {
	for _, tymdot := range tym.savedTimedots {
		if !matchtimedot(tymdot, filtr) {
			continue
		}

		err := stream.Send(tymdot)
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

func matchtimedot(tymdot *timedot.Record, filtr *timedot.Record) bool {
	return true
}

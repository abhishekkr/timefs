package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/abhishekkr/gol/golenv"

	timedot "../timedot"
)

var (
	TIMEDOT_PORT = golenv.OverrideIfEnv("TIMEDOT_PORT", ":7999")
)

type Timedots struct {
	savedTimedots []*timedot.Record
}

func main() {
	conn, err := net.Listen("tcp", TIMEDOT_PORT)
	if err != nil {
		log.Fatalln("failed to bind at", TIMEDOT_PORT)
		return
	}

	log.Println("starting server... @", TIMEDOT_PORT)
	svr := grpc.NewServer()
	timedot.RegisterTimeFSServer(svr, &Timedots{})
	svr.Serve(conn)
}

func (tym *Timedots) CreateTimedot(c context.Context, input *timedot.Record) (*timedot.TimedotSave, error) {
	tym.savedTimedots = append(tym.savedTimedots, input)
	log.Println("CreateTimedot", tym.savedTimedots)
	return &timedot.TimedotSave{
		Success: true,
	}, nil
}

func (tym *Timedots) ReadTimedot(filtr *timedot.Record, stream timedot.TimeFS_ReadTimedotServer) error {
	log.Println("GetTimedot")
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

package timefsClient

import (
	"context"
	"fmt"
	"io"
	"log"

	"google.golang.org/grpc"

	timedot "github.com/abhishekkr/timefs/timedot"
)

func createTimeFS(client timedot.TimeFSClient, l *timedot.Record) {
	resp, err := client.CreateTimedot(context.Background(), l)

	if err != nil || !resp.Success {
		log.Printf("create timedot failed\nerr: %q\nresponse: %q", err.Error(), resp.Success)
	}
}

func getTimeFS(client timedot.TimeFSClient, filtr *timedot.Record) {
	stream, err := client.ReadTimedot(context.Background(), filtr)
	if err != nil {
		log.Println("error on get timedot: ", err.Error())
		return
	}

	for {
		l, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("%v.GetLogs(_) = _, %v", client, err)
		}
		fmt.Println("timedot: ", l)
	}
}

func pushSomeDummyTimeFS(client timedot.TimeFSClient) {
	tymdot := &timedot.Record{
		TopicKey: "appX",
		TopicId:  "x.cpu",
		Value:    "99",
		Time: []*timedot.Timedot{
			&timedot.Timedot{
				Year:        2017,
				Month:       3,
				Date:        17,
				Hour:        1,
				Minute:      20,
				Second:      30,
				Microsecond: 7,
			},
		},
	}
	createTimeFS(client, tymdot)
}

func GrpcClient(port string) {
	conn, err := grpc.Dial(port, grpc.WithInsecure())
	if err != nil {
		log.Fatalln("did not connect: ", err.Error())
		return
	}
	defer conn.Close()
	log.Println("starting client...")

	client := timedot.NewTimeFSClient(conn)
	pushSomeDummyTimeFS(client)

	log.Println("--- all")
	filterx := &timedot.Record{
		TopicKey: "appX",
		TopicId:  "x.cpu",
	}
	getTimeFS(client, filterx)
}

package timefsClient

import (
	"context"
	"fmt"
	"io"
	"log"

	"google.golang.org/grpc"

	timedot "github.com/abhishekkr/timefs/timedot"
)

func createTimeFS(client *timedot.TimeFSClient, l *timedot.Record) {
	resp, err := (*client).CreateTimedot(context.Background(), l)

	if err != nil {
		log.Printf("create timedot failed\nerr: %q\n",
			err.Error())
	} else if !resp.Success {
		log.Printf("create timedot failed\nresponse: %q", resp)
	}
}

func getTimeFS(client *timedot.TimeFSClient, filtr *timedot.Record) {
	stream, err := (*client).ReadTimedot(context.Background(), filtr)
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
			log.Printf("%v.ReadRecord(_) = _, %v", client, err)
		}
		fmt.Println("timedot: ", l)
	}
}

func LinkOpen(port string) *grpc.ClientConn {
	conn, err := grpc.Dial(port, grpc.WithInsecure())
	if err != nil {
		log.Fatalln("did not connect:", err.Error())
	}
	return conn
}

func LinkClose(conn *grpc.ClientConn) {
	conn.Close()
}

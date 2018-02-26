package timefsSplitter

import (
	"log"
	"strings"

	"google.golang.org/grpc"

	timefsClient "github.com/abhishekkr/timefs/client/timefsClient"
	timedot "github.com/abhishekkr/timefs/timedot"
)

var (
	Clients     []*timedot.TimeFSClient
	GrpcClients []*grpc.ClientConn
)

func ConnectBackends(backendCSV string) {
	links := strings.Split(backendCSV, ",")
	Clients = make([]*timedot.TimeFSClient, len(links))
	GrpcClients = make([]*grpc.ClientConn, len(links))

	for idx, link := range links {
		GrpcClients[idx] = timefsClient.LinkOpen(link)

		client := timedot.NewTimeFSClient(GrpcClients[idx])
		log.Println("connecting to backend:", link)
		Clients[idx] = &client
	}
}

func CloseBackends() {
	for _, client := range GrpcClients {
		timefsClient.LinkClose(client)
	}
}

func CreateTimeFS(record *timedot.Record) {
	timefsClient.CreateTimeFS(Clients[0], record)
}

func GetTimeFS(recordChan chan timedot.Record, filtr *timedot.Record) {
	recordChanX := make(chan timedot.Record)

	go timefsClient.GetTimeFS(Clients[0], filtr, recordChanX)

	for _record := range recordChanX {
		recordChan <- _record
	}
	close(recordChan)
}

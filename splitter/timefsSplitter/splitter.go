package timefsSplitter

import (
	"log"
	"strconv"
	"strings"

	golenv "github.com/abhishekkr/gol/golenv"
	grpc "google.golang.org/grpc"

	timefsClient "github.com/abhishekkr/timefs/client/timefsClient"
	timedot "github.com/abhishekkr/timefs/timedot"
)

var (
	clientByChannelCount = golenv.OverrideIfEnv("TIMEFS_CLIENTBYCHANNEL_COUNT", "100")
	timefsBackends       = golenv.OverrideIfEnv("TIMEFS_BACKENDS", "127.0.0.1:7999")

	clients     []*timedot.TimeFSClient
	grpcClients []*grpc.ClientConn
)

func init() {
	RegisterCreateEngines()

	if timefsBackends == "" {
		log.Println("no backends provided as of now, skipping connect backend and create-by-channel goroutines")
		return
	}

	ConnectBackends(timefsBackends)

	clientByChannelCountInt, err := strconv.Atoi(clientByChannelCount)
	if err != nil {
		log.Fatalln("failed to convert TIMEFS_CLIENTBYCHANNEL_COUNT value to int", err)
	}
	LaunchCreateByChannel(clientByChannelCountInt)
}

func ConnectBackends(backendCSV string) {
	links := strings.Split(backendCSV, ",")
	createRR.ClientsSize = uint64(len(links))
	clients = make([]*timedot.TimeFSClient, createRR.ClientsSize)
	grpcClients = make([]*grpc.ClientConn, createRR.ClientsSize)

	for idx, link := range links {
		grpcClients[idx] = timefsClient.LinkOpen(link)

		client := timedot.NewTimeFSClient(grpcClients[idx])
		log.Println("connecting to backend:", link)
		clients[idx] = &client
	}
}

func LaunchCreateByChannel(clientByChannelCountInt int) {
	createChannel = make(chan *timedot.Record)

	for idx := 0; idx < clientByChannelCountInt; idx++ {
		log.Println("creating background channel#", idx)
		go createTimeFSByChannel()
	}
}

func CloseBackends() {
	for _, client := range grpcClients {
		timefsClient.LinkClose(client)
	}
}

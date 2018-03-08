package timefsSplitter

import (
	"io"
	"sync"

	timefsClient "github.com/abhishekkr/timefs/client/timefsClient"
	timedot "github.com/abhishekkr/timefs/timedot"
)

func getTimeFS(client *timedot.TimeFSClient, recordChan chan timedot.Record, filtr *timedot.Record, wg *sync.WaitGroup) {
	defer wg.Done()
	timefsClient.GetTimeFS(client, filtr, recordChan)
}

func GetTimeFS(recordChan chan timedot.Record, filtr *timedot.Record, stream timedot.TimeFS_ReadTimedotServer) {
	var wg sync.WaitGroup

	wg.Add(len(clients))
	recordChanX := make(chan timedot.Record)
	for _, client := range clients {
		go getTimeFS(client, recordChanX, filtr, &wg)
	}

	var err error
	for record := range recordChanX {
		err = stream.Send(&record)
		if err == io.EOF {
			break
		}
		if err != nil {
			break
		}
	}
	close(recordChanX)

	wg.Wait()
}

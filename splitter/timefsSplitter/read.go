package timefsSplitter

import (
	"log"
	"sync"

	timefsClient "github.com/abhishekkr/timefs/client/timefsClient"
	timedot "github.com/abhishekkr/timefs/timedot"
)

func getTimeFS(client *timedot.TimeFSClient, recordChan chan timedot.Record, recordChanX chan timedot.Record, filtr *timedot.Record, wgGetTimeFS *sync.WaitGroup) {
	defer wgGetTimeFS.Done()

	go timefsClient.GetTimeFS(client, filtr, recordChanX)

	for _record := range recordChanX {
		recordChan <- _record
	}
	log.Println("[+] done with", client)
}

func GetTimeFS(recordChan chan timedot.Record, filtr *timedot.Record) {
	var wgGetTimeFS sync.WaitGroup

	wgGetTimeFS.Add(len(clients))
	for _, client := range clients {
		recordChanX := make(chan timedot.Record)
		go getTimeFS(client, recordChan, recordChanX, filtr, &wgGetTimeFS)
	}

	wgGetTimeFS.Wait()
	close(recordChan)
	log.Println("[+] stream ended for", filtr)
}

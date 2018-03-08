package timefsSplitter

import (
	"sync"

	timefsClient "github.com/abhishekkr/timefs/client/timefsClient"
	timedot "github.com/abhishekkr/timefs/timedot"
)

/*
create allows couple strategies

* CreateTimeFSByChannel
allows round-robin using tail-recursion~like~goroutine and global Channel
but it's just one goroutine at a time handling all Creates

* CreateTimeFS : allows round-robin
allows round-robin using global rotating backend counter
*/

type CreateEngine interface {
	CreateTimeFS(*timedot.Record)
}

type createRoundRobinByChannel struct {
}

type createRoundRobin struct {
	sync.RWMutex
	ClientIndex uint64
	ClientsSize uint64
}

var (
	CreateEngines  map[string]CreateEngine
	createRR       createRoundRobin
	createRRByChan createRoundRobinByChannel

	createChannel chan *timedot.Record
)

func RegisterCreateEngines() {
	CreateEngines = map[string]CreateEngine{
		"mutex":     &createRR,
		"recursion": &createRRByChan,
	}
}

func GetCreateEngine(name string) CreateEngine {
	return CreateEngines[name]
}

func (crr *createRoundRobin) Index() uint64 {
	crr.RLock()
	defer crr.RUnlock()
	if crr.ClientIndex >= crr.ClientsSize {
		crr.ClientIndex = 0
	}
	return crr.ClientIndex
}

func (crr *createRoundRobin) CreateTimeFS(record *timedot.Record) {
	timefsClient.CreateTimeFS(clients[createRR.Index()], record)
}

func createTimeFSByChannel() {
	/*
		doing it this way will keep rotating between clients
		without keeping rotating index explicitly
	*/
	var record *timedot.Record
	for _, client := range clients {
		record = <-createChannel
		timefsClient.CreateTimeFS(client, record)
	}
	go createTimeFSByChannel()
}

func (crr *createRoundRobinByChannel) CreateTimeFS(record *timedot.Record) {
	createChannel <- record
}

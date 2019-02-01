package fs

import (
	"github.com/abhishekkr/gol/golenv"
	timedot "github.com/abhishekkr/timefs/timedot"
)

/*
tips from:
	* https://stackimpact.com/blog/practical-golang-benchmarks/
*/

var (
	TIMEFS_DIR_ROOT = golenv.OverrideIfEnv("TIMEFS_DIR_ROOT", "/tmp/timefs")

	ValueFile = "value"
)

type StoreEngine interface {
	CreateRecord(record *timedot.Record)
	ReadRecords(recordChan chan timedot.Record, record *timedot.Record)
}

/*
StoreEngines acts as map for all available data storeX engines
*/
var StoreEngines = make(map[string]StoreEngine)

/*
RegisterStoreEngine gets used by adapters to register themselves.
*/
func RegisterStoreEngine(name string, store StoreEngine) {
	StoreEngines[name] = store
}

/*
GetStoreEngine gets used by client to fetch a required store.
*/
func GetStoreEngine(name string) StoreEngine {
	return StoreEngines[name]
}

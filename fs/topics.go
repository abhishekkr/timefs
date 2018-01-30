package fs

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"

	"github.com/abhishekkr/gol/golenv"
	timedot "github.com/abhishekkr/timefs/timedot"
)

/*
	tymdot := &timedot.Record{
		TopicKey: "appX",
		TopicId:  "x.cpu",
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
*/

var (
	TIMEFS_DIR_ROOT = golenv.OverrideIfEnv("TIMEFS_DIR_ROOT", "/tmp/timefs")
)

func int32ToStr(n int32) string {
	return strconv.Itoa(int(n))
}

func timedotPath(record *timedot.Record) string {
	return path.Join(TIMEFS_DIR_ROOT,
		record.TopicKey,
		record.TopicId,
		int32ToStr(record.Time[0].Year),
		int32ToStr(record.Time[0].Month),
		int32ToStr(record.Time[0].Date),
		int32ToStr(record.Time[0].Hour),
		int32ToStr(record.Time[0].Minute),
		int32ToStr(record.Time[0].Second),
		int32ToStr(record.Time[0].Microsecond),
	)
}

func timedotFile(dirname string) string {
	return path.Join(dirname, "value")
}

func CreateRecord(record *timedot.Record) {
	dirname := timedotPath(record)
	filepath := timedotFile(dirname)

	if _, err := os.Stat(filepath); !os.IsNotExist(err) {
		log.Println("[fs.CreateRecord] pre-existing", filepath)
		return
	}

	if _, err := os.Stat(dirname); os.IsNotExist(err) {
		err = os.MkdirAll(dirname, 0750)
		if err != nil {
			log.Println("[fs.CreateRecord] failed creating dir", dirname)
			return
		}
	}

	if err := ioutil.WriteFile(filepath, []byte(record.Value), 0644); err != nil {
		log.Println("[fs.CreateRecord] failed creating value file ", filepath)
	}
	return
}

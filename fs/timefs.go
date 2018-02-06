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

var (
	TIMEFS_DIR_ROOT = golenv.OverrideIfEnv("TIMEFS_DIR_ROOT", "/tmp/timefs")
	ValueFile       = "value"
)

func Int32ToStr(n int32) string {
	return strconv.Itoa(int(n))
}

func StrToInt32(s string) int32 {
	n, err := strconv.Atoi(s)
	if err != nil {
		log.Println("[fs.StrToInt32] failure to convert", s)
	}
	return int32(n)
}

func timedotPath(record *timedot.Record) string {
	return path.Join(TIMEFS_DIR_ROOT,
		record.TopicKey,
		record.TopicId,
		Int32ToStr(record.Time[0].Year),
		Int32ToStr(record.Time[0].Month),
		Int32ToStr(record.Time[0].Date),
		Int32ToStr(record.Time[0].Hour),
		Int32ToStr(record.Time[0].Minute),
		Int32ToStr(record.Time[0].Second),
		Int32ToStr(record.Time[0].Microsecond),
	)
}

func timedotFile(dirname string) string {
	return path.Join(dirname, ValueFile)
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

func ReadRecords(recordChan chan timedot.Record, record *timedot.Record) {
	var err error
	if len(record.Time) > 0 {
		return
	}

	record.Time = append(record.Time, &timedot.Timedot{})
	if len(record.Time) == 0 {
		return
	}

	err = globTimedots(recordChan, record)
	if err != nil {
		log.Println(err)
	}
	return
}

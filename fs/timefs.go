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

func int32ToStr(n int32) string {
	return strconv.Itoa(int(n))
}

func strToInt32(s string) int32 {
	n, err := strconv.Atoi(s)
	if err != nil {
		log.Println("[fs.strToInt32] failure to convert", s)
	}
	return int32(n)
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

func ReadRecords(record *timedot.Record) (records []*timedot.Record, err error) {
	if err = globRecordTopicKey(record); err != nil {
		return
	}
	if err = globRecordTopicId(record); err != nil {
		return
	}
	if err = globRecordTime(record); err != nil {
		return
	}
	return
}

func ReadRecord(record *timedot.Record) (value string) {
	dirname := timedotPath(record)
	filepath := timedotFile(dirname)
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		log.Println("[fs.ReadRecord] missing", filepath)
		return
	}

	if valueByte, err := ioutil.ReadFile(filepath); err == nil {
		value = string(valueByte)
		return
	}
	log.Println("[fs.ReadRecord] failed reading value file ", filepath)
	return
}

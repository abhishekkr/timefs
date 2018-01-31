package fs

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"

	timedot "github.com/abhishekkr/timefs/timedot"
)

func globRecordTopicKey(record *timedot.Record) error {
	if record.TopicKey != "" {
		return nil
	}
	topics, err := filepath.Glob(path.Join(TIMEFS_DIR_ROOT, "[-_a-Z0-9]"))
	if err != nil {
		return err
	}
	record.TopicKey = topics[0]
	return nil
}

func globRecordTopicId(record *timedot.Record) error {
	if record.TopicId != "" {
		return nil
	}

	dirname := path.Join(TIMEFS_DIR_ROOT, record.TopicKey, "[-_a-Z0-9]")
	ids, err := filepath.Glob(dirname)
	if err != nil {
		return err
	}
	record.TopicId = ids[0]
	return nil
}

func globRecordTime(record *timedot.Record) (err error) {
	if len(record.Time) > 0 {
		return
	}

	record.Time = append(record.Time, &timedot.Timedot{})
	if len(record.Time) == 0 {
		return
	}

	err = globTimedots(record)
	return
}

func globTimedots(record *timedot.Record) (err error) {
	dirname := path.Join(TIMEFS_DIR_ROOT,
		record.TopicKey,
		record.TopicId,
	)

	updatePathForTimedot(&dirname, record.Time[0].Year)
	updatePathForTimedot(&dirname, record.Time[0].Month)
	updatePathForTimedot(&dirname, record.Time[0].Date)
	updatePathForTimedot(&dirname, record.Time[0].Hour)
	updatePathForTimedot(&dirname, record.Time[0].Minute)
	updatePathForTimedot(&dirname, record.Time[0].Second)
	updatePathForTimedot(&dirname, record.Time[0].Microsecond)

	//bring in chan here to stream over grpc
	requiredRecords := []*timedot.Record{}
	err = filepath.Walk(dirname, func(_path string, f os.FileInfo, err error) error {
		if filepath.Base(_path) != ValueFile {
			return nil
		}
		_record := pathToTimedot(_path)
		requiredRecords = append(requiredRecords, _record)
		return nil
	})
	if err != nil {
		return
	}

	log.Println(requiredRecords)
	return
}

func updatePathForTimedot(dirname *string, dot int32) {
	if dot == 0 {
		return
	}
	*dirname = path.Join(*dirname, int32ToStr(dot))
}

func pathToTimedot(dotpath string) (record *timedot.Record) {
	dotvalue, err := ioutil.ReadFile(dotpath)
	if err != nil {
		log.Println(err)
		return
	}

	record = &timedot.Record{
		TopicKey: "-",
		TopicId:  "-",
		Value:    string(dotvalue),
		Time: []*timedot.Timedot{
			&timedot.Timedot{},
		},
	}

	return
}

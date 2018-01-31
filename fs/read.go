package fs

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	timedot "github.com/abhishekkr/timefs/timedot"
)

var (
	PathSlash = filepath.ToSlash("/")
)

func globTimedots(recordChan chan timedot.Record, record *timedot.Record) (err error) {
	dirname := TIMEFS_DIR_ROOT

	updatePathForTimedot(&dirname, record.TopicKey)
	updatePathForTimedot(&dirname, record.TopicId)
	updatePathForTimedot(&dirname, int32ToStr(record.Time[0].Year))
	updatePathForTimedot(&dirname, int32ToStr(record.Time[0].Month))
	updatePathForTimedot(&dirname, int32ToStr(record.Time[0].Date))
	updatePathForTimedot(&dirname, int32ToStr(record.Time[0].Hour))
	updatePathForTimedot(&dirname, int32ToStr(record.Time[0].Minute))
	updatePathForTimedot(&dirname, int32ToStr(record.Time[0].Second))
	updatePathForTimedot(&dirname, int32ToStr(record.Time[0].Microsecond))

	err = filepath.Walk(dirname, func(_path string, f os.FileInfo, err error) error {
		if filepath.Base(_path) != ValueFile {
			return nil
		}
		_record := pathToTimedot(_path)
		recordChan <- _record
		return nil
	})
	close(recordChan)

	if err != nil {
		return
	}

	return
}

func updatePathForTimedot(dirname *string, dot string) {
	if dot == "0" || dot != "" {
		return
	}
	*dirname = path.Join(*dirname, dot)
}

func pathToTimedot(dotpath string) (record timedot.Record) {
	relativeDotPath := strings.Split(dotpath, TIMEFS_DIR_ROOT)[1]
	dotpathSplit := strings.Split(relativeDotPath, PathSlash)
	if len(dotpathSplit) != 11 {
		log.Println("[fs.pathToTimedot] dotpath don't convert to dot", dotpath)
		return
	}

	dotvalue, err := ioutil.ReadFile(dotpath)
	if err != nil {
		log.Println(err)
		return
	}

	record = timedot.Record{
		TopicKey: dotpathSplit[1],
		TopicId:  dotpathSplit[2],
		Value:    string(dotvalue),
		Time: []*timedot.Timedot{
			&timedot.Timedot{
				Year:        strToInt32(dotpathSplit[3]),
				Month:       strToInt32(dotpathSplit[4]),
				Date:        strToInt32(dotpathSplit[5]),
				Hour:        strToInt32(dotpathSplit[6]),
				Minute:      strToInt32(dotpathSplit[7]),
				Second:      strToInt32(dotpathSplit[8]),
				Microsecond: strToInt32(dotpathSplit[9]),
			},
		},
	}

	return
}

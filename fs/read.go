package fs

import (
	"bufio"
	"io"
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

func readfile(filepath string) string {
	var data, _data string
	f, err := os.Open(filepath)
	if err != nil {
		log.Println(err)
		return ""
	}
	r := bufio.NewReader(f)

	data, err = r.ReadString('\n')
	for err == nil {
		_data, err = r.ReadString('\n')
		data += _data
	}
	if err != io.EOF {
		log.Println(err)
	}
	f.Close()
	return data
}

func globTimedots(recordChan chan timedot.Record, record *timedot.Record) (err error) {
	dirname := TIMEFS_DIR_ROOT

	updatePathForTimedot(&dirname, record.TopicKey)
	updatePathForTimedot(&dirname, record.TopicId)
	updatePathForTimedot(&dirname, Int32ToStr(record.Time[0].Year))
	updatePathForTimedot(&dirname, Int32ToStr(record.Time[0].Month))
	updatePathForTimedot(&dirname, Int32ToStr(record.Time[0].Date))
	updatePathForTimedot(&dirname, Int32ToStr(record.Time[0].Hour))
	updatePathForTimedot(&dirname, Int32ToStr(record.Time[0].Minute))
	updatePathForTimedot(&dirname, Int32ToStr(record.Time[0].Second))
	updatePathForTimedot(&dirname, Int32ToStr(record.Time[0].Microsecond))

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

	dotvalue := readfile(dotpath)

	record = timedot.Record{
		TopicKey: dotpathSplit[1],
		TopicId:  dotpathSplit[2],
		Value:    dotvalue,
		Time: []*timedot.Timedot{
			&timedot.Timedot{
				Year:        StrToInt32(dotpathSplit[3]),
				Month:       StrToInt32(dotpathSplit[4]),
				Date:        StrToInt32(dotpathSplit[5]),
				Hour:        StrToInt32(dotpathSplit[6]),
				Minute:      StrToInt32(dotpathSplit[7]),
				Second:      StrToInt32(dotpathSplit[8]),
				Microsecond: StrToInt32(dotpathSplit[9]),
			},
		},
	}

	return
}

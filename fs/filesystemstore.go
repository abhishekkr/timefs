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

type FilesystemStore struct {
}

/*
init registers filesystem store
*/
func init() {
	RegisterStoreEngine("filesystem", new(FilesystemStore))
}

func (store FilesystemStore) CreateRecord(record *timedot.Record) {
	dirname := store.timedotPath(record)
	filepath := store.timedotFile(dirname)

	if _, err := os.Stat(filepath); !os.IsNotExist(err) {
		log.Println("[fs.CreateRecord] pre-existing", filepath)
		return
	}

	store.writefile(dirname, filepath, record.Value)
	return
}

func (store FilesystemStore) ReadRecords(recordChan chan timedot.Record, record *timedot.Record) {
	var err error
	if len(record.Time) > 0 {
		return
	}

	record.Time = append(record.Time, &timedot.Timedot{})
	if len(record.Time) == 0 {
		return
	}

	err = store.globTimedots(recordChan, record)
	if err != nil {
		log.Println(err)
	}
	return
}

func (store FilesystemStore) timedotPath(record *timedot.Record) string {
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

func (store FilesystemStore) timedotFile(dirname string) string {
	return path.Join(dirname, ValueFile)
}

func (store FilesystemStore) writefile(dirname, filepath, data string) {
	var err error
	_, err = os.Stat(dirname)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dirname, 0750)
		if err != nil {
			log.Println("[fs.CreateRecord] failed creating dir", dirname)
			return
		}
	}

	f, err := os.Create(filepath)
	if err != nil {
		log.Println("[fs.CreateRecord] failed creating value file ", filepath)
		return
	}
	w := bufio.NewWriter(f)
	w.WriteString(data)
	w.Flush()
	f.Close()
}

func (store FilesystemStore) globTimedots(recordChan chan timedot.Record, record *timedot.Record) (err error) {
	dirname := TIMEFS_DIR_ROOT

	store.updatePathForTimedot(&dirname, record.TopicKey)
	store.updatePathForTimedot(&dirname, record.TopicId)
	store.updatePathForTimedot(&dirname, Int32ToStr(record.Time[0].Year))
	store.updatePathForTimedot(&dirname, Int32ToStr(record.Time[0].Month))
	store.updatePathForTimedot(&dirname, Int32ToStr(record.Time[0].Date))
	store.updatePathForTimedot(&dirname, Int32ToStr(record.Time[0].Hour))
	store.updatePathForTimedot(&dirname, Int32ToStr(record.Time[0].Minute))
	store.updatePathForTimedot(&dirname, Int32ToStr(record.Time[0].Second))
	store.updatePathForTimedot(&dirname, Int32ToStr(record.Time[0].Microsecond))

	err = filepath.Walk(dirname, func(_path string, f os.FileInfo, err error) error {
		if filepath.Base(_path) != ValueFile {
			return nil
		}
		_record := store.pathToTimedot(_path)
		recordChan <- _record
		return nil
	})
	close(recordChan)

	if err != nil {
		return
	}

	return
}

func (store FilesystemStore) updatePathForTimedot(dirname *string, dot string) {
	if dot == "0" || dot != "" {
		return
	}
	*dirname = path.Join(*dirname, dot)
}

func (store FilesystemStore) pathToTimedot(dotpath string) (record timedot.Record) {
	relativeDotPath := strings.Split(dotpath, TIMEFS_DIR_ROOT)[1]
	dotpathSplit := strings.Split(relativeDotPath, PathSlash)
	if len(dotpathSplit) != 11 {
		log.Println("[fs.pathToTimedot] dotpath don't convert to dot", dotpath)
		return
	}

	dotvalue := store.readfile(dotpath)

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

func (store FilesystemStore) readfile(filepath string) string {
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

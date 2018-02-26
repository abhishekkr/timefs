package timefsClient

import (
	"fmt"

	golenv "github.com/abhishekkr/gol/golenv"
	fs "github.com/abhishekkr/timefs/fs"
	timedot "github.com/abhishekkr/timefs/timedot"
)

type dummy struct {
	dot    timedot.Timedot
	client *timedot.TimeFSClient
}

var (
	DUMMY_TOPICKEY = golenv.OverrideIfEnv("DUMMY_TOPICKEY", "appX")
	DUMMY_TOPICID  = golenv.OverrideIfEnv("DUMMY_TOPICID", "x.cpu")
	DUMMY_VALUE    = golenv.OverrideIfEnv("DUMMY_VALUE", "99")
)

func (d *dummy) dummyRecord() *timedot.Timedot {
	return &(d.dot)
}

func (d *dummy) pushDummyRecord() {
	tymdot := &timedot.Record{
		TopicKey: DUMMY_TOPICKEY,
		TopicId:  DUMMY_TOPICID,
		Value:    DUMMY_VALUE,
		Time: []*timedot.Timedot{
			d.dummyRecord(),
		},
	}

	CreateTimeFS((*d).client, tymdot)
}

func (d *dummy) pushSomeDummyMicros() {
	ms := fs.StrToInt32(golenv.OverrideIfEnv("DUMMY_MS", "20"))
	for d.dot.Microsecond = int32(1); d.dot.Microsecond <= ms; d.dot.Microsecond++ {
		d.pushDummyRecord()
	}
}

func (d *dummy) pushSomeDummySec() {
	s := fs.StrToInt32(golenv.OverrideIfEnv("DUMMY_SEC", "40"))
	for d.dot.Second = int32(1); d.dot.Second <= s; d.dot.Second++ {
		d.pushSomeDummyMicros()
	}
}

func (d *dummy) pushSomeDummyMin() {
	m := fs.StrToInt32(golenv.OverrideIfEnv("DUMMY_MIN", "60"))
	for d.dot.Minute = int32(1); d.dot.Minute <= m; d.dot.Minute++ {
		d.pushSomeDummySec()
	}
}

func (d *dummy) pushSomeDummyHr() {
	h := fs.StrToInt32(golenv.OverrideIfEnv("DUMMY_HOUR", "12"))
	for d.dot.Minute = int32(1); d.dot.Minute <= h; d.dot.Minute++ {
		d.pushSomeDummyMin()
	}
}

func (d *dummy) pushSomeDummyTimeFS() {
	d.dot.Year = fs.StrToInt32(golenv.OverrideIfEnv("DUMMY_YEAR", "2018"))
	d.dot.Month = fs.StrToInt32(golenv.OverrideIfEnv("DUMMY_YEAR", "1"))
	d.dot.Date = fs.StrToInt32(golenv.OverrideIfEnv("DUMMY_YEAR", "30"))
	d.pushSomeDummyHr()
}

func DummyCreate(client *timedot.TimeFSClient) {
	d := dummy{client: client}
	d.pushSomeDummyTimeFS()
}

func DummyRead(client *timedot.TimeFSClient) {
	recordChan := make(chan timedot.Record)

	filterx := &timedot.Record{
		TopicKey: DUMMY_TOPICKEY,
		TopicId:  DUMMY_TOPICID,
	}

	go GetTimeFS(client, filterx, recordChan)

	for _record := range recordChan {
		fmt.Println("timedot: ", &_record)
	}
}

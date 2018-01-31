package timefsClient

import timedot "github.com/abhishekkr/timefs/timedot"

func dummyRecord(yy, mm, dd, hr, min, sec, ms int32) *timedot.Timedot {
	return &timedot.Timedot{
		Year:        yy,
		Month:       mm,
		Date:        dd,
		Hour:        hr,
		Minute:      min,
		Second:      sec,
		Microsecond: ms,
	}
}

func pushDummyRecord(client timedot.TimeFSClient, yy, mm, dd, hr, min, sec, ms int32) {
	tymdot := &timedot.Record{
		TopicKey: "appX",
		TopicId:  "x.cpu",
		Value:    "99",
		Time: []*timedot.Timedot{
			dummyRecord(yy, mm, dd, hr, min, sec, ms),
		},
	}
	createTimeFS(client, tymdot)
}

func pushSomeDummyTimeFS(client timedot.TimeFSClient) {
	yy := int32(2017)
	mm := int32(3)
	dd := int32(17)
	for hr := int32(1); hr <= 6; hr++ {
		for min := int32(1); min <= 6; min++ {
			for sec := int32(1); sec <= 6; sec++ {
				for ms := int32(1); ms <= 6; ms++ {
					pushDummyRecord(client, yy, mm, dd, hr, min, sec, ms)
				}
			}
		}
	}
}

package bitmex

import (
	"testing"
	"time"

	"github.com/go-openapi/strfmt"
	log "github.com/sirupsen/logrus"
)

type TestParam struct {
	Start int32
	Count int32
}

func NewTestParam() DownParam {
	return new(TestParam)
}

func (p *TestParam) SetStart(start *int32) {
	p.Start = *start
}

func (p *TestParam) SetCount(count *int32) {
	p.Count = *count
}

func (p *TestParam) SetStartTime(startTime *strfmt.DateTime) {
}

func (p *TestParam) SetEndTime(endTime *strfmt.DateTime) {
}

type TestDataCenter struct {
	Datas []interface{}
	// Offset      int
	// offsetMutex sync.Mutex
	// Once        int
}

func (d *TestDataCenter) SampleDownImpl(param DownParam) (data []interface{}, err error) {
	nLen := len(d.Datas)
	var end int
	tp := param.(*TestParam)

	if int(tp.Start) >= nLen {
		return
	}
	end = int(tp.Start + tp.Count)
	if end >= nLen {
		end = nLen
	}
	data = make([]interface{}, end-int(tp.Start))
	copy(data, d.Datas[tp.Start:end])
	time.Sleep(time.Second)
	return
}

func TestDownload(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	d := new(TestDataCenter)
	nTotal := 1024
	for i := 0; i != nTotal; i++ {
		d.Datas = append(d.Datas, i)
	}
	var once int32
	once = 13
	tmStart := strfmt.DateTime(time.Now().Add(0 - time.Hour))
	tmEnd := strfmt.DateTime(time.Now())
	down := NewDataDownload(tmStart, tmEnd, NewTestParam, d.SampleDownImpl, int32(once), 5)
	dataChan := down.Start()
	nCount := 0
	for d := range dataChan {
		log.Println(d)
		nCount += len(d)
	}
	if nCount != nTotal {
		log.Fatal("count not match")
	}
}

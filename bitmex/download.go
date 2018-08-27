package bitmex

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-openapi/strfmt"
	log "github.com/sirupsen/logrus"
)

type DownParam interface {
	SetStart(start *int32)
	SetCount(count *int32)
	SetStartTime(startTime *strfmt.DateTime)
	SetEndTime(endTime *strfmt.DateTime)
}

type NewParamFunc func() DownParam

type DownFunc func(DownParam) ([]interface{}, error)

type DataDownload struct {
	nRoutine  int
	dataCh    chan []interface{}
	paramFunc NewParamFunc
	downFunc  DownFunc
	onceCount int32
	startTime strfmt.DateTime
	endTime   strfmt.DateTime
	nFinish   int32 // 0 not finish, 1 finish
	start     chan int
	wg        sync.WaitGroup
}

func NewDataDownload(start, end strfmt.DateTime, paramFunc NewParamFunc, downFunc DownFunc, onceCount int32, nRoutine int) (d *DataDownload) {
	d = new(DataDownload)
	d.paramFunc = paramFunc
	d.nRoutine = nRoutine
	d.dataCh = make(chan []interface{}, 1024)
	d.downFunc = downFunc
	d.onceCount = onceCount
	d.startTime = start
	d.endTime = end
	d.start = make(chan int, nRoutine)
	return
}

func (d *DataDownload) Start() (dataCh chan []interface{}) {
	go d.Run()
	dataCh = d.dataCh
	return
}
func (d *DataDownload) Run() {
	d.wg.Add(d.nRoutine)
	for i := 0; i != d.nRoutine; i++ {
		go d.routine(i)
	}
	nStart := 0
	for {
		if d.IsFinish() {
			close(d.start)
			break
		}
		d.start <- nStart
		nStart += int(d.onceCount)
	}
	d.wg.Wait()
	close(d.dataCh)
	log.Info("DataDownload finished...")
	return
}

func (d *DataDownload) routine(nIndex int) {
	defer d.wg.Done()
	var err error
	var bFinish bool
Outer:
	for {
		select {
		case nStart, ok := <-d.start:
			if !ok {
				break Outer
			}
			bFinish, err = d.RunOnce(nStart)
			if err != nil {
				log.Error("routine Runonce error:", nStart, err.Error())
				break Outer
			}
			if bFinish {
				break Outer
			}
		default:
			if d.IsFinish() {
				break Outer
			}
		}
		time.Sleep(time.Microsecond)
	}
	log.Info("routine finish:", nIndex)
}

func (d *DataDownload) RunOnce(start int) (bFinish bool, err error) {
	log.Debug("RunOnce:", start)
	nStart := int32(start)
	params := d.paramFunc()
	params.SetStartTime(&d.startTime)
	params.SetStartTime(&d.endTime)
	params.SetCount(&d.onceCount)
	params.SetStart(&nStart)
	ret, err := d.downFunc(params)
	if err != nil {
		return
	}
	if len(ret) > 0 {
		d.dataCh <- ret
	}
	if len(ret) < int(d.onceCount) {
		d.SetFinish(1)
		bFinish = true
	}
	return
}

func (d *DataDownload) SetFinish(nFinish int32) {
	atomic.StoreInt32(&d.nFinish, nFinish)
}

func (d *DataDownload) IsFinish() (bFinish bool) {
	nFinish := atomic.LoadInt32(&d.nFinish)
	bFinish = nFinish == 1
	return
}

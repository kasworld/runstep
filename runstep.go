// goroutine stepwise exec package
// embed in struct to make
package runstep

import (
	"time"

	"github.com/kasworld/runstate"
)

type RunStep struct {
	*runstate.RunState
	startStep chan interface{}
	resultCh  chan interface{}
}

func New(bufcount int) *RunStep {
	return &RunStep{
		runstate.New(),
		make(chan interface{}, bufcount),
		make(chan interface{}, bufcount),
	}
}

// change for shared result ch
func (fs *RunStep) SetResultCh(ch chan interface{}) {
	fs.resultCh = ch
}

// for externel code
func (fs *RunStep) StartStepCh() chan<- interface{} {
	return fs.startStep
}
func (fs *RunStep) ResultCh() <-chan interface{} {
	return fs.resultCh
}
func (fs *RunStep) Run(stepfn func(d interface{}) interface{}) {
	for stepdata := range fs.startStep {
		if !fs.CanRun() {
			fs.TryStop()
			break
		}
		fs.resultCh <- stepfn(stepdata)
	}
	fs.SetBit(1)
}

func (fs *RunStep) Stop() {
	fs.TryStop()
	time.Sleep(0)
	for !fs.IsStopped() {
		select {
		case <-fs.resultCh:
		default:
		}
		select {
		case fs.startStep <- nil:
		default:
		}
	}
}

func (fs *RunStep) IsStopped() bool {
	return fs.GetBit(1)
}

// for embeding struct method
// use when custom Run method
// not included in runstepi

func (fs *RunStep) RecvStepArg() interface{} {
	return <-fs.startStep
}
func (fs *RunStep) SendStepResult(d interface{}) {
	fs.resultCh <- d
}

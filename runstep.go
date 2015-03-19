// goroutine stepwise exec package
// embed in struct to make
package runstep

import (
	"time"

	"github.com/kasworld/runstate"
)

type RunStep struct {
	startStep chan interface{}
	resultCh  chan interface{}
	rs        *runstate.RunState
}

func New(bufcount int) *RunStep {
	return &RunStep{
		make(chan interface{}, bufcount),
		make(chan interface{}, bufcount),
		runstate.New(),
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
		if !fs.rs.CanRun() {
			fs.rs.TryStop()
			break
		}
		fs.resultCh <- stepfn(stepdata)
	}
	fs.rs.SetBit(1)
}

func (fs *RunStep) Quit() {
	fs.rs.TryStop()
	time.Sleep(0)
	for !fs.rs.GetBit(1) {
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

func (fs *RunStep) IsQuit() bool {
	return fs.rs.GetBit(1)
}

// for embeding struct method
// use when custom Run method
func (fs *RunStep) RecvStepArg() interface{} {
	return <-fs.startStep
}
func (fs *RunStep) SendStepResult(d interface{}) {
	fs.resultCh <- d
}

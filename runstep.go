// goroutine stepwise exec package
// embed in struct to make
package runstep

import (
// "time"
)

type RunStep struct {
	Dummyint  int // for disk save
	startStep chan interface{}
	resultCh  chan interface{}
	ended     bool
}

func New(bufcount int) *RunStep {
	return &RunStep{
		0,
		make(chan interface{}, bufcount),
		make(chan interface{}, bufcount),
		false,
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
		fs.resultCh <- stepfn(stepdata)
	}
	fs.ended = true
}

func (fs *RunStep) Quit() {
	close(fs.startStep)
	if len(fs.resultCh) > 1 { // if shared ch
		return
	}
	for !fs.ended {
		select {
		case <-fs.resultCh:
		default:
		}
	}
}

// for embeding struct method
// use when custom Run method
func (fs *RunStep) RecvStepArg() interface{} {
	return <-fs.startStep
}
func (fs *RunStep) SendStepResult(d interface{}) {
	fs.resultCh <- d
}

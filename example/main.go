package main

import (
	"github.com/kasworld/runstep"
	"time"
)

//start
type StepRunObj struct {
	*runstep.RunStep
	Mode int
}

func NewStepRunObj() *StepRunObj {
	aib := &StepRunObj{}
	aib.RunStep = runstep.New(0)
	return aib
}
func (aib *StepRunObj) Step(datain interface{}) interface{} {
	// do step work
	println("step")
	time.Sleep(1 * time.Second)
	return 0
}

func main() {
	// init objs
	objs := [10]*StepRunObj{}
	for i, _ := range objs {
		objs[i] = NewStepRunObj()
		go objs[i].Run(objs[i].Step)
	}
	// run step
	for _, v := range objs {
		v.StartStepCh() <- 0
	}
	// do other work
	time.Sleep(3 * time.Second)
	// confirm all end
	for _, v := range objs {
		<-v.ResultCh()
	}
}

//end

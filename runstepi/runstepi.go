// interface of runstep package
package runstepi

type RunStepI interface {
	StartStepCh() chan<- interface{}
	ResultCh() <-chan interface{}
	Run(stepfn func(d interface{}) interface{})
	Quit()
}

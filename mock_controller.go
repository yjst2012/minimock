package minimock

import (
	"sync"
	"time"
)

//Mocker describes common interface for all mocks generated by minimock
type Mocker interface {
	MinimockFinish()
	MinimockWait(time.Duration)
}

// Tester contains subset of the testing.T methods used by the generated code
type Tester interface {
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Error(...interface{})
}

//MockController can be passed to mocks generated by minimock
type MockController interface {
	Tester

	RegisterMocker(Mocker)
}

//Controller implements MockController interface and has to be used in your tests:
//mockController := minimock.NewController(t)
//defer mockController.Finish()
//stringerMock := NewStringerMock(mockController)
type Controller struct {
	Tester
	sync.Mutex

	mockers []Mocker
}

//NewController returns an instance of Controller
func NewController(t Tester) Controller {
	return Controller{Tester: t}
}

//RegisterMocker puts mocker to the list of controller mockers
func (c *Controller) RegisterMocker(m Mocker) {
	c.Lock()
	c.mockers = append(c.mockers, m)
	c.Unlock()
}

//Finish calls to MinimockFinish method for all registered mockers
func (c *Controller) Finish() {
	c.Lock()
	for _, m := range c.mockers {
		m.MinimockFinish()
	}
	c.Unlock()
}

//Wait calls to MinimockWait method for all registered mockers
func (c *Controller) Wait(d time.Duration) {
	wg := sync.WaitGroup{}
	wg.Add(len(c.mockers))
	for _, m := range c.mockers {
		go func(m Mocker) {
			defer wg.Done()
			m.MinimockWait(d)
		}(m)
	}

	wg.Wait()
}

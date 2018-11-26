package checker

import (
	"github.com/localghost/healthy/utils"
	"github.com/spf13/viper"
	"time"
)

type Checker struct {
	tasks     map[string]*Task
	metrics   map[string]error
	request   chan string
	responses map[string]chan error
	interval  time.Duration
}

type Spec struct {
	Type string
	Interval time.Duration
}

type Task struct {
	spec  Spec
	check Check
}

type metric struct {
	name string
	value error
}

func New(checks interface{}) (*Checker, error) {
	result := &Checker{
		tasks:     make(map[string]*Task),
		metrics:   make(map[string]error),
		request:   make(chan string),
		responses: make(map[string]chan error),
		interval:  viper.GetDuration("checker.interval"),
	}
	if err := result.parseChecks(checks); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Checker) Start() {
	c.startChecks()
}

func (c *Checker) parseChecks(checks interface{}) error {
	var specs = make(map[string]Spec)
	if err := utils.Decode(checks, &specs); err != nil {
		return err
	}

	for name, check := range checks.(map[string]interface{}) {
		options := check.(map[string]interface{})
		if check, err := registry.CreateAndConfigure(specs[name].Type, options); err != nil {
			return err
		} else {
			c.tasks[name] = &Task{
				spec:  specs[name],
				check: check,
			}
			c.responses[name] = make(chan error)
		}
	}
	return nil
}

func (c *Checker) startChecks() {
	receiver := make(chan metric)
	for name, tasks := range c.tasks {
		var interval time.Duration
		if tasks.spec.Interval != time.Duration(0) {
			interval = tasks.spec.Interval
		} else {
			interval = c.interval
		}
		go func(name string, check Check, interval time.Duration) {
			for {
				select {
				case <- time.After(interval):
					receiver <- metric{name, check.Run()}
				}
			}
		}(name, tasks.check, interval)
	}
	go func() {
		for {
			select {
			case m := <-receiver:
				c.metrics[m.name] = m.value
			case name := <- c.request:
				err, ok := c.metrics[name]
				if !ok {
					c.responses[name] <- utils.NewNoSuchCheckError(name)
				} else {
					c.responses[name] <- err
				}
			}
		}
	}()
}

func (c *Checker) Get(name string) error {
	c.request <- name
	return <-c.responses[name]
}

func (c* Checker) GetAll() error {
	for name := range c.tasks {
		if err := c.Get(name); err != nil {
			return err
		}
	}
	return nil
}

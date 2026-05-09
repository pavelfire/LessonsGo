package pkg

import (
	"errors"
	"fmt"
)

var ErrDoSomething = errors.New("doSomething error I want to catch")

func DoSomething() error {
	err := DoSomethingElse()
	if err != nil {
		return fmt.Errorf("doSometingElse failed: %w", err)
	}
	return nil
}

var ErrDoSomethingElse = errors.New("doSomethingElse error I want to catch")

func DoSomethingElse() error {
	//return errors.New("ooops")
	return ErrDoSomethingElse
}

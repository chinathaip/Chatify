package error

import (
	"fmt"
	"log"
)

type HubError struct {
	Err error
}

func (e *HubError) Error() string {
	return fmt.Sprintf("Error In Hub.go: %s", e.Err.Error())
}

func (e *HubError) Log(err error) {
	e.Err = err
	log.Println(e.Error())
}

type HandlerError struct {
	Err error
}

func (e *HandlerError) Error() string {
	return fmt.Sprintf("Error In handler.go: %s", e.Err.Error())
}

func (e *HandlerError) Log(err error) {
	e.Err = err
	log.Println(e.Error())
}

type ChatError struct {
	Err error
}

func (e *ChatError) Error() string {
	return fmt.Sprintf("Error In chat.go: %s", e.Err.Error())
}

func (e *ChatError) Log(err error) {
	e.Err = err
	log.Println(e.Error())
}

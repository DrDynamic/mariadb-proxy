package errors

import (
	"log"
)

func HandleErrors(errors chan error, quit chan bool) {
	for {
		select {
		case <-quit:
			return
		default:
			WriteError(<-errors)
		}
	}
}

func WriteError(err error) {
	log.Fatal("Error: ", err)
}

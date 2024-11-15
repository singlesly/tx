package p2p

import (
	"errors"
	"fmt"
)

func CheckMessageSize(msg []byte) error {
	actual := len(msg)
	maximum := 500000
	if actual > maximum {
		return errors.New(fmt.Sprintf("Message size is %d but max can be %d", actual, maximum))
	}

	return nil
}

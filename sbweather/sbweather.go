package sbweather

import (
	"errors"
	"strconv"
)

type routine struct {
	err error
}

func New(zip string) *routine {
	var r routine

	if len(zip) != 5 {
		r.err = errors.New("Invalid Zip Code length")
		return &r
	}

	_, err := strconv.Atoi(zip)
	if err != nil {
		r.err = err
		return &r
	}

	return &r
}

func (r *routine) Update() {
}

func (r *routine) String() string {
	return "weather"
}

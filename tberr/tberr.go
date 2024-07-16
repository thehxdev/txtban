package tberr

import (
	"encoding/json"
)

type TbError struct {
	Err  string `json:"error"`
	Help string `json:"help,omitempty"`
}

func (t *TbError) Error() string {
	bytes, err := json.Marshal(t)
	if err != nil {
		return err.Error()
	}
	return string(bytes)
}

func New(msg string, help ...string) *TbError {
	var helpMsg string
	if len(help) > 0 {
		helpMsg = help[0]
	}

	return &TbError{
		Err:  msg,
		Help: helpMsg,
	}
}

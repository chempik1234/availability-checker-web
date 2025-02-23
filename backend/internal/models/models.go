package models

import "time"

type LogRecord struct {
	Name     string    `json:"name"`
	Result   bool      `json:"result"`
	Datetime time.Time `json:"datetime"`
}

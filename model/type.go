package models

import (
	"fmt"
	"time"
)

type ConfDB struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
	MaxIdle  int
	MaxOpen  int
	MaxLive  int
}

type TimeNormal struct {
	time.Time
}

func (t TimeNormal) MarshalJSON() ([]byte, error) {
	tune := t.Format(`"2006-01-02 15:04:05"`)
	return []byte(tune), nil
}

// Scan valueof time.Time
func (t *TimeNormal) Scan(v any) error {
	value, ok := v.(time.Time) // NOT directly assertion v.(TimeNormal)
	if ok {
		*t = TimeNormal{
			Time: value,
		}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

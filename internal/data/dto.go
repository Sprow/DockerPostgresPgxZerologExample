package data

import (
	"errors"
	"time"
)

type Obj struct {
	ID        int       `json:"id,omitempty"`
	Data1     string    `json:"data1"`
	Data2     string    `json:"data2"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

func (obj *Obj) IsValid() error {
	if len([]rune(obj.Data1)) < 5 {
		return errors.New("BadRequest Data1")
	}
	if len(obj.Data2) < 1 || len(obj.Data2) > 50 {
		return errors.New("BadRequest Data2")
	}
	return nil
}

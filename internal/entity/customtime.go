package entity

import (
	"bytes"
	"fmt"
	"time"
)

const (
	defaultTimeLayout = "2006-01-02T15:04:05"
)

type CustomTime struct {
	time.Time
}

func (t CustomTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, t.Time.Format(defaultTimeLayout))), nil
}

func (t *CustomTime) UnmarshalJSON(data []byte) error {
	// from json doc: by convention, unmarshalers implement UnmarshalJSON([]byte("null")) as a no-op
	if bytes.Equal(data, []byte("null")) {
		return nil
	}

	time, err := time.Parse(`"`+defaultTimeLayout+`"`, string(data))
	if err != nil {
		return err
	}

	t.Time = time
	return nil
}

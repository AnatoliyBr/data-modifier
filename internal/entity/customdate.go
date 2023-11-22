package entity

import (
	"bytes"
	"fmt"
	"time"
)

const (
	defaultDateLayout = "2006-01-02"
)

type CustomDate struct {
	time.Time
}

func (t CustomDate) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, t.Time.Format(defaultDateLayout))), nil
}

func (t *CustomDate) UnmarshalJSON(data []byte) error {
	// from json doc: by convention, unmarshalers implement UnmarshalJSON([]byte("null")) as a no-op
	if bytes.Equal(data, []byte("null")) {
		return nil
	}

	time, err := time.Parse(`"`+defaultDateLayout+`"`, string(data))
	if err != nil {
		return err
	}

	t.Time = time
	return nil
}

package entity

import (
	"bytes"
	"fmt"
	"time"
)

const (
	defaultDateLayout = "2006-01-02"
)

// CustomDate contains date in a custom format.
type CustomDate struct {
	time.Time
}

// MarshalJSON is a redefined method for correct json.Encoder work.
func (t CustomDate) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, t.Time.Format(defaultDateLayout))), nil
}

// UnmarshalJSON is a redefined method for correct json.Decoder work.
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

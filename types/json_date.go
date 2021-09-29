package types

import (
	"encoding/json"
	"strings"
	"time"
)

// JsonDate is a date in YYYY-MM-dd format
// swagger:model
// swagger:strfmt date
type JsonDate time.Time

func (j *JsonDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	*j = JsonDate(t)
	return nil
}

func (j JsonDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(j)
}

// Format for printing your date
func (j JsonDate) Format(layout string) string {
	t := time.Time(j)
	return t.Format(layout)
}

func (j JsonDate) String() string {
	return j.Format("2006-01-02")
}

func (j JsonDate) Time() time.Time {
	return time.Time(j)
}

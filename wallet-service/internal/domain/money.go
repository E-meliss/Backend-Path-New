package domain

import (
	"encoding/json"
	"fmt"
)

type Money int64

func (m Money) IsPositive() bool { return m > 0 }

func (m Money) MarshalJSON() ([]byte, error) {
	sign := ""
	v := m
	if v < 0 {
		sign = "-"
		v = -v
	}
	major := v / 100
	minor := v % 100
	s := fmt.Sprintf(`"%s%d.%02d"`, sign, major, minor)
	return []byte(s), nil
}

func (m *Money) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	var sign int64 = 1
	if len(s) > 0 && s[0] == '-' {
		sign = -1
		s = s[1:]
	}
	var major, minor int64
	_, err := fmt.Sscanf(s, "%d.%02d", &major, &minor)
	if err != nil {
		return err
	}
	*m = Money(sign * (major*100 + minor))
	return nil
}

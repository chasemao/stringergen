package example

import (
	"encoding/json"
)

// String Used in fmt to generate string
func (s *someStruct) String() string {
	v, _ := json.Marshal(s)
	return string(v)
}

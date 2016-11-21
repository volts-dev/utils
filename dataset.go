package utils

import (
	"time"
)

type (
	TDataSet struct {
		name string
		data map[string]string // TODO []byte
	}
)

func (self *TDataSet) AsString(name string, dlt ...string) string {
	return self.data[name]
}

func (self *TDataSet) AsInteger(name string) int64 {
	return StrToInt64(self.data[name])
}

func (self *TDataSet) AsBoolean(name string) bool {
	return StrToBool(self.data[name])
}

func (self *TDataSet) AsDateTime(name string) (t time.Time) {
	t, _ = time.Parse(time.RFC3339, self.data[name])
	return
}

func (self *TDataSet) AsFloat(name string) float64 {
	return StrToFloat(self.data[name])
}

package shutterfly

import (
	"testing"
)

const (
	APP_ID = "179e2c6d91b95bbc8abf2673c181ee70"
)

type Shutterfly struct {
	AuthToken string
	Test      *testing.T
}

func (self *Shutterfly) TestLog(msg string) {
	if self.Test != nil {
		self.Test.Log(msg)
	}
}

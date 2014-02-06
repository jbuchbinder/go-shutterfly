package shutterfly

import (
	"testing"
)

const (
	APP_ID = "179e2c6d91b95bbc8abf2673c181ee70"
)

// Shutterfly is the base object for the go-shutterfly library.
type Shutterfly struct {
	AuthToken string
	Test      *testing.T
}

// TestLog prints a logging statement, if a testing.T object has been
// passed to the Shutterfly object. If not, it does nothing.
func (self *Shutterfly) TestLog(msg string) {
	if self.Test != nil {
		self.Test.Log(msg)
	}
}

package connect_proxy_scheme_test

import (
	"testing"

	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type DialerSuite struct{}

var _ = Suite(&DialerSuite{})

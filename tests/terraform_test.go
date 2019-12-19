package testing

import (
	"testing"
	"time"
)

func TestNow(t *testing.T) {
	t.Logf("Now: %s\n", time.Now())
}

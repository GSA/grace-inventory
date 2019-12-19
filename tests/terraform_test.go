package testing

import (
    "time"
)

func TestNow(t *testing.T) {
    t.Logf("Now: %s\n", time.Now())
}
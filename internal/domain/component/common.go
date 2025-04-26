package component

import (
	"time"

	"github.com/benbjohnson/clock"
)

func timeNowUTC(c clock.Clock) time.Time {
	return c.Now().UTC()
}

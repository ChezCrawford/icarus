package utils

import "time"

func ParseUnixSeconds(unixSeconds int64) time.Time {
	return time.Unix(unixSeconds, 0)
}

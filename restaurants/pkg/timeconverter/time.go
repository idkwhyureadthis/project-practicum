package timeconverter

import "time"

func TimeToMicro(t time.Time) int64 {
	h, m, s := t.Hour(), t.Minute(), t.Second()
	total_seconds := int64(s) + int64(60*m) + int64(3600*h)
	return total_seconds * 1_000_000
}

package limiter

import "time"

type RequestData struct {
	Endpoint string
	IP       string
	Time     time.Time
}

type Limiter struct {
}

func (l *Limiter) Start() {}

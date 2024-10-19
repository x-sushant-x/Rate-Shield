package models

import "time"

type ErrorNotification struct {
	Error     string    `json:"error"`
	Timestamp time.Time `json:"timestamp"`
	IP        string    `json:"ip"`
	Endpoint  string    `json:"endpoint"`
}

func (e *ErrorNotification) CreateErrorNotificationObject(err string, timestamp time.Time, ip, endpoint string) ErrorNotification {
	return ErrorNotification{
		Error:     err,
		Timestamp: timestamp,
		IP:        ip,
		Endpoint:  endpoint,
	}
}

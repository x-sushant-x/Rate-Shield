package service

import (
	"time"

	"github.com/x-sushant-x/RateShield/models"
	"github.com/x-sushant-x/RateShield/utils"
)

type ErrorNotificationSVC struct {
	slackSVC            SlackService
	notificationHistory map[string]time.Time
}

func NewErrorNotificationSVC(slackService SlackService) ErrorNotificationSVC {
	return ErrorNotificationSVC{
		slackSVC:            slackService,
		notificationHistory: make(map[string]time.Time),
	}
}

func (e *ErrorNotificationSVC) SendErrorNotification(systemError string, timestamp time.Time, ip string, endpoint string) {
	if !e.CanSendNotification(ip, endpoint) {
		return
	}

	notification := models.ErrorNotification{}
	notification = notification.CreateErrorNotificationObject(systemError, timestamp, ip, endpoint)

	e.SendNotification(notification)
	e.notificationHistory[ip+":"+endpoint] = time.Now()

}

func (e *ErrorNotificationSVC) CanSendNotification(ip, endpoint string) bool {
	key := ip + ":" + endpoint

	lastNotifiedTime, ok := e.notificationHistory[key]
	if !ok {
		return true
	}

	sinceTime := time.Since(lastNotifiedTime)
	return sinceTime.Minutes() >= 5
}

func (e *ErrorNotificationSVC) SendNotification(notification models.ErrorNotification) {
	b, _ := utils.MarshalJSON(notification)
	e.slackSVC.SendSlackMessage(string(b))
}

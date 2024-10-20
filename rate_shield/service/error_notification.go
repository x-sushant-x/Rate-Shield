package service

import (
	"fmt"
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

func (e *ErrorNotificationSVC) SendErrorNotification(systemError string, timestamp time.Time, ip string, endpoint string, rule models.Rule) {
	if !e.canSendNotification(ip, endpoint) {
		return
	}

	ruleString, _ := utils.MarshalJSON(rule)

	notificationString := fmt.Sprintf("Error: %s,\n IP: %s,\n Endpoint: %s,\n Rule: %s,\n Timestamp: %s", systemError, ip, endpoint, ruleString, timestamp)

	e.sendNotification(notificationString)
	e.notificationHistory[ip+":"+endpoint] = time.Now()

}

func (e *ErrorNotificationSVC) canSendNotification(ip, endpoint string) bool {
	key := ip + ":" + endpoint

	lastNotifiedTime, ok := e.notificationHistory[key]
	if !ok {
		return true
	}

	sinceTime := time.Since(lastNotifiedTime)
	return sinceTime.Seconds() >= 30
}

func (e *ErrorNotificationSVC) sendNotification(notification string) {
	e.slackSVC.SendSlackMessage(notification)
}

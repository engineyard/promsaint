package models

import (
	"time"

	prometheus "github.com/prometheus/common/model"
)

// Alert represents a Nagios alert that is to be converted to a Prometheus alert
//
// Type: host or service or something custom
// Host: hostname
// Service: servicename
// AlertName: custom alert name, if Type is neither host or service
// Notify: notify string ** overwriten by notifyLabel if present **
// Notification types:
//   PROBLEM / ACKNOWLEDGEMENT / RECOVERY
// State:
//   Host states:
//     UP / DOWN
//   Service states:
//     CRITICAL / WARNING / UNKNOWN / OK
// Message: Optional message
// Note: Reference URL
type Alert struct {
	Type             string `json:"type"`
	Host             string `json:"host"`
	Service          string `json:"service"`
	AlertName        string `json:"alert-name"`
	Notify           string `json:"notify"`
	NotificationType string `json:"notification-type"`
	State            string `json:"state"`
	Message          string `json:"message"`
	Note             string `json:"note"`
	FirePeriod       string `json:"fire-period"`
}

type AlertMetadata struct {
	LastUpdate      time.Time
	DontForgetUntil time.Time
}

type InternalAlert struct {
	PrometheusAlert prometheus.Alert
	Metadata        AlertMetadata
}

// A function type sent to export function
type NotificationSender interface {
	Send([]prometheus.Alert)
}

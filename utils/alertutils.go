package utils

import (
	"crypto/sha1"
	"flag"
	"net/url"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/engineyard/promsaint/models"
	prometheus "github.com/prometheus/common/model"
)

var (
	generatorUrl  = flag.String("generator_url", "https://nagios.example.com/nagios/cgi-bin/status.cgi", "Breakdown Source for nagios")
	customerLabel = flag.String("customer-label", "", "The customer label for the forwarded alerts")
)

func Key(alert *models.Alert) string {
	sum := sha1.Sum(append([]byte(alert.Service), []byte(alert.Host)...))
	str := string(sum[:])
	log.Debugf("AlertKey: %s", str)
	return str
}

/* https://github.com/prometheus/prometheus/blob/master/vendor/github.com/prometheus/common/model/alert.go#L29
 * // Alert is a generic representation of an alert in the Prometheus eco-system.
 * type Alert struct {
 * 	// Label value pairs for purpose of aggregation, matching, and disposition
 * 	// dispatching. This must minimally include an "alertname" label.
 * 	Labels LabelSet `json:"labels"`
 *
 * 	// Extra key/value information which does not define alert identity.
 * 	Annotations LabelSet `json:"annotations"`
 *
 * 	// The known time range for this alert. Both ends are optional.
 * 	StartsAt     time.Time `json:"startsAt,omitempty"`
 * 	EndsAt       time.Time `json:"endsAt,omitempty"`
 * 	GeneratorURL string    `json:"generatorURL"`
 * }
 *
 * type LabelSet map[LabelName]LabelValue
 * type LabelName string
 * type LabelValue string
 */

// Merge a new nagios alert into an existing prometheus alert (or an empty prometheus struct if the alert doesn't already exist)
func Merge(pAlert *models.InternalAlert, alert *models.Alert) {
	var alertname string
	if alert.Type == "host" {
		alertname = "Host Down"
	} else if alert.Type == "service" {
		alertname = alert.Service
	} else {
		alertname = alert.AlertName
	}

	log.Debugf("NOTIFY: %s -> %s", string(pAlert.PrometheusAlert.Labels["notify"]), alert.Notify)
	notifyMap := map[string]bool{}
	if v := pAlert.PrometheusAlert.Labels["notify"]; v != "" {
		for _, value := range strings.Split(string(v), " ") {
			notifyMap[value] = true
		}
	}

	if alert.Notify != "" {
		notifyMap[alert.Notify] = true
	}

	var notifySlice []string
	for key, _ := range notifyMap {
		notifySlice = append(notifySlice, key)
	}

	labels := prometheus.LabelSet{
		"alertname": prometheus.LabelValue(alertname),
		"notify":    prometheus.LabelValue(strings.Join(notifySlice, " ")),
	}
	if *customerLabel != "" {
		labels["customer"] = prometheus.LabelValue(*customerLabel)
	}

	annotations := prometheus.LabelSet{
		"location":  prometheus.LabelValue(alert.Host),
		"component": prometheus.LabelValue(alert.Service),
		"type":      prometheus.LabelValue(alert.Type),
		"severity":  prometheus.LabelValue(alert.State),
		"creator":   prometheus.LabelValue("nagios"),
	}
	if alert.Message != "" {
		annotations["message"] = prometheus.LabelValue(alert.Message)
	}

	if alert.Note != "" {
		annotations["link"] = prometheus.LabelValue(alert.Note)
	}

	pAlert.PrometheusAlert.Labels = labels
	pAlert.PrometheusAlert.Annotations = annotations

	u, err := url.Parse(*generatorUrl)
	if err != nil {
		log.Fatal(err)
	}
	q := u.Query()
	if alert.Type == "host" {
		q.Set("host", alert.Host)
	} else if alert.Type == "service" {
		q.Set("servicegroup", alert.Service)
	}
	q.Set("style", "detail")
	q.Set("limit", "1000")
	u.RawQuery = q.Encode()

	pAlert.PrometheusAlert.GeneratorURL = u.String()

	pAlert.Metadata.LastUpdate = time.Now().UTC()

	if pAlert.PrometheusAlert.StartsAt.IsZero() {
		// Set the start at time to 1s ago to avoid conflicts with recoveries on an empty database
		pAlert.PrometheusAlert.StartsAt = time.Now().UTC().Add(-time.Second)
	}

	if alert.NotificationType == "RECOVERY" {
		pAlert.PrometheusAlert.EndsAt = time.Now().UTC()
		pAlert.Metadata.DontForgetUntil = time.Time{}
	} else {
		// Odd case that we have a recovered alert that fires before we do a prune
		pAlert.PrometheusAlert.EndsAt = time.Time{}
		pAlert.Metadata.DontForgetUntil = time.Time{}
		if firePeriod, err := time.ParseDuration(alert.FirePeriod); err == nil && firePeriod > 0 {
			pAlert.Metadata.DontForgetUntil = pAlert.Metadata.LastUpdate.Add(firePeriod)
		} else {
			log.Error(err)
		}
	}
}

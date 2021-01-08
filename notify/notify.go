package notify

import (
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/fanaticscripter/EggContractor/config"
)

type Notifier interface {
	Notify(m Notification) error
}

type Notification interface {
	Title() string
	Message() string
	URL() string
	Timestamp() time.Time
	// Returns notifier-specific parameters.
	NotifierParams(notifierId string) map[string]interface{}
}

func NotificationWorker(conf config.NotificationConfig, notifications <-chan Notification) {
	notifiers := make([]Notifier, 0)
	if conf.Pushover.On {
		notifiers = append(notifiers, NewPushoverNotifier(conf))
	}
	var wg sync.WaitGroup
	for m := range notifications {
		log.Debugf("queued up notification '%s'", m.Title())
		for _, n := range notifiers {
			wg.Add(1)
			go func(notifier Notifier, notification Notification) {
				if err := notifier.Notify(notification); err != nil {
					log.Error(err)
				}
				wg.Done()
			}(n, m)
		}
	}
	wg.Wait()
}

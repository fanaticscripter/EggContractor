package notify

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"

	"github.com/fanaticscripter/EggContractor/config"
)

type PushoverSound string

const _pushoverApiUrl = "https://api.pushover.net/1/messages.json"

const (
	PushoverSoundUserDefault  PushoverSound = ""
	PushoverSoundPushover     PushoverSound = "pushover"
	PushoverSoundBike         PushoverSound = "bike"
	PushoverSoundBugle        PushoverSound = "bugle"
	PushoverSoundCashRegister PushoverSound = "cashregister"
	PushoverSoundClassical    PushoverSound = "classical"
	PushoverSoundCosmic       PushoverSound = "cosmic"
	PushoverSoundFalling      PushoverSound = "falling"
	PushoverSoundGamelan      PushoverSound = "gamelan"
	PushoverSoundIncoming     PushoverSound = "incoming"
	PushoverSoundIntermission PushoverSound = "intermission"
	PushoverSoundMagic        PushoverSound = "magic"
	PushoverSoundMechanical   PushoverSound = "mechanical"
	PushoverSoundPianobar     PushoverSound = "pianobar"
	PushoverSoundSiren        PushoverSound = "siren"
	PushoverSoundSpaceAlarm   PushoverSound = "spacealarm"
	PushoverSoundTugBoat      PushoverSound = "tugboat"
	PushoverSoundAlien        PushoverSound = "alien"
	PushoverSoundClimb        PushoverSound = "climb"
	PushoverSoundPersistent   PushoverSound = "persistent"
	PushoverSoundEcho         PushoverSound = "echo"
	PushoverSoundUpDown       PushoverSound = "updown"
	PushoverSoundVibrate      PushoverSound = "vibrate"
	PushoverSoundNone         PushoverSound = "none"
)

type pushoverNotifier struct {
	apiKey     string
	userKey    string
	httpClient *http.Client
}

func NewPushoverNotifier(c config.NotificationConfig) Notifier {
	return &pushoverNotifier{
		apiKey:  c.Pushover.APIKey,
		userKey: c.Pushover.UserKey,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (n *pushoverNotifier) Notify(m Notification) error {
	notifierParams := m.NotifierParams("pushover")
	sound := PushoverSoundUserDefault
	if isound, ok := notifierParams["sound"]; ok {
		s, ok := isound.(PushoverSound)
		if ok {
			sound = s
		}
	}
	params := url.Values{}
	params.Set("token", n.apiKey)
	params.Set("user", n.userKey)
	params.Set("message", m.Message())
	if title := m.Title(); title != "" {
		params.Set("title", title)
	}
	if url := m.URL(); url != "" {
		params.Set("url", url)
	}
	if sound != PushoverSoundUserDefault {
		params.Set("sound", string(sound))
	}
	if timestamp := m.Timestamp(); !timestamp.IsZero() {
		params.Set("timestamp", fmt.Sprintf("%d", timestamp.Unix()))
	}
	params.Set("html", "1")
	resp, err := n.httpClient.PostForm(_pushoverApiUrl, params)
	if err != nil {
		return errors.Wrapf(err, "error sending Pushover notification \"%s\"", m.Title())
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrapf(err, "error sending Pushover notification \"%s\": HTTP %d",
				m.Title(), resp.StatusCode)
		}
		return errors.Errorf("error sending Pushover notification \"%s\": HTTP %d: %s",
			m.Title(), resp.StatusCode, string(body))
	}
	return nil
}

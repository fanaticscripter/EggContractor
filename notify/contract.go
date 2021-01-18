package notify

import (
	"bytes"
	"fmt"
	"html/template"
	"time"

	"github.com/fanaticscripter/EggContractor/api"
	"github.com/fanaticscripter/EggContractor/util"
)

const _contractMessageTextTmpl = `
{{- define "rewards" -}}
{{- range .}}
- <b>{{.Goal | numfmtwhole}}</b>: {{.Type}} {{.Name}} x{{.Count}};
{{- end}}
{{- end -}}

<b>{{.Name}}</b> ({{.Id}})

Egg: <b>{{.EggType.Display}}</b>
Max coop size: <b>{{.MaxCoopSize}}</b>
Time to complete: <b>{{.Duration | days}}d</b>
Token interval: <b>{{.TokenIntervalMinutes}}m</b>
Expires: <b>{{.ExpiryTime | fmtdate}}</b>

{{if eq (.RewardTiers | len) 2 -}}
{{$eliterewards := (index .RewardTiers 0).Rewards -}}
{{$standardrewards := (index .RewardTiers 1).Rewards -}}

Elite tier:{{template "rewards" $eliterewards}}
Required rate: <b>{{hourlyrate (finalgoal $eliterewards) .Duration | numfmt}}/hr</b>

Standard tier:{{template "rewards" $standardrewards}}
Required rate: <b>{{hourlyrate (finalgoal $standardrewards) .Duration | numfmt}}/hr</b>
{{else -}}
{{template "rewards" .Rewards}}
Required rate: <b>{{hourlyrate (finalgoal .Rewards) .Duration | numfmt}}/hr</b>
{{end -}}
`

var _contractMessageTmpl *template.Template

type ContractNotification struct {
	title     string
	message   string
	timestamp time.Time
}

func init() {
	_contractMessageTmpl = template.Must(template.New("").Funcs(template.FuncMap{
		"days":        func(d time.Duration) int { return int(d.Hours() / 24) },
		"finalgoal":   func(r []*api.Reward) float64 { return r[len(r)-1].Goal },
		"hourlyrate":  func(goal float64, d time.Duration) float64 { return goal / d.Hours() },
		"fmtdate":     util.FormatDate,
		"numfmt":      util.Numfmt,
		"numfmtwhole": util.NumfmtWhole,
	}).Parse(_contractMessageTextTmpl))
}

func NewContractNotification(c *api.ContractProperties) (*ContractNotification, error) {
	title := fmt.Sprintf("EggContractor: new contract \"%s\"", c.Name)
	var buf bytes.Buffer
	if err := _contractMessageTmpl.Execute(&buf, c); err != nil {
		return nil, err
	}
	message := buf.String()
	return &ContractNotification{
		title:     title,
		message:   message,
		timestamp: time.Now(),
	}, nil
}

func (n ContractNotification) Title() string {
	return n.title
}

func (n ContractNotification) Message() string {
	return n.message
}

func (n ContractNotification) URL() string {
	return ""
}

func (n ContractNotification) Timestamp() time.Time {
	return n.timestamp
}

func (n ContractNotification) NotifierParams(notifierId string) map[string]interface{} {
	switch {
	case notifierId == "pushover":
		return map[string]interface{}{
			"sound": PushoverSoundMagic,
		}
	default:
		return nil
	}
}

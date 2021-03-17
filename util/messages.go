package util

import "html/template"

const MsgNoActiveContracts = "No active contracts. Focus on growing the home farm!"

var HTMLMsgNoActiveContracts = HTMLMsg(MsgNoActiveContracts)

func HTMLMsg(m string) template.HTML {
	return template.HTML(template.HTMLEscapeString(m))
}

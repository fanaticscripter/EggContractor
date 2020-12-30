package api

import (
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

const _apiPrefix = "http://www.auxbrain.com/ei"

var _client *http.Client

func init() {
	_client = &http.Client{
		Timeout: 5 * time.Second,
	}
}

func Request(endpoint string, reqMsg proto.Message, respMsg proto.Message) error {
	apiUrl := _apiPrefix + endpoint
	reqBin, err := proto.Marshal(reqMsg)
	if err != nil {
		return errors.Wrapf(err, "marshaling payload %+v for %s", reqMsg, apiUrl)
	}
	log.Infof("POST %s: %+v", apiUrl, reqMsg)
	enc := base64.StdEncoding
	resp, err := _client.PostForm(apiUrl, url.Values{"data": {enc.EncodeToString(reqBin)}})
	if err != nil {
		return errors.Wrapf(err, "POST %s", apiUrl)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrapf(err, "POST %s", apiUrl)
	}
	if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
		return errors.Errorf("POST %s: HTTP %d: %#v", apiUrl, resp.StatusCode, string(body))
	}
	respBinBuf := make([]byte, enc.DecodedLen(len(body)))
	n, err := enc.Decode(respBinBuf, body)
	if err != nil {
		return errors.Wrapf(err, "base64 decoding %s reponse (%#v)", apiUrl, string(body))
	}
	err = proto.Unmarshal(respBinBuf[:n], respMsg)
	if err != nil {
		return errors.Wrapf(err, "unmarshaling %s response (%#v)", apiUrl, string(body))
	}
	return nil
}

func RequestFirstContact(payload *FirstContactRequestPayload) (*FirstContact, error) {
	resp := &FirstContact{}
	err := Request("/first_contact", payload, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func RequestCoopStatus(payload *CoopStatusRequestPayload) (*CoopStatus, error) {
	resp := &CoopStatus{}
	err := Request("/coop_status", payload, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

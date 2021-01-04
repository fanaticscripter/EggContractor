package api

import (
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"net/url"
	"runtime"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

const ClientVersion = 26

var _apiPrefix = "http://www.auxbrain.com/ei"

var _client *http.Client

func init() {
	timeout := 5 * time.Second
	if runtime.GOOS == "js" && runtime.GOARCH == "wasm" {
		// Use CORS proxy in the browser setting.
		_apiPrefix = "https://cors-anywhere.herokuapp.com/http://www.auxbrain.com/ei"
		// cors-anywhere may respond very slowly.
		timeout = 20 * time.Second
	}
	_client = &http.Client{
		Timeout: timeout,
	}
}

func Request(endpoint string, reqMsg proto.Message, respMsg proto.Message) error {
	apiUrl := _apiPrefix + endpoint
	reqBin, err := proto.Marshal(reqMsg)
	if err != nil {
		return errors.Wrapf(err, "marshaling payload %+v for %s", reqMsg, apiUrl)
	}
	enc := base64.StdEncoding
	reqDataEncoded := enc.EncodeToString(reqBin)
	log.Infof("POST %s: %+v", apiUrl, reqMsg)
	log.Debugf("POST %s data=%s", apiUrl, reqDataEncoded)
	resp, err := _client.PostForm(apiUrl, url.Values{"data": {reqDataEncoded}})
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

func RequestPeriodicals(payload *GetPeriodicalsRequestPayload) (*Periodicals, error) {
	if payload.ClientVersion == 0 {
		payload.ClientVersion = ClientVersion
	}
	resp := &Periodicals{}
	err := Request("/get_periodicals", payload, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

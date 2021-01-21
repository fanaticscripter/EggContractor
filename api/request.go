package api

import (
	"context"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"net/url"
	"runtime"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context/ctxhttp"
	"google.golang.org/protobuf/proto"
)

const (
	ClientVersion = 27
	AppVersion    = "1.20.0"
)

var _apiPrefix = "http://afx-2-dot-auxbrainhome.appspot.com"

var _client *http.Client

func init() {
	if runtime.GOOS == "js" && runtime.GOARCH == "wasm" {
		// Use CORS proxy in the browser setting.
		_apiPrefix = "https://wasmegg.zw.workers.dev/?url=http://afx-2-dot-auxbrainhome.appspot.com"
	}
	_client = &http.Client{
		Timeout: 5 * time.Second,
	}
}

func Request(endpoint string, reqMsg proto.Message, respMsg proto.Message) error {
	return RequestWithContext(context.Background(), endpoint, reqMsg, respMsg)
}

func RequestWithContext(ctx context.Context, endpoint string, reqMsg proto.Message, respMsg proto.Message) error {
	apiUrl := _apiPrefix + endpoint
	reqBin, err := proto.Marshal(reqMsg)
	if err != nil {
		return errors.Wrapf(err, "marshaling payload %+v for %s", reqMsg, apiUrl)
	}
	enc := base64.StdEncoding
	reqDataEncoded := enc.EncodeToString(reqBin)
	log.Infof("POST %s: %+v", apiUrl, reqMsg)
	log.Debugf("POST %s data=%s", apiUrl, reqDataEncoded)
	resp, err := ctxhttp.PostForm(ctx, _client, apiUrl, url.Values{"data": {reqDataEncoded}})
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
	return RequestFirstContactWithContext(context.Background(), payload)
}

func RequestFirstContactWithContext(ctx context.Context, payload *FirstContactRequestPayload) (*FirstContact, error) {
	if payload.ClientVersion == 0 {
		payload.ClientVersion = ClientVersion
	}
	if payload.DeviceId == "" {
		payload.DeviceId = uuid.New().String()
	}
	resp := &FirstContactResponsePayload{}
	err := RequestWithContext(ctx, "/ei/first_contact", payload, resp)
	if err != nil {
		return nil, err
	}
	return resp.Payload, nil
}

func RequestCoopStatus(payload *CoopStatusRequestPayload) (*CoopStatus, error) {
	return RequestCoopStatusWithContext(context.Background(), payload)
}

func RequestCoopStatusWithContext(ctx context.Context, payload *CoopStatusRequestPayload) (*CoopStatus, error) {
	resp := &CoopStatusResponsePayload{}
	err := RequestWithContext(ctx, "/ei/coop_status", payload, resp)
	if err != nil {
		return nil, err
	}
	return resp.Status, nil
}

func RequestPeriodicals(payload *GetPeriodicalsRequestPayload) (*Periodicals, error) {
	return RequestPeriodicalsWithContext(context.Background(), payload)
}

func RequestPeriodicalsWithContext(ctx context.Context, payload *GetPeriodicalsRequestPayload) (*Periodicals, error) {
	if payload.CurrentClientVersion == 0 {
		payload.CurrentClientVersion = ClientVersion
	}
	resp := &GetPeriodicalsResponsePayload{}
	err := RequestWithContext(ctx, "/ei/get_periodicals", payload, resp)
	if err != nil {
		return nil, err
	}
	return resp.Periodicals, nil
}

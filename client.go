package twiliogae

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/appengine/urlfetch"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client interface {
	AccountSid() string
	AuthToken() string
	RootUrl() string
	post(context.Context, url.Values, string) ([]byte, error)
}

type TwilioClient struct {
	accountSid string
	authToken  string
	rootUrl    string
}

func NewClient(accountSid, authToken string) *TwilioClient {
	rootUrl := fmt.Sprintf("%s/%s", "/2010-04-01/Accounts", accountSid)
	return &TwilioClient{accountSid, authToken, rootUrl}
}

func (client *TwilioClient) post(c context.Context, values url.Values, uri string) ([]byte, error) {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", "https://api.twilio.com", uri), strings.NewReader(values.Encode()))

	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(client.AccountSid(), client.AuthToken())
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	tr := &urlfetch.Transport{
		Context:  c,
		Deadline: time.Duration(30) * time.Second,
	}

	res, err := tr.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return body, err
	}

	if res.StatusCode != 200 && res.StatusCode != 201 {
		if res.StatusCode == 500 {
			return body, Error{"Server Error"}
		} else {
			twilioError := new(TwilioError)
			json.Unmarshal(body, twilioError)
			return body, twilioError
		}
	}

	return body, err
}

func (client *TwilioClient) AccountSid() string {
	return client.accountSid
}

func (client *TwilioClient) AuthToken() string {
	return client.authToken
}

func (client *TwilioClient) RootUrl() string {
	return client.rootUrl
}

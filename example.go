package twiliogae

import (
	"flag"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"html/template"
	"net/http"
)

var (
	twilioSID    = flag.String("twilioSID", "111111111", "twilio sid auth")
	twilioToken  = flag.String("twilioToken", "111111111", "twilio token auth")
	twilioNumber = flag.String("twilioNumber", "+111111111", "twilio sms capable number")
	twilioClient = NewClient(*twilioSID, *twilioToken)
)

var templates = template.Must(template.ParseGlob("templates/*"))

func init() {
	http.HandleFunc("/", root)
	http.HandleFunc("/sms", SmsHandler)
}

func root(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "test.html", nil)
}

func SmsHandler(w http.ResponseWriter, r *http.Request) {
	qs := r.URL.Query()
	to := qs.Get("to")
	body := qs.Get("body")
	c := appengine.NewContext(r)
	if err := SendSMS(c, to, body); err != nil {
		log.Errorf(c, err.Error())
	} else {
		log.Infof(c, "sms sent succesfully to %s", to)
	}
}

func SendSMS(c context.Context, to string, body string) error {
	if len(body) == 0 {
		return fmt.Errorf("must contain 'body'")
	}
	if len(to) == 0 {
		return fmt.Errorf("must contain 'to' number")
	}
	message, err := NewMessage(c, twilioClient, *twilioNumber, to, Body(body))

	if err != nil {
		return err
	} else {
		fmt.Println(message.Status)
	}
	return nil
}

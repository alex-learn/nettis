package transport

import (
	"github.com/laher/nettis/config"
	"github.com/laher/nettis/responsebuilders"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

//initiate only
func ConnectHttp(settings config.Settings) {
	addr := GetAddress(settings, "127.0.0.1")
	httpAddr := "http://" + addr + "/"
	log.Printf("Starting connection to %s", httpAddr)
	i := -1
	message := settings.InitiateMessage
	//only one at once. -1 represents 'forever'
	for i < settings.MaxReconnects || settings.MaxReconnects < 0 {
		resp, err := http.Post(httpAddr, "text/plain", strings.NewReader(message))
		if err != nil {
			log.Printf("Error connecting: %s", err.Error())
		} else {
			if settings.Verbose {
				defer resp.Body.Close()
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Printf("Error getting response: %s", err.Error())
				} else {
					//feeding response into next request
					log.Printf("Response body: %v", body)
					message = string(body)
				}
			}
		}
	}
	log.Printf("Finished after %d connections", i+1)
}


func setupHandler(settings config.Settings) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Duration(settings.Delay) * time.Second)
		// TODO echo content type
		w.Header().Set("Content-Type", "text/plain")
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error getting request: %s", err.Error())
		} else {
			log.Printf("Received request: %s", body)
			rg := settings.ResponseGenerator
			if rg == nil {
				log.Fatalf("NO response generator set!!")
			}
			resp, err := rg.GetResponse(responsebuilders.ResponseBuilderParams{body})
			if err != nil {
				log.Printf("Error getting request: %s", err.Error())
			} else {
				w.Write([]byte(resp))
			}
		}
	})

}

// initiate not suppported
// TODO: message logging - verbose and otherwise
func ListenHttp(settings config.Settings) {
	setupHandler(settings)
	addr := GetAddress(settings, "0.0.0.0")
	log.Printf("About to listen on http://" + addr + "/")
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func ListenHttps(settings config.Settings) {
	setupHandler(settings)
	addr := GetAddress(settings, "0.0.0.0")
	log.Printf("About to listen on https://" + addr + "/")
	//TODO cert and key names ...
	err := http.ListenAndServeTLS(addr, "cert.pem", "key.pem", nil)
	if err != nil {
		log.Fatal(err)
	}
}

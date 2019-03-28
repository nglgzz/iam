package api

import (
	"log"
	"net/http"
	"net/url"

	"github.com/getsentry/raven-go"
	jsoniter "github.com/json-iterator/go"
	"github.com/julienschmidt/httprouter"
	"gitlab.skypicker.com/cs-devs/governant/services/okta"
	"gitlab.skypicker.com/cs-devs/governant/storage"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// GetOktaUserByEmail : Look up Okta user by email
func GetOktaUserByEmail(cache *storage.Cache) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var values, err = url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
		}
		var email = values.Get("email")

		userData, err := okta.GetUser(cache, email)
		if err != nil {
			http.Error(w, "Service unavailable", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		je := json.NewEncoder(w)
		if err := je.Encode(userData); err != nil {
			log.Println("[ERROR]", err.Error())
			raven.CaptureError(err, nil)
		}
	}
}

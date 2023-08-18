package api

import (
	"net/http"
	"io/ioutil"
	"time"
	"redirect-service/config"
	"redirect-service/kafka"
	"redirect-service/redis"
)

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/r/"):]

	url, err := redis.GetURLFromCache(id)
	if err != nil {
		resp, err := http.Get(config.ServiceURL + id)
		if err != nil {
			http.Error(w, "Error fetching URL", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			http.Error(w, "Error retrieving URL", resp.StatusCode)
			return
		}

		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "Error reading response body", http.StatusInternalServerError)
			return
		}

		url = string(bodyBytes)
		redis.SetURLInCache(id, url, time.Hour*1)
	}

	kafka.PublishMessage(id, url)

	http.Redirect(w, r, url, http.StatusSeeOther)
}

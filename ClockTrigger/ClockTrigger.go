package main

import (
	"net/http"
	"time"
)

func doEvery(d time.Duration, f func()) {
	for range time.Tick(d) {
		f()
	}
}

func triggerWebhook() {
	http.Get("https://imt2681-paragliding-api.herokuapp.com/paragliding/admin/api/webhook")
}

func main() {
	doEvery(10*time.Minute, triggerWebhook)
}

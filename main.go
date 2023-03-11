package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	normalizeurl "github.com/sekimura/go-normalize-url"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	// Set routing rules
	http.HandleFunc("/present", present)
	http.HandleFunc("/cleanup", cleanup)

	//Use the default DefaultServeMux.
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

var AGH_URL = os.Getenv("ADGUARD_URL")
var AGH_USER = os.Getenv("ADGUARD_USER")
var AGH_PASS = os.Getenv("ADGUARD_PASS")

type defaultPayload struct {
	Fqdn  string `json:"fqdn"`
	Value string `json:"value"`
}

type setRuleP struct {
	Rules []string `json:"rules"`
}

func getPayloadOfRawCall(w http.ResponseWriter, r *http.Request) (defaultPayload, error) {
	if r.Method == http.MethodPost {
		d := json.NewDecoder(r.Body)
		d.DisallowUnknownFields() // catch unwanted fields
		var dP defaultPayload
		err := d.Decode(&dP)
		if err != nil {
			// bad JSON or unrecognized json field
			http.Error(w, err.Error(), http.StatusBadRequest)
			return defaultPayload{}, err
		}
		return dP, nil
	} else {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return defaultPayload{}, errors.New("HTTP method not allowed")
	}
}

func generateFilterRule(dp defaultPayload) string {
	return fmt.Sprintf("|%s^$dnsrewrite=NOERROR;TXT;%s", dp.Fqdn, dp.Value)
}

func getFilters(w http.ResponseWriter) ([]string, error) {
	get, err := callUrl(http.MethodGet, "/control/filtering/status", http.NoBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotImplemented)
		return nil, err
	}
	d := json.NewDecoder(get.Body)
	b := struct {
		UserRules []string `json:"user_rules"`
	}{}
	e := d.Decode(&b)
	if e != nil {
		// bad JSON or unrecognized json field
		http.Error(w, e.Error(), http.StatusBadRequest)
		return nil, e
	}
	return b.UserRules, nil
}

func callUrl(method string, url string, body io.Reader) (*http.Response, error) {
	n, _ := normalizeurl.Normalize(AGH_URL + url)
	client := http.Client{Timeout: 2 * time.Second}
	req, err := http.NewRequest(method, n, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(AGH_USER, AGH_PASS)
	get, err := client.Do(req)
	return get, err
}

func findAndDeleteAll(s []string, item string) []string {
	index := 0
	for _, i := range s {
		if i != item {
			s[index] = i
			index++
		}
	}
	return s[:index]
}

func present(w http.ResponseWriter, r *http.Request) {
	dp, err := getPayloadOfRawCall(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	filterRule := generateFilterRule(dp)
	f, err := getFilters(w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	filterWithNew := append(f, filterRule)
	b, err := json.Marshal(setRuleP{Rules: filterWithNew})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	_, err = callUrl(http.MethodPost, "/control/filtering/set_rules", bytes.NewReader(b))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusCreated)
}

func cleanup(w http.ResponseWriter, r *http.Request) {
	dp, err := getPayloadOfRawCall(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	filterRule := generateFilterRule(dp)
	f, err := getFilters(w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	filterWithout := findAndDeleteAll(f, filterRule)
	b, err := json.Marshal(setRuleP{Rules: filterWithout})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	_, err = callUrl(http.MethodPost, "/control/filtering/set_rules", bytes.NewReader(b))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

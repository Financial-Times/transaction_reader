package main

import (
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

type httpHandlers struct {
	ers transactionReaderService
}

// e.g. everything with type Article and no end time - IN PROGRESS
func (hh *httpHandlers) getInProgressTransactionsForType(w http.ResponseWriter, r *http.Request) {
	// expect a type param
	publishType := r.URL.Query().Get("type")
	if publishType == "" {
		log.Errorf("No type param supplied")
		return
	}
	log.Infof("Get In Progress Events For Type %s", publishType)

	jason, err := ioutil.ReadFile("exampleInProgressTransactions.json")
	if err != nil {
		log.Errorf("Failed to read in file, %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write([]byte(jason))
	return
}

// everything with no end time and start time is more than two minutes ago - EXCEEDED SLA
// (will only be fixed up by republishing)
func (hh *httpHandlers) getTransactionsExceedingSLA(w http.ResponseWriter, r *http.Request) {

}

// for a given article uuid look for records with a start time after the given time and has an end time
// - ANOTHER PUBLISH HAS SUCCEEDED (we then set the end time for the previous publish to match this end time)
func (hh *httpHandlers) getTransactionsForUUID(w http.ResponseWriter, r *http.Request) {

	// expect a uuid param
	uuid := r.URL.Query().Get("uuid")
	if uuid == "" {
		log.Errorf("No uuid param supplied")
		return
	}
	after := r.URL.Query().Get("after")
	if after == "" {
		log.Errorf("No after param supplied")
		return
	}
	log.Infof("Get transactions for uuid %s which started after %s", uuid, after)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if (uuid == "GHI") {
		jason, err := ioutil.ReadFile("exampleTransactionsForGHI.json")
		if err != nil {
			log.Errorf("Failed to read in file, %v", err)
			return
		}
		w.Write([]byte(jason))
	} else {
		w.Write([]byte("[]"))
	}
	return
}

// get the event for a particular transaction id (only for debugging really)
func (hh *httpHandlers) getTransactionForTransactionID(w http.ResponseWriter, r *http.Request) {

}

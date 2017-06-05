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

}

// get the event for a particular transaction id (only for debugging really)
func (hh *httpHandlers) getTransactionForTransactionID(w http.ResponseWriter, r *http.Request) {

}

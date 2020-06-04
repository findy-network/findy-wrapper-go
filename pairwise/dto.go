package pairwise

import (
	"encoding/json"
	"strings"

	"github.com/golang/glog"
)

// Data is a pairwise data saved to indy wallet.
type Data struct {
	MyDid    string `json:"my_did,omitempty"`
	TheirDid string `json:"their_did,omitempty"`
	Metadata string `json:"metadata,omitempty"`
}

// NewData creates a pairwise data from a string returned by indy pairwise
// function. Note! The indy's string is buggy, that's why this function removes
// unnecessary \ characters from it before JSON unmarshal.
func NewData(js string) []Data {
	a := make([]Data, 0)
	if js == "" {
		js = "[]"
	}
	js = strings.Replace(js, "\"", "", -1)
	js = strings.Replace(js, "\\", "\"", -1)
	err := json.Unmarshal([]byte(js), &a)
	if err != nil {
		glog.Error("err marshalling from JSON: ", err.Error())
		return nil
	}
	return a
}

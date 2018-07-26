package main

import (
	"db"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Property struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type"` // int float string
	Units string `json:"units"`
}

type Measurement struct {
	MeasuredUnixTime float64    `json:"unixtime"`
	Properties       []Property `json:"properties"`
	ReceivedUnixTime float64
	Ip               string
}

// sends the current unix time
func sendTime(w http.ResponseWriter, r *http.Request) {
	// Accept GET requests only
	if r.Method != http.MethodGet {
		status := http.StatusMethodNotAllowed
		msg := http.StatusText(status)
		http.Error(w, msg, status)
		return
	}
	// get the current time
	now := float64(time.Now().UnixNano()) / 1e9
	log.Printf("Sending time %f", now)
	// send to the client
	nowString := strconv.FormatFloat(now, 'f', -1, 64)
	w.Write([]byte(nowString))
}

// saves a []Measurement to database
func saveMeasurements(w http.ResponseWriter, r *http.Request) {
	// Accept POST requests only
	if r.Method != http.MethodPost {
		status := http.StatusMethodNotAllowed
		msg := http.StatusText(status) + ", use POST"
		http.Error(w, msg, status)
		return
	}
	// read request body
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Println("Error reading request body")
		log.Println(err)
		status := http.StatusInternalServerError
		msg := http.StatusText(status)
		http.Error(w, msg, status)
		return
	}
	// parse json data
	measurements := []Measurement{}
	err = json.Unmarshal(body, &measurements)
	if err != nil {
		log.Println("Error parsing request body")
		log.Println(err)
		status := http.StatusBadRequest
		msg := http.StatusText(status) + ", error parsing json"
		http.Error(w, msg, status)
		return
	}
	// get metadata
	now := float64(time.Now().UnixNano()) / 1e9
	ip := strings.Split(r.RemoteAddr, ":")[0]
	// save each measurement
	log.Println("Saving", len(measurements), "measurements from", ip)
	for _, measurement := range measurements {
		// measurement metadata
		measurement.Ip = ip
		measurement.ReceivedUnixTime = now
		r, err := db.WriteDb(`
		INSERT INTO measurement (
			measured_unix_time,
			received_unix_time,
			ip
		) VALUES (?,?,?);`,
			measurement.MeasuredUnixTime,
			measurement.ReceivedUnixTime,
			measurement.Ip)
		if err != nil {
			log.Println("Error saving measurement to database")
			log.Println(err)
			status := http.StatusInternalServerError
			msg := http.StatusText(status)
			http.Error(w, msg, status)
			return
		}
		// get measurement id
		measurementId, err := r.LastInsertId()
		if err != nil {
			log.Println("Error getting measurement id")
			log.Println(err)
			status := http.StatusInternalServerError
			msg := http.StatusText(status)
			http.Error(w, msg, status)
			return
		}
		// measurement properties
		for _, property := range measurement.Properties {
			db.WriteDb(`
			INSERT INTO measurement_property (
				measurement_id,
				key,
				value,
				type,
				units
			) VALUES (?,?,?,?,?);`,
				measurementId,
				property.Key,
				property.Value,
				property.Type,
				property.Units)
		}
	}
}

func main() {
	http.HandleFunc("/measurements", saveMeasurements)
	http.HandleFunc("/time", sendTime)
	log.Println("Started on port 5678")
	http.ListenAndServe(":5678", nil)
}

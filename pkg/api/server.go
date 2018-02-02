package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"bitbucket.org/bdsengineering/perceptor/pkg/common"
	log "github.com/sirupsen/logrus"
)

func SetupHTTPServer(responder Responder) {
	// state of the program
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			responder.GetMetrics(w, r)
		} else {
			responder.NotFound(w, r)
		}
	})
	http.HandleFunc("/model", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			fmt.Fprint(w, responder.GetModel())
		} else {
			responder.NotFound(w, r)
		}
	})

	// for receiving data from perceiver
	http.HandleFunc("/pod", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Errorf("unable to read body for pod POST: %s", err.Error())
				responder.Error(w, r, err, 400)
				return
			}
			var pod common.Pod
			err = json.Unmarshal(body, &pod)
			if err != nil {
				log.Infof("unable to ummarshal JSON for pod POST: %s", err.Error())
				responder.Error(w, r, err, 400)
				return
			}
			responder.AddPod(pod)
			fmt.Fprint(w, "")
		case "PUT":
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				responder.Error(w, r, err, 400)
				return
			}
			var pod common.Pod
			err = json.Unmarshal(body, &pod)
			if err != nil {
				responder.Error(w, r, err, 400)
				return
			}
			responder.UpdatePod(pod)
			fmt.Fprint(w, "")
		case "DELETE":
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				responder.Error(w, r, err, 400)
				return
			}
			responder.DeletePod(string(body))
			fmt.Fprint(w, "")
		default:
			responder.NotFound(w, r)
		}
	})
	http.HandleFunc("/allpods", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "PUT" {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				responder.Error(w, r, err, 400)
				return
			}
			var allPods AllPods
			err = json.Unmarshal(body, &allPods)
			if err != nil {
				responder.Error(w, r, err, 400)
				return
			}
			responder.UpdateAllPods(allPods)
		} else {
			responder.NotFound(w, r)
		}
	})
	http.HandleFunc("/image", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				responder.Error(w, r, err, 400)
				return
			}
			var image common.Image
			err = json.Unmarshal(body, &image)
			if err != nil {
				responder.Error(w, r, err, 400)
				return
			}
			responder.AddImage(image)
		} else {
			responder.NotFound(w, r)
		}
	})

	// for providing data to perceiver
	http.HandleFunc("/scanresults", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			scanResults := responder.GetScanResults()
			jsonBytes, err := json.Marshal(scanResults)
			if err != nil {
				responder.Error(w, r, err, 500)
				return
			}
			fmt.Fprint(w, string(jsonBytes))
		} else {
			responder.NotFound(w, r)
		}
	})

	// for providing data to scanners
	http.HandleFunc("/nextimage", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			var wg sync.WaitGroup
			wg.Add(1)
			responder.GetNextImage(func(nextImage NextImage) {
				jsonBytes, err := json.Marshal(nextImage)
				if err != nil {
					responder.Error(w, r, err, 500)
				} else {
					fmt.Fprint(w, string(jsonBytes))
				}
				wg.Done()
			})
			wg.Wait()
		} else {
			responder.NotFound(w, r)
		}
	})

	http.HandleFunc("/finishedscan", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				responder.Error(w, r, err, 400)
				return
			}
			var scanResults FinishedScanClientJob
			err = json.Unmarshal(body, &scanResults)
			responder.PostFinishScan(scanResults)
			fmt.Fprint(w, "")
		} else {
			responder.NotFound(w, r)
		}
	})
}

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

type pipelines struct {
	Pipelines []pipeline `json:"pipelines"`
	TotalSize int        `json:"total_size"`
}

type pipeline struct {
	ID          string       `json:"id,omitempty"`
	CreatedAt   time.Time    `json:"created_at,omitempty"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Parameters  []parameters `json:"parameters,omitempty"`
}

type parameters struct {
	Name  string `json:"name"`
	Value string `json:"value,omitempty"`
}

type experiments struct {
	Experiments []experiment `json:"experiments"`
	TotalSize   int          `json:"total_size"`
}

type experiment struct {
	ID          string    `json:"id,omitempty"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
}

const baseURL string = "http://188.40.161.63:8888/apis/v1beta1/"
const url string = "http://188.40.161.63:8888/apis/v1beta1/pipelines/020356d7-a13c-41e2-8c35-c98e7c1ea65d"

func main() {
	getAllRuns()
}

func getAllPipelines() pipelines {
	var p pipelines
	url := baseURL + "pipelines"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Failed to send request to the backend: %s", err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	json.Unmarshal(body, &p)
	return p
}

func getPipeline(id string) pipeline {
	var p pipeline
	url := baseURL + "pipelines/" + id
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Failed to send request to the backend: %s", err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	json.Unmarshal(body, &p)
	return p
}

func getAllExperiments() experiments {
	var e experiments
	url := baseURL + "experiments"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Failed to send request to the backend: %s", err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	json.Unmarshal(body, &e)
	return e
}

func getExperiment(id string) experiment {
	var e experiment
	url := baseURL + "experiments/" + id
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Failed to send request to the backend: %s", err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	json.Unmarshal(body, &e)
	return e
}

func uploadPipeline(filename string, name string) pipeline {
	var p pipeline
	url := baseURL + "pipelines/upload" + "?name=" + name
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	// this step is very important
	fileWriter, err := bodyWriter.CreateFormFile("uploadfile", filename)
	if err != nil {
		log.Fatalf("error writing to buffer: %s", err.Error())
	}

	// open file handle
	fh, err := os.Open(filename)
	if err != nil {
		log.Fatalf("error opening file: %s", err.Error())
	}
	defer fh.Close()

	//iocopy
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		log.Fatalf("error: %s", err.Error())
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post(url, contentType, bodyBuf)
	if err != nil {
		log.Fatalf("error sending post: %s", err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("error reading body: %s", err.Error())
	}

	json.Unmarshal(body, &p)

	return p
}

func deleteExperiment(id string) error {
	url := baseURL + "experiments/" + id
	hc := http.Client{}
	req, err := http.NewRequest("DELETE", url, nil)
	resp, err := hc.Do(req)
	if err != nil {
		log.Fatalf("Failed to send request to the backend: %s", err.Error())
		return err
	}
	defer resp.Body.Close()

	return err
}

func createExperiment(name string, description string) experiment {
	e := experiment{
		Name:        name,
		Description: description,
	}
	url := baseURL + "experiments"

	body, err := json.Marshal(e)
	if err != nil {
		log.Fatalf("error marshaling: %s", err.Error())
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Fatalf("error sending post: %s", err.Error())
	}
	defer resp.Body.Close()

	var newE experiment
	respBody, err := ioutil.ReadAll(resp.Body)
	json.Unmarshal(respBody, &newE)

	return newE
}

func deletePipeline(id string) error {
	url := baseURL + "pipelines/" + id
	hc := http.Client{}
	req, err := http.NewRequest("DELETE", url, nil)
	resp, err := hc.Do(req)
	if err != nil {
		log.Fatalf("Failed to send request to the backend: %s", err.Error())
		return err
	}
	defer resp.Body.Close()

	return err
}

func getAllRuns() {
	//var e experiment
	url := baseURL + "runs"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Failed to send request to the backend: %s", err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	fmt.Println(string(body))
	// json.Unmarshal(body, &e)
	// return e
}

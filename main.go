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
	ID          string      `json:"id,omitempty"`
	CreatedAt   time.Time   `json:"created_at,omitempty"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Parameters  []parameter `json:"parameters,omitempty"`
}

type parameter struct {
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

type runs struct {
	Runs      []run `json:"runs"`
	TotalSize int   `json:"total_size"`
}

type run struct {
	ID                 string              `json:"id,omitempty"`
	Name               string              `json:"name,omitempty"`
	StorageState       string              `json:"storage_state,omitempty"` // STORAGESTATE_AVAILABLE (default) or STORAGESTATE_ARCHIVED
	Description        string              `json:"description,omitemtpy"`
	PipelineSpec       pipelineSpec        `json:"pipeline_spec,omitemtpy"`
	ResourceReferences []resourceReference `json:"resource_references,omitemtpy"`
	//output
	CreatedAt   time.Time `json:"created_at,omitempty"`
	ScheduledAt time.Time `json:"scheduled_at,omitempty"`
	FinishedAt  time.Time `json:"finished_at,omitempty"`
	Status      string    `json:"status,omitempty"`
	Error       string    `json:"error,omitempty"`
	Metrics     string    `json:"metrics,omitempty"` //CHANGE
}

type pipelineSpec struct {
	PipelineID       string      `json:"pipeline_id,omitempty"`
	WorkflowManifest string      `json:"workflow_manifest,omitempty"`
	PipelineManifest string      `json:"pipeline_manifest,omitempty"`
	Parameters       []parameter `json:"parameters,omitempty"`
}

type resourceReference struct {
	Key          resourceKey `json:"key,omitempty"`
	Relationship string      `json:"relationship,omitempty"` //UNKNOWN_RELATIONSHIP (default), OWNER, CREATOR
}

type resourceKey struct {
	Type string `json:"type,omitempty"` //UNKNOWN_RESOURCE_TYPE (default), EXPERIMENT, JOB
	ID   string `json:"id,omitempty"`
}

type runDetail struct {
	Run             run             `json:"run,omitempty"`
	PipelineRuntime pipelineRuntime `json:"pipeline_runtime,omitempty"`
}

type pipelineRuntime struct {
	WorkflowManifest string `json:"workflow_manifest,omitempty"`
}

const baseURL string = "http://188.40.161.99:8888/apis/v1beta1/"
const url string = "http://188.40.161.99:8888/apis/v1beta1/pipelines/020356d7-a13c-41e2-8c35-c98e7c1ea65d"
const pipelineID string = "f761e123-3637-4b8d-bd28-82abfd253e49"
const experimentID string = "04ccd65d-f6eb-47f8-ae0d-bddae719699c"
const runID string = "4dfc6e6e-e5be-11e9-908e-fa163ed02a23"

func main() {
	// p := uploadPipeline("pipeline.yaml", "MNIST")
	// fmt.Printf("%s", p.Name)

	// pipelines := getAllPipelines()
	// fmt.Println(pipelines)

	mnistPipeline := getPipeline(pipelineID)
	fmt.Println(mnistPipeline)

	// experiments := getAllExperiments()
	// fmt.Println(experiments)

	e := getExperiment(experimentID)
	fmt.Println(e)

	// allRuns := getAllRuns()
	// fmt.Println(allRuns)

	r := run{
		Name:        "Test run from sdk",
		Description: "Test run from sdk - description",
		PipelineSpec: pipelineSpec{
			PipelineID: pipelineID,
			Parameters: []parameter{
				{
					Name:  "model-export-dir",
					Value: "/mnt/export",
				},
				{
					Name:  "train-steps",
					Value: "50",
				},
				{
					Name:  "learning-rate",
					Value: "0.01",
				},
				{
					Name:  "batch-size",
					Value: "100",
				},
				{
					Name:  "pvc-name",
					Value: "local-storage",
				},
			},
		},
		ResourceReferences: []resourceReference{
			{
				Key: resourceKey{
					ID:   experimentID,
					Type: "EXPERIMENT",
				},
				Relationship: "OWNER",
			},
		},
	}
	fmt.Println(r)

	rDetail := createRun(r)
	fmt.Printf("Name: %s", rDetail.Run.Name)

	// rDetail := getRun(runID)
	// fmt.Println(rDetail)

	// _ = getRun(runID)
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
	if err != nil {
		log.Fatalf("error reading from resp.Body: %s", err.Error())
	}

	err = json.Unmarshal(body, &p)
	if err != nil {
		log.Printf("\noutput: %v", string(body))
		log.Fatalf("error unmarshaling: %s", err.Error())
	}
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
	if err != nil {
		log.Fatalf("error reading from resp.Body: %s", err.Error())
	}

	err = json.Unmarshal(body, &p)
	if err != nil {
		log.Printf("\noutput: %v", string(body))
		log.Fatalf("error unmarshaling: %s", err.Error())
	}
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
	if err != nil {
		log.Fatalf("error reading from resp.Body: %s", err.Error())
	}

	err = json.Unmarshal(body, &e)
	if err != nil {
		log.Printf("\noutput: %v", string(body))
		log.Fatalf("error unmarshaling: %s", err.Error())
	}
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
	if err != nil {
		log.Fatalf("error reading from resp.Body: %s", err.Error())
	}

	err = json.Unmarshal(body, &e)
	if err != nil {
		log.Printf("\noutput: %v", string(body))
		log.Fatalf("error unmarshaling: %s", err.Error())
	}
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
		log.Fatalf("error reading from resp.Body: %s", err.Error())
	}

	err = json.Unmarshal(body, &p)
	if err != nil {
		log.Printf("\noutput: %v", string(body))
		log.Fatalf("error unmarshaling: %s", err.Error())
	}

	return p
}

func deleteExperiment(id string) error {
	url := baseURL + "experiments/" + id
	hc := http.Client{}
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		log.Fatalf("error creating new request: %s", err.Error())
	}
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
	if err != nil {
		log.Fatalf("error reading from resp.Body: %s", err.Error())
	}

	err = json.Unmarshal(respBody, &newE)
	if err != nil {
		log.Printf("\noutput: %v", string(respBody))
		log.Fatalf("error unmarshaling: %s", err.Error())
	}

	return newE
}

func deletePipeline(id string) error {
	url := baseURL + "pipelines/" + id
	hc := http.Client{}
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		log.Fatalf("error creating new request: %s", err.Error())
	}
	resp, err := hc.Do(req)
	if err != nil {
		log.Fatalf("Failed to send request to the backend: %s", err.Error())
		return err
	}
	defer resp.Body.Close()

	return err
}

func getAllRuns() runs {
	var r runs
	url := baseURL + "runs"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Failed to send request to the backend: %s", err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("error reading from resp.Body: %s", err.Error())
	}

	// fmt.Println(string(body))
	err = json.Unmarshal(body, &r)
	if err != nil {
		log.Printf("\noutput: %v", string(body))
		log.Fatalf("error unmarshaling: %s", err.Error())
	}
	return r
}

func createRun(r run) runDetail {
	url := baseURL + "runs"
	var rDetail runDetail
	// fmt.Println(url)

	body, err := json.Marshal(r)
	if err != nil {
		log.Fatalf("error marshaling: %s", err.Error())
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Fatalf("error sending post: %s", err.Error())
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("error reading from resp.Body: %s", err.Error())
	}

	// fmt.Printf("\nResponse:\n%s", string(respBody))

	err = json.Unmarshal(respBody, &rDetail)
	if err != nil {
		log.Printf("\noutput: %v", string(respBody))
		log.Fatalf("error unmarshaling: %s", err.Error())
	}

	return rDetail
}

func getRun(id string) runDetail {
	var r runDetail
	url := baseURL + "runs/" + id
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Failed to send request to the backend: %s", err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("error reading from resp.Body: %s", err.Error())
	}
	//fmt.Println(string(body))

	err = json.Unmarshal(body, &r)
	if err != nil {
		log.Printf("\noutput: %v", string(body))
		log.Fatalf("error unmarshaling: %s", err.Error())
	}
	return r
}

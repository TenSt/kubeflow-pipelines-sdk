package sdk

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

// Pipelines describe a list of pipelines
type Pipelines struct {
	Pipelines []Pipeline `json:"pipelines"`
	TotalSize int        `json:"total_size"`
}

// Pipeline describe a single pipeline
type Pipeline struct {
	ID          string      `json:"id,omitempty"`
	CreatedAt   time.Time   `json:"created_at,omitempty"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Parameters  []Parameter `json:"parameters,omitempty"`
}

// Parameter describe pipeline parameter
type Parameter struct {
	Name  string `json:"name"`
	Value string `json:"value,omitempty"`
}

// Experiments describe a list of experiments
type Experiments struct {
	Experiments []Experiment `json:"experiments"`
	TotalSize   int          `json:"total_size"`
}

// Experiment describe a single experiment
type Experiment struct {
	ID          string    `json:"id,omitempty"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
}

// Runs describe a list of runs
type Runs struct {
	Runs      []Run `json:"runs"`
	TotalSize int   `json:"total_size"`
}

// Run describe a single run
type Run struct {
	ID                 string              `json:"id,omitempty"`
	Name               string              `json:"name,omitempty"`
	StorageState       string              `json:"storage_state,omitempty"` // STORAGESTATE_AVAILABLE (default) or STORAGESTATE_ARCHIVED
	Description        string              `json:"description,omitemtpy"`
	PipelineSpec       PipelineSpec        `json:"pipeline_spec,omitemtpy"`
	ResourceReferences []ResourceReference `json:"resource_references,omitemtpy"`
	//output
	CreatedAt   time.Time `json:"created_at,omitempty"`
	ScheduledAt time.Time `json:"scheduled_at,omitempty"`
	FinishedAt  time.Time `json:"finished_at,omitempty"`
	Status      string    `json:"status,omitempty"`
	Error       string    `json:"error,omitempty"`
	Metrics     string    `json:"metrics,omitempty"` //CHANGE
}

// PipelineSpec describe a single specification of pipeline
type PipelineSpec struct {
	PipelineID       string      `json:"pipeline_id,omitempty"`
	WorkflowManifest string      `json:"workflow_manifest,omitempty"`
	PipelineManifest string      `json:"pipeline_manifest,omitempty"`
	Parameters       []Parameter `json:"parameters,omitempty"`
}

// ResourceReference describe a single reference to parent resource
type ResourceReference struct {
	Key          ResourceKey `json:"key,omitempty"`
	Relationship string      `json:"relationship,omitempty"` //UNKNOWN_RELATIONSHIP (default), OWNER, CREATOR
}

// ResourceKey describe a single key of a parent resource
type ResourceKey struct {
	Type string `json:"type,omitempty"` //UNKNOWN_RESOURCE_TYPE (default), EXPERIMENT, JOB
	ID   string `json:"id,omitempty"`
}

// RunDetail describe an output details of a run
type RunDetail struct {
	Run             Run             `json:"run,omitempty"`
	PipelineRuntime PipelineRuntime `json:"pipeline_runtime,omitempty"`
}

// PipelineRuntime describe a workflow manifest as an JSON string of a pipeline
type PipelineRuntime struct {
	WorkflowManifest string `json:"workflow_manifest,omitempty"`
}

//KfPipelineClient is a client struct
type KfPipelineClient struct {
	BaseURL string `json:"baseurl"`
}

// GetClient returns client with BaseURL
// Please mind that "/apis/v1beta1/" will be added at the end of BaseURL
func GetClient(url string) KfPipelineClient {
	client := KfPipelineClient{
		BaseURL: os.Getenv("KUBEFLOW_PIPELINE_API") + "/apis/v1beta1/",
	}
	return client
}

// GetAllPipelines will return all pipelines
func (c *KfPipelineClient) GetAllPipelines() Pipelines {
	var p Pipelines
	url := c.BaseURL + "pipelines"
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

// GetPipeline will return pipeline based on provided ID
func (c *KfPipelineClient) GetPipeline(id string) Pipeline {
	var p Pipeline
	url := c.BaseURL + "pipelines/" + id
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

// GetAllExperiments will return all experiments
func (c *KfPipelineClient) GetAllExperiments() Experiments {
	var e Experiments
	url := c.BaseURL + "experiments"
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

// GetExperiment will return experiment based on ID
func (c *KfPipelineClient) GetExperiment(id string) Experiment {
	var e Experiment
	url := c.BaseURL + "experiments/" + id
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

// UploadPipeline will upload pipeline using provided path to a file
func (c *KfPipelineClient) UploadPipeline(filename string, name string) Pipeline {
	var p Pipeline
	url := c.BaseURL + "pipelines/upload" + "?name=" + name
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

// DeleteExperiment will delete experiment based on ID
func (c *KfPipelineClient) DeleteExperiment(id string) error {
	url := c.BaseURL + "experiments/" + id
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

// CreateExperiment will create new experiment
func (c *KfPipelineClient) CreateExperiment(name string, description string) Experiment {
	e := Experiment{
		Name:        name,
		Description: description,
	}
	url := c.BaseURL + "experiments"

	body, err := json.Marshal(e)
	if err != nil {
		log.Fatalf("error marshaling: %s", err.Error())
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Fatalf("error sending post: %s", err.Error())
	}
	defer resp.Body.Close()

	var newE Experiment
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

// DeletePipeline will delete pipeline based on ID
func (c *KfPipelineClient) DeletePipeline(id string) error {
	url := c.BaseURL + "pipelines/" + id
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

// GetAllRuns will return all runs
func (c *KfPipelineClient) GetAllRuns() Runs {
	var r Runs
	url := c.BaseURL + "runs"
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

// CreateRun will create new Run of a pipeline based on the provided Run struct and return RunDetail
func (c *KfPipelineClient) CreateRun(r Run) RunDetail {
	url := c.BaseURL + "runs"
	var rDetail RunDetail
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

// GetRun return RunDetail based on ID
func (c *KfPipelineClient) GetRun(id string) RunDetail {
	var r RunDetail
	url := c.BaseURL + "runs/" + id
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

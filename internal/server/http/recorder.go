package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var (
	ErrCreatingRecordDir       = errors.New("impossible create record directory")
	ErrOpenRecordFile          = errors.New("impossible open record file")
	ErrReadingOutputRecordFile = errors.New("error trying to parse the record file")
	ErrMarshallingRecordFile   = errors.New("error during the marshalling process of the record file")
	ErrWritingRecordFile       = errors.New("error trying to write on the record file")
)

// RecorderHTTP service to Record the return output of the request
type RecorderHTTP interface {
	// Record save the return output from the missing request on the imposters
	Record(req *http.Request, resp ResponseRecorder) error
}

// Recorder implementation of the RecorderHTTP
type Recorder struct {
	outputPathFile string
}

// ResponseRecorder response data transfer object
type ResponseRecorder struct {
	Status  int
	Headers http.Header
	Body    string
}

// NewRecorder initialise the Recorder
func NewRecorder(outputPathFile string) Recorder {
	return Recorder{
		outputPathFile: outputPathFile,
	}
}

func (r Recorder) Record(req *http.Request, resp ResponseRecorder) error {
	imposterRequest := r.prepareImposterRequest(req)
	imposterResponse, err := r.prepareImposterResponse(resp)
	if err != nil {
		return err
	}

	imposter := Imposter{
		Request:  imposterRequest,
		Response: imposterResponse,
	}

	f, err := r.prepareOutputFile()
	if err != nil {
		return err
	}
	defer f.Close()

	var b []byte
	switch {
	case strings.HasSuffix(r.outputPathFile, jsonImposterExtension):
		b, err = r.recordOnJSON(f, imposter)
		if err != nil {
			return err
		}
	case strings.HasSuffix(r.outputPathFile, yamlImposterExtension) || strings.HasSuffix(r.outputPathFile, ymlImposterExtension) :
		b, err = r.recordOnYAML(f, imposter)
		if err != nil {
			return err
		}
	default:
		return errors.New("file extension not recognized")
	}

	_ = f.Truncate(0)
	_, _ = f.Seek(0, 0)

	if _, err := f.Write(b); err != nil {
		return fmt.Errorf("%v: %w", err, ErrWritingRecordFile)
	}

	return nil
}

// RecorderNoop an implementation of the RecorderHTTP without any functionality
type RecorderNoop struct{}

func (r RecorderNoop) Record(req *http.Request, resp ResponseRecorder) error {
	return nil
}

func (r Recorder) prepareOutputFile() (*os.File, error) {
	dir := filepath.Dir(r.outputPathFile)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.Mkdir(dir, 0755)
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, ErrCreatingRecordDir)
		}
	}

	f, err := os.OpenFile(r.outputPathFile, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, ErrOpenRecordFile)
	}

	return f, nil
}

func (r Recorder) prepareImposterRequest(req *http.Request) Request {
	headers := make(map[string]string, len(req.Header))
	for k, v := range req.Header {
		for _, val := range v {
			// TODO: configure which headers don't you want to store or more commons like Postman??
			headers[k] = val
		}
	}

	params := make(map[string]string, len(req.URL.Query()))
	query := req.URL.Query()
	for k, v := range query {
		params[k] = v[0]
	}

	imposterRequest := Request{
		Method:   req.Method,
		Endpoint: req.URL.Path,
		Headers:  &headers,
		Params:   &params,
	}

	return imposterRequest
}

func (r Recorder) prepareImposterResponse(resp ResponseRecorder) (Response, error) {
	headers := make(map[string]string, len(resp.Headers))
	for k, v := range resp.Headers {
		for _, val := range v {
			headers[k] = val
		}
	}

	response := Response{
		Status:  resp.Status,
		Headers: &headers,
		Body: resp.Body,
	}

	return response, nil
}

func (r Recorder) recordOnJSON(file *os.File, imposter Imposter) ([]byte,error) {
	var imposters []Imposter
	bytes, _ := ioutil.ReadAll(file)
	if err := json.Unmarshal(bytes, &imposters); err != nil && len(bytes) > 0 {
		return nil, fmt.Errorf("%v: %w", err, ErrReadingOutputRecordFile)
	}

	//TODO: create an inMemory to store which imposters are saved during this session to avoid duplicated
	imposters = append(imposters, imposter)
	b, err := json.Marshal(imposters)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, ErrMarshallingRecordFile)
	}

	return b, nil
}

func (r Recorder) recordOnYAML(file *os.File, imposter Imposter) ([]byte,error) {
	var imposters []Imposter
	bytes, _ := ioutil.ReadAll(file)
	if err := yaml.Unmarshal(bytes, &imposters); err != nil && len(bytes) > 0 {
		return nil, fmt.Errorf("%v: %w", err, ErrReadingOutputRecordFile)
	}

	//TODO: create an inMemory to store which imposters are saved during this session to avoid duplicated
	imposters = append(imposters, imposter)
	b, err := yaml.Marshal(imposters)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, ErrMarshallingRecordFile)
	}

	return b, nil
}


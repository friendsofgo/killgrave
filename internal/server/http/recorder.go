package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

var (
	ErrCreatingRecordDir       = errors.New("impossible create record directory")
	ErrCreatingRecordFile      = errors.New("impossible create record file")
	ErrOpenRecordFile          = errors.New("impossible open record file")
	ErrTryingToReadBody        = errors.New("impossible read the body response")
	ErrReadingOutputRecordFile = errors.New("error trying to parse the record file")
	ErrMarshallingRecordFile   = errors.New("error during the marshalling process of the record file")
	ErrWritingRecordFile       = errors.New("error trying to write on the record file")
)

type Recorder struct {
	outputPathFile string
}

func NewRecorder(outputPathFile string) Recorder {
	return Recorder{
		outputPathFile: outputPathFile,
	}
}

func (r Recorder) Record(req *http.Request, resp *http.Response) error {
	f, err := r.prepareOutputFile()
	if err != nil {
		return err
	}
	defer f.Close()

	var imposters []Imposter
	bytes, _ := ioutil.ReadAll(f)
	if err := json.Unmarshal(bytes, &imposters); err != nil && len(bytes) > 0 {
		return fmt.Errorf("%v: %w", err, ErrReadingOutputRecordFile)
	}

	imposterRequest := r.prepareImposterRequest(req)
	imposterResponse, err := r.prepareImposterResponse(resp)
	if err != nil {
		return err
	}

	imposter := Imposter{
		Request:  imposterRequest,
		Response: imposterResponse,
	}

	//TODO: create an inMemory to store which imposters are saved during this session to avoid duplicated
	imposters = append(imposters, imposter)
	b, err := json.Marshal(imposters)
	if err != nil {
		return fmt.Errorf("%v: %w", err, ErrMarshallingRecordFile)
	}

	if _, err := f.Write(b); err != nil {
		return fmt.Errorf("%v: %w", err, ErrWritingRecordFile)
	}

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

	var f *os.File
	if _, err := os.Stat(r.outputPathFile); os.IsNotExist(err) {
		f, err = os.Create(r.outputPathFile)
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, ErrCreatingRecordFile)
		}
	} else {
		f, err = os.OpenFile(r.outputPathFile, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			return nil, fmt.Errorf("%v: %w", err, ErrOpenRecordFile)
		}
	}

	return f, nil
}

func (r Recorder) prepareImposterRequest(req *http.Request) Request {
	headers := make(map[string]string, len(req.Header))
	for k, v := range req.Header {
		for _, val := range v {
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

func (r Recorder) prepareImposterResponse(resp *http.Response) (Response, error) {
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Response{}, fmt.Errorf("%v: %w", err, ErrTryingToReadBody)
	}
	defer resp.Body.Close()

	headers := make(map[string]string, len(resp.Header))
	for k, v := range resp.Header {
		for _, val := range v {
			headers[k] = val
		}
	}

	response := Response{
		Status:  resp.StatusCode,
		Body:    string(b),
		Headers: &headers,
	}

	return response, nil
}

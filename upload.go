package upload

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/levigross/grequests"
)

const version = "0.1.0"

// File is the File's return type
type File struct {
	ObjectID string `json:"objectID"`
	URL      string `json:"url"`
}

// Error is the LeanCloud API Server API common error format
type Error struct {
	Code         int    `json:"code"`
	Content      string `json:"error"`
	ErrorEventID string `json:"errorEventID"`
}

func (err Error) Error() string {
	return fmt.Sprintf("LeanCloud API error %d: %s", err.Code, err.Content)
}

// Upload upload specific file to LeanCloud
func Upload(appID string, appKey string, filename string, serverURL string, reader io.Reader, contentType string) (*File, error) {
	opts := &grequests.RequestOptions{
		Headers: map[string]string{
			"X-LC-Id":      appID,
			"X-LC-Key":     appKey,
			"Content-Type": contentType,
		},
		UserAgent:   "LeanCloud-Go-Upload/" + version,
		RequestBody: reader,
	}

	resp, err := grequests.Post(serverURL+"/1.1/files/"+filename, opts)
	if err != nil {
		return nil, err
	}
	if !resp.Ok {
		return nil, newErrorFromBody(resp.Bytes())
	}

	result := new(File)
	err = resp.JSON(result)
	if result.URL == "" {
		return nil, errors.New("Upload file failed")
	}
	return result, err
}

func newErrorFromBody(body []byte) error {
	var err Error
	err2 := json.Unmarshal([]byte(body), &err)
	if err2 != nil {
		return errors.New("Upload failed")
	}
	return err
}

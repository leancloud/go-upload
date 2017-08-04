package upload

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var client = &http.Client{
	Timeout: 14*time.Second + 1*time.Second,
}

// Upload upload specific file to LeanCloud
func Upload(name string, mimeType string, reader io.ReadSeeker, opts *Options) (*File, error) {
	// if opts.serverURL() == "https://api.leancloud.cn" || opts.serverURL() == "https://leancloud.cn" {
	// 	size, err := reader.Seek(0, io.SeekEnd)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	tokens, err := getFileTokens(name, mimeType, size, opts)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	putRet := new(qiniu.PutRet)
	// 	err = qiniu.Put(nil, putRet, tokens.Token, tokens.Key, reader, &qiniu.PutExtra{
	// 		MimeType: mimeType,
	// 	})
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	file := &File{
	// 		ObjectID: tokens.ObjectID,
	// 		URL:      tokens.URL,
	// 	}
	// 	return file, nil
	// }

	url := opts.serverURL() + "/1.1/files/" + name
	println(url)
	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		return nil, err
	}

	request.Header.Set("X-LC-Id", opts.AppID)
	request.Header.Set("X-LC-Key", opts.AppKey)
	request.Header.Set("User-Agent", "LeanCloud-Go-Upload/"+version)

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 201 {
		return nil, newErrorFromBody(body)
	}

	result := new(File)
	err = json.Unmarshal(body, result)
	if result.URL == "" {
		return nil, errors.New("Upload file failed")
	}
	return result, err
}

// UploadFileVerbose will open an file and upload it
func UploadFileVerbose(name string, mimeType string, path string, opts *Options) (*File, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return Upload(name, mimeType, f, opts)
}

// UploadFile will open an file and upload it. the file name and mime type is autodetected
func UploadFile(path string, opts *Options) (*File, error) {
	_, name := filepath.Split(path)
	mimeType := mime.TypeByExtension(filepath.Ext(path))
	return UploadFileVerbose(name, mimeType, path, opts)
}

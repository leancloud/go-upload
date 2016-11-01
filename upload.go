package upload

import (
	"errors"
	qiniu "github.com/qiniu/api.v6/io"
	"io"

	"github.com/levigross/grequests"
)

// Upload upload specific file to LeanCloud
func Upload(name string, mime string, reader io.Reader, opts *Options) (*File, error) {
	if opts.serverURL() == "https://api.leancloud.cn" {
		tokens, err := getFileTokens(name, mime, opts)
		if err != nil {
			return nil, err
		}
		putRet := new(qiniu.PutRet)
		err = qiniu.Put(nil, putRet, tokens.Token, tokens.Key, reader, &qiniu.PutExtra{
			MimeType: mime,
		})
		if err != nil {
			return nil, err
		}
		file := &File{
			ObjectID: tokens.ObjectID,
			URL:      tokens.URL,
		}
		return file, nil
	}

	reqOpts := &grequests.RequestOptions{
		Headers: map[string]string{
			"X-LC-Id":  opts.AppID,
			"X-LC-Key": opts.AppKey,
		},
		UserAgent:   "LeanCloud-Go-Upload/" + version,
		RequestBody: reader,
	}

	resp, err := grequests.Post(opts.serverURL()+"/1.1/files/"+name, reqOpts)
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

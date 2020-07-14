package requests

import (
	"bytes"
	"errors"
	"net/http"
	"reflect"
)

func buildBodyRequest(wf *Temporary) *http.Request {
	var req *http.Request
	var err error
	contentType := ""

	if wf.Body.GetIOBody() == nil {
		req, err = http.NewRequest(wf.Method, wf.GetRawURL(), nil)
	} else {
		var bodybuf *bytes.Buffer
		switch wf.Body.GetIOBody().(type) {
		case []byte:
			bodybuf = bytes.NewBuffer(wf.Body.GetIOBody().([]byte))
		case string:
			bodybuf = bytes.NewBuffer([]byte(wf.Body.GetIOBody().(string)))
		case *bytes.Buffer:
			bodybuf = bytes.NewBuffer(wf.Body.GetIOBody().(*bytes.Buffer).Bytes())
		default:
			panic(errors.New("the type is not exist, type is " + reflect.TypeOf(wf.Body.GetIOBody()).String()))
		}
		req, err = http.NewRequest(wf.Method, wf.GetRawURL(), bodybuf)
	}

	if err != nil {
		panic(err)
	}

	if wf.Body.ContentType() != "" {
		contentType = wf.Body.ContentType()
	} else {
		contentType = ""
		if wf.Method == "POST" || wf.Method == "PUT" || wf.Method == "PATCH" {
			contentType = TypeURLENCODED
		}
	}

	if contentType != "" {
		req.Header.Set(HeaderKeyContentType, contentType)
	}

	return req

}

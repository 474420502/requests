package requests

import (
	"net/http"
)

func buildBodyRequest(tp *Temporary) *http.Request {
	var req *http.Request
	var err error

	// contentType := ""
	// if err != nil {
	// 	panic(err)
	// }

	// if tp.mwriter != nil {
	// 	err = tp.mwriter.mwriter.Close()
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	defer func() { tp.mwriter = nil }()
	// 	tp.Body.SetPrefix(tp.mwriter.mwriter.FormDataContentType())
	// }

	if tp.Body == nil {
		req, err = http.NewRequest(tp.Method, tp.GetRawURL(), nil)
		if err != nil {
			panic(err)
		}
	} else {

		req, err = http.NewRequest(tp.Method, tp.GetRawURL(), tp.Body)
		if err != nil {
			panic(err)
		}
		// req.ContentLength = int64(bodybuf.Len())
	}

	// if tp.Body.ContentType() != "" {
	// 	contentType = tp.Body.ContentType()
	// } else {
	// 	contentType = ""
	// 	if tp.Method == "POST" || tp.Method == "PUT" || tp.Method == "PATCH" {
	// 		contentType = TypeURLENCODED
	// 	}
	// }

	// if contentType != "" {
	// 	req.Header.Set(HeaderKeyContentType, contentType)
	// }

	return req

}

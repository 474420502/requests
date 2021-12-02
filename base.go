package requests

import (
	"net/http"
)

func buildBodyRequest(tp *Temporary) *http.Request {
	var req *http.Request
	// var err error

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

	// if tp.Body.GetIOBody() == nil {
	// 	req, err = http.NewRequest(tp.Method, tp.GetRawURL(), nil)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// } else {
	// 	var bodybuf *bytes.Buffer
	// 	switch tp.Body.GetIOBody().(type) {
	// 	case []byte:
	// 		bodybuf = bytes.NewBuffer(tp.Body.GetIOBody().([]byte))
	// 	case string:
	// 		bodybuf = bytes.NewBuffer([]byte(tp.Body.GetIOBody().(string)))
	// 	case *bytes.Buffer:
	// 		bodybuf = bytes.NewBuffer(tp.Body.GetIOBody().(*bytes.Buffer).Bytes())
	// 	default:
	// 		panic(errors.New("the type is not exist, type is " + reflect.TypeOf(tp.Body.GetIOBody()).String()))
	// 	}
	// 	req, err = http.NewRequest(tp.Method, tp.GetRawURL(), bodybuf)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	// req.ContentLength = int64(bodybuf.Len())
	// }

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

package requests

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"errors"
	"net/http"

	"github.com/andybalholm/brotli"
)

type M map[string]interface{}

func buildBodyRequest(tp *Temporary) (*http.Request, error) {
	var req *http.Request
	var err error

	var cts [5]bool = [5]bool{false, false, false, false, false}
	for _, typ := range tp.session.acceptEncoding {
		cts[typ] = true
	}

	for _, typ := range tp.acceptEncoding {
		cts[typ] = true
	}

	if cts[AcceptEncodingGzip] {
		tp.Header.Add("Accept-Encoding", "gzip")
	}

	if cts[AcceptEncodingDeflate] {
		tp.Header.Add("Accept-Encoding", "deflate")
	}

	if cts[AcceptEncodingBr] {
		tp.Header.Add("Accept-Encoding", "br")
	}

	if tp.Body == nil {
		req, err = http.NewRequest(tp.Method, tp.GetRawURL(), nil)
		if err != nil {
			return req, err
		}
	} else {
		var buf = bytes.NewBuffer(nil)
		var ct ContentEncodingType
		if tp.contentEncoding != ContentEncodingNoCompress {
			ct = tp.contentEncoding
		} else if tp.session.contentEncoding != ContentEncodingNoCompress {
			ct = tp.session.contentEncoding
		}
		switch ct {
		case ContentEncodingNoCompress:
			_, err := buf.Write(tp.Body.Bytes())
			if err != nil {
				return nil, err
			}
		case ContentEncodingGzip:
			tp.Header.Add("Content-Encoding", "gzip")
			w := gzip.NewWriter(buf)
			w.Write(tp.Body.Bytes())
			err = w.Close()
			if err != nil {
				panic(err)
			}
		case ContentEncodingDeflate:
			tp.Header.Add("Content-Encoding", "deflate")
			w, err := flate.NewWriter(buf, flate.DefaultCompression)
			if err != nil {
				return nil, err
			}
			w.Write(tp.Body.Bytes())
			err = w.Close()
			if err != nil {
				return nil, err
			}
		case ContentEncodingBr:
			tp.Header.Add("Content-Encoding", "br")
			w := brotli.NewWriter(buf)
			w.Write(tp.Body.Bytes())
			err = w.Close()
			if err != nil {
				return nil, err
			}
		default:
			return nil, errors.New("compress type not support")
		}

		req, err = http.NewRequest(tp.Method, tp.GetRawURL(), buf)
		if err != nil {
			return req, err
		}
	}

	return req, nil
}

// Head 请求
func Head(url string) *Temporary {
	return NewSession().Head(url)
}

// Get 请求
func Get(url string) *Temporary {
	return NewSession().Get(url)
}

// Post 请求
func Post(url string) *Temporary {
	return NewSession().Post(url)
}

// Put 请求
func Put(url string) *Temporary {
	return NewSession().Put(url)
}

// Patch 请求
func Patch(url string) *Temporary {
	return NewSession().Patch(url)
}

// Delete 请求
func Delete(url string) *Temporary {
	return NewSession().Delete(url)
}

// Connect 请求
func Connect(url string) *Temporary {
	return NewSession().Connect(url)
}

// Options 请求
func Options(url string) *Temporary {
	return NewSession().Options(url)
}

// Trace 请求
func Trace(url string) *Temporary {
	return NewSession().Trace(url)
}

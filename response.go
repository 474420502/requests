package requests

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/andybalholm/brotli"
	"github.com/tidwall/gjson"
)

// Response Response from Execute()
type Response struct {
	readBytes    []byte
	readResponse *http.Response
}

// FromHTTPResponse Response . isDecompressNoAccept auto Decompress. like python requests
func FromHTTPResponse(resp *http.Response, isDecompressNoAccept bool) (*Response, error) {
	var err error
	var rbuf []byte

	defer resp.Body.Close()

	ContentEncoding := resp.Header.Get("Content-Encoding")
	switch ContentEncoding {
	case "gzip":
		r, err := gzip.NewReader(resp.Body)
		if err != nil {
			return nil, err
		}
		rbuf, err = ioutil.ReadAll(r)
		if err != nil {
			return nil, err
		}
	case "deflate":
		r := flate.NewReader(resp.Body)
		rbuf, err = ioutil.ReadAll(r)
		if err != nil {
			return nil, err
		}
	case "br":
		r := brotli.NewReader(resp.Body)
		rbuf, err = ioutil.ReadAll(r)
		if err != nil {
			return nil, err
		}
	default:
		srcbuf, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		if isDecompressNoAccept { // 在某个已经遗忘的网页测试过, 为了兼容 Python requests
			srcReader := bytes.NewReader(srcbuf)
			var reader io.ReadCloser
			if reader, err = gzip.NewReader(srcReader); err == nil {
				defer reader.Close()
				rbuf, err = ioutil.ReadAll(reader)
				if err != nil {
					return nil, err
				}

			} else if reader, err = zlib.NewReader(srcReader); err == nil {
				defer reader.Close()
				rbuf, err = ioutil.ReadAll(reader)
				if err != nil {
					return nil, err
				}
			} else {
				rbuf = srcbuf
			}

		} else {
			rbuf = srcbuf
		}
	}

	return &Response{readBytes: rbuf, readResponse: resp}, nil
}

// ContentString return string(Content())
func (gresp *Response) ContentString() string {
	return string(gresp.readBytes)
}

// Content return Response Bytes
func (gresp *Response) Content() []byte {
	return gresp.readBytes
}

// GetResponse  get golang http.Response
func (gresp *Response) GetResponse() *http.Response {
	return gresp.readResponse
}

// GetStatus get Statue String
func (gresp *Response) GetStatus() string {
	return gresp.readResponse.Status
}

// GetStatusCode  get Statue int
func (gresp *Response) GetStatusCode() int {
	return gresp.readResponse.StatusCode
}

// GetHeader Header map[string][]string
func (gresp *Response) GetHeader() http.Header {
	return gresp.readResponse.Header
}

// GetCookie get response cookies
func (gresp *Response) GetCookie() []*http.Cookie {
	return gresp.readResponse.Cookies()
}

// GetContentLength get Content length, if exists IsDecompressNoAccept, data is the length that compressed.
func (gresp *Response) GetContentLength() int64 {
	return gresp.readResponse.ContentLength
}

// Json  return gjson.Parse(jsonBody)
func (gresp *Response) Json() gjson.Result {
	return gjson.ParseBytes(gresp.readBytes)
}

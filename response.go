package requests

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"io"
	"io/ioutil"
	"net/http"
)

// Response 响应内容包含http.Response 已读
type Response struct {
	readBytes    []byte
	readResponse *http.Response
}

// FromHTTPResponse 生成Response 从标准http.Response
func FromHTTPResponse(resp *http.Response, isDecompressNoAccept bool) (*Response, error) {
	var err error
	var rbuf []byte

	// 复制response 返回内容 并且测试是否有解压的需求
	srcbuf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()

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

	return &Response{readBytes: rbuf, readResponse: resp}, nil
}

// Content 返回解压后的内容
func (gresp *Response) Content() string {
	return string(gresp.readBytes)
}

// ContentBytes 返回解压后的内容Bytes
func (gresp *Response) ContentBytes() []byte {
	return gresp.readBytes
}

// GetSrcResponse  获取原生golang http.Response
func (gresp *Response) GetSrcResponse() *http.Response {
	return gresp.readResponse
}

// GetStatue 获取Statue String
func (gresp *Response) GetStatue() string {
	return gresp.readResponse.Status
}

// GetStatueCode 获取Statue int
func (gresp *Response) GetStatueCode() int {
	return gresp.readResponse.StatusCode
}

// GetHeader Header map[string][]string
func (gresp *Response) GetHeader() http.Header {
	return gresp.readResponse.Header
}

// GetContentLength 获取Content的内容长度, 如果存在 IsDecompressNoAccept 可能是压缩级别的长度, 非GetContent长度
func (gresp *Response) GetContentLength() int64 {
	return gresp.readResponse.ContentLength
}

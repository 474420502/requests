package requests

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

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

// DecodeJSON 将响应体反序列化到给定的结构体中
func (gresp *Response) DecodeJSON(v interface{}) error {
	return json.Unmarshal(gresp.readBytes, v)
}

// BindJSON DecodeJSON的别名，更符合现代Go API命名习惯
func (gresp *Response) BindJSON(v interface{}) error {
	return gresp.DecodeJSON(v)
}

// IsJSON 检查响应是否为JSON类型
func (gresp *Response) IsJSON() bool {
	contentType := gresp.readResponse.Header.Get("Content-Type")
	return strings.Contains(contentType, "application/json") || strings.Contains(contentType, "text/json")
}

// GetJSONField 获取JSON中的特定字段值（使用gjson）
func (gresp *Response) GetJSONField(path string) gjson.Result {
	return gjson.GetBytes(gresp.readBytes, path)
}

// GetJSONString 获取JSON中的字符串字段
func (gresp *Response) GetJSONString(path string) (string, error) {
	result := gjson.GetBytes(gresp.readBytes, path)
	if !result.Exists() {
		return "", fmt.Errorf("JSON field '%s' does not exist", path)
	}
	return result.String(), nil
}

// GetJSONInt 获取JSON中的整数字段
func (gresp *Response) GetJSONInt(path string) (int64, error) {
	result := gjson.GetBytes(gresp.readBytes, path)
	if !result.Exists() {
		return 0, fmt.Errorf("JSON field '%s' does not exist", path)
	}
	if result.Type != gjson.Number {
		return 0, fmt.Errorf("JSON field '%s' is not a number", path)
	}
	return result.Int(), nil
}

// GetJSONFloat 获取JSON中的浮点数字段
func (gresp *Response) GetJSONFloat(path string) (float64, error) {
	result := gjson.GetBytes(gresp.readBytes, path)
	if !result.Exists() {
		return 0, fmt.Errorf("JSON field '%s' does not exist", path)
	}
	if result.Type != gjson.Number {
		return 0, fmt.Errorf("JSON field '%s' is not a number", path)
	}
	return result.Float(), nil
}

// GetJSONBool 获取JSON中的布尔字段
func (gresp *Response) GetJSONBool(path string) (bool, error) {
	result := gjson.GetBytes(gresp.readBytes, path)
	if !result.Exists() {
		return false, fmt.Errorf("JSON field '%s' does not exist", path)
	}
	if result.Type != gjson.True && result.Type != gjson.False {
		return false, fmt.Errorf("JSON field '%s' is not a boolean", path)
	}
	return result.Bool(), nil
}

// MustBindJSON 绑定JSON，如果失败则panic（用于必须成功的场景）
//
// Deprecated: 此方法将在 v3.0.0 中移除。推荐使用 BindJSON 方法并适当处理错误。
// panic模式违反了Go的错误处理最佳实践。
//
// 迁移示例:
//
//	旧: response.MustBindJSON(&data)
//	新: if err := response.BindJSON(&data); err != nil {
//	      return fmt.Errorf("failed to bind JSON: %w", err)
//	    }
func (gresp *Response) MustBindJSON(v interface{}) {
	if err := gresp.BindJSON(v); err != nil {
		panic(fmt.Sprintf("MustBindJSON failed: %v", err))
	}
}

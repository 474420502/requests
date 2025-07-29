package requests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/url"
	"reflect"
	"strconv"
)

func writeFormUploadFile(mwriter *multipart.Writer, ufile *UploadFile) error {
	part, err := mwriter.CreateFormFile(ufile.FieldName, ufile.FileName)
	if err != nil {
		return err
	}
	_, err = io.Copy(part, ufile.FileReader)
	return err
}

// // *multipart.Writer 需要 Close()
// func createMultipart(params ...interface{}) (*bytes.Buffer, *multipart.Writer) {
// 	plen := len(params)

// 	body := &bytes.Buffer{}
// 	mwriter := multipart.NewWriter(body)

// 	for i, iparam := range params[0 : plen-1] {
// 		switch param := iparam.(type) {
// 		case *UploadFile:
// 			if param.FieldName == "" {
// 				param.FieldName = fmt.Sprintf("file%d", i)
// 			}
// 			writeFormUploadFile(mwriter, param)
// 		case UploadFile:
// 			if param.FieldName == "" {
// 				param.FieldName = fmt.Sprintf("file%d", i)
// 			}
// 			writeFormUploadFile(mwriter, &param)
// 		case []*UploadFile:
// 			for i, p := range param {
// 				if p.FieldName == "" {
// 					p.FieldName = "file" + strconv.Itoa(i)
// 				}
// 				writeFormUploadFile(mwriter, p)
// 			}
// 		case []UploadFile:
// 			for i, p := range param {
// 				if p.FieldName == "" {
// 					p.FieldName = "file" + strconv.Itoa(i)
// 				}
// 				writeFormUploadFile(mwriter, &p)
// 			}
// 		case string:
// 			uploadFiles, err := UploadFileFromGlob(param)
// 			if err != nil {
// 				log.Println(err)
// 			} else {
// 				for i, p := range uploadFiles {
// 					if p.FieldName == "" {
// 						p.FieldName = "file" + strconv.Itoa(i)
// 					}
// 					writeFormUploadFile(mwriter, p)
// 				}
// 			}

// 		case []string:
// 			for i, glob := range param {
// 				uploadFiles, err := UploadFileFromGlob(glob)
// 				if err != nil {
// 					log.Println(err)
// 				} else {
// 					for ii, p := range uploadFiles {
// 						if p.FieldName == "" {
// 							p.FieldName = "file" + strconv.Itoa(i) + "_" + strconv.Itoa(ii)
// 						}
// 						writeFormUploadFile(mwriter, p)
// 					}
// 				}
// 			}
// 		case map[string]string:
// 			for k, v := range param {
// 				mwriter.WriteField(k, v)
// 			}
// 		case map[string][]string:
// 			for k, vs := range param {
// 				for _, v := range vs {
// 					mwriter.WriteField(k, v)
// 				}
// 			}
// 		case url.Values:
// 			for k, vs := range param {
// 				for _, v := range vs {
// 					mwriter.WriteField(k, v)
// 				}
// 			}
// 		case map[string]interface{}:
// 			for k, v := range param {
// 				data, err := json.Marshal(v)
// 				if err != nil {
// 					log.Println(err)
// 				} else {
// 					mwriter.WriteField(k, string(data))
// 				}
// 			}
// 		default:
// 			if reflect.TypeOf(param).ConvertibleTo(compatibleType) {
// 				cparam := reflect.ValueOf(param).Convert(compatibleType)
// 				for k, v := range cparam.Interface().(map[string]interface{}) {
// 					switch cv := v.(type) {
// 					case string:
// 						mwriter.WriteField(k, cv)
// 					case []byte:
// 						mwriter.WriteField(k, string(cv))
// 					case []rune:
// 						mwriter.WriteField(k, string(cv))
// 					default:
// 						data, err := json.Marshal(v)
// 						if err != nil {
// 							log.Println(err)
// 						} else {
// 							mwriter.WriteField(k, string(data))
// 						}
// 					}
// 				}
// 			}
// 		}
// 	}

// 	// postParams.AddContentType("boundary=" + b)

// 	err := mwriter.Close()
// 	if err != nil {
// 		panic(err)
// 	}

// 	// log.Println(string(body.Bytes()))
// 	return body, mwriter
// }

// createMultipartExSafe *multipart.Writer 需要 Close() - 安全版本，返回错误而不是panic
func createMultipartExSafe(params ...interface{}) (*bytes.Buffer, *multipart.Writer, error) {
	body := &bytes.Buffer{}
	mwriter := multipart.NewWriter(body)

	for i, iparam := range params {
		switch param := iparam.(type) {
		case *UploadFile:
			if param.FieldName == "" {
				param.FieldName = fmt.Sprintf("file%d", i)
			}
			if err := writeFormUploadFile(mwriter, param); err != nil {
				return nil, nil, fmt.Errorf("failed to write upload file: %w", err)
			}
		case UploadFile:
			if param.FieldName == "" {
				param.FieldName = fmt.Sprintf("file%d", i)
			}
			if err := writeFormUploadFile(mwriter, &param); err != nil {
				return nil, nil, fmt.Errorf("failed to write upload file: %w", err)
			}
		case []*UploadFile:
			for i, p := range param {
				if p.FieldName == "" {
					p.FieldName = "file" + strconv.Itoa(i)
				}
				if err := writeFormUploadFile(mwriter, p); err != nil {
					return nil, nil, fmt.Errorf("failed to write upload file: %w", err)
				}
			}
		case []UploadFile:
			for i, p := range param {
				if p.FieldName == "" {
					p.FieldName = "file" + strconv.Itoa(i)
				}
				if err := writeFormUploadFile(mwriter, &p); err != nil {
					return nil, nil, fmt.Errorf("failed to write upload file: %w", err)
				}
			}
		case string:
			uploadFiles, err := UploadFileFromGlob(param)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to process glob pattern: %w", err)
			}
			for i, p := range uploadFiles {
				if p.FieldName == "" {
					p.FieldName = "file" + strconv.Itoa(i)
				}
				if err := writeFormUploadFile(mwriter, p); err != nil {
					return nil, nil, fmt.Errorf("failed to write upload file: %w", err)
				}
			}

		case []string:
			for i, glob := range param {
				uploadFiles, err := UploadFileFromGlob(glob)
				if err != nil {
					return nil, nil, fmt.Errorf("failed to process glob pattern: %w", err)
				}
				for ii, p := range uploadFiles {
					if p.FieldName == "" {
						p.FieldName = "file" + strconv.Itoa(i) + "_" + strconv.Itoa(ii)
					}
					if err := writeFormUploadFile(mwriter, p); err != nil {
						return nil, nil, fmt.Errorf("failed to write upload file: %w", err)
					}
				}
			}
		case map[string]string:
			for k, v := range param {
				if err := mwriter.WriteField(k, v); err != nil {
					return nil, nil, fmt.Errorf("failed to write form field: %w", err)
				}
			}
		case map[string][]string:
			for k, vs := range param {
				for _, v := range vs {
					if err := mwriter.WriteField(k, v); err != nil {
						return nil, nil, fmt.Errorf("failed to write form field: %w", err)
					}
				}
			}
		case url.Values:
			for k, vs := range param {
				for _, v := range vs {
					if err := mwriter.WriteField(k, v); err != nil {
						return nil, nil, fmt.Errorf("failed to write form field: %w", err)
					}
				}
			}
		case map[string]interface{}:
			for k, v := range param {
				switch v.(type) {
				case map[string]interface{}, []interface{}, []map[string]interface{}:
					data, err := json.Marshal(v)
					if err != nil {
						return nil, nil, fmt.Errorf("failed to marshal JSON: %w", err)
					}
					if err := mwriter.WriteField(k, string(data)); err != nil {
						return nil, nil, fmt.Errorf("failed to write form field: %w", err)
					}
				default:
					// TODO: 处理json的基础类型到 WriteField 要求都转字符串
					var str string
					switch t := v.(type) {
					case int:
						str = strconv.Itoa(t)
					case float64:
						str = strconv.FormatFloat(t, 'f', -1, 64)
					case bool:
						str = strconv.FormatBool(t)
					case string:
						str = t
					default:
						str = fmt.Sprintf("%v", t)
					}
					if err := mwriter.WriteField(k, str); err != nil {
						return nil, nil, fmt.Errorf("failed to write form field: %w", err)
					}
				}
			}
		case *multipart.Writer, multipart.Writer:
			return nil, nil, fmt.Errorf("only accept single (*)multipart.Writer")
		default:
			if reflect.TypeOf(param).ConvertibleTo(compatibleType) {
				cparam := reflect.ValueOf(param).Convert(compatibleType)
				for k, v := range cparam.Interface().(map[string]interface{}) {
					switch cv := v.(type) {
					case string:
						if err := mwriter.WriteField(k, cv); err != nil {
							return nil, nil, fmt.Errorf("failed to write form field: %w", err)
						}
					case []byte:
						if err := mwriter.WriteField(k, string(cv)); err != nil {
							return nil, nil, fmt.Errorf("failed to write form field: %w", err)
						}
					case []rune:
						if err := mwriter.WriteField(k, string(cv)); err != nil {
							return nil, nil, fmt.Errorf("failed to write form field: %w", err)
						}
					default:
						data, err := json.Marshal(v)
						if err != nil {
							return nil, nil, fmt.Errorf("failed to marshal JSON: %w", err)
						}
						if err := mwriter.WriteField(k, string(data)); err != nil {
							return nil, nil, fmt.Errorf("failed to write form field: %w", err)
						}
					}
				}
			}
		}
	}

	err := mwriter.Close()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	return body, mwriter, nil
}

// *multipart.Writer 需要 Close()
func createMultipartEx(params ...interface{}) (*bytes.Buffer, *multipart.Writer) {

	body := &bytes.Buffer{}
	mwriter := multipart.NewWriter(body)

	for i, iparam := range params {
		switch param := iparam.(type) {
		case *UploadFile:
			if param.FieldName == "" {
				param.FieldName = fmt.Sprintf("file%d", i)
			}
			if err := writeFormUploadFile(mwriter, param); err != nil {
				log.Println(err)
			}
		case UploadFile:
			if param.FieldName == "" {
				param.FieldName = fmt.Sprintf("file%d", i)
			}
			if err := writeFormUploadFile(mwriter, &param); err != nil {
				log.Println(err)
			}
		case []*UploadFile:
			for i, p := range param {
				if p.FieldName == "" {
					p.FieldName = "file" + strconv.Itoa(i)
				}
				if err := writeFormUploadFile(mwriter, p); err != nil {
					log.Println(err)
				}
			}
		case []UploadFile:
			for i, p := range param {
				if p.FieldName == "" {
					p.FieldName = "file" + strconv.Itoa(i)
				}
				if err := writeFormUploadFile(mwriter, &p); err != nil {
					log.Println(err)
				}
			}
		case string:
			uploadFiles, err := UploadFileFromGlob(param)
			if err != nil {
				log.Println(err)
			} else {
				for i, p := range uploadFiles {
					if p.FieldName == "" {
						p.FieldName = "file" + strconv.Itoa(i)
					}
					if err := writeFormUploadFile(mwriter, p); err != nil {
						log.Println(err)
					}
				}
			}

		case []string:
			for i, glob := range param {
				uploadFiles, err := UploadFileFromGlob(glob)
				if err != nil {
					log.Println(err)
				} else {
					for ii, p := range uploadFiles {
						if p.FieldName == "" {
							p.FieldName = "file" + strconv.Itoa(i) + "_" + strconv.Itoa(ii)
						}
						if err := writeFormUploadFile(mwriter, p); err != nil {
							log.Println(err)
						}
					}
				}
			}
		case map[string]string:
			for k, v := range param {
				mwriter.WriteField(k, v)
			}
		case map[string][]string:
			for k, vs := range param {
				for _, v := range vs {
					mwriter.WriteField(k, v)
				}
			}
		case url.Values:
			for k, vs := range param {
				for _, v := range vs {
					mwriter.WriteField(k, v)
				}
			}
		case map[string]interface{}:
			for k, v := range param {
				switch v.(type) {
				case map[string]interface{}, []interface{}, []map[string]interface{}:
					data, err := json.Marshal(v)
					if err != nil {
						log.Println(err)
					} else {
						mwriter.WriteField(k, string(data))
					}
				default:
					// TODO: 处理json的基础类型到 WriteField 要求都转字符串
					var str string
					switch t := v.(type) {
					case int:
						str = strconv.Itoa(t)
					case float64:
						str = strconv.FormatFloat(t, 'f', -1, 64)
					case bool:
						str = strconv.FormatBool(t)
					case string:
						str = t
					default:
						str = fmt.Sprintf("%v", t)
					}
					mwriter.WriteField(k, str)
				}
			}
		case *multipart.Writer, multipart.Writer:
			panic("only accept single (*)multipart.Writer")
		default:
			if reflect.TypeOf(param).ConvertibleTo(compatibleType) {
				cparam := reflect.ValueOf(param).Convert(compatibleType)
				for k, v := range cparam.Interface().(map[string]interface{}) {
					switch cv := v.(type) {
					case string:
						mwriter.WriteField(k, cv)
					case []byte:
						mwriter.WriteField(k, string(cv))
					case []rune:
						mwriter.WriteField(k, string(cv))
					default:
						data, err := json.Marshal(v)
						if err != nil {
							log.Println(err)
						} else {
							mwriter.WriteField(k, string(data))
						}
					}
				}
			}
		}
	}

	// postParams.AddContentType("boundary=" + b)

	err := mwriter.Close()
	if err != nil {
		panic(err)
	}

	// log.Println(string(body.Bytes()))
	return body, mwriter
}

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

func writeFormUploadFile(mwriter *multipart.Writer, ufile *UploadFile) {
	part, err := mwriter.CreateFormFile(ufile.FieldName, ufile.FileName)
	if err != nil {
		log.Panic(err)
	}
	io.Copy(part, ufile.FileReader)
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
			writeFormUploadFile(mwriter, param)
		case UploadFile:
			if param.FieldName == "" {
				param.FieldName = fmt.Sprintf("file%d", i)
			}
			writeFormUploadFile(mwriter, &param)
		case []*UploadFile:
			for i, p := range param {
				if p.FieldName == "" {
					p.FieldName = "file" + strconv.Itoa(i)
				}
				writeFormUploadFile(mwriter, p)
			}
		case []UploadFile:
			for i, p := range param {
				if p.FieldName == "" {
					p.FieldName = "file" + strconv.Itoa(i)
				}
				writeFormUploadFile(mwriter, &p)
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
					writeFormUploadFile(mwriter, p)
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
						writeFormUploadFile(mwriter, p)
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

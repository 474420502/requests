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

// // MultipartWriter Only Write data. Execute() will with multipart data
// type MultipartWriter struct {
// 	fileindex int
// 	mwriter   *multipart.Writer
// }

// // SetBoundary overrides the Writer's default randomly-generated boundary separator with an explicit value.
// // SetBoundary must be called before any parts are created, may only contain certain ASCII characters, and must be non-empty and at most 70 bytes long.
// func (mw *MultipartWriter) SetBoundary(boundary string) error {
// 	return mw.mwriter.SetBoundary(boundary)
// }

// // Boundary returns the Writer's boundary.
// func (mw *MultipartWriter) Boundary() string {
// 	return mw.mwriter.Boundary()
// }

// // AddField write name value with boundary
// func (mw *MultipartWriter) AddField(name, value string) error {
// 	return mw.mwriter.WriteField(name, value)
// }

// // AddFile write name value with boundary
// func (mw *MultipartWriter) AddFile(filename string, dataReader io.Reader) error {
// 	fn := fmt.Sprintf("file%d", mw.fileindex)
// 	w, err := mw.mwriter.CreateFormFile(fn, filename)
// 	if err != nil {
// 		return err
// 	}
// 	mw.fileindex++
// 	_, err = io.Copy(w, dataReader)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// // AddFileEx write name value with boundary, fieldname
// func (mw *MultipartWriter) AddFileEx(fieldname string, filename string, dataReader io.Reader) error {
// 	w, err := mw.mwriter.CreateFormFile(fieldname, filename)
// 	if err != nil {
// 		return err
// 	}
// 	_, err = io.Copy(w, dataReader)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func writeFormUploadFile(mwriter *multipart.Writer, ufile *UploadFile) {
	part, err := mwriter.CreateFormFile(ufile.FieldName, ufile.FileName)
	if err != nil {
		log.Panic(err)
	}
	io.Copy(part, ufile.FileReaderCloser)

	err = ufile.FileReaderCloser.Close()
	if err != nil {
		panic(err)
	}
}

func createMultipart(params ...interface{}) (*bytes.Buffer, *multipart.Writer) {
	plen := len(params)

	body := &bytes.Buffer{}
	mwriter := multipart.NewWriter(body)

	for i, iparam := range params[0 : plen-1] {
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
				data, err := json.Marshal(v)
				if err != nil {
					log.Println(err)
				} else {
					mwriter.WriteField(k, string(data))
				}
			}
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

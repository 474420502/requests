package requests

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/url"
	"strconv"
)

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

func createMultipart(postParams IBody, params []interface{}) {
	plen := len(params)

	body := &bytes.Buffer{}
	mwriter := multipart.NewWriter(body)

	for _, iparam := range params[0 : plen-1] {
		switch param := iparam.(type) {
		case *UploadFile:
			if param.FieldName == "" {
				param.FieldName = "file0"
			}
			writeFormUploadFile(mwriter, param)
		case UploadFile:
			if param.FieldName == "" {
				param.FieldName = "file0"
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
		}
	}

	postParams.AddContentType(mwriter.FormDataContentType())
	postParams.SetIOBody(body)

	err := mwriter.Close()
	if err != nil {
		panic(err)
	}
}

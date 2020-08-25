package requests

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
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

func postFile(filename string, target_url string) *http.Request {
	bodybuf := bytes.NewBufferString("")
	bodywriter := multipart.NewWriter(bodybuf)

	// use the body_writer to write the Part headers to the buffer
	_, err := bodywriter.CreateFormFile("userfile", filename)
	if err != nil {
		fmt.Println("error writing to buffer")
		return nil
	}

	// the file data will be the second part of the body
	fh, err := os.Open(filename)
	if err != nil {
		fmt.Println("error opening file")
		return nil
	}
	// need to know the boundary to properly close the part myself.
	boundary := bodywriter.Boundary()
	//close_string := fmt.Sprintf("\r\n--%s--\r\n", boundary)
	closebuf := bytes.NewBufferString(fmt.Sprintf("\r\n--%s--\r\n", boundary))

	// use multi-reader to defer the reading of the file data until
	// writing to the socket buffer.
	requestreader := io.MultiReader(bodybuf, fh, closebuf)
	fi, err := fh.Stat()
	if err != nil {
		fmt.Printf("Error Stating file: %s", filename)
		return nil
	}
	req, err := http.NewRequest("POST", target_url, requestreader)
	if err != nil {
		return nil
	}

	// Set headers for multipart, and Content Length
	req.Header.Add("Content-Type", "multipart/form-data; boundary="+boundary)
	req.ContentLength = fi.Size() + int64(bodybuf.Len()) + int64(closebuf.Len())

	return req
}

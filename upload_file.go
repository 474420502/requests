package requests

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

// UploadFile 上传文件的结构
type UploadFile struct {
	FileName   string
	FieldName  string
	FileReader io.Reader
	fileCloser io.Closer // 关闭文件
}

// SetFileName 设置FileName属性
func (ufile *UploadFile) SetFileName(filename string) {
	ufile.FileName = filename
}

// GetFileName 设置FileName属性
func (ufile *UploadFile) GetFileName() string {
	return ufile.FileName
}

// SetFile 设置FileName属性
func (ufile *UploadFile) SetFile(reader io.Reader) {
	ufile.FileReader = reader
	if ufile.fileCloser != nil {
		err := ufile.fileCloser.Close()
		if err != nil {
			panic(err)
		}
	}
	ufile.fileCloser = nil
}

// SetFileFromPath 设置FileName属性
func (ufile *UploadFile) SetFileFromPath(filename string) error {
	fd, err := os.Open(filename)
	if err != nil {
		return err
	}
	ufile.SetFile(fd)
	ufile.fileCloser = fd
	return nil
}

// GetFile 设置FileName属性
func (ufile *UploadFile) GetFile() io.Reader {
	return ufile.FileReader
}

// SetFieldName 设置FileName属性
func (ufile *UploadFile) SetFieldName(fieldname string) {
	ufile.FieldName = fieldname
}

// GetFieldName 设置FileName属性
func (ufile *UploadFile) GetFieldName() string {
	return ufile.FieldName
}

// NewUploadFile 创建一个空的UploadFile, 必须设置 FileName FieldName FileReader  三个属性
func NewUploadFile() *UploadFile {
	ufile := &UploadFile{}
	runtime.SetFinalizer(ufile, func(ufile *UploadFile) {
		if ufile.fileCloser != nil {
			err := ufile.fileCloser.Close()
			if err != nil {
				log.Println(err)
			}
		}
	})
	return ufile
}

// UploadFileFromPath 从本地文件获取上传文件
func UploadFileFromPath(fileName string) (*UploadFile, error) {
	fd, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	ufile := NewUploadFile()
	ufile.FileReader = fd
	ufile.FileName = fileName
	ufile.fileCloser = fd

	return ufile, nil

}

// UploadFileFromGlob 根据Glob从本地文件获取上传文件
func UploadFileFromGlob(glob string) ([]*UploadFile, error) {
	files, err := filepath.Glob(glob)

	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("UploadFileFromGlob: len(files) == 0")
	}

	var ufiles []*UploadFile

	for _, f := range files {
		if s, err := os.Stat(f); err != nil || s.IsDir() {
			continue
		}

		fd, err := os.Open(f)
		if err != nil {
			// log.Println(fd.Name(), err)
			return nil, fmt.Errorf("%s error: %s", fd.Name(), err)
		} else {

			ufile := NewUploadFile()
			ufile.FileReader = fd
			ufile.FileName = filepath.Base(fd.Name())
			ufile.fileCloser = fd

			ufiles = append(ufiles, ufile)
		}
	}

	return ufiles, nil

}

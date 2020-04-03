package requests

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

// UploadFile 上传文件的结构
type UploadFile struct {
	FileName         string
	FieldName        string
	FileReaderCloser io.ReadCloser
}

// SetFileName 设置FileName属性
func (ufile *UploadFile) SetFileName(filename string) {
	ufile.FileName = filename
}

// GetFileName 设置FileName属性
func (ufile *UploadFile) GetFileName() string {
	return ufile.FileName
}

// SetFileReaderCloser 设置FileName属性
func (ufile *UploadFile) SetFileReaderCloser(readerCloser io.ReadCloser) {
	ufile.FileReaderCloser = readerCloser
}

// SetFileReaderCloserFromFile 设置FileName属性
func (ufile *UploadFile) SetFileReaderCloserFromFile(filename string) error {
	fd, err := os.Open(filename)
	if err != nil {
		return err
	}
	ufile.SetFileReaderCloser(fd)
	return nil
}

// GetFileReaderCloser 设置FileName属性
func (ufile *UploadFile) GetFileReaderCloser() io.ReadCloser {
	return ufile.FileReaderCloser
}

// SetFieldName 设置FileName属性
func (ufile *UploadFile) SetFieldName(fieldname string) {
	ufile.FieldName = fieldname
}

// GetFieldName 设置FileName属性
func (ufile *UploadFile) GetFieldName() string {
	return ufile.FieldName
}

// NewUploadFile 创建一个空的UploadFile, 必须设置 FileName FieldName FileReaderCloser 三个属性
func NewUploadFile() *UploadFile {
	return &UploadFile{}
}

// UploadFileFromPath 从本地文件获取上传文件
func UploadFileFromPath(fileName string) (*UploadFile, error) {
	fd, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	return &UploadFile{FileReaderCloser: fd, FileName: fileName}, nil
}

// UploadFileFromGlob 根据Glob从本地文件获取上传文件
func UploadFileFromGlob(glob string) ([]*UploadFile, error) {
	files, err := filepath.Glob(glob)

	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		log.Println("UploadFileFromGlob: len(files) == 0")
	}

	var ufiles []*UploadFile

	for _, f := range files {
		if s, err := os.Stat(f); err != nil || s.IsDir() {
			continue
		}

		fd, err := os.Open(f)
		if err != nil {
			log.Println(fd.Name(), err)
		} else {
			ufiles = append(ufiles, &UploadFile{FileReaderCloser: fd, FileName: filepath.Base(fd.Name())})
		}
	}

	return ufiles, nil

}

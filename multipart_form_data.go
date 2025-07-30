package requests

import (
	"bytes"
	"io"
	"mime/multipart"
)

// MultipartFormData 用于构建multipart表单数据
//
// Deprecated: 此结构体及其方法将在 v3.0.0 中移除。
// 推荐使用以下现代API：
//   - request.SetFormFields(map[string]string) - 设置表单字段
//   - request.AddFormFile(fieldName, fileName, content) - 添加文件
//   - request.SetFormFileFromPath(fieldName, filePath) - 从路径添加文件
//
// 迁移示例:
//
//	旧代码:
//	  mpfd := &MultipartFormData{}
//	  mpfd.AddField("name", "value")
//	  mpfd.AddFile("file", "test.txt", []byte("content"))
//	  request.SetBody(mpfd.Data())
//
//	新代码:
//	  request.SetFormFields(map[string]string{"name": "value"})
//	  request.AddFormFile("file", "test.txt", []byte("content"))
type MultipartFormData struct {
	data   bytes.Buffer
	writer *multipart.Writer
}

// Data 返回构建的multipart数据
func (mpfd *MultipartFormData) Data() *bytes.Buffer {
	return &mpfd.data
}

// Writer 返回multipart writer，用于手动添加字段
func (mpfd *MultipartFormData) Writer() *multipart.Writer {
	return mpfd.writer
}

// AddField 添加表单字段
func (mpfd *MultipartFormData) AddField(name, value string) error {
	return mpfd.writer.WriteField(name, value)
}

// AddFile 添加文件字段
func (mpfd *MultipartFormData) AddFile(fieldName, fileName string, content []byte) error {
	part, err := mpfd.writer.CreateFormFile(fieldName, fileName)
	if err != nil {
		return err
	}
	_, err = part.Write(content)
	return err
}

// AddFieldFile 添加文件字段（支持io.Reader）
func (mpfd *MultipartFormData) AddFieldFile(fieldName, fileName string, reader io.Reader) error {
	part, err := mpfd.writer.CreateFormFile(fieldName, fileName)
	if err != nil {
		return err
	}
	_, err = io.Copy(part, reader)
	return err
}

// Close 关闭writer，必须在使用完毕后调用
func (mpfd *MultipartFormData) Close() error {
	return mpfd.writer.Close()
}

// ContentType 返回Content-Type头的值
func (mpfd *MultipartFormData) ContentType() string {
	return mpfd.writer.FormDataContentType()
}

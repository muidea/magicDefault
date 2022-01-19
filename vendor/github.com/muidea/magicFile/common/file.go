package common

import (
	commonDef "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicFile/model"
)

const FileModule = "/module/file"

// FileView 文件视图
type FileView struct {
	ID          int    `json:"id"`
	Token       string `json:"token"`
	Name        string `json:"name"`
	Description string `json:"description"`
	UploadDate  int64  `json:"uploadDate"`
	Validity    int64  `json:"validity"`
	ReserveFlag int    `json:"reserveFlag"`
}

// FromFileDetail from fileDetail
func (s *FileView) FromFileDetail(detail *model.FileDetail) {
	s.ID = detail.ID
	s.Token = detail.Token
	s.Name = detail.Name
	s.Description = detail.Description
	s.UploadDate = detail.UploadDate
	s.Validity = detail.Validity
	s.ReserveFlag = detail.ReserveFlag
}

// FileParam 文件信息
type FileParam struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Validity    int64    `json:"validity"`
	ReserveFlag int      `json:"reserveFlag"`
	Tags        []string `json:"tags"`
}

// UploadFileResult upload file result
type UploadFileResult struct {
	commonDef.Result
	File *model.FileDetail `json:"file"`
}

// ViewFileResult view file result
type ViewFileResult UploadFileResult

// UpdateFileResult update file result
type UpdateFileResult UploadFileResult

// DeleteFileResult delete file result
type DeleteFileResult UploadFileResult

// QueryFileResult query file result
type QueryFileResult UploadFileResult

// QueryFilesResult query files result
type QueryFilesResult struct {
	commonDef.Result
	Total int64               `json:"total"`
	Files []*model.FileDetail `json:"file"`
}

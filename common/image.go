package common

import (
	commonDef "github.com/muidea/magicCommon/def"
	fileCommon "github.com/muidea/magicFile/common"
)

const (
	FilterImage = "/image/query/"
	QueryImage  = "/image/query/:id"
	UpdateImage = "/image/update/:id"
	DeleteImage = "/image/destroy/:id"
)

const ImageModule = "/module/image"

type ImageResult struct {
	commonDef.Result
	Image *fileCommon.FileView `json:"image"`
}

type ImageListResult struct {
	commonDef.Result
	Total int64                  `json:"total"`
	Image []*fileCommon.FileView `json:"image"`
}

type ImageStatisticResult struct {
	ImageListResult
}

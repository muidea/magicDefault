package model

// FileDetail 文件信息
// @ID ID
// @Token 访问Token
// @Name 文件名
// @Path 文件访问路径
// @UploadData 上传日期
// @Validity 有效期
// @ReserveFlag 预留标识(共享、私有)
// @Source 来源
type FileDetail struct {
	ID          int      `json:"id" orm:"id key auto"`
	Token       string   `json:"token" orm:"token"`
	Name        string   `json:"name" orm:"name"`
	Path        string   `json:"path" orm:"path"`
	Description string   `json:"description" orm:"description"`
	UploadDate  int64    `json:"uploadDate" orm:"uploadDate"`
	Validity    int64    `json:"validity" orm:"validity"`
	ReserveFlag int      `json:"reserveFlag" orm:"reserveFlag"`
	Tags        []string `json:"tags" orm:"tags"`
	Source      string   `json:"source" orm:"source"`
}

package model

// Article 文章
type Article struct {
	ID         int        `json:"id" orm:"id key auto"`
	Title      string     `json:"title" orm:"title"`
	Content    string     `json:"content" orm:"content"`
	Catalog    []*Catalog `json:"catalog" orm:"catalog"`
	Creater    int        `json:"creater" orm:"creater"`
	CreateTime int64      `json:"createTime" orm:"createTime"`
	Namespace  string     `json:"namespace" orm:"namespace"`
}

// Catalog 分类详细信息
type Catalog struct {
	ID          int      `json:"id" orm:"id key auto"`
	Name        string   `json:"name" orm:"name"`
	Description string   `json:"description" orm:"description"`
	Catalog     *Catalog `json:"catalog" orm:"catalog"`
	Creater     int      `json:"creater" orm:"creater"`
	CreateTime  int64    `json:"createTime" orm:"createTime"`
	Namespace   string   `json:"namespace" orm:"namespace"`
}

// Link 链接
type Link struct {
	ID          int        `json:"id" orm:"id key auto"`
	Name        string     `json:"name" orm:"name"`
	Description string     `json:"description" orm:"description"`
	URL         string     `json:"url" orm:"url"`
	Logo        string     `json:"logo" orm:"logo"`
	Catalog     []*Catalog `json:"catalog" orm:"catalog"`
	Creater     int        `json:"creater" orm:"creater"`
	CreateTime  int64      `json:"createTime" orm:"createTime"`
	Namespace   string     `json:"namespace" orm:"namespace"`
}

// Media 文件信息
type Media struct {
	ID          int        `json:"id" orm:"id key auto"`
	Name        string     `json:"name" orm:"name"`
	Description string     `json:"description" orm:"description"`
	FileToken   string     `json:"fileToken" orm:"fileToken"`
	Expiration  int        `json:"expiration" orm:"expiration"`
	Tags        []string   `json:"tags" orm:"tags"`
	Catalog     []*Catalog `json:"catalog" orm:"catalog"`
	Creater     int        `json:"creater" orm:"creater"`
	CreateTime  int64      `json:"createTime" orm:"createTime"`
	Namespace   string     `json:"namespace" orm:"namespace"`
}

// Comment 注释
type Comment struct {
	ID         int    `json:"id" orm:"id key auto"`
	Content    string `json:"content" orm:"content"`
	Flag       int    `json:"flag" orm:"flag"`
	Creater    int    `json:"creater" orm:"creater"`
	CreateTime int64  `json:"createTime" orm:"createTime"`
	Namespace  string `json:"namespace" orm:"namespace"`
}

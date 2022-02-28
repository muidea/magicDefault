package model

// Shelf 货架
// @Code 编号
// @Description 描述
// @Used 已用用容量
// @Capacity 设定容量
// @Warehouse 所在仓库
// @Creater 创建人
// @CreateTime 创建时间
// @Namespace 所属空间
type Shelf struct {
	ID          int        `json:"id" orm:"id key auto"`
	Code        string     `json:"code" orm:"code"`
	Description string     `json:"description" orm:"description"`
	Used        int        `json:"used" orm:"used"`
	Capacity    int        `json:"capacity" orm:"capacity" validate:"required"`
	Warehouse   *Warehouse `json:"warehouse" orm:"warehouse" validate:"required"`
	Creater     int        `json:"creater" orm:"creater"`
	CreateTime  int64      `json:"createTime" orm:"createTime"`
	Namespace   string     `json:"namespace" orm:"namespace"`
}

// Warehouse 仓库
// @Code 编号
// @Name 名称
// @Description 描述
// @Creater 创建人
// @CreateTime 创建时间
// @Namespace 所属空间
type Warehouse struct {
	ID          int    `json:"id" orm:"id key auto"`
	Code        string `json:"code" orm:"code"`
	Name        string `json:"name" orm:"name" validate:"required"`
	Description string `json:"description" orm:"description"`
	Creater     int    `json:"creater" orm:"creater"`
	CreateTime  int64  `json:"createTime" orm:"createTime"`
	Namespace   string `json:"namespace" orm:"namespace"`
}

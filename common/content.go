package common

import (
	cc "github.com/muidea/magicCas/common"
	cd "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicDefault/model"
)

const ContentModule = "/module/content"

const (
	FilterContentArticle = "/content/article/query/"
	QueryContentArticle  = "/content/article/query/:id"
	CreateContentArticle = "/content/article/create/"
	UpdateContentArticle = "/content/article/update/:id"
	DeleteContentArticle = "/content/article/delete/:id"
	NotifyContentArticle = "/content/article/notify/:id"

	FilterContentCatalog = "/content/catalog/query/"
	QueryContentCatalog  = "/content/catalog/query/:id"
	CreateContentCatalog = "/content/catalog/create/"
	UpdateContentCatalog = "/content/catalog/update/:id"
	DeleteContentCatalog = "/content/catalog/delete/:id"
	NotifyContentCatalog = "/content/catalog/notify/:id"

	FilterContentComment = "/content/comment/query/"
	QueryContentComment  = "/content/comment/query/:id"
	CreateContentComment = "/content/comment/create/"
	UpdateContentComment = "/content/comment/update/:id"
	DeleteContentComment = "/content/comment/delete/:id"
	NotifyContentComment = "/content/comment/notify/:id"

	FilterContentLink = "/content/link/query/"
	QueryContentLink  = "/content/link/query/:id"
	CreateContentLink = "/content/link/create/"
	UpdateContentLink = "/content/link/update/:id"
	DeleteContentLink = "/content/link/delete/:id"
	NotifyContentLink = "/content/link/notify/:id"

	FilterContentMedia = "/content/media/query/"
	QueryContentMedia  = "/content/media/query/:id"
	CreateContentMedia = "/content/media/create/"
	UpdateContentMedia = "/content/media/update/:id"
	DeleteContentMedia = "/content/media/delete/:id"
	NotifyContentMedia = "/content/media/notify/:id"
)

type ArticleView struct {
	ID         int            `json:"id"`
	Title      string         `json:"title"`
	Content    string         `json:"content"`
	Catalog    []*CatalogLite `json:"catalog"`
	Creater    *cc.EntityView `json:"creater"`
	UpdateTime int64          `json:"updateTime"`
}

func (s *ArticleView) FromArticle(ptr *model.Article, entityPtr *cc.EntityView) {
	if ptr == nil {
		return
	}
	s.ID = ptr.ID
	s.Title = ptr.Title
	s.Content = ptr.Content
	for _, val := range ptr.Catalog {
		lite := &CatalogLite{}
		lite.FromCatalog(val)
		s.Catalog = append(s.Catalog, lite)
	}
	s.Creater = entityPtr
	s.UpdateTime = ptr.UpdateTime
	return
}

type ArticleLite struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (s *ArticleLite) FromArticle(ptr *model.Article) {
	if ptr == nil {
		return
	}
	s.ID = ptr.ID
	s.Title = ptr.Title
	s.Content = ptr.Content
	return
}

func (s *ArticleLite) FromView(ptr *ArticleView) {
	if ptr == nil {
		return
	}
	s.ID = ptr.ID
	s.Title = ptr.Title
	s.Content = ptr.Content
	return
}

func (s *ArticleLite) ToArticle() (ret *model.Article) {
	ptr := &model.Article{Catalog: []*model.Catalog{}}
	if s.ID != 0 {
		ptr.ID = s.ID
	}
	if s.Title != "" {
		ptr.Title = s.Title
	}
	if s.Content != "" {
		ptr.Content = s.Content
	}
	ret = ptr
	return
}

type ArticleParam struct {
	ID      int            `json:"id"`
	Title   string         `json:"title"`
	Content string         `json:"content"`
	Catalog []*CatalogLite `json:"catalog"`
}

func (s *ArticleParam) ToArticle(ptr *model.Article) (ret *model.Article) {
	if ptr == nil {
		ptr = &model.Article{Catalog: []*model.Catalog{}}
	}

	if s.Title != "" {
		ptr.Title = s.Title
	}
	if s.Content != "" {
		ptr.Content = s.Content
	}
	for _, val := range s.Catalog {
		ptr.Catalog = append(ptr.Catalog, val.ToCatalog())
	}

	ret = ptr
	return
}

type ArticleResult struct {
	cd.Result
	Article *ArticleView `json:"article"`
}

type ArticleLiteListResult struct {
	cd.Result
	Total   int64          `json:"total"`
	Article []*ArticleLite `json:"article"`
}

type ArticleListResult struct {
	cd.Result
	Total   int64          `json:"total"`
	Article []*ArticleView `json:"article"`
}

type ArticleStatisticResult struct {
	ArticleListResult
	Summary []*TotalizerView `json:"summary"`
}

type CatalogView struct {
	ID          int            `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Catalog     *CatalogLite   `json:"catalog"`
	Creater     *cc.EntityView `json:"creater"`
	UpdateTime  int64          `json:"updateTime"`
}

func (s *CatalogView) FromCatalog(ptr *model.Catalog, entityPtr *cc.EntityView) {
	if ptr == nil {
		return
	}
	s.ID = ptr.ID
	s.Name = ptr.Name
	s.Description = ptr.Description
	s.Catalog = &CatalogLite{}
	s.Catalog.FromCatalog(ptr.Catalog)
	s.Creater = entityPtr
	s.UpdateTime = ptr.UpdateTime
	return
}

type CatalogLite struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (s *CatalogLite) FromCatalog(ptr *model.Catalog) {
	if ptr == nil {
		return
	}
	s.ID = ptr.ID
	s.Name = ptr.Name
	s.Description = ptr.Description
	return
}

func (s *CatalogLite) FromView(ptr *CatalogView) {
	if ptr == nil {
		return
	}
	s.ID = ptr.ID
	s.Name = ptr.Name
	s.Description = ptr.Description
	return
}

func (s *CatalogLite) ToCatalog() (ret *model.Catalog) {
	ptr := &model.Catalog{Catalog: &model.Catalog{}}
	if s.ID != 0 {
		ptr.ID = s.ID
	}
	if s.Name != "" {
		ptr.Name = s.Name
	}
	if s.Description != "" {
		ptr.Description = s.Description
	}
	ret = ptr
	return
}

type CatalogParam struct {
	ID          int          `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Catalog     *CatalogLite `json:"catalog"`
}

func (s *CatalogParam) ToCatalog(ptr *model.Catalog) (ret *model.Catalog) {
	if ptr == nil {
		ptr = &model.Catalog{Catalog: &model.Catalog{}}
	}

	if s.Name != "" {
		ptr.Name = s.Name
	}
	if s.Description != "" {
		ptr.Description = s.Description
	}
	if s.Catalog != nil {
		ptr.Catalog = s.Catalog.ToCatalog()
	}

	ret = ptr
	return
}

type CatalogResult struct {
	cd.Result
	Catalog *CatalogView `json:"catalog"`
}

type CatalogLiteListResult struct {
	cd.Result
	Total   int64          `json:"total"`
	Catalog []*CatalogLite `json:"catalog"`
}

type CatalogListResult struct {
	cd.Result
	Total   int64          `json:"total"`
	Catalog []*CatalogView `json:"catalog"`
}

type CatalogStatisticResult struct {
	CatalogListResult
	Summary []*TotalizerView `json:"summary"`
}

type CommentView struct {
	ID         int            `json:"id"`
	Content    string         `json:"content"`
	Flag       int            `json:"flag"`
	Creater    *cc.EntityView `json:"creater"`
	UpdateTime int64          `json:"updateTime"`
}

func (s *CommentView) FromComment(ptr *model.Comment, entityPtr *cc.EntityView) {
	if ptr == nil {
		return
	}
	s.ID = ptr.ID
	s.Content = ptr.Content
	s.Flag = ptr.Flag
	s.Creater = entityPtr
	s.UpdateTime = ptr.UpdateTime
	return
}

type CommentLite struct {
	ID      int    `json:"id"`
	Content string `json:"content"`
	Flag    int    `json:"flag"`
}

func (s *CommentLite) FromComment(ptr *model.Comment) {
	if ptr == nil {
		return
	}
	s.ID = ptr.ID
	s.Content = ptr.Content
	s.Flag = ptr.Flag
	return
}

func (s *CommentLite) FromView(ptr *CommentView) {
	if ptr == nil {
		return
	}
	s.ID = ptr.ID
	s.Content = ptr.Content
	s.Flag = ptr.Flag
	return
}

func (s *CommentLite) ToComment() (ret *model.Comment) {
	ptr := &model.Comment{}
	if s.ID != 0 {
		ptr.ID = s.ID
	}
	if s.Content != "" {
		ptr.Content = s.Content
	}
	if s.Flag != 0 {
		ptr.Flag = s.Flag
	}
	ret = ptr
	return
}

type CommentParam struct {
	ID      int    `json:"id"`
	Content string `json:"content"`
	Flag    int    `json:"flag"`
}

func (s *CommentParam) ToComment(ptr *model.Comment) (ret *model.Comment) {
	if ptr == nil {
		ptr = &model.Comment{}
	}

	if s.Content != "" {
		ptr.Content = s.Content
	}
	if s.Flag != 0 {
		ptr.Flag = s.Flag
	}

	ret = ptr
	return
}

type CommentResult struct {
	cd.Result
	Comment *CommentView `json:"comment"`
}

type CommentLiteListResult struct {
	cd.Result
	Total   int64          `json:"total"`
	Comment []*CommentLite `json:"comment"`
}

type CommentListResult struct {
	cd.Result
	Total   int64          `json:"total"`
	Comment []*CommentView `json:"comment"`
}

type CommentStatisticResult struct {
	CommentListResult
	Summary []*TotalizerView `json:"summary"`
}

type LinkView struct {
	ID          int            `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	URL         string         `json:"url"`
	Logo        string         `json:"logo"`
	Catalog     []*CatalogLite `json:"catalog"`
	Creater     *cc.EntityView `json:"creater"`
	UpdateTime  int64          `json:"updateTime"`
}

func (s *LinkView) FromLink(ptr *model.Link, entityPtr *cc.EntityView) {
	if ptr == nil {
		return
	}
	s.ID = ptr.ID
	s.Name = ptr.Name
	s.Description = ptr.Description
	s.URL = ptr.URL
	s.Logo = ptr.Logo
	for _, val := range ptr.Catalog {
		lite := &CatalogLite{}
		lite.FromCatalog(val)
		s.Catalog = append(s.Catalog, lite)
	}
	s.Creater = entityPtr
	s.UpdateTime = ptr.UpdateTime
	return
}

type LinkLite struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (s *LinkLite) FromLink(ptr *model.Link) {
	if ptr == nil {
		return
	}
	s.ID = ptr.ID
	s.Name = ptr.Name
	s.Description = ptr.Description
	return
}

func (s *LinkLite) FromView(ptr *LinkView) {
	if ptr == nil {
		return
	}
	s.ID = ptr.ID
	s.Name = ptr.Name
	s.Description = ptr.Description
	return
}

func (s *LinkLite) ToLink() (ret *model.Link) {
	ptr := &model.Link{Catalog: []*model.Catalog{}}
	if s.ID != 0 {
		ptr.ID = s.ID
	}
	if s.Name != "" {
		ptr.Name = s.Name
	}
	if s.Description != "" {
		ptr.Description = s.Description
	}
	ret = ptr
	return
}

type LinkParam struct {
	ID          int            `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	URL         string         `json:"url"`
	Logo        string         `json:"logo"`
	Catalog     []*CatalogLite `json:"catalog"`
}

func (s *LinkParam) ToLink(ptr *model.Link) (ret *model.Link) {
	if ptr == nil {
		ptr = &model.Link{Catalog: []*model.Catalog{}}
	}

	if s.Name != "" {
		ptr.Name = s.Name
	}
	if s.Description != "" {
		ptr.Description = s.Description
	}
	if s.URL != "" {
		ptr.URL = s.URL
	}
	if s.Logo != "" {
		ptr.Logo = s.Logo
	}
	for _, val := range s.Catalog {
		ptr.Catalog = append(ptr.Catalog, val.ToCatalog())
	}

	ret = ptr
	return
}

type LinkResult struct {
	cd.Result
	Link *LinkView `json:"link"`
}

type LinkLiteListResult struct {
	cd.Result
	Total int64       `json:"total"`
	Link  []*LinkLite `json:"link"`
}

type LinkListResult struct {
	cd.Result
	Total int64       `json:"total"`
	Link  []*LinkView `json:"link"`
}

type LinkStatisticResult struct {
	LinkListResult
	Summary []*TotalizerView `json:"summary"`
}

type MediaView struct {
	ID          int            `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	FileToken   string         `json:"fileToken"`
	Expiration  int            `json:"expiration"`
	Tags        []string       `json:"tags"`
	Catalog     []*CatalogLite `json:"catalog"`
	Creater     *cc.EntityView `json:"creater"`
	UpdateTime  int64          `json:"updateTime"`
}

func (s *MediaView) FromMedia(ptr *model.Media, entityPtr *cc.EntityView) {
	if ptr == nil {
		return
	}
	s.ID = ptr.ID
	s.Name = ptr.Name
	s.Description = ptr.Description
	s.FileToken = ptr.FileToken
	s.Expiration = ptr.Expiration
	s.Tags = ptr.Tags
	for _, val := range ptr.Catalog {
		lite := &CatalogLite{}
		lite.FromCatalog(val)
		s.Catalog = append(s.Catalog, lite)
	}
	s.Creater = entityPtr
	s.UpdateTime = ptr.UpdateTime
	return
}

type MediaLite struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (s *MediaLite) FromMedia(ptr *model.Media) {
	if ptr == nil {
		return
	}
	s.ID = ptr.ID
	s.Name = ptr.Name
	s.Description = ptr.Description
	return
}

func (s *MediaLite) FromView(ptr *MediaView) {
	if ptr == nil {
		return
	}
	s.ID = ptr.ID
	s.Name = ptr.Name
	s.Description = ptr.Description
	return
}

func (s *MediaLite) ToMedia() (ret *model.Media) {
	ptr := &model.Media{Catalog: []*model.Catalog{}}
	if s.ID != 0 {
		ptr.ID = s.ID
	}
	if s.Name != "" {
		ptr.Name = s.Name
	}
	if s.Description != "" {
		ptr.Description = s.Description
	}
	ret = ptr
	return
}

type MediaParam struct {
	ID          int            `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	FileToken   string         `json:"fileToken"`
	Expiration  int            `json:"expiration"`
	Tags        []string       `json:"tags"`
	Catalog     []*CatalogLite `json:"catalog"`
}

func (s *MediaParam) ToMedia(ptr *model.Media) (ret *model.Media) {
	if ptr == nil {
		ptr = &model.Media{Catalog: []*model.Catalog{}}
	}

	if s.Name != "" {
		ptr.Name = s.Name
	}
	if s.Description != "" {
		ptr.Description = s.Description
	}
	if s.FileToken != "" {
		ptr.FileToken = s.FileToken
	}
	if s.Expiration != 0 {
		ptr.Expiration = s.Expiration
	}
	if len(s.Tags) > 0 {
		ptr.Tags = s.Tags
	}
	for _, val := range s.Catalog {
		ptr.Catalog = append(ptr.Catalog, val.ToCatalog())
	}

	ret = ptr
	return
}

type MediaResult struct {
	cd.Result
	Media *MediaView `json:"media"`
}

type MediaLiteListResult struct {
	cd.Result
	Total int64        `json:"total"`
	Media []*MediaLite `json:"media"`
}

type MediaListResult struct {
	cd.Result
	Total int64        `json:"total"`
	Media []*MediaView `json:"media"`
}

type MediaStatisticResult struct {
	MediaListResult
	Summary []*TotalizerView `json:"summary"`
}

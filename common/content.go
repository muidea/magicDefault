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

	FilterContentComment = "/content/comment/query/"
	QueryContentComment  = "/content/comment/query/:id"
	CreateContentComment = "/content/comment/create/"
	UpdateContentComment = "/content/comment/update/:id"
	DeleteContentComment = "/content/comment/delete/:id"
	NotifyContentComment = "/content/comment/notify/:id"
)

type ArticleView struct {
	// TODO
}

func (s *ArticleView) FromArticle(ptr *model.Article, entityPtr *cc.EntityView) {
	// TODO
	return
}

type ArticleLite struct {
	// TODO
}

func (s *ArticleLite) FromArticle(ptr *model.Article) {
	// TODO
	return
}

func (s *ArticleLite) FromView(ptr *ArticleView) {
	// TODO
	return
}

func (s *ArticleLite) ToArticle() (ret *model.Article) {
	// TODO
	return
}

type ArticleParam struct {
	// TODO
}

func (s *ArticleParam) ToArticle() (ret *model.Article) {
	// TODO
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
	// TODO
}

func (s *CatalogView) FromCatalog(ptr *model.Catalog, entityPtr *cc.EntityView) {
	// TODO
	return
}

type CatalogLite struct {
	// TODO
}

func (s *CatalogLite) FromCatalog(ptr *model.Catalog) {
	// TODO
	return
}

func (s *CatalogLite) FromView(ptr *CatalogView) {
	// TODO
	return
}

func (s *CatalogLite) ToCatalog() (ret *model.Catalog) {
	// TODO
	return
}

type CatalogParam struct {
	// TODO
}

func (s *CatalogParam) ToCatalog() (ret *model.Catalog) {
	// TODO
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

type LinkView struct {
	// TODO
}

func (s *LinkView) FromLink(ptr *model.Link, entityPtr *cc.EntityView) {
	// TODO
	return
}

type LinkLite struct {
	// TODO
}

func (s *LinkLite) FromLink(ptr *model.Link) {
	// TODO
	return
}

func (s *LinkLite) FromView(ptr *LinkView) {
	// TODO
	return
}

func (s *LinkLite) ToLink() (ret *model.Link) {
	// TODO
	return
}

type LinkParam struct {
	// TODO
}

func (s *LinkParam) ToLink() (ret *model.Link) {
	// TODO
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
	// TODO
}

func (s *MediaView) FromMedia(ptr *model.Media, entityPtr *cc.EntityView) {
	// TODO
	return
}

type MediaLite struct {
	// TODO
}

func (s *MediaLite) FromMedia(ptr *model.Media) {
	// TODO
	return
}

func (s *MediaLite) FromView(ptr *MediaView) {
	// TODO
	return
}

func (s *MediaLite) ToMedia() (ret *model.Media) {
	// TODO
	return
}

type MediaParam struct {
	// TODO
}

func (s *MediaParam) ToMedia() (ret *model.Media) {
	// TODO
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

type CommentView struct {
	// TODO
}

func (s *CommentView) FromComment(ptr *model.Comment, entityPtr *cc.EntityView) {
	// TODO
	return
}

type CommentLite struct {
	// TODO
}

func (s *CommentLite) FromComment(ptr *model.Comment) {
	// TODO
	return
}

func (s *CommentLite) FromView(ptr *CommentView) {
	// TODO
	return
}

func (s *CommentLite) ToComment() (ret *model.Comment) {
	// TODO
	return
}

type CommentParam struct {
	// TODO
}

func (s *CommentParam) ToComment() (ret *model.Comment) {
	// TODO
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

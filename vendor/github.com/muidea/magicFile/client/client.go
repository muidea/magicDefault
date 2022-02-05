package client

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	commonDef "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/foundation/net"
	fu "github.com/muidea/magicCommon/foundation/util"
	"github.com/muidea/magicCommon/session"

	fileCommon "github.com/muidea/magicFile/common"
	fileModel "github.com/muidea/magicFile/model"
)

const fileItem = "fileItem"

// Client client interface
type Client interface {
	UploadFile(filePath string) (*fileModel.FileDetail, error)
	ViewFile(fileToken string) (*fileModel.FileDetail, error)
	DownloadFile(fileToken, filePath string) (string, error)
	UpdateFile(id int, param *fileCommon.FileParam) (*fileModel.FileDetail, error)
	DeleteFile(id int) (*fileModel.FileDetail, error)
	QueryFile(id int) (*fileModel.FileDetail, error)
	FilterFile(filter *fu.ContentFilter) ([]*fileModel.FileDetail, int64, error)

	AttachSource(source string)

	BindSession(sessionInfo *session.SessionInfo)
	UnBindSession()
	Release()
}

// NewClient new client
func NewClient(serverURL string) Client {
	clnt := &impl{serverURL: serverURL, httpClient: &http.Client{}}

	return clnt
}

type impl struct {
	serverURL   string
	source      string
	sessionInfo *session.SessionInfo
	httpClient  *http.Client
}

func (s *impl) UploadFile(filePath string) (*fileModel.FileDetail, error) {
	result := &fileCommon.UploadFileResult{}
	vals := url.Values{}
	vals.Set("key-name", fileItem)

	url, _ := url.ParseRequestURI(s.serverURL)
	url.Path = strings.Join([]string{url.Path, fileCommon.ApiVersion, fileCommon.UploadFileURL}, "")
	url.RawQuery = vals.Encode()
	err := net.HTTPUpload(s.httpClient, url.String(), fileItem, filePath, result, s.getContextValues())
	if err != nil {
		return nil, err
	}

	if result.ErrorCode != commonDef.Success {
		err = fmt.Errorf("upload file failed, reason:%s", result.Reason)
	}

	return result.File, err
}

func (s *impl) ViewFile(fileToken string) (*fileModel.FileDetail, error) {
	result := &fileCommon.ViewFileResult{}
	vals := url.Values{}
	vals.Set("fileToken", fileToken)

	url, _ := url.ParseRequestURI(s.serverURL)
	url.Path = strings.Join([]string{url.Path, fileCommon.ApiVersion, fileCommon.ViewFileURL}, "")
	url.RawQuery = vals.Encode()
	_, err := net.HTTPGet(s.httpClient, url.String(), result, s.getContextValues())
	if err != nil {
		return nil, err
	}

	if result.ErrorCode != commonDef.Success {
		err = fmt.Errorf("view file failed, reason:%s", result.Reason)
	}

	return result.File, err
}

func (s *impl) DownloadFile(fileToken, filePath string) (string, error) {
	vals := url.Values{}
	vals.Set("fileToken", fileToken)

	url, _ := url.ParseRequestURI(s.serverURL)
	url.Path = strings.Join([]string{url.Path, fileCommon.ApiVersion, fileCommon.DownloadFileURL}, "")
	url.RawQuery = vals.Encode()

	filePath, err := net.HTTPDownload(s.httpClient, url.String(), filePath, s.getContextValues())
	if err != nil {
		return "", err
	}

	return filePath, err
}

func (s *impl) UpdateFile(id int, param *fileCommon.FileParam) (*fileModel.FileDetail, error) {
	result := &fileCommon.UpdateFileResult{}
	vals := url.Values{}
	url, _ := url.ParseRequestURI(s.serverURL)
	url.Path = strings.Join([]string{url.Path, fileCommon.ApiVersion, fileCommon.UpdateFileURL}, "")
	url.Path = strings.ReplaceAll(url.Path, ":id", fmt.Sprintf("%d", id))
	url.RawQuery = vals.Encode()
	_, err := net.HTTPPut(s.httpClient, url.String(), param, result, s.getContextValues())
	if err != nil {
		return nil, err
	}

	if result.ErrorCode != commonDef.Success {
		err = fmt.Errorf("update file failed, reason:%s", result.Reason)
	}

	return result.File, err
}

func (s *impl) DeleteFile(id int) (*fileModel.FileDetail, error) {
	result := &fileCommon.DeleteFileResult{}
	vals := url.Values{}
	url, _ := url.ParseRequestURI(s.serverURL)
	url.Path = strings.Join([]string{url.Path, fileCommon.ApiVersion, fileCommon.DeleteFileURL}, "")
	url.Path = strings.ReplaceAll(url.Path, ":id", fmt.Sprintf("%d", id))
	url.RawQuery = vals.Encode()
	_, err := net.HTTPDelete(s.httpClient, url.String(), result, s.getContextValues())
	if err != nil {
		return nil, err
	}

	if result.ErrorCode != commonDef.Success {
		err = fmt.Errorf("delete file failed, reason:%s", result.Reason)
	}

	return result.File, err
}

func (s *impl) QueryFile(id int) (*fileModel.FileDetail, error) {
	result := &fileCommon.QueryFileResult{}
	vals := url.Values{}
	url, _ := url.ParseRequestURI(s.serverURL)
	url.Path = strings.Join([]string{url.Path, fileCommon.ApiVersion, fileCommon.QueryFileURL}, "")
	url.Path = strings.ReplaceAll(url.Path, ":id", fmt.Sprintf("%d", id))
	url.RawQuery = vals.Encode()
	_, err := net.HTTPGet(s.httpClient, url.String(), result, s.getContextValues())
	if err != nil {
		return nil, err
	}

	if result.ErrorCode != commonDef.Success {
		err = fmt.Errorf("query file failed, reason:%s", result.Reason)
	}

	return result.File, err
}

func (s *impl) FilterFile(filter *fu.ContentFilter) (ret []*fileModel.FileDetail, total int64, err error) {
	result := &fileCommon.QueryFilesResult{}
	vals := url.Values{}
	if filter != nil {
		vals = filter.Encode(vals)
	}

	url, _ := url.ParseRequestURI(s.serverURL)
	url.Path = strings.Join([]string{url.Path, fileCommon.ApiVersion, fileCommon.FilterFileURL}, "")
	url.RawQuery = vals.Encode()
	_, err = net.HTTPGet(s.httpClient, url.String(), result, s.getContextValues())
	if err != nil {
		return
	}

	if result.ErrorCode != commonDef.Success {
		err = fmt.Errorf("query all account failed, reason:%s", result.Reason)
		return
	}

	ret = result.Files
	total = result.Total
	return
}

func (s *impl) getContextValues() url.Values {
	ret := url.Values{}
	if s.source != "" {
		ret.Set("source", s.source)
	}
	if s.sessionInfo != nil {
		ret = s.sessionInfo.Encode(ret)
	}

	return ret
}

func (s *impl) AttachSource(source string) {
	s.source = source
}

func (s *impl) BindSession(sessionInfo *session.SessionInfo) {
	s.sessionInfo = sessionInfo
}

func (s *impl) UnBindSession() {
	s.sessionInfo = nil
}

func (s *impl) Release() {
	if s.httpClient != nil {
		s.httpClient.CloseIdleConnections()
		s.httpClient = nil
	}
}

package client

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/muidea/magicBatis/common"
	commonDef "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/foundation/net"

	"github.com/muidea/magicOrm/provider"
	"github.com/muidea/magicOrm/provider/remote"
)

// Client client interface
type Client interface {
	RegisterService(instancePtr *common.ServiceParam) (err error)
	UnregisterService() (err error)
	RegisterModel(entityPtr interface{}) (err error)
	UnregisterModel(entityPtr interface{}) (err error)
	CreateSchema(entityPtr interface{}) (err error)
	DropSchema(entityPtr interface{}) (err error)
	InsertEntity(entityPtr interface{}) (err error)
	DeleteEntity(entityPtr interface{}) (err error)
	UpdateEntity(entityPtr interface{}) (err error)
	QueryEntity(entityPtr interface{}) (err error)
	BatchQueryEntity(entitySlicePtr interface{}, filter *common.QueryFilter) (ret int64, err error)
	CountEntity(entityPtr interface{}, filter *common.QueryFilter) (ret int64, err error)

	Release()
}

// NewClient new client
func NewClient(serverURL, serviceName string) Client {
	clnt := &impl{serverURL: serverURL, serviceName: serviceName, httpClient: &http.Client{}}

	return clnt
}

type impl struct {
	serverURL   string
	serviceName string
	httpClient  *http.Client
}

func (s *impl) RegisterService(instancePtr *common.ServiceParam) (err error) {
	if instancePtr == nil {
		err = fmt.Errorf("illegal instance")
		return
	}

	result := &common.Result{}
	url := strings.Join([]string{s.serverURL, net.JoinPrefix(common.RegisterService, common.ApiVersion)}, "")
	if s.serviceName != "" {
		url = fmt.Sprintf("%s?source=%s", url, s.serviceName)
	}

	_, err = net.HTTPPost(s.httpClient, url, instancePtr, result)
	if err != nil {
		return
	}

	if result.ErrorCode != commonDef.Success {
		err = fmt.Errorf("register model failed, reason:%s", result.Reason)
	}

	return
}

func (s *impl) UnregisterService() (err error) {
	result := &common.Result{}
	url := strings.Join([]string{s.serverURL, net.JoinPrefix(common.UnregisterService, common.ApiVersion)}, "")
	if s.serviceName != "" {
		url = fmt.Sprintf("%s?source=%s", url, s.serviceName)
	}

	_, err = net.HTTPDelete(s.httpClient, url, result)
	if err != nil {
		return
	}

	if result.ErrorCode != commonDef.Success {
		err = fmt.Errorf("unregister model failed, reason:%s", result.Reason)
	}

	return
}

func (s *impl) RegisterModel(entityPtr interface{}) (err error) {
	entityObj, entityErr := remote.GetObject(entityPtr)
	if entityErr != nil {
		err = entityErr
		return
	}

	result := &common.Result{}
	url := strings.Join([]string{s.serverURL, net.JoinPrefix(common.RegisterModel, common.ApiVersion)}, "")
	if s.serviceName != "" {
		url = fmt.Sprintf("%s?source=%s", url, s.serviceName)
	}

	_, err = net.HTTPPost(s.httpClient, url, entityObj, result)
	if err != nil {
		return
	}

	if result.ErrorCode != commonDef.Success {
		err = fmt.Errorf("register model failed, reason:%s", result.Reason)
	}

	return
}

func (s *impl) UnregisterModel(entityPtr interface{}) (err error) {
	entityObj, entityErr := remote.GetObject(entityPtr)
	if entityErr != nil {
		err = entityErr
		return
	}

	result := &common.Result{}
	url := strings.Join([]string{s.serverURL, net.JoinPrefix(common.UnregisterModel, common.ApiVersion)}, "")
	if s.serviceName != "" {
		url = fmt.Sprintf("%s?source=%s", url, s.serviceName)
	}

	_, err = net.HTTPPost(s.httpClient, url, entityObj, result)
	if err != nil {
		return
	}

	if result.ErrorCode != commonDef.Success {
		err = fmt.Errorf("unregister model failed, reason:%s", result.Reason)
	}

	return
}

func (s *impl) CreateSchema(entityPtr interface{}) (err error) {
	entityObj, entityErr := remote.GetObject(entityPtr)
	if entityErr != nil {
		err = entityErr
		return
	}

	result := &common.Result{}
	url := strings.Join([]string{s.serverURL, net.JoinPrefix(common.CreateSchema, common.ApiVersion)}, "")
	if s.serviceName != "" {
		url = fmt.Sprintf("%s?source=%s", url, s.serviceName)
	}
	_, err = net.HTTPPost(s.httpClient, url, entityObj, result)
	if err != nil {
		return
	}

	if result.ErrorCode != commonDef.Success {
		err = fmt.Errorf("create schema failed, reason:%s", result.Reason)
	}

	return
}

func (s *impl) DropSchema(entityPtr interface{}) (err error) {
	entityObj, entityErr := remote.GetObject(entityPtr)
	if entityErr != nil {
		err = entityErr
		return
	}

	result := &common.Result{}
	url := strings.Join([]string{s.serverURL, net.JoinPrefix(common.DropSchema, common.ApiVersion)}, "")
	if s.serviceName != "" {
		url = fmt.Sprintf("%s?source=%s", url, s.serviceName)
	}
	_, err = net.HTTPPost(s.httpClient, url, entityObj, result)
	if err != nil {
		return
	}

	if result.ErrorCode != commonDef.Success {
		err = fmt.Errorf("drop schema failed, reason:%s", result.Reason)
	}

	return
}

func (s *impl) InsertEntity(entityPtr interface{}) (err error) {
	entityVal, entityErr := remote.GetObjectValue(entityPtr)
	if entityErr != nil {
		err = entityErr
		return
	}

	result := &common.ObjectValueResult{}
	url := strings.Join([]string{s.serverURL, net.JoinPrefix(common.InsertValue, common.ApiVersion)}, "")
	if s.serviceName != "" {
		url = fmt.Sprintf("%s?source=%s", url, s.serviceName)
	}
	_, err = net.HTTPPost(s.httpClient, url, entityVal, result)
	if err != nil {
		return
	}

	if result.ErrorCode != commonDef.Success {
		err = fmt.Errorf("insert value failed, reason:%s", result.Reason)
		return
	}

	val, valErr := remote.ConvertObjectValue(result.Value)
	if valErr != nil {
		err = valErr
		return
	}

	err = provider.UpdateEntity(val, entityPtr)
	return
}

func (s *impl) DeleteEntity(entityPtr interface{}) (err error) {
	entityVal, entityErr := remote.GetObjectValue(entityPtr)
	if entityErr != nil {
		err = entityErr
		return
	}

	result := &common.ObjectValueResult{}
	url := strings.Join([]string{s.serverURL, net.JoinPrefix(common.DeleteValue, common.ApiVersion)}, "")
	if s.serviceName != "" {
		url = fmt.Sprintf("%s?source=%s", url, s.serviceName)
	}
	_, err = net.HTTPPost(s.httpClient, url, entityVal, result)
	if err != nil {
		return
	}

	if result.ErrorCode != commonDef.Success {
		err = fmt.Errorf("delete value failed, reason:%s", result.Reason)
		return
	}

	return
}

func (s *impl) UpdateEntity(entityPtr interface{}) (err error) {
	entityVal, entityErr := remote.GetObjectValue(entityPtr)
	if entityErr != nil {
		err = entityErr
		return
	}

	result := &common.ObjectValueResult{}
	url := strings.Join([]string{s.serverURL, net.JoinPrefix(common.UpdateValue, common.ApiVersion)}, "")
	if s.serviceName != "" {
		url = fmt.Sprintf("%s?source=%s", url, s.serviceName)
	}
	_, err = net.HTTPPost(s.httpClient, url, entityVal, result)
	if err != nil {
		return
	}

	if result.ErrorCode != commonDef.Success {
		err = fmt.Errorf("update value failed, reason:%s", result.Reason)
		return
	}

	return
}

func (s *impl) QueryEntity(entityPtr interface{}) (err error) {
	entityVal, entityErr := remote.GetObjectValue(entityPtr)
	if entityErr != nil {
		err = entityErr
		return
	}

	result := &common.ObjectValueResult{}
	url := strings.Join([]string{s.serverURL, net.JoinPrefix(common.QueryValue, common.ApiVersion)}, "")
	if s.serviceName != "" {
		url = fmt.Sprintf("%s?source=%s", url, s.serviceName)
	}
	_, err = net.HTTPPost(s.httpClient, url, entityVal, result)
	if err != nil {
		return
	}

	if result.ErrorCode != commonDef.Success {
		err = fmt.Errorf("query value failed, reason:%s", result.Reason)
		return
	}

	val, valErr := remote.ConvertObjectValue(result.Value)
	if valErr != nil {
		err = valErr
		return
	}

	err = provider.UpdateEntity(val, entityPtr)
	return
}

func (s *impl) BatchQueryEntity(entitySlicePtr interface{}, filter *common.QueryFilter) (ret int64, err error) {
	entityVal, entityErr := remote.GetSliceObjectValue(entitySlicePtr)
	if entityErr != nil {
		err = entityErr
		return
	}

	result := &common.SliceObjectValueResult{}
	url := strings.Join([]string{s.serverURL, net.JoinPrefix(common.QueryValues, common.ApiVersion)}, "")
	if s.serviceName != "" {
		url = fmt.Sprintf("%s?source=%s", url, s.serviceName)
	}

	param := common.ObjectValueFilter{TypeName: entityVal.GetName(), PkgPath: entityVal.GetPkgPath(), ValueFilter: filter}
	_, err = net.HTTPPost(s.httpClient, url, param, result)
	if err != nil {
		return
	}

	if result.ErrorCode != commonDef.Success {
		err = fmt.Errorf("batch query value failed, reason:%s", result.Reason)
		return
	}

	sliceVal, sliceValErr := remote.ConvertSliceObjectValue(result.Value)
	if sliceValErr != nil {
		err = sliceValErr
		return
	}

	err = provider.UpdateSliceEntity(sliceVal, entitySlicePtr)
	ret = result.Total

	return
}

func (s *impl) CountEntity(entityPtr interface{}, filter *common.QueryFilter) (ret int64, err error) {
	entityVal, entityErr := remote.GetObjectValue(entityPtr)
	if entityErr != nil {
		err = entityErr
		return
	}

	result := &common.ObjectValueCountResult{}
	url := strings.Join([]string{s.serverURL, net.JoinPrefix(common.QueryCount, common.ApiVersion)}, "")
	if s.serviceName != "" {
		url = fmt.Sprintf("%s?source=%s", url, s.serviceName)
	}

	param := common.ObjectValueFilter{TypeName: entityVal.GetName(), PkgPath: entityVal.GetPkgPath(), ValueFilter: filter}
	_, err = net.HTTPPost(s.httpClient, url, param, result)
	if err != nil {
		return
	}

	if result.ErrorCode != commonDef.Success {
		err = fmt.Errorf("count value failed, reason:%s", result.Reason)
		return
	}

	ret = result.Total

	return
}

func (s *impl) Release() {
	if s.httpClient != nil {
		s.httpClient.CloseIdleConnections()
		s.httpClient = nil
	}
}

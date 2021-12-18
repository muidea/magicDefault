package service

import (
	"encoding/json"
	"net/http"

	bc "github.com/muidea/magicBatis/common"

	commonDef "github.com/muidea/magicCommon/def"
	fu "github.com/muidea/magicCommon/foundation/util"

	"github.com/muidea/magicDefault/common"
)

func (s *Setting) ViewSettingProfile(res http.ResponseWriter, req *http.Request) {
	pageFilter := fu.NewPageFilter()
	pageFilter.Decode(req)

	filter := bc.NewFilter()
	filter.Page(pageFilter)

	result := &common.ProfileResult{}
	for {
		curSession := s.sessionRegistry.GetSession(res, req)
		curNamespace := s.getCurrentNamespace(res, req)
		curEntity, curErr := s.getCurrentEntity(res, req)
		if curErr != nil {
			result.ErrorCode = commonDef.InvalidAuthority
			result.Reason = "invalid authority"
			break
		}

		curAccount, curErr := s.settingBiz.QuerySettingProfile(curSession.GetSessionInfo(), curEntity, curNamespace)
		if curErr != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = curErr.Error()
			break
		}

		filter.Equal("Creater", curEntity.ID)
		logList, logErr := s.settingBiz.QueryOperateLog(curSession.GetSessionInfo(), filter, curNamespace)
		if logErr != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = logErr.Error()
			break
		}
		for _, val := range logList {
			view := &common.LogView{}
			view.FromLog(val, curEntity)
			result.OperateLog = append(result.OperateLog, view)
		}

		result.Profile = curAccount
		result.ErrorCode = commonDef.Success
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)
}

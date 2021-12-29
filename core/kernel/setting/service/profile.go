package service

import (
	"context"
	"encoding/json"
	"net/http"

	bc "github.com/muidea/magicBatis/common"

	cd "github.com/muidea/magicCommon/def"
	fu "github.com/muidea/magicCommon/foundation/util"
	"github.com/muidea/magicCommon/session"

	"github.com/muidea/magicDefault/common"
)

func (s *Setting) ViewSettingProfile(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	curSession := ctx.Value(session.AuthSession).(session.Session)
	filter := fu.NewFilter()
	filter.Decode(req)

	queryFilter := bc.NewFilter()
	queryFilter.Page(filter.Pagination)

	result := &common.ProfileResult{}
	for {
		curNamespace := s.getCurrentNamespace(ctx, res, req)
		curEntity, curErr := s.getCurrentEntity(ctx, res, req)
		if curErr != nil {
			result.ErrorCode = cd.InvalidAuthority
			result.Reason = "invalid authority"
			break
		}

		curAccount, curErr := s.settingBiz.QuerySettingProfile(curSession.GetSessionInfo(), curEntity, curNamespace)
		if curErr != nil {
			result.ErrorCode = cd.Failed
			result.Reason = curErr.Error()
			break
		}

		queryFilter.Equal("Creater", curEntity.ID)
		logList, logErr := s.settingBiz.QueryOperateLog(curSession.GetSessionInfo(), queryFilter, curNamespace)
		if logErr != nil {
			result.ErrorCode = cd.Failed
			result.Reason = logErr.Error()
			break
		}
		for _, val := range logList {
			view := &common.LogView{}
			view.FromLog(val, curEntity)
			result.OperateLog = append(result.OperateLog, view)
		}

		result.Profile = curAccount
		result.ErrorCode = cd.Success
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)
}

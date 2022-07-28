package toolkit

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/muidea/magicCommon/def"
	engine "github.com/muidea/magicEngine"

	"github.com/muidea/magicCas/common"
)

// RoleVerifier role verifier
type RoleVerifier interface {
	CasVerifier
	VerifyRole(ctx context.Context, res http.ResponseWriter, req *http.Request) (*common.RoleView, error)
}

// RoleRegistry role route registry
type RoleRegistry interface {
	SetApiVersion(version string)

	AddHandler(pattern, method string, privilegeValue int, handler func(context.Context, http.ResponseWriter, *http.Request))

	AddRoute(route engine.Route, privilegeValue int, filters ...engine.MiddleWareHandler)

	GetAllPrivilege() []*common.Privilege
}

// NewRoleRegistry create routeRegistry
func NewRoleRegistry(verifier RoleVerifier, router engine.Router) (ret RoleRegistry) {
	ret = &roleRegistryImpl{roleVerifier: verifier, router: router, privilegeItemSlice: privilegeItemSlice{}}
	return
}

type privilegeItem struct {
	patternFilter  *engine.PatternFilter
	privilegeValue int
	patternPath    string
}

type privilegeItemSlice []*privilegeItem

// roleRegistryImpl cas route registry
type roleRegistryImpl struct {
	roleVerifier       RoleVerifier
	router             engine.Router
	privilegeItemSlice privilegeItemSlice
}

func (s *roleRegistryImpl) SetApiVersion(version string) {
	s.router.SetApiVersion(version)
}

// AddHandler add route handler
func (s *roleRegistryImpl) AddHandler(
	pattern, method string,
	privilegeValue int,
	handler func(context.Context, http.ResponseWriter, *http.Request)) {

	rtPattern := pattern
	apiVersion := s.router.GetApiVersion()
	if apiVersion != "" {
		rtPattern = fmt.Sprintf("%s%s", apiVersion, rtPattern)
	}

	privilegeItem := &privilegeItem{
		patternFilter:  engine.NewPatternFilter(rtPattern),
		privilegeValue: privilegeValue,
		patternPath:    rtPattern,
	}

	s.privilegeItemSlice = append(s.privilegeItemSlice, privilegeItem)

	s.router.AddRoute(engine.CreateRoute(pattern, method, handler), s)
}

func (s *roleRegistryImpl) AddRoute(route engine.Route, privilegeValue int, filters ...engine.MiddleWareHandler) {
	privilegeItem := &privilegeItem{
		patternFilter:  engine.NewPatternFilter(route.Pattern()),
		privilegeValue: privilegeValue,
		patternPath:    route.Pattern(),
	}

	s.privilegeItemSlice = append(s.privilegeItemSlice, privilegeItem)

	filters = append(filters, s)
	s.router.AddRoute(route, filters...)
}

func (s *roleRegistryImpl) GetAllPrivilege() (ret []*common.Privilege) {
	for _, val := range s.privilegeItemSlice {
		item := &common.Privilege{Path: val.patternPath, Value: common.GetPermission(val.privilegeValue)}

		ret = append(ret, item)
	}

	return
}

// Handle middleware handler
func (s *roleRegistryImpl) Handle(ctx engine.RequestContext, res http.ResponseWriter, req *http.Request) {
	result := &def.Result{ErrorCode: def.Success}
	for {
		// must verify cas
		_, casErr := s.roleVerifier.Verify(ctx.Context(), res, req)
		if casErr != nil {
			result.ErrorCode = def.InvalidAuthority
			result.Reason = casErr.Error()
			break
		}

		//casCtx := context.WithValue(ctx.Context(), session.AuthAccount, casEntity)
		casRole, casErr := s.roleVerifier.VerifyRole(ctx.Context(), res, req)
		if casErr != nil {
			result.ErrorCode = def.InvalidAuthority
			result.Reason = casErr.Error()
			break
		}

		patternPath := ""
		privilegeValue := 0
		for _, val := range s.privilegeItemSlice {
			if val.patternFilter.Match(req.URL.Path) {
				patternPath = val.patternPath
				privilegeValue = val.privilegeValue
				break
			}
		}

		err := s.verifyRole(casRole, patternPath, privilegeValue)
		if err != nil {
			result.ErrorCode = def.InvalidAuthority
			result.Reason = err.Error()
			break
		}

		//roleCtx := context.WithValue(casCtx, session.AuthRole, casRole)
		//ctx.Update(roleCtx)
		break
	}

	if result.Fail() {
		block, err := json.Marshal(result)
		if err == nil {
			res.Write(block)
			return
		}

		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx.Next()
}

func (s *roleRegistryImpl) verifyRole(rolePtr *common.RoleView, privilegePath string, privilegeValue int) (err error) {
	var privilegeLite *common.Privilege
	//for {
	// 如果是处于初始化状态的administrator账号，则认为有权限(特殊判断)
	//if accountInfoVal.Account == "administrator" && accountInfoVal.Status.IsInitStatus() && accountInfoVal.Role == nil {
	//	return nil
	//}

	privilegeLite = s.checkPrivilege(privilegePath, rolePtr)
	//	break
	//}
	if privilegeLite == nil {
		return fmt.Errorf("无效权限组")
	}

	if privilegeLite.Value.Value >= privilegeValue {
		return nil
	}

	return fmt.Errorf("当前账号无操作权限")
}

func (s *roleRegistryImpl) checkPrivilege(privilegePath string, rolePtr *common.RoleView) (ret *common.Privilege) {
	if rolePtr == nil {
		return
	}

	for ii := range rolePtr.Privilege {
		val := rolePtr.Privilege[ii]
		if val.Path == "*" {
			ret = val
			break
		}

		if val.Path == privilegePath {
			ret = val
			break
		}
	}

	return
}

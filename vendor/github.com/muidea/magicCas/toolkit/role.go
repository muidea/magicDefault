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

	AddHandler(pattern, method string, privateValue int, handler func(context.Context, http.ResponseWriter, *http.Request))

	AddRoute(route engine.Route, privateValue int, filters ...engine.MiddleWareHandler)

	GetAllPrivateItem() []*common.PrivateItem
}

// NewRoleRegistry create routeRegistry
func NewRoleRegistry(verifier RoleVerifier, router engine.Router) (ret RoleRegistry) {
	ret = &roleRegistryImpl{roleVerifier: verifier, router: router, privateItemSlice: privateItemSlice{}}
	return
}

type privateItem struct {
	patternFilter *engine.PatternFilter
	privateValue  int
	patternPath   string
}

type privateItemSlice []*privateItem

// roleRegistryImpl cas route registry
type roleRegistryImpl struct {
	roleVerifier     RoleVerifier
	router           engine.Router
	privateItemSlice privateItemSlice
}

func (s *roleRegistryImpl) SetApiVersion(version string) {
	s.router.SetApiVersion(version)
}

// AddHandler add route handler
func (s *roleRegistryImpl) AddHandler(
	pattern, method string,
	privateValue int,
	handler func(context.Context, http.ResponseWriter, *http.Request)) {

	rtPattern := pattern
	apiVersion := s.router.GetApiVersion()
	if apiVersion != "" {
		rtPattern = fmt.Sprintf("%s%s", apiVersion, rtPattern)
	}

	privateItem := &privateItem{
		patternFilter: engine.NewPatternFilter(rtPattern),
		privateValue:  privateValue,
		patternPath:   rtPattern,
	}

	s.privateItemSlice = append(s.privateItemSlice, privateItem)

	s.router.AddRoute(engine.CreateRoute(pattern, method, handler), s)
}

func (s *roleRegistryImpl) AddRoute(route engine.Route, privateValue int, filters ...engine.MiddleWareHandler) {
	privateItem := &privateItem{
		patternFilter: engine.NewPatternFilter(route.Pattern()),
		privateValue:  privateValue,
		patternPath:   route.Pattern(),
	}

	s.privateItemSlice = append(s.privateItemSlice, privateItem)

	filters = append(filters, s)
	s.router.AddRoute(route, filters...)
}

func (s *roleRegistryImpl) GetAllPrivateItem() (ret []*common.PrivateItem) {
	for _, val := range s.privateItemSlice {
		item := &common.PrivateItem{Path: val.patternPath, Value: common.GetPrivateInfo(val.privateValue)}

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

		privatePattern := ""
		privateValue := 0
		for _, val := range s.privateItemSlice {
			if val.patternFilter.Match(req.URL.Path) {
				privatePattern = val.patternPath
				privateValue = val.privateValue
				break
			}
		}

		err := s.verifyRole(casRole, privatePattern, privateValue)
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

func (s *roleRegistryImpl) verifyRole(privateRole *common.RoleView, privatePath string, privateValue int) (err error) {
	var privateLite *common.PrivateItem
	//for {
	// ?????????????????????????????????administrator???????????????????????????(????????????)
	//if accountInfoVal.Account == "administrator" && accountInfoVal.Status.IsInitStatus() && accountInfoVal.Role == nil {
	//	return nil
	//}

	privateLite = s.checkPrivate(privatePath, privateRole)
	//	break
	//}
	if privateLite == nil {
		return fmt.Errorf("???????????????")
	}

	if privateLite.Value.Value >= privateValue {
		return nil
	}

	return fmt.Errorf("???????????????????????????")
}

func (s *roleRegistryImpl) checkPrivate(privatePath string, privateRole *common.RoleView) (ret *common.PrivateItem) {
	if privateRole == nil {
		return
	}

	for ii := range privateRole.Private {
		val := privateRole.Private[ii]
		if val.Path == "*" {
			ret = val
			break
		}

		if val.Path == privatePath {
			ret = val
			break
		}
	}

	return
}

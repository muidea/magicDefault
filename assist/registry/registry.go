package registry

import (
	"net/http"
	"sync"

	"github.com/muidea/magicCas/toolkit"
	"github.com/muidea/magicCommon/session"
	engine "github.com/muidea/magicEngine"
)

var registryOnce sync.Once
var registryImpl *impl

func New(
	sessionCallBack session.CallBack,
	casVerifier toolkit.CasVerifier,
	roleVerifier toolkit.RoleVerifier,
	router engine.Router,
) {
	registryOnce.Do(func() {
		sessionRegistry := session.CreateRegistry(sessionCallBack)
		casRegistry := toolkit.NewCasRegistry(casVerifier, router)
		roleRegistry := toolkit.NewRoleRegistry(roleVerifier, router)

		registryImpl = &impl{
			sessionRegistry:   sessionRegistry,
			casRouteRegistry:  casRegistry,
			roleRouteRegistry: roleRegistry,
			router:            router,
		}
	})
}

func GetSession(res http.ResponseWriter, req *http.Request) session.Session {
	return registryImpl.sessionRegistry.GetSession(res, req)
}

func GetCasRegistry() toolkit.CasRegistry {
	return registryImpl.casRouteRegistry
}

func GetRoleRegistry() toolkit.RoleRegistry {
	return registryImpl.roleRouteRegistry
}

func GetRouter() engine.Router {
	return registryImpl.router
}

type impl struct {
	sessionRegistry   session.Registry
	casRouteRegistry  toolkit.CasRegistry
	roleRouteRegistry toolkit.RoleRegistry
	router            engine.Router
}

func (s *impl) SessionRegistry() session.Registry {
	return registryImpl.sessionRegistry
}

func (s *impl) CasRegistry() toolkit.CasRegistry {
	return registryImpl.casRouteRegistry
}

func (s *impl) RoleRegistry() toolkit.RoleRegistry {
	return registryImpl.roleRouteRegistry
}

func (s *impl) Router() engine.Router {
	return registryImpl.router
}

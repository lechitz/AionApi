package rotas

import (
	"github.com/lechitz/AionApi/src/controllers"
	"net/http"
)

var routeLogin = Route{
	URI:          "/login",
	Method:       http.MethodPost,
	Function:     controllers.Login,
	AuthRequired: false,
}

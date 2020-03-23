package app

import (
	"github.com/shohruk56/BookStorage/pkg/core/token"
	"github.com/shohruk56/BookStorage/pkg/mux/middleware/authenticated"
	"github.com/shohruk56/BookStorage/pkg/mux/middleware/authorized"
	"github.com/shohruk56/BookStorage/pkg/mux/middleware/jwt"
	"github.com/shohruk56/BookStorage/pkg/mux/middleware/logger"
	"reflect"
)

func (s Server) InitRoutes() {

	s.router.GET(
		"/",
		s.handleIndex(),
		logger.Logger("Index"),
	)

	s.router.GET(
		"/api/products",
		s.handleProductList(),
		authenticated.Authenticated(jwt.IsContextNonEmpty),
		jwt.JWT(reflect.TypeOf((*token.Payload)(nil)).Elem(), s.secret),
		logger.Logger("get list"),
	)

	s.router.GET(
		"/api/products/{id}",
		s.handleProductByID(),
		authenticated.Authenticated(jwt.IsContextNonEmpty),
		jwt.JWT(reflect.TypeOf((*token.Payload)(nil)).Elem(), s.secret),
		logger.Logger("get product by id"),
	)

	s.router.POST(
		"/api/products/{id}",
		s.handProduct(),
		authenticated.Authenticated(jwt.IsContextNonEmpty),
		authorized.Authorized([]string{"Admin"}, jwt.FromContext),
		jwt.JWT(reflect.TypeOf((*token.Payload)(nil)).Elem(), s.secret),
		logger.Logger("post product"),
	)

	s.router.DELETE(
		"/api/products/{id}",
		s.handleDeleteProduct(),
		authenticated.Authenticated(jwt.IsContextNonEmpty),
		authorized.Authorized([]string{"Admin"}, jwt.FromContext),
		jwt.JWT(reflect.TypeOf((*token.Payload)(nil)).Elem(), s.secret),
		logger.Logger("delete product"),
	)


}
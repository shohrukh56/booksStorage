package app

import (
	"github.com/shohruk56/BookStorage/pkg/core/product"
	"github.com/shohrukh56/jwt/pkg/jwt"
	"github.com/shohrukh56/mux/pkg/mux"
	"github.com/shohrukh56/rest/pkg/rest"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
)

type Server struct {
	router     *mux.ExactMux
	productSvc *product.Service
	secret     jwt.Secret
}

func NewServer(router *mux.ExactMux, productSvc *product.Service, secret jwt.Secret) *Server {
	return &Server{router: router, productSvc: productSvc, secret: secret}
}

func (s Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.router.ServeHTTP(writer, request)
}

func (s Server) Start() {
	s.InitRoutes()
}

func (s *Server) handleIndex() http.HandlerFunc {

	var (
		tpl *template.Template
		err error
	)
	tpl, err = template.ParseFiles(
		filepath.Join("web/templates", "index.gohtml"),
	)
	if err != nil {
		panic(err)
	}
	return func(writer http.ResponseWriter, request *http.Request) {
		// executes in many goroutines
		// TODO: fetch data from multiple upstream services
		err = tpl.Execute(writer, struct{ Title string }{Title: "auth",})
		if err != nil {
			log.Printf("error while executing template %s %v", tpl.Name(), err)
		}
	}

}
func (s Server) handleProductList() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		list, err := s.productSvc.ProductList(request.Context())
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Print(err)
			return
		}
		err = rest.WriteJSONBody(writer, &list)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Print(err)
		}
	}
}

func (s Server) handleProductByID() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		idFromCTX, ok := mux.FromContext(request.Context(), "id")
		if !ok {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(idFromCTX)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		prod, err := s.productSvc.ProductByID(request.Context(), int64(id))
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Print(err)
			return
		}
		err = rest.WriteJSONBody(writer, &prod)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Print(err)
		}
	}
}

//func (s Server) handleNewProduct() http.HandlerFunc {
//	return func(writer http.ResponseWriter, request *http.Request) {
//		get := request.Header.Get("Content-Type")
//		if get != "application/json" {
//			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
//			return
//		}
//		prod := product.Product{}
//		err := rest.ReadJSONBody(request, &prod)
//		if err != nil {
//			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
//			return
//		}
//		err = s.productSvc.AddNewProduct(request.Context(), prod)
//		if err != nil {
//			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
//			log.Print(err)
//			return
//		}
//		_, err = writer.Write([]byte("New Product Added!"))
//		if err != nil {
//			log.Print(err)
//		}
//	}
//}

func (s Server) handleDeleteProduct() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		idFromCTX, ok := mux.FromContext(request.Context(), "id")
		if !ok {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(idFromCTX)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		err = s.productSvc.RemoveByID(request.Context(), int64(id))
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Print(err)
			return
		}
		writer.WriteHeader(http.StatusNoContent)
	}
}

func (s Server) handProduct() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		context, ok := mux.FromContext(request.Context(), "id")
		if !ok {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(context)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		get := request.Header.Get("Content-Type")
		if get != "application/json" {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		prod := product.Product{}
		err = rest.ReadJSONBody(request, &prod)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		if id == 0 {
			err = s.productSvc.AddNewProduct(request.Context(), prod)
			if err != nil {
				http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				log.Print(err)
				return
			}
			writer.WriteHeader(http.StatusNoContent)
			return
		}
		if id > 0 {
			err = s.productSvc.UpdateProduct(request.Context(), int64(id), prod)
			if err != nil {
				http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				log.Print(err)
				return
			}
			writer.WriteHeader(http.StatusNoContent)
			return
		}
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

	}
}

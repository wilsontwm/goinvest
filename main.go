package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"goinvest/controllers"
	"goinvest/middleware"
	"log"
	"net/http"
	"os"
)

type route struct {
	Router *mux.Router
	Path   string
	Func   func(http.ResponseWriter, *http.Request)
	Method string
}

var routes []route

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	router := mux.NewRouter()
	setupRoutes(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server started and running at port", port)

	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{"*"})
	log.Fatal(http.ListenAndServe(":"+port, handlers.CORS(headers, methods, origins)(router)))

}

// Setup REST routes
func setupRoutes(router *mux.Router) {
	apiRoutes := router.PathPrefix("/api").Subrouter()
	routes = append(routes, route{Router: apiRoutes, Path: "/articles", Func: controllers.ArticleList, Method: "GET"})
	routes = append(routes, route{Router: apiRoutes, Path: "/articles/crawl", Func: controllers.ArticleCrawl, Method: "POST"})
	routes = append(routes, route{Router: apiRoutes, Path: "/login", Func: controllers.UserLogin, Method: "POST"})

	apiAuthenticatedRoutes := apiRoutes.PathPrefix("/auth").Subrouter()
	apiAuthenticatedRoutes.Use(middleware.JwtAuthentication())

	// Account routes
	routes = append(routes, route{Router: apiAuthenticatedRoutes, Path: "/account/list", Func: controllers.AccountList, Method: "GET"})
	routes = append(routes, route{Router: apiAuthenticatedRoutes, Path: "/account/create", Func: controllers.AccountCreate, Method: "POST"})
	routes = append(routes, route{Router: apiAuthenticatedRoutes, Path: "/account/update", Func: controllers.AccountUpdate, Method: "POST"})
	routes = append(routes, route{Router: apiAuthenticatedRoutes, Path: "/account/delete", Func: controllers.AccountDelete, Method: "POST"})

	// Fund flow routes
	routes = append(routes, route{Router: apiAuthenticatedRoutes, Path: "/fundflow/list", Func: controllers.FundFlowList, Method: "GET"})
	// routes = append(routes, route{Router: apiAuthenticatedRoutes, Path: "/fundflow/create", Func: controllers.AccountCreate, Method: "POST"})
	// routes = append(routes, route{Router: apiAuthenticatedRoutes, Path: "/fundflow/update", Func: controllers.AccountUpdate, Method: "POST"})
	// routes = append(routes, route{Router: apiAuthenticatedRoutes, Path: "/fundflow/delete", Func: controllers.AccountDelete, Method: "POST"})

	for _, r := range routes {
		r.Router.HandleFunc(r.Path, r.Func).Methods(r.Method)
	}
}

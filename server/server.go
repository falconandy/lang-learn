package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

const (
	defaultHTTPPort = 8080
)

type Server struct {
	port        int
	frontendDir string
}

func NewServer() *Server {
	return &Server{
		frontendDir: filepath.Join("front", "build"),
		port:        defaultHTTPPort,
	}
}

func (s *Server) Start() {
	r := chi.NewRouter()

	s.setupMiddlewares(r)

	r.Route("/api", func(r chi.Router) {
		// TODO
	})

	s.serveFrontendFiles(r, "/static/", filepath.Join(s.frontendDir, "static"))
	r.Get("/*", s.serveIndex)

	log.Printf("started on :%d\n", s.port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", s.port), r))
}

func (s *Server) setupMiddlewares(r *chi.Mux) {
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
}

func (s *Server) serveFrontendFiles(r *chi.Mux, path string, rootDir string) {
	root := http.Dir(rootDir)
	fs := http.StripPrefix(path, http.FileServer(root))

	r.Get(path+"*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}

func (s *Server) serveIndex(w http.ResponseWriter, _ *http.Request) {
	data, err := ioutil.ReadFile(filepath.Join(s.frontendDir, "index.html"))
	if err != nil {
		panic(err)
	}
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	_, _ = w.Write(data)
}

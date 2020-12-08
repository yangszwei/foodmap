package delivery

import (
	"github.com/gin-gonic/gin"
)

// NewServer create an instance of Server
func NewServer(data string) *Server {
	s := new(Server)
	s.Engine = gin.Default()
	s.root = s.Engine.Group("/")
	return s
}

// Server wrapper of gin.Engine
type Server struct {
	Engine *gin.Engine
	root   *gin.RouterGroup // for middlewares
	Web    *gin.RouterGroup // front-end web
	API    *gin.RouterGroup // api
	data   string           // folder path for storing data files
}

// SetupRouterGroups setup router groups, should be called AFTER middlewares
// are set on server.root so that the middlewares are applied to all route
// groups
func (s *Server) SetupRouterGroups() {
	s.Web = s.root.Group("/")
	s.API = s.root.Group("/api")
}

// Run start server
func (s *Server) Run(addr string) error {
	return s.Engine.Run(addr)
}

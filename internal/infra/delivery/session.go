package delivery

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
)

// Setup session setup session middleware
func SetupSession(s *Server, secret []byte) {
	store := cookie.NewStore(secret)
	s.root.Use(sessions.Sessions("session", store))
}

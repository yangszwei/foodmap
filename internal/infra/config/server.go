package config

// ServerConfig server config ()
type ServerConfig struct {
	Addr          string `json:"ADDR" validate:"required"`           // passed to gin.Server.Run
	SessionSecret string `json:"SESSION_SECRET" validate:"required"` // passed to gin-contrib/sessions
	DataPath      string `json:"DATA_PATH" validate:"required"`      // folder path for storing files
}

package config

// Root is root config
type Root struct {
	Device *Device `toml:"device" validate:"required"`
	Logger *Logger `toml:"logger" validate:"required"`
	Server *Server `toml:"server"`
}

// Device is device property
type Device struct {
	Address string `toml:"address"`
	Lang    string `toml:"lang"`
	Timeout string `toml:"timeout" validate:"required"`
}

// Logger is zap logger property
type Logger struct {
	Service      string `toml:"service" validate:"required"`
	Env          string `toml:"env" validate:"oneof=dev prod custom"`
	Level        string `toml:"level" validate:"required"`
	IsStackTrace bool   `toml:"is_stacktrace"`
}

// Server is server property
type Server struct {
	Port int `toml:"port"`
}

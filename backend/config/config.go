package config

import "alexlupatsiy.com/personal-website/backend/db"

// Contains all configurations of our backend by combining smaller configurations of other subsystems.
// Like DbConfiguration, later InfluxConfiguration etc.
type Config struct {
	DbConfig db.Config
	DevMode  bool `env:"DEV_MODE, default=false"`
}

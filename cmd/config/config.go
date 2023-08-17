package config

import (
	"log"
)

// AppConfig is a instance of Config that is accessible in entire application
var AppConfig Config

// Config is our wide configuration for the application
type Config struct {
	HTTPPort   int
	Domain     string
	DSN        string
	ErrorLog   *log.Logger
	InfoLog    *log.Logger
	WarningLog *log.Logger
}

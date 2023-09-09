package env

import (
	"os"
	"strconv"
)

// GetEnvOrDefaultString read string data from env file
func GetEnvOrDefaultString(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetEnvOrDefaultBool read bool data from env file
func GetEnvOrDefaultBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	boolVal, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return boolVal
}

// GetEnvOrDefaultInt read int data from env file
func GetEnvOrDefaultInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	intVal, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intVal
}

// GetEnvOrDefaultFloat32 read float32 type from env file
func GetEnvOrDefaultFloat32(key string, defaultValue float32) float32 {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	float32Val, err := strconv.ParseFloat(value, 32)
	if err != nil {
		return defaultValue
	}
	return float32(float32Val)
}

// GetEnvOrDefaultFloat64 read float64 file from env file
func GetEnvOrDefaultFloat64(key string, defaultValue float64) float64 {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	float64Val, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return defaultValue
	}
	return float64Val
}

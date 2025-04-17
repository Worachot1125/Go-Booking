package config

import (
	"log"
	"os"
	"strconv"

	"github.com/spf13/viper"
)

// conf is a utility function to get the value from environment variables or use the default value
// confString is a utility function to get the value from environment variables or use the default string value
func confString(key string, defaultValue string) string {
	viper.AutomaticEnv()

	if viper.IsSet(key) {
		return viper.GetString(key)
	}

	return defaultValue
}

// confInt64 is a utility function to get the value from environment variables or use the default int64 value
func confInt64(key string, defaultValue int64) int64 {
	viper.AutomaticEnv()

	if viper.IsSet(key) {
		valueStr := viper.GetString(key)
		valueInt, err := strconv.ParseInt(valueStr, 10, 64)
		if err != nil {
			log.Printf("Error converting %s to int64: %v", key, err)
			return defaultValue
		}
		return valueInt
	}

	return defaultValue
}

func conf[T string | ~int | ~int32 | ~int64 | bool](key string, fallback T) T {
	viper.Set(key, fallback)
	if value, ok := os.LookupEnv(key); ok {
		viper.Set(key, value)
	}
	err := viper.UnmarshalKey(key, &fallback)
	if err != nil {
		panic(err)
	}
	return fallback
}

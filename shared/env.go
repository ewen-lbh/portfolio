package shared

import "os"

// IsDev returns true if the environment variable "ENV" is set to "development"
func IsDev() bool {
	return os.Getenv("ENV") == "development"
}

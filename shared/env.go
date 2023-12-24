package shared

import "os"

// IsDev returns true if the environment variable "ENV" is set to "development"
func IsDev() bool {
	return os.Getenv("ENV") == "development"
}

// WantToRemoveUnusedMessages returns true if the user wants to remove unused translation messages from the .po file (even when an associated non-empty msgstr exists)
func WantToRemoveUnusedMessages() bool {
	return os.Getenv("REMOVE_UNUSED_MESSAGES") == "1"
}

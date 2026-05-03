package lib

import (
	"os"
)

// Appending :tags creates an alternate data stream in NTFS
const streamSuffix string = ":tags"

// GetTags reads tags from the file's alternate data stream
func GetTags(path string) (string, error) {
	data, err := os.ReadFile(path + streamSuffix)
	if err != nil {
		// If the stream doesn't exist, return empty (Windows throws PathError)
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	return string(data), nil
}

// SetTag writes tags to the file's alternate data stream
func SetTag(path, tag string) error {
	// Overwrites or creates the stream with the tag string
	err := os.WriteFile(path+streamSuffix, []byte(tag), 0644)
	return err
}

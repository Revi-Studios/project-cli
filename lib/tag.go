package lib

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/pkg/xattr"
	"howett.net/plist"
)

const attrname string = "com.apple.metadata:_kMDItemUserTags"

func GetTags(path string) ([]string, error) {

	data, err := xattr.Get(path, attrname)
	if err != nil {
		return nil, err
	}
	var tags []string
	decoder := plist.NewDecoder(bytes.NewReader(data))
	decoder.Decode(&tags)

	for i, t := range tags {
		tag := strings.Split(t, "\n")
		tags[i] = tag[0]
	}

	return tags, nil
}

func SetTags(path string, tags ...string) error {
	if len(tags) == 0 {
		return nil
	}
	var buf bytes.Buffer

	if err := plist.NewEncoder(&buf).Encode(tags); err != nil {
		return fmt.Errorf("failed to encode tags: %w", err)
	}
	if err := xattr.Set(path, attrname, buf.Bytes()); err != nil {
		return fmt.Errorf("failed to set tags: %w", err)
	}
	return nil
}

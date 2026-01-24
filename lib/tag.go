package lib

import (
	"bytes"
	"strings"

	"github.com/pkg/xattr"
	"howett.net/plist"
)

const attrname string = "com.apple.metadata:_kMDItemUserTags"

func GetTags(path string) (string, error) {

	data, err := xattr.Get(path, attrname)
	if err != nil {
		return "", err
	}
	var tags []string
	decoder := plist.NewDecoder(bytes.NewReader(data))
	decoder.Decode(&tags)

	for i, t := range tags {
		tag := strings.Split(t, "\n")
		tags[i] = tag[0]
	}

	return strings.Join(tags, ", "), nil
}

func SetTag(path, tag string) error {
	var buf bytes.Buffer

	plist.NewEncoder(&buf).Encode([]string{tag})
	if err := xattr.Set(path, attrname, buf.Bytes()); err != nil {
		return err
	}
	return nil
}

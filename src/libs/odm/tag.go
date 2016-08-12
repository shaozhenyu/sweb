package odm

import (
	"encoding/json"
	"strings"

	"gopkg.in/mgo.v2"
)

const (
	mgoIndexTagName = "mgoIndex"
)

type MgoIndexTag struct {
	Index *mgo.Index
}

func stringsContains(arr []string, s string) bool {
	return stringsIndexOf(arr, s) != -1
}

func stringsIndexOf(arr []string, s string) int {
	for i, str := range arr {
		if strings.Trim(str, " ") == s {
			return i
		}
	}
	return -1
}

func newMgoIndexTag(index string) (*MgoIndexTag, error) {

	v := mgo.Index{}
	if err := json.Unmarshal([]byte(index), &v); err != nil {
		indexStr := strings.Split(index, ",")
		fields := strings.Split(indexStr[0], "+")
		v.Key = fields
		v.Unique = stringsContains(indexStr[1:], "unique")
		v.DropDups = stringsContains(indexStr[1:], "drop_dups")
		v.Background = stringsContains(indexStr[1:], "background")
		v.Sparse = stringsContains(indexStr[1:], "sparse")
	}

	return &MgoIndexTag{&v}, nil
}

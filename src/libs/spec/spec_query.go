package spec

import (
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"gopkg.in/mgo.v2/bson"
)

var (
	MaxListSize = 100
	defaultSize = 20
)

type Paging struct {
	From  int `json:"from"`
	Size  int `json:"size"`
	Total int `json:"total"`
}

type ListQuery struct {
	Cond   map[string]interface{}
	Paging Paging
	Sorts  []string
}

type ListResult struct {
	Result interface{} `json:"result"`
	Paging Paging      `json:"paging"`
}

func NewListQueryWithReq(typ_ reflect.Type, req *http.Request) *ListQuery {
	q := NewListQuery(typ_, req.URL.Query())
	return q
}

func NewListQuery(typ_ reflect.Type, vals url.Values) *ListQuery {

	getInt := func(key string) int {
		if str, ok := vals[key]; ok && len(str) > 0 {
			if val, err := strconv.ParseInt(str[0], 10, 64); err == nil {
				return int(val)
			}
		}
		return 0
	}

	getString := func(key string) []string {
		if str, ok := vals[key]; ok {
			return strings.Split(str[0], ",")
		}
		return []string{}
	}

	cond := map[string]interface{}{}

	for i := 0; i < typ_.NumField(); i++ {
		f := typ_.Field(i)
		bsonTag := f.Tag.Get("bson")
		if bsonTag == "" {
			continue
		}
		if bsonTag = strings.Split(bsonTag, ",")[0]; bsonTag == "-" {
			continue
		}

		jsonTag := f.Tag.Get("json")
		if jsonTag == "" {
			continue
		}
		if jsonTag = strings.Split(jsonTag, ",")[0]; jsonTag == "-" {
			continue
		}

		strs, ok := vals[jsonTag]
		if !ok {
			continue
		}
		if len(strs) == 0 {
			continue
		}

		kind := f.Type.Kind()
		if kind == reflect.Slice {
			kind = f.Type.Elem().Kind()
		}

		switch kind {
		case reflect.String:
			{
				cond[bsonTag] = map[string]interface{}{"$regex": bson.RegEx{strs[0], "i"}}
			}
		case reflect.Int, reflect.Int64:
			{
				if len(strs) > 1 {
					iarr := []int64{}
					for _, s := range strs {
						if i_, err := strconv.ParseInt(s, 10, 64); err == nil {
							iarr = append(iarr, i_)
						}
					}
					cond[bsonTag] = map[string]interface{}{"$in": iarr}
				} else if len(strs) > 0 && len(strs[0]) > 0 {
					if strs[0][0] == '!' {
						if i_, err := strconv.ParseInt(strs[0][1:], 10, 64); err == nil {
							cond[bsonTag] = map[string]interface{}{"$ne": i_}
						}
					} else {
						cond[bsonTag] = parseInt(strs)
					}
				}
			}
		case reflect.Bool:
			{
				if strs[0] == "true" {
					cond[bsonTag] = true
				}

			}
		default:
			{
				if len(strs) > 1 {
					cond[bsonTag] = map[string]interface{}{"$in": strs}
				} else if len(strs) > 0 {
					cond[bsonTag] = strs[0]
				}
			}
		}
	}

	size := getInt("size")
	if size < 1 {
		size = defaultSize
	} else if size > MaxListSize {
		size = MaxListSize
	}

	from := getInt("from")
	if from < 0 {
		from = 0
	}

	page := Paging{
		From:  from,
		Size:  size,
		Total: 0,
	}

	return &ListQuery{
		Cond:   cond,
		Paging: page,
		Sorts:  getString("sort"),
	}
}

func parseInt(vals []string) map[string]interface{} {

	if len(vals) == 0 {
		return nil
	}

	prefixs := []string{"$gte", "$lte", "$gt", "$lt"}

	m := map[string]interface{}{}
	in := []int64{}

	conds := strings.Split(vals[0], ",")
	for _, v := range conds {
		parsed := false
		for _, pre := range prefixs {
			if strings.HasPrefix(v, pre) {
				parsed = true
				if i, err := strconv.ParseInt(v[len(pre):], 10, 64); err == nil {
					m[pre] = i
				}
				break
			}
		}

		if parsed == false {
			if i, err := strconv.ParseInt(v, 10, 64); err == nil {
				in = append(in, i)
			}
		}
	}

	if len(in) == 0 {
		return m
	}

	r := map[string]interface{}{}

	for k, v := range m {
		r[k] = v
	}

	r["$in"] = in

	return r
}

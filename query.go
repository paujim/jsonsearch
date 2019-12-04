package jsonsearch

import (
	"fmt"
	"strings"
)

type JsonObjectQuery struct {
	fieldName string
	jsonObj   *JsonObject
	with      map[string]*JsonObject
}

func (obj *JsonObject) With(included map[string]*JsonObject) *JsonObjectQuery {
	return &JsonObjectQuery{with: included, jsonObj: obj}
}

func (q *JsonObjectQuery) Where(key string) *JsonObjectQuery {
	return &JsonObjectQuery{with: q.with, jsonObj: q.jsonObj, fieldName: key}
}

func (obj *JsonObject) Where(fieldName string) *JsonObjectQuery {
	return &JsonObjectQuery{fieldName: fieldName, jsonObj: obj}
}

func (q *JsonObjectQuery) IsEqualTo(value string) []RawJson {
	var results []RawJson
	if q.jsonObj.IsIdField(q.fieldName) {
		val, ok := q.jsonObj.GetById(value)
		if !ok {
			return []RawJson{}
		}
		results = []RawJson{val}
	}
	results = q.jsonObj.SearchAll(q.fieldName, value, func(v1 string, v2 string) bool {
		return v1 == v2
	})
	return q.combine(results)
}

func (q *JsonObjectQuery) Contains(value string) []RawJson {
	results := q.jsonObj.SearchAll(q.fieldName, value, func(v1 string, v2 string) bool {
		return strings.Contains(v1, v2)
	})
	return q.combine(results)
}

func (q *JsonObjectQuery) combine(results []RawJson) []RawJson {
	if q.with == nil || len(q.with) == 0 {
		return results
	}

	for _, raw := range results {
		for fieldName, fieldValue := range raw {
			fieldVal := fmt.Sprintf("%v", fieldValue)
			obj, ok := q.with[fieldName]
			if ok {
				idfound, found := obj.GetById(fieldVal)
				if found {
					raw[fieldName] = idfound
				}
			}
		}
	}
	return results
}

package jsonsearch

import (
	"encoding/json"
	"fmt"
	"log"
)

type SetItem struct{}
type Set map[int]SetItem

type RawJson map[string]interface{}

type Comparable func(string, string) bool

func (obj *RawJson) String() string {
	json, err := json.MarshalIndent(obj, "", "\t")
	if err != nil {
		return err.Error()
	}
	return string(json)
}

type JsonObject struct {
	logger      *log.Logger
	idFieldName string
	name        string
	indices     map[string]int
	keys        map[string]Set
	jsonData    []RawJson
}

func logMessage(logger *log.Logger, message string) {
	if logger != nil {
		if logger != nil {
			logger.Println(message)
		}
	}
}

func buildKeysAndIndices(logger *log.Logger, data []RawJson, idFieldName string) (map[string]Set, map[string]int) {
	keys := make(map[string]Set)
	indices := make(map[string]int)
	for index, raw := range data {
		val, ok := raw[idFieldName]
		if !ok {
			logMessage(logger, "The idkey is not present on the object")
			return nil, nil
		}
		keyValue := fmt.Sprintf("%v", val)
		indices[keyValue] = index
		for k, _ := range raw {
			keys[k] = Set{}
		}
	}
	return keys, indices
}

func NewJsonObject(logger *log.Logger, data []byte, name, idFieldName string) *JsonObject {
	logMessage(logger, "Creating "+name)

	var raw []RawJson
	err := json.Unmarshal(data, &raw)
	if err != nil {
		logMessage(logger, err.Error())
		return nil
	}
	keys, indices := buildKeysAndIndices(logger, raw, idFieldName)
	if keys == nil && indices == nil {
		return nil
	}
	return &JsonObject{logger: logger, jsonData: raw, keys: keys, idFieldName: idFieldName, indices: indices, name: name}
}

func (obj *JsonObject) IsIdField(fieldName string) bool {
	return fieldName == obj.idFieldName
}

func (obj *JsonObject) GetById(id string) (RawJson, bool) {
	logMessage(obj.logger, "Get by "+id)
	index, ok := obj.indices[id]
	if !ok {
		logMessage(obj.logger, "id Not found")
		return RawJson{}, false
	}
	return obj.jsonData[index], true
}

func (obj *JsonObject) Keys() []string {
	keys := []string{}
	for key, _ := range obj.keys {
		keys = append(keys, key)
	}
	return keys
}

func (obj *JsonObject) SearchAll(fieldName, fieldValue string, compareAction Comparable) []RawJson {
	response := []RawJson{}
	for _, v := range obj.jsonData {
		addToResponse := false
		field := v[fieldName]
		switch x := field.(type) {
		case []interface{}:
			for _, i := range x {
				if hasFieldValue(i, fieldValue, compareAction) {
					addToResponse = true
					break
				}
			}
		case interface{}:
			addToResponse = hasFieldValue(x, fieldValue, compareAction)
		}

		if addToResponse {
			response = append(response, v)
		}
	}
	return response
}

func hasFieldValue(field interface{}, fieldValue string, compareAction Comparable) bool {
	value := fmt.Sprintf("%v", field)
	return compareAction(value, fieldValue)
}

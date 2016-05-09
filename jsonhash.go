package jsonhash

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"hash"
	"sort"
	"strings"
)

const toescape = `x`
const escape = `_x`

// Hash a valid json map.
func Hash(json map[string]interface{}) []byte {
	hash := md5.New()
	hashObject(json, hash)

	return hash.Sum(nil)
}

// HashS generates a base64 encoded json hash.
// Can be used in filenames and urls.
func HashS(json map[string]interface{}) string {

	h := Hash(json)
	return base64.URLEncoding.EncodeToString(h)
}

func hashObject(json map[string]interface{}, hash hash.Hash) {

	var keys = sortedKeys(json)
	for _, key := range keys {
		v := json[key]

		hash.Write([]byte(escapeKey(key)))

		switch o := v.(type) {
		case map[string]interface{}:
			hashObject(o, hash)
		case []interface{}:
			hashArray(o, hash)
		default:
			hash.Write([]byte(escapeValue(fmt.Sprintf("%T", o), fmt.Sprintf("%v", o))))
		}
	}
}

func hashArray(jsonArray []interface{}, hash hash.Hash) {

	for _, v := range jsonArray {

		switch o := v.(type) {
		case map[string]interface{}:
			hashObject(o, hash)
		case []interface{}:
			hashArray(o, hash)
		default:
			hash.Write([]byte(escapeValue(fmt.Sprintf("%T", o), fmt.Sprintf("%v", o))))
		}
	}
}

func escapeKey(key string) string {

	return toescape + escapeString(key)
}

func escapeValue(vtype string, value string) string {
	return toescape + escapeString(vtype) + toescape + toescape + escapeString(value)
}

func escapeString(value string) string {

	return strings.Replace(value, toescape, escape, -1)
}

func sortedKeys(json map[string]interface{}) []string {
	var keys []string
	for k := range json {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	return keys
}

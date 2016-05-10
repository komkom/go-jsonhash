package jsonhash

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"hash"
	"sort"
)

var escsq = [...]rune{'[', ']', '{', '}', '#', '/'}

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

	hash.Write([]byte(`{`))
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
			hash.Write([]byte(escapeValue(o)))
		}
	}

	hash.Write([]byte(`}`))
}

func hashArray(jsonArray []interface{}, hash hash.Hash) {

	hash.Write([]byte(`[`))
	for _, v := range jsonArray {

		switch o := v.(type) {
		case map[string]interface{}:
			hashObject(o, hash)
		case []interface{}:
			hashArray(o, hash)
		default:
			hash.Write([]byte(escapeValue(o)))
		}
	}

	hash.Write([]byte(`]`))
}

func escapeKey(key string) string {
	return escapeString(key)
}

func escapeValue(value interface{}) string {
	return `/` + escapeString(fmt.Sprintf(`%T`, value)) + `#` + escapeString(fmt.Sprintf("%v", value)) + `/`
}

func escapeString(value string) string {

	b := bytes.Buffer{}

	for _, r := range value {
		check := false
		for _, er := range escsq {
			if er == r {
				check = true
				b.WriteRune('_')
				b.WriteRune(r)
			}
		}

		if !check {
			b.WriteRune(r)
		}
	}

	//fmt.Println(b.String())

	return b.String()
}

func sortedKeys(json map[string]interface{}) []string {
	var keys []string
	for k := range json {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	return keys
}

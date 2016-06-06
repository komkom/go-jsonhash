package jsonhash

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"hash"
	"sort"
)

// Hash a valid json map.
func Hash(json map[string]interface{}) []byte {
	hash := md5.New()
	hashObject(json, &hash)

	return hash.Sum(nil)
}

// HashS generates a base64 encoded json hash.
// Can be used in filenames and urls.
func HashS(json map[string]interface{}) string {

	h := Hash(json)
	return base64.URLEncoding.EncodeToString(h)
}

func hashObject(json map[string]interface{}, hash *hash.Hash) {

	barrier(hash)
	var keys = sortedKeys(json)
	for _, key := range keys {
		v := json[key]

		escapeKey(key, hash)
		switch o := v.(type) {
		case map[string]interface{}:
			hashObject(o, hash)
		case []interface{}:
			hashArray(o, hash)
		default:
			escapeValue(o, hash)
		}
	}
	barrier(hash)
}

func hashArray(jsonArray []interface{}, hash *hash.Hash) {

	barrier(hash)
	for _, v := range jsonArray {

		switch o := v.(type) {
		case map[string]interface{}:
			hashObject(o, hash)
		case []interface{}:
			hashArray(o, hash)
		default:
			escapeValue(o, hash)
		}
	}
	barrier(hash)
}

func escapeKey(key string, writer *hash.Hash) {
	(*writer).Write([]byte(key))
	barrier(writer)
}

func escapeValue(value interface{}, writer *hash.Hash) {
	(*writer).Write([]byte(fmt.Sprintf(`%T`, value)))
	barrier(writer)
	(*writer).Write([]byte(fmt.Sprintf("%v", value)))
	barrier(writer)
}

// an invalide utf8 byte is used as a barrier
var barrierValue = []byte{byte('\255')}

func barrier(writer *hash.Hash) {
	(*writer).Write(barrierValue)
}

func sortedKeys(json map[string]interface{}) []string {
	var keys []string
	for k := range json {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	return keys
}

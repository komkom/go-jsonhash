package jsonhash

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"hash"
	"io"
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

	addCtrlSeq(`{`, hash)
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

	addCtrlSeq(`}`, hash)
}

func hashArray(jsonArray []interface{}, hash hash.Hash) {

	addCtrlSeq(`[`, hash)
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
	addCtrlSeq(`]`, hash)
}

func escapeKey(key string, writer io.Writer) {
	escapeString(key, writer)
}

func escapeValue(value interface{}, writer io.Writer) {

	addCtrlSeq(`/`, writer)
	escapeString(fmt.Sprintf(`%T`, value), writer)
	addCtrlSeq(`#`, writer)
	escapeString(fmt.Sprintf("%v", value), writer)
	addCtrlSeq(`/`, writer)
}

func addCtrlSeq(c string, writer io.Writer) {
	writer.Write([]byte(`*`))
	writer.Write([]byte(c))
}

func escapeString(value string, writer io.Writer) {

	for _, r := range value {
		for _, er := range escsq {
			if er == r {
				writer.Write([]byte(`_`))
				break
			}
		}

		writer.Write([]byte(string(r)))
	}
}

func sortedKeys(json map[string]interface{}) []string {
	var keys []string
	for k := range json {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	return keys
}

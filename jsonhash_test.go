package jsonhash

import (
	"strings"
	"testing"
	"tutti-golang/utils"
)

func TestJson(t *testing.T) {

	resources := []string{
		"resources/test.json",
		"resources/test2.json",
		"resources/test4.json",
		"resources/test5.json",
		"resources/test7.json",
		"resources/typetest.json",
		"resources/typetest2.json",
		"resources/typetest3.json",
		"resources/typetest4.json",
	}

	hmap := map[string]string{}

	// hash the test jsons.
	for _, r := range resources {
		split := strings.Split(r, "/")
		if len(split) != 2 {
			t.Fatal("illegal resource: " + r)
		}

		fn := split[1]

		j, err := utils.ParseJsonPath(r)
		if err != nil {
			t.Fatal(err.Error())
		}

		h := HashS(j)
		hmap[fn] = h
	}

	// assert hash(test.json) != hash(test2.json)
	Eval("test.json", "test2.json", hmap, false, t)

	// assert hash(test4.json) != hash(test5.json)
	Eval("test4.json", "test5.json", hmap, false, t)

	//Eval("test4.json", "test6.json", hmap, true, t)
	Eval("test6.json", "test7.json", hmap, false, t)

	Eval("typetest.json", "typetest2.json", hmap, false, t)

	Eval("typetest3.json", "typetest4.json", hmap, false, t)
}

func Eval(ljson, rjson string, hmap map[string]string, shouldBeEqual bool, t *testing.T) {

	h1 := hmap[ljson]
	h2 := hmap[rjson]

	areEqual := h1 == h2

	if areEqual != shouldBeEqual {

		if shouldBeEqual {
			t.Errorf(`hash(` + ljson + `) != hash(` + rjson + `) should be equal.")`)
		} else {
			t.Errorf(`hash(` + ljson + `) == hash(` + rjson + `) should be different.")`)
		}
	}
}
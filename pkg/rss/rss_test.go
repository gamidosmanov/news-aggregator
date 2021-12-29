package rss

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func Test_parse(t *testing.T) {
	absPath, _ := filepath.Abs("../../testdata/raw_test_feed.xml")
	rawFile, err := os.Open(absPath)
	if err != nil {
		t.Error(err)
	}
	raw, err := ioutil.ReadAll(rawFile)
	if err != nil {
		t.Error(err)
	}
	posts, err := parse(raw)
	if err != nil {
		t.Error(err)
	}
	if len(posts) == 0 {
		t.Fatal("Failed to parse data")
	}
	posts2, err := Parse("https://habr.com/ru/rss/hub/go/all/?fl=ru")
	if err != nil {
		t.Error(err)
	}
	if len(posts2) == 0 {
		t.Fatal("Failed to parse data")
	}
	_, err = Parse("bad_url")
	if err != nil {
		t.Log("Bad url is not parsed")
	}
}

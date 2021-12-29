package storage

import (
	"fmt"
	"testing"
	"time"
)

func TestDB(t *testing.T) {
	db, err := New()
	if err != nil {
		t.Error(err)
	}
	post := Post{
		Title:   "Test post",
		Content: "Test content",
		PubTime: time.Now().Unix(),
		Link:    fmt.Sprintf("test_url %d", time.Now().Unix()),
	}
	id, err := db.SavePost(post)
	if err != nil {
		t.Error(err)
	}
	posts, err := db.Posts(10)
	if err != nil {
		t.Error(err)
	}
	for _, p := range posts {
		if p.ID == id {
			t.Log(id)
			return
		}
	}
	t.Error("incorrect post")
}

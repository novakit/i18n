package i18n_test

import (
	"testing"

	"github.com/novakit/binfs"
	"github.com/novakit/i18n"
)

func TestSource_Get(t *testing.T) {
	root := &binfs.Node{}
	root.Load(&binfs.Chunk{
		Path: []string{"locales", "zh-CN.a.yml"},
		Data: []byte("Key1: Value1\nKey2:\n  Key22: Value22"),
	})
	root.Load(&binfs.Chunk{
		Path: []string{"locales", "zh-CN.b.yml"},
		Data: []byte("Key3:\n  Key31:\n    Key311: Value311"),
	})
	src := i18n.NewSource(root.Find("locales").FileSystem())
	src.Reload()
	if src.Get("zh-CN.Key3.Key31.Key311") != "Value311" {
		t.Error("failed 1")
	}
	if src.Get("zh-CN.Key2.Key22") != "Value22" {
		t.Error("failed 2")
	}
	if src.Get("zh-CN.Key1") != "Value1" {
		t.Error("failed 3")
	}
}

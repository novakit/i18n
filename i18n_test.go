package i18n_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/novakit/i18n"
	"github.com/novakit/nova"
	"github.com/novakit/router"
	"github.com/novakit/testkit"
	"github.com/novakit/view"
)

func TestI18n_All(t *testing.T) {
	n := nova.New()
	n.Use(i18n.Handler(i18n.Options{
		Directory:   "testdata/locales",
		Locales:     []string{"en-GB", "zh-CN"},
		LocaleNames: []string{"正宗伦敦腔", "中国话"},
	}))
	n.Use(view.Handler(view.Options{
		Directory: "testdata/views",
	}))
	router.Route(n).Get("/hello").Use(func(c *nova.Context) error {
		v := view.Extract(c)
		v.HTML("index")
		return nil
	})
	var req *http.Request
	var res *testkit.DummyResponse

	req, _ = http.NewRequest(http.MethodGet, "/hello", nil)
	req.Header.Set("Accept-Language", "zh")
	res = testkit.NewDummyResponse()
	n.ServeHTTP(res, req)

	if !strings.Contains(res.String(), "中国话") || !strings.Contains(res.String(), "zhHELLOworld") || !strings.Contains(res.String(), "zhworldhello") {
		t.Error("failed 1")
	}

	req, _ = http.NewRequest(http.MethodGet, "/hello", nil)
	req.Header.Set("Accept-Language", "zh")
	req.AddCookie(&http.Cookie{Name: "locale", Value: "en"})
	res = testkit.NewDummyResponse()
	n.ServeHTTP(res, req)

	if !strings.Contains(res.String(), "正宗伦敦腔") || !strings.Contains(res.String(), "enHELLOworld") || !strings.Contains(res.String(), "enworldhello") {
		t.Error("failed 2")
	}
	req, _ = http.NewRequest(http.MethodGet, "/hello?locale=zh", nil)
	req.Header.Set("Accept-Language", "zh")
	req.AddCookie(&http.Cookie{Name: "locale", Value: "en"})
	res = testkit.NewDummyResponse()
	n.ServeHTTP(res, req)

	if !strings.Contains(res.String(), "中国话") || !strings.Contains(res.String(), "zhHELLOworld") || !strings.Contains(res.String(), "zhworldhello") {
		t.Error("failed 3", res.String())
	}
}

func TestI18n_AllBinFS(t *testing.T) {
	n := nova.New()
	n.Use(i18n.Handler(i18n.Options{
		Directory:   "testdata/locales",
		Locales:     []string{"en-GB", "zh-CN"},
		LocaleNames: []string{"正宗伦敦腔", "中国话"},
		BinFS:       true,
	}))
	n.Use(view.Handler(view.Options{
		Directory: "testdata/views",
		BinFS:     true,
	}))
	router.Route(n).Get("/hello").Use(func(c *nova.Context) error {
		v := view.Extract(c)
		v.HTML("index")
		return nil
	})
	var req *http.Request
	var res *testkit.DummyResponse

	req, _ = http.NewRequest(http.MethodGet, "/hello", nil)
	req.Header.Set("Accept-Language", "zh")
	res = testkit.NewDummyResponse()
	n.ServeHTTP(res, req)

	if !strings.Contains(res.String(), "中国话") || !strings.Contains(res.String(), "zhHELLOworld") || !strings.Contains(res.String(), "zhworldhello") {
		t.Error("failed 1")
	}

	req, _ = http.NewRequest(http.MethodGet, "/hello", nil)
	req.Header.Set("Accept-Language", "zh")
	req.AddCookie(&http.Cookie{Name: "locale", Value: "en"})
	res = testkit.NewDummyResponse()
	n.ServeHTTP(res, req)

	if !strings.Contains(res.String(), "正宗伦敦腔") || !strings.Contains(res.String(), "enHELLOworld") || !strings.Contains(res.String(), "enworldhello") {
		t.Error("failed 2")
	}
	req, _ = http.NewRequest(http.MethodGet, "/hello?locale=zh", nil)
	req.Header.Set("Accept-Language", "zh")
	req.AddCookie(&http.Cookie{Name: "locale", Value: "en"})
	res = testkit.NewDummyResponse()
	n.ServeHTTP(res, req)

	if !strings.Contains(res.String(), "中国话") || !strings.Contains(res.String(), "zhHELLOworld") || !strings.Contains(res.String(), "zhworldhello") {
		t.Error("failed 3", res.String())
	}
}

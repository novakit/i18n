package i18n // import "github.com/novakit/i18n"

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/novakit/binfs"
	"github.com/novakit/nova"
	"golang.org/x/text/language"
)

// ContextKey key in nova.Context, same value is hard coded in "view" package
const ContextKey = "_i18n"

// Options i18n options
type Options struct {
	// Directory directory contains i18n files, default "locales"
	Directory string
	// BinFS using binfs filesystem
	BinFS bool
	// Locales locales, first is default
	Locales []string
	// LocaleNames locale names
	LocaleNames []string
	// CookieName cookie name for locale overriding, default to "lang"
	CookieName string
	// QueryName query name for locale overriding, default to "lang"
	QueryName string
}

func sanitizeOptions(opts ...Options) (opt Options) {
	if len(opts) > 0 {
		opt = opts[0]
	}
	if len(opt.Directory) == 0 {
		opt.Directory = "locales"
	}
	if len(opt.Locales) != len(opt.LocaleNames) {
		panic("i18n.Options len(opt.Locales) != len(opt.LocaleNames)")
	}
	if len(opt.Locales) == 0 {
		opt.Locales = []string{"en-US"}
		opt.LocaleNames = []string{"English"}
	}
	if len(opt.QueryName) == 0 {
		opt.QueryName = "locale"
	}
	if len(opt.CookieName) == 0 {
		opt.CookieName = "locale"
	}
	return
}

func createMatcher(opt Options) language.Matcher {
	var tags []language.Tag
	for _, l := range opt.Locales {
		tags = append(tags, language.MustParse(l))
	}
	return language.NewMatcher(tags)
}

// I18n the i18n instance
type I18n struct {
	// Source source of locales
	Source *Source
	// Locale active locale
	Locale string
	// LocaleName active locale name
	LocaleName string
}

// T find translation with key and render with optional arguments
func (n *I18n) T(key string, args ...string) string {
	fk := n.Locale + "." + key
	v := n.Source.Get(fk)
	if len(v) == 0 {
		return "[i18n missing: " + fk + "]"
	}
	if len(args) > 0 {
		for i, a := range args {
			v = strings.Replace(v, fmt.Sprintf("{{%d}}", i+1), a, -1)
		}
	}
	return v
}

// Handler create nova.HandlerFunc
func Handler(opts ...Options) nova.HandlerFunc {
	opt := sanitizeOptions(opts...)
	mch := createMatcher(opt)
	// build filesystem
	var fs http.FileSystem
	if opt.BinFS {
		var n *binfs.Node
		n = binfs.Find(strings.Split(opt.Directory, "/")...)
		if n == nil {
			panic(fmt.Errorf("cann't find directory " + opt.Directory + " from binfs"))
		}
		fs = n.FileSystem()
	} else {
		fs = http.Dir(opt.Directory)
	}
	// build source
	src := NewSource(fs)
	return func(c *nova.Context) error {
		if c.Env.IsDevelopment() {
			src.Reload()
		}
		locales := make([]string, 0, 3)
		// extract query
		if qr := c.Req.URL.Query(); qr != nil {
			if q := qr.Get(opt.QueryName); len(q) > 0 {
				locales = append(locales, q)
			}
		}
		// extract cookie
		if cookie, err := c.Req.Cookie(opt.CookieName); err == nil && cookie != nil {
			locales = append(locales, cookie.Value)
		}
		// extract Accept-Language
		if a := c.Req.Header.Get("Accept-Language"); len(a) > 0 {
			locales = append(locales, a)
		}
		locale, _ := language.MatchStrings(mch, locales...)
		// create i18n
		n := &I18n{Source: src, Locale: locale.String()}
		for i, l := range opt.Locales {
			if l == n.Locale {
				n.LocaleName = opt.LocaleNames[i]
			}
		}
		c.Values[ContextKey] = n
		// invoke next handler
		c.Next()
		return nil
	}
}

// Extract extract I18n from nova.Context
func Extract(c *nova.Context) (n *I18n) {
	n, _ = c.Values[ContextKey].(*I18n)
	return
}

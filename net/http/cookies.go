package http

import (
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/dgrr/cookiejar"
	"github.com/sudosz/go-utils/bytes"
	"github.com/valyala/fasthttp"
)

func isCookieDomainValid(cookieDomain, urlHost string) bool {
	// Normalize domain and host for comparison
	cookieDomain = strings.ToLower(cookieDomain)
	urlHost = strings.ToLower(urlHost)

	// If cookie domain is equal to URL host, it's valid
	if cookieDomain == urlHost {
		return true
	}

	// If cookie domain is a suffix of URL host, preceded by a '.'
	if strings.HasSuffix(urlHost, "."+cookieDomain) {
		return true
	}

	return false
}

// HTTPCookieJar wraps cookiejar.CookieJar to provide a thread-safe cookie jar
// that can convert between fasthttp and net/http cookies. It implements the
// http.CookieJar interface while internally using fasthttp cookies for better
// performance.
type HTTPCookieJar struct {
	*cookiejar.CookieJar
	mux *sync.RWMutex
}

var cookieJarPool = &sync.Pool{
	New: func() any {
		return &HTTPCookieJar{
			CookieJar: cookiejar.AcquireCookieJar(),
			mux:       &sync.RWMutex{},
		}
	},
}

func AcquireCookieJar() *HTTPCookieJar {
	return cookieJarPool.Get().(*HTTPCookieJar)
}

func ReleaseCookieJar(cj *HTTPCookieJar) {
	cj.CookieJar.Release()
	cookieJarPool.Put(cj)
}

func (cj *HTTPCookieJar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	cj.mux.Lock()
	defer cj.mux.Unlock()
	for _, c := range cookies {
		cj.Put(HTTPCookie2FastHTTPCookie(c))
	}
}

func (cj *HTTPCookieJar) Cookies(u *url.URL) []*http.Cookie {
	cj.mux.RLock()
	defer cj.mux.RUnlock()

	cs := make([]*http.Cookie, 0)

	for _, cookie := range *(cj.CookieJar) {
		d, p := bytes.B2s(cookie.Domain()), bytes.B2s(cookie.Path())
		if d != "" && isCookieDomainValid(d, u.Hostname()) {
			continue
		}
		if p != "" && p != "/" && p != u.Path {
			continue
		}
		sameSite := http.SameSiteDefaultMode
		switch cookie.SameSite() {
		case fasthttp.CookieSameSiteLaxMode:
			sameSite = http.SameSiteLaxMode
		case fasthttp.CookieSameSiteStrictMode:
			sameSite = http.SameSiteStrictMode
		case fasthttp.CookieSameSiteNoneMode:
			sameSite = http.SameSiteNoneMode
		}
		cs = append(cs, &http.Cookie{
			Name:     bytes.B2s(cookie.Key()),
			Value:    url.QueryEscape(bytes.B2s(cookie.Value())),
			Domain:   bytes.B2s(cookie.Domain()),
			Path:     bytes.B2s(cookie.Path()),
			Expires:  cookie.Expire(),
			MaxAge:   cookie.MaxAge(),
			Secure:   cookie.Secure(),
			HttpOnly: cookie.HTTPOnly(),
			SameSite: sameSite,
			Raw:      cookie.String(),
		})

	}

	return cs
}

func HTTPCookie2FastHTTPCookie(c *http.Cookie) *fasthttp.Cookie {

	cookie := fasthttp.AcquireCookie()

	cookie.SetKey(c.Name)
	cookie.SetValue(c.Value)
	cookie.SetPath(c.Path)
	cookie.SetDomain(c.Domain)
	cookie.SetExpire(c.Expires)
	cookie.SetMaxAge(c.MaxAge)
	cookie.SetSecure(c.Secure)
	cookie.SetHTTPOnly(c.HttpOnly)

	sameSite := fasthttp.CookieSameSiteDefaultMode
	switch c.SameSite {
	case http.SameSiteLaxMode:
		sameSite = fasthttp.CookieSameSiteLaxMode
	case http.SameSiteStrictMode:
		sameSite = fasthttp.CookieSameSiteStrictMode
	case http.SameSiteNoneMode:
		sameSite = fasthttp.CookieSameSiteNoneMode
	}

	cookie.SetSameSite(sameSite)
	return cookie
}

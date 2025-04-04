package httputils

import (
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/dgrr/cookiejar"
	bytesutils "github.com/sudosz/go-utils/bytes"
	"github.com/valyala/fasthttp"
)

// isCookieDomainValid checks if the cookie domain is valid for the URL host.
// Optimization: Simple string comparison with early returns.
func isCookieDomainValid(cookieDomain, urlHost string) bool {
	cookieDomain = strings.ToLower(cookieDomain)
	urlHost = strings.ToLower(urlHost)
	if cookieDomain == urlHost {
		return true
	}
	if strings.HasSuffix(urlHost, "."+cookieDomain) {
		return true
	}
	return false
}

// HTTPCookieJar provides a thread-safe cookie jar bridging fasthttp and net/http.
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

// AcquireCookieJar retrieves a cookie jar from the pool.
// Optimization: Pooling reduces allocations.
func AcquireCookieJar() *HTTPCookieJar {
	return cookieJarPool.Get().(*HTTPCookieJar)
}

// ReleaseCookieJar returns a cookie jar to the pool.
// Optimization: Enables reuse of cookie jars.
func ReleaseCookieJar(cj *HTTPCookieJar) {
	cj.CookieJar.Release()
	cookieJarPool.Put(cj)
}

// SetCookies sets cookies for the given URL, thread-safely.
// Optimization: Uses write lock for safe concurrent modification.
func (cj *HTTPCookieJar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	cj.mux.Lock()
	defer cj.mux.Unlock()
	for _, c := range cookies {
		cj.Put(HTTPCookie2FastHTTPCookie(c))
	}
}

// Cookies retrieves cookies for the given URL, thread-safely.
// Optimization: Uses read lock for concurrent read safety.
func (cj *HTTPCookieJar) Cookies(u *url.URL) []*http.Cookie {
	cj.mux.RLock()
	defer cj.mux.RUnlock()
	cs := make([]*http.Cookie, 0)
	for _, cookie := range *(cj.CookieJar) {
		d, p := bytesutils.B2s(cookie.Domain()), bytesutils.B2s(cookie.Path())
		if d != "" && !isCookieDomainValid(d, u.Hostname()) {
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
			Name:     bytesutils.B2s(cookie.Key()),
			Value:    url.QueryEscape(bytesutils.B2s(cookie.Value())),
			Domain:   bytesutils.B2s(cookie.Domain()),
			Path:     bytesutils.B2s(cookie.Path()),
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

// HTTPCookie2FastHTTPCookie converts an http.Cookie to a fasthttp.Cookie.
// Optimization: Uses fasthttp pooling for cookie objects.
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

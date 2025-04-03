package useragent

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/malisit/kolpa"
	"github.com/mileusna/useragent"
	"github.com/sudosz/go-utils/ints"
)

var random = rand.New(rand.NewSource(time.Now().UnixNano()))
var (
	kpg   = kolpa.C()
	funcs = []func() string{
		kpg.Chrome,
		kpg.Firefox,
		kpg.Safari,
		kpg.Opera,
	}
	funcsLen = int32(len(funcs))
)

// GetRandomUserAgentKolpa returns a random user agent string using kolpa.
// Optimization: Pre-seeded random source and fixed function list.
func GetRandomUserAgentKolpa() string {
	return funcs[random.Int31n(funcsLen)]()
}

// GetRandomUserAgentJSON returns a random user agent from FakeUserAgents.
// Optimization: Efficient random selection from pre-defined list.
func GetRandomUserAgentJSON() string {
	return FakeUserAgents[random.Int31n(int32(len(FakeUserAgents)))]
}

// GenerateUserAgent generates a random user agent string (alias for JSON method).
// Optimization: Direct call to optimized GetRandomUserAgentJSON.
func GenerateUserAgent() string {
	return GetRandomUserAgentJSON()
}

// BrowserUserAgent represents parsed user agent information.
type BrowserUserAgent struct {
	Platform        string
	PlatformVersion string
	Architecture    string
	Model           string
	IsMobile        bool
	Brands          []*BrowserUserAgentBrand
}

func (u BrowserUserAgent) String() string {
	return fmt.Sprintf("<BrowserUserAgent platform=%q platform_version=%q architecture=%q model=%q brands=%v>", u.Platform, u.PlatformVersion, u.Architecture, u.Model, u.Brands)
}

// BrowserUserAgentBrand represents a browser brand and version.
type BrowserUserAgentBrand struct {
	Brand   string
	Version string
}

func (u BrowserUserAgentBrand) String() string {
	return fmt.Sprintf("<Brand name=%q version=%q>", u.Brand, u.Version)
}

// ParseUserAgentForBrowser parses a user agent string into BrowserUserAgent.
// Optimization: Efficient parsing with minimal allocations.
func ParseUserAgentForBrowser(userAgent string) *BrowserUserAgent {
	ua := useragent.Parse(userAgent)
	platform := ua.OS
	platformVersion := ints.Itoa(ua.OSVersionNo.Major)
	arch := ""
	model := ""
	isMobile := ua.Mobile
	major := ints.Itoa(ua.VersionNo.Major)
	brands := []*BrowserUserAgentBrand{
		{Brand: "Chromium", Version: ints.Itoa(117 + random.Intn(5))},
		{Brand: "Not(A:Brand", Version: "24"},
		{Brand: ua.Name, Version: major},
	}
	return &BrowserUserAgent{
		Platform:        platform,
		PlatformVersion: platformVersion,
		Architecture:    arch,
		Model:           model,
		IsMobile:        isMobile,
		Brands:          brands,
	}
}

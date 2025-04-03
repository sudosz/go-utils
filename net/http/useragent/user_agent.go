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

func GetRandomUserAgentKolpa() string {
	return funcs[random.Int31n(funcsLen)]()
}

func GetRandomUserAgentJSON() string {
	return FakeUserAgents[random.Int31n(int32(len(FakeUserAgents)))]
}

func GenerateUserAgent() string {
	return GetRandomUserAgentJSON()
}

// BrowserUserAgent represents the user agent information for a browser.
// It contains details about the platform, architecture, model, and brands.
type BrowserUserAgent struct {
	// Platform is the operating system platform.
	Platform string
	// PlatformVersion is the version of the operating system platform.
	PlatformVersion string
	// Architecture is the hardware architecture.
	Architecture string
	// Model is the device model.
	Model string
	// IsMobile indicates whether the user agent is from a mobile device.
	IsMobile bool
	// Brands is a list of browser brands and their versions.
	Brands []*BrowserUserAgentBrand
}

func (u BrowserUserAgent) String() string {
	return fmt.Sprintf("<BrowserUserAgent platform=%q platform_version=%q architecture=%q model=%q brands=%v>", u.Platform, u.PlatformVersion, u.Architecture, u.Model, u.Brands)
}

// BrowserUserAgentBrand represents a brand and version of a browser user agent.
type BrowserUserAgentBrand struct {
	Brand   string // The name of the browser brand.
	Version string // The version of the browser brand.
}

func (u BrowserUserAgentBrand) String() string {
	return fmt.Sprintf("<Brand name=%q version=%q>", u.Brand, u.Version)
}

func ParseUserAgentForBrowser(userAgent string) *BrowserUserAgent {
	ua := useragent.Parse(userAgent)
	platform := ua.OS
	platformVersion := ints.Itoa(ua.OSVersionNo.Major)
	arch := ""
	model := ""
	isMobile := ua.Mobile
	major := ints.Itoa(ua.VersionNo.Major)
	_ = major
	brands := []*BrowserUserAgentBrand{
		{
			Brand:   "Chromium",
			Version: ints.Itoa(117 + random.Intn(5)),
		},
		{
			Brand:   "Not(A:Brand",
			Version: "24",
		},
		{
			Brand:   ua.Name,
			Version: major,
		},
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

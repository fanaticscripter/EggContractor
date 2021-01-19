package web

import (
	"encoding/json"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/dustin/go-humanize"
	"github.com/fanaticscripter/EggContractor/api"
	"github.com/fanaticscripter/EggContractor/util"
)

const (
	_devStaticDir                 = "public"
	_prodStaticDir                = "static"
	_devStaticCacheControlHeader  = "no-cache"
	_prodStaticCacheControlHeader = "max-age=2592000"
)

var (
	_staticMap   map[string]string
	_iconPathSet map[string]struct{}
)

type ServerOptions struct {
	BindAddr string
	Dev      bool
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func Serve(opts ServerOptions) {
	if !opts.Dev {
		if err := preloadStaticManifests(); err != nil {
			log.Fatal(err)
		}
	}
	if err := preloadIcons(); err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	// For a self-hosted web app, let's leave debug mode enabled even in production.
	e.Debug = true
	e.HideBanner = true
	e.Renderer = loadTemplates(opts)
	e.IPExtractor = echo.ExtractIPFromXFFHeader()

	e.Pre(middleware.AddTrailingSlashWithConfig(middleware.TrailingSlashConfig{
		Skipper: func(c echo.Context) bool {
			path := c.Request().URL.Path
			if path == "/favicon.ico" {
				return true
			}
			if strings.HasPrefix(path, "/static/") {
				return true
			}
			return false
		},
		RedirectCode: http.StatusFound,
	}))
	// Common Log Format
	logFormat := `${remote_ip} - - [${time_custom}] "${method} ${path} ${protocol}" ${status} ${bytes_out}`
	customTimeFormat := "2/Jan/2006:15:04:05 -0700"
	if !opts.Dev {
		logFormat += ` "${referer}" "${user_agent}"`
	}
	logFormat += "\n"
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format:           logFormat,
		CustomTimeFormat: customTimeFormat,
	}))

	e.GET("/", indexHandler)
	e.GET("/peek/:contractId/:code/", peekHandler)
	e.GET("/peeked/", peekedHandler)
	e.GET("/events/", eventsHandler)

	staticDir := _prodStaticDir
	staticCacheControlHeader := _prodStaticCacheControlHeader
	if opts.Dev {
		staticDir = _devStaticDir
		staticCacheControlHeader = _devStaticCacheControlHeader
	}
	e.Group("/static", func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("Cache-Control", staticCacheControlHeader)
			return next(c)
		}
	}).Static("/", staticDir)
	e.File("/favicon.ico", "static/favicon.ico")

	go dbPeekedWorker()

	log.Fatal(e.Start(opts.BindAddr))
}

func preloadStaticManifests() error {
	_staticMap = make(map[string]string)
	manifestPaths, _ := filepath.Glob("static/manifest*.json")
	for _, path := range manifestPaths {
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		var payload map[string]string
		if err := json.Unmarshal(data, &payload); err != nil {
			return errors.Wrapf(err, "failed to JSON decode %s", path)
		}
		for k, v := range payload {
			vv, exists := _staticMap[k]
			if exists {
				return errors.Errorf("duplicate key %#v in manifests mapped to %#v and %#v", k, vv, v)
			}
			_staticMap[k] = v
		}
	}
	return nil
}

func preloadIcons() error {
	_iconPathSet = make(map[string]struct{})
	iconPaths, _ := filepath.Glob("static/egginc/*.png")
	additional, _ := filepath.Glob("static/egginc-extras/*.png")
	iconPaths = append(iconPaths, additional...)
	for _, p := range iconPaths {
		rp := strings.TrimPrefix(p, "static/")
		_iconPathSet[rp] = struct{}{}
	}
	return nil
}

func loadTemplates(opts ServerOptions) *Template {
	staticAssetURL := prodStaticAssetURL
	if opts.Dev {
		staticAssetURL = devStaticAssetURL
	}
	return &Template{
		templates: template.Must(template.New("").Funcs(template.FuncMap{
			"css":              func(s string) template.CSS { return template.CSS(s) },
			"eggiconpath":      eggIconPath,
			"eggname":          func(e api.EggType) string { return e.Display() },
			"eggvalue":         func(e api.EggType) string { return e.ValueDisplay() },
			"fmtcountdown":     util.FormatCountdown,
			"fmtdate":          util.FormatDate,
			"fmtdatecasual":    util.FormatDateCasual,
			"fmtdatetime":      util.FormatDatetime,
			"fmtduration":      util.FormatDuration,
			"fmtdurationGe0":   util.FormatDurationNonNegative,
			"fmttimecasual":    util.FormatTimeCasual,
			"fmtpercent":       util.FormatPercentage,
			"hasactivitystats": hasActivityStats,
			"increment":        func(x int) int { return x + 1 },
			"iseven":           func(x int) bool { return x%2 == 0 },
			"islastindex":      func(index int, length int) bool { return index == length-1 },
			"isodd":            func(x int) bool { return x%2 == 1 },
			"json":             marshalJSON,
			"members":          getMemberPayloads,
			"numfmt":           util.Numfmt,
			"numfmtWhole":      util.NumfmtWhole,
			"reltime":          humanize.Time,
			"rewardiconpath":   rewardIconPath,
			"static":           staticAssetURL,
		}).ParseGlob("templates/*/*.html")),
	}
}

func devStaticAssetURL(pth string) string {
	return path.Join("/static", pth)
}

func prodStaticAssetURL(pth string) string {
	realpath, ok := _staticMap[pth]
	if ok {
		return path.Join("/static", realpath)
	}
	return path.Join("/static", pth)
}

func eggIconPath(e api.EggType) string {
	var path string
	switch e {
	case api.EggType_EDIBLE:
		path = "egginc/egg_edible.png"
	case api.EggType_SUPERFOOD:
		path = "egginc/egg_superfood.png"
	case api.EggType_MEDICAL:
		path = "egginc/egg_medical2.png"
	case api.EggType_ROCKET_FUEL:
		path = "egginc/egg_rocketfuel.png"
	case api.EggType_SUPER_MATERIAL:
		path = "egginc/egg_supermaterial.png"
	case api.EggType_FUSION:
		path = "egginc/egg_fusion.png"
	case api.EggType_QUANTUM:
		path = "egginc/egg_quantum.png"
	case api.EggType_IMMORTALITY:
		path = "egginc/egg_immortality.png"
	case api.EggType_TACHYON:
		path = "egginc/egg_tachyon.png"
	case api.EggType_GRAVITON:
		path = "egginc/egg_graviton.png"
	case api.EggType_DILITHIUM:
		path = "egginc/egg_dilithium.png"
	case api.EggType_PRODIGY:
		path = "egginc/egg_prodigy.png"
	case api.EggType_TERRAFORM:
		path = "egginc/egg_terraform.png"
	case api.EggType_ANTIMATTER:
		path = "egginc/egg_antimatter.png"
	case api.EggType_DARK_MATTER:
	case api.EggType_AI:
		path = "egginc/egg_ai.png"
	case api.EggType_NEBULA:
		path = "egginc/egg_vision.png"
	case api.EggType_UNIVERSE:
		path = "egginc/egg_universe.png"
	case api.EggType_ENLIGHTENMENT:
		path = "egginc/egg_enlightenment.png"
	// Contract-only eggs.
	case api.EggType_CHOCOLATE:
		path = "egginc/egg_chocolate.png"
	case api.EggType_EASTER:
		path = "egginc/egg_easter.png"
	case api.EggType_WATERBALLOON:
		path = "egginc/egg_waterballoon.png"
	case api.EggType_FIREWORK:
		path = "egginc/egg_firework.png"
	case api.EggType_PUMPKIN:
		path = "egginc/egg_pumpkin.png"
	}
	_, ok := _iconPathSet[path]
	if ok {
		return path
	}
	path = "egginc/egg_unknown.png"
	_, ok = _iconPathSet[path]
	if ok {
		return path
	}
	return "egginc/icon_help.png"
}

func rewardIconPath(r *api.Reward) string {
	var path string
	switch r.Type {
	case api.RewardType_GOLDEN_EGG:
		path = "egginc-extras/icon_golden_egg.png"
	case api.RewardType_SOUL_EGG:
		path = "egginc/egg_soul.png"
	case api.RewardType_PROPHECY_EGG:
		path = "egginc/egg_of_prophecy.png"
	case api.RewardType_EPIC_RESEARCH:
		path = "egginc/r_icon_" + r.Name + ".png"
	case api.RewardType_PIGGY_GOLDEN_EGG:
		path = "egginc-extras/icon_piggy_golden_egg.png"
	case api.RewardType_PIGGY_MULTIPLY:
		path = "egginc-extras/icon_piggy_multiply.png"
	case api.RewardType_PIGGY_LEVEL_UP:
		path = "egginc-extras/icon_piggy_level_up.png"
	case api.RewardType_BOOST:
		path = "egginc/b_icon_" + r.Name + ".png"
	}
	_, ok := _iconPathSet[path]
	if ok {
		return path
	}
	return "egginc/icon_help.png"
}

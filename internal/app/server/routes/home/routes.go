// Package home contains the route handlers related to the app's home screen.
package home

import (
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/sargassum-world/godest"

	sh "github.com/openUC2/machine-portal/internal/app/server/handling"
	"github.com/openUC2/machine-portal/internal/clients/machinename"
)

type Handlers struct {
	r   godest.TemplateRenderer
	mnc *machinename.Client
}

func New(r godest.TemplateRenderer, mnc *machinename.Client) *Handlers {
	return &Handlers{
		r:   r,
		mnc: mnc,
	}
}

func (h *Handlers) Register(er godest.EchoRouter) {
	er.GET("/", h.HandleHomeGet())
}

type HomeViewData struct {
	Hostname    string
	Port        string
	MachineName string
}

func getHomeViewData(host, machineName string) (vd HomeViewData, err error) {
	split := strings.Split(host, ":")
	const expectedComponents = 2
	if len(split) > expectedComponents {
		return HomeViewData{}, errors.Errorf(
			"unable to split host '%s' into a hostname and a port", host,
		)
	}
	vd.Hostname = split[0]
	if len(split) == expectedComponents {
		vd.Port = split[expectedComponents-1]
	}
	vd.MachineName = machineName
	return vd, nil
}

func (h *Handlers) HandleHomeGet() echo.HandlerFunc {
	t := "home/index.page.tmpl"
	h.r.MustHave(t)
	ta := "home/index.advanced.page.tmpl"
	h.r.MustHave(ta)
	return func(c echo.Context) error {
		// Parse params
		mode := c.QueryParam("mode")

		// Run queries
		machineName, err := h.mnc.GetName()
		if err != nil {
			return err
		}
		vd, err := getHomeViewData(c.Request().Host, machineName)
		if err != nil {
			return err
		}
		// Produce output
		switch mode {
		default:
			return h.r.CacheablePage(c.Response(), c.Request(), t, vd, struct{}{})
		case sh.ViewModeAdvanced:
			return h.r.CacheablePage(c.Response(), c.Request(), ta, vd, struct{}{})
		}
	}
}

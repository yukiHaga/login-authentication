package controller

import (
	"time"

	"github.com/leekchan/timeutil"
	"github.com/yukiHaga/web_server/src/pkg/henagin/http"
	"github.com/yukiHaga/web_server/src/pkg/henagin/view"
)

type Now struct{}

func NewNow() *Now {
	return &Now{}
}

func (c *Now) Action(request *http.Request) *http.Response {
	t := time.Now()
	formatedTime := timeutil.Strftime(&t, "%a, %d %b %Y %H:%M:%S")
	body := view.Render("now.html", formatedTime)

	return http.NewResponse(
		http.VersionsFor11,
		http.StatusSuccessCode,
		http.StatusReasonOk,
		request.TargetPath,
		body,
	)
}

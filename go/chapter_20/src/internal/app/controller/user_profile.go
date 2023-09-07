package controller

import (
	"github.com/yukiHaga/web_server/src/pkg/henagin/http"
	"github.com/yukiHaga/web_server/src/pkg/henagin/view"
)

type UserProfile struct{}

func NewUserProfile() *UserProfile {
	return &UserProfile{}
}

func (c *UserProfile) Action(request *http.Request) *http.Response {
	id := request.Params["id"]
	body := view.Render("user_profile.html", id)

	return http.NewResponse(
		http.VersionsFor11,
		http.StatusSuccessCode,
		http.StatusReasonOk,
		request.TargetPath,
		body,
	)
}

package controller

import (
	"log"
	"net/url"

	"github.com/yukiHaga/web_server/src/internal/app/model"
	"github.com/yukiHaga/web_server/src/internal/app/view"
	"github.com/yukiHaga/web_server/src/pkg/henagin/http"
)

type MyPage struct{}

func NewMyPage() *MyPage {
	return &MyPage{}
}

func (c *MyPage) Action(request *http.Request) *http.Response {
	if cookie, isThere := request.GetCookieByName("user_id"); isThere {
		id, _ := url.QueryUnescape(cookie.Value)
		user, err := model.FindUserById(id)
		if err != nil {
			log.Printf("fail to find error: %v", err)
			response := http.NewResponse(
				http.VersionsFor11,
				http.StatusRedirectCode,
				http.StatusReasonRedirect,
				request.TargetPath,
				[]byte{},
			)
			response.SetHeader("Location", "/sign_up")
			return response
		}
		body := view.Render("my_page.html", user.Name)
		return http.NewResponse(
			http.VersionsFor11,
			http.StatusSuccessCode,
			http.StatusReasonOk,
			request.TargetPath,
			body,
		)
	} else {
		response := http.NewResponse(
			http.VersionsFor11,
			http.StatusRedirectCode,
			http.StatusReasonRedirect,
			request.TargetPath,
			[]byte{},
		)
		response.SetHeader("Location", "/sign_up")
		return response
	}
}

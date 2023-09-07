package controller

import (
	"net/url"

	"github.com/yukiHaga/web_server/src/internal/app/model"
	"github.com/yukiHaga/web_server/src/pkg/henagin/http"
	"github.com/yukiHaga/web_server/src/pkg/henagin/view"
)

type MyPage struct{}

func NewMyPage() *MyPage {
	return &MyPage{}
}

func (c *MyPage) Action(request *http.Request) *http.Response {
	// クッキーのユーザーがいるかどうかのミドルウェアを通過してきたので、クッキーは必ず存在する
	// ここの処理がシンプルになる
	cookie, _ := request.GetCookieByName("user_id")
	id, _ := url.QueryUnescape(cookie.Value)
	user, _ := model.FindUserById(id)
	body := view.Render("my_page.html", user.Name)
	return http.NewResponse(
		http.VersionsFor11,
		http.StatusSuccessCode,
		http.StatusReasonOk,
		request.TargetPath,
		body,
	)
}

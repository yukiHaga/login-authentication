package middleware

import (
	"log"

	"github.com/yukiHaga/web_server/src/internal/app/auth"
	"github.com/yukiHaga/web_server/src/internal/app/controller"
	"github.com/yukiHaga/web_server/src/internal/app/model"
	"github.com/yukiHaga/web_server/src/pkg/henagin/http"
)

type CheckLoginController struct {
	nextAction func(request *http.Request) *http.Response
}

// ログインユーザー専用のページをリクエストした場合、ユーザーidのクッキーがないならログインページにリダイレクトさせる
func (c CheckLoginController) Action(request *http.Request) *http.Response {
	if cookie, isThere := request.GetCookieByName("session_id"); isThere {
		session := auth.NewSession()
		userId, _ := session.Load(cookie.Value)
		_, err := model.FindUserById(userId)
		// user_idクッキーはあるけど不正な場合
		// ログインページにリダイレクトさせる
		if err != nil {
			log.Printf("fail to find error: %v", err)
			response := http.NewResponse(
				http.VersionsFor11,
				http.StatusRedirectCode,
				http.StatusReasonRedirect,
				request.TargetPath,
				[]byte{},
			)
			response.SetRedirectHeader("/login")
			return response

		} else {
			// user_idクッキーがあってかつ有効な場合
			// ログインユーザー専用ページのアクションを呼び出す
			c.nextAction(request)
		}
	} else {
		// user_idクッキーがそもそもない場合
		// ログインページにリダイレクトさせる
		response := http.NewResponse(
			http.VersionsFor11,
			http.StatusRedirectCode,
			http.StatusReasonRedirect,
			request.TargetPath,
			[]byte{},
		)
		response.SetRedirectHeader("/login")
		return response
	}
	return c.nextAction(request)
}

// ミドルウェア
// / ダミーコントローラとダミーアクションを作って、ダミーアクションの中で元々のコントローラのアクションを呼び出して、最終的にレスポンスを返せばOK
// goではメソッドを書き換えるのはできなかった
func CheckLogin(c controller.Controller) controller.Controller {
	return CheckLoginController{nextAction: c.Action}
}

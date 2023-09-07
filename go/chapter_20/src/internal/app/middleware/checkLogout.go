package middleware

import (
	"log"

	"github.com/yukiHaga/web_server/src/internal/app/auth"
	"github.com/yukiHaga/web_server/src/internal/app/controller"
	"github.com/yukiHaga/web_server/src/internal/app/model"
	"github.com/yukiHaga/web_server/src/pkg/henagin/http"
)

type CheckLogoutController struct {
	nextAction func(request *http.Request) *http.Response
}

// このミドルウェアは/loginと/sign_upのエンドポイントで使う
// ログインユーザーがログインページをやサインアップページリクエストした場合、ユーザーidのクッキーがあるなら、マイページへリダイレクトさせる
func (c CheckLogoutController) Action(request *http.Request) *http.Response {
	if cookie, isThere := request.GetCookieByName("session_id"); isThere {
		session := auth.NewSession()
		userId, _ := session.Load(cookie.Value)
		_, err := model.FindUserById(userId)
		// user_idクッキーはあるけど不正な場合(ユーザーが存在しない)
		// そのまま/loginと/sing_upのアクションにリクエストオブジェクトを渡す
		if err != nil {
			log.Printf("fail to find error: %v", err)
			return c.nextAction(request)
		} else {
			// user_idクッキーがあってかつ有効な場合
			// マイページにリダイレクトさせる
			response := http.NewResponse(
				http.VersionsFor11,
				http.StatusRedirectCode,
				http.StatusReasonRedirect,
				request.TargetPath,
				[]byte{},
			)
			response.SetRedirectHeader("/mypage")
			return response
		}
	} else {
		// user_idクッキーがそもそもない場合
		// そのまま/loginと/sign_upのアクションにリクエストオブジェクトを渡す
		return c.nextAction(request)
	}
}

// ミドルウェア
// / ダミーコントローラとダミーアクションを作って、ダミーアクションの中で元々のコントローラのアクションを呼び出して、最終的にレスポンスを返せばOK
// goではメソッドを書き換えるのはできなかった
func CheckLogout(c controller.Controller) controller.Controller {
	return CheckLogoutController{nextAction: c.Action}
}

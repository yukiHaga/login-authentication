package controller

import (
	"fmt"
	"log"
	"net/url"

	"github.com/yukiHaga/web_server/src/internal/app/auth"
	"github.com/yukiHaga/web_server/src/internal/app/model"
	"github.com/yukiHaga/web_server/src/pkg/henagin/http"
	"github.com/yukiHaga/web_server/src/pkg/henagin/view"
)

type SignUp struct{}

func NewSignUp() *SignUp {
	return &SignUp{}
}

func (c *SignUp) Action(request *http.Request) *http.Response {
	var statusCode string
	var reasonPhrase string
	var body []byte
	cookieHeaders := map[string]string{}

	if request.Method == http.Get {
		// ミドルウェアを通ってきたので、クッキーが存在しないことは確定
		// ゆえにクッキーを取得する処理は書かなくて良い
		// 何も埋め込んでいなないなら、普通にファイルを開くのか。
		body = view.Render("sign_up_form.html")
		statusCode = http.StatusSuccessCode
		reasonPhrase = http.StatusReasonOk
	} else if request.Method == http.Post {
		// ユーザー登録をここでする
		// emailがユニークであることを確認する。あとパスワードがちゃんと一致しているか
		// okなら、ユーザーデータをインサートして、クッキーにuser_idを入れる(本当はあんま良くない。セッションidにした方が良い)、
		decodedBody, _ := url.QueryUnescape(string(request.Body))
		// 正規表現使わなくて済む
		values, _ := url.ParseQuery(decodedBody)
		name := values.Get("name")
		email := values.Get("email")
		password := values.Get("password")
		passwordConfirmation := values.Get("password_confirmation")

		user := model.NewUser(name, email)

		if err := user.SignUp(password, passwordConfirmation); err != nil {
			log.Printf("fail to save user: %v\n", err)
			body = view.Render("sign_up_form.html")
			statusCode = http.StatusInternalServerErrorCode
			reasonPhrase = http.StatusReasonInternalServerError
		} else {
			// サインアップに成功した
			statusCode = http.StatusRedirectCode
			reasonPhrase = http.StatusReasonRedirect
			session := auth.NewSession()
			fmt.Println("サインアップ成功時のユーザーid", user.Id)
			sessionId, _ := session.Save(user.Id)
			cookieHeaders["session_id"] = string(sessionId)
		}
	}

	// headerがオプショナル引数ならどんなだけ楽か。
	response := http.NewResponse(
		http.VersionsFor11,
		statusCode,
		reasonPhrase,
		request.TargetPath,
		body,
	)

	if request.Method == http.Post && statusCode == http.StatusRedirectCode {
		for key, value := range cookieHeaders {
			response.SetCookieHeader(key, value)
		}
		response.SetRedirectHeader("/mypage")
	}

	return response
}

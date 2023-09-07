package controller

import (
	"log"
	"net/url"

	"github.com/yukiHaga/web_server/src/internal/app/auth"
	"github.com/yukiHaga/web_server/src/internal/app/model"
	"github.com/yukiHaga/web_server/src/pkg/henagin/http"
	"github.com/yukiHaga/web_server/src/pkg/henagin/view"
)

type Login struct{}

func NewLogin() *Login {
	return &Login{}
}

func (c *Login) Action(request *http.Request) *http.Response {
	var statusCode string
	var reasonPhrase string
	var body []byte
	cookieHeaders := map[string]string{}

	// ミドルウェアでクッキーの存在確認はすでにしている
	// ミドルウェアを通ってきたので、クッキーが存在しないことは確定
	// ゆえにクッキーを取得する処理は書かなくて良い
	if request.Method == http.Get {
		body = view.Render("login_form.html")
		statusCode = http.StatusSuccessCode
		reasonPhrase = http.StatusReasonOk
	} else if request.Method == http.Post {
		// ユーザーログインをここでする
		// emailでユーザーを特定
		// その後Loginメソッド内で比較する。OKならクッキーを返す。そして mypageにリダイレクトする
		// ログインに失敗したなら、
		decodedBody, _ := url.QueryUnescape(string(request.Body))
		values, _ := url.ParseQuery(decodedBody)
		email := values.Get("email")
		password := values.Get("password")

		user, err := model.FindUserByEmail(email)
		if err != nil {
			log.Printf("fail to find user by email: %v\n", err)
			body = view.Render("login_form.html")
			statusCode = http.StatusInternalServerErrorCode
			reasonPhrase = http.StatusReasonInternalServerError
		}

		if err := user.Login(password); err != nil {
			log.Printf("fail to find user by email: %v\n", err)
			body = view.Render("login_form.html")
			statusCode = http.StatusInternalServerErrorCode
			reasonPhrase = http.StatusReasonInternalServerError
		} else {
			// ログインに成功した
			statusCode = http.StatusRedirectCode
			reasonPhrase = http.StatusReasonRedirect
			session := auth.NewSession()
			sessionId, _ := session.Save(user.Id)
			// ここでレスポンスにヘッダーセットできたなら最高
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

	// 一回でもクライアントにクッキーを送れば、毎回リクエストごとにクッキーを送ってくる
	if request.Method == http.Post && statusCode == http.StatusRedirectCode {
		for key, value := range cookieHeaders {
			response.SetCookieHeader(key, value)
		}
		response.SetHeader("Location", "/mypage")
	}

	return response
}

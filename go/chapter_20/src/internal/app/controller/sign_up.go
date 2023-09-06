package controller

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path"

	"github.com/yukiHaga/web_server/src/internal/app/config/settings"
	"github.com/yukiHaga/web_server/src/internal/app/model"
	"github.com/yukiHaga/web_server/src/pkg/henagin/http"
)

type SignUp struct{}

func NewSignUp() *SignUp {
	return &SignUp{}
}

func (c *SignUp) Action(request *http.Request) *http.Response {
	STATIC_ROOT, _ := settings.GetStaticRoot()
	var statusCode string
	var reasonPhrase string
	var body []byte
	cookieHeaders := map[string]string{}

	if request.Method == http.Get {
		if cookie, isThere := request.GetCookieByName("user_id"); isThere {
			id, _ := url.QueryUnescape(cookie.Value)
			_, err := model.FindUserById(id)
			if err != nil {
				log.Printf("fail to find error: %v", err)
				body, _ = os.ReadFile(path.Join(STATIC_ROOT, "sign_up_form.html"))
				statusCode = http.StatusInternalServerErrorCode
				reasonPhrase = http.StatusReasonInternalServerError
			} else {
				statusCode = http.StatusRedirectCode
				reasonPhrase = http.StatusReasonRedirect
			}
		} else {
			body, _ = os.ReadFile(path.Join(STATIC_ROOT, "sign_up_form.html"))
			statusCode = http.StatusSuccessCode
			reasonPhrase = http.StatusReasonOk
		}
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
			body, _ = os.ReadFile(path.Join(STATIC_ROOT, "sign_up_form.html"))
			statusCode = http.StatusInternalServerErrorCode
			reasonPhrase = http.StatusReasonInternalServerError
		} else {
			statusCode = http.StatusRedirectCode
			reasonPhrase = http.StatusReasonRedirect
			cookieHeaders["user_id"] = fmt.Sprintf("%v", user.Id)
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

	if statusCode == http.StatusRedirectCode {
		if request.Method == http.Post {
			for key, value := range cookieHeaders {
				response.SetCookieHeader(key, value)
			}
		}
		response.SetHeader("Location", "/mypage")
	}

	return response
}

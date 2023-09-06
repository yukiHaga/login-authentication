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

type Login struct{}

func NewLogin() *Login {
	return &Login{}
}

func (c *Login) Action(request *http.Request) *http.Response {
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
				body, _ = os.ReadFile(path.Join(STATIC_ROOT, "login_form.html"))
				statusCode = http.StatusInternalServerErrorCode
				reasonPhrase = http.StatusReasonInternalServerError
			} else {
				// このif文の中に入ったってことは、ユーザーがdbに存在していて、かつクライアントにクッキーを送っている
				statusCode = http.StatusRedirectCode
				reasonPhrase = http.StatusReasonRedirect
			}
		} else {
			body, _ = os.ReadFile(path.Join(STATIC_ROOT, "login_form.html"))
			statusCode = http.StatusSuccessCode
			reasonPhrase = http.StatusReasonOk
		}
	} else if request.Method == http.Post {
		// ユーザーログインをここでする
		// emailでユーザーを特定
		// その後Loginメソッド内で比較する。OKならクッキーを返す。そして mypageにリダイレクトする
		// ログインに失敗したなら、
		decodedBody, _ := url.QueryUnescape(string(request.Body))
		values, _ := url.ParseQuery(decodedBody)
		fmt.Println("values", values)
		email := values.Get("email")
		password := values.Get("password")

		user, err := model.FindUserByEmail(email)
		if err != nil {
			log.Printf("fail to find user by email: %v\n", err)
			body, _ = os.ReadFile(path.Join(STATIC_ROOT, "login_form.html"))
			statusCode = http.StatusInternalServerErrorCode
			reasonPhrase = http.StatusReasonInternalServerError
		}

		if err := user.Login(password); err != nil {
			log.Printf("fail to find user by email: %v\n", err)
			body, _ = os.ReadFile(path.Join(STATIC_ROOT, "login_form.html"))
			statusCode = http.StatusInternalServerErrorCode
			reasonPhrase = http.StatusReasonInternalServerError
		} else {
			statusCode = http.StatusRedirectCode
			reasonPhrase = http.StatusReasonRedirect
			// ここでレスポンスにヘッダーセットできたなら最高
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

	// 一回でもクライアントにクッキーを送れば、毎回リクエストごとにクッキーを送ってくる
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

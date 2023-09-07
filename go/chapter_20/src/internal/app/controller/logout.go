package controller

import (
	"context"
	"fmt"
	"time"

	"github.com/yukiHaga/web_server/src/internal/app/auth"
	"github.com/yukiHaga/web_server/src/pkg/henagin/http"
)

type Logout struct{}

func NewLogout() *Logout {
	return &Logout{}
}

var ctx = context.Background()

func (c *Logout) Action(request *http.Request) *http.Response {
	var statusCode string
	var reasonPhrase string
	body := []byte{}
	cookieHeaders := map[string]string{}

	if request.Method == http.DELETE {
		cookie, _ := request.GetCookieByName("session_id")
		session := auth.NewSession()
		sessionId := cookie.Value
		// 循環インポートマジでうざいな
		// ログアウトに成功した
		session.Store.Del(ctx, sessionId)
		statusCode = http.StatusSuccessCode
		reasonPhrase = http.StatusReasonOk
		// 同じ名前のクッキーを送信し、期限を過去に設定した
		// クライアントは期日が過去日のクッキーを受け取ると、同じ名前だと勝手に消してくれる
		// timeの部分は1日前の日付を生成している
		gmtTime := time.Now().AddDate(0, 0, -1).UTC()
		formattedTime := gmtTime.Format("Mon, 02 Jan 2006 15:04:05 GMT")
		cookieHeaders["session_id"] = fmt.Sprintf("%s; Expires=%v", "deactivated-session", formattedTime)
	} else {
		statusCode = http.StatusInternalServerErrorCode
		reasonPhrase = http.StatusReasonInternalServerError
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
	if request.Method == http.DELETE && statusCode == http.StatusSuccessCode {
		for key, value := range cookieHeaders {
			response.SetCookieHeader(key, value)
		}
		response.SetHeader("Location", "/login")
	}

	return response
}

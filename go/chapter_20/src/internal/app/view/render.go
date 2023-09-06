package view

import (
	"fmt"
	"os"
	"path"
	"regexp"

	"github.com/yukiHaga/web_server/src/internal/app/config/settings"
)

// テンプレートエンジンのコア機能
// 呼び出し時にキーワード引数で渡せたら最高だけど、なんかめんどくさそうだからやめた。ビューで使われる変数の順番を覚えていないと使えないのがネック
func Render(templatePath string, context ...any) []byte {
	root, _ := settings.GetTemplateRoot()
	bytes, _ := os.ReadFile(path.Join(root, templatePath))
	fileContent := string(bytes)

	re := regexp.MustCompile(`{(.+?)}`)
	// fileContentの中の全ての変数展開を%vにリプレイスする
	replacedFileContent := re.ReplaceAllString(fileContent, "%v")
	view := fmt.Sprintf(replacedFileContent, context...)
	body := []byte(view)

	return body
}

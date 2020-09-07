// @Author:DaiPengyuan
// @Date:2018/10/9
// @Desc:

package stringx

import (
	"fmt"
	"testing"
)

func TestIsChineseChar(t *testing.T) {
	fmt.Println(IsChineseChar("ddd"))
	fmt.Println(IsChineseChar("，。‘”"))
	fmt.Println(IsChineseChar("你好"))
}

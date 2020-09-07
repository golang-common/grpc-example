// @Author: Perry
// @Date  : 2020/5/11
// @Desc  : json快捷操作

package common

import "encoding/json"

/*将内容转化为json字符串,转化失败则返回错误字符串*/
func GetIndentJson(obj interface{}) string {
	ret, err := json.MarshalIndent(obj, "", "\t")
	if err != nil {
		return err.Error()
	}
	return string(ret)
}



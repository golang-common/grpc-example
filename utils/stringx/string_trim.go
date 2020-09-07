// @Author:DaiPengyuan
// @Date:2019/4/14
// @Desc:

package stringx

import "strings"

func TrimBlockString(strBlock string, cutset string) string {
	var tmp []string
	for _, line := range strings.Split(strBlock, "\n") {
		tmp = append(tmp, strings.Trim(line, cutset))
	}
	return strings.Join(tmp, "\n")
}

func CutBlockStringLeft(strBlock string, cutset string) string {
	var tmp []string
	for _, line := range strings.Split(strBlock, "\n") {
		if strings.HasPrefix(line, cutset) {
			line = line[len(cutset):]
		}
		tmp = append(tmp, line)
	}
	return strings.Join(tmp, "\n")
}
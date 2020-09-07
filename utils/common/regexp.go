// @Author: Perry
// @Date  : 2020/5/11
// @Desc  : 正则表达式快捷操作

package common

import "regexp"

/*
功能与RegStringByGroupMap类似,类似匹配多次s中命中的reg,
r中的数组长度为reg匹配的次数
可以设定n的值来限制最大匹配次数,默认匹配全部
*/
func RegStringAllByGroup(s string, reg *regexp.Regexp, n ...int) (r []map[string]string) {
	var (
		matchTime  = -1
		groups     []string
		groupMList [][]string
	)
	if len(n) != 0 {
		matchTime = n[0]
	}
	groups = reg.SubexpNames()
	groupMList = reg.FindAllStringSubmatch(s, matchTime)
	for _, groupR := range groupMList {
		match := make(map[string]string)
		for gnK, gName := range groups {
			if gName !=""{
				match[gName] = groupR[gnK]
			}
		}
		r = append(r, match)
	}
	return
}

/*
用已经初始化规则的正则表达式reg检查字符串s中相应匹配的字符
如果正则命中的字符串的两处位置,则保留首先匹配到的值
r中的key对应正则中的<?P>组名,如果组名为空则忽略
r中对应的value对应正则中组对应的名称
*/
func RegStringByGroupMap(s string, reg *regexp.Regexp) (r map[string]string) {
	var (
		groups       []string
		groupVarList []string
	)
	r = make(map[string]string)
	groups = reg.SubexpNames()
	groupVarList = reg.FindStringSubmatch(s)
	for k, v := range groupVarList {
		gName := groups[k]
		if gName!=""{
			r[groups[k]] = v
		}
	}
	return
}

/*
输入待匹配字符s,组名groupName,以及正则reg,返回组名匹配的字符串
*/
func RegStringByGroup(s, groupName string, reg *regexp.Regexp) (r string) {
	for k, v := range reg.SubexpNames() {
		if groupName == v {
			l := reg.FindStringSubmatch(s)
			if k < len(l) {
				r = l[k]
			}
		}
	}
	return
}

/*输入待匹配字符串s,组名groupName,以及正则reg,返回首个命中的匹配字符的首尾字符位置,如果不匹配返回空*/
func RegStringIndexByGroup(s, groupName string, reg *regexp.Regexp) (r []int) {
	for k, v := range reg.SubexpNames() {
		if groupName == v {
			l := reg.FindStringSubmatchIndex(s)
			if k < len(l)/2 {
				r = append(r, l[k*2])
				r = append(r, l[k*2+1])
			}
		}
	}
	return
}

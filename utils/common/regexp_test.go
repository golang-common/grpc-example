// @Author: Perry
// @Date  : 2020/5/11
// @Desc  : 

package common

import (
	"regexp"
	"testing"
)

var (
	testStr = `record1:dpy's ip address is 192.168.1.1 and next hop is 192.168.1.254
record2:bob's ip address is 172.16.1.1 and next hop is 172.16.1.254`
	testReg = regexp.MustCompile(`(?P<ipaddr>[\d]{1,3}\.[\d]{1,3}\.[\d]{1,3}\.[\d]{1,3}).+next hop is (?P<nexthop>[\d]{1,3}\.[\d]{1,3}\.[\d]{1,3}\.[\d]{1,3})`)
)

func TestRegStringAllByGroup(t *testing.T) {
	r := RegStringAllByGroup(testStr, testReg)
	t.Log(r)
}
// Output:[map[ipaddr:192.168.1.1 nexthop:192.168.1.254] map[ipaddr:172.16.1.1 nexthop:172.16.1.254]]
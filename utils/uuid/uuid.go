// @Author: Perry
// @Date  : 2020/4/28
// @Desc  : 

package uuid

import uuid "github.com/satori/go.uuid"

func NewUUIDV1() string {
	return uuid.NewV1().String()
}

func NewUUIDV4() string {
	return uuid.NewV4().String()
}

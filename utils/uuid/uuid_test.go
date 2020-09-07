// @Author: Perry
// @Date  : 2020/4/28
// @Desc  : 

package uuid

import (
	"fmt"
	"testing"
)

func TestNewUUIDV1(t *testing.T) {
	fmt.Println(NewUUIDV1())
	fmt.Println(NewUUIDV1())
	fmt.Println(NewUUIDV1())
	fmt.Println(NewUUIDV1())
	fmt.Println(len(NewUUIDV1()))
}
/*
* @Author : DAIPENGYUAN
* @File : encrypt_test
* @Time : 2020/8/4 17:55
* @Description :
 */

package encrypt

import "testing"

func TestDesEncrypt(t *testing.T) {
	encryptedStr, err := DesEncrypt("teststring")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(encryptedStr)
}

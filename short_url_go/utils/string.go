package utils

import (
	"crypto/md5"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"testing"
)

// 把输入的字符计算为 MD5
func MD5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has) //将[]byte转成16进制
	return strings.ToUpper(md5str)
}

// 测试uuid，通过os自动生成
func TestUUID1(t *testing.T) {
	out, err := exec.Command("uuidgen").Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s \n", out)
}

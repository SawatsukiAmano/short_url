package utils

import (
	"fmt"

	"github.com/beego/beego/v2/core/config"
)

// 全局配置
var INIconf config.Configer

func init() {
	var err error
	INIconf, err = config.NewConfig("ini", "../data/go/conf/secret.conf")
	if err != nil {
		fmt.Println(err)
		panic(err) //https://zhuanlan.zhihu.com/p/373653492
	}
}

package controllers

import (
	"short_url_go/models"
	"time"

	"github.com/beego/beego/v2/client/cache"
)

var DoaminUser map[string]uint

var Short2Long map[uint]map[string]string

var AC cache.Cache

var LatestDoaminUserTime time.Time

// 刷新 域名——用户 内存缓存
func RefreshDoaminUser(isforce bool) {
	if LatestDoaminUserTime.AddDate(0, 0, 1).Before(time.Now()) || isforce {
		DoaminUser = models.QueryUsersDomainID()
	}
}

// https://www.cnblogs.com/hei-ma/articles/13847724.html
func init() {
	// var err error
	// AC, err = cache.NewCache("memory", `{"interval":"86400"}`)
	// if err != nil {
	// 	fmt.Println("NewCache failed, err:", err)
	// }
	// RefreshDoaminUser(true)
}

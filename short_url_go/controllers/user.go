﻿package controllers

import (
	"encoding/json"
	"fmt"
	"short_url_go/models"
	"short_url_go/utils"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

// Operations about Users
type UserController struct {
	BaseController
}

//https://www.cnblogs.com/arestrack/p/7799425.html#%E7%9B%B4%E6%8E%A5%E8%BE%93%E5%87%BA%E5%AD%97%E7%AC%A6%E4%B8%B2
//模板
//https://cloud.tencent.com/developer/article/1557075

// @Title Register
// @Summary 注册
// @Description  user into the system
// @Param	body		body 	models.User	true	"body for user"
// @Success 200 {string}	register success
// @Failure 403	{string}	Prohibiting the registration
// @router /register [post]
func (u *UserController) Register() {
	// Infos(u)
	var user models.User
	json.Unmarshal(u.Ctx.Input.RequestBody, &user)
	user.Role = 0
	user.DefaultURLLength = 6

	if _, err := models.CreateUser(user); err != nil {
		u.Ctx.WriteString("注册成功")
	} else {
		u.Ctx.WriteString(err.Error())
	}
}

// @Title Login
// @Description user into the system
// @Summary 登录
// @Param	body		body 	controllers.Login	true		"The username for login"
// @Success 200 {models.User}	Login success
// @Failure 401 string	The user does not exist.
// @router /login [post]
func (u *UserController) Login() {
	requestBody := u.JsonData()
	username := requestBody["name"].(string)
	password := requestBody["pwd"].(string)
	user, err := models.Login(username, password)
	fmt.Println("进入login")
	if user.ID > 0 && len(user.Name) > 0 && err == nil {
		u.Ctx.WriteString(generateRefreshJWT(user.ID))
		return
	} else {
		u.Ctx.ResponseWriter.WriteHeader(401)
		u.Ctx.WriteString(err.Error())
	}
}

type Login struct {
	Name     string `json:"name"`
	Password string `json:"pwd"`
}

// @Title account tocken
// @Description user into the system
// @Summary 刷新 account tocken
// @Param	jwt		body 	string	true	"The refresh jwt tocken"
// @Success 200	{string}	Refresh success
// @Failure 401	{string}	refresh token 失效
// @router /tocken/account [post]
func (u *UserController) RefreshTocken() {
	tokenString := string(u.Ctx.Input.RequestBody)
	tokenString = strings.Trim(tokenString, "\"")
	if len(tokenString) > 0 {
		token, err := jwt.ParseWithClaims(tokenString, &RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
			key, _ := utils.INIconf.String("JWT::RefreshTokenKey")
			return []byte(key), nil
		})
		if err != nil {
			u.Ctx.ResponseWriter.WriteHeader(401)
			u.Ctx.WriteString("refresh token 失效，请重新登录")
			return
		}
		if claims, ok := token.Claims.(*RefreshClaims); ok && token.Valid {
			fmt.Println(claims.ID)
			user, err := models.QueryUserByID(claims.ID)
			if err != nil {
				u.Ctx.ResponseWriter.WriteHeader(400)
				u.Ctx.WriteString(err.Error())
				return
			}
			u.Ctx.WriteString(generateAccountJWT(user))
			return
		} else {
			u.Ctx.ResponseWriter.WriteHeader(401)
			u.Ctx.WriteString("refresh token 失效，请重新登录")
		}
	}
}

// @Title user
// @Summary 新增一个用户
// @Description  user into the system
// @Param	body	body 	models.User	true	"body for user"
// @Param Authorization header string false "Bearer 用户令牌"
// @Success 200	{string}	Create success
// @Failure 403	{string}	Insufficient user permissions
// @router / [post]
func (u *UserController) CreateUser() {
	var user models.User
	json.Unmarshal(u.Ctx.Input.RequestBody, &user)
	accInfo := u.analysisAccountClaims()
	if accInfo.Role == 1 || accInfo.ID == user.ID {
		if _, err := models.CreateUser(user); err == nil {
			u.Ctx.WriteString("创建成功")
		} else {
			u.Ctx.ResponseWriter.WriteHeader(400)
			u.Ctx.WriteString(err.Error())
			return
		}
	} else {
		u.Ctx.ResponseWriter.WriteHeader(403)
		u.Ctx.WriteString("无权创建用户")
	}
}

// @Title Delete
// @Summary 删除一个用户
// @Description delete the user
// @Param	uid		path 	unit	true	"The uid you want to delete"
// @Success 200	{string}	delete success!
// @Failure 403	{string}	Insufficient user permissions
// @router /:uid [delete]
func (u *UserController) DeleteUser() {
	uid, err := u.GetUint64(":uid")
	if err != nil {
		u.Ctx.ResponseWriter.WriteHeader(400)
		u.Ctx.WriteString("参数错误")
		return
	}
	accInfo := u.analysisAccountClaims()
	if accInfo.Role == 1 || accInfo.ID == uint(uid) {
		if models.DeleteUser(uint(uid)) == nil {
			u.Ctx.WriteString("删除成功!")
			return
		} else {
			u.Ctx.WriteString(err.Error())
			u.Ctx.ResponseWriter.WriteHeader(400)
			return
		}
	} else {
		u.Ctx.ResponseWriter.WriteHeader(403)
		u.Ctx.WriteString("无权删除用户")
		return
	}
}

// @Title	UpdateUser
// @Summary 修改一个用户
// @Description update the user
// @Param	user	body 	models.User true	"body for user"
// @Success 200 {string} "update success!"
// @Failure 403 {string} "Insufficient user permissions"
// @router / [put]
func (u *UserController) UpdateUser() {
	var user models.User
	json.Unmarshal(u.Ctx.Input.RequestBody, &user)
	accInfo := u.analysisAccountClaims()
	if accInfo.Role == 1 || accInfo.ID == user.ID {
		existUser, err := models.QueryUserByID(user.ID)
		if err != nil {
			u.Ctx.WriteString(err.Error())
			u.Ctx.ResponseWriter.WriteHeader(400)
			return
		}
		existUser.Name = user.Name
		existUser.Nickname = user.Nickname
		existUser.DefaultURLLength = user.DefaultURLLength
		existUser.Phone = user.Phone
		existUser.Group = user.Group
		existUser.I18n = user.I18n
		existUser.Remarks = user.Remarks
		existUser.AutoInsertSpace = user.AutoInsertSpace
		existUser.Domain = user.Domain
		existUser.Author = user.Author
		existUser.Role = user.Role
		if models.UpdateUser(existUser) == nil {
			u.Ctx.WriteString("修改成功")
			return
		} else {
			u.Ctx.WriteString(err.Error())
			u.Ctx.ResponseWriter.WriteHeader(400)
			return
		}
	} else {
		u.Ctx.ResponseWriter.WriteHeader(403)
		u.Ctx.WriteString("无权修改用户")
		return
	}
}

// @Title Update password
// @Summary 修改一个用户的密码
// @Description update the user's password
// @Param	uid	path 	int true "用户id"
// @Param	body	body 	string true "body for password"
// @Success 200	{string}	update password success!
// @Failure 403	{string}	Insufficient user permissions
// @router /:uid/pwd [patch]
func (u *UserController) UpdateUserPassword() {
	id, err := u.GetUint64(":uid")
	if err != nil {
		u.Ctx.ResponseWriter.WriteHeader(400)
		u.Ctx.WriteString("参数类型错误")
		return
	}
	var pwd string
	accInfo := u.analysisAccountClaims()
	if uint(id) == accInfo.ID || accInfo.Role == 1 {
		json.Unmarshal(u.Ctx.Input.RequestBody, &pwd)
		err = models.UpdateUserPassword(uint(id), pwd)
		if err == nil {
			u.Ctx.WriteString("修改成功")
			return
		} else {
			u.Ctx.WriteString(err.Error())
			u.Ctx.ResponseWriter.WriteHeader(400)
			return
		}
	} else {
		u.Ctx.WriteString("无权修改他人密码")
		u.Ctx.ResponseWriter.WriteHeader(403)
		return
	}
}

// @Title	GetUsersByPage
// @Summary	user分页查询
// @Date	2022-11-18
// @Auth	sfhj
// @Param	offset	query	int	true	偏移量
// @Param	limit	query	int	true	指定返回记录的数量
// @Param	sort	query	string	true	排序
// @Param	name	query	string	false	账号
// @Param	nickname	query	string	false	昵称
// @Param	group	query	string	false	分组
// @Param	role	query	string	false	权限
// @Param	phone	query	string	false	手机号
// @Param	domain	query	string flase	域名
// @Param	crt	query	string	false	创建时间
// @Param	upt	query	string	false	修改时间
// @Param	det	query	string	false	删除时间
// @Success	200
// @Failure 403	{string}	Insufficient user permissions
// @router / [get]
func (u *UserController) GetUsersByPage() {
	accInfo := u.analysisAccountClaims()
	if accInfo.Role == 0 {
		u.Ctx.ResponseWriter.WriteHeader(403)
		u.Ctx.WriteString("无权查询其他用户")
		return
	}
	var err error
	var page models.Page
	page.Offset, err = u.GetInt("offset")
	if err != nil {
		u.Ctx.ResponseWriter.WriteHeader(400)
		u.Ctx.WriteString("参数错误")
		return
	}
	page.Lmit, err = u.GetInt("limit")
	if err != nil {
		u.Ctx.ResponseWriter.WriteHeader(400)
		u.Ctx.WriteString("参数错误")
		return
	}
	page.Sort = analysisOrderBy(u.GetString("sort"))
	result, count, err := models.QueryUsersPage(page, u.GetString("name"), u.GetString("nickname"), u.GetString("role"), u.GetString("group"), u.GetString("phone"), u.GetString("domain"))
	if err != nil {
		u.Ctx.ResponseWriter.WriteHeader(400)
		u.Ctx.WriteString("参数错误")
		return
	}
	u.Data["json"] = map[string]interface{}{
		"count": count,
		"data":  result,
	}
	u.ServeJSON()
}

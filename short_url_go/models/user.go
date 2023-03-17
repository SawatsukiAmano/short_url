package models

import (
	"errors"
	"time"

	linq "github.com/ahmetb/go-linq/v3"

	uuid "github.com/satori/go.uuid"
)

type User struct {
	ID               uint      `json:"id" gorm:"primaryKey;<-:create"`       //id
	CreatedAt        time.Time `json:"crt" gorm:"<-:create"`                 //创建时间
	UpdatedAt        time.Time `json:"upt" gorm:"<-"`                        //最后更新时间
	Name             string    `json:"name" gorm:"not null"`                 //用户名，登录名称
	Nickname         string    `json:"nickname" gorm:"not null;uniqueIndex"` //昵称
	Password         string    `json:"pwd" gorm:"not null"`                  //密码
	Role             int8      `json:"role" gorm:"not null"`                 //角色
	DefaultURLLength uint8     `json:"urlLength" gorm:"not null"`            //配置项：url默认长度
	Author           []byte    `json:"author"`                               //头像地址
	Phone            string    `json:"phone"`                                //手机号
	Group            string    `json:"group"`                                //分组
	Remarks          string    `json:"remarks"`                              //备注
	I18n             string    `json:"i18n"`                                 //国际化
	AutoInsertSpace  bool      `json:"autoInsertSpace"`                      //盘古之白
	Domain           string    `json:"domain" gorm:"uniqueIndex"`            //域名
}

// 登录
func Login(username, password string) (User, error) {
	var user User
	password = uuid.NewV5(U5Seed, password).String()
	result := DB.Model(&User{}).Where(&User{Name: username, Password: password}).First(&user)
	if result.RowsAffected == 0 {
		return user, errors.New("登录失败，用户名或密码错误")
	}
	return user, nil
}

// 创建用户
func CreateUser(user User) (uint, error) {
	user.Password = uuid.NewV5(U5Seed, user.Password).String()
	result := DB.Create(&user)
	if result.RowsAffected == 0 {
		return 0, errors.New("创建失败，用户名或域名重复")
	}
	return user.ID, nil
}

// 删除用户，同时删除短链接
func DeleteUser(id uint) error {
	DB.Unscoped().Delete(&Short{}, "fk_user = ?", id)
	result := DB.Delete(&User{}, id)
	if result.RowsAffected == 0 {
		return errors.New("删除失败，数据库错误")
	}
	return nil
}

// 删除多个用户，同时删除短链接
func DeleteUsers(ids []uint) error {
	DB.Delete(&Short{}, "fk_user IN ?", ids)
	result := DB.Delete(&User{}, "id IN ? ", ids)
	if result.RowsAffected == 0 {
		return errors.New("删除失败，数据库错误")
	}
	return nil
}

// 更新用户，不更新密码
func UpdateUser(user User) error {
	var existUser User
	if DB.Unscoped().Where("(name = ? OR domain = ?) AND id != ?", user.Name, user.Domain, user.ID).First(&existUser).RowsAffected > 0 {
		return errors.New("未查找到改用户")
	}
	result := DB.Model(&user).Updates(User{Name: user.Name, Nickname: user.Nickname, Role: user.Role, Author: user.Author, Phone: user.Phone, Group: user.Group, I18n: user.I18n, AutoInsertSpace: user.AutoInsertSpace, Remarks: user.Remarks, DefaultURLLength: user.DefaultURLLength, Domain: user.Domain})
	if result.RowsAffected == 0 {
		return errors.New("修改失败,账号名或域名重复")
	}
	return nil
}

// 更新密码
func UpdateUserPassword(id uint, pwd string) error {
	if len(pwd) < 8 {
		return errors.New("密码长度过短")
	} else if len(pwd) > 32 {
		return errors.New("密码长度过长")
	}
	password := uuid.NewV5(U5Seed, pwd).String()
	result := DB.Model(&User{}).Where("id = ?", id).Update("password", password)
	if result.RowsAffected == 0 {
		return errors.New("更新失败，数据库错误")
	}
	return nil
}

// 分页查询
func QueryUsersPage(page Page, name string, nickname string, role string, group string, phone string, domain string) (result []User, count int64, err error) {
	express := DB.Model(&User{})
	if analysisRestfulRHS(express, "name", name) &&
		analysisRestfulRHS(express, "nickname", nickname) &&
		analysisRestfulRHS(express, "role", role) &&
		analysisRestfulRHS(express, "phone", phone) &&
		analysisRestfulRHS(express, "domain", domain) &&
		analysisRestfulRHS(express, "group", group) {
		express.Count(&count)
		express = express.Order(page.Sort).Limit(page.Lmit).Offset((page.Offset - 1) * page.Lmit)
		express.Select("id", "created_at", "updated_at", "name", "nickname", "role", "default_url_length", "group", "i18n", "auto_insert_space", "remarks", "domain", "phone").Find(&result)
	} else {
		err = errors.New("查询参数错误")
	}
	return
}

// 根据id查询用户
func QueryUserByID(id uint) (User, error) {
	var user User
	result := DB.Unscoped().First(&user, id)
	if result.RowsAffected == 0 {
		return user, errors.New("未查询到用户")
	}
	return user, nil
}

// 查询所有用户的 域名-id 键值对
func QueryUsersDomainID() map[string]uint {
	var users []User
	DB.Find(&users)
	result := make(map[string]uint, len(users))
	linq.From(users).SelectT(func(e User) map[string]uint {
		return map[string]uint{e.Domain: e.ID}
	}).ToMap(&result)
	return result
}

// func QueryUserByDomain(domain string) User {
// 	var user User
// 	DB.Unscoped().Where("domain = ?", domain).First(&user)
// 	return user
// }

package msp

import (
	"errors"
	"strings"
)

// IsEmail 判断是否是邮箱格式 测试成功
func IsEmail(email string) bool {
	return strings.Contains(email, "@")
}

// IsNameAndPwd 用户名或密码5-13位之间
func IsNameAndPwd(Name, Pwd string) (bool, error) {

	if len(Name) < 5 || len(Name) > 13 {
		return false, errors.New("用户名或密码长度5-13位")
	}
	if len(Pwd) < 5 || len(Pwd) > 13 {
		return false, errors.New("用户名或密码长度5-13位")
	}
	return true, nil
}

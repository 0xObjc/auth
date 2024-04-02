package sdk

import (
	"fmt"
	"os"
	"time"
)

var MS MspKey

type MspkeyUI struct {
	IP       string   //服务器IP
	ExeID    string   //软件ID
	DevID    string   //设备ID
	Version  string   //当前版本
	exeInfo  ExeInfo  //软件信息
	userInfo UserInfo //用户信息
	Fc       func()   //退出回调
}

func (c *MspkeyUI) UIStart() {
	err := MS.Init(c.IP, c.ExeID, c.DevID)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	c.exeInfo, err = MS.GetExeInfo()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	c.home()
}

// Home 首页
func (c *MspkeyUI) home() {
	for {
		//清屏
		Clear()
		fmt.Println(c.exeInfo.Title + "    当前版本:" + c.Version)
		fmt.Println()
		fmt.Println("----------------公告-------------------")
		fmt.Println(c.exeInfo.Notice)
		var st int
		str := `----------------菜单-------------------
1.账号登录		2.卡密登录
3.注册账号		4.账号充值
5.修改密码		6.设备换绑
7.购买卡密		8.扫码支付
0.退出`
		fmt.Println(str)
		fmt.Print("请选择(0-8):")
		_, _ = fmt.Scan(&st)
		switch st {
		case 0:
			os.Exit(0)
		case 1:
			c.login()
		case 2:
			c.carLogin()
		}

		if MS.IsLogin {
			return
		}
	}
}

// login 账号登录
func (c *MspkeyUI) login() {
	for {
		Clear()
		var Name, pwd string
		fmt.Println("账号登录")
		fmt.Print("请输入账号:")
		_, _ = fmt.Scan(&Name)
		fmt.Print("请输入密码:")
		_, _ = fmt.Scan(&pwd)
		err := MS.Login(Name, pwd)
		if err != nil {
			fmt.Println(err)
			time.Sleep(time.Second * 1)
			continue
		} else {
			fmt.Println("登录成功!")
			c.userInfo, err = MS.GetUserInfo()
			if err != nil {
				fmt.Println(err)
			}
			break
		}

	}

}

// carLogin
func (c *MspkeyUI) carLogin() {
	for {
		Clear()
		var Serial string
		fmt.Println("卡密模式")
		fmt.Print("请输入卡密:")
		_, _ = fmt.Scan(&Serial)
		err := MS.CarLogin(Serial)
		if err != nil {
			fmt.Println(err)
			time.Sleep(time.Second * 1)
			continue
		} else {
			fmt.Println("登录成功!")
			c.userInfo, err = MS.GetUserInfo()
			if err != nil {
				fmt.Println(err)
			}
			break
		}
	}

}

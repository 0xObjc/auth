package sdk

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type ResJson struct {
	Tag  string `json:"Tag" bson:"Tag"`
	Code int    `json:"Code" bson:"Code"`
	Msg  string `json:"Msg" bson:"Msg"`
	Data any    `json:"Data" bson:"Data"`
	Time string `json:"Time" bson:"Time"`
}

type SendJson struct {
	Type string `json:"Type" bson:"Type"`
	Data any    `json:"Data" bson:"Data"`
}

// ExeInfo 软件
type ExeInfo struct {
	ID              primitive.ObjectID `json:"ID" bson:"_id"`                           //软件ID
	AdminID         primitive.ObjectID `json:"AdminID"  bson:"AdminID"`                 //绑定管理员ID
	Title           string             `json:"Title"  bson:"Title"`                     //软件标题
	Versions        string             `json:"Versions"  bson:"Versions"`               //版本
	State           bool               `json:"State"  bson:"State"`                     //状态 正常/禁用
	Notice          string             `json:"Notice"  bson:"Notice"`                   //公告
	Address         string             `json:"Address"  bson:"Address"`                 //更新地址
	Md5             string             `json:"Md5"  bson:"Md5"`                         //软件MD5
	Data            string             `json:"Data"  bson:"Data"`                       //软件核心数据
	Key             string             `json:"Key"  bson:"Key"`                         //密钥
	IsDK            bool               `json:"IsDK"  bson:"IsDK"`                       //是否多开
	IsReg           bool               `json:"IsReg"  bson:"IsReg"`                     //是否允许注册
	IsDbg           bool               `json:"IsDbg"  bson:"IsDbg"`                     //是否开启检测
	IsBindIP        bool               `json:"IsBindIP"  bson:"IsBindIP"`               //是否绑定IP
	IsDeviceID      bool               `json:"IsDeviceID"  bson:"IsDeviceID"`           //绑定设备ID
	GiveTime        int64              `json:"GiveTime"  bson:"GiveTime"`               //软件注册赠送时间 分钟
	BindCount       int64              `json:"BindCount"  bson:"BindCount"`             //设置换绑次数 次/月
	SubTime         int64              `json:"SubTime" bson:"SubTime"`                  //换绑扣时间 单位/小时
	BindDeviceIDNum int64              `json:"BindDeviceIDNum"  bson:"BindDeviceIDNum"` //同一软件同一设备用户注册数量限制
	LoginMod        int64              `json:"LoginMod"  bson:"LoginMod"`               //登录模式0=单卡+用户 1=用户登录 2=单卡登录
	ExePrice        ExePrice           `json:"ExePrice" bson:"ExePrice"`                //软件价格
}

// ExePrice 软件价格
type ExePrice struct {
	HourCar      float64 `json:"HourCar" bson:"HourCar"`           //小时卡价格
	DayCar       float64 `json:"DayCar" bson:"DayCar"`             //天卡价格
	WeekCar      float64 `json:"WeekCar" bson:"WeekCar"`           //周卡价格
	MonthCar     float64 `json:"MonthCar" bson:"MonthCar"`         //月卡价格
	SeasonCar    float64 `json:"SeasonCar" bson:"SeasonCar"`       //季卡价格
	HalfYearCar  float64 `json:"HalfYearCar" bson:"HalfYearCar"`   //半年卡价格
	YearCar      float64 `json:"YearCar" bson:"YearCar"`           //年卡价格
	PermanentCar float64 `json:"PermanentCar" bson:"PermanentCar"` //永久卡价格
}

// Car   卡密结构
type Car struct {
	ID        primitive.ObjectID ` bson:"_id"`      //卡ID
	AdminID   primitive.ObjectID ` bson:"AdminID"`  //绑定管理员ID
	ExeID     primitive.ObjectID ` bson:"ExeID"`    //软件ID
	Serial    string             ` bson:"Serial"`   //卡号
	State     bool               ` bson:"State"`    //状态 正常/禁用
	TyCar     int64              ` bson:"TyCar"`    //卡密类型  0=小时卡 1=天卡 2=周卡 3=月卡 4=季卡 5=半年卡 6=年卡 7=永久卡
	Price     float64            ` bson:"Price"`    //售价
	Bak       string             ` bson:"Bak"`      //备注
	Lock      bool               `bson:"Lock"`      //锁定卡密不能被删除和获取
	BillID    string             `bson:"BillID"`    //支付宝交易号
	CreatTime time.Time          `bson:"CreatTime"` //制卡时间
}

// UserInfo 用户
type UserInfo struct {
	ID            primitive.ObjectID `bson:"_id"`           //用户ID
	AdminID       primitive.ObjectID `bson:"AdminID"`       //管理员ID
	AgentID       primitive.ObjectID `bson:"AgentID"`       //代理ID
	DeviceID      string             `bson:"DeviceID"`      //设备ID
	ExeID         primitive.ObjectID `bson:"ExeID"`         //绑定的软件
	Name          string             `bson:"Name"`          //用户名
	Pwd           string             `bson:"Pwd"`           //密码
	Level         int64              `bson:"Level"`         //用户等级 0=小时卡 1=天卡 2=周卡 3=月卡 4=季卡 5=半年卡 6=年卡 7=永久卡
	Serial        string             `bson:"Serial"`        //最后一次充值的卡号
	State         bool               `bson:"State"`         //状态
	Online        bool               `bson:"Online"`        //在线状态
	RegIP         string             `bson:"RegIP"`         //注册ip
	RegTime       time.Time          `bson:"RegTime"`       //注册时间
	EndTime       time.Time          `bson:"EndTime"`       //到期时间
	LoginIP       string             `bson:"LoginIP"`       //登录ip
	LastLoginIP   string             `bson:"LastLoginIP"`   //上一次登录ip
	LoginTime     time.Time          `bson:"LoginTime"`     //登录时间
	LastLoginTime time.Time          `bson:"LastLoginTime"` //上一次登录时间
	Bak           string             `bson:"Bak"`           //备注
	BindCont      int64              `bson:"BindCont"`      //当前用户换绑次数
	BindTime      time.Time          `bson:"BindTime"`      //换绑的时间
	Conf          string             `bson:"Conf"`          //用户配置信息
}

// PayLink 支付连接
type PayLink struct {
	PayLink     string `json:"PayLink"`     //微信连接
	PayCarLink  string `json:"PayCarLink"`  //发卡连接
	PayAliState bool   `json:"PayAliState"` //支付宝当面付是否启用
}

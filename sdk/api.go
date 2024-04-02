package sdk

import (
	"crypto/rc4"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Rc4EncryptString 用rc4进行加密 返回base64 格式数据
func Rc4EncryptString(key, strData string) string {
	cipher, err := rc4.NewCipher([]byte(key))
	if err != nil {
		return ""
	}
	data := make([]byte, len(strData))
	cipher.XORKeyStream(data, []byte(strData))
	encoded := base64.StdEncoding.EncodeToString(data)
	return encoded
}

// Rc4DecodeString 用rc4进行解密
func Rc4DecodeString(key string, Base64Data string) ([]byte, error) {
	decodeText, err := base64.StdEncoding.DecodeString(Base64Data)
	if err != nil {
		return nil, errors.New("base64 解码失败")
	}
	cipher, err := rc4.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}
	output := make([]byte, len(decodeText))
	cipher.XORKeyStream(output, decodeText)
	return output, nil
}

type MspKey struct {
	conn     *websocket.Conn
	response *http.Response
	key      string
	IsLogin  bool   //是否登录
	IsDug    bool   //是否调试信息输出
	Msg      string //获取Msg 消息
}

var ResData []ResJson

func (c *MspKey) auto() {

}

// GetData 服务器消息返回事件
func (c *MspKey) onMessage() {
	for {
		var temp ResJson
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("Error in receive:", err)
			//os.Exit(1)
			return
		}
		res := strings.ReplaceAll(string(msg), "\"", "")
		msg, err = Rc4DecodeString(c.key, res)
		if err != nil {
			continue
		}
		if c.IsDug {
			log.Printf("接受数据:<- %s\n", msg)
		}
		err = json.Unmarshal(msg, &temp)
		if err != nil {
			log.Println("Error:", err)
			continue
		}
		//防止攻击回放 用于时间戳判断
		p, _ := strconv.ParseInt(temp.Time, 10, 64)
		if p <= time.Now().Unix() {
			fmt.Println("校验失败,数据已过期")
			//os.Exit(1)
			_ = c.AddBlack("玄月检测:检测到内存被修改")
		}

		//动态密钥替换操作
		if c.key == "mspkey" && temp.Tag == "DevKey" {
			c.key = fmt.Sprintf("%s", temp.Data)
			continue
		}

		//实时消息
		if temp.Tag == "SendMsg" && temp.Code == 1 {
			fmt.Println("收到一条实时消息:" + temp.Msg)
			continue
		}
		//主动下线
		if temp.Tag == "OffLine" && temp.Code == 1 {
			if temp.Msg != "主动下线" {
				fmt.Println(temp.Msg)
			}
			//os.Exit(1)
			_ = c.AddBlack("玄月检测:检测到内存被修改")
		}

		//支付
		if temp.Tag == "CompletePayment" && temp.Code == 1 {
			fmt.Println(temp.Msg)
			continue
		}

		if temp.Tag != "Null" {
			ResData = append(ResData, temp)
		}

	}

}

// sendData 发送数据并接受数据
func (c *MspKey) sendData(data SendJson) (ResJson, error) {
	marshal, err := json.Marshal(data)
	if err != nil {
		return ResJson{}, err
	}
	if c.IsDug {
		log.Println("发送数据:->" + string(marshal))
	}

	msg := Rc4EncryptString(c.key, string(marshal))
	err = c.conn.WriteMessage(1, []byte(msg))
	if err != nil {
		return ResJson{}, err
	}
	//并等待数据返回
	for {
		for i, datum := range ResData {
			if datum.Tag == data.Type {
				c.Msg = datum.Msg
				ResData = append(ResData[:i], ResData[i+1:]...)
				if datum.Code == 1 {
					return datum, nil
				}
				return ResJson{}, errors.New(c.Msg)
			}
		}
	}

}

// Init 验证初始化
func (c *MspKey) Init(IP string, ExeID, DeviceID string) error {
	var err error
	IP = fmt.Sprintf("ws://%s/api/user/ws?ExeID=%s&DevID=%s", IP, ExeID, DeviceID)
	c.key = "mspkey" //默认密钥
	c.conn, c.response, err = websocket.DefaultDialer.Dial(IP, nil)
	if err != nil {
		return errors.New("服务器连接失败")
	}
	go c.onMessage()
	//启动程序监听程序
	count := 0
	for {
		if c.key != "mspkey" {
			//启动心跳包
			go func() {
				for {
					time.Sleep(time.Second * 60)
					c.ping()
				}
			}()

			return nil
		}
		if count >= 5 {
			break
		}
		count++
		time.Sleep(time.Second)
	}
	_ = c.conn.Close()
	return errors.New("等待超时")
}

// GetCode 获取验证码
func (c *MspKey) GetCode() (string, error) {
	var p SendJson
	p.Type = "GetCode"
	data, err := c.sendData(p)
	if err != nil {
		return "", err
	}
	return data.Data.(string), nil
}

// GetExeInfo 获取软件基本信息
func (c *MspKey) GetExeInfo() (ExeInfo, error) {
	var p SendJson
	p.Type = "GetExeInfo"
	data, err := c.sendData(p)
	if err != nil {
		return ExeInfo{}, err
	}

	var st struct {
		Exe ExeInfo
	}
	marshal, _ := json.Marshal(data.Data)
	_ = json.Unmarshal(marshal, &st)
	return st.Exe, nil

}

// Register 用户注册
func (c *MspKey) Register(Name, Pwd, Code string) error {
	var p SendJson
	p.Type = "Register"
	p.Data = bson.M{"Name": Name, "Pwd": Pwd, "Code": Code}
	_, err := c.sendData(p)
	if err != nil {
		return err
	}
	return nil
}

// Login 用户登录
func (c *MspKey) Login(Name, Pwd string) error {
	var p SendJson
	p.Type = "Login"
	p.Data = bson.M{"Name": Name, "Pwd": Pwd}
	_, err := c.sendData(p)
	if err != nil {
		return err
	}
	c.IsLogin = true
	return nil

}

// CarLogin 卡密登录
func (c *MspKey) CarLogin(Serial string) error {
	var p SendJson
	p.Type = "CarLogin"
	p.Data = bson.M{"Serial": Serial}
	_, err := c.sendData(p)
	if err != nil {
		return err
	}
	c.IsLogin = true
	return nil
}

// UserPay 用户卡密充值
func (c *MspKey) UserPay(Name, Serial string) error {
	var p SendJson
	p.Type = "UserPay"
	p.Data = bson.M{"Name": Name, "Serial": Serial}
	_, err := c.sendData(p)
	if err != nil {
		return err
	}
	return nil
}

// UpUserPwd 修改密码
func (c *MspKey) UpUserPwd(Name, OldPwd, NewPwd string) error {
	var p SendJson
	p.Type = "UpUserPwd"
	p.Data = bson.M{"Name": Name, "OldPwd": OldPwd, "NewPwd": NewPwd}
	_, err := c.sendData(p)
	if err != nil {
		return err
	}
	return nil
}

// BindDeviceID 换绑
func (c *MspKey) BindDeviceID(Name, Pwd string) error {
	var p SendJson
	p.Type = "BindDeviceID"
	p.Data = bson.M{"Name": Name, "Pwd": Pwd}
	_, err := c.sendData(p)
	if err != nil {
		return err
	}
	return nil
}

// AddBlack 加入黑名单
func (c *MspKey) AddBlack(Bak string) error {
	var p SendJson
	p.Type = "BindDeviceID"
	p.Data = bson.M{"Bak": Bak}
	_, err := c.sendData(p)
	if err != nil {
		return err
	}
	return nil

}

// GetUserInfo 获取用户信息
func (c *MspKey) GetUserInfo() (UserInfo, error) {
	var p SendJson
	p.Type = "GetUserInfo"
	data, err := c.sendData(p)
	if err != nil {
		return UserInfo{}, err
	}
	var st struct {
		UserInfo UserInfo
	}
	marshal, _ := json.Marshal(data.Data)
	_ = json.Unmarshal(marshal, &st)
	return st.UserInfo, nil
}

// SetUerConf 设置用户配置信息
func (c *MspKey) SetUerConf(Conf string) error {
	var p SendJson
	p.Type = "SetUerConf"
	p.Data = bson.M{"Conf": Conf}
	_, err := c.sendData(p)
	if err != nil {
		return err
	}
	return nil

}

// GetExeData 获取核心数据
func (c *MspKey) GetExeData() (string, error) {
	var p SendJson
	p.Type = "GetExeData"
	data, err := c.sendData(p)
	if err != nil {
		return "", err

	}
	var st struct {
		ExeData string
	}
	marshal, _ := json.Marshal(data.Data)
	_ = json.Unmarshal(marshal, &st)
	return st.ExeData, nil

}

// GetVariable 获取远程变量
func (c *MspKey) GetVariable(Key string) (string, error) {
	var p SendJson
	p.Type = "GetVariable"
	p.Data = bson.M{"Key": Key}
	data, err := c.sendData(p)
	if err != nil {
		return "", err
	}
	var st struct {
		ExeData string
	}
	marshal, _ := json.Marshal(data.Data)
	_ = json.Unmarshal(marshal, &st)
	return st.ExeData, nil
}

// GetAdminPay 支付_获取微信二维码链接和发卡链接
func (c *MspKey) GetAdminPay() (PayLink, error) {
	var p SendJson
	p.Type = "GetAdminPay"
	data, err := c.sendData(p)
	if err != nil {
		return PayLink{}, err
	}
	var Link PayLink
	marshal, _ := json.Marshal(data.Data)
	_ = json.Unmarshal(marshal, &Link)
	return Link, nil

}

// FindCarInfo 支付_查询卡密  判断仓库里是否有对应类型的卡密 返回卡ID  	0=小时卡 1=天卡 2=周卡 3=月卡 4=季卡 5=半年卡 6=年卡 7=永久卡
func (c *MspKey) FindCarInfo(CarType int64) (Car, error) {
	var p SendJson
	p.Type = "FindCarInfo"
	p.Data = bson.M{"CarType": CarType}
	data, err := c.sendData(p)
	if err != nil {
		return Car{}, err
	}
	var Car Car
	marshal, _ := json.Marshal(data.Data)
	_ = json.Unmarshal(marshal, &Car)
	return Car, nil

}

// AliPayCreate 支付_创建支付订单  创建支付宝订单 返回base64 图片二维码
func (c *MspKey) AliPayCreate(CarID string) (string, error) {
	var p SendJson
	p.Type = "AliPayCreate"
	p.Data = bson.M{"CarID": CarID}
	data, err := c.sendData(p)
	if err != nil {
		return "", err
	}
	var imgStr string
	marshal, _ := json.Marshal(data.Data)
	_ = json.Unmarshal(marshal, &imgStr)
	return imgStr, nil

}

// IsCompletePayment 是否完成支付
func (c *MspKey) IsCompletePayment() error {
	var p SendJson
	p.Type = "CompletePayment"
	for i, temp := range ResData {
		if temp.Tag == p.Type && temp.Code == 1 {
			c.Msg = temp.Msg
			ResData = append(ResData[:i], ResData[i+1:]...)
			return nil
		} else {
			return errors.New(temp.Msg)
		}
	}
	return errors.New("等待买家付款")
}

// PutCar 支付_提卡 返回卡密信息
func (c *MspKey) PutCar(BuillID string) (Car, error) {
	var p SendJson
	p.Type = "FindCarInfo"
	p.Data = bson.M{"BuillID": BuillID}
	data, err := c.sendData(p)
	if err != nil {
		return Car{}, err
	}
	var Car Car
	marshal, _ := json.Marshal(data.Data)
	_ = json.Unmarshal(marshal, &Car)
	return Car, nil
}

// ping 发送心跳包
func (c *MspKey) ping() {
	var p SendJson
	p.Type = "Ping"
	_, err := c.sendData(p)
	if err != nil {
		log.Println(err)
	}
}

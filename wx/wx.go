package wx

import (
	"crypto/sha1"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/keienWang/wechat/gongNeng"
	"io/ioutil"
	"log"
	"net/http"
	"sort"

	"github.com/clbanning/mxj"
)

var (
	daohang string

	MapZiLiao  = make(map[string]string, 0)
	GuanJianCi = make(map[string]string, 0)
)

func init() {

	MapZiLiao["教师资格证"] = "\n链接: https://pan.baidu.com/s/1kQeqDBm_mTIPPhmsvSmCkA \n提取码: yyu2 \n\n"
	MapZiLiao["Pr教程"] = " \nhttps://pan.baidu.com/s/1pAyzJZFC1mT3BzTqJ6RPLQ \n提取码: h1zw \n\n "
	MapZiLiao["Ps资料"] = "\nhttps://pan.baidu.com/s/1RJ1nz4gv1L-c3RI72oe8XA \n提取码: 0oc9 \n\n"
	MapZiLiao["无水印工具"] = "\nhttps://pan.baidu.com/s/1zx2hUiXH-pI8sLc31uO3Ug \n提取码: x29k  \n\n"
	MapZiLiao["销售人员鸡汤"] = "\nhttps://pan.baidu.com/s/171_NdVbKtHBp-vferjkzLQ \n提取码：4yft \n\n"
	MapZiLiao["英语六级资料"] = "\nhttps://pan.baidu.com/s/193N5nSCXEhHbzp-ohfANbg \n提取码: l4s1  \n\n"
	MapZiLiao["计算机等级考试系统"] = "\n链接: https://pan.baidu.com/s/1Xt5FUgKT2AcpvZajI1eQJg \n提取码: 9u57\n\n"

	daohang = "-----欢迎您关注本公众号-----\n\n想要教程资料或工具的请回复:\"导航\"\n\n也可以回复相应关键词:\n\n"
	for i, _ := range MapZiLiao {
		daohang = daohang + "\"" + i + "\"\n"
	}

	daohang = daohang + "\n流行关键词:\n\n" + "\"舔狗日记\"\n\"渣男语录\"\n\"鸡汤\"\n\"人间凑数\"\n\n\n更多功能正在开发中！"

}

type weixinQuery struct {
	Signature    string `json:"signature"`
	Timestamp    string `json:"timestamp"`
	Nonce        string `json:"nonce"`
	EncryptType  string `json:"encrypt_type"`
	MsgSignature string `json:"msg_signature"`
	Echostr      string `json:"echostr"`
}

type WeixinClient struct {
	Token          string
	Query          weixinQuery
	Message        map[string]interface{}
	Request        *http.Request
	ResponseWriter http.ResponseWriter
	Methods        map[string]func() bool
}

func NewClient(r *http.Request, w http.ResponseWriter, token string) (*WeixinClient, error) {

	weixinClient := new(WeixinClient)

	weixinClient.Token = token
	weixinClient.Request = r
	weixinClient.ResponseWriter = w

	weixinClient.initWeixinQuery()

	if weixinClient.Query.Signature != weixinClient.signature() {
		return nil, errors.New("Invalid Signature.")
	}

	return weixinClient, nil
}

func (this *WeixinClient) initWeixinQuery() {

	var q weixinQuery

	q.Nonce = this.Request.URL.Query().Get("nonce")
	q.Echostr = this.Request.URL.Query().Get("echostr")
	q.Signature = this.Request.URL.Query().Get("signature")
	q.Timestamp = this.Request.URL.Query().Get("timestamp")
	q.EncryptType = this.Request.URL.Query().Get("encrypt_type")
	q.MsgSignature = this.Request.URL.Query().Get("msg_signature")

	this.Query = q
}

func (this *WeixinClient) signature() string {

	strs := sort.StringSlice{this.Token, this.Query.Timestamp, this.Query.Nonce}
	sort.Strings(strs)
	str := ""
	for _, s := range strs {
		str += s
	}
	h := sha1.New()
	h.Write([]byte(str))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (this *WeixinClient) initMessage() error {

	body, err := ioutil.ReadAll(this.Request.Body)

	if err != nil {
		return err
	}

	m, err := mxj.NewMapXml(body)

	if err != nil {
		return err
	}

	if _, ok := m["xml"]; !ok {
		return errors.New("Invalid Message.")
	}

	message, ok := m["xml"].(map[string]interface{})

	if !ok {
		return errors.New("Invalid Field `xml` Type.")
	}

	this.Message = message

	log.Println(this.Message)

	return nil
}

func (this *WeixinClient) text() {

	inMsg, ok := this.Message["Content"].(string)
	fmt.Println("接收到的消息：", inMsg)

	if !ok {
		return
	}

	var reply TextMessage

	if res, ok := MapZiLiao[inMsg]; ok {
		reply.InitBaseData(this, "text")
		reply.Content = value2CDATA(fmt.Sprintf("%s", inMsg+": "+res))

	} else {
		switch inMsg {
		case "舔狗日记":

			dog, err := gongNeng.TianDog()
			if err != nil {
				log.Println(err)
				this.ResponseWriter.WriteHeader(403)
				return
			}

			reply.InitBaseData(this, "text")
			reply.Content = value2CDATA(fmt.Sprintf("%s", dog))
		case "渣男语录":
			words, err := gongNeng.ZhaManWords()
			if err != nil {
				log.Println(err)
				this.ResponseWriter.WriteHeader(403)
				return
			}
			reply.InitBaseData(this, "text")
			reply.Content = value2CDATA(fmt.Sprintf("%s", words))

		case "导航":
			reply.InitBaseData(this, "text")
			var ziliao string

			for i, v := range MapZiLiao {
				ziliao = ziliao + i + ":" + v
			}
			ziliao = ziliao + "\n\n" + "\n流行关键词:\n\n" + "\"舔狗日记\"\n\"渣男语录\"\n\"鸡汤\"\n\"人间凑数\"\n\n\n"

			reply.Content = value2CDATA(fmt.Sprintf("%s", ziliao))
		case "鸡汤":
			tang, err := gongNeng.DuJiTang()
			if err != nil {
				log.Println(err)
				this.ResponseWriter.WriteHeader(403)
				return
			}
			reply.InitBaseData(this, "text")
			reply.Content = value2CDATA(fmt.Sprintf("%s", tang))

		case "人间凑数":

			shu, err := gongNeng.RenJianCouShu()
			if err != nil {
				log.Println(err)
				this.ResponseWriter.WriteHeader(403)
				return
			}
			reply.InitBaseData(this, "text")
			reply.Content = value2CDATA(fmt.Sprintf("%s", shu))
		default:

			ren, err := gongNeng.TencentJiQiRen(inMsg)
			if err != nil {
				log.Println(err)
				this.ResponseWriter.WriteHeader(403)
				return
			}
			reply.InitBaseData(this, "text")
			reply.Content = value2CDATA(fmt.Sprintf("%s", ren))
		}
	}

	replyXml, err := xml.Marshal(reply)

	if err != nil {
		log.Println(err)
		this.ResponseWriter.WriteHeader(403)
		return
	}

	this.ResponseWriter.Header().Set("Content-Type", "text/xml")
	this.ResponseWriter.Write(replyXml)
}

func (this *WeixinClient) event() {

	inMsg, ok := this.Message["Event"].(string)

	if !ok {
		return
	}

	var reply TextMessage
	fmt.Println("\ninMsg : ", inMsg, "\n")

	//订阅回复
	if inMsg == "subscribe" {

		reply.InitBaseData(this, "text")
		reply.Content = value2CDATA(fmt.Sprintf("%s", daohang))

	} else {
		reply.InitBaseData(this, "text")
		reply.Content = value2CDATA(fmt.Sprintf("%s", "不好意思，主人还未解析本操作！"))
	}

	replyXml, err := xml.Marshal(reply)

	if err != nil {
		log.Println(err)
		this.ResponseWriter.WriteHeader(403)
		return
	}

	this.ResponseWriter.Header().Set("Content-Type", "text/xml")
	this.ResponseWriter.Write(replyXml)
}

func (this *WeixinClient) Run() {

	err := this.initMessage()

	if err != nil {

		log.Println(err)
		this.ResponseWriter.WriteHeader(403)
		return
	}

	MsgType, ok := this.Message["MsgType"].(string)

	if !ok {
		this.ResponseWriter.WriteHeader(403)
		return
	}

	switch MsgType {
	case "text":
		this.text()
		break
	case "event":
		this.event()
		break

	default:
		break
	}

	return
}

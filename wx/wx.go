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

var daohang = "欢迎你关注本公众号，想要教程资料或工具的请回复“资料”\t\t\t\t" + "流行关键词分别为:“舔狗日记”；“渣男语录”；“鸡汤”；“人间凑数”，更多功能正在开发中！"

var ziliao = "Pr教程链接: https://pan.baidu.com/s/1pAyzJZFC1mT3BzTqJ6RPLQ 提取码: h1zw \n " +
	"Ps资料链接：https://pan.baidu.com/s/1RJ1nz4gv1L-c3RI72oe8XA 提取码: 0oc9 \n" +
	"无水印工具：https://pan.baidu.com/s/1zx2hUiXH-pI8sLc31uO3Ug 提取码: x29k  \n" +
	"销售必备心灵鸡汤：https://pan.baidu.com/s/171_NdVbKtHBp-vferjkzLQ 提取码：4yft \n" +
	" 英语六级资料：https://pan.baidu.com/s/193N5nSCXEhHbzp-ohfANbg 提取码: l4s1  \n" +
	"计算机等级考试系统：链接: https://pan.baidu.com/s/1Xt5FUgKT2AcpvZajI1eQJg 提取码: 9u57"

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

	case "资料":
		reply.InitBaseData(this, "text")
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
		fmt.Println("1232213123213213")
		shu, err := gongNeng.RenJianCouShu()
		if err != nil {
			log.Println(err)
			this.ResponseWriter.WriteHeader(403)
			return
		}
		reply.InitBaseData(this, "text")
		reply.Content = value2CDATA(fmt.Sprintf("%s", shu))
	default:
		reply.InitBaseData(this, "text")
		reply.Content = value2CDATA(fmt.Sprintf("我收到的是：%s", inMsg))
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

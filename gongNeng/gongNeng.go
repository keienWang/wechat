package gongNeng

import (
	"encoding/json"
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"

	tenerrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"io/ioutil"

	nlp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/nlp/v20190408"

	"net/http"
)

var credential = common.NewCredential(
	"AKIDSGDtJTAKtBev63gESQ84XjfKjmb9nljQ",
	"73h6JNIAmBUVJrAFb6KvFVVCTTCrITWK",
)

func TianDog() (res string, err error) {

	get, err := http.Get("https://v1.alapi.cn/api/dog?format=json")
	if err != nil {

		return
	}
	defer get.Body.Close()

	var data TianGouRiJi
	all, err := ioutil.ReadAll(get.Body)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(all, &data)
	if err != nil {
		return "", err
	}
	return data.Data.Content, nil
}

func ZhaManWords() (res string, err error) {
	var get *http.Response

	get, err = http.Get("https://api.lovelive.tools/api/SweetNothings/Serialization/:serializationType")
	if err != nil {
		get, err = http.Get("https://api.lovelive.tools/api/SweetNothings/Serialization/:serializationType/:count")
		if err != nil {
			get, err = http.Get("https://api.lovelive.tools/api/SweetNothings/:count/Serialization/:serializationType")
			if err != nil {
				return "", err
			}
		}
	}
	defer get.Body.Close()

	var data ZhaManWord
	all, err := ioutil.ReadAll(get.Body)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(all, &data)
	if err != nil {
		return "", err
	}

	return data.ReturnObj[0], nil
}

func DuJiTang() (res string, err error) {

	get, err := http.Get("https://v1.alapi.cn//api/soul")
	if err != nil {
		return
	}
	defer get.Body.Close()

	var data JiTang
	all, err := ioutil.ReadAll(get.Body)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(all, &data)
	if err != nil {
		return "", err
	}
	return data.Data.Title, nil

}

func RenJianCouShu() (res string, err error) {
	get, err := http.Get("http://api.hanpi.top/")
	if err != nil {
		return
	}
	defer get.Body.Close()

	var data RenJian
	all, err := ioutil.ReadAll(get.Body)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(all, &data)
	if err != nil {
		return "", err
	}
	return data.Say, nil
}

func TencentJiQiRen(msg string) (res string, err error) {

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "nlp.tencentcloudapi.com"
	client, _ := nlp.NewClient(credential, "ap-guangzhou", cpf)

	request := nlp.NewChatBotRequest()

	params := "{\"Query\":\"" + msg + "\"}"
	err = request.FromJsonString(params)
	if err != nil {
		return "", err
	}
	response, err := client.ChatBot(request)
	if _, ok := err.(*tenerrors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return "", err
	}
	if err != nil {
		return "", err
	}

	return *response.Response.Reply, nil

}

package gongNeng

import (
	"encoding/json"
	"io/ioutil"

	"net/http"
)

func TianDog()(res string,err error)  {

	get, err := http.Get("https://v1.alapi.cn/api/dog?format=json")
	if err!=nil{

		return
	}
	defer get.Body.Close()


	var data TianGouRiJi
	all, err := ioutil.ReadAll(get.Body)
	if err!=nil{
		return "" ,err
	}
	err= json.Unmarshal(all, &data)
	if err!=nil{
		return "",err
	}
	return  data.Data.Content,nil
}


func ZhaManWords()(res string,err error){
	var get *http.Response

	get, err = http.Get("https://api.lovelive.tools/api/SweetNothings/Serialization/:serializationType")
	if err!=nil{
		get, err = http.Get("https://api.lovelive.tools/api/SweetNothings/Serialization/:serializationType/:count")
		if err!=nil{
			get, err = http.Get("https://api.lovelive.tools/api/SweetNothings/:count/Serialization/:serializationType")
			if err!=nil{
				return "",err
			}
		}
	}
	defer get.Body.Close()

	var data ZhaManWord
	all, err := ioutil.ReadAll(get.Body)
	if err!=nil{
		return "" ,err
	}
	err= json.Unmarshal(all, &data)
	if err!=nil{
		return "",err
	}

	return data.ReturnObj[0],nil
}

func DuJiTang()(res string,err error){

	get, err := http.Get("https://v1.alapi.cn//api/soul")
	if err!=nil{
		return
	}
	defer get.Body.Close()


	var data JiTang
	all, err := ioutil.ReadAll(get.Body)
	if err!=nil{
		return "" ,err
	}
	err= json.Unmarshal(all, &data)
	if err!=nil{
		return "",err
	}
	return  data.Data.Title,nil

}

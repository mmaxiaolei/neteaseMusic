package main

import (
	"fmt"
)

func main() {
	/*	formatStr := fmt.Sprintf("%0256x", 123)
		fmt.Println("format str:", formatStr)
		text := []byte{1, 2, 3, 4, 5}
		api.ReverseBytes(text)
		fmt.Println("text:", text)*/
	//api.Login("张三","123")
	//BpLnfgDsc2WD8F2q
	//var secKey = "BpLnfgDsc2WD8F2q"
	/*	var secKey = "BpLnfgDsc2WD8F2q"
		myParam := api.LoginUser{"412171946@qq.com", "123456", false}
		jsons, _ := json.Marshal(myParam)
		fmt.Println("jsons:", jsons)
		key := api.AesEncrypt(string(jsons), secKey)
		fmt.Println("key:", key)
		originJsons := api.AesDecrypt(key, secKey)
		fmt.Println("origin text:", originJsons)
		var myLoginUser api.LoginUser
		err := json.Unmarshal([]byte(strings.TrimSpace(originJsons)), &myLoginUser)
		if err != nil {
			fmt.Println("json parse error:", err)
		} else {
			fmt.Println("myLoginUser:", myLoginUser)
		}*/

	/*	secc := api.RsaEncrypt(secKey, api.PubKey, api.Modulus)
		fmt.Println("secc:", secc)*/
	/*	key := api.AesEncrypt("test aes enctypt", "BpLnfgDsc2WD8F2q")
		fmt.Println("key:", key)
		originText := api.AesDecrypt(key, "BpLnfgDsc2WD8F2q")
		fmt.Println("text:", originText)*/

	//fmt.Println("-------------request login")
	//api.Login("412171946@qq.com", "a66054577")
	jsonByte := []byte{123, 10, 9, 34, 108, 111, 103, 105, 110, 84, 121, 112, 101, 34, 58, 32, 48, 44, 10, 9, 34, 99, 111, 100, 101, 34, 58, 32, 50, 48, 48, 44, 10, 9, 34, 97, 99, 99, 111, 117, 110, 116, 34, 58, 32, 123, 10, 9, 9, 34, 105, 100, 34, 58, 32, 56, 49, 55, 54, 48, 48, 51, 51, 44, 10, 9, 9, 34, 117, 115, 101, 114, 78, 97, 109, 101, 34, 58, 32, 34, 48, 95, 52, 49, 50, 49, 55, 49, 57, 52, 54, 64, 113, 113, 46, 99, 111, 109, 34, 44, 10, 9, 9, 34, 116, 121, 112, 101, 34, 58, 32, 48, 44, 10, 9, 9, 34, 115, 116, 97, 116, 117, 115, 34, 58, 32, 48, 44, 10, 9, 9, 34, 119, 104, 105, 116, 101, 108, 105, 115, 116, 65, 117, 116, 104, 111, 114, 105, 116, 121, 34, 58, 32, 48, 44, 10, 9, 9, 34, 99, 114, 101, 97, 116, 101, 84, 105, 109, 101, 34, 58, 32, 49, 52, 51, 55, 54, 49, 55, 50, 52, 57, 48, 54, 56, 44, 10, 9, 9, 34, 115, 97, 108, 116, 34, 58, 32, 34, 34, 44, 10, 9, 9, 34, 116, 111, 107, 101, 110, 86, 101, 114, 115, 105, 111, 110, 34, 58, 32, 48, 44, 10, 9, 9, 34, 98, 97, 110, 34, 58, 32, 48, 44, 10, 9, 9, 34, 98, 97, 111, 121, 117, 101, 86, 101, 114, 115, 105, 111, 110, 34, 58, 32, 48, 44, 10, 9, 9, 34, 100, 111, 110, 97, 116, 101, 86, 101, 114, 115, 105, 111, 110, 34, 58, 32, 48, 44, 10, 9, 9, 34, 118, 105, 112, 84, 121, 112, 101, 34, 58, 32, 48, 44, 10, 9, 9, 34, 118, 105, 112, 116, 121, 112, 101, 86, 101, 114, 115, 105, 111, 110, 34, 58, 32, 48, 44, 10, 9, 9, 34, 97, 110, 111, 110, 105, 109, 111, 117, 115, 85, 115, 101, 114, 34, 58, 32, 102, 97, 108, 115, 101, 10, 9, 125, 10, 125}
	fmt.Println(string(jsonByte))
}

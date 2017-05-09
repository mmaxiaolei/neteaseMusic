package main

import (
	"github.com/mmaxiaolei/neteaseMusic/api"
	"fmt"
)

func main() {
	//api.Login("张三","123")
	//BpLnfgDsc2WD8F2q
	key := api.AesEncrypt("test aes enctypt", "BpLnfgDsc2WD8F2q")
	fmt.Println("key:", key)
	originText := api.AesDecrypt(key, "BpLnfgDsc2WD8F2q")
	fmt.Println("text:", originText)
}

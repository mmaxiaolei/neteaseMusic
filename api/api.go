package api

import (
	"net/http"
	"fmt"
	"encoding/json"
	"math/rand"
	"time"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

var versionCookie = &http.Cookie{Name: "appver", Value: "1.5.2"}
var client = http.DefaultClient
var header http.Header = http.Header{
	"Accept":          {"*/*"},
	"Accept-Encoding": {"gzip", "deflate", "sdch"},
	"Accept-Language": {"zh-CN", "zh;q=0.8", "gl;q=0.6", "zh-TW;q=0.4"},
	"Connection":      {"keep-alive"},
	"Content-Type":    {"application/x-www-form-urlencoded"},
	"Host":            {"music.163.com"},
	"Referer":         {"http://music.163.com/search/"},
	"User-Agent":      {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/33.0.1750.152 Safari/537.36"}}

func Login(username string, password string) {
	//loginUser := LoginUser{username, password, true}

	request, _ := http.NewRequest("POST",
		"https://music.163.com/weapi/login?csrf_token=",
		nil)
	request.AddCookie(versionCookie)
	request.Header = header
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(response)
}

// 本地生成指定长度私钥
func CreateSecretKey(size int) string {
	length := len(keys)
	key := make([]byte, size)
	for i := range key {
		// 每次使用不同的随机种子，以避免多次执行的随机数一样
		index := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(length)
		key[i] = keys[index]
	}
	return string(key)
}

const (
	keys   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789" // 生成本地随机密钥的原字符串
	iv     = "0102030405060708"                                               // aes加密用的iv
	nonce  = "0CoJUm6Qyw8W8jud"
	pubKey = "010001"
)

// AES-128-CBC加密
func AesEncrypt(text string, secKey string) string {
	// 加密原文
	plaintext := []byte(text)
	// 加密密钥
	key := []byte(secKey)
	// 填充原文--aes加密的原文必须填充至blocksize的整数倍
	//plaintext = ZeroPadding(plaintext, aes.BlockSize)
	if len(plaintext)%aes.BlockSize != 0 {
		panic("plaintext is not a multiple of the block size")
	}
	// 生成加密用的block
	block, _ := aes.NewCipher(key)
	// aes cbc加密模式
	mode := cipher.NewCBCEncrypter(block, []byte(iv))
	// 密文
	cipherText := make([]byte, len(plaintext))
	mode.CryptBlocks(cipherText, plaintext)
	//return hex.EncodeToString(cipherText)
	//base64.NewEncoding(string(cipherText)).EncodeToString()
	return string(cipherText)
}

// AES-128-CBC解密
func AesDecrypt(cipgerText string, secKey string) string {
	// 加密密钥
	key := []byte(secKey)
	// 密文
	//cipgerTextBytes, _ := hex.DecodeString(cipgerText)
	cipgerTextBytes := []byte(cipgerText)
	if len(cipgerTextBytes) < aes.BlockSize {
		panic("ciphertext too short")
	}
	// CBC mode always works in whole blocks.
	if len(cipgerTextBytes)%aes.BlockSize != 0 {
		panic("ciphertext is not a multiple of the block size")
	}
	// 生成解密的block
	block, _ := aes.NewCipher(key)
	// 解密
	mode := cipher.NewCBCDecrypter(block, []byte(iv))
	// CryptBlocks可以原地更新
	mode.CryptBlocks(cipgerTextBytes, cipgerTextBytes)
	return string(cipgerTextBytes)
}

// aes加密填充：0值填充
func ZeroPadding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{0}, padding) //用0去填充
	return append(cipherText, padText...)
}

// 登录用户信息
type LoginUser struct {
	username      string
	password      string
	rememberLogin bool
}

// 登录参数信息
type LoginParam struct {
	params    string
	encSecKey string
}

// 加密生成登录请求入参
func (loginUser LoginUser) encrypted() LoginParam {
	text, _ := json.Marshal(loginUser)

	return LoginParam{"", string(text)}
}

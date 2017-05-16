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
	"encoding/hex"
	"crypto/md5"
	"math/big"
	"net/url"
	"compress/gzip"
)

var versionCookie = &http.Cookie{Name: "appver", Value: "2.0.2"}
var client = http.DefaultClient
var csrf string

func Login(username string, password string) {
	// 先对password MD5加密
	md5Crypt := md5.New()
	md5Crypt.Write([]byte(password))
	password = hex.EncodeToString(md5Crypt.Sum(nil))
	// 构造明文请求对象
	loginUser := LoginUser{username, true, password}
	// 对请求对象加密
	param := encrypted(loginUser)
	fmt.Println("param", param)
	params := url.Values{
		"params":    {param.Params},
		"encSecKey": {param.EncSecKey}}
	loginUrl := "https://music.163.com/weapi/login"
	request, _ := http.NewRequest("POST",
		loginUrl,
		bytes.NewBufferString(params.Encode()))
	InitHeader(request)
	fmt.Println("request:", request)
	fmt.Println("cookies:", request.Cookies())
	response, _ := client.Do(request)
	var responseAccount ResponseAccount
	rawBody, _ := gzip.NewReader(response.Body)
	json.NewDecoder(rawBody).Decode(&responseAccount)
	// 格式化输出响应json
	body, _ := json.MarshalIndent(responseAccount, "", "\t")
	fmt.Println("body:", string(body))
	// 保存登录信息
}

// 初始化request
func InitHeader(request *http.Request) *http.Request {
	request.Header.Set("Accept", "*/*")
	request.Header.Set("Accept-Encoding", "gzip,deflate,sdch")
	request.Header.Set("Accept-Language", "zh-CN,zh;q=0.8,gl;q=0.6,zh-TW;q=0.4")
	request.Header.Set("Connection", "keep-alive")
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Host", "music.163.com")
	request.Header.Set("Referer", "http://music.163.com/")
	request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/33.0.1750.152 Safari/537.36")
	request.AddCookie(versionCookie)
	return request
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
	keys    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789" // 生成本地随机密钥的原字符串
	iv      = "0102030405060708"                                               // aes加密用的iv
	nonce   = "0CoJUm6Qyw8W8jud"                                               // 初始私钥
	PubKey  = "010001"                                                         // 公钥
	radix   = 16                                                               // 基数
	Modulus = "00e0b509f6259df8642dbc35662901477df22677ec152b5ff68ace615bb7b725152b3ab17a876aea8a5aa76d2e417629ec4ee341f56135fccf695280104e0312ecbda92557c93870114af6c9d05c4f7f0c3685b7a46bee255932575cce10b424d813cfe4875d3e82047b97ddef52741d546b8e289dc6935b3ece0462db0a22b8e7"
)

// AES-128-CBC加密
func AesEncrypt(text string, secKey string) string {
	// 加密原文
	plaintext := []byte(text)
	// 加密密钥
	key := []byte(secKey)
	// 填充原文--aes加密的原文必须填充至blocksize的整数倍
	plaintext = PKCS5Padding(plaintext, aes.BlockSize)
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
	// base64输出
	return base64.StdEncoding.EncodeToString(cipherText)
}

// AES-128-CBC解密
func AesDecrypt(cipgerText string, secKey string) string {
	// 加密密钥
	key := []byte(secKey)
	// 密文
	//cipgerTextBytes, _ := hex.DecodeString(cipgerText)
	cipgerTextBytes, _ := base64.StdEncoding.DecodeString(cipgerText)
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
	// 去掉填充
	cipgerTextBytes = PKCS5UnPadding(cipgerTextBytes)
	return string(cipgerTextBytes)
}

// rsa加密算法：reverse(text)^pubKey%modulus
func RsaEncrypt(text string, exponent string, modulus string) string {
	// 反转text
	textByte := ReverseBytes([]byte(text))
	bigText, _ := big.NewInt(0).SetString(hex.EncodeToString(textByte), 16)
	bigEx, _ := big.NewInt(0).SetString(PubKey, 16)
	bigMod, _ := big.NewInt(0).SetString(Modulus, 16)
	// reverse(text)^pubKey%modulus
	bigRet := big.NewInt(1)
	bigRet.Exp(bigText, bigEx, nil).Mod(bigRet, bigMod)
	// 填充返回
	return fmt.Sprintf("%0256x", bigRet)
}

// aes加密填充：0值填充
func ZeroPadding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{0}, padding) //用0去填充
	return append(cipherText, padText...)
}

// 默认的aes加密填充算法
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(string(ciphertext))%blockSize //需要padding的数目
	// 只要少于256就能放到一个byte中，默认的blockSize=16(即采用16*8=128, AES-128长的密钥)
	// 最少填充1个byte，如果原文刚好是blocksize的整数倍，则再填充一个blocksize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding) //生成填充的文本
	return append(ciphertext, padtext...)
}

// 默认的aes解密去填充算法
func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// 反转byte切片
func ReverseBytes(text []byte) []byte {
	for i, j := 0, len(text)-1; i < j; i, j = i+1, j-1 {
		text[i], text[j] = text[j], text[i]
	}
	return text
}

// 登录用户信息
type LoginUser struct {
	Username      string `json:"username"`
	RememberLogin bool `json:"rememberLogin"`
	Password      string `json:"password"`
}

// 请求参数信息
type Param struct {
	Params    string `json:"params"`
	EncSecKey string `json:"encSecKey"`
}

func (param Param) String() string {
	return "params=>" + param.Params + "======encSecKey=>" + param.EncSecKey
}

// 加密生成请求入参
func encrypted(data interface{}) Param {
	text, _ := json.Marshal(data)
	secKey := CreateSecretKey(16)
	encText := AesEncrypt(AesEncrypt(string(text), nonce), secKey)
	encSecKey := RsaEncrypt(secKey, PubKey, Modulus)
	return Param{encText, encSecKey}
}

// NetEase Music账户信息
type NmAccount struct {
	Id                 int `json:"id"`
	UserName           string `json:"userName"`
	Type               int `json:"type"`
	Status             int `json:"status"`
	WhitelistAuthority int `json:"whitelistAuthority"`
	CreateTime         int `json:"createTime"`
	Salt               string `json:"salt"`
	TokenVersion       int `json:"tokenVersion"`
	Ban                int `json:"ban"`
	BaoyueVersion      int `json:"baoyueVersion"`
	DonateVersion      int `json:"donateVersion"`
	VipType            int `json:"vipType"`
	ViptypeVersion     int `json:"viptypeVersion"`
	AnonimousUser      bool `json:"anonimousUser"`
}

// 登录响应账户信息
type ResponseAccount struct {
	LoginType int `json:"loginType"`
	Code      int `json:"code"`
	Account   NmAccount `json:"account"`
}

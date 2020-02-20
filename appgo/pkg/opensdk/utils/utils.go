package utils

// 通用工具类
import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/i2eco/egoshop/appgo/pkg/opensdk/common"
)

// NewRequest 请求包装
func NewRequest(method, url string, data []byte) (body []byte, err error) {

	if method == "GET" {
		url = fmt.Sprint(url, "?", string(data))
		data = nil
	}

	client := http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		return body, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err = ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return body, err
	}

	return body, err
}

// Request 请求包装体
type Request struct {
	Client *http.Client
}

// NewCertRequest 双向安全证书请求
func NewCertRequest(certFile, keyFile, rootCaFile string) (*Request, error) {

	if certFile == "" || keyFile == "" || rootCaFile == "" {
		return nil, errors.New(common.ErrCertCertEmpty)
	}

	cert, err := ioutil.ReadFile(certFile)
	if err != nil {
		return nil, err
	}

	key, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}

	rootCa, err := ioutil.ReadFile(rootCaFile)
	if err != nil {
		return nil, err
	}

	tlsCert, err := tls.X509KeyPair(cert, key)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	ok := certPool.AppendCertsFromPEM(rootCa)
	if !ok {
		return nil, errors.New("failed to parse root certificate")
	}

	conf := &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		RootCAs:      certPool,
	}
	trans := &http.Transport{
		TLSClientConfig: conf,
	}
	client := &http.Client{
		Transport: trans,
	}

	return &Request{Client: client}, nil
}

// NewRequest 发送请求
func (m *Request) NewRequest(method, url string, data []byte) (body []byte, err error) {
	if method == "GET" {
		url = fmt.Sprint(url, "?", string(data))
		data = nil
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		return body, err
	}
	resp, err := m.Client.Do(req)
	if err != nil {
		return body, err
	}

	body, err = ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return body, err
	}

	return body, err
}

// Struct2Map struct to map，依赖 json tab
func Struct2Map(r interface{}) (s map[string]string, err error) {
	var temp map[string]interface{}
	var result = make(map[string]string)

	bin, err := json.Marshal(r)
	if err != nil {
		return result, err
	}
	if err := json.Unmarshal(bin, &temp); err != nil {
		return nil, err
	}
	for k, v := range temp {
		result[k], err = ToStringE(v)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

// ToStringE interface to string
func ToStringE(i interface{}) (string, error) {
	switch s := i.(type) {
	case string:
		return s, nil
	case bool:
		return strconv.FormatBool(s), nil
	case float64:
		return strconv.FormatFloat(s, 'f', -1, 64), nil
	case float32:
		return strconv.FormatFloat(float64(s), 'f', -1, 32), nil
	case int:
		return strconv.Itoa(s), nil
	case int64:
		return strconv.FormatInt(s, 10), nil
	case int32:
		return strconv.Itoa(int(s)), nil
	case int16:
		return strconv.FormatInt(int64(s), 10), nil
	case int8:
		return strconv.FormatInt(int64(s), 10), nil
	case uint:
		return strconv.FormatInt(int64(s), 10), nil
	case uint64:
		return strconv.FormatInt(int64(s), 10), nil
	case uint32:
		return strconv.FormatInt(int64(s), 10), nil
	case uint16:
		return strconv.FormatInt(int64(s), 10), nil
	case uint8:
		return strconv.FormatInt(int64(s), 10), nil
	case []byte:
		return string(s), nil
	case nil:
		return "", nil
	case fmt.Stringer:
		return s.String(), nil
	case error:
		return s.Error(), nil
	default:
		return "", fmt.Errorf("unable to cast %#v of type %T to string", i, i)
	}
}

// GenWeChatPaySign 生成微信签名
func GenWeChatPaySign(m map[string]string, payKey string) (string, error) {
	delete(m, "sign")
	var signData []string
	for k, v := range m {
		if v != "" {
			signData = append(signData, fmt.Sprintf("%s=%s", k, v))
		}
	}

	sort.Strings(signData)
	signStr := strings.Join(signData, "&")
	signStr = signStr + "&key=" + payKey

	c := md5.New()
	_, err := c.Write([]byte(signStr))
	if err != nil {
		return "", err
	}
	signByte := c.Sum(nil)
	if err != nil {
		return "", err
	}

	sign := strings.ToUpper(hex.EncodeToString(signByte))
	return sign, nil
}

// GetTradeNO 生成订单号，不推荐直接使用
func GetTradeNO(prefix string) string {
	now := time.Now()
	strTime := fmt.Sprintf("%04d%02d%02d%02d%02d%02d", now.Year(), now.Month(), now.Day(), now.Hour(),
		now.Minute(),
		now.Second())
	return prefix + strTime + RandomNumString(100000, 999999)
}

// GetBillNo 生成订单号，指定位数
func GetBillNo(prefix string, length int) string {
	now := time.Now()
	strTime := fmt.Sprintf("%04d%02d%02d%02d%02d%02d", now.Year(), now.Month(), now.Day(), now.Hour(),
		now.Minute(),
		now.Second())

	str := fmt.Sprint(prefix, strTime)
	if a := length - len(str); a > 0 {
		str = str + RandomLenNum(a)
	}
	fmt.Println(str)
	return str[:length]
}

// RandomNum 随机数
func RandomNum(min int64, max int64) int64 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	num := min + r.Int63n(max-min+1)
	return num
}

// RandomNumString 随机数字字符串
func RandomNumString(min int64, max int64) string {
	num := RandomNum(min, max)
	return strconv.FormatInt(num, 10)
}

// RandomLenNum 指定随机数字字符串
func RandomLenNum(length int) string {
	str := ""
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		str += strconv.Itoa(r.Intn(10))
	}
	return str
}

const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

// RandomString 随机字符串
func RandomString(length int) string {
	b := make([]byte, length)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i, cache, remain := length-1, r.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = r.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}

// PKCS7Padding Aes 加密 PKCS7填充
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// PKCS7UnPadding Aes 解密去除PKCS7填充
func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// AesEncrypt Aes 加密
func AesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

// AesDecrypt Aes 解密
func AesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS7UnPadding(origData)
	return origData, nil
}

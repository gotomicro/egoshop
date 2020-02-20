package wechatauth

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/i2eco/egoshop/appgo/pkg/opensdk/common"
	"github.com/i2eco/egoshop/appgo/pkg/opensdk/utils"
)

// refer https://github.com/aimuz/wechat-sdk/blob/master/login/login.go
type (

	// WxConfig 微信配置类
	WxConfig struct {
		AppID  string `json:"appid"`  // 微信APPID
		Secret string `json:"secret"` // 微信Secret
	}

	// WxAccessToken 微信授权Token
	WxAccessToken struct {
		AccessToken  string `json:"access_token,omitempty"`
		ExpiresIn    uint   `json:"expires_in,omitempty"`
		RefreshToken string `json:"refresh_token,omitempty"`
		OpenID       string `json:"openid,omitempty"`
		Scope        string `json:"scope,omitempty"`
		ErrCode      uint   `json:"errcode,omitempty"`
		ErrMsg       string `json:"errmsg,omitempty"`
		ExpiredAt    time.Time
	}

	// WxUserInfo 微信用户资料
	WxUserInfo struct {
		OpenID     string `json:"openid,omitempty"`     // 授权用户唯一标识
		NickName   string `json:"nickname,omitempty"`   // 普通用户昵称
		Sex        uint32 `json:"sex,omitempty"`        // 普通用户性别，1为男性，2为女性
		Province   string `json:"province,omitempty"`   // 普通用户个人资料填写的省份
		City       string `json:"city,omitempty"`       // 普通用户个人资料填写的城市
		Country    string `json:"country,omitempty"`    // 国家，如中国为CN
		HeadImgURL string `json:"headimgurl,omitempty"` // 用户头像，最后一个数值代表正方形头像大小（有0、46、64、96、132数值可选，0代表640*640正方形头像），用户没有头像时该项为空
		//Privilege  string `json:"privilege"`
		Privilege []string `json:"privilege,omitempty"` // 用户特权信息，json数组，如微信沃卡用户为（chinaunicom）
		UnionID   string   `json:"unionid,omitempty"`   // 普通用户的标识，对当前开发者帐号唯一
		ErrCode   uint     `json:"errcode,omitempty"`
		ErrMsg    string   `json:"errmsg,omitempty"`
	}

	// WechatEncryptedData 小程序解密后结构
	WechatEncryptedData struct {
		OpenID    string          `json:"openId"`
		NickName  string          `json:"nickName"`
		Gender    int             `json:"gender"` //性别，0-未知，1-男，2-女
		City      string          `json:"city"`
		Province  string          `json:"province"`
		Country   string          `json:"country"`
		AvatarURL string          `json:"avatarUrl"`
		UnionID   string          `json:"unionId"`
		WaterMark WechatWaterMark `json:"watermark"` //水印
	}

	// WechatWaterMark 加密验证信息
	WechatWaterMark struct {
		Appid     string `json:"appid"`
		Timestamp uint64 `json:"timestamp"`
	}

	// WXBizDataCrypt 小程序解密密钥信息
	WXBizDataCrypt struct {
		Openid     string `json:"openid"`
		SessionKey string `json:"session_key"`
		UnionID    string `json:"unionid"`
		ErrCode    int    `json:"errcode"`
		ErrMsg     string `json:"errmsg"`
	}
)

// AppLogin 微信APP登录 直接登录获取用户信息
func (m *WxConfig) AppLogin(code string) (wxUserInfo *WxUserInfo, err error) {
	return m.LoginCode(code)
}

// WexLogin 微信小程序登录 直接登录获取用户信息
func (m *WxConfig) WexLogin(code, encryptedData, iv string) (wxUserInfo *WechatEncryptedData, err error) {
	wXBizDataCrypt, err := m.GetJsCode2Session(code)
	if err != nil {
		return nil, err
	}
	fmt.Println("fuckkkkkkkkkkk")
	return wXBizDataCrypt.WeDecryptData(encryptedData, iv)
}

// WemLogin 微信网页登录，在微信网页授权，需要认证公众号
func (m *WxConfig) WemLogin(code string) (wxUserInfo *WxUserInfo, err error) {
	return m.LoginCode(code)
}

// LoginCode 通过Code登录
func (m *WxConfig) LoginCode(code string) (wxUserInfo *WxUserInfo, err error) {
	accessToken, err := m.GetWxAccessToken(code)

	if err != nil {
		return wxUserInfo, err
	}
	return accessToken.GetUserInfo()
}

// GetWxAccessToken 通过code获取AccessToken
func (m *WxConfig) GetWxAccessToken(code string) (accessToken *WxAccessToken, err error) {

	if code == "" {
		return accessToken, errors.New("GetWxAccessToken error: code is null")
	}

	params := url.Values{
		"code":       []string{code},
		"grant_type": []string{"authorization_code"},
	}

	t, err := utils.Struct2Map(m)
	if err != nil {
		return accessToken, err
	}

	for k, v := range t {
		params.Set(k, v)
	}

	body, err := utils.NewRequest("GET", common.AccessTokenURL, []byte(params.Encode()))
	if err != nil {
		return accessToken, err
	}

	err = json.Unmarshal(body, &accessToken)
	if err != nil {
		return accessToken, err
	}
	if accessToken.ErrMsg != "" {
		return accessToken, errors.New(accessToken.ErrMsg)
	}

	return
}

// GetUserInfo 获取用户资料
func (m *WxAccessToken) GetUserInfo() (wxUserInfo *WxUserInfo, err error) {
	if m.AccessToken == "" {
		return nil, errors.New("GetUserInfo error: accessToken is null")
	}

	if m.OpenID == "" {
		return nil, errors.New("GetUserInfo error: openID is null")
	}

	params := url.Values{
		"access_token": []string{m.AccessToken},
		"openid":       []string{m.OpenID},
	}
	body, err := utils.NewRequest("GET", common.UserInfoURL, []byte(params.Encode()))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &wxUserInfo)
	if err != nil {
		return nil, err
	}
	if wxUserInfo.OpenID == "" {
		return wxUserInfo, errors.New(wxUserInfo.ErrMsg)
	}

	return
}

// GetRefreshToken 重新获取AccessToken
func (m *WxAccessToken) GetRefreshToken(appid string) error {

	if appid == "" {
		return errors.New(common.ErrAppIDEmpty)
	}

	if m.RefreshToken == "" {
		return errors.New(common.ErrRefreshTokenEmpty)
	}

	if m.OpenID == "" {
		return errors.New(common.ErrOpenIDEmpty)
	}

	params := url.Values{
		"appid":         []string{appid},
		"grant_type":    []string{"refresh_token"},
		"refresh_token": []string{m.RefreshToken},
	}
	body, err := utils.NewRequest("GET", common.RefreshTokenURL, []byte(params.Encode()))
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &m)
	if err != nil {
		return err
	}
	if m.ErrMsg != "" {
		return errors.New(m.ErrMsg)
	}

	return nil
}

// CheckAccessToken 校验AccessToken
func (m *WxAccessToken) CheckAccessToken() (ok bool, err error) {

	if m.AccessToken == "" {
		return ok, errors.New(common.ErrAccessTokenEmpty)
	}

	if m.OpenID == "" {
		return ok, errors.New(common.ErrOpenIDEmpty)
	}

	params := url.Values{
		"openid":       []string{m.OpenID},
		"access_token": []string{m.AccessToken},
	}
	body, err := utils.NewRequest("GET", common.CheckAccessTokenURL, []byte(params.Encode()))
	if err != nil {
		return ok, err
	}

	result := struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return ok, err
	}
	if result.ErrMsg != "ok" {
		return ok, errors.New(result.ErrMsg)
	}

	ok = true
	return
}

// GetJsCode2Session 获取
func (m *WxConfig) GetJsCode2Session(code string) (wXBizDataCrypt *WXBizDataCrypt, err error) {

	if code == "" {
		return wXBizDataCrypt, errors.New("GetJsCode2Session error: code is null")
	}

	params := url.Values{
		"js_code":    []string{code},
		"grant_type": []string{"authorization_code"},
	}

	t, err := utils.Struct2Map(m)
	if err != nil {
		return wXBizDataCrypt, err
	}

	for k, v := range t {
		params.Set(k, v)
	}
	body, err := utils.NewRequest("GET", common.JsCode2SessionURL, []byte(params.Encode()))
	if err != nil {
		return wXBizDataCrypt, err
	}
	err = json.Unmarshal(body, &wXBizDataCrypt)
	if err != nil {
		return wXBizDataCrypt, err
	}

	if wXBizDataCrypt.ErrMsg != "" {
		return wXBizDataCrypt, errors.New(wXBizDataCrypt.ErrMsg)
	}

	return
}

const (
	// IllegalAesKey 解密错误信息
	IllegalAesKey = "encodingAesKey illegal"
)

// WeDecryptData 微信小程序登录数据解密
func (m *WXBizDataCrypt) WeDecryptData(encryptedData, iv string) (wechatEncryptedData *WechatEncryptedData, err error) {

	if len(m.SessionKey) != 24 {
		return wechatEncryptedData, errors.New(IllegalAesKey)
	}

	aesKey, err := base64.StdEncoding.DecodeString(m.SessionKey)
	if err != nil {
		return wechatEncryptedData, err
	}
	aesCipher, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return wechatEncryptedData, err
	}
	aesIV, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		return wechatEncryptedData, err
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		errMsg := fmt.Sprintf("aes new cipher error: %#v", err)
		return wechatEncryptedData, errors.New(errMsg)
	}

	c := cipher.NewCBCDecrypter(block, aesIV)
	resBytes := make([]byte, len(aesCipher))
	c.CryptBlocks(resBytes, aesCipher)
	resBytes = utils.PKCS7UnPadding(resBytes)

	//解密后的byte数组数据做json解析
	wechatEncryptedData = &WechatEncryptedData{}
	err = json.Unmarshal(resBytes, &wechatEncryptedData)
	if err != nil {
		errMsg := fmt.Sprintf("json unmarshal data error: %#v", err)
		return wechatEncryptedData, errors.New(errMsg)
	}

	return
}

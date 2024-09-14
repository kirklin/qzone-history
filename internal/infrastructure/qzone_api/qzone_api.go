package qzone_api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"qzone-history/internal/domain/entity"
	"qzone-history/internal/infrastructure/config"
	"qzone-history/pkg/utils"
	"regexp"
	"strings"
	"time"
)

type QzoneAPIClient interface {
	GetLoginQRCode() ([]byte, string, error)
	CheckLoginStatus(qrsig string) (entity.LoginStatus, string, error)
	CompleteLogin(responseText string) (map[string]string, error)
	GetUserInfo(cookies map[string]string) (*entity.User, error)
}

type qzoneAPIClient struct {
	httpClient *http.Client
	config     *config.Config
}

func NewQzoneAPIClient(config *config.Config) QzoneAPIClient {
	return &qzoneAPIClient{
		httpClient: &http.Client{
			Timeout: time.Second * 10,
		},
		config: config,
	}
}

func (c *qzoneAPIClient) GetLoginQRCode() ([]byte, string, error) {
	resp, err := c.httpClient.Get(c.config.QzoneAPI.QRCodeURL)
	if err != nil {
		return nil, "", fmt.Errorf("获取二维码失败: %w", err)
	}
	defer resp.Body.Close()

	var qrsig string
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "qrsig" {
			qrsig = cookie.Value
			break
		}
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("读取二维码数据失败: %w", err)
	}

	return body, qrsig, nil
}

func (c *qzoneAPIClient) CheckLoginStatus(qrsig string) (entity.LoginStatus, string, error) {
	ptqrtoken := utils.GeneratePtqrToken(qrsig)

	loginURL := fmt.Sprintf("%s?u1=https%%3A%%2F%%2Fqzs.qq.com%%2Fqzone%%2Fv5%%2Floginsucc.html%%3Fpara%%3Dizone&ptqrtoken=%s&ptredirect=0&h=1&t=1&g=1&from_ui=1&ptlang=2052&action=0-0-%d&js_ver=20032614&js_type=1&login_sig=&pt_uistyle=40&aid=549000912&daid=5&",
		c.config.QzoneAPI.LoginURL, ptqrtoken, time.Now().Unix())

	req, err := http.NewRequest("GET", loginURL, nil)
	if err != nil {
		return entity.LoginStatusWaiting, "", fmt.Errorf("创建登录请求失败: %w", err)
	}

	req.AddCookie(&http.Cookie{Name: "qrsig", Value: qrsig})
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return entity.LoginStatusWaiting, "", fmt.Errorf("发送登录请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return entity.LoginStatusWaiting, "", fmt.Errorf("读取登录响应失败: %w", err)
	}

	responseText := string(body)
	if strings.Contains(responseText, "二维码未失效") {
		return entity.LoginStatusWaiting, "", nil
	} else if strings.Contains(responseText, "二维码认证中") {
		return entity.LoginStatusScanning, "", nil
	} else if strings.Contains(responseText, "二维码已失效") {
		return entity.LoginStatusExpired, "", nil
	} else if strings.Contains(responseText, "登录成功") {
		return entity.LoginStatusSuccess, responseText, nil
	}

	return entity.LoginStatusWaiting, "", nil
}

func (c *qzoneAPIClient) CompleteLogin(responseText string) (map[string]string, error) {
	re := regexp.MustCompile(`ptsigx=(.*?)&`)
	matches := re.FindStringSubmatch(responseText)
	if len(matches) < 2 {
		return nil, fmt.Errorf("无法提取 ptsigx")
	}
	sigx := matches[1]

	re = regexp.MustCompile(`uin=(\d+)`)
	matches = re.FindStringSubmatch(responseText)
	if len(matches) < 2 {
		return nil, fmt.Errorf("无法获取 uin")
	}
	uin := matches[1]

	checkSigURL := fmt.Sprintf("https://ptlogin2.qzone.qq.com/check_sig?pttype=1&uin=%s&service=ptqrlogin&nodirect=0&ptsigx=%s&s_url=https%%3A%%2F%%2Fqzs.qq.com%%2Fqzone%%2Fv5%%2Floginsucc.html%%3Fpara%%3Dizone&f_url=&ptlang=2052&ptredirect=100&aid=549000912&daid=5&j_later=0&low_login_hour=0&regmaster=0&pt_login_type=3&pt_aid=0&pt_aaid=16&pt_light=0&pt_3rd_aid=0",
		uin, sigx)

	req, err := http.NewRequest("GET", checkSigURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建 check_sig 请求失败: %w", err)
	}

	c.httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送 check_sig 请求失败: %w", err)
	}
	defer resp.Body.Close()

	cookies := make(map[string]string)
	for _, cookie := range resp.Cookies() {
		if cookie.Value != "" || cookies[cookie.Name] == "" {
			cookies[cookie.Name] = cookie.Value
		}
	}

	c.httpClient.CheckRedirect = nil
	return cookies, nil
}

func (c *qzoneAPIClient) GetUserInfo(cookies map[string]string) (*entity.User, error) {
	uin := utils.ExtractUin(cookies)
	skey := cookies["p_skey"]
	g_tk := utils.GenerateGTK(skey)

	url := fmt.Sprintf("https://r.qzone.qq.com/fcg-bin/cgi_get_portrait.fcg?g_tk=%s&uins=%s", g_tk, uin)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建获取用户信息请求失败: %w", err)
	}

	for name, value := range cookies {
		req.AddCookie(&http.Cookie{Name: name, Value: value})
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送获取用户信息请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %w", err)
	}

	info := string(body)
	info = strings.TrimSpace(info)

	if strings.HasPrefix(info, "portraitCallBack(") {
		info = strings.TrimPrefix(info, "portraitCallBack(")
		info = strings.TrimSuffix(info, ")")
	}

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(info), &result); err != nil {
		return nil, fmt.Errorf("解析用户信息响应失败: %w", err)
	}

	userData, ok := result[uin].([]interface{})
	if !ok || len(userData) < 7 {
		return nil, fmt.Errorf("用户信息格式不正确")
	}

	nickname, ok := userData[6].(string)
	if !ok {
		return nil, fmt.Errorf("无法获取昵称")
	}

	return &entity.User{
		QQ:       uin,
		Nickname: nickname,
		Cookies:  cookies,
	}, nil
}

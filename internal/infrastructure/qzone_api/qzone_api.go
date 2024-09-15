package qzone_api

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/schollz/progressbar/v3"
	"io"
	"net/http"
	"qzone-history/internal/domain/entity"
	"qzone-history/internal/infrastructure/config"
	"qzone-history/pkg/utils"
	"regexp"
	"strings"
	"time"
)

type QzoneAPIClient interface {
	// GetLoginQRCode 获取登录二维码
	GetLoginQRCode() ([]byte, string, error)

	// CheckLoginStatus 检查登录二维码状态
	CheckLoginStatus(qrsig string) (entity.LoginStatus, string, error)

	// CompleteLogin 完成登录并返回 cookies
	CompleteLogin(responseText string) (map[string]string, error)

	// GetUserInfo 获取用户信息
	GetUserInfo(cookies map[string]string) (*entity.User, error)

	// GetActivities 获取用户活动
	GetActivities(cookies map[string]string, offset, count int) ([]*entity.Activity, error)

	// GetAllActivities 获取全部用户活动
	GetAllActivities(cookies map[string]string) ([]*entity.Activity, error)
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

	body, err := io.ReadAll(resp.Body)
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

	body, err := io.ReadAll(resp.Body)
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

	body, err := io.ReadAll(resp.Body)
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

func (c *qzoneAPIClient) getActivityCount(cookies map[string]string) (int, error) {
	lowerBound := 0
	upperBound := 10000000
	total := upperBound / 2

	for lowerBound <= upperBound {
		activities, err := c.GetActivities(cookies, total, 100)
		if err != nil {
			return 0, fmt.Errorf("获取活动失败: %w", err)
		}

		if len(activities) > 0 {
			lowerBound = total + 1
		} else {
			upperBound = total - 1
		}

		total = (lowerBound + upperBound) / 2
	}

	return total, nil
}

func (c *qzoneAPIClient) GetAllActivities(cookies map[string]string) ([]*entity.Activity, error) {
	var allActivities []*entity.Activity
	totalCount, err := c.getActivityCount(cookies)
	if err != nil {
		return nil, fmt.Errorf("获取活动总数失败: %w", err)
	}

	bar := progressbar.Default(int64(totalCount), "正在获取活动")

	for i := 0; i <= totalCount/100; i++ {
		activities, err := c.GetActivities(cookies, i*100, 100)
		if err != nil {
			return nil, fmt.Errorf("获取活动失败 (offset %d): %w", i*100, err)
		}
		allActivities = append(allActivities, activities...)

		// 更新进度条
		bar.Add(len(activities))

		time.Sleep(200 * time.Millisecond)
	}

	fmt.Println("\n获取完成!")
	return allActivities, nil
}

func (c *qzoneAPIClient) GetActivities(cookies map[string]string, offset, count int) ([]*entity.Activity, error) {
	// 提取用户 QQ 和 g_tk
	uin := utils.ExtractUin(cookies)
	g_tk := utils.GenerateGTK(cookies["p_skey"])

	// 构造请求 URL，使用 offset 和 count 来控制分页
	url := fmt.Sprintf(
		"https://user.qzone.qq.com/proxy/domain/ic2.qzone.qq.com/cgi-bin/feeds/feeds2_html_pav_all?uin=%s&begin_time=0&end_time=0&getappnotification=1&getnotifi=1&has_get_key=0&offset=%d&set=0&count=%d&useutf8=1&outputhtmlfeed=1&scope=1&format=jsonp&g_tk=%s",
		uin, offset, count, g_tk,
	)

	// 创建请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 添加请求 Cookies
	for name, value := range cookies {
		req.AddCookie(&http.Cookie{Name: name, Value: value})
	}

	// 发送请求
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	// 处理响应 HTML
	processedHTML := utils.ProcessOldHTML(string(body))
	if !strings.Contains(processedHTML, "li") {
		// 如果没有活动返回空
		return nil, nil
	}

	// 解析 HTML
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(processedHTML))
	if err != nil {
		return nil, fmt.Errorf("解析HTML失败: %w", err)
	}

	var activities []*entity.Activity

	// 遍历每个活动条目
	doc.Find("li.f-single.f-s-s").Each(func(i int, s *goquery.Selection) {
		activity := &entity.Activity{}
		activity.ReceiverQQ = uin
		// 提取发送者的信息（昵称、QQ号、链接）
		senderElement := s.Find("a.f-name.q_namecard")
		if senderElement.Length() > 0 {
			activity.SenderName = senderElement.Text()
			activity.SenderQQ = strings.TrimPrefix(senderElement.AttrOr("link", ""), "nameCard_")
			activity.SenderLink = senderElement.AttrOr("href", "")
		}

		// 提取时间信息
		timeElement := s.Find("div.info-detail")
		if timeElement.Length() > 0 {
			activity.TimeText = strings.TrimSpace(timeElement.Text())
			activity.Timestamp = parseTime(activity.TimeText)
		}

		// 提取内容
		contentElement := s.Find("p.txt-box-title.ellipsis-one")
		if contentElement.Length() > 0 {
			activity.Content = strings.TrimSpace(contentElement.Text())
			activity.Content = strings.ReplaceAll(contentElement.Text(), "\u00a0", " ")
		}

		// 提取图片 URL
		imgElements := s.Find("a.img-item img")
		imgElements.Each(func(i int, img *goquery.Selection) {
			if src, exists := img.Attr("src"); exists {
				activity.ImageURLs = append(activity.ImageURLs, src)
			}
		})

		// 判断活动类型
		stateElement := s.Find("span.state")
		stateText := stateElement.Text()

		switch {
		case strings.Contains(stateText, "赞了我的说说"):
			activity.Type = entity.TypeLike
		case strings.Contains(stateText, "查看了我的说说"):
			activity.Type = entity.TypeView
		case strings.Contains(stateText, "评论"):
			activity.Type = entity.TypeComment
		case strings.Contains(stateText, "留言"):
			activity.Type = entity.TypeBoardMessage
		case strings.Contains(stateText, "回复"):
			activity.Type = entity.TypeReply
		case s.Find("div.f-reprint").Length() > 0:
			activity.Type = entity.TypeForward
			// 提取转发内容
			forwardContent := s.Find("div.f-reprint div.f-info").Text()
			activity.Content = strings.TrimSpace(forwardContent)
		default:
			activity.Type = entity.TypeOther
		}

		// 保存当前解析的活动
		activities = append(activities, activity)
	})

	return activities, nil
}

func parseTime(timeStr string) time.Time {
	now := time.Now()
	layouts := []string{
		"2006年1月2日 15:04",
		"2006年01月02日 15:04",
		"1月2日 15:04",
		"01月02日 15:04",
		"昨天 15:04",
		"15:04",
	}

	for _, layout := range layouts {
		t, err := time.ParseInLocation(layout, timeStr, time.Local)
		if err == nil {
			switch layout {
			case "2006年1月2日 15:04", "2006年01月02日 15:04":
				// 完整日期，直接返回
				return t
			case "1月2日 15:04", "01月02日 15:04":
				// 只有月日，使用当前年份
				return time.Date(now.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, time.Local)
			case "昨天 15:04":
				// 昨天，使用当前年月，日期减一
				yesterday := now.AddDate(0, 0, -1)
				return time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), t.Hour(), t.Minute(), 0, 0, time.Local)
			case "15:04":
				// 只有时间，使用当前年月日
				return time.Date(now.Year(), now.Month(), now.Day(), t.Hour(), t.Minute(), 0, 0, time.Local)
			}
		}
	}

	// 如果所有格式都无法解析，返回零值
	return time.Time{}
}

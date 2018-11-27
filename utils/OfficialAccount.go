package utils

import (
	"encoding/json"
	"strconv"
	"strings"
)

var (
	DefaultPage    = 1  // 获取哪页的文章
	DefaultPerPage = 50 // 每页显示多少文章
)

type OfficialAccount struct {
	*QingboClient
}

type AccountData struct {
	VerifyName    string   `json:"verify_name"`
	WxName        string   `json:"wx_name"`
	AddTime       string   `json:"add_time"`
	WxVip         string   `json:"wx_vip"`
	WxNote        string   `json:"wx_note"`
	WxLogo        string   `json:"wx_logo"`
	Wci           float64  `json:"wci,omitempty"`
	WxNickname    string   `json:"wx_nickname,omitempty"`
	NicknameId    string   `json:"nickname_id"`
	WxQrcode      string   `json:"wx_qrcode"`
	WxBiz         string   `json:"wx_biz"`
	WxAccountTags []string `json:"wx_account_tags,omitempty"`
}

type AccountResponse struct {
	DataResp    []*AccountData `json:"data"`
	Url         string         `json:"url"`
	Application string         `json:"application"`
}

type ArticleData struct {
	Id            string `json:"id"`              // 微信文章ID
	Top           int64  `json:"top"`             // 文章位置
	Url           string `json:"url"`             // 微信地址
	Name          string `json:"name"`            // 微信昵称
	Type          string `json:"type"`            // 文章类型
	Title         string `json:"title"`           // 标题
	Author        string `json:"author"`          // 作者
	Picurl        string `json:"picurl"`          // 图片链接
	Digest        string `json:"digest"`          // 描述
	WxName        string `json:"wx_name"`         // 微信账号
	ReadCount     int64  `json:"read_count"`      // 阅读数
	LikeCount     int64  `json:"like_count"`      // 点赞数
	CreatedAt     string `json:"created_at"`      // 发布时间
	OriginalUrl   string `json:"original_url"`    // 原始微信链接
	WeekReadCount int64  `json:"week_read_count"` // 周阅读数
	WeekLikeCount int64  `json:"week_like_count"` // 周点赞数
}

type ArticleResponse struct {
	DataResp    []*ArticleData `json:"data"`
	Url         string         `json:"url"`
	Application string         `json:"application"`
}

func (a *OfficialAccount) GetAccount(accountName string) (*AccountResponse, error) {
	resp, err := a.QingboClient.get("users", "wx_name="+accountName, "weixin")
	if err != nil {
		return nil, err
	}
	wechatResp := &AccountResponse{}
	err = json.Unmarshal([]byte(resp), wechatResp)
	if err != nil {
		return nil, err
	}
	return wechatResp, nil
}

func (a *OfficialAccount) GetArticles(wxName, startDate string, perPage, page int) (*ArticleResponse, error) {
	sb := strings.Builder{}
	sb.WriteString("wx_name=" + wxName)
	if startDate != "" {
		sb.WriteString("&start_date=" + startDate)
	}
	if perPage > 0 {
		sb.WriteString("&per-page=" + strconv.Itoa(perPage))
	} else {
		sb.WriteString("&per-page=" + strconv.Itoa(DefaultPerPage))
	}
	if page > 0 {
		sb.WriteString("&page=" + strconv.Itoa(page))
	} else {
		sb.WriteString("&page=" + strconv.Itoa(DefaultPage))
	}
	resp, err := a.QingboClient.get("articles", sb.String(), "weixin")
	if err != nil {
		return nil, err
	}
	articlesResp := &ArticleResponse{}
	err = json.Unmarshal([]byte(resp), articlesResp)
	if err != nil {
		return nil, err
	}
	return articlesResp, nil
}

func NewOfficialAccount(client *QingboClient) *OfficialAccount {
	return &OfficialAccount{client}
}

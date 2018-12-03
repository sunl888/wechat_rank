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

type ArticleResponse struct {
	Id        string `json:"id"`         // 微信文章ID
	Title     string `json:"title"`      // 标题
	Url       string `json:"url"`        // 微信url
	WxName    string `json:"wx_name"`    // 微信账号
	Top       int64  `json:"top"`        // 文章位置
	ReadCount int64  `json:"read_count"` // 阅读数
	LikeCount int64  `json:"like_count"` // 点赞数
	CreatedAt string `json:"created_at"` // 发布时间
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

func (a *OfficialAccount) GetArticles(wxName, startDate, endDate string, perPage, page int) ([]*ArticleResponse, error) {
	sb := strings.Builder{}
	if wxName != "" {
		sb.WriteString("wx_name=" + wxName)
	}
	if startDate != "" {
		sb.WriteString("&start_date=" + startDate)
	}
	if endDate != "" {
		sb.WriteString("&end_date=" + endDate)
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
	articlesResp := &struct {
		Data []*ArticleResponse `json:"data"`
	}{}
	err = json.Unmarshal([]byte(resp), articlesResp)
	if err != nil {
		return nil, err
	}
	return articlesResp.Data, nil
}

func NewOfficialAccount(client *QingboClient) *OfficialAccount {
	return &OfficialAccount{client}
}

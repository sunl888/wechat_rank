package qingbo

import (
	"code.aliyun.com/zmdev/wechat_rank/errors"
	"encoding/json"
	"strconv"
)

var (
	Page    = 1
	PerPage = 50
)

type WxAccount struct {
	*client
}

type AccountData struct {
	VerifyName string `json:"verify_name"`
	WxName     string `json:"wx_name"`
	AddTime    string `json:"add_time"`
	WxNickname string `json:"wx_nickname"`
	NicknameId string `json:"nickname_id"`
	WxBiz      string `json:"wx_biz"`
	WxVip      string `json:"wx_vip"`
	WxNote     string `json:"wx_note"`
	WxLogo     string `json:"wx_logo"`
	WxQrcode   string `json:"wx_qrcode"`
}

type AccountResponse struct {
	Data        []*AccountData `json:"data"`
	Url         string         `json:"url"`
	Application string         `json:"application"`
}

type ArticleResponse struct {
	Id        string `json:"id"`         // 微信文章ID
	Title     string `json:"title"`      // 标题
	Digest    string `json:"digest"`     // 描述(摘要)
	Url       string `json:"url"`        // 微信url
	WxName    string `json:"wx_name"`    // 微信账号
	Top       int64  `json:"top"`        // 文章位置
	ReadCount int64  `json:"read_count"` // 阅读数
	LikeCount int64  `json:"like_count"` // 点赞数
	CreatedAt string `json:"created_at"` // 发布时间
}

type RankDay struct {
	WxNickname string `json:"wx_nickname"`
	// ...
}

type RankDayResponse struct {
	Data        []*RankDay `json:"data"`
	Url         string     `json:"url"`
	Application string     `json:"application"`
}

type ErrorResponse struct {
	Name    string
	Message string
	Code    int
	Status  int
}

func (a *WxAccount) GetAccount(accountName string) (*AccountResponse, error) {
	resp, err := a.client.get("users", map[string]string{
		"wx_name": accountName,
	}, "weixin")
	if err != nil {
		return nil, err
	}
	wechatResp := &AccountResponse{}
	err = json.Unmarshal([]byte(resp), wechatResp)
	if err != nil {
		return nil, err
	}
	if wechatResp.Data == nil {
		errorResp := &ErrorResponse{}
		err = json.Unmarshal([]byte(resp), errorResp)
		if err != nil {
			return nil, err
		}
		return nil, errors.QingboError(errorResp.Name, errorResp.Message, errorResp.Code, errorResp.Status)
	}
	return wechatResp, nil
}

func (a *WxAccount) GetArticles(wxName, startDate, endDate string, perPage, page int) ([]*ArticleResponse, error) {
	var params map[string]string
	params = make(map[string]string, 3)
	if wxName != "" {
		params["wx_name"] = wxName
	} else {
		return nil, errors.QingboError("", "wx_name 不能为空", 400, 400)
	}
	if startDate != "" {
		params["start_date"] = startDate
	}
	if endDate != "" {
		params["end_date"] = endDate
	}
	if perPage > 0 {
		params["per-page"] = strconv.Itoa(perPage)
	} else {
		params["per-page"] = strconv.Itoa(PerPage)
	}
	if page > 0 {
		params["page"] = strconv.Itoa(page)
	} else {
		params["page"] = strconv.Itoa(Page)
	}
	resp, err := a.client.get("articles", params, "weixin")
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
	if articlesResp.Data == nil {
		errorResp := &ErrorResponse{}
		err = json.Unmarshal([]byte(resp), errorResp)
		if err != nil {
			return nil, err
		}
		return nil, errors.QingboError(errorResp.Name, errorResp.Message, errorResp.Code, errorResp.Status)
	}
	return articlesResp.Data, nil
}

func (a *WxAccount) GetRankDays(wxName, startDate string) (*RankDayResponse, error) {
	params := make(map[string]string, 2)
	if wxName != "" {
		params["wx_name"] = wxName
	}
	if startDate != "" {
		params["start_date"] = startDate
	}
	resp, err := a.client.get("users/rank-days", params, "weixin")
	if err != nil {
		return nil, err
	}
	rankDayResp := &RankDayResponse{}
	err = json.Unmarshal([]byte(resp), rankDayResp)
	if err != nil {
		return nil, err
	}
	if rankDayResp.Data == nil {
		errorResp := &ErrorResponse{}
		err = json.Unmarshal([]byte(resp), errorResp)
		if err != nil {
			return nil, err
		}
		return nil, errors.QingboError(errorResp.Name, errorResp.Message, errorResp.Code, errorResp.Status)
	}
	return rankDayResp, nil
}

func NewWxAccount(client *client) *WxAccount {
	return &WxAccount{client}
}

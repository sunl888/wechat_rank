package utils

import (
	"encoding/json"
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
	DataResp    []AccountData `json:"data"`
	Url         string        `json:"url"`
	Application string        `json:"application"`
}

func (a *OfficialAccount) GetOfficialAccount(accountName string) (*AccountResponse, error) {
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

func NewOfficialAccount(client *QingboClient) *OfficialAccount {
	return &OfficialAccount{client}
}

package config

import (
	"github.com/TangSengDaoDao/TangSengDaoDaoServerLib/common"
	"github.com/TangSengDaoDao/TangSengDaoDaoServerLib/pkg/util"
)

// SendTradeSystemNotifyTemplate 交易通知模版
func (c *Context) SendTradeSystemNotifyTemplate(msg MsgTradeSystemNotifyTemplate) error {
	return c.SendMessage(&MsgSendReq{
		Header: MsgHeader{
			RedDot: 1,
		},
		ChannelID:   msg.ChannelID,
		ChannelType: msg.ChannelType,
		FromUID:     c.cfg.Account.SystemUID,
		Payload: []byte(util.ToJson(map[string]interface{}{
			"left_title":      msg.LeftTitle,
			"left_subtitle":   msg.LeftSubtitle,
			"center_title":    msg.CenterTitle,
			"center_subtitle": msg.CenterSubtitle,
			"url_title":       msg.URLTitle,
			"notice":          msg.Notice,
			"imprest_code":    msg.TradeNo,
			// "type":            common.TradeSystemNotifyTemplate,
			"attrs":   msg.Attrs,
			"type":    common.RedpacketReceive,
			"content": "{0}的红包24小时未领取，已退回。",
			"extra": []UserBaseVo{
				{
					UID:  msg.Attrs["creater"],
					Name: msg.Attrs["creater_name"],
				},
			},
		})),
	})
}

// MsgTradeSystemNotifyTemplate 交易通知模版
type MsgTradeSystemNotifyTemplate struct {
	ChannelID      string
	ChannelType    uint8
	LeftTitle      string // 左边标题
	LeftSubtitle   string // 左边子标题
	CenterTitle    string // 中间标题
	CenterSubtitle string // 中间子标题
	URLTitle       string // url标题
	Notice         string // 通知（列表里的提示）
	TradeNo        string // 交易编号
	Attrs          map[string]string
}

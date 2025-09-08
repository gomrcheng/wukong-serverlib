package config

import (
	"fmt"
	"time"

	"github.com/TangSengDaoDao/TangSengDaoDaoServerLib/common"
	"github.com/TangSengDaoDao/TangSengDaoDaoServerLib/pkg/util"
)

func (c *Context) sendRedpacketReceiveForPerson(msg MsgRedpacketReceive) error {

	// 对方领取了红包，发送给创建者
	messageMap := c.getRecveivePayload(msg, `“{0}“领取了“{1}“的红包`, msg.Receiver, msg.ReceiverName, msg.Creater, msg.CreaterName, []string{msg.CreaterName, msg.ReceiverName})
	err := c.SendMessage(&MsgSendReq{
		Header: MsgHeader{
			RedDot: 1,
		},
		ChannelID:   msg.Creater,
		ChannelType: common.ChannelTypePerson.Uint8(),
		FromUID:     msg.Receiver,
		Payload:     []byte(util.ToJson(messageMap)),
	})
	if err != nil {
		return err
	}

	// // 对方领取了红包，发送给领取者 你领取了xxx的红包
	// messageMap = c.getRecveivePayload(msg, `你领取了“{0}“的红包`, msg.Creater, msg.CreaterName, []string{msg.Receiver})
	// err = c.SendMessage(&MsgSendReq{
	// 	Header: MsgHeader{
	// 		RedDot: 1,
	// 	},
	// 	ChannelID:   msg.Receiver,
	// 	ChannelType: common.ChannelTypePerson.Uint8(),
	// 	FromUID:     msg.Creater,
	// 	Payload:     []byte(util.ToJson(messageMap)),
	// })
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (c *Context) sendRedpacketReceiveForGroup(msg MsgRedpacketReceive) error {
	// 自己领取了自己的红包
	// if msg.Receiver == msg.Creater {
	// 	messageMap := c.getRecveivePayload(msg, `你领取了自己的红包`, msg.Receiver, msg.ReceiverName, []string{msg.Receiver})
	// 	return c.SendMessage(&MsgSendReq{
	// 		Header: MsgHeader{
	// 			RedDot: 1,
	// 		},
	// 		ChannelID:   msg.ChannelID,
	// 		ChannelType: msg.ChannelType,
	// 		// Subscribers: []string{msg.Receiver},
	// 		Subscribers: []string{},
	// 		Payload:     []byte(util.ToJson(messageMap)),
	// 	})

	// }
	// messageMap := c.getRecveivePayload(msg, `你领取了“{0}“的红包`, msg.Creater, msg.CreaterName, []string{msg.Receiver})
	// err := c.SendMessage(&MsgSendReq{
	// 	Header: MsgHeader{
	// 		RedDot: 1,
	// 	},
	// 	ChannelID:   msg.ChannelID,
	// 	ChannelType: msg.ChannelType,
	// 	// Subscribers: []string{msg.Receiver},
	// 	Subscribers: []string{},
	// 	Payload:     []byte(util.ToJson(messageMap)),
	// })
	// if err != nil {
	// 	return err
	// }

	messageMap := c.getRecveivePayload(msg, `“{0}“领取了“{1}“的红包`, msg.Receiver, msg.ReceiverName, msg.Creater, msg.CreaterName, []string{msg.CreaterName, msg.ReceiverName})
	return c.SendMessage(&MsgSendReq{
		Header: MsgHeader{
			RedDot: 1,
		},
		ChannelID:   msg.ChannelID,
		ChannelType: msg.ChannelType,
		FromUID:     msg.Receiver,
		// Subscribers: []string{msg.Creater},
		Subscribers: []string{},
		Payload:     []byte(util.ToJson(messageMap)),
	})
}

func (c *Context) getRecveivePayload(msg MsgRedpacketReceive, content string, uid string, name string, creater string, createrName string, visibles []string) map[string]interface{} {
	messageMap := map[string]interface{}{
		"redpacket_no": msg.RecordNo,
		"type":         common.RedpacketReceive,
		"content":      content,
		"extra": []UserBaseVo{
			{
				UID:  uid,
				Name: name,
			},
			{
				UID:  creater,
				Name: createrName,
			},
		},
	}
	if len(visibles) > 0 {
		messageMap["visibles"] = visibles
	}
	return messageMap
}

// SendRedpacketReceive 发送红包领取消息
func (c *Context) SendRedpacketReceive(msg MsgRedpacketReceive) error {
	if msg.ChannelType == common.ChannelTypePerson.Uint8() {
		return c.sendRedpacketReceiveForPerson(msg)
	} else if msg.ChannelType == common.ChannelTypeGroup.Uint8() {
		return c.sendRedpacketReceiveForGroup(msg)
	} else {
		return fmt.Errorf("不支持的频道类型: %d", msg.ChannelType)
	}
}

// SendRedpacketReceive 发送红包领取消息
// func (c *Context) SendRedpacketReceive(msg MsgRedpacketReceive) error {

// 	content := fmt.Sprintf(`“{0}“领取了你的红包`)
// 	if msg.Receiver == msg.Creater {
// 		content = fmt.Sprintf(`你领取了自己发的红包`)
// 	}
// 	return c.SendMessage(&MsgSendReq{
// 		Header: MsgHeader{
// 			RedDot: 1,
// 		},
// 		ChannelID:   msg.ChannelID,
// 		ChannelType: msg.ChannelType,
// 		Subscribers: []string{msg.Creater},
// 		Payload: []byte(util.ToJson(map[string]interface{}{
// 			"redpacket_no": msg.RecordNo,
// 			"type":         common.RedpacketReceive,
// 			"content":      content,
// 			"extra": []UserBaseVo{
// 				{
// 					UID:  msg.Receiver,
// 					Name: msg.ReceiverName,
// 				},
// 			},
// 		})),
// 	})
// }

// SendRedpacketRecover 发送红包回收消息
func (c *Context) SendRedpacketRecover(msg MsgRedpacketRecover) error {
	return c.SendTradeSystemNotifyTemplate(MsgTradeSystemNotifyTemplate{
		ChannelID:      msg.ChannelID,
		ChannelType:    msg.ChannelType,
		LeftTitle:      "红包退款到账通知",
		CenterTitle:    fmt.Sprintf("¥%0.2f", util.CentToYuan(msg.Amount)),
		CenterSubtitle: "退款金额",
		URLTitle:       "查看详情",
		Notice:         "红包退款",
		Attrs: map[string]string{
			"退款方式":         "退回零钱",
			"退款原因":         "红包超过24小时未被领取",
			"到账时间":         util.ToyyyyMMddHHmmss(time.Unix(msg.ExpiredAt, 0)),
			"备注":           "退款金额已到账",
			"creater_name": msg.CreaterName,
			"creater":      msg.Creater,
		},
	})
}

// MsgRedpacketReceive 红包领取
type MsgRedpacketReceive struct {
	ChannelID    string `json:"channel_id"`
	ChannelType  uint8  `json:"channel_type"`
	RecordNo     string `json:"record_no"`     // 转账记录编号
	Creater      string `json:"creater"`       // 红包创建者
	CreaterName  string `json:"create_name"`   // 红包创建者名称
	Receiver     string `json:"receiver"`      // 领取者uid
	ReceiverName string `json:"receiver_name"` // 领取者名称uid
}

// MsgRedpacketRecover 红包退款
type MsgRedpacketRecover struct {
	Creater     string `json:"creater"`      // 红包创建者
	CreaterName string `json:"create_name"`  // 红包创建者名称
	Amount      int64  `json:"amount"`       // 退款金额
	ExpiredAt   int64  `json:"expired_at"`   // 过期时间
	ChannelID   string `json:"channel_id"`   // 频道ID
	ChannelType uint8  `json:"channel_type"` // 频道类型
}

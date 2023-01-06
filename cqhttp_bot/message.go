package cqhttp_bot

import (
	"fmt"
	"strings"
)

func ParseMessage(raw string) (m Messages) {
	var index int
	fmt.Println(raw)
	for {
		index = strings.Index(raw, "[CQ:")
		if index > 0 {
			m = append(m, TextMessage(raw[:index]))
			raw = raw[index:]
		} else if index == 0 {
			feet := strings.Index(raw, "]")
			if feet == -1 {
				m = append(m, TextMessage(raw))
				break
			}
			data := parse(raw[:feet])
			cq := strings.Split(raw[:feet], ",")
			switch cq[0] {
			case "[CQ:image": // 图片
				m = append(m, ImageMessage(data["file"], data["url"], data["type"] == "flash"))
			case "[CQ:face": // 表情
			case "[CQ:reply": // 回复
				m = append(m, ReplyMessage(data["id"]))
			case "[CQ:at": // @
				m = append(m, AtMessage(data["qq"]))
			}
			raw = raw[feet+1:]
		} else if index == -1 {
			if raw != "" {
				m = append(m, TextMessage(raw))
			}
			break
		}
	}
	return
}

func EncodingMessage(msgs Messages) string {
	var out strings.Builder
	for _, m := range msgs {
		switch m.Type {
		case Reply:
			out.WriteString(fmt.Sprintf("[CQ:reply,id=%s]", m.Text))
		case IMAGE:
			if m.Image == nil {
				continue
			}
			var imageType = "show"
			if m.Image.Flash {
				imageType = "flash"
			}
			var file = m.Image.File
			if file == "" && m.Image.Url != "" {
				file = m.Image.Url
			}
			out.WriteString(fmt.Sprintf("[CQ:image,file=%s,type=%s]", file, imageType))
		case TEXT:
			out.WriteString(m.Text)
		case At:
			if m.At == nil {
				continue
			}
			out.WriteString(fmt.Sprintf("[CQ:at,qq=%v", m.At.Qid))
			if m.At.Name != "" {
				out.WriteString(fmt.Sprintf(",name=%s", m.At.Name))
			}
			out.WriteString("]")
		}
	}
	return out.String()
}
func MessageArray(msg ...Message) Messages {
	fmt.Println(msg)
	return msg
}
func TextMessage(text string) Message {
	return Message{
		Type: TEXT,
		Text: text,
	}
}

func AtMessage(qid any, name ...string) Message {
	m := Message{
		Type: At,
		At: &messageAt{
			Qid: qid,
		},
	}
	if len(name) != 0 && name[0] != "" {
		m.At.Name = name[0]
	}
	return m
}

func ImageMessage(file, url string, flash bool) Message {
	return Message{
		Type: IMAGE,
		Image: &messageImage{
			Url:   url,
			File:  file,
			Flash: flash,
		},
	}
}
func ReplyMessage(messageId string) Message {
	return Message{
		Type: Reply,
		Text: messageId,
	}
}

func (m *Messages) Add(msg Message) Messages {
	*m = append(*m, msg)
	return *m
}
func (m *Messages) AddTextMessage(text string) Messages {
	return m.Add(TextMessage(text))
}
func (m *Messages) AtMessage(qid any, name ...string) Messages {
	return m.Add(AtMessage(qid, name...))
}

func (m *Messages) ImageMessage(file, url string, flash bool) Messages {
	return m.Add(ImageMessage(file, url, flash))
}
func (m *Messages) ReplyMessage(messageId string) Messages {
	return m.Add(ReplyMessage(messageId))
}

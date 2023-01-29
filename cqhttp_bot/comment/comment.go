package comment

import (
	"github.com/Rehtt/qbot/cqhttp_bot"
	"strings"
)

type Cmd struct {
	bot      *cqhttp_bot.Bot
	comments map[string]*Comment
}
type Comment struct {
	Name string
	Help string
	Run  commentRunFunc
	sub  map[string]*Comment
}

type commentRunFunc func(thisComment *Comment, paramete string, bot *cqhttp_bot.Bot, messageType cqhttp_bot.EventMessageType, messageId int32, senderQid, groupId int64, message *cqhttp_bot.EventMessage)

func New(bot *cqhttp_bot.Bot) (c *Cmd) {
	c = new(Cmd)
	c.bot = bot
	c.comments = DefaultComment(c)
	return
}
func (c *Cmd) Run() {
	c.bot.OnPrivateMessage(func(messageId int32, userId int64, message *cqhttp_bot.EventMessage) {
		c.parseMessage(cqhttp_bot.Private, messageId, userId, 0, message)
	})
	c.bot.OnGroupMessage(func(messageId int32, senderQid, groupId int64, message *cqhttp_bot.EventMessage) {
		c.parseMessage(cqhttp_bot.Group, messageId, senderQid, groupId, message)
	})

}
func (c *Cmd) parseMessage(messageType cqhttp_bot.EventMessageType, messageId int32, senderQid, groupId int64, message *cqhttp_bot.EventMessage) {
	for _, m := range message.Messages {
		if m.Type == cqhttp_bot.TEXT && m.Text[0] == '/' {
			com, p := parseComment(m.Text[1:], c.comments)
			if com != nil {
				com.Run(com, p, c.bot, messageType, messageId, senderQid, groupId, message)
			}
		}
	}
}
func parseComment(text string, comments map[string]*Comment) (c *Comment, p string) {
	texts := strings.SplitN(text, " ", 2)
	if com, ok := comments[texts[0]]; !ok {
		return nil, text
	} else {
		c = com
		if len(com.sub) != 0 && len(texts) == 2 && texts[1] != "" {
			c, p = parseComment(texts[1], com.sub)
			if c == nil {
				c = com
			}
		}
		return
	}
}
func (c *Cmd) AddComment(comment ...*Comment) {
	if len(comment) == 0 {
		return
	}
	if c.comments == nil {
		c.comments = make(map[string]*Comment, len(comment))
	}
	for _, sub := range comment {
		c.comments[sub.Name] = sub
	}
}
func (c *Comment) AddSubComment(comment ...*Comment) {
	if len(comment) == 0 {
		return
	}
	if c.sub == nil {
		c.sub = make(map[string]*Comment, len(comment))
	}
	for _, sub := range comment {
		c.sub[sub.Name] = sub
	}
}
func (c *Cmd) help() string {
	var tmp strings.Builder
	for _, com := range c.comments {
		if tmp.Len() != 0 {
			tmp.WriteString("\n")
		}
		tmp.WriteString(com.Name)
		tmp.WriteString(" ")
		tmp.WriteString(com.Help)
	}
	return tmp.String()
}

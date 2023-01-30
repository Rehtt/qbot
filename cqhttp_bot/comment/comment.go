package comment

import (
	"github.com/Rehtt/qbot/cqhttp_bot"
	"strings"
)

type Cmd struct {
	bot      *cqhttp_bot.Bot
	comments comments
}
type comments map[string]*Comment
type Comment struct {
	Name  string
	Usage string
	flag  *Flag
	Run   RunFunc
	sub   comments
}

type RunFunc func(paramete string, flag Flag, bot *cqhttp_bot.Bot, messageType cqhttp_bot.EventMessageType, messageId int32, senderQid, groupId int64, message *cqhttp_bot.EventMessage)

func New(bot *cqhttp_bot.Bot) (c *Cmd) {
	c = new(Cmd)
	c.bot = bot
	c.comments = DefaultComment(c)
	c.run()
	return
}
func (c *Cmd) run() {
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
			com, arg, p := c.comments.parse(m.Text[1:])
			//com, arg, p := parseComment(m.Text[1:], c.comments, nil)
			if com != nil {
				com.Run(p, arg, c.bot, messageType, messageId, senderQid, groupId, message)
			}
		}
	}
}
func (c comments) parse(str string) (co *Comment, f Flag, p string) {
	strArr := strings.Split(str, " ")
	var flagName string
	for i := 0; i < len(strArr); i++ {
		s := strArr[i]
		if flagName != "" {
			f.set(flagName, s)
			flagName = ""
			continue
		}
		if com, ok := c[s]; ok {
			co = com
		} else {
			if co != nil && co.flag != nil && s[0] == '-' {
				flagName = s[1:]
				if _, ok = co.flag.Get(flagName); ok {
					cqhttp_bot.DeepCopy(&f.args, co.flag.args)
					continue
				}
			}
			return co, f, strings.Join(strArr[i:], " ")
		}
		flagName = ""
	}
	return
}
func (c *Cmd) AddComment(comments ...*Comment) {
	if len(comments) == 0 {
		return
	}
	if c.comments == nil {
		c.comments = make(map[string]*Comment, len(comments))
	}
	for _, sub := range comments {
		c.comments[sub.Name] = sub
	}
}
func (c *Comment) AddSubComment(comments ...*Comment) {
	if len(comments) == 0 {
		return
	}
	if c.sub == nil {
		c.sub = make(map[string]*Comment, len(comments))
	}
	for _, sub := range comments {
		c.sub[sub.Name] = sub
	}
}
func (c *Comment) Flag() *Flag {
	if c.flag == nil {
		c.flag = new(Flag)
	}
	return c.flag
}

func (c *Cmd) help() string {
	var tmp strings.Builder
	for _, com := range c.comments {
		if tmp.Len() != 0 {
			tmp.WriteString("\n")
		}
		tmp.WriteString(com.Name)
		tmp.WriteString(" ")
		tmp.WriteString(com.Usage)
	}
	return tmp.String()
}

package command

import (
	"fmt"
	"github.com/Rehtt/qbot/cqhttp_bot"
	"strings"
)

type Cmd struct {
	bot      *cqhttp_bot.Bot
	Commands Commands
}
type Commands map[string]*Command
type Command struct {
	Name  string
	Usage string
	flag  *Flag
	Run   RunFunc
	sub   Commands
}

type RunFunc func(paramete string, flag Flag, bot *cqhttp_bot.Bot, messageType cqhttp_bot.EventMessageType, messageId int32, senderQid, groupId int64, message *cqhttp_bot.EventMessage)

func New(bot *cqhttp_bot.Bot) (c *Cmd) {
	c = new(Cmd)
	c.bot = bot
	c.Commands = DefaultCommand(c)
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
			com, arg, p := c.Commands.Parse(m.Text[1:])
			//com, arg, p := parseCommand(m.Text[1:], c.Commands, nil)
			if com != nil {
				com.Run(p, arg, c.bot, messageType, messageId, senderQid, groupId, message)
			}
		}
	}
}
func (c Commands) Parse(str string) (co *Command, f Flag, p string) {
	strArr := strings.Split(str, " ")
	var flagName string
	for i := 0; i < len(strArr); i++ {
		s := strArr[i]
		if s == "" {
			continue
		}
		if flagName != "" {
			f.set(flagName, s)
			flagName = ""
			continue
		}
		if co == nil {
			if com, ok := c[s]; ok {
				co = com
				flagName = ""
				continue
			}
		} else {
			if com, ok := co.sub[s]; ok {
				co = com
				flagName = ""
				continue
			}
			if co.flag != nil && s[0] == '-' {
				flagName = s[1:]
				if _, ok := co.flag.Get(flagName); ok {
					cqhttp_bot.DeepCopy(&f.args, co.flag.args)
					continue
				}
			}
			return co, f, strings.Join(strArr[i:], " ")
		}
	}
	return
}
func (c *Cmd) AddCommand(Commands ...*Command) {
	if len(Commands) == 0 {
		return
	}
	if c.Commands == nil {
		c.Commands = make(map[string]*Command, len(Commands))
	}
	for _, sub := range Commands {
		c.Commands[sub.Name] = sub
	}
}
func (c *Command) AddSubCommand(Commands ...*Command) {
	if len(Commands) == 0 {
		return
	}
	if c.sub == nil {
		c.sub = make(map[string]*Command, len(Commands))
	}
	for _, sub := range Commands {
		c.sub[sub.Name] = sub
	}
}
func (c *Command) Flag() *Flag {
	if c.flag == nil {
		c.flag = new(Flag)
	}
	return c.flag
}
func (c *Cmd) GetCommand(name string) *Command {
	return c.Commands[name]
}

func (c *Cmd) Help() string {
	var tmp strings.Builder
	for _, com := range c.Commands {
		if tmp.Len() != 0 {
			tmp.WriteString("\n")
		}
		tmp.WriteString(fmt.Sprintf("/%s %s", com.Name, com.Usage))
	}
	return tmp.String()
}
func (c *Command) Help() string {
	var tmp strings.Builder
	//tmp.WriteString(fmt.Sprintf("%s %s\n", c.Name, c.Usage))
	if len(c.sub) != 0 {
		tmp.WriteString("子命令：\n")
		for _, sub := range c.sub {
			tmp.WriteString(fmt.Sprintf(" %s %s\n", sub.Name, sub.Usage))
		}
	}
	if c.flag != nil && len(c.flag.args) != 0 {
		tmp.WriteString("选项：\n")
		for _, f := range c.flag.args {
			tmp.WriteString(fmt.Sprintf(" -%s %s [默认：%s]\n", f.name, f.usage, f.value))
		}
	}
	return tmp.String()
}

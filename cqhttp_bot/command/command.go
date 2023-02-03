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

type RunFunc func(paramete string, flag Flag, bot *cqhttp_bot.Bot, ctx *cqhttp_bot.EventMessageContext)

func New(bot *cqhttp_bot.Bot) (c *Cmd) {
	c = new(Cmd)
	c.bot = bot
	c.run()
	return
}
func (c *Cmd) run() {
	c.bot.OnMessage(func(ctx *cqhttp_bot.EventMessageContext) {
		c.parseMessage(ctx)
	})
}
func (c *Cmd) parseMessage(ctx *cqhttp_bot.EventMessageContext) {
	for _, m := range ctx.Message.Messages {
		if m.Type == cqhttp_bot.TEXT && m.Text[0] == '/' {
			com, arg, p := c.Commands.Parse(m.Text[1:])
			//com, arg, p := parseCommand(m.Text[1:], c.Commands, nil)
			if com != nil {
				com.Run(p, arg, c.bot, ctx)
			}
		}
	}
}
func (c Commands) Parse(str string) (co *Command, f Flag, p string) {
	fields := strings.Fields(str)
	var flagName string
	defer func() {
		if co != nil && co.flag != nil {
			co.flag.Range(func(name, value, usage string) {
				if _, ok := f.Get(name); !ok {
					f.set(name, value)
				}
			})
		}
	}()
	for i := 0; i < len(fields); i++ {
		s := fields[i]
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
			if s[0] == '-' {
				flagName = s[1:]
				if co.flag != nil {
					if v, ok := co.flag.Get(flagName); ok {
						f.set(flagName, v)
					}
				}
				continue
			}
			return co, f, strings.Join(fields[i:], " ")
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

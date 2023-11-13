/**
 * Copyright (c) 2023 Rehtt
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package command

import (
	"fmt"
	"strings"

	"github.com/Rehtt/qbot/cqhttp_bot"
)

type Cmd struct {
	bot      *cqhttp_bot.Bot
	trigger  string
	Commands Commands
	SelfQid  int64
}
type (
	Commands map[string]*Command
	Command  struct {
		Name  string
		Usage string
		flag  *Flag
		Run   RunFunc
		sub   Commands
	}
)

type RunFunc func(paramete string, flag Flag, bot *cqhttp_bot.Bot, ctx *cqhttp_bot.EventMessageContext)

func New(bot *cqhttp_bot.Bot) (c *Cmd) {
	c = new(Cmd)
	c.bot = bot
	c.trigger = "/"
	c.run()
	return
}

// Trigger 触发解析命令的字段，默认：/
func (c *Cmd) Trigger(t string) {
	c.trigger = t
}

func (c *Cmd) run() {
	c.bot.OnMessage(func(ctx *cqhttp_bot.EventMessageContext) {
		if c.SelfQid == 0 {
			info, err := c.bot.Action.GetSelfInfo()
			if err != nil {
				c.bot.Options.Log().ErrorContext("command GetSelfInfo", "error", err)
				return
			}
			c.SelfQid = info.UserId
		}
		if ctx.SenderId == c.SelfQid {
			return
		}
		c.parseMessage(ctx)
	})
}

func (c *Cmd) parseMessage(ctx *cqhttp_bot.EventMessageContext) {
	for _, m := range ctx.Message.Messages {
		if m.Type == cqhttp_bot.TEXT && strings.HasPrefix(m.Text, c.trigger) {
			com, arg, p := c.Commands.Parse(m.Text[len(c.trigger):])
			// com, arg, p := parseCommand(m.Text[1:], c.Commands, nil)
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
	var flagValueTmp strings.Builder
	for i := 0; i < len(fields); i++ {
		s := fields[i]
		if flagName != "" {
			if l := len(s); l > 1 && s[0] == '"' && flagValueTmp.Len() == 0 {
				flagValueTmp.WriteString(s[1:])
				flagValueTmp.WriteString(" ")
				continue
			} else if s[l-1] == '"' && flagValueTmp.Len() != 0 {
				flagValueTmp.WriteString(s[:l-1])
				s = flagValueTmp.String()
				flagValueTmp.Reset()
			} else if flagValueTmp.Len() != 0 {
				flagValueTmp.WriteString(s)
				flagValueTmp.WriteString(" ")
				continue
			}

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
		tmp.WriteString(fmt.Sprintf("%s%s %s", c.trigger, com.Name, com.Usage))
	}
	return tmp.String()
}

func (c *Command) Help() string {
	var tmp strings.Builder
	// tmp.WriteString(fmt.Sprintf("%s %s\n", c.Name, c.Usage))
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

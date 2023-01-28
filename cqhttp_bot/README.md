对[go-cqhttp](https://github.com/Mrs4s/go-cqhttp) api简单封装

[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/Rehtt/qbot/cqhttp_bot)

## Installation

```bash
go get -u github.com/Rehtt/qbot/cqhttp_bot
```

## Use

### 初始化
```go
bot := cqhttp_bot.New("ws://127.0.0.1:8080")
```

### 事件
```go
// 群消息事件
bot.Event.OnGroupMessage(func (senderQid, groupId int64, message *EventMessage) {
    fmt.Println(senderQid, groupId, message.Messages, message.RawMessage)
})

// 私人消息事件
bot.OnPrivateMessage(func (userId int64, message *EventMessage) {
    fmt.Println(userId, message.Messages)
})

```

### 请求
```go
// 获取好友列表
bot.GetFriendsList()

// 发送消息
bot.SendMsg(<目标qq>,<消息>,<EventMessageType>)
```

对[go-cqhttp](https://github.com/Mrs4s/go-cqhttp) api简单封装

> 由于[GO-CQHTTP#2471](https://github.com/Mrs4s/go-cqhttp/issues/2471)原因，决定逐渐放缓cqhttp_bot开发

[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/Rehtt/qbot/cqhttp_bot)

## Installation

```bash
go get -u github.com/Rehtt/qbot/cqhttp_bot
```

## Use

### 初始化

```go
bot := cqhttp_bot.New("ws://127.0.0.1:8080")
//bot.Start() //  已非阻塞的方式运行
bot.Run()   //  已阻塞的方式运行
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

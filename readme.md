# mwchat
微信接入chatgpt

/chat 前缀的消息会自动透传到chatgpt
上下文可用
私信、群聊可用

运行：`go run . --key="your open ai key" --proxy="your proxy address"`<br>
eg: `go run . --key="xxxxxx123xxxsds" --proxy="127.0.0.1:7890"`<br>
不使用代理：`go run . --key="your open ai key"`<br>
# mwchat
微信接入chatgpt

/chat 前缀的消息会自动透传到chatgpt</br>
上下文可用</br>
私信、群聊可用</br>


```
# 获取项目
git clone git@github.com:heiwa9/mwchat.git

# 进入项目目录
cd mwchat

# 启动项目
运行：`go run . --key="your open ai key" --proxy="your proxy address"`<br>
eg: `go run . --key="xxxxxx123xxxsds" --proxy="127.0.0.1:7890"`<br>
不使用代理：`go run . --key="your open ai key"`<br>
```
package main

import (
	"sync"
	"time"

	"github.com/DanPlayer/chatgpt-sdk/v1"
)

type MsgCtx struct {
    HistoryMsgs []v1.ChatMessage
    LastTime    time.Time
    lock        bool
}

func (mc *MsgCtx) Add(cm v1.ChatMessage){
    mc.HistoryMsgs = append(mc.HistoryMsgs,cm)
    if len(mc.HistoryMsgs) >= 6 {
        mc.HistoryMsgs = mc.HistoryMsgs[len(mc.HistoryMsgs) - 5:]
    }
}

func (mc *MsgCtx) Clear() {
//    mc.HistoryMsgs = mc.HistoryMsgs[:0]  //保留空间
    mc.HistoryMsgs = []v1.ChatMessage{}  //不保留空间
}

func (mc *MsgCtx) Lock() {
    mc.lock = true
}

func (mc *MsgCtx) UnLock(){
    mc.lock = false
}

func (mc *MsgCtx) IsLock() bool {
    return mc.lock
}

type ChatMsgCtx struct {
	chatMsgCtx     map[string]*MsgCtx
	chatMsgCtxLock sync.Mutex
}

func NewChatMsgCtx() *ChatMsgCtx {
	return &ChatMsgCtx{
        chatMsgCtx: map[string]*MsgCtx{},
        chatMsgCtxLock: sync.Mutex{},
    }
}

func (cmc *ChatMsgCtx) FindCtx(uuid string) *MsgCtx {
	cmc.chatMsgCtxLock.Lock()
	if mc, ok := cmc.chatMsgCtx[uuid]; ok {
        mc.LastTime = time.Now()
		return mc
	}
	ctx := &MsgCtx{
        HistoryMsgs: []v1.ChatMessage{},
        LastTime:    time.Now(),
        lock:        false,
    }
	cmc.chatMsgCtx[uuid] = ctx
	return ctx
}

func (cmc *ChatMsgCtx) Clear(uuid string) {
	cmc.chatMsgCtxLock.Lock()
	if v,ok := cmc.chatMsgCtx[uuid];ok{
        v.Clear()
    }
	cmc.chatMsgCtxLock.Unlock()
}

func (cmc *ChatMsgCtx) ClearCtxTask() {
	go func() {
        ticker := time.NewTicker(time.Minute)
		for range ticker.C {
            cmc.chatMsgCtxLock.Lock()
            for _,v := range cmc.chatMsgCtx {
                if time.Since(v.LastTime) > time.Minute * 10 && !v.IsLock() {
                    v.Clear()
                }
            }
            cmc.chatMsgCtxLock.Unlock()
        }
	}()
}

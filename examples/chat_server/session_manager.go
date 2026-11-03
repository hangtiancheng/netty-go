package main

import (
	"sync"

	"github.com/yukino-js/netty-go"
)

type Manager interface {
	netty.ActiveHandler
	netty.InactiveHandler
	Size() int
	Context(id int64) netty.HandlerContext
	ForEach(func(netty.HandlerContext) bool)
	Broadcast(message netty.Message)
	BroadcastIf(message netty.Message, fn func(netty.HandlerContext) bool)
}

func NewManager() Manager {
	return &sessionManager{
		_sessions: make(map[int64]netty.HandlerContext, 64),
	}
}

type sessionManager struct {
	_sessions map[int64]netty.HandlerContext
	_mutex    sync.RWMutex
}

func (s *sessionManager) Size() int {
	s._mutex.RLock()
	size := len(s._sessions)
	s._mutex.RUnlock()
	return size
}

func (s *sessionManager) Context(id int64) netty.HandlerContext {
	s._mutex.RLock()
	ctx, _ := s._sessions[id]
	s._mutex.RUnlock()
	return ctx
}

func (s *sessionManager) ForEach(fn func(netty.HandlerContext) bool) {
	s._mutex.RLock()
	defer s._mutex.RUnlock()

	for _, ctx := range s._sessions {
		fn(ctx)
	}
}

func (s *sessionManager) Broadcast(message netty.Message) {
	s.ForEach(func(ctx netty.HandlerContext) bool {
		ctx.Write(message)
		return true
	})
}

func (s *sessionManager) BroadcastIf(message netty.Message, fn func(netty.HandlerContext) bool) {
	s.ForEach(func(ctx netty.HandlerContext) bool {
		if fn(ctx) {
			ctx.Write(message)
		}
		return true
	})
}

func (s *sessionManager) HandleActive(ctx netty.ActiveContext) {

	s._mutex.Lock()
	s._sessions[ctx.Channel().ID()] = ctx
	s._mutex.Unlock()

	ctx.HandleActive()
}

func (s *sessionManager) HandleInactive(ctx netty.InactiveContext, ex netty.Exception) {
	s._mutex.Lock()
	delete(s._sessions, ctx.Channel().ID())
	s._mutex.Unlock()

	ctx.HandleInactive(ex)
}

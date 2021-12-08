package sHelper

import (
	"gitlab.com/soteapps/packages/v2021/sError"
	"gitlab.com/soteapps/packages/v2021/sLogger"
)

var (
	intitalized = false
)

type Helper struct {
	Env              Environment
	r                *Run
	CreateSubscriber func(consumerName, subject string, streamName ...string) *Subscriber
	CreateDatabase   func() sError.SoteError
	InitApp          func() sError.SoteError
	AddSubscriber    func(consumerName, subject string, listener MessageListener, schema *Schema, streamName ...string) sError.SoteError
	Run              func(isGoroutine bool)
}

func NewHelper(env Environment) *Helper {
	var (
		h Helper
	)
	h = Helper{
		Env:              env,
		r:                NewRun(env),
		CreateSubscriber: h.createSubscriber,
		CreateDatabase:   h.createDatabase,
		InitApp:          h.initApp,
		AddSubscriber:    h.addSubscriber,
		Run:              h.run,
	}
	return &h
}

func (h *Helper) addSubscriber(consumerName, subject string, listener MessageListener, schema *Schema, streamName ...string) (soteErr sError.SoteError) {
	sLogger.DebugMethod()
	s := h.CreateSubscriber(consumerName, subject, streamName...)
	if schema != nil {
		s.Schema = schema
		soteErr = schema.Validate()
	}
	if soteErr.ErrCode == nil && !intitalized {
		if soteErr = h.InitApp(); soteErr.ErrCode == nil {
			soteErr = h.CreateDatabase()
		}
		intitalized = true
	}
	if soteErr.ErrCode == nil {
		if soteErr = s.PullSubscribe(); soteErr.ErrCode == nil {
			h.r.AddSubscriber(s, listener)
		}
	}
	return
}

func (h *Helper) run(isGoroutine bool) {
	sLogger.DebugMethod()
	h.r.Listen(func(s *Subscriber) (soteErr sError.SoteError) {
		var (
			messages []Msg
		)
		if messages, soteErr = s.Fetch(); soteErr.ErrCode == nil {
			for _, message := range messages {
				sLogger.DebugMethod()
				s.Start(&message)
				if isGoroutine {
					go func(s *Subscriber, msg Msg) {
						soteErr := s.Listener(s, &msg)
						s.End(&msg, soteErr)
					}(s, message)
				} else {
					soteErr := s.Listener(s, &message)
					s.End(&message, soteErr)
				}
			}
		}
		return
	})
}

func (h *Helper) createSubscriber(consumerName, subject string, streamName ...string) *Subscriber {
	sLogger.DebugMethod()
	return NewSubscriber(h.r, consumerName, subject, streamName...)
}

func (h *Helper) createDatabase() sError.SoteError {
	sLogger.DebugMethod()
	return NewDatabase(h.r)
}

func (h *Helper) initApp() sError.SoteError {
	sLogger.DebugMethod()
	return h.r.InitApp()
}

package sHelper

import (
	"fmt"
	"time"

	"github.com/nats-io/nats.go"

	"gitlab.com/soteapps/packages/v2021/sConfigParams"
	"gitlab.com/soteapps/packages/v2021/sDatabase"
	"gitlab.com/soteapps/packages/v2021/sError"
	"gitlab.com/soteapps/packages/v2021/sLogger"
	"gitlab.com/soteapps/packages/v2021/sMessage"
)

type ConsumerInfo = nats.ConsumerInfo

type MessageListener func(*Subscriber, *Msg) sError.SoteError

type Run struct {
	Env                 Environment
	Nats                *natsConfig
	Subscribers         []*Subscriber
	ValidateEnvironment func(environment string) sError.SoteError
	GetNATSURL          func(application, environment string) (string, sError.SoteError)
	NewMessage          func(env Environment, natsURL string) (*sMessage.MessageManager, sError.SoteError)
	GetConnection       func(dbName, user, password, host, sslMode string, port, timeout int) (sDatabase.ConnInfo, sError.SoteError)
	Listen              func(listener func(*Subscriber) sError.SoteError)
	myMMPtr             *sMessage.MessageManager
	dbHelper            *DatabaseHelper
	returnChain         chan *ReturnChain
}

type natsConfig struct {
	Secure             bool
	MaxReconnect       int
	ReconnectWait      time.Duration
	ConnectionName     string
	CredentialFileName string
}

type Msg struct {
	Subject string
	Header  nats.Header
	Data    []byte
	index   int
	uuid    string
}

type ReturnChain struct {
	msg     *Msg
	soteErr sError.SoteError
	s       *Subscriber
}

func (m *Msg) Index() int {
	return m.index
}

func (m *Msg) Id() string {
	return m.uuid
}

func NewRun(env Environment) *Run {
	sLogger.DebugMethod()
	var (
		run Run
	)
	run = Run{
		Env:                 env,
		Subscribers:         []*Subscriber{},
		GetConnection:       sDatabase.GetConnection,
		ValidateEnvironment: sConfigParams.ValidateEnvironment,
		GetNATSURL:          sConfigParams.GetNATSURL,
		NewMessage:          run.newMessage,
		Listen:              run.listen,
		Nats: &natsConfig{
			Secure:             true,
			MaxReconnect:       5,
			ReconnectWait:      250 * time.Millisecond,
			ConnectionName:     "myNATS",
			CredentialFileName: "",
		},
	}
	return &run
}

func (r *Run) InitApp() (soteErr sError.SoteError) {
	sLogger.DebugMethod()
	var (
		natsURL string
	)
	env := r.Env
	if soteErr = r.ValidateEnvironment(env.TargetEnvironment); soteErr.ErrCode == nil {
		if natsURL, soteErr = r.GetNATSURL(env.ApplicationName, env.TargetEnvironment); soteErr.ErrCode == nil {
			r.myMMPtr, soteErr = r.NewMessage(env, natsURL)
		}
	}
	return
}

func (r *Run) newMessage(env Environment, natsURL string) (*sMessage.MessageManager, sError.SoteError) {
	sLogger.DebugMethod()
	return sMessage.New(env.ApplicationName, env.TargetEnvironment, r.Nats.CredentialFileName, natsURL,
		r.Nats.ConnectionName, r.Nats.Secure, r.Nats.MaxReconnect, r.Nats.ReconnectWait, env.TestMode)
}

func (r *Run) AddSubscriber(s *Subscriber, listener MessageListener) {
	sLogger.DebugMethod()
	s.Listener = listener
	r.Subscribers = append(r.Subscribers, s)
}

func (r *Run) listen(listener func(*Subscriber) (soteErr sError.SoteError)) {
	sLogger.DebugMethod()
	if len(r.Subscribers) > 0 {
		sLogger.Info(fmt.Sprintf("Listening Subscribers: %v", len(r.Subscribers)))
		r.returnChain = make(chan *ReturnChain, 1)
		go func() {
			//Listen error(s) from goroutine
			for rc := range r.returnChain {
				if rc.soteErr.ErrCode != nil {
					r.Error(rc.soteErr, rc.msg)
				}
				sLogger.Info(fmt.Sprintf("End Subscription[%v] Subject: %s, Index: %v", rc.msg.Id(), rc.msg.Subject, rc.msg.Index()))
			}
		}()
		for {
			for _, s := range r.Subscribers {
				soteErr := listener(s)
				if soteErr.ErrCode != nil {
					r.PanicService(soteErr)
				}
			}
			time.Sleep(250 * time.Millisecond)
		}
	}
}

func (r *Run) Error(soteErr sError.SoteError, msg *Msg) {
	sLogger.DebugMethod()
	sLogger.Info(soteErr.FmtErrMsg)
	sLogger.Info(fmt.Sprintf("MESSAGE - Subject: %v, Data: %v", msg.Subject, string(msg.Data)))
	if r.Env.TestMode {
		r.PanicService(soteErr)
	}
}

func (r *Run) PanicService(soteErr sError.SoteError) {
	sLogger.DebugMethod()
	panic(soteErr.FmtErrMsg)
}

### Removing unused dependencies
go mod tidy

### Run
go run main.go
or
go run main.go --targetEnv production

### is like all but adds stack frames for run-time functions and shows goroutines created internally by the run-time.
GOTRACEBACK=system

### Run tests
go test -v ./sHelper/...

### Run Code cover
go test ./sHelper/... -cover

go test ./sHelper/ -coverprofile=coverage.out && go tool cover -func=coverage.out


### main.go
```
package main

import (
	"fmt"
	"notification/packages"
	"notification/sHelper"
)

func main() {
	env := sHelper.Parameter{
		Version:     "v2021.1.0",
		AppName:     "notification",
		Description: `This business service adds notification details to the database and sends an email as well.`,
	}.Init()

	if soteErr := packages.Run(env); soteErr.ErrCode != nil {
		panic(soteErr.FmtErrMsg)
	}
}
```

### run.go  - sHelper.Helper
```
package packages

import (
	"notification/sHelper"

	"gitlab.com/soteapps/packages/v2021/sError"
	"gitlab.com/soteapps/packages/v2021/sLogger"
)

func Run(env sHelper.Environment) (soteErr sError.SoteError) {
	sLogger.DebugMethod()
	helper := sHelper.NewHelper(env)
	if soteErr = helper.AddSubscriber("bsl-notification-wildcard", "bsl.notification.add", getNotification, nil); soteErr.ErrCode == nil {
		helper.Run(false) //synchronous
	}
	return
}

func getNotification(s *sHelper.Subscriber, message *sHelper.Msg) (soteErr sError.SoteError) {
	sLogger.DebugMethod()
	return
}
```

### run.go  - sHelper.Helper and schema validator
```
package packages

import (
	"notification/sHelper"

	"gitlab.com/soteapps/packages/v2021/sError"
	"gitlab.com/soteapps/packages/v2021/sLogger"
)

type Notification struct {
	sHelper.BaseSchema
	Subject     string `json:"subject"`
	TextPayload string `json:"text-payload"`
	HtmlPayload string `json:"html-payload"`
	CommChannel string `json:"comm-channel"`
}

func Run(env sHelper.Environment) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	schema := sHelper.Schema{
		FileName:  "packages/resources/schema-v1.json",
		StructRef: &Notification{},
	}
	helper := sHelper.NewHelper(env)
	if soteErr = helper.AddSubscriber("bsl-notification-wildcard", "bsl.notification.add", getNotification, &schema); soteErr.ErrCode == nil {
		helper.Run(true) //asynchronously using goroutine
	}
	return
}

func getNotification(s *sHelper.Subscriber, message *sHelper.Msg) (soteErr sError.SoteError) {
	sLogger.DebugMethod()
	body := Notification{}
	soteErr = s.Schema.Parse(message.Data, &body)
	sLogger.Info(body.JsonWebToken)
	return
}
```

### run.go  - sHelper.Run
```
package packages

import (
	"notification/sHelper"

	"gitlab.com/soteapps/packages/v2021/sError"
	"gitlab.com/soteapps/packages/v2021/sLogger"
)

func Run(env sHelper.Environment) (soteErr sError.SoteError) {
	sLogger.DebugMethod()
	run := sHelper.NewRun(env)
	if soteErr = run.InitApp(); soteErr.ErrCode == nil {
		if soteErr = sHelper.NewDatabase(run); soteErr.ErrCode == nil {
			s := sHelper.NewSubscriber(run, "bsl-notification-wildcard", "bsl.notification.add")
			if soteErr = s.PullSubscribe(); soteErr.ErrCode == nil {
				run.AddSubscriber(s, getNotification)
			}
			run.Listen(listenNotifications)
		}
	}
	return
}

func listenNotifications(s *sHelper.Subscriber) (soteErr sError.SoteError) {
	var (
		messages []sHelper.Msg
	)
	if messages, soteErr = s.Fetch(); soteErr.ErrCode == nil {
		for _, message := range messages {
			s.Listener(s, &message)
		}
	}
	return
}

func getNotification(s *sHelper.Subscriber, msg *sHelper.Msg) (soteErr sError.SoteError) {
	sLogger.DebugMethod()
	s.Start(msg)
	s.End(msg, soteErr)
	return sError.SoteError{}
}
```

### run_test.go  - unittest template
```
package packages

import (
	"notification/sHelper"
	"gitlab.com/soteapps/packages/v2021/sError"
	"testing"
	"bou.ke/monkey"
)

var (
	AssertEqual        = sHelper.AssertEqual
	verifyConsumerName = "bsl-notification-wildcard"
	verifySubject      = "bsl.notification.add"
)

func TestRun(t *testing.T) {
	monkey.Patch(sHelper.NewHelper, func(env sHelper.Environment) *sHelper.Helper {
		helper := sHelper.Helper{
			Env: env,
		}
		helper.Run = func(isGoroutine bool) {}
		helper.AddSubscriber = func(consumerName, subject string, _ sHelper.MessageListener, _ *sHelper.Schema) sError.SoteError {
			AssertEqual(t, consumerName, verifyConsumerName)
			AssertEqual(t, subject, verifySubject)
			return sError.SoteError{}
		}
		return &helper
	})
	env, _ := sHelper.NewEnvironment(sHelper.ENVDEFAULTAPPNAME, sHelper.ENVDEFAULTTARGET, sHelper.ENVDEFAULTTARGET)
	Run(env)
}
```
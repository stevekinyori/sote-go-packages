package main

import "gitlab.com/soteapps/packages/v2021/sHelper"

func main() {
	mail := sHelper.NewMail("staging", "Test", "David Gofman")
	mail.To("dgofman@gmail.com").Name("Personal Email")
	mail.Bcc("david.gofman@getsote.com").Name("Company Email")
	soteErr := mail.Attachment("C:\\git_sote\\flutter-core\\android\\app\\src\\main\\res\\drawable\\app_icon.png")
	if soteErr.ErrCode != nil {
		panic(soteErr.FmtErrMsg)
	}
	soteErr = mail.Send("Hello World1", "<h1>Hello World2<h1>")
	if soteErr.ErrCode != nil {
		panic(soteErr.FmtErrMsg)
	}
}

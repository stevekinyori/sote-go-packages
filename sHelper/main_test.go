package sHelper

import "testing"

func TestMail(t *testing.T) {
	mail := NewMail("staging", "Test Subject", "David Gofman")
	mail.To("dgofman@gmail.com").Name("Personal Email")
	mail.Bcc("david.gofman@getsote.com").Name("Company Email")
	soteErr := mail.Attachment("./schema_test.json")
	AssertEqual(t, soteErr.FmtErrMsg, "")
	//soteErr = mail.Send("Hello World1", "<h1>Hello World2<h1>")
	AssertEqual(t, soteErr.FmtErrMsg, "")
}

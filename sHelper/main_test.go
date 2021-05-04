package sHelper

import "testing"

func TestMail(t *testing.T) {
	mail := NewMail("staging", "Test", "David Gofman")
	mail.To("dgofman@gmail.com <Personal Email>")
	//mail.Cs("david_gofman@getsote.com <Company Email>")
	mail.Send("Hello World1", "<b>Hello World2<b>")
}

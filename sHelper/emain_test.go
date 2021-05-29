package sHelper

import (
	"bytes"
	"fmt"
	"net/smtp"
	"regexp"
	"strings"
	"testing"

	"gitlab.com/soteapps/packages/v2021/sError"
)

func newTestEmail(environment, subject string, from ...string) *Email {
	email := NewEmail(environment, subject, from...)
	email.GetSmtpUsername = func(application, environment string) (string, sError.SoteError) {
		return "USERNAME", sError.SoteError{}
	}
	email.GetSmtpPassword = func(application, environment string) (string, sError.SoteError) {
		return "PASSWORD", sError.SoteError{}
	}
	return email
}

func TestEMailSubject(t *testing.T) {
	email := newTestEmail("staging", "Test Subject")
	AssertEqual(t, email.environment, "staging")
	AssertEqual(t, email.subject, "Test Subject")
	AssertEqual(t, email.from.String(), DEFAULTEMAIL)
}

func TestEMailSubjectFrom(t *testing.T) {
	email := newTestEmail("staging", "Test Subject", "David Gofman")
	AssertEqual(t, email.environment, "staging")
	AssertEqual(t, email.subject, "Test Subject")
	AssertEqual(t, email.from.String(), fmt.Sprintf("%s <%s>", "David Gofman", DEFAULTEMAIL))
}

func TestEmailTo(t *testing.T) {
	email := newTestEmail("staging", "Test Subject")

	email.To("admin@getsote.com")
	email.rcpt = []string{}
	list, soteErr := email.addRcpt("To", email.to)
	AssertEqual(t, soteErr.FmtErrMsg, "")
	AssertEqual(t, strings.Join(list, ";"), "admin@getsote.com")
	AssertEqual(t, strings.Join(email.rcpt, ";"), "admin@getsote.com")

	email.To("sales@getsote.com").Name("Customer Support")
	email.rcpt = []string{}
	list, soteErr = email.addRcpt("To", email.to)
	AssertEqual(t, soteErr.FmtErrMsg, "")
	AssertEqual(t, strings.Join(list, ";"), "admin@getsote.com;Customer Support <sales@getsote.com>")
	AssertEqual(t, strings.Join(email.rcpt, ";"), "admin@getsote.com;sales@getsote.com")

	email.To("invalid.com")
	_, soteErr = email.addRcpt("To", email.to)
	AssertEqual(t, soteErr.FmtErrMsg, "207050: invalid.com (To) is not a valid email address ERROR DETAILS: >>Key: EMAIL_PARSE Value: mail: misformatted email address")
}

func TestEmailCc(t *testing.T) {
	email := newTestEmail("staging", "Test Subject")

	email.Cc("admin@getsote.com")
	email.rcpt = []string{}
	list, soteErr := email.addRcpt("Cc", email.cc)
	AssertEqual(t, soteErr.FmtErrMsg, "")
	AssertEqual(t, strings.Join(list, ";"), "admin@getsote.com")
	AssertEqual(t, strings.Join(email.rcpt, ";"), "admin@getsote.com")

	email.Cc("sales@getsote.com").Name("Customer Support")
	email.rcpt = []string{}
	list, soteErr = email.addRcpt("Cc", email.cc)
	AssertEqual(t, soteErr.FmtErrMsg, "")
	AssertEqual(t, strings.Join(list, ";"), "admin@getsote.com;Customer Support <sales@getsote.com>")
	AssertEqual(t, strings.Join(email.rcpt, ";"), "admin@getsote.com;sales@getsote.com")
}

func TestEmailBcc(t *testing.T) {
	email := newTestEmail("staging", "Test Subject")

	email.Bcc("admin@getsote.com")
	email.rcpt = []string{}
	list, soteErr := email.addRcpt("Bcc", email.bcc)
	AssertEqual(t, soteErr.FmtErrMsg, "")
	AssertEqual(t, strings.Join(list, ";"), "admin@getsote.com")
	AssertEqual(t, strings.Join(email.rcpt, ";"), "admin@getsote.com")

	email.Bcc("sales@getsote.com").Name("Customer Support")
	email.rcpt = []string{}
	list, soteErr = email.addRcpt("Bcc", email.bcc)
	AssertEqual(t, soteErr.FmtErrMsg, "")
	AssertEqual(t, strings.Join(list, ";"), "admin@getsote.com;Customer Support <sales@getsote.com>")
	AssertEqual(t, strings.Join(email.rcpt, ";"), "admin@getsote.com;sales@getsote.com")
}

func TestEmailAttachment(t *testing.T) {
	email := newTestEmail("staging", "Test Subject")
	soteErr := email.Attachment("foo.txt")
	// Multiple operating systems generate different error messages
	AssertEqual(t, strings.Split(soteErr.FmtErrMsg, ": open foo.txt:")[0], "209010: foo.txt file was not found. Message return")
	soteErr = email.Attachment("/")
	AssertEqual(t, strings.Split(soteErr.FmtErrMsg, " >>Key:")[0], "210599: Business Service error has occurred that is not expected. ERROR DETAILS:")
	soteErr = email.Attachment("schema_test.json")
	AssertEqual(t, len(email.attachments), 1)
	AssertEqual(t, len(email.attachments[0].buffer) > 0, true)
	AssertEqual(t, email.attachments[0].filepath, "schema_test.json")
	AssertEqual(t, email.attachments[0].contentType, "text/plain; charset=utf-8")
}

func TestEmailValidateEmail(t *testing.T) {
	email := newTestEmail("staging", "Test Subject")
	soteErr := email.isEmailValid("To", "us")
	AssertEqual(t, soteErr.FmtErrMsg, "207050: us (To) is not a valid email address ERROR DETAILS: >>Key: EMAIL_LENGTH Value: mail: invalid string")
	soteErr = email.isEmailValid("To", "us.com")
	AssertEqual(t, soteErr.FmtErrMsg, "207050: us.com (To) is not a valid email address ERROR DETAILS: >>Key: EMAIL_PARSE Value: mail: misformatted email address")
	soteErr = email.isEmailValid("To", "invalid@ABCD.EFG")
	AssertEqual(t, soteErr.FmtErrMsg, "207050: invalid@ABCD.EFG (To) is not a valid email address ERROR DETAILS: >>Key: EMAIL_LOOKUP Value: mail: invalid domain in addr-spec")
}

func TestEmailAddAttachments(t *testing.T) {
	var buffer bytes.Buffer
	email := newTestEmail("staging", "Test Subject")
	email.attachments = append(email.attachments, attachment{
		filepath:    "schema_test.json",
		contentType: "text/plain; charset=utf-8",
		buffer:      []byte("Hello World"),
	})
	soteErr := email.addAttachments("CONTENT_BOUNDARY", &buffer)
	AssertEqual(t, soteErr.FmtErrMsg, "")
	re := regexp.MustCompile(`\r?\n`)
	AssertEqual(t, re.ReplaceAllString(buffer.String(), " "), re.ReplaceAllString(`
--MIX_CONTENT_BOUNDARY
Content-Type: text/plain; charset=utf-8; name="schema_test.json"
Content-Disposition: attachment; filename="schema_test.json"
Content-Transfer-Encoding: base64
SGVsbG8gV29ybGQ=`, " "))
}

func TestEmailInit(t *testing.T) {
	email := newTestEmail("staging", "Test Subject")
	auth, soteErr := email.initAuth()
	AssertEqual(t, soteErr.FmtErrMsg, "")
	AssertEqual(t, fmt.Sprintf("%v", auth), "&{ USERNAME PASSWORD email-smtp.eu-west-1.amazonaws.com}")
}

func TestEmailSendMail(t *testing.T) {
	email := newTestEmail("staging", "Test Subject")
	email.To("toAdmin@getsote.com")
	email.To("toSales@getsote.com").Name("To Customer Support")
	email.Cc("ccAdmin@getsote.com")
	email.Cc("ccSales@getsote.com").Name("CC Customer Support")
	email.Bcc("bccAdmin@getsote.com")
	email.Bcc("bccSales@getsote.com").Name("BCC Customer Support")
	email.GetSmtpUsername = func(application, environment string) (string, sError.SoteError) {
		return "USERNAME", sError.SoteError{}
	}
	email.GetSmtpPassword = func(application, environment string) (string, sError.SoteError) {
		return "PASSWORD", sError.SoteError{}
	}
	email.attachments = append(email.attachments, attachment{
		filepath:    "schema_test.json",
		contentType: "text/plain; charset=utf-8",
		buffer:      []byte("Hello World"),
	})
	email.sendMail = func(addr string, auth smtp.Auth, from string, to []string, msg []byte) error {
		AssertEqual(t, addr, "email-smtp.eu-west-1.amazonaws.com:587")
		AssertEqual(t, fmt.Sprintf("%v", auth), "&{ USERNAME PASSWORD email-smtp.eu-west-1.amazonaws.com}")
		AssertEqual(t, from, "support@getsote.com")
		AssertEqual(t, strings.Join(to, ";"), "toAdmin@getsote.com;toSales@getsote.com;ccAdmin@getsote.com;ccSales@getsote.com;bccAdmin@getsote.com;bccSales@getsote.com")
		re := regexp.MustCompile(`\r?\n`)
		AssertEqual(t, re.ReplaceAllString(string(msg), " "), re.ReplaceAllString(`MIME-Version: 1.0
Subject: Test Subject
Content-Type: multipart/mixed; boundary=MIX_CONTENT_BOUNDARY
From: support@getsote.com
To: toAdmin@getsote.com;To Customer Support <toSales@getsote.com>
Cc: ccAdmin@getsote.com;CC Customer Support <ccSales@getsote.com>
Bcc: bccAdmin@getsote.com;BCC Customer Support <bccSales@getsote.com>

--MIX_CONTENT_BOUNDARY
Content-Type: multipart/alternative; boundary=ALT_CONTENT_BOUNDARY

--ALT_CONTENT_BOUNDARY
Content-Type: text/plain; charset="UTF-8"

Hello World

--ALT_CONTENT_BOUNDARY
Content-Type: text/html; charset="UTF-8"

<p>Hello World</p>


--ALT_CONTENT_BOUNDARY--

--MIX_CONTENT_BOUNDARY
Content-Type: text/plain; charset=utf-8; name="schema_test.json"
Content-Disposition: attachment; filename="schema_test.json"
Content-Transfer-Encoding: base64
SGVsbG8gV29ybGQ=

--MIX_CONTENT_BOUNDARY--
`, " "))
		return nil
	}
	email.Send("Hello World", "<p>Hello World</p>")
}

func TestEmailSendMailFailed(t *testing.T) {
	email := newTestEmail("staging", "Test Subject")
	soteErr := email.Send("Hello World")
	AssertEqual(t, soteErr.FmtErrMsg, "210599: Business Service error has occurred that is not expected. ERROR DETAILS: >>Key: SMTP_ERROR Value: 535 Authentication Credentials Invalid")
}

func TestEmailSendMailInvalidFrom(t *testing.T) {
	email := newTestEmail("staging", "Test Subject")
	email.from.address = "invalid.com"
	soteErr := email.Send("Hello World")
	AssertEqual(t, soteErr.FmtErrMsg, "207050: invalid.com (From) is not a valid email address ERROR DETAILS: >>Key: EMAIL_PARSE Value: mail: misformatted email address")
}

func TestEmailSendMailInvalidTo(t *testing.T) {
	email := newTestEmail("staging", "Test Subject")
	email.To("invalid.com")
	soteErr := email.Send("Hello World")
	AssertEqual(t, soteErr.FmtErrMsg, "207050: invalid.com (To) is not a valid email address ERROR DETAILS: >>Key: EMAIL_PARSE Value: mail: misformatted email address")
}

func TestEmailSendMailInvalidCc(t *testing.T) {
	email := newTestEmail("staging", "Test Subject")
	email.Cc("invalid.com")
	soteErr := email.Send("Hello World")
	AssertEqual(t, soteErr.FmtErrMsg, "207050: invalid.com (Cc) is not a valid email address ERROR DETAILS: >>Key: EMAIL_PARSE Value: mail: misformatted email address")
}

func TestEmailSendMailInvalidBcc(t *testing.T) {
	email := newTestEmail("staging", "Test Subject")
	email.Bcc("invalid.com")
	soteErr := email.Send("Hello World")
	AssertEqual(t, soteErr.FmtErrMsg, "207050: invalid.com (Bcc) is not a valid email address ERROR DETAILS: >>Key: EMAIL_PARSE Value: mail: misformatted email address")
}

func TestEmailGetSmtpUsername(t *testing.T) {
	email := NewEmail("staging", "Test Subject")
	_, soteErr := email.getSmtpUsername("application", "environment")
	AssertEqual(t, soteErr.FmtErrMsg, "209110: environment value (environment) is invalid")
}

func TestEmailGetSmtpPassword(t *testing.T) {
	email := NewEmail("staging", "Test Subject")
	_, soteErr := email.GetSmtpPassword("application", "environment")
	AssertEqual(t, soteErr.FmtErrMsg, "209110: environment value (environment) is invalid")
}

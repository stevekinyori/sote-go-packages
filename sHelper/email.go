package sHelper

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/smtp"
	"os"
	"regexp"
	"strings"

	"gitlab.com/soteapps/packages/v2021/sConfigParams"
	"gitlab.com/soteapps/packages/v2021/sError"
	"gitlab.com/soteapps/packages/v2021/sLogger"
)

const (
	DEFAULTEMAIL = "support@getsote.com"
	SMTPHOST     = "email-smtp.eu-west-1.amazonaws.com"
	SMTPPORT     = 587
)

var (
	emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

type Email struct {
	subject     string
	environment string
	rcpt        []string
	from        *emailItem
	to          []*emailItem
	cc          []*emailItem
	bcc         []*emailItem
	attachments []attachment

	GetSmtpUsername func(application, environment string) (string, sError.SoteError)
	GetSmtpPassword func(application, environment string) (string, sError.SoteError)
	Send            func(text string, htmls ...string) sError.SoteError
	sendMail        func(addr string, a smtp.Auth, from string, to []string, msg []byte) error
}

type emailItem struct {
	name    string
	address string
}

type attachment struct {
	filepath    string
	contentType string
	buffer      []byte
}

func (e *emailItem) Name(name string) *emailItem {
	e.name = name
	return e
}

func (e *emailItem) String() string {
	if e.name != "" {
		return fmt.Sprintf("%s <%s>", e.name, e.address)
	}
	return e.address
}

func NewEmail(environment, subject string, from ...string) *Email {
	sLogger.DebugMethod()
	var email Email
	email = Email{
		subject:     subject,
		environment: environment,
		from: &emailItem{
			address: DEFAULTEMAIL,
		},
		GetSmtpUsername: email.getSmtpUsername,
		GetSmtpPassword: email.getSmtpPassword,
		Send:            email.send,
		sendMail:        smtp.SendMail,
	}
	if len(from) > 0 {
		email.from.name = from[0]
	}
	return &email
}

func (m *Email) To(address string) *emailItem {
	return m.addAddress(&m.to, address)
}

func (m *Email) Cc(address string) *emailItem {
	return m.addAddress(&m.cc, address)
}

func (m *Email) Bcc(address string) *emailItem {
	return m.addAddress(&m.bcc, address)
}

func (m *Email) Attachment(filepath string) (soteErr sError.SoteError) {
	f, err := os.Open(filepath)
	if err != nil {
		soteErr = NewError().FileNotFound(filepath, err.Error())
	} else {
		reader := bufio.NewReader(f)
		buffer, err := ioutil.ReadAll(reader)
		defer f.Close()
		if err != nil {
			soteErr = NewError(map[string]string{filepath: err.Error()}).InternalError()
		} else {
			contentType := http.DetectContentType(buffer)
			attach := attachment{
				filepath:    filepath,
				contentType: contentType,
				buffer:      buffer,
			}
			m.attachments = append(m.attachments, attach)
		}
	}
	return
}

func (m *Email) send(text string, htmls ...string) (soteErr sError.SoteError) {
	sLogger.DebugMethod()
	var (
		auth     smtp.Auth
		emails   []string
		buffer   bytes.Buffer
		boundary = "CONTENT_BOUNDARY"
	)
	auth, soteErr = m.initAuth()
	if soteErr.ErrCode == nil {
		buffer.WriteString(fmt.Sprintf("%s: %s\r\n", "MIME-Version", "1.0"))
		buffer.WriteString(fmt.Sprintf("%s: %s\r\n", "Subject", m.subject))
		buffer.WriteString(fmt.Sprintf("%s: %s\r\n", "Content-Type", "multipart/mixed; boundary=MIX_"+boundary))

		if soteErr = m.isEmailValid("From", m.from.address); soteErr.ErrCode != nil {
			return
		}
		buffer.WriteString(fmt.Sprintf("%s: %s\r\n", "From", m.from.String()))

		if emails, soteErr = m.addRcpt("To", m.to); soteErr.ErrCode != nil {
			return
		}
		buffer.WriteString(fmt.Sprintf("%s: %s\r\n", "To", strings.Join(emails, ";")))

		if emails, soteErr = m.addRcpt("Cc", m.cc); soteErr.ErrCode != nil {
			return
		}
		buffer.WriteString(fmt.Sprintf("%s: %s\r\n", "Cc", strings.Join(emails, ";")))

		if emails, soteErr = m.addRcpt("Bcc", m.bcc); soteErr.ErrCode != nil {
			return
		}
		buffer.WriteString(fmt.Sprintf("%s: %s\r\n", "Bcc", strings.Join(emails, ";")))

		buffer.WriteString("\r\n--MIX_" + boundary + "\r\n")
		buffer.WriteString(fmt.Sprintf("%s: %s\r\n", "Content-Type", "multipart/alternative; boundary=ALT_"+boundary))

		buffer.WriteString("\r\n--ALT_" + boundary + "\r\n")
		buffer.WriteString("Content-Type: text/plain; charset=\"UTF-8\"\r\n\r\n")
		buffer.WriteString(text)

		if len(htmls) > 0 {
			for _, html := range htmls {
				buffer.WriteString("\r\n\r\n--ALT_" + boundary + "\r\n")
				buffer.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n\r\n")
				buffer.WriteString(html)
			}
		}

		buffer.WriteString("\r\n")
		buffer.WriteString("\r\n\r\n--ALT_" + boundary + "--\r\n")

		soteErr = m.addAttachments(boundary, &buffer)
		buffer.WriteString("\r\n\r\n--MIX_" + boundary + "--\r\n")

		if soteErr.ErrCode == nil {
			err := m.sendMail(
				fmt.Sprintf("%s:%v", SMTPHOST, SMTPPORT),
				auth,
				m.from.address,
				m.rcpt,
				buffer.Bytes(),
			)
			if err != nil {
				soteErr = NewError(map[string]string{"SMTP_ERROR": err.Error()}).InternalError()
			}
		}
	}
	return
}

func (m *Email) initAuth() (auth smtp.Auth, soteErr sError.SoteError) {
	sLogger.DebugMethod()
	var (
		username string
		password string
	)
	username, soteErr = m.GetSmtpUsername("api", m.environment)
	if soteErr.ErrCode == nil {
		password, soteErr = m.GetSmtpPassword("api", m.environment)
		if soteErr.ErrCode == nil {
			auth = smtp.PlainAuth("", username, password, SMTPHOST)
		}
	}
	return
}

func (m *Email) addAttachments(boundary string, buffer *bytes.Buffer) (soteErr sError.SoteError) {
	sLogger.DebugMethod()
	for _, attach := range m.attachments {
		path := strings.Split(strings.ReplaceAll(attach.filepath, "\\", "/"), "/")
		name := path[len(path)-1]
		buffer.WriteString("\r\n--MIX_" + boundary + "\r\n")
		buffer.WriteString("Content-Type: " + attach.contentType + "; name=\"" + name + "\"\r\n")
		buffer.WriteString("Content-Disposition: attachment; filename=\"" + name + "\"\r\n")
		buffer.WriteString("Content-Transfer-Encoding: base64\r\n")
		buffer.WriteString(base64.StdEncoding.EncodeToString(attach.buffer))
	}
	return
}

func (m *Email) addAddress(emails *[]*emailItem, address string) *emailItem {
	item := &emailItem{
		address: address,
	}
	*emails = append(*emails, item)
	return item
}

func (m *Email) addRcpt(fieldName string, addresses []*emailItem) (emails []string, soteErr sError.SoteError) {
	for _, item := range addresses {
		if soteErr = m.isEmailValid(fieldName, item.address); soteErr.ErrCode != nil {
			return
		}
		m.rcpt = append(m.rcpt, item.address)
		emails = append(emails, item.String())
	}
	return
}

func (m *Email) isEmailValid(fieldName string, e string) sError.SoteError {
	if len(e) < 3 || len(e) > 254 {
		return NewError(map[string]string{"EMAIL_LENGTH": "mail: invalid string"}).InvalidEmailAddress(fieldName, e)
	}
	if !emailRegex.MatchString(e) {
		return NewError(map[string]string{"EMAIL_PARSE": "mail: misformatted email address"}).InvalidEmailAddress(fieldName, e)
	}
	parts := strings.Split(e, "@")
	mx, err := net.LookupMX(parts[1])
	if err != nil || len(mx) == 0 {
		return NewError(map[string]string{"EMAIL_LOOKUP": "mail: invalid domain in addr-spec"}).InvalidEmailAddress(fieldName, e)
	}
	return sError.SoteError{}
}

func (m *Email) getSmtpUsername(application, environment string) (string, sError.SoteError) {
	return sConfigParams.GetSmtpUsername(application, environment)
}

func (m *Email) getSmtpPassword(application, environment string) (string, sError.SoteError) {
	return sConfigParams.GetSmtpPassword(application, environment)
}

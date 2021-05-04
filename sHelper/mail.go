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
	"time"

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

type Mail struct {
	subject     string
	environment string
	rcpt        []string
	from        *emailItem
	to          []*emailItem
	cc          []*emailItem
	bcc         []*emailItem
	attachments []attachment
	auth        smtp.Auth
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

func NewMail(environment, subject string, from ...string) Mail {
	sLogger.DebugMethod()
	soteMail := Mail{
		subject:     subject,
		environment: environment,
		from: &emailItem{
			address: DEFAULTEMAIL,
		},
	}
	if len(from) > 0 {
		soteMail.from.name = from[0]
	}
	return soteMail
}

func (m *Mail) To(address string) *emailItem {
	return m.addAddress(&m.to, address)
}

func (m *Mail) Cs(address string) *emailItem {
	return m.addAddress(&m.cc, address)
}

func (m *Mail) Bcc(address string) *emailItem {
	return m.addAddress(&m.bcc, address)
}

func (m *Mail) Attachment(filepath string) (soteErr sError.SoteError) {
	f, err := os.Open(filepath)
	if err != nil {
		soteErr = NewError().FileNotFound(filepath, err.Error())
	} else {
		reader := bufio.NewReader(f)
		buffer, err := ioutil.ReadAll(reader)
		defer f.Close()
		if err != nil {
			soteErr = NewError().FileNotFound(filepath, err.Error())
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

func (m *Mail) Send(text string, htmls ...string) (soteErr sError.SoteError) {
	sLogger.DebugMethod()
	soteErr = m.initAuth()
	var (
		emails   []string
		buffer   bytes.Buffer
		boundary = "CONTENT_BOUNDARY"
	)
	if soteErr.ErrCode == nil {
		buffer.WriteString(fmt.Sprintf("%s: %s\r\n", "MIME-Version", "1.0"))
		buffer.WriteString(fmt.Sprintf("%s: %s\r\n", "Date", time.Now().String()))
		buffer.WriteString(fmt.Sprintf("%s: %s\r\n", "Subject", m.subject))
		buffer.WriteString(fmt.Sprintf("%s: %s\r\n", "Content-Type", "multipart/related;boundary="+boundary))

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

		buffer.WriteString("\r\n--" + boundary + "\r\n")
		buffer.WriteString("Content-Type: text/plain; charset=\"UTF-8\"\r\n\r\n")
		buffer.WriteString(text)

		for _, html := range htmls {
			buffer.WriteString("\r\n\r\n--" + boundary + "\r\n")
			buffer.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n\r\n")
			buffer.WriteString(html)
		}

		buffer.WriteString("\r\n")
		soteErr = m.addAttachments(boundary, &buffer)
		buffer.WriteString("\r\n\r\n--" + boundary + "--\r\n")

		if soteErr.ErrCode == nil {
			err := smtp.SendMail(
				fmt.Sprintf("%s:%v", SMTPHOST, SMTPPORT),
				m.auth,
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

func (m *Mail) initAuth() (soteErr sError.SoteError) {
	sLogger.DebugMethod()
	var (
		username string
		password string
	)
	if soteErr = sConfigParams.ValidateEnvironment(m.environment); soteErr.ErrCode == nil {
		username, password, soteErr = sConfigParams.GetSmtpUsernameAndPassword("api", m.environment)
		if soteErr.ErrCode == nil {
			m.auth = smtp.PlainAuth("", username, password, SMTPHOST)
		}
	}
	return
}

func (m *Mail) addAttachments(boundary string, buffer *bytes.Buffer) (soteErr sError.SoteError) {
	sLogger.DebugMethod()
	for _, attach := range m.attachments {
		path := strings.Split(strings.ReplaceAll(attach.filepath, "\\", "/"), "/")
		name := path[len(path)-1]
		buffer.WriteString("\r\n--" + boundary + "\r\n")
		buffer.WriteString("Content-Type: " + attach.contentType + "; name=\"" + name + "\"\r\n")
		buffer.WriteString("Content-Disposition: attachment; filename=\"" + name + "\"\r\n")
		buffer.WriteString("Content-Transfer-Encoding: base64\r\n")
		buffer.WriteString("Content-ID: <" + UUID(UUIDKind.Short) + "> \r\n")
		buffer.WriteString(base64.StdEncoding.EncodeToString(attach.buffer))
	}
	return
}

func (m *Mail) addAddress(emails *[]*emailItem, address string) *emailItem {
	item := &emailItem{
		address: address,
	}
	*emails = append(*emails, item)
	return item
}

func (m *Mail) addRcpt(fieldName string, addresses []*emailItem) (emails []string, soteErr sError.SoteError) {
	for _, item := range addresses {
		if m.isEmailValid(fieldName, item.address); soteErr.ErrCode != nil {
			return
		}
		m.rcpt = append(m.rcpt, item.address)
		if item.name != "" {
			emails = append(emails, item.String())
		}
	}
	return
}

func (m *Mail) isEmailValid(fieldName string, e string) sError.SoteError {
	if len(e) < 3 && len(e) > 254 {
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

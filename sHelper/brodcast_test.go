package sHelper

import (
	"fmt"
	"regexp"
	"testing"

	"gitlab.com/soteapps/packages/v2021/sDatabase"
	"gitlab.com/soteapps/packages/v2021/sError"
)

func createBrodcast(t *testing.T) *Brodcast {
	env, _ := NewEnvironment(ENVDEFAULTAPPNAME, ENVDEFAULTTARGET, ENVDEFAULTTARGET)
	run := NewRun(env)
	s := NewSubscriber(run, "test-consumer", "test-subject")
	b := &Brodcast{
		subscriber: s,
	}
	s.Publish = func(message interface{}, subject ...string) sError.SoteError {
		re := regexp.MustCompile(`\r?\n|\s+`)
		AssertEqual(t, subject[0], "10000.5d5147e2-57fc-48a6-b493-1783931ae9c0")
		AssertEqual(t, re.ReplaceAllString(fmt.Sprint(message), " "), `{  "id": 123,  "message-id": "12345",  "row-version": 10,  "subject": "bsl.organization.remove" }`)
		return sError.SoteError{}
	}
	return b
}

func TestBroadcastMessage(t *testing.T) {
	var (
		queryClose *PatchGuard
		queryExec  *PatchGuard
		b          = createBrodcast(t)
	)
	header := RequestHeaderSchema{
		MessageId:      "12345",
		AwsUserName:    "soteuser",
		OrganizationId: 10000,
	}
	queryClose = Patch(Query.Close, func(Query, sDatabase.SRows, *sError.SoteError) { queryClose.Unpatch() })
	queryExec = Patch(Query.Exec, func(q Query, r *Run) (sDatabase.SRows, sError.SoteError) {
		queryExec.Unpatch()
		AssertEqual(t, q.Sql.String(), "SELECT DISTINCT cognito_username")
		AssertEqual(t, q.Join, "INNER JOIN sote.usermanagementhistory AS H ON M.usermanagement_id = H.usermanagement_id")
		AssertEqual(t, q.Where, fmt.Sprintf("H.cognito_username != '%v' AND M.user_organizations_id = %v", header.AwsUserName, header.OrganizationId))
		rows := sDatabase.Rows{}
		index := 0
		rows.IValues = func() ([]interface{}, error) {
			return []interface{}{"5d5147e2-57fc-48a6-b493-1783931ae9c0"}, nil
		}
		rows.INext = func() bool {
			index++
			return index == 1
		}
		return rows, sError.SoteError{}
	})
	b.Message(sError.SoteError{}, header, &Msg{Subject: "bsl.organization.remove"}, 123, 10)
}

func TestBroadcastMessageError(t *testing.T) {
	var (
		queryClose *PatchGuard
		queryExec  *PatchGuard
		b          = createBrodcast(t)
	)
	queryClose = Patch(Query.Close, func(Query, sDatabase.SRows, *sError.SoteError) { queryClose.Unpatch() })
	queryExec = Patch(Query.Exec, func(q Query, r *Run) (sDatabase.SRows, sError.SoteError) {
		queryExec.Unpatch()
		return nil, NewError().NoDbConnection()
	})
	defer func() {
		r := recover()
		AssertEqual(t, fmt.Sprint(r), "209299: No database connection has been established")
	}()
	r := RequestHeaderSchema{
		MessageId:      "12345",
		AwsUserName:    "soteuser",
		OrganizationId: 10000,
	}
	b.Message(sError.SoteError{}, r, &Msg{Subject: "bsl.organization.remove"}, 123)
}

package sHelper

import (
	"encoding/json"
	"fmt"

	"gitlab.com/soteapps/packages/v2021/sError"
	"gitlab.com/soteapps/packages/v2021/sLogger"
)

type Brodcast struct {
	subscriber *Subscriber
}

func (b *Brodcast) Message(soteErr sError.SoteError, header RequestHeaderSchema, message *Msg, id int64, row_version ...int64) {
	sLogger.DebugMethod()
	var (
		s    = b.subscriber
		data = map[string]interface{}{
			"subject":    message.Subject,
			"message-id": header.MessageId,
			"id":         id,
		}
	)
	if len(row_version) > 0 {
		data["row-version"] = row_version[0]
	}
	str_data, _ := json.MarshalIndent(data, "", "\t")
	if soteErr.ErrCode == nil {
		query := Query{
			Table:   "usermanagement AS M",
			Columns: []string{"DISTINCT cognito_username"},
			Join:    "INNER JOIN sote.usermanagementhistory AS H ON M.usermanagement_id = H.usermanagement_id",
			Where:   fmt.Sprintf("H.cognito_username != '%v' AND M.user_organizations_id = %v", header.AwsUserName, header.OrganizationId),
		}
		tRows, soteErr := query.Select().Exec(s.Run)
		if soteErr.ErrCode == nil {
			for tRows.Next() {
				tCols, _ := tRows.Values()
				s.Publish(string(str_data), fmt.Sprintf("%v.%v", header.OrganizationId, tCols[0]))
			}
		}
		query.Close(tRows, &soteErr)
		if soteErr.ErrCode != nil {
			s.Run.PanicService(soteErr)
		}
	}
}

package sDatabase

import (
	"github.com/jackc/pgconn"
	"github.com/jackc/pgproto3/v2"
)

type Rows struct {
	IEerr              func() error
	ICommandTag        func() pgconn.CommandTag
	IFieldDescriptions func() []pgproto3.FieldDescription
	INext              func() bool
	IScan              func(dest ...interface{}) error
	IValues            func() ([]interface{}, error)
	IRawValues         func() [][]byte
}

func (r Rows) Close()                                         {}
func (r Rows) Err() error                                     { return r.IEerr() }
func (r Rows) CommandTag() pgconn.CommandTag                  { return r.ICommandTag() }
func (r Rows) FieldDescriptions() []pgproto3.FieldDescription { return r.IFieldDescriptions() }
func (r Rows) Next() bool                                     { return r.INext() }
func (r Rows) Scan(dest ...interface{}) error                 { return r.IScan(dest...) }
func (r Rows) Values() ([]interface{}, error)                 { return r.IValues() }
func (r Rows) RawValues() [][]byte                            { return r.IRawValues() }

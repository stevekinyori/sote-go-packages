package sHelper

// FilterFields describes the  format for filter expression in the shipment fintrans request json payload
type FilterFields struct {
	FieldName string      `json:"field-name" validate:"required,oneof=shipments-id client-company-id fintrans-type currency-type transactions-amount transactions-timestamp created-by-requestor-username bill-of-lading memo cost-is-unexpected"`
	Operator  string      `json:"operator" mod:"ucase" validate:"required,oneof=> < = >= <= != <> IN 'NOT IN'"`
	Value     interface{} `json:"value" validate:"required"`
}

// ArrFilterParam filter parameters for a slice/array
type ArrFilterParam struct {
	FilterFields
	Prefix            string
	InitialParamCount int
	CaseInsensitive   bool
}

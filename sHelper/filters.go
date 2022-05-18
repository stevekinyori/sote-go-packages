package sHelper

// FilterCommon describes the  format for common filter expression in the request json payload
type FilterCommon struct {
	Operator string      `json:"operator" mod:"ucase" validate:"required,oneof=> < = >= <= != <> IN 'NOT IN'"`
	Value    interface{} `json:"value" validate:"required_unless=Operator != Operator ="`
}

// FilterFields describes the  format for all filter expression in the request json payload
type FilterFields struct {
	FilterCommon
	FieldName string `json:"field-name"`
}

// ShipmentFinTransFilterFields describes the  format for filter expression in the shipment fintrans request json payload
type ShipmentFinTransFilterFields struct {
	FieldName string `json:"field-name" validate:"required,oneof=shipments-id client-company-id fintrans-type currency-type transactions-amount transactions-timestamp created-by-requestor-username bill-of-lading memo cost-is-unexpected"`
	FilterCommon
}

// ArrFilterParam filter parameters for a slice/array
type ArrFilterParam struct {
	FilterFields
	Prefix            string
	InitialParamCount int
	CaseInsensitive   bool
}

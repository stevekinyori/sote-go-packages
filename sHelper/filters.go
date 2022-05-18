package sHelper

import (
	"gitlab.com/soteapps/packages/v2021/sLogger"
	"reflect"
)

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
	FilterCommon
	FieldName string `json:"field-name" validate:"required,oneof=shipments-id client-company-id fintrans-type currency-type transactions-amount transactions-timestamp created-by-requestor-username bill-of-lading memo cost-is-unexpected"`
}

// ShipmentsFilterFields describes the format for filter expression in the shipment's request json payload
type ShipmentsFilterFields struct {
	FilterCommon
	FieldName string `json:"field-name" validate:"required,oneof=shipments-id row-version organizations-id client-company-id client-company-name suppliers-id supplier-name shipment-status bill-of-lading port-of-entry port-of-origin country-destination shipment-port-of-entry-eta consignment-wgt consignment-wgt-uom through-bill-of-lading created-by-date updated-by-date created-by-aws-user-name updated-by-aws-user-name"`
}

// ContainerFilterFields describes the format for filter expression in the container's request json payload
type ContainerFilterFields struct {
	FilterCommon
	FieldName string `json:"field-name" validate:"required,oneof=containers-id organizations-id row-version container-tag container-form-factor shipments-id trips-id cargo-gross-wgt cargo-gross-uom created-by-date updated-by-date created-by-aws-user-name updated-by-aws-user-name"`
}

// ArrFilterParam filter parameters for a slice/array
type ArrFilterParam struct {
	FilterFields
	Prefix            string
	InitialParamCount int
	CaseInsensitive   bool
}

// SetFilters set convert filters of map[string][]ShipmentsFilterFields | map[string][] ContainerFilterFields to map[string][]FilterFields
func SetFilters[F ShipmentsFilterFields | ContainerFilterFields | ShipmentFinTransFilterFields](tFilters map[string][]F) (filters map[string][]FilterFields) {
	sLogger.DebugMethod()

	if len(tFilters) > 0 {
		filters = make(map[string][]FilterFields, 0)
		for operand, filterValues := range tFilters {
			for _, v := range filterValues {
				field := FilterFields{}
				f := reflect.ValueOf(v).Interface()
				switch f.(type) {
				case ShipmentsFilterFields:
					field = FilterFields(f.(ShipmentsFilterFields))
				case ContainerFilterFields:
					field = FilterFields(f.(ContainerFilterFields))
				case ShipmentFinTransFilterFields:
					field = FilterFields(f.(ShipmentFinTransFilterFields))
				}
				filters[operand] = append(filters[operand], field)
			}
		}
	}
	return
}

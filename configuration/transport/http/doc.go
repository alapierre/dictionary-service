package http

import "dictionaries-service/configuration"

// swagger:parameters deleteConfigurationEntry
//goland:noinspection GoUnusedType
type deleteConfigurationEntryRequestWrapper struct {

	// optional tenant id
	// in:header
	Tenant string `json:"X-Tenant-ID"`

	// Configuration body
	// in:body
	Body deleteConfigurationRequest
}

// swagger:parameters updateConfiguration
//goland:noinspection GoUnusedType
type updateConfigurationRequestWrapper struct {

	// optional tenant id
	// in:header
	Tenant string `json:"X-Tenant-ID"`

	// Configuration body
	// in:body
	Body saveConfigurationRequest
}

// swagger:parameters saveConfiguration
//goland:noinspection GoUnusedType
type saveConfigurationRequestWrapper struct {

	// optional tenant id
	// in:header
	Tenant string `json:"X-Tenant-ID"`

	// Configuration body
	// in:body
	Body saveConfigurationRequest
}

// swagger:parameters loadAllConfigKeys
//goland:noinspection GoUnusedType
type loadAllShortRequest struct {

	// optional tenant id
	// in:header
	Tenant string `json:"X-Tenant-ID"`
}

// swagger:response loadConfigurationResponse
//goland:noinspection ALL
type loadConfigurationResponseWrapper struct {

	// in:body
	Body []loadConfigurationResponse
}

//goland:noinspection ALL

// swagger:response loadShortResponseWrapper
//goland:noinspection ALL
type loadShortResponseWrapper struct {

	// in:body
	Body []configuration.Short
}

// swagger:response loadValuesResponseWrapper
//goland:noinspection ALL
type loadValuesResponseWrapper struct {

	// in:body
	Body []loadValueResponse
}

// swagger:response loadConfigurationOneResponseWrapper
//goland:noinspection ALL
type loadConfigurationOneResponseWrapper struct {

	// in:body
	Body loadConfigurationResponse
}

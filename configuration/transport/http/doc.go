package http

import "dictionaries-service/configuration"

// swagger:parameters updateConfigValue
type updateValueInTimeRequestWrapper struct {

	// optional tenant id
	// in:header
	Tenant string `json:"X-Tenant-ID"`

	// Config entry key
	// in:path
	Key string `json:"key"`

	// date from
	// in:path
	From string `json:"from"`

	// Configuration value
	// in:body
	Body updateValueInTimeRequest
}

// swagger:parameters addNewConfigEntry
type addNewConfigurationEntryRequestWrapper struct {

	// optional tenant id
	// in:header
	Tenant string `json:"X-Tenant-ID"`

	// in:path
	Key string `json:"key"`

	// Configuration body
	// in:body
	Body addNewConfigurationValueRequest
}

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
	Body loadValueResponseFull
}

// swagger:response loadConfigurationOneResponseWrapper
//goland:noinspection ALL
type loadConfigurationOneResponseWrapper struct {

	// in:body
	Body loadConfigurationResponse
}

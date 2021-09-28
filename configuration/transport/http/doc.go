package http

import "dictionaries-service/configuration"

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

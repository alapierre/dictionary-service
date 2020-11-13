package service

import (
	"database/sql"
	"dictionaries-service/model"
	slog "github.com/go-eden/slf4go"
	"time"
)

func NewConfigurationService(configurationRepository ConfigurationRepository) ConfigurationService {
	return &configurationService{configurationRepository: configurationRepository}
}

type configurationService struct {
	configurationRepository ConfigurationRepository
}

type ConfigurationService interface {
	LoadForDay(key, tenant string, day time.Time) (*model.Configuration, error)
	Load(key, tenant string, from, to time.Time) ([]model.Configuration, error)
	Save(conf *model.Configuration) error
	NewConfigValue(key, tenant, configType, name, value string) error
	LoadMany(tenant string, day time.Time, keys ...string) []model.Configuration
}

func (c *configurationService) LoadForDay(key, tenant string, day time.Time) (*model.Configuration, error) {
	return c.configurationRepository.LoadForDay(key, tenant, day)
}

func (c *configurationService) Load(key, tenant string, from, to time.Time) ([]model.Configuration, error) {
	return c.configurationRepository.Load(key, tenant, from, to)
}

func (c *configurationService) LoadMany(tenant string, day time.Time, keys ...string) []model.Configuration {

	var res []model.Configuration

	for _, key := range keys {
		conf, err := c.LoadForDay(key, tenant, day)

		if err != nil {
			slog.Warnf("can't load key = %s tenant = %s cause: %v", key, tenant, err)
		}
		res = append(res, *conf)
	}

	return res
}

func (c *configurationService) Save(conf *model.Configuration) error {
	return c.configurationRepository.Save(conf)
}

func (c *configurationService) NewConfigValue(key, tenant, configType, name, value string) error {
	conf := &model.Configuration{
		Key:      key,
		Tenant:   tenant,
		Type:     configType,
		Name:     name,
		Value:    sql.NullString{String: value},
		DateFrom: BeginOfTheWorld,
		DateTo:   EndOfTheWorld,
	}

	return c.configurationRepository.Save(conf)
}

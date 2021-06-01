package confguration

import (
	"database/sql"
	slog "github.com/go-eden/slf4go"
	"time"
)

func NewService(configurationRepository Repository) Service {
	return &service{configurationRepository: configurationRepository}
}

type service struct {
	configurationRepository Repository
}

type Service interface {
	LoadForDay(key, tenant string, day time.Time) (*Configuration, error)
	Load(key, tenant string, from, to time.Time) ([]Configuration, error)
	Save(conf *Configuration) error
	NewConfigValue(key, tenant, configType, name, value string) error
	LoadMany(tenant string, day time.Time, keys ...string) []Configuration
}

func (c *service) LoadForDay(key, tenant string, day time.Time) (*Configuration, error) {
	return c.configurationRepository.LoadForDay(key, tenant, day)
}

func (c *service) Load(key, tenant string, from, to time.Time) ([]Configuration, error) {
	return c.configurationRepository.Load(key, tenant, from, to)
}

func (c *service) LoadMany(tenant string, day time.Time, keys ...string) []Configuration {

	var res []Configuration

	for _, key := range keys {
		conf, err := c.LoadForDay(key, tenant, day)

		if err != nil {
			slog.Warnf("can't load key = %s tenant = %s cause: %v", key, tenant, err)
		}
		res = append(res, *conf)
	}

	return res
}

func (c *service) Save(conf *Configuration) error {
	return c.configurationRepository.Save(conf)
}

func (c *service) NewConfigValue(key, tenant, configType, name, value string) error {
	conf := &Configuration{
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

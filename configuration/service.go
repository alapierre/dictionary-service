package configuration

import (
	"context"
	"database/sql"
	"dictionaries-service/tenant"
	"fmt"
	slog "github.com/go-eden/slf4go"
	"time"
)

func NewConfigurationService(configurationRepository Repository) Service {
	return &service{repository: configurationRepository}
}

type service struct {
	repository Repository
}

type Service interface {
	LoadForDay(key, tenant string, day time.Time) (*Configuration, error)
	Load(key, tenant string, from, to time.Time) ([]Configuration, error)
	Save(conf *Configuration) error
	NewConfigValue(key, tenant, configType, name, value string) error
	LoadMany(tenant string, day time.Time, keys ...string) []Configuration
	LoadAllShort(ctx context.Context) ([]Short, error)
}

func (c *service) LoadForDay(key, tenant string, day time.Time) (*Configuration, error) {
	return c.repository.LoadForDay(key, tenant, day)
}

func (c *service) Load(key, tenant string, from, to time.Time) ([]Configuration, error) {
	return c.repository.Load(key, tenant, from, to)
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

func (c *service) LoadAllShort(ctx context.Context) ([]Short, error) {

	t, ok := tenant.FromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("can't read tenant from context")
	}

	return c.repository.LoadAllShort(t.Name)
}

func (c *service) Save(conf *Configuration) error {
	return c.repository.Save(conf)
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

	return c.repository.Save(conf)
}

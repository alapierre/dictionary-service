# dictionary-service

[![Sonarcloud Status](https://sonarcloud.io/api/project_badges/measure?project=alapierre_dictionary-service&metric=alert_status)](https://sonarcloud.io/dashboard?id=alapierre_dictionary-service)
[![Renovate enabled](https://img.shields.io/badge/renovate-enabled-brightgreen.svg)](https://renovatebot.com/)

Fast and simple dictionary service on PostgreSQL database and JSON content

## Main goals

Almost any system needs to store and manage flexible dictionary values. Some of these need variable over time configuration. 

- provide fast REST microservice for storing dictionary in JSON format 
- dictionaries should be describable - for automatic GUI build
- i18n support for dictionary names
- store configuration values changeable over time (e.g. what will be configuration value of X in 2025-12-31?)
- store calendar like dictionaries (the dictionary key is a date) e.g. holiday calendar
- integrate with Netflix Eureka and oAtuh
- no additional, other than PostgreSQL database, for store dictionaries and configuration data - in cloud operators noSQL database is additional cost

## Current status

In production in several commercial projects. 

## The Latest news

- Configuration load (store and update - soon)
- Calendar dictionary load, store and update
- swagger support (if you want serve swagger gui or generate open API file, you need to install goswagger: https://goswagger.io/install.html)

## Required environment variables

- DICT_DATASOURCE_USER
- DICT_DATASOURCE_PASSWORD

...and more not mandatory settings can be found in main.go, eg:

- INIT_DB_CONNECTION_RETS - how many time try to connect to database - it is useful in development on docker-compose 
- DATASOURCE_SCHEMA - name of database schema to use
- DICT_SHOW_SQL - show SQL in DEBUG logs or not

Default service connects to database on localhost:5432 with schema dictionary and 100 retries

## How it works

Do not forget to set environment variables (check Makefile or docker-compose.yml).

### Run in bash

```
$ make run
```

### Run with Docker

```yaml
version: '3.1'
services:

  db:
    image: postgres:12
    volumes:
      - pg_data:/var/lib/postgresql/data
    environment:
      - POSTGRES_PASSWORD=qwedsazxc
      - POSTGRES_USER=app
    ports:
      - "5432:5432"

  eureka:
    image: lapierre/eureka:1.0.1
    ports:
      - "8761:8761"

  dict:
    image: lapierre/dictionary-service
    environment:
      - DICT_DATASOURCE_HOST=db:5432
      - DICT_DATASOURCE_PASSWORD=qwedsazxc
      - DICT_DATASOURCE_USER=app
      - DICT_EUREKA_SERVICE_URL=http://eureka:8761/eureka
    ports:
      - "9098:9098"

volumes:
  pg_data:
```

### Calendar dictionaries

When we need store e.g., list of free days like holiday calendar, the best suitable will be dictionary with date as a key 
and possibility to get all entries in date range. Then we can get all holidays e.g., in 2021 year. 
Some time we need extra information about free day, e.g., to count salary components for work on a holiday.
This is the reason the calendar API created.

Because in many cases we have common holiday calendar for more than one tenant (application, module) it is possible to 
use default tenant as a base calendar and override it in specific tenant if needed. 

In database:

| tenant     | type     | name         | kind           | labels                            | day        |
|------------|----------|--------------|----------------|-----------------------------------|------------|
|            | holidays | Nowy Rok     | public holiday | "holliday_type"=>"country"        | 2021-01-01 |
|            | holidays | Trzech Króli | public holiday | "holliday_type"=>"church holiday" | 2021-01-06 | 
| tenant1    | holidays | Barburka     | company holiday| "holliday_type"=>"company day"    | 2021-12-04 |

Load from default tenant (or without tenant at all - no `X-Tenant-ID` header):

````http request
GET http://localhost:9098/api/calendar/holidays/2021-01-01/2021-12-31
Content-Type: application/json
X-Tenant-ID: default
Accept-Language: en-EN
````

result

````json
[
  {
    "day": "2021-01-01",
    "tenant": "",
    "name": "Nowy Rok",
    "kind": "public holiday",
    "labels": {
      "holliday_type": "country"
    }
  },
  {
    "day": "2021-01-06",
    "tenant": "",
    "name": "Trzech Króli",
    "kind": "public holiday",
    "labels": {
       "holliday_type": "church holiday"
    }
  }
]
````

Load from tenant1:

````http request
GET http://localhost:9098/api/calendar/holidays/2021-01-01/2021-12-31
Content-Type: application/json
X-Tenant-ID: tenant1
Accept-Language: en-EN
````

````json
[
  {
    "day": "2021-12-04",
    "tenant": "tenant1",
    "name": "Barburka",
    "kind": "company holiday",
    "labels": {
       "holliday_type": "company day"
    }
  }
]
````

If in request header we put `X-Tenant-Merge-Default` with `true` value:

````http request
GET http://localhost:9098/api/calendar/holidays/2021-01-01/2021-12-31
Content-Type: application/json
X-Tenant-ID: tenant1
X-Tenant-Merge-Default: true
Accept-Language: en-EN
````

we will get:

````json
[
  {
    "day": "2021-01-01",
    "tenant": "",
    "name": "Nowy Rok",
    "kind": "public holiday",
    "labels": {
      "holliday_type": "country"
    }
  },
  {
    "day": "2021-01-06",
    "tenant": "",
    "name": "Trzech Króli",
    "kind": "public holiday",
    "labels": {
       "holliday_type": "church holiday"
    }
  },
 {
  "day": "2021-12-04",
  "tenant": "tenant1",
  "name": "Barburka",
  "kind": "company holiday",
  "labels": {
      "holliday_type": "company day"
  }
 }
]
````

#### Create new, edit and delete calendar entry

create

````http request
POST http://localhost:9098/api/calendar/holidays
Content-Type: application/json
X-Tenant-ID: default

{
  "day": "2021-04-05",
  "name": "Poniedziałek Wielkanocny",
  "kind": "public holiday"
}
````

update

````http request
PUT http://localhost:9098/api/calendar/holidays
Content-Type: application/json
X-Tenant-ID: default

{
  "day": "2021-04-05",
  "name": "Poniedziałek Wielkanocny 1234",
  "kind": "public holiday",
  "labels":  {
    "holliday_type": "church holiday"
  }
}
````

delete

````http request
DELETE http://localhost:9098/api/calendar/holidays/2021-04-05
Content-Type: application/json
X-Tenant-ID: default
````


### Configuration entries store and load

Every config parameter has own period of validity, e.g., in January smtp.server.host was localhost, 
but from next year it will have other value, and we want to set it now. So, we need two records in database:
 
 | value     | date_form | date_to  |
 |-----------|-----------|----------|
 | 127.0.1.1 |2020-01-01 |2020-12-31|
 | 34.0.11.1 |2021-01-01 |2999-12-31|

#### Get one value for given key, tenant and day

```http request
GET http://localhost:9098/api/config/smtp.server.host/2020-01-01
X-Tenant-ID: default
Accept-Language: en-EN
Cache-Control: no-cache
Accept: application/json
```

result:

```json
{
  "key": "smtp.server.host",
  "value": "127.0.0.1",
  "type": "String"
}
```

#### Get many keys, same tenant and day

```http request
GET http://localhost:9098/api/configs/2020-01-01?key=smtp.server.port&key=smtp.server.host
X-Tenant-ID: default
Accept-Language: en-EN
Cache-Control: no-cache
Accept: application/json
```

result

```json
[
  {
    "key": "smtp.server.port",
    "value": "25",
    "type": "Int"
  },
  {
    "key": "smtp.server.host",
    "value": "127.0.0.1",
    "type": "String"
  }
]
```

#### Store config values 

not implemented yet - put your config into database


### Create dictionary entry metadata

```http request
POST http://localhost:9098/api/metadata/AbsenceType
X-Tenant-ID: default
Accept-Language: en-EN
Cache-Control: no-cache
Content-Type: application/json

{
  "$id": "https://alapierre.io/dictionary.schema.json",
  "$schema": "http://json-schema.org/draft-07/schema#",
  "title": "DictionaryAbsenceType",
  "type": "object",
  "required": [ "onlyOnBeginOrEnd", "needDeliveryDateConfirmation", "needConfirmationDocumentNumber" ],
  "properties": {
    "onlyOnBeginOrEnd": {
      "type": "boolean",
      "description": "Absence can only start on beginning or finish on end of work day",
      "default": false
    },
    "needDeliveryDateConfirmation": {
      "type": "boolean",
      "default": false,
      "description": "Is proof of absence delivery date required - should field be visible on form"
    },
    "needConfirmationDocumentNumber": {
      "description": "Is absence confirmation document number needed",
      "type": "boolean",
      "default": false
    }
  }
}
```

### Save new dictionary entry

```http request
POST http://localhost:9098/api/dictionary
Content-Type: application/json
X-Tenant-ID: default
Accept-Language: en-EN

{
  "children": [
    {
      "key": "newCh1",
      "label": "wwww",
      "name": "Child 1"
    },
    {
      "key": "newCh2",
      "label": "qqqq",
      "name": "Child 2"
    }
  ],
  "key": "HollidayLeave",
  "name": "Holliday Leave",
  "needConfirmationDocumentNumber": false,
  "needDeliveryDateConfirmation": true,
  "onlyOnBeginOrEnd": true,
  "type": "AbsenceType"
}
```


### Loading stored dictionary entry

`GET /api/dictionary/{type}/{key}`


#### Load parent dictionary entry with children

```http request
GET /api/dictionary/AbsenceType/HollidayLeave
X-Tenant-ID: default
Accept-Language: en-EN
```

Result

```json
{
  "children": [
    {
      "key": "newCh1",
      "label": "wwww",
      "name": "Child 1"
    },
    {
      "key": "newCh2",
      "label": "qqqq",
      "name": "Child 2"
    }
  ],
  "key": "HollidayLeave",
  "name": "Holliday Leave",
  "needConfirmationDocumentNumber": false,
  "needDeliveryDateConfirmation": true,
  "onlyOnBeginOrEnd": true,  
  "type": "AbsenceType"
}
```

#### Load children only 

```http request
GET /api/dictionary/AbsenceType/HollidayLeave/children
X-Tenant-ID: default
Accept-Language: en-EN
```

Result

```json
[
    {
          "key": "newCh1",
          "label": "wwww",
          "name": "Child 1",
          "parent_key": "HollidayLeave"
        },
        {
          "key": "newCh2",
          "label": "qqqq",
          "name": "Child 2",
          "parent_key": "HollidayLeave"
        }
]
```

### Update existing dictionary entry

```http request
PUT http://localhost:9098/api/dictionary
Content-Type: application/json
X-Tenant-ID: default
Accept-Language: en-EN

{
  "children": [
    {
      "key": "newCh1",
      "label": "bbb",
      "name": "Child 1 updated",
      "type": "AbsenceType"
    },
    {
      "key": "newCh3",
      "label": "ttt",
      "name": "Child 3 new",
      "type": "AbsenceType"
    }
  ],
  "key": "HollidayLeave",
  "name": "New parent updated",
  "needConfirmationDocumentNumber": false,
  "needDeliveryDateConfirmation": false,
  "onlyOnBeginOrEnd": true,
  "type": "AbsenceType"
}
```

Result

```json
{
  "children": [    
    {
      "key": "newCh1",
      "label": "bbb",
      "name": "Child 1 updated",
      "type": "AbsenceType"
    },
    {
      "key": "newCh3",
      "label": "ttt",
      "name": "Child 3 new",
      "type": "AbsenceType"
    }
  ],
  "key": "HollidayLeave",
  "name": "New parent updated",
  "needConfirmationDocumentNumber": false,
  "needDeliveryDateConfirmation": false,
  "onlyOnBeginOrEnd": true,
  "tenant": "default",
  "type": "AbsenceType"
}
```

#### Load all available dictionary metadata

```http request
GET /api/dictionary/metadata
X-Tenant-ID: default
Accept-Language: en-EN
```
Result
```
["type1", "type2", "type3"]
```

#### Load dictionary metadata by type

```http request
GET /api/dictionary/metadata/{type}
X-Tenant-ID: default
Accept-Language: en-EN
```

Example result

```json
{
  "$id": "https://alapierre.io/dictionary.schema.json",
  "type": "object",
  "title": "DictionaryAbsenceTypeTitle",
  "$schema": "http://json-schema.org/draft-07/schema#",
  "properties": {
    "onlyOnBeginOrEnd": {
      "type": "boolean",
      "description": "Absence can only start on beginning or finish on end of work day"
    },
    "needDeliveryDateConfirmation": {
      "type": "boolean",
      "description": "Is proof of absence delivery date required - should field be visible on form"
    },
    "needConfirmationDocumentNumber": {
      "type": "boolean",
      "description": "Is absence confirmation document number needed"
    }
  }
}
```

### Update existing dictionary metadata

```http request
PUT http://localhost:9098/api/metadata/DictionaryAbsenceType
X-Tenant-ID: default
Accept-Language: en-EN
Cache-Control: no-cache
Content-Type: application/json

{
  "$id": "https://alapierre.io/dictionary.schema.json",
  "$schema": "http://json-schema.org/draft-07/schema#",
  "title": "DictionaryAbsenceType",
  "type": "object",
  "required": [ "onlyOnBeginOrEnd", "needDeliveryDateConfirmation", "needConfirmationDocumentNumber" ],
  "properties": {
    "onlyOnBeginOrEnd": {
      "type": "boolean",
      "description": "Absence can only start on beginning or finish on end of work day",
      "default": false
    },
    "needDeliveryDateConfirmation": {
      "type": "boolean",
      "default": false,
      "description": "Is proof of absence delivery date required - should field be visible on form"
    },
    "needConfirmationDocumentNumber": {
      "description": "Is absence confirmation document number needed",
      "type": "boolean",
      "default": false
    }
  }
}
```

Result
```
Empty body
```

### Next steps, coming soon

- True unit tests
- Integration tests

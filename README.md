# dictionary-service

[![Sonarcloud Status](https://sonarcloud.io/api/project_badges/measure?project=alapierre_dictionary-service&metric=alert_status)](https://sonarcloud.io/dashboard?id=alapierre_dictionary-service)
[![Renovate enabled](https://img.shields.io/badge/renovate-enabled-brightgreen.svg)](https://renovatebot.com/)

Fast and simple dictionary service on PostgreSQL database and JSON content

## Main goals

Almost any system needs to store and manage flexible dictionary values. Some of these need variable over time configuration. 

- provide fast REST microservice for storing dictionary in JSON format 
- store configuration values changeable over time (eg. what will be configuration value of X in 2025-12-31?)
- dictionaries should be describable - for automatic GUI build
- i18n support for dictionary names
- integrate with Netflix Eureka and oAtuh
- no additional, other than PostgreSQL database, for store dictionaries and configuration data - in cloud operators noSQL database is additional cost

## Current status

In tests on production in several commertial projects. 

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
    image: lapierre/dictionary-service:0.0.11
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

### Create dictionary entry metadata

```
POST http://localhost:9098/api/metadata/AbsenceType
X-Tenant: default
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

```
POST http://localhost:9098/api/dictionary
Content-Type: application/json
X-Tenant: default
Accept-Language: en-EN

{
  "children": [
    {
      "key": "newCh1",
      "label": "wwww",
      "name": "Child 1",
      "type": "AbsenceType"
    },
    {
      "key": "newCh2",
      "label": "qqqq",
      "name": "Child 2",
      "type": "AbsenceType"
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

```
GET /api/dictionary/AbsenceType/HollidayLeave
X-Tenant: default
Accept-Language: en-EN
```

Result

```json
{
  "children": [
    {
      "key": "newCh1",
      "label": "wwww",
      "name": "Child 1",
      "type": "AbsenceType"
    },
    {
      "key": "newCh2",
      "label": "qqqq",
      "name": "Child 2",
      "type": "AbsenceType"
    }
  ],
  "key": "HollidayLeave",
  "name": "Holliday Leave",
  "needConfirmationDocumentNumber": false,
  "needDeliveryDateConfirmation": true,
  "onlyOnBeginOrEnd": true,
  "tenant": "default",
  "type": "AbsenceType"
}
```

#### Load child only dictionary entry

```
GET /api/dictionary/AbsenceType/HollidayLeave/newCh1
X-Tenant: default
Accept-Language: en-EN
```

Result

```json
{
  "key": "newCh1",
  "label": "wwww",
  "name": "Child 1",
  "parent_key": "newP",
  "tenant": "default",
  "type": "AbsenceType"
}
```

### Update existing dictionary entry

```
PUT http://localhost:9098/api/dictionary
Content-Type: application/json
X-Tenant: default
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

```
GET /api/dictionary/metadata
X-Tenant: default
Accept-Language: en-EN
```
Result
```
["type1", "type2", "type3"]
```

#### Load dictionary metadata by type

```
GET /api/dictionary/metadata/{type}
X-Tenant: default
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

### Update existing dictionary metadata

```
###
PUT http://localhost:9098/api/metadata/DictionaryAbsenceType
X-Tenant: default
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

- Translate dictionary name base on Accept-Language header
- Configuration store, load and update
- True unit tests
- Integration tests

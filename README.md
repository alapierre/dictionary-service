# dictionary-service

Project is in heavy development status!

fast and simple dictionary service on PostgreSQL database and JSON content

## Main goals

Almost any system needs to store and manage flexible dictionary values. Some of these need variable over time configuration. 

- provide fast REST microservice for storing dictionary in JSON format 
- store configuration values changeable over time (eg. what will be configuration value of X in 2025-12-31?)
- dictionaries should be describable - for automatic GUI build
- i18n support for dictionary names
- integrate with Netflix Eureka and oAtuh
- no additional, other than PostgreSQL database, for store dictionaries and configuration data - in cloud operators noSQL database is additional cost

## Required environment variables

- DICT_DATASOURCE_USER
- DICT_DATASOURCE_PASSWORD

...and more can be found in main.go

Default service connects to database on localhost:5432

## How it works

Do not forget to set environment variables.

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
  "key": "newP",
  "name": "New parent",
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
GET /api/dictionary/AbsenceType/newP
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
  "key": "newP",
  "name": "New parent",
  "needConfirmationDocumentNumber": false,
  "needDeliveryDateConfirmation": true,
  "onlyOnBeginOrEnd": true,
  "tenant": "default",
  "type": "AbsenceType"
}
```

#### Load child only dictionary entry

```
GET /api/dictionary/AbsenceType/newCh1
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
  "key": "newP",
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
  "key": "newP",
  "name": "New parent updated",
  "needConfirmationDocumentNumber": false,
  "needDeliveryDateConfirmation": false,
  "onlyOnBeginOrEnd": true,
  "tenant": "default",
  "type": "AbsenceType"
}
```
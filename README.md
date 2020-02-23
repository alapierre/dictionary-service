# dictionary-service

Project is in havy development status!

fast and simple dictionary service on PostgreSQL database and JSON content

## Main goals

Almost any system needs to store and manage flexible dictionary values. Some of these need variable over time configuration. 

- provide fast REST microservice for stoing dictionary in JSON format 
- store configuration values changable over time (eg. what will be configuration value of X in 2025-12-31?)
- dictionaries shoud be describable - for automatic GUI build
- i18n support for dictionary names
- integrate with Netflix Eureka and oAtuh
- no additional, other than PostgreSQ database, for store dictionarys and configuration data - in cloud operators noSQL database is additional cost

## Reqired envirnoment variables

- DICT_DATASOURCE_USER
- DICT_DATASOURCE_PASSWORD

...and more can be found in main.go

Default service connects to database on localhost:5432

## How it works

You need put somethind into database first. Try SQL from testdata.
Do not forget to set envirnoment variables.


### Loading stored dictionary 

`GET /api/dictionary/{type}/{key}`

```
{
"children": [
 {
   "key": "uw1",
   "label": "ala",
   "name": "urlop 1",
   "type": "AbsenceType"
 },
 {
   "key": "uw2",
   "label": "ola",
   "name": "urlop 2",
   "type": "AbsenceType"
 }
],
"key": "uw",
"name": "Urlop wypoczynkowy",
"needConfirmationDocumentNumber": false,
"needDeliveryDateConfirmation": true,
"onlyOnBeginOrEnd": true,
"tenant": "",
"type": "AbsenceType"
}
```

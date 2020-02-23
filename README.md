# dictionary-service

Project is in havy development status!

fast and simple dictionary service on PostgreSQL database and JSON content

## Main goals

Almost any system needs to store and manage flexible dictionary values. Some of these need variable over time configuration. 

- provide fast REST microservice for stoing dictionary in JSON format and configuration values
- dictionaries shoud be describable - for automatic GUI build
- i18n support for dictionary names
- integrate with Netflix Eureka and oAtuh
- no additional, other than PostgreSQ database for store dictionarys and configuration data - in cloud operators it is additional cost

## How it works

You need put somethind into database first. Try SQL from testdata.

### Loading stored dictionary 

`GET /api/dictionary/{type}/{key}`


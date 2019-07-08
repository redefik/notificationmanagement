# notificationmanagement
Il microservizio si occupa di inviare avvisi alla mailing list di un corso, leggendoli da una coda SQS sulla quale vengono pubblicati dal microservizio [Course Management](https://github.com/tommasoVilla/Course_Management_Microservice).
Oltre a ciò, il microservizio espone un'interfaccia REST per la creazione e la gestione delle mailing list. Gli endpoint dell'interfaccia sono documentati su [SwaggerHub](https://app.swaggerhub.com/apis-docs/tommasoVilla/NotificationManagement/1.1) e nella cartella [api](api/).

## Linguaggio
Go

## Strato di persistenza
DynamoDB

# mail-api
API для получения архива писем по уникальной ссылке

## Особенности

* Статическая авторизация (basic auth) - username и пароль
* Данные пользователей указываются в JSON, в программе пароли хешируются
* Стандартный PDF readme файл - в папке ./static
* По умолчанию make выполняет полный цикл сборки - тесты, build, запуск. Все описано в Makefile

## Запуск

    ./bin/mailapi

## Справка по ключам

    ./bin/mailapi -h
    Usage of ./bin/mailapi:
    -config string
        path to JSON file for server configuration (default "configs/mailapi.json")
    -pdf string
        path to PDF file (default "static/readme.pdf")
    
### Информация о настройке конфига приведена в configs/README.md

## Работа с API
### POST /*username
#### Требуется передать имя пользователя и пароль

### Пример:

Request

    curl -X POST --user "beck.cierra@wb.ru:beck" localhost:8080/beck.cierra@wb.ru
Response - link

    /get/ae039c01-5076-4068-a62d-36adde319d29

### GET /get/*link

### Пример:

Request

    curl --output beck.cierra@wb.ru.zip -X GET localhost:8080/get/ae039c01-5076-4068-a62d-36adde319d29

Response

    ZIP архив с почтой + PDF

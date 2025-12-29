# Название проекта
Реализация pet project для онлайн библиотеки написанный RESTful API, имеет функционал регистрации , авторизации , JWT токенов и API-ручки для добавления книг, изменений, удаления и просмотра, база данных является PostgreSQL.
Дальше реализовано в виде gRPC-сервера логгер, который записывает все происходящие событие и записывает в NoSQL Mongo.

## Содержание
- [Технологии](#технологии)
- [Как запустить](#как-запустить)
- [Тестирование](#тестирование)


## Технологии
- [Echo](https://echo.labstack.com/)
- [Golang](https://go.dev/doc/install)
- [PostgreSQL](https://www.postgresql.org/)
- [Mongo](https://www.mongodb.com/)

## Как запустить
Сначала скачиваем [log_grpc сервер](https://github.com/CryptoGu1/books-grpc-log) и [books-rest-clean-arch](https://github.com/CryptoGu1/books-rest-clean-arch) командами:

```sh
$ git clone https://github.com/CryptoGu1/books-grpc-log
$ git clone https://github.com/CryptoGu1/books-rest-clean-arch
```

Начинаем запуск проекта с gRPC_log server, скачиваем все необходимые зависимости и поднимаем docker-ом compose весь проект и базу данных в докере:

```sh
$ cd books-grpc-log
$ go mod download
$ docker-compose up -d --build
```

Дальше запускаем основной REST-API сервис, в котором и реализован gRPC-клиент, который посылает запросы на gRPC-сервер:

```sh
$ cd books-rest-clean-arch
$ go mod download
$ make run
```

чтобы остановить и удалить контейнеры , вводим: 
```sh
$ docker-compose down -v --remove-orphans //удалит контейнеры и базы данных для gRPC_log server
$ make clean // тоже самое но для реста, или же можно просто make stop для просто остановки
```

## Разработка

### Требования
Для установки и запуска проекта, необходим [Golang](https://go.dev/doc/install)
 v25.0+. 

И [Docker Dekstop](https://www.docker.com/products/docker-desktop/) для запуска всего проекта в контейнерах.

## Тестирование
Наш проект пока что покрыт юнит тестами для одного хендлера , но в будущем планируются добить покрытие до 70%+ на хендлеры, сервисы и написать пару интеграционных тестирований.



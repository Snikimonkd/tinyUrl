# tinyUrl

## Задание (Стажер-разработчик)

Укорачиватель ссылок

Необходимо реализовать сервис, который должен предоставлять API по созданию сокращенных ссылок следующего формата:

 - Ссылка должна быть уникальной и на один оригинальный URL должна ссылаться только одна сокращенная ссылка.
 - Ссылка должна быть длинной 10 символов.
 - Ссылка должна состоять из символов латинского алфавита в нижнем и верхнем регистре, цифр и символа _ (подчеркивание).
 - Сервис должен быть написан на Go и принимать следующие запросы по gRPC:
1. Метод Create, который будет сохранять оригинальный URL в базе и возвращать сокращённый
2. Метод Get, который будет принимать сокращённый URL и возвращать оригинальный URL

Решение должно быть предоставлено в «конечном виде», а именно:
Сервис должен быть распространён в виде Docker-образа. Ожидается, что сервис позволяет использовать для хранения postgresql*. И in-memory хранилище, содержащее данные в памяти сервиса (т.е. Redis или любое другое внешнее хранилище не подойдет). Какое хранилище использовать, указывается параметром при запуске сервиса. API должно быть описано в proto файле.

Покрыть реализованный функционал Unit-тестами

Результат предоставить в виде публичного репозитория на http://github.com

* предпочтительней использовать postgresql

# Иснтрукция по запуску при помощи скрипта:

Запустить программу с sql базой данных:

```bash
./run.sh
```

Запустить программу с in memory базой данных:

```bash
./run.sh -n
```

# Как запустить самому:

Сборка с sql базой данных:

```bash
sudo docker-compose build
sudo docker-compose up -d
```

Сборка с in memory хранилищем:

```bash
sudo docker build -t tinyurl .
sudo docker run -d -p 5000:5000 tinyurl
```

# Команды для генерации файлов

Генерация файлов протобафа:

``` bash
protoc --go_out=/tinyUrl/internal/pkg/tinyUrl/delivery/server --go_opt=paths=source_relative --go-grpc_out=/tinyUrl/internal/pkg/tinyUrl/delivery/server --go-grpc_opt=paths=source_relative /tinyUrl/internal/pkg/tinyUrl/delivery/server/proto/server.proto --proto_path=/tinyUrl/internal/pkg/tinyUrl/delivery/server/proto
```

Генерация моков:

```bash
mockgen -source=./internal/pkg/tinyUrl/usecase/tinyUrl.go -destination=./internal/pkg/tinyUrl/usecase/mocks/tinyUrl_mock.go
```

Генерация html-файла с покрытием:

```bash
go test ./... -v -coverpkg=./... -coverprofile=cover.out.tmp && cat cover.out.tmp | grep -v "mock.go" | grep -v "pb.go" > cover.out && go tool cover -func=cover.out && go tool cover -html=cover.out
```
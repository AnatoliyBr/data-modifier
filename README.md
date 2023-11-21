## Data Modifier
Data Modifier - gRPC-сервер для расширения/модификации данных о пользователе на основе сторонней системы.

## Decisions
Чтобы организовать очередь запросов, можно использовать буферизированный канал, как показано [здесь](https://eli.thegreenplace.net/2019/on-concurrency-in-go-http-servers/). Либо можно ограничить число подключений на уровне слушателя с помощью [netutil.LimitListener](https://pkg.go.dev/golang.org/x/net/netutil#LimitListener).

Чтобы организовать воркер-пул, можно использовать цикл горутин, считывающих задачи из общего канала, как показано [здесь](https://gobyexample.com/worker-pools). Либо можно задать количество воркеров для gRPC сервера из пакета [grpc](https://pkg.go.dev/google.golang.org/grpc) с помощью экспериментальной функции [NumStreamWorkers](https://pkg.go.dev/google.golang.org/grpc#NumStreamWorkers).

## Примеры запросов
Тестирование осуществлялось при помощи Postman.

0. Эхо-версия AddAbsenceStatus
1. 

### Эхо-версия AddAbsenceStatus
Сначала реализация метода AddAbsenceStatus не изменяла данные, а просто возвращала их.

```json
{
    "user_data": {
        "id": 123,
        "display_name": "Иванов Семен Петрович",
        "email": "petrovich@mail.ru",
        "mobile_phone": "+71234567890",
        "work_phone": "1234"
        
    },
    "time_period": {
        "date_from": "2022-07-01T00:00:00Z",
        "date_to": "2022-09-01T23:59:59Z"
    }
}
```

Пример ответа:

```json
{
    "modified_user_data": {
        "id": 123,
        "display_name": "Иванов Семен Петрович",
        "email": "petrovich@mail.ru",
        "mobile_phone": "+71234567890",
        "work_phone": "1234"
    }
}
```

## Используемые команды
```bash
go get github.com/BurntSushi/toml

go get google.golang.org/grpc
go get google.golang.org/protobuf

go get golang.org/x/net

curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.55.2

go get google.golang.org/grpc
go get google.golang.org/protobuf

protoc --proto_path=api/proto --go_out=pkg --go-grpc_out=pkg api/proto/adder.proto

go get go.uber.org/zap
go get github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery
go get github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging

go get github.com/go-ozzo/ozzo-validation
go get github.com/go-ozzo/ozzo-validation/is

go get github.com/stretchr/testify
```

## Полезные ссылки
* [Пишем gRPC сервис на Go — сервис авторизации](https://habr.com/ru/articles/774796/)
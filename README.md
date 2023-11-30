# Data Modifier
Data Modifier - gRPC-сервер для расширения/модификации данных о пользователе на основе сторонней системы.

Используемые технологии:
* toml, godotenv (конфигурация)
* ozzo-validation (валидация)
* go.uber.org/zap (логирование)
* google.golang.org/grpc (gRPC сервер)
* grpc-ecosystem/go-grpc-middleware/v2/interceptors (gRPC интерсепторы)
* testify, go.uber.org/mock/gomock (unit-тесты)

Код организован согласно Clean Architecture, инъекция зависимостей обеспечивает низкую связанность слоев и упрощает тестирование.

Также реализован Graceful Shutdown для корректного завершения работы сервера. 

# Getting Started
Перед запуском приложения необходимо добавить в директорию проекта `.env` файл с логином и паролем от сторонней системы, как в примере [.env.example](.env.example).

А также настроить конфигурацию приложения через `.toml` файл, как в примере [config.toml](configs/config.toml). Путь к этому файлу можно передать с флагом `-config-path` при запуске приложения.

Если требуется записать логи и ошибки в файлы, необходимо создать эти файлы и передавать путь к ним в опциях `output_path`/`error_output_path` (подробнее см. ниже).

<details>
    <summary> Описание всех опций с примерами: </summary>

* `tcp_ip`/`tcp_port` - IP-адрес и порт, на который gRPC сервер будет принимать запросы . Заметьте, что порт указывается с `:`.
* `max_clients` - максимальное количество подключений (размер очереди запросов). 
* `num_pool_workers` - количество горутин для обработки запросов (размер пула обработчиков).

Пример конфигурации gRPC сервера:
```toml
tcp_ip="127.0.0.1"
tcp_port=":8081"
max_clients=10
num_pool_workers=5
```

* `web_api_ip`/`web_api_port` - IP-адрес и порт сторонней системы. Заметьте, что порт указывается с `:`.
* `protocol_type` - HTTP протокол ("http", "https").
* `employee_path` - путь для формирования запроса на получение данных о пользователе со сторонней системы. Заметьте, что путь указывается **без** `/`.
* `absence_path` - путь для формирования запроса на получение данных о статусе отсутствия пользователя со сторонней системы. Заметьте, что путь указывается **без** `/`.

Пример конфигурации для работы со сторонней системой:
```toml
web_api_ip="127.0.0.1"
web_api_port=":8082"
protocol_type="http"
employee_path="employees"
absence_path="absences"
```

* `log_format` - формат кодировки логов ("json", "console").

Пример лога в формате "json" и "console":
```bash
{"L":"INFO","T":"2023-11-29T15:24:59.674+0300","C":"app/app.go:61","M":"Initializing UserWebAPI..."}

2023-11-29T13:22:41.524+0300    INFO    app/app.go:61   Initializing UserWebAPI...
```

* `log_level` - уровень логирования ("debug", "info", "warn", "error"). Заметьте, что если поднять `log_level` до уровня "error", то логер пропустит сообщение уровня "info", так как все сообщения ниже установленного уровня игнорируются.
* `encoder_type` - определяет формат ключей у полей логов ("dev", "prod"), если установлен `log_format`="json". Описание для [dev](https://pkg.go.dev/go.uber.org/zap#NewDevelopmentEncoderConfig) и [prod](https://pkg.go.dev/go.uber.org/zap#NewProductionEncoderConfig).

Пример лога типа "dev" и "prod":
```bash
{"L":"INFO","T":"2023-11-29T16:36:59.832+0300","C":"app/app.go:62","M":"Initializing UserWebAPI..."}

{"level":"info","ts":1701264983.1611516,"caller":"app/app.go:62","msg":"Initializing UserWebAPI..."}
```

* `output_path`/`error_output_path` - **списки** путей к файлам для записи выходных логов и ошибок.

Пример конфигурации логера:
```toml
log_format="console"
log_level="debug"
encoder_type="dev"
output_path=[ "stdout", "./tmp/logs/rpc_traffic.txt" ]
error_output_path=[ "stderr", "./tmp/logs/rpc_traffic_errors.txt" ]
```

</details>

# Usage
Для сборки проекта необходимо выполнить команду `make`.

Для запуска сервиса необходимо выполнить команду `make up`. Для запуска сервиса с тестовым сервером - `make test-up`. 

Для запуска unit-тестов необходимо выполнить команду `make test`.

Для запуска линтера необходимо выполнить команду `make linter`.

## Examples
Для тестирования сервиса использовался **Postman** и **тестовый сервер** для имитации сторонней системы, который запускается с помощью флага `-test-server` (именно это делается при выполнении команды `make test-up`).

По заданию требовалось:

> ... по входящему gRPC запросу с информацией о пользователе по email найти того же пользователя на внешнем HTTP сервере, и обогатить входящее имя пользователя (ФИО) статусом отсутствия, если оно есть (в конец ФИО дописать emoji с соответствующим статусом).

Был описан интерфейс сервиса `DataModifier` с Unary RPC методом `AddAbsenceStatus`.

### Добавление emoji, соответствующего статусу отсутствия, в поле display_name
Для добавления **emoji**, соответствующего статусу отсутствия, в поле `display_name` необходимо передать информацию о пользователе (все поля обязательные и валидируются системой) и интервал времени (формат времени по умолчанию "2006-01-02T15:04:05"):

```json
{
    "user_data": {
        "display_name": "Иванов Семен Петрович",
        "email": "petrovich@mail.ru",
        "mobile_phone": "+71234567890",
        "work_phone": "1234"
        
    },
    "time_period": {
        "date_from": "2022-07-01T00:00:00",
        "date_to": "2022-09-01T23:59:59"
    }
}
```

Пример ответа:

```json
{
    "modified_user_data": {
        "display_name": "Иванов Семен Петрович 🏠",
        "email": "petrovich@mail.ru",
        "mobile_phone": "+71234567890",
        "work_phone": "1234"
    }
}
```

**Правила валидации** для каждой сущности определены в соответствующих методах (метод [Validate](/internal/entity/user.go) для структуры пользователя).

* `display_name` - обязательное поле, должно состоять только из букв unicode, может содержать составные имена/фамилии (например, Иванов-Сидоров).
* `email` - обязательное поле, должно быть валидной электронной почтой.
* `mobile_phone` - обязательное поле, должно состоять только из символов [0-9], может начинаться на "+", длина от 10 до 12 символов.
* `work_phone` - обязательное поле, должно состоять только из символов [0-9], может начинаться на "+", длина от 1 до 12 символов.

# Decisions
Чтобы организовать очередь запросов, можно использовать буферизированный канал, как показано [здесь](https://eli.thegreenplace.net/2019/on-concurrency-in-go-http-servers/). Либо можно ограничить число подключений на уровне слушателя с помощью [netutil.LimitListener](https://pkg.go.dev/golang.org/x/net/netutil#LimitListener).

Чтобы организовать воркер-пул, можно использовать цикл горутин, считывающих задачи из общего канала, как показано [здесь](https://gobyexample.com/worker-pools). Либо можно задать количество воркеров для gRPC сервера из пакета [grpc](https://pkg.go.dev/google.golang.org/grpc) с помощью экспериментальной функции [NumStreamWorkers](https://pkg.go.dev/google.golang.org/grpc#NumStreamWorkers).

<details>
<summary> Подходит для Unary RPC запросов (под капотом в исходном коде) </summary>

Заходим в [grpc-go/server.go](https://github.com/grpc/grpc-go/blob/v1.59.0/server.go). Находим функцию `Serve`, которая на каждое новое соединение запускает горутину `handleRowConn`. `handleRowConn` запускает горутину, в которой последовательно вызываются методы `serveStreams` и `removeConn`. В `serveStreams` вызывается `HandleStreams`, в которой в канал `s.serverWorkerChannel` передается функция с методом `handleStream`.

Далее в функции `serverWorker` эта функция считывается из канала и вызывается. Если сервис и метод существуют, `handleStream` вызывает метод `processUnaryRPC` или `processStreamingRPC`.

</details>

# Commands
<details>
<summary> Использованные команды</summary>

```bash
go get github.com/BurntSushi/toml

go get google.golang.org/grpc
go get google.golang.org/protobuf

go get golang.org/x/net

curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.55.2

protoc --proto_path=api/proto --go_out=pkg --go-grpc_out=pkg api/proto/datamodifier.proto

go get go.uber.org/zap
go get github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery
go get github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging

go get github.com/go-ozzo/ozzo-validation
go get github.com/go-ozzo/ozzo-validation/is

go get github.com/stretchr/testify

go install go.uber.org/mock/mockgen
go get go.uber.org/mock/mockgen

mockgen -source=./internal/webapi/interfaces.go -destination=./internal/webapi/webapi_mocks.go -package=webapi

go get github.com/gorilla/mux
```

</details>

## Полезные ссылки
* [Документация gRPC](https://grpc.io/) - фреймворк для удаленного вызова процедур, разработанный компанией Google.
    * [Инструкция Go Quick Start](https://grpc.io/docs/languages/go/quickstart/).
* [Protocol Buffers](https://protobuf.dev/) - инструмент сериализации (перевода в поток байтов) и язык описания gRPC, для начала работы необходимо установить компилятор.
    * [Компилятор](https://github.com/protocolbuffers/protobuf/releases/latest) Protocol Buffers, чтобы компилировать из proto-файлов, которые описывают интерфейс сервиса, код на каком-то языке.
    * [Инструкция Protocol Buffer Basics: Go](https://protobuf.dev/getting-started/gotutorial/)
* [Пишем gRPC сервис на Go — сервис авторизации](https://habr.com/ru/articles/774796/)
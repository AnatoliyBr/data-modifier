## Data Modifier
Data Modifier - gRPC-сервер для расширения/модификации данных о пользователе на основе сторонней системы.

## Decisions
Чтобы организовать очередь запросов, можно использовать буферизированный канал, как показано [здесь](https://eli.thegreenplace.net/2019/on-concurrency-in-go-http-servers/). Либо можно ограничить число подключений на уровне слушателя с помощью [netutil.LimitListener](https://pkg.go.dev/golang.org/x/net/netutil#LimitListener).

Чтобы организовать воркер-пул, можно использовать цикл горутин, считывающих задачи из общего канала, как показано [здесь](https://gobyexample.com/worker-pools). Либо можно задать количество воркеров для gRPC сервера из пакета [grpc](https://pkg.go.dev/google.golang.org/grpc) с помощью экспериментальной функции [NumStreamWorkers](https://pkg.go.dev/google.golang.org/grpc#NumStreamWorkers).

```bash
go get github.com/BurntSushi/toml

go get google.golang.org/grpc
go get google.golang.org/protobuf

go get golang.org/x/net
```
## 1. 다중 인스턴스를 위한 Redis Pub/Sub 도입
- 현재 채팅 서버 구조는 단일 인스턴스를 가장하고 제작된거기에 인스턴스 서버가 확장될 경우 각 A, B서버는
각각의 서버 Room Register 상태를 모르기에 메세지를 실시간으로 주고 받을 수 없다.
- Redis를 도입하여 다중 서버 인스턴스를 단일 Redis에 구독함으로써 서버가 확장되어도 실시간으로 브로드 캐스팅이 가능하도록 확장

## 2. docker-compose.yml Redis 추가
```yml
redis:
    image: redis
    container_name: gochat-redis
    restart: unless-stopped
    ports:
      - "6379:6379"
```

## 2. broker 패키지에 broker 인터페이스 추가
- 일반 Redis Broker만 사용하는게 아닌 broker로 추상하여 redis가 아닌 Pub/Sub Broker를 사용할 수 있도록 진행
```go
type Broker interface {
	Publish(ctx context.Context, roomName string, payload []byte) error
	Subscribe(ctx context.Context, roomName string) (<-chan []byte, func(), error)
}
```
- Publish : 해당 방 채널에 메세지 발행
- Subscribe : 해당 방 채널 구독 및 메세지가 올 때마다 채널 송신
  - Subscrbie 반환값에서 <-chan(수신전용)으로 사용해야하는 이유는 redis Broker 구현체가 받기만 하도록 되어있다
  - func()은 연결이 종료되었을 시 해당 커넥션을 종료할 때 구독을 해체해야하기에 상속 구현체에서 클로저로 구현해야한다.

## 3. 
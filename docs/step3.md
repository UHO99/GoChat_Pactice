## 1. JSON 포맷팅
- 기존의 에코 서버에서는 단순 텍스트 바이트 스트림을 받아 그대로 출력을 해주었습니다.
- messages/message.go 추가
```go
type Incoming struct {
	Content string `json:"content"`
}

type Outgoing struct {
	Nickname string    `json:"nickname,omitempty"`
	Content  string    `json:"content,omitempty"`
	SentAt   time.Time `json:"sentAt,omitempty"`
}
```
- 송/수신별로 JSON 객체를 제작
```go
		// messages.Incoming으로 Unmarshal-역직렬화
		var in messages.Incoming
		if err := json.Unmarshal(data, &in); err != nil {
			log.Fatal("역직렬화 실패", err)
			continue
		}

		// messages.Outgoing으로 marshal-직렬화
		out := messages.Outgoing{
			Nickname: nickname,
			Content:  in.Content,
			SentAt:   time.Now(),
		}
		payload, err := json.Marshal(out)
		if err != nil {
			log.Fatal("직렬화 실패", err)
			continue
		}

		log.Printf("[%s] %s", nickname, in.Content)

		if err := s.hub.Broadcast(ctx, payload); err != nil {
			log.Fatal("broadcast error : ", err)
			return
		}
	}
```
- handleWS에서 수신 데이터를 GO 객체로 역직렬화 후 해당 데이터에 time, nickname을 추가 후 JSON 객체로 직렬화하여 브로드 캐스팅

## 2. 결과
```bash
서버
2026/07/28 17:30:47 [lee] Hello
2026/07/28 17:31:11 [uho] Hello nice meet u

클라이언트 1
> {"content":"Hello"}
< {"nickname":"lee","content":"Hello","sentAt":"2026-07-28T17:30:47.836217241+09:00"}
< {"nickname":"uho","content":"Hello nice meet u","sentAt":"2026-07-28T17:31:11.78539828+09:00"}

클라이언트 2
< {"nickname":"lee","content":"Hello","sentAt":"2026-07-28T17:30:47.836217241+09:00"}
> {"content":"Hello nice meet u"}
< {"nickname":"uho","content":"Hello nice meet u","sentAt":"2026-07-28T17:31:11.78539828+09:00"}
```
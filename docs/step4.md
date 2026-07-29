## 1. Hub에서 단독 서버 입장 -> ROOM 입장
- Hub는 이제 클라이언트를 직접 들고 있지않고 Room을 관리하도록 변경했습니다.
- Room이 없으면 만들고 있으면 재사용 lazy-creation
```go
room.go
type Room struct {
	name    string
	mu      sync.Mutex
	clients map[*websocket.Conn]struct{}
}

func newRoom(name string) *Room {
	return &Room{
		name:    name,
		clients: make(map[*websocket.Conn]struct{}),
	}
}

func (r *Room) Register(c *websocket.Conn) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.clients[c] = struct{}{}
}

func (r *Room) Unregister(c *websocket.Conn) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.clients, c)
}

func (r *Room) Broadcast(ctx context.Context, payload []byte) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for c := range r.clients {
		if err := c.Write(ctx, websocket.MessageText, payload); err != nil {
			return err
		}
	}

	return nil
}

hub.go
type Hub struct {
	mu    sync.Mutex
	rooms map[string]*Room
}

func NewHub() *Hub {
	return &Hub{rooms: make(map[string]*Room)}
}

func (h *Hub) Room(name string) *Room {
	h.mu.Lock()
	defer h.mu.Unlock()

	r, ok := h.rooms[name]
	if !ok {
		r = newRoom(name)
		h.rooms[name] = r
	}

	return r
}
```

- handleWS 수정
```go
func (s *Server) handleWS(w http.ResponseWriter, r *http.Request) {
	nickname := r.URL.Query().Get("nickname")
	roomName := r.URL.Query().Get("room")
	if roomName == "" {
		roomName = "general"
	}

	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		return
	}
	defer conn.CloseNow()

	room := s.hub.Room(roomName)

	room.Register(conn)
	defer room.Unregister(conn)

	ctx := r.Context()
	for {
		_, data, err := conn.Read(ctx)
		if err != nil {
			log.Fatal("Cannot Connection Failed : ", err)
			return
		}

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

		if err := room.Broadcast(ctx, payload); err != nil {
			log.Fatal("broadcast error : ", err)
			return
		}
	}
```

## 2. 결과
```bash
서버
2026/07/28 19:29:55 방 생성 : lobbi
2026/07/28 19:29:55 방 입장 : lobbi
2026/07/28 19:29:58 방 생성 : lobby
2026/07/28 19:29:58 방 입장 : lobby
2026/07/28 19:30:09 방 입장 : lobby
2026/07/28 19:31:07 [lee] hello
2026/07/28 19:31:21 [uho] GOOD!
2026/07/28 19:31:37 [kim] Any body here?

클라이언트 1
< {"nickname":"lee","content":"hello","sentAt":"2026-07-28T19:31:07.768238752+09:00"}
> {"content":"GOOD!"}
< {"nickname":"uho","content":"GOOD!","sentAt":"2026-07-28T19:31:21.348480829+09:00"}

클라이언트 2
> {"content":"hello"}
< {"nickname":"lee","content":"hello","sentAt":"2026-07-28T19:31:07.768238752+09:00"}
< {"nickname":"uho","content":"GOOD!","sentAt":"2026-07-28T19:31:21.348480829+09:00"}

클라이언트 3
> {"content":"Any body here?"}
< {"nickname":"kim","content":"Any body here?","sentAt":"2026-07-28T19:31:37.284403956+09:00"}
```
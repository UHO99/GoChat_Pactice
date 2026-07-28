## 1. 다중 사용자로 확장
- 현재는 handleWS에서 하나의 서버에서 한명의 사용자가 송/수신하는 에코 서버입니다.
- servers/hub.go 추가
```go
type Hub struct {
	mu      sync.Mutex
	clients map[*websocket.Conn]struct{}
}

func NewHub() *Hub {
	return &Hub{clients: make(map[*websocket.Conn]struct{})}
}

func (h *Hub) Register(c *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[c] = struct{}{}
}

func (h *Hub) Unregister(c *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.clients, c)
}

func (h *Hub) Broadcast(ctx context.Context, payload []byte) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	for c := range h.clients {
		if err := c.Write(context.Background(), websocket.MessageText, payload); err != nil {
			return err
		}
	}

	return nil
}
```
- Hub를 제작해 clients를 제작 clients는 websocket에 커넥션된 client, 구조체를 가지는 map으로 각 client를 돌며 write를 진행 후 해당 write 도중 error 발생시 error 반환
- register, unregister를 제작해 connection별로 hub를 제작하는게 아닌 해당 hub하나에 다중 사용자를 커넥션시켜서 브로드캐스트시 해당 Hub 내의 사용자들에게 브로드캐스트 진행

## 2. 핸들러 변경
```go
func (s *Server) handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		return
	}
	defer conn.CloseNow()

	s.hub.Register(conn)
	defer s.hub.Unregister(conn)

	ctx := r.Context()
	for {
		_, data, err := conn.Read(ctx)
		if err != nil {
			return
		}
		log.Printf("[%s] %s", r.RemoteAddr, string(data))

		if err := s.hub.Broadcast(ctx, data); err != nil {
			log.Fatal("broadcast error : ", err)
			return
		}
	}
}
```
- 핸들러에 단일 사용자 작성을 hub 구조체를 통해 각 연결된 client별로 개별 송신
- Accept시 Hub에 Register 해당 커넥션 종료시 Unregister

## 3. 서버 변경
```go
type Options struct {
	Addr string
}

type Server struct {
	httpServer *http.Server
	hub        *servers.Hub
}

func New(opts Options) *Server {
	s := &Server{hub: servers.NewHub()}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /ws", s.handleWS)

	s.httpServer = &http.Server{
		Addr:    opts.Addr,
		Handler: mux,
	}

	return s
}

func (s *Server) ListenAndServe() error {
	return s.httpServer.ListenAndServe()
}
```
- 서버 코드 또한 hub에 맞춰서 각 서버에 Hub를 개별로 생성자 초기화

```
서버
2026/07/28 16:46:20 [127.0.0.1:36646] asd
2026/07/28 16:46:23 [127.0.0.1:58834] qwe
2026/07/28 16:46:31 [127.0.0.1:36646] wwww
2026/07/28 16:46:35 [127.0.0.1:36646] sdf

클라이언트 1
< asd
> qwe
< qwe
< wwww
< sdf

클라이언트 2
> asd
< asd
< qwe
> wwww
< wwww
> sdf
< sdf
```
## 1. 사용된 외부 라이브러리
```
websocket
viper
```

## 2. 환경변수 체크
- app.env 파일의 PORT 환경변수를 config.go에서 viper로 읽는다.
```go
type Config struct {
	Port string `mapstructure:"PORT"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
```

## 3. handleWS 작성
```go
func (s *Server) handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		return
	}
	defer conn.CloseNow()

	ctx := r.Context()
	for {
		_, data, err := conn.Read(ctx)
		if err != nil {
			return
		}
		log.Println(string(data))

		if err := conn.Write(ctx, websocket.MessageText, data); err != nil {
			return
		}
	}
}
```
1. websocket.Accept : HTTP 요청을 WebSocket으로 업그레이드. 클라이언트가 보낸 `Upgrade: websocket`, `Sec-WebSocket-Key` 헤더를 검증, 101 Switching Protocols 응답을 줘서 해당 시점부터 TCP 연결이 HTTP가 아닌 WebSocket 프레임 프로토콜로 전환
2. conn.Read(ctx) - 메세지 하나가 도착할때까지 블로킹, 반환값 3개중 첫번째 _ 는 `websocket.MessageType`으로 텍스트인지 바이너리인지 확인 가능, `ctx`는 `r.Context()` - 클라이언트가 연결을 종료하거나 셧다운 시 해당 컨텍스트 캔슬되며 `Read`에러 반환
3. conn.Write(...) - 읽은 `data`를 텍스트 메세지로 반환
루프와 goroutine - `for` 루프는 에러가 날때까지 계속 루프 `net/http`가 요청마다 별도 goroutine 실행, 때문에 해당 함수가 무한 루프를 돌아도 다른 클라이언트를 막지 않고 접속한 클라이언트 수만큼 `handleWS` goroutine이 실 행

## 4. server.go 작성
```go
type Options struct {
	Addr string
}

type Server struct {
	httpServer *http.Server
}

func New(opts Options) *Server {
	s := &Server{}

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

1. http.NewServeMux() : go 표준 라이브러리의 라우터(멀티 플렉서), 들러온 HTTP 요청 메서드 + 경로를 보고 핸들러 함수로 라우팅
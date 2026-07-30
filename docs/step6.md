## 1. DB 메세지, 방, 유저 영속성을 위한 스키마 생성
- migrate create -ext sql -dir db/migration -seq init_schema
```
db/migration/000001_init_schema.down.sql 
db/migration/000001_init_schema.up.sql
```

## 2. SQLC 사용
- https://docs.sqlc.dev/en/stable/tutorials/getting-started-postgresql.html
- 해당 URL sqlc 프레임워크의 사용법 명시
- `sqlc.yaml`로 sqlc 프레임워크 INIT
- 매개변수
	- `engine`: `postgresql`
	- `query` : `db/query` -> 해당 쿼리 폴더에 postgresql 문법을 혼합하여 쿼리 삽입
	- `schema`: `db/migration`
	- `out`   : `db/sqlc` -> 쿼리문을 토대로 sqlc가 go 문법으로 작성
- `sqlc.yaml` 설정 후 하단 명령어 실행
```bash
sqlc generate
```

## 3. store 작성
- `db/sqlc`에 room.sql.go, user.sql.go, message.sql.go에서 각각의 SQL문 존재
- 별도의 store.go를 제작해 각 SQL문 함수를 래핑 Java에서 Repository같은 존재
```go
type Store struct {
	query *db.Queries
}

type Message struct {
	ID       int64
	RoomID   int64
	Nickname string
	Content  string
	CreateAt time.Time
}

func New(pool *pgxpool.Pool) *Store {
	return &Store{query: db.New(pool)}
}

func (s *Store) CreateUser(ctx context.Context, nickname string) (db.User, error) {
	return s.query.CreateUser(ctx, nickname)
}

func (s *Store) GetUser(ctx context.Context, nickname string) (db.User, error) {
	return s.query.GetUser(ctx, nickname)
}

func (s *Store) CreateRoom(ctx context.Context, name string) (db.Room, error) {
	return s.query.CreateRoom(ctx, name)
}

func (s *Store) GetRoomByName(ctx context.Context, name string) (db.Room, error) {
	return s.query.GetRoomByName(ctx, name)
}

func (s *Store) InsertMessage(ctx context.Context, roomID, userID int64, content string) (db.Message, error) {
	return s.query.InsertMessage(ctx, db.InsertMessageParams{
		RoomID:  roomID,
		UserID:  userID,
		Content: content,
	})
}

func (s *Store) ListRooms(ctx context.Context) ([]db.Room, error) {
	return s.query.ListRooms(ctx)
}
```

## 4. http_handler.go 작성
- 현재 room/user/message 등의 DB persist를 별도 API로 분리
```go
type roomResponse struct {
	Name string `json:"name"`
}

type userResponse struct {
	Nickname string `json:"nickname"`
}

type createRoomRequest struct {
	Name string `json:"name"`
}

type createUserRequest struct {
	Nickname string `json:"name"`
}

const pgUniqueViolation = "23505"

func (s *Server) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var req createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Nickname == "" {
		http.Error(w, "Nickname parameter Error", http.StatusBadRequest)
		return
	}

	user, err := s.store.CreateUser(r.Context(), req.Nickname)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
			http.Error(w, "room already exists", http.StatusConflict)
			return
		}
		http.Error(w, "cannot create user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(userResponse{Nickname: user.Nickname})
}

func (s *Server) handleCreateRoom(w http.ResponseWriter, r *http.Request) {
	var req createRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
		http.Error(w, "room name required", http.StatusBadRequest)
		return
	}

	room, err := s.store.CreateRoom(r.Context(), req.Name)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
			http.Error(w, "room already exists", http.StatusConflict)
			return
		}
		http.Error(w, "cannot create room", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(roomResponse{Name: room.Name})
}

func (s *Server) handleListRooms(w http.ResponseWriter, r *http.Request) {
	rooms, err := s.store.ListRooms(r.Context())
	if err != nil {
		http.Error(w, "cannot list rooms", http.StatusInternalServerError)
		return
	}

	resp := make([]roomResponse, 0, len(rooms))
	for _, room := range rooms {
		resp = append(resp, roomResponse{Name: room.Name})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
```
- 분리 이유는 각 ws에서 WebSocket 핸드셰이크(Upgrade) 이전에 room/user를 조회 후 해당 방에 register를 하기 위해서

## 5. ws_handler.go 수정
```go
user, err := s.store.GetUser(ctx, nickname)
if err != nil {
	http.Error(w, "user get fail", http.StatusNotFound)
	return
}

room, err := s.store.GetRoomByName(ctx, roomName)
if err != nil {
	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "room does not exist", http.StatusNotFound)
		return
	}
	http.Error(w, "cannot look up room", http.StatusInternalServerError)
	return
}
```
- WebSocket 핸드셰이크(websocket.Accept)를 보내기 이전에 GetUser/Room을 진행하여 DB조회

## 5. 결과
```bash
서버
2026/07/30 17:38:06 방 생성 : general
2026/07/30 17:38:06 [ROOM : general / USER : lee] : {"type":"join","room":"general","nickname":"lee","sentAt":"2026-07-30T17:38:06.951956976+09:00"}
2026/07/30 17:38:25 [ROOM : general / USER : lee] : {"type":"join","room":"general","nickname":"uho","sentAt":"2026-07-30T17:38:25.997600835+09:00"}
2026/07/30 17:38:25 [ROOM : general / USER : uho] : {"type":"join","room":"general","nickname":"uho","sentAt":"2026-07-30T17:38:25.997600835+09:00"}

클라이언트 1
:~/gochat$ make createroom ROOMNAME=general
curl -X POST http://localhost:8086/rooms -d '{"name":"general"}'
{"name":"general"}
:~/gochat$ make listroom
curl -X GET http://localhost:8086/rooms
[{"name":"general"},{"name":"lobby"},{"name":"loby"},{"name":"world"},{"name":"worlds"}]
:~/gochat$ make createuser USERNAME=lee
curl -X POST http://localhost:8086/users -d '{"name":"lee"}'
{"nickname":"lee"}
:~/gochat$ make send SERVER=6 NICKNAME=lee ROOM=genera
wscat -c "ws://localhost:8086/ws?nickname=lee&room=genera"
error: Unexpected server response: 404
> make: *** [scripts/server.mk:18: send] Error 255
:~/gochat$ make send SERVER=6 NICKNAME=lee ROOM=general
wscat -c "ws://localhost:8086/ws?nickname=lee&room=general"
Connected (press CTRL+C to quit)
< {"type":"join","room":"general","nickname":"lee","sentAt":"2026-07-30T17:38:06.951956976+09:00"}
< {"type":"join","room":"general","nickname":"uho","sentAt":"2026-07-30T17:38:25.997600835+09:00"}

클라이언트 2
:~/gochat$ make send SERVER=6 NICKNAME=uho ROOM=general
wscat -c "ws://localhost:8086/ws?nickname=uho&room=general"
Connected (press CTRL+C to quit)
< {"type":"join","room":"general","nickname":"uho","sentAt":"2026-07-30T17:38:25.997600835+09:00"}
```
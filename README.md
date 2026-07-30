## 1. 각 STEP별 폴더
- /docs 에 각 STEP별 폴더 설명 포함

## 2. 실행방법
### 1. 사전 준비
- Go 1.26 이상
- (선택) wscat: WebSocket 클라이언트 테스트용 (`npm install -g wscat`)

### 2. 환경변수
- `app.env`에 각 STEP 서버가 사용할 포트가 정의되어 있다.
- `DB_URL`도 세팅
```
DATABASE_URL=postgres://root:1234@localhost:5432/gochat?sslmode=disable
STEP1_PORT=8080
STEP2_PORT=8081
```

### 3. 마이그레이션 파일 생성
- `up` `down` 마이그레이션 SQL 생성
```bash
migrate create -ext sql -dir [생성 폴더] -seq [생성 이름]
```

### 4. Postgres Docker Image 실행
```bash
make postgres
```

### 5. DB 스키마 생성 (선택)
- `docker-compose.yml`에서 enviroment로 POSTGRES_DB: gochat을 줄 경우 따로 스키마 생성할 필요 없음
```bash
make createdb
```

### 6. 테이블 생성
```bash
make migrateup
```

### 서버 실행
```bash
# STEP1 (echo 서버)
make server SERVER=1
# 또는
go run ./step1_simple_echo/cmd

# STEP2 (다중 사용자 브로드캐스트 서버)
make server SERVER=2
# 또는
go run ./step2_simple_multi_broadcast/cmd
```

### 클라이언트 접속 테스트
```bash
# STEP1 서버(8080) 접속
wscat -c ws://localhost:8080/ws
# 또는
make send SERVER=1

# STEP2 서버(8081) 접속 - 여러 터미널에서 접속하면 브로드캐스트 동작을 확인 가능
wscat -c ws://localhost:8081/ws
# 또는
make send SERVER=2
```

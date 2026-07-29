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
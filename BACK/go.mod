module BACK

go 1.24.4

// 세션 의존성 추가
require github.com/gorilla/sessions v1.2.2

require (
	github.com/go-sql-driver/mysql v1.9.3
	github.com/joho/godotenv v1.5.1
	golang.org/x/crypto v0.40.0
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/golang-jwt/jwt/v5 v5.2.3 // indirect
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/gorilla/securecookie v1.1.2 // indirect
)

//프로젝트의 의존성 관리 파일
//프로젝트 모듈 선언 + 모든 외부 패키지 버전 관리 파일

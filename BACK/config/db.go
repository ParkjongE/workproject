package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

// 전역 변수로 DB 연결 객체 선언
var DB *sql.DB

// InitDB는 MySQL 데이터베이스 연결을 초기화합니다
func InitDB() {
	// 현재 실행 경로 확인
	currentDir, err := os.Getwd()
	if err != nil {
		log.Println("현재 디렉토리 확인 실패:", err)
	} else {
		log.Printf("현재 작업 디렉토리: %s", currentDir)
	}

	// 여러 경로에서 .env 파일 로드 시도
	envPaths := []string{
		filepath.Join(currentDir, ".env"),
		filepath.Join(currentDir, "BACK", ".env"),
		filepath.Join(currentDir, "..", "BACK", ".env"),
		filepath.Join(currentDir, "..", "..", "BACK", ".env"),
		filepath.Join("C:", "project", ".env"),
		filepath.Join("C:", "project", "BACK", ".env"),
	}

	envLoaded := false
	for _, path := range envPaths {
		log.Printf(".env 파일 시도: %s", path)
		if _, err := os.Stat(path); err == nil {
			log.Printf(".env 파일 발견: %s", path)
			if err := godotenv.Load(path); err == nil {
				log.Printf("데이터베이스: .env 파일을 로드했습니다: %s", path)
				envLoaded = true
				break
			} else {
				log.Printf(".env 파일 로드 실패: %s - %v", path, err)
			}
		}
	}

	if !envLoaded {
		log.Println("경고: 데이터베이스 설정을 위한 .env 파일을 찾을 수 없습니다. 시스템 환경변수를 사용합니다.")
		// 환경변수 직접 설정
		os.Setenv("DB_HOST", "localhost")
		os.Setenv("DB_PORT", "3306")
		os.Setenv("DB_USER", "root")
		os.Setenv("DB_PASSWORD", "qkrwhdgml123")
		os.Setenv("DB_NAME", "project")
		log.Println("데이터베이스 환경변수를 직접 설정했습니다.")
	}

	// 환경변수에서 데이터베이스 연결 정보 가져오기
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "3306")
	dbUser := getEnv("DB_USER", "root")
	dbPassword := getEnv("DB_PASSWORD", "qkrwhdgml123")
	dbName := getEnv("DB_NAME", "project")

	// MySQL 연결 문자열 생성
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
	log.Printf("데이터베이스 연결 시도: %s@%s:%s/%s", dbUser, dbHost, dbPort, dbName)

	// 데이터베이스 연결
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("데이터베이스 연결 실패:", err)
	}

	// 연결 확인
	err = DB.Ping()
	if err != nil {
		log.Fatal("데이터베이스 ping 실패:", err)
	}

	log.Println("MySQL 데이터베이스 연결 성공")
}

// getEnv는 환경변수 값을 가져오고, 없으면 기본값을 반환합니다
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

//프로젝트 전역에서 쓸 MySQL DB 연결을 초기화하는 코드

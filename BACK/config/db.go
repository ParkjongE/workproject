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

// Database 데이터베이스 인터페이스
type Database interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}

// DB 전역 데이터베이스 인스턴스
var DB Database

// GetDB *sql.DB 타입으로 반환 (타입 단언)
func GetDB() *sql.DB {
	if db, ok := DB.(*sql.DB); ok {
		return db
	}
	return nil
}

// SetDB 데이터베이스 인스턴스 설정
func SetDB(db Database) {
	DB = db
}

// loadEnvFile .env 파일 로드
func loadEnvFile() error {
	// 현재 디렉토리에서 .env 파일 찾기
	envPath := ".env"

	// BACK 디렉토리에서 .env 파일 찾기
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		// 상위 디렉토리에서 .env 파일 찾기
		envPath = filepath.Join("..", ".env")
		if _, err := os.Stat(envPath); os.IsNotExist(err) {
			// 프로젝트 루트에서 .env 파일 찾기
			envPath = filepath.Join("..", "..", ".env")
		}
	}

	// .env 파일 로드
	if err := godotenv.Load(envPath); err != nil {
		log.Printf("경고: .env 파일을 찾을 수 없습니다 (%s): %v", envPath, err)
		// .env 파일이 없어도 계속 진행 (환경 변수 사용)
		return nil
	}

	log.Printf(".env 파일을 성공적으로 로드했습니다: %s", envPath)
	return nil
}

// InitDB 데이터베이스 초기화
func InitDB() error {
	// .env 파일 로드
	if err := loadEnvFile(); err != nil {
		log.Printf("경고: .env 파일 로드 실패: %v", err)
	}

	// 환경 변수에서 데이터베이스 연결 정보 가져오기
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "172.30.1.20" // 서버 컴퓨터 IP로 변경
	}

	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "3306"
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "root"
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = ""
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "hr_system"
	}

	log.Printf("데이터베이스 연결 정보: %s:%s@%s:%s/%s", dbUser, "***", dbHost, dbPort, dbName)

	// DSN (Data Source Name) 생성
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&parseTime=true&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	// 데이터베이스 연결
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("데이터베이스 연결 실패: %v", err)
	}

	// 연결 테스트
	if err := db.Ping(); err != nil {
		return fmt.Errorf("데이터베이스 연결 테스트 실패: %v", err)
	}

	// 전역 변수에 설정
	SetDB(db)

	log.Println("데이터베이스 연결 성공")
	return nil
}

// CloseDB 데이터베이스 연결 종료
func CloseDB() error {
	if DB != nil {
		if db, ok := DB.(*sql.DB); ok {
			return db.Close()
		}
	}
	return nil
}

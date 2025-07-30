-- 데이터베이스 생성 
CREATE DATABASE project;

-- 데이터베이스 선택
USE project;

-- 사용자 테이블 생성
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    userid VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    salt VARCHAR(64) NOT NULL,
    name VARCHAR(100) NOT NULL,
    role VARCHAR(20),
    department VARCHAR(50),
    hire_date DATE NOT NULL DEFAULT '2024-01-01',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
); 

-- 일정(스케줄) 테이블 생성
CREATE TABLE schedules (
    schedule_id INT AUTO_INCREMENT PRIMARY KEY, -- 일정번호
    user_id INT NOT NULL,                      -- 사원id
    title VARCHAR(200) NOT NULL,               -- 일정제목
    content TEXT,                              -- 일정 내용
    project_id INT,                            -- 프로젝트번호(옵션)
    start_time DATETIME NOT NULL,              -- 일정 시작 시간
    end_time DATETIME NOT NULL,                -- 일정 종료 시간
    color VARCHAR(20) DEFAULT '#4CAF50',       -- 일정 색상
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
); 

-- 출퇴근 기록 테이블 생성
CREATE TABLE attendance_records (
    record_id INT AUTO_INCREMENT PRIMARY KEY,  -- 기록번호
    user_id INT NOT NULL,                      -- 사원id
    work_date DATE NOT NULL,                   -- 출퇴근 날짜
    check_in_time DATETIME,                    -- 출근 시간
    check_out_time DATETIME,                   -- 퇴근 시간
    status VARCHAR(20),                        -- 상태 (정상, 지각, 조퇴 등)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 휴가 신청 테이블 생성
CREATE TABLE leave_requests (
    request_id INT AUTO_INCREMENT PRIMARY KEY,      -- 휴가신청 id
    user_id INT NOT NULL,                           -- 신청 사원번호
    leave_type VARCHAR(50) NOT NULL,                -- 휴가 유형 (예: 연차, 병가 등)
    start_date DATE NOT NULL,                       -- 휴가 시작일
    end_date DATE NOT NULL,                         -- 휴가 종료일
    status VARCHAR(20) NOT NULL,                    -- 신청 상태 (예: 대기, 승인, 반려 등)
    reason TEXT,                                    -- 신청 사유
    requested_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- 신청 일시
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 사용자별 연차 현황 테이블 생성
CREATE TABLE user_leave_status (
    status_id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,                          -- 사용자 ID
    total_earned_days INT NOT NULL DEFAULT 0,      -- 총 발생 연차 일수
    used_days INT NOT NULL DEFAULT 0,              -- 사용한 연차 일수
    remaining_days INT NOT NULL DEFAULT 0,         -- 잔여 연차 일수
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE KEY unique_user (user_id)
);
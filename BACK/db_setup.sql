-- 데이터베이스 생성 
CREATE DATABASE project;

-- 데이터베이스 선택
USE project;

-- users 테이블 생성
CREATE TABLE users (
    id INT PRIMARY KEY AUTO_INCREMENT,
    userid VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    salt VARCHAR(64) NOT NULL,
    name VARCHAR(100) NOT NULL,
    role VARCHAR(20),
    department VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    hire_date DATE DEFAULT '2024-01-01' NOT NULL
);

-- user_leave_status 테이블 생성
CREATE TABLE user_leave_status (
    status_id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT NOT NULL,
    total_earned_days INT DEFAULT 0 NOT NULL,
    used_days INT DEFAULT 0 NOT NULL,
    remaining_days INT DEFAULT 0 NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- schedules 테이블 생성
CREATE TABLE schedules (
    schedule_id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT NOT NULL,
    title VARCHAR(200) NOT NULL,
    content TEXT,
    project_id INT,
    start_time DATETIME NOT NULL,
    end_time DATETIME NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    color VARCHAR(20) DEFAULT '#4CAF50',
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- leave_requests 테이블 생성
CREATE TABLE leave_requests (
    request_id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT NOT NULL,
    leave_type VARCHAR(50) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    status VARCHAR(20) NOT NULL,
    reason TEXT,
    requested_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- attendance_records 테이블 생성
CREATE TABLE attendance_records (
    record_id INT PRIMARY KEY AUTO_INCREMENT, 
    work_date DATE NOT NULL,
    check_in_time DATETIME,
    check_out_time DATETIME, 
    status VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    UNIQUE KEY unique_user_date (user_id, work_date)
);

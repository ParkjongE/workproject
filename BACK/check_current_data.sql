-- 현재 데이터베이스 상태 확인

USE project;

-- 1. 현재 사용자 목록 확인
SELECT id, userid, name, department, hire_date FROM users;

-- 2. 연차 사용 내역 확인
SELECT 
    lr.user_id,
    u.name,
    lr.leave_type,
    lr.start_date,
    lr.end_date,
    lr.status,
    DATEDIFF(lr.end_date, lr.start_date) + 1 as days
FROM leave_requests lr
JOIN users u ON lr.user_id = u.id
WHERE lr.leave_type = '연차' AND lr.status = '승인'
ORDER BY lr.user_id;

-- 3. 현재 연차 현황 테이블 상태 확인
SELECT * FROM user_leave_status; 
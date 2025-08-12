package service

import (
	"BACK/config"
	"BACK/models"
	"fmt"
	"time"
)

// AttendanceServiceImpl 출근 서비스 구현체
type AttendanceServiceImpl struct {
	db config.Database
}

// NewAttendanceService 출근 서비스 생성자
func NewAttendanceService(db config.Database) *AttendanceServiceImpl {
	return &AttendanceServiceImpl{db: db}
}

// ==================== 서비스 구현체 메서드들 ====================

// GetAttendance 서비스 구현
func (s *AttendanceServiceImpl) GetAttendance(userID int, startDate, endDate time.Time) ([]models.Attendance, error) {
	// 실제 테이블명과 컬럼명 사용
	query := "SELECT record_id, user_id, work_date, check_in_time, check_out_time, status, created_at FROM attendance_records WHERE user_id = ? AND work_date BETWEEN ? AND ? ORDER BY work_date DESC"

	rows, err := s.db.Query(query, userID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []models.Attendance
	for rows.Next() {
		var attendance models.Attendance
		err := rows.Scan(
			&attendance.RecordID,     // record_id
			&attendance.UserID,       // user_id
			&attendance.WorkDate,     // work_date
			&attendance.CheckInTime,  // check_in_time
			&attendance.CheckOutTime, // check_out_time
			&attendance.Status,       // status
			&attendance.CreatedAt,    // created_at
		)
		if err != nil {
			return nil, err
		}
		history = append(history, attendance)
	}

	return history, nil
}

// PostAttendance 서비스 구현
func (s *AttendanceServiceImpl) PostAttendance(userID int, date time.Time, action string, timeStr string) error {
	fmt.Printf("[PostAttendance] %s 처리 - userID: %d, date: %s\n", action, userID, date.Format("2006-01-02"))

	var status string
	var query string
	var args []interface{}

	switch action {
	case "checkin":
		// 출근 시에는 status를 저장하지 않음 (NULL로 유지)
		// timeStr가 제공되면 파싱해서 사용, 아니면 현재 시간 사용
		var checkInTime time.Time
		if timeStr != "" {
			// ISO 형식인지 HH:MM 형식인지 확인
			if len(timeStr) > 10 && (timeStr[4] == '-' || timeStr[10] == 'T') {
				// ISO 형식 (2024-01-01T10:30:00.000Z)
				parsedTime, err := time.Parse(time.RFC3339, timeStr)
				if err != nil {
					return fmt.Errorf("invalid ISO time format: %s", timeStr)
				}
				checkInTime = parsedTime
			} else {
				// HH:MM 형식
				parsedTime, err := time.Parse("15:04", timeStr)
				if err != nil {
					return fmt.Errorf("invalid time format: %s", timeStr)
				}
				// 오늘 날짜와 결합
				checkInTime = time.Date(date.Year(), date.Month(), date.Day(),
					parsedTime.Hour(), parsedTime.Minute(), 0, 0, time.Local)
			}
		} else {
			checkInTime = time.Now()
		}

		query = `INSERT INTO attendance_records (user_id, work_date, check_in_time, created_at) 
				 VALUES (?, ?, ?, NOW()) 
				 ON DUPLICATE KEY UPDATE 
				 check_in_time = VALUES(check_in_time), 
				 created_at = NOW()`
		args = []interface{}{userID, date.Format("2006-01-02"), checkInTime}

	case "checkout":
		// 퇴근 시에는 status를 저장하지 않음 (NULL로 유지)
		// timeStr가 제공되면 파싱해서 사용, 아니면 현재 시간 사용
		var checkOutTime time.Time
		if timeStr != "" {
			// ISO 형식인지 HH:MM 형식인지 확인
			if len(timeStr) > 10 && (timeStr[4] == '-' || timeStr[10] == 'T') {
				// ISO 형식 (2024-01-01T10:30:00.000Z)
				parsedTime, err := time.Parse(time.RFC3339, timeStr)
				if err != nil {
					return fmt.Errorf("invalid ISO time format: %s", timeStr)
				}
				checkOutTime = parsedTime
			} else {
				// HH:MM 형식
				parsedTime, err := time.Parse("15:04", timeStr)
				if err != nil {
					return fmt.Errorf("invalid time format: %s", timeStr)
				}
				// 오늘 날짜와 결합
				checkOutTime = time.Date(date.Year(), date.Month(), date.Day(),
					parsedTime.Hour(), parsedTime.Minute(), 0, 0, time.Local)
			}
		} else {
			checkOutTime = time.Now()
		}

		query = `INSERT INTO attendance_records (user_id, work_date, check_out_time, created_at) 
				 VALUES (?, ?, ?, NOW()) 
				 ON DUPLICATE KEY UPDATE 
				 check_out_time = VALUES(check_out_time), 
				 created_at = NOW()`
		args = []interface{}{userID, date.Format("2006-01-02"), checkOutTime}

	case "record":
		status = "기록"
		query = `INSERT INTO attendance_records (user_id, work_date, status, created_at) 
				 VALUES (?, ?, ?, NOW()) 
				 ON DUPLICATE KEY UPDATE 
				 status = VALUES(status), 
				 created_at = NOW()`
		args = []interface{}{userID, date.Format("2006-01-02"), status}

	default:
		return fmt.Errorf("invalid action: %s", action)
	}

	result, err := s.db.Exec(query, args...)
	if err != nil {
		fmt.Printf("[PostAttendance] 실패: %v\n", err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	fmt.Printf("[PostAttendance] 성공 (영향받은 행: %d)\n", rowsAffected)

	return nil
}

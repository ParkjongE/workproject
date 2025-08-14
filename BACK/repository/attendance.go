package repository

import (
	"BACK/config"
	"BACK/model"
	"time"
)

// AttendanceRepository 인터페이스
type AttendanceRepository interface {
	FindByUserAndDateRange(userID int, startDate, endDate time.Time) ([]model.Attendance, error)
	CreateOrUpdate(userID int, date time.Time, action string, actionTime time.Time, status string) error
}

// AttendanceRepositoryImpl 구현체
type AttendanceRepositoryImpl struct {
	db config.Database
}

// NewAttendanceRepository 생성자
func NewAttendanceRepository(db config.Database) AttendanceRepository {
	return &AttendanceRepositoryImpl{db: db}
}

// FindByUserAndDateRange 날짜 범위로 출근 기록 조회
func (r *AttendanceRepositoryImpl) FindByUserAndDateRange(userID int, startDate, endDate time.Time) ([]model.Attendance, error) {
	query := "SELECT record_id, user_id, work_date, check_in_time, check_out_time, status, created_at FROM attendance_records WHERE user_id = ? AND work_date BETWEEN ? AND ? ORDER BY work_date DESC"

	rows, err := r.db.Query(query, userID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []model.Attendance
	for rows.Next() {
		var attendance model.Attendance
		err := rows.Scan(
			&attendance.RecordID,
			&attendance.UserID,
			&attendance.WorkDate,
			&attendance.CheckInTime,
			&attendance.CheckOutTime,
			&attendance.Status,
			&attendance.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		records = append(records, attendance)
	}

	return records, nil
}

// CreateOrUpdate 출퇴근 기록 생성/업데이트 통합
func (r *AttendanceRepositoryImpl) CreateOrUpdate(userID int, date time.Time, action string, actionTime time.Time, status string) error {
	var query string
	var args []interface{}

	switch action {
	case "checkin":
		query = `INSERT INTO attendance_records (user_id, work_date, check_in_time, created_at) 
				 VALUES (?, ?, ?, NOW()) 
				 ON DUPLICATE KEY UPDATE 
				 check_in_time = VALUES(check_in_time), 
				 created_at = NOW()`
		args = []interface{}{userID, date.Format("2006-01-02"), actionTime}

	case "checkout":
		query = `INSERT INTO attendance_records (user_id, work_date, check_out_time, created_at) 
				 VALUES (?, ?, ?, NOW()) 
				 ON DUPLICATE KEY UPDATE 
				 check_out_time = VALUES(check_out_time), 
				 created_at = NOW()`
		args = []interface{}{userID, date.Format("2006-01-02"), actionTime}

	case "record":
		query = `INSERT INTO attendance_records (user_id, work_date, status, created_at) 
				 VALUES (?, ?, ?, NOW()) 
				 ON DUPLICATE KEY UPDATE 
				 status = VALUES(status), 
				 created_at = NOW()`
		args = []interface{}{userID, date.Format("2006-01-02"), status}
	}

	_, err := r.db.Exec(query, args...)
	return err
}

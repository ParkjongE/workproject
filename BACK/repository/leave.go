package repository

import (
	"BACK/config"
	"time"
)

// LeaveRepository 인터페이스
type LeaveRepository interface {
	FindUserHireDate(userID int) (time.Time, error)
	GetUsedLeaveDays(userID int) (int, error)
	UpdateUserLeaveStatus(userID int, totalEarned, usedDays, remainingDays int) error
	UpdateAllUsersLeaveStatus() (int, error)
}

// LeaveRepositoryImpl 구현체
type LeaveRepositoryImpl struct {
	db config.Database
}

// NewLeaveRepository 생성자
func NewLeaveRepository(db config.Database) LeaveRepository {
	return &LeaveRepositoryImpl{db: db}
}

// FindUserHireDate 사용자 입사일 조회
func (r *LeaveRepositoryImpl) FindUserHireDate(userID int) (time.Time, error) {
	var hireDate time.Time
	query := "SELECT hire_date FROM users WHERE id = ?"
	err := r.db.QueryRow(query, userID).Scan(&hireDate)
	return hireDate, err
}

// GetUsedLeaveDays 사용자의 사용된 연차 일수 조회
func (r *LeaveRepositoryImpl) GetUsedLeaveDays(userID int) (int, error) {
	var usedDays int
	query := `
		SELECT COALESCE(SUM(DATEDIFF(end_date, start_date) + 1), 0) 
		FROM leave_requests 
		WHERE user_id = ? AND leave_type = '연차' AND status = '승인' 
		AND YEAR(start_date) = YEAR(CURDATE())
	`
	err := r.db.QueryRow(query, userID).Scan(&usedDays)
	if err != nil {
		usedDays = 0 // 사용된 연차가 없으면 0
	}
	return usedDays, nil
}

// UpdateUserLeaveStatus 사용자 연차 현황 업데이트/삽입
func (r *LeaveRepositoryImpl) UpdateUserLeaveStatus(userID int, totalEarned, usedDays, remainingDays int) error {
	query := `
		INSERT INTO user_leave_status (user_id, total_earned_days, used_days, remaining_days, updated_at)
		VALUES (?, ?, ?, ?, NOW())
		ON DUPLICATE KEY UPDATE
		total_earned_days = VALUES(total_earned_days),
		used_days = VALUES(used_days),
		remaining_days = VALUES(remaining_days),
		updated_at = NOW()
	`
	_, err := r.db.Exec(query, userID, totalEarned, usedDays, remainingDays)
	return err
}

// UpdateAllUsersLeaveStatus 모든 사용자 연차 현황 업데이트
func (r *LeaveRepositoryImpl) UpdateAllUsersLeaveStatus() (int, error) {
	query := "UPDATE leave_status SET remaining_leave = total_leave - used_leave WHERE year = YEAR(CURDATE())"

	result, err := r.db.Exec(query)
	if err != nil {
		return 0, err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

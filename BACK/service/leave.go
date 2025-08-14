package service

import (
	"BACK/repository"
	"BACK/utils"
	"time"
)

// LeaveServiceImpl 연차 서비스 구현체
type LeaveServiceImpl struct {
	repo repository.LeaveRepository
	lc   *utils.LeaveCalculator
}

// NewLeaveService 연차 서비스 생성자
func NewLeaveService(repo repository.LeaveRepository) *LeaveServiceImpl {
	return &LeaveServiceImpl{
		repo: repo,
		lc:   utils.NewLeaveCalculator(),
	}
}

// GetUserLeaveStatus 사용자 연차 현황 조회 (비즈니스 로직)
func (s *LeaveServiceImpl) GetUserLeaveStatus(userID int) (*utils.LeaveStatus, error) {
	// 비즈니스 검증
	if userID <= 0 {
		return nil, &LeaveError{Message: "유효하지 않은 사용자 ID입니다", Code: "INVALID_USER_ID"}
	}

	// 1. 사용자 정보 가져오기 (hire_date 포함)
	hireDate, err := s.repo.FindUserHireDate(userID)
	if err != nil {
		return nil, &LeaveError{Message: "사용자 정보를 찾을 수 없습니다", Code: "USER_NOT_FOUND"}
	}

	// 2. 사용된 연차 계산 (승인된 연차 신청 기준)
	usedDays, err := s.repo.GetUsedLeaveDays(userID)
	if err != nil {
		return nil, &LeaveError{Message: "연차 사용 정보 조회 실패", Code: "DB_ERROR"}
	}

	// 3. 연차 계산기로 총 연차 계산 (비즈니스 로직)
	leaveStatus := s.lc.CalculateUserLeaveStatus(userID, hireDate, usedDays)

	// 4. DB에 현재 연차 상태 업데이트/삽입
	err = s.repo.UpdateUserLeaveStatus(userID, leaveStatus.TotalEarned, leaveStatus.UsedDays, leaveStatus.RemainingDays)
	if err != nil {
		return nil, &LeaveError{Message: "연차 현황 업데이트 실패", Code: "DB_ERROR"}
	}

	return leaveStatus, nil
}

// GetLeavePolicy 연차 정책 조회 (비즈니스 로직)
func (s *LeaveServiceImpl) GetLeavePolicy() *utils.LeavePolicy {
	return s.lc.GetLeavePolicy()
}

// CalculateLeaveForDate 특정 날짜 연차 계산 (비즈니스 로직)
func (s *LeaveServiceImpl) CalculateLeaveForDate(userID int, targetDate time.Time) (map[string]interface{}, error) {
	// 비즈니스 검증
	if userID <= 0 {
		return nil, &LeaveError{Message: "유효하지 않은 사용자 ID입니다", Code: "INVALID_USER_ID"}
	}

	// 연차 계산 로직 구현 (비즈니스 로직)
	result := map[string]interface{}{
		"user_id":     userID,
		"target_date": targetDate.Format("2006-01-02"),
		"is_holiday":  false,
		"leave_type":  "none",
	}

	// TODO: 실제 연차 계산 로직 추가
	// - 공휴일 확인
	// - 연차 신청 가능 여부 확인
	// - 연차 잔여일 확인

	return result, nil
}

// UpdateAllUsersLeaveStatus 모든 사용자 연차 현황 업데이트 (비즈니스 로직)
func (s *LeaveServiceImpl) UpdateAllUsersLeaveStatus() (int, error) {
	count, err := s.repo.UpdateAllUsersLeaveStatus()
	if err != nil {
		return 0, &LeaveError{Message: "연차 현황 업데이트 실패", Code: "DB_ERROR"}
	}

	return count, nil
}

// LeaveError 연차 관련 에러 구조체
type LeaveError struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

func (e *LeaveError) Error() string {
	return e.Message
}

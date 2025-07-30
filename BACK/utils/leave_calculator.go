package utils

import (
	"BACK/models"
	"time"
)

// LeaveCalculator 연차 계산기
type LeaveCalculator struct{}

// NewLeaveCalculator 새로운 연차 계산기 생성
func NewLeaveCalculator() *LeaveCalculator {
	return &LeaveCalculator{}
}

// CalculateYearsOfService 근무 연수 계산
func (lc *LeaveCalculator) CalculateYearsOfService(hireDate, targetDate time.Time) int {
	years := targetDate.Year() - hireDate.Year()

	// 입사일이 지나지 않았으면 1년 빼기
	if targetDate.YearDay() < hireDate.YearDay() {
		years--
	}

	// 음수인 경우 0으로 설정
	if years < 0 {
		years = 0
	}

	return years
}

// CalculateMonthsOfService 근무 개월 수 계산
func (lc *LeaveCalculator) CalculateMonthsOfService(hireDate, targetDate time.Time) int {
	months := (targetDate.Year()-hireDate.Year())*12 + (int(targetDate.Month()) - int(hireDate.Month()))

	// 일자가 입사일보다 작으면 1개월 빼기
	if targetDate.Day() < hireDate.Day() {
		months--
	}

	// 음수인 경우 0으로 설정
	if months < 0 {
		months = 0
	}

	return months
}

// CalculateLeaveUnderOneYear 1년 미만 연차 계산 (개월 단위)
func (lc *LeaveCalculator) CalculateLeaveUnderOneYear(hireDate, targetDate time.Time) int {
	months := lc.CalculateMonthsOfService(hireDate, targetDate)

	// 1개월당 1일, 최대 11일
	if months > 11 {
		return 11
	}
	return months
}

// CalculateLeaveOverOneYear 1년 이상 연차 계산
func (lc *LeaveCalculator) CalculateLeaveOverOneYear(hireDate, targetDate time.Time) int {
	years := lc.CalculateYearsOfService(hireDate, targetDate)

	// 1년 미만인 경우 0 반환
	if years < 1 {
		return 0
	}

	// 연차 발생 규칙에 따른 일수 계산
	leaveDays := lc.getRegularLeaveDays(years)

	// 최대 25일 제한
	if leaveDays > 25 {
		return 25
	}

	return leaveDays
}

// getRegularLeaveDays 정규 연차 일수 계산
func (lc *LeaveCalculator) getRegularLeaveDays(yearsOfService int) int {
	// 연차 발생 규칙
	switch {
	case yearsOfService <= 2:
		return 15
	case yearsOfService <= 4:
		return 16
	case yearsOfService <= 6:
		return 17
	case yearsOfService <= 8:
		return 18
	case yearsOfService <= 10:
		return 19
	case yearsOfService <= 12:
		return 20
	case yearsOfService <= 14:
		return 21
	case yearsOfService <= 16:
		return 22
	case yearsOfService <= 18:
		return 23
	case yearsOfService <= 20:
		return 24
	default:
		return 25 // 21년차 이상: 최대 25일
	}
}

// CalculateTotalLeave 총 연차 계산
func (lc *LeaveCalculator) CalculateTotalLeave(hireDate, targetDate time.Time) int {
	yearsOfService := lc.CalculateYearsOfService(hireDate, targetDate)

	// 1년 미만인 경우: 개월 단위로 계산 (1개월당 1일, 최대 11일)
	if yearsOfService < 1 {
		return lc.CalculateLeaveUnderOneYear(hireDate, targetDate)
	}

	// 1년 이상인 경우: 정규 연차 정책 적용
	return lc.CalculateLeaveOverOneYear(hireDate, targetDate)
}

// CalculateUserLeaveStatus 사용자별 연차 현황 계산
func (lc *LeaveCalculator) CalculateUserLeaveStatus(userID int, hireDate time.Time, usedDays int) *models.LeaveStatus {
	// 현재 날짜를 기준으로 계산
	targetDate := time.Now()

	// 총 발생 연차 계산
	totalEarned := lc.CalculateTotalLeave(hireDate, targetDate)

	// 잔여 연차 계산
	remainingDays := totalEarned - usedDays
	if remainingDays < 0 {
		remainingDays = 0
	}

	// 근무 연수
	yearsOfService := lc.CalculateYearsOfService(hireDate, targetDate)

	return &models.LeaveStatus{
		UserID:         userID,
		TotalEarned:    totalEarned,
		UsedDays:       usedDays,
		RemainingDays:  remainingDays,
		YearsOfService: yearsOfService,
		HireDate:       hireDate,
	}
}

// GetLeavePolicy 연차 정책 정보 반환
func (lc *LeaveCalculator) GetLeavePolicy() map[string]interface{} {
	return map[string]interface{}{
		"under_one_year": map[string]interface{}{
			"description": "1년 미만 근무자 (개월 단위 계산)",
			"rule":        "1개월당 1일, 최대 11일",
			"max_days":    11,
			"note":        "정규 연차 정책 미적용",
		},
		"over_one_year": map[string]interface{}{
			"description": "1년 이상 근무자 (정규 연차 정책 적용)",
			"rule":        "근무 연수에 따른 정규 연차 발생",
			"max_days":    25,
			"policy": map[string]int{
				"1-2년차":   15,
				"3-4년차":   16,
				"5-6년차":   17,
				"7-8년차":   18,
				"9-10년차":  19,
				"11-12년차": 20,
				"13-14년차": 21,
				"15-16년차": 22,
				"17-18년차": 23,
				"19-20년차": 24,
				"21년차이상":  25,
			},
		},
	}
}

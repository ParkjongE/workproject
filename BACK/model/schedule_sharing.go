package model

import "time"

// ScheduleSharingUser 스케줄 공유에서 사용할 사용자 정보
type ScheduleSharingUser struct {
	ID         int       `json:"id" db:"id"`
	UserID     string    `json:"userid" db:"userid"`
	Name       string    `json:"name" db:"name"`
	Role       string    `json:"role" db:"role"`
	Department string    `json:"department" db:"department"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	HireDate   time.Time `json:"hire_date" db:"hire_date"`
}

// ScheduleSharingUserWithSchedules 일정이 포함된 사용자 정보
type ScheduleSharingUserWithSchedules struct {
	ScheduleSharingUser
	Schedules []Schedule `json:"schedules"`
}

// DepartmentUsers 부서별 사용자 목록
type DepartmentUsers struct {
	Department string                             `json:"department"`
	Users      []ScheduleSharingUserWithSchedules `json:"users"`
}

// ScheduleSharingResponse 스케줄 공유 API 응답
type ScheduleSharingResponse struct {
	Success bool                               `json:"success"`
	Message string                             `json:"message"`
	Data    []ScheduleSharingUserWithSchedules `json:"data,omitempty"`
	Error   string                             `json:"error,omitempty"`
}

// ScheduleResponse 일정 조회 API 응답
type ScheduleResponse struct {
	Success bool       `json:"success"`
	Message string     `json:"message"`
	Data    []Schedule `json:"data,omitempty"`
	Error   string     `json:"error,omitempty"`
}

// ScheduleSharingRequest 스케줄 공유 요청 (필요시)
type ScheduleSharingRequest struct {
	UserID   int    `json:"user_id"`
	Date     string `json:"date"`
	Duration int    `json:"duration"` // 분 단위
}

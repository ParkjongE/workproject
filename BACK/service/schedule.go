package service

import (
	"BACK/model"
	"BACK/repository"
	"fmt"
)

// ScheduleService 스케줄 관련 서비스 인터페이스
type ScheduleService interface {
	CreateSchedule(schedule *model.Schedule) error
	GetSchedulesByUser(userID string) ([]model.Schedule, error)
	UpdateSchedule(scheduleID string, schedule *model.Schedule) error
	DeleteSchedule(scheduleID string) error
}

// ScheduleServiceImpl 스케줄 서비스 구현체
type ScheduleServiceImpl struct {
	repo repository.ScheduleRepository
}

// NewScheduleService 스케줄 서비스 생성자
func NewScheduleService(repo repository.ScheduleRepository) ScheduleService {
	return &ScheduleServiceImpl{repo: repo}
}

// CreateSchedule 스케줄 생성 (비즈니스 로직)
func (s *ScheduleServiceImpl) CreateSchedule(schedule *model.Schedule) error {
	// 비즈니스 로직: 입력값 검증
	if schedule.Title == "" {
		return fmt.Errorf("title is required")
	}
	if schedule.UserID == 0 {
		return fmt.Errorf("user ID is required")
	}

	return s.repo.Create(schedule)
}

// GetSchedulesByUser 사용자별 스케줄 조회
func (s *ScheduleServiceImpl) GetSchedulesByUser(userID string) ([]model.Schedule, error) {
	if userID == "" {
		return nil, fmt.Errorf("user ID is required")
	}

	return s.repo.FindByUserID(userID)
}

// UpdateSchedule 스케줄 수정 (비즈니스 로직)
func (s *ScheduleServiceImpl) UpdateSchedule(scheduleID string, schedule *model.Schedule) error {
	// 비즈니스 로직: 입력값 검증
	if scheduleID == "" {
		return fmt.Errorf("schedule ID is required")
	}
	if schedule.Title == "" {
		return fmt.Errorf("title is required")
	}

	return s.repo.Update(scheduleID, schedule)
}

// DeleteSchedule 스케줄 삭제 (비즈니스 로직)
func (s *ScheduleServiceImpl) DeleteSchedule(scheduleID string) error {
	if scheduleID == "" {
		return fmt.Errorf("schedule ID is required")
	}

	return s.repo.Delete(scheduleID)
}

package service

import (
	"BACK/model"
	"BACK/repository"
	"fmt"
)

// ScheduleSharingService 스케줄 공유 관련 비즈니스 로직
type ScheduleSharingService interface {
	GetUsersByDepartment(department string) ([]model.ScheduleSharingUserWithSchedules, error)
	GetAllDepartmentUsers() ([]model.ScheduleSharingUserWithSchedules, error)
	GetUserSchedulesByID(userID int) ([]model.Schedule, error)
}

// ScheduleSharingServiceImpl ScheduleSharingService 구현체
type ScheduleSharingServiceImpl struct {
	repo repository.ScheduleSharingRepository
}

// NewScheduleSharingService ScheduleSharingService 생성자
func NewScheduleSharingService(repo repository.ScheduleSharingRepository) ScheduleSharingService {
	return &ScheduleSharingServiceImpl{repo: repo}
}

// GetUsersByDepartment 특정 부서의 사용자와 일정 정보 조회
func (s *ScheduleSharingServiceImpl) GetUsersByDepartment(department string) ([]model.ScheduleSharingUserWithSchedules, error) {
	// 1. 부서별 사용자 조회
	users, err := s.repo.GetUsersByDepartment(department)
	if err != nil {
		return nil, fmt.Errorf("부서별 사용자 조회 실패: %v", err)
	}

	// 2. 각 사용자의 일정 정보 조회
	var usersWithSchedules []model.ScheduleSharingUserWithSchedules
	for _, user := range users {
		schedules, err := s.repo.GetUserSchedules(user.ID)
		if err != nil {
			// 일정 조회 실패해도 사용자 정보는 반환 (빈 일정 배열)
			schedules = []model.Schedule{}
		}

		userWithSchedules := model.ScheduleSharingUserWithSchedules{
			ScheduleSharingUser: user,
			Schedules:           schedules, // 빈 배열이라도 설정
		}

		usersWithSchedules = append(usersWithSchedules, userWithSchedules)
	}

	return usersWithSchedules, nil
}

// GetAllDepartmentUsers 모든 부서의 사용자와 일정 정보 조회
func (s *ScheduleSharingServiceImpl) GetAllDepartmentUsers() ([]model.ScheduleSharingUserWithSchedules, error) {
	// 1. 모든 부서의 사용자 조회
	departmentUsers, err := s.repo.GetAllDepartmentUsers()
	if err != nil {
		return nil, fmt.Errorf("전체 부서 사용자 조회 실패: %v", err)
	}

	// 2. 모든 사용자를 하나의 배열로 합치기
	var allUsers []model.ScheduleSharingUserWithSchedules
	for _, dept := range departmentUsers {
		allUsers = append(allUsers, dept.Users...)
	}

	return allUsers, nil
}

// GetUserSchedulesByID 특정 사용자의 일정만 조회
func (s *ScheduleSharingServiceImpl) GetUserSchedulesByID(userID int) ([]model.Schedule, error) {
	schedules, err := s.repo.GetUserSchedules(userID)
	if err != nil {
		return nil, err
	}

	// 일정이 없는 경우 빈 배열 반환
	if schedules == nil {
		return []model.Schedule{}, nil
	}

	return schedules, nil
}

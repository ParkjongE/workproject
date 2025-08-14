package repository

import (
	"BACK/config"
	"BACK/model"
)

// ScheduleRepository 인터페이스
type ScheduleRepository interface {
	Create(schedule *model.Schedule) error
	FindByUserID(userID string) ([]model.Schedule, error)
	Update(scheduleID string, schedule *model.Schedule) error
	Delete(scheduleID string) error
}

// ScheduleRepositoryImpl 구현체
type ScheduleRepositoryImpl struct {
	db config.Database
}

// NewScheduleRepository 생성자
func NewScheduleRepository(db config.Database) ScheduleRepository {
	return &ScheduleRepositoryImpl{db: db}
}

// Create 스케줄 생성
func (r *ScheduleRepositoryImpl) Create(schedule *model.Schedule) error {
	query := "INSERT INTO schedules (user_id, title, content, start_time, end_time, color, created_at) VALUES (?, ?, ?, ?, ?, ?, NOW())"
	_, err := r.db.Exec(query, schedule.UserID, schedule.Title, schedule.Content, schedule.StartTime, schedule.EndTime, schedule.Color)
	return err
}

// FindByUserID 사용자별 스케줄 조회
func (r *ScheduleRepositoryImpl) FindByUserID(userID string) ([]model.Schedule, error) {
	query := "SELECT schedule_id, user_id, title, content, start_time, end_time, created_at, color FROM schedules WHERE user_id = ? ORDER BY start_time"

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []model.Schedule
	for rows.Next() {
		var schedule model.Schedule
		err := rows.Scan(
			&schedule.ScheduleID, &schedule.UserID, &schedule.Title, &schedule.Content,
			&schedule.StartTime, &schedule.EndTime, &schedule.CreatedAt, &schedule.Color,
		)
		if err != nil {
			return nil, err
		}
		schedules = append(schedules, schedule)
	}

	return schedules, nil
}

// Update 스케줄 수정
func (r *ScheduleRepositoryImpl) Update(scheduleID string, schedule *model.Schedule) error {
	query := "UPDATE schedules SET title = ?, content = ?, start_time = ?, end_time = ?, color = ? WHERE schedule_id = ?"
	_, err := r.db.Exec(query, schedule.Title, schedule.Content, schedule.StartTime, schedule.EndTime, schedule.Color, scheduleID)
	return err
}

// Delete 스케줄 삭제
func (r *ScheduleRepositoryImpl) Delete(scheduleID string) error {
	query := "DELETE FROM schedules WHERE schedule_id = ?"
	_, err := r.db.Exec(query, scheduleID)
	return err
}

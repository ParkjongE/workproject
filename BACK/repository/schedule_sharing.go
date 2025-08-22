package repository

import (
	"BACK/model"
	"database/sql"
	"fmt"
	"log"
)

// ScheduleSharingRepository 스케줄 공유 관련 데이터베이스 작업
type ScheduleSharingRepository interface {
	GetUsersByDepartment(department string) ([]model.ScheduleSharingUser, error)
	GetUserSchedules(userID int) ([]model.Schedule, error)
	GetUserSchedulesByID(userID int) ([]model.Schedule, error) // 추가
	GetAllDepartmentUsers() ([]model.DepartmentUsers, error)
	GetUserByID(userID int) (*model.ScheduleSharingUser, error)
}

// ScheduleSharingRepositoryImpl 스케줄 공유 Repository 구현체
type ScheduleSharingRepositoryImpl struct {
	db *sql.DB
}

// NewScheduleSharingRepository ScheduleSharingRepository 생성자
func NewScheduleSharingRepository(db *sql.DB) ScheduleSharingRepository {
	return &ScheduleSharingRepositoryImpl{db: db}
}

// GetUsersByDepartment 특정 부서의 사용자 목록 조회
func (r *ScheduleSharingRepositoryImpl) GetUsersByDepartment(department string) ([]model.ScheduleSharingUser, error) {
	query := `
		SELECT id, userid, name, role, department, created_at, hire_date
		FROM users 
		WHERE department = ? 
		ORDER BY 
			CASE role
				WHEN '부장' THEN 1
				WHEN '차장' THEN 2
				WHEN '과장' THEN 3
				WHEN '대리' THEN 4
				WHEN '사원' THEN 5
				ELSE 6
			END,
			name
	`

	rows, err := r.db.Query(query, department)
	if err != nil {
		return nil, fmt.Errorf("부서별 사용자 조회 실패: %v", err)
	}
	defer rows.Close()

	var users []model.ScheduleSharingUser
	for rows.Next() {
		var user model.ScheduleSharingUser
		err := rows.Scan(
			&user.ID,
			&user.UserID,
			&user.Name,
			&user.Role,
			&user.Department,
			&user.CreatedAt,
			&user.HireDate,
		)
		if err != nil {
			log.Printf("사용자 데이터 스캔 실패: %v", err)
			continue
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("행 처리 중 오류: %v", err)
	}

	log.Printf("부서 %s에서 %d명의 사용자를 조회했습니다", department, len(users))
	return users, nil
}

// GetUserSchedules 특정 사용자의 일정 목록 조회
func (r *ScheduleSharingRepositoryImpl) GetUserSchedules(userID int) ([]model.Schedule, error) {
	query := `
		SELECT schedule_id, user_id, title, start_time, end_time, color, created_at
		FROM schedules 
		WHERE user_id = ? 
		ORDER BY start_time
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("사용자 %d 일정 조회 실패: %v", userID, err)
	}
	defer rows.Close()

	var schedules []model.Schedule
	for rows.Next() {
		var schedule model.Schedule
		err := rows.Scan(
			&schedule.ScheduleID, // id 대신 ScheduleID 사용
			&schedule.UserID,
			&schedule.Title,
			&schedule.StartTime,
			&schedule.EndTime,
			&schedule.Color,
			&schedule.CreatedAt,
		)
		if err != nil {
			log.Printf("일정 데이터 스캔 실패: %v", err)
			continue
		}
		schedules = append(schedules, schedule)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("행 처리 중 오류: %v", err)
	}

	log.Printf("사용자 %d의 %d개 일정을 조회했습니다", userID, len(schedules))
	return schedules, nil
}

// GetAllDepartmentUsers 모든 부서의 사용자와 일정 정보 조회
func (r *ScheduleSharingRepositoryImpl) GetAllDepartmentUsers() ([]model.DepartmentUsers, error) {
	departments := []string{"개발부", "디자인부", "마케팅부"}
	var result []model.DepartmentUsers

	for _, dept := range departments {
		// 1. 부서별 사용자 조회
		users, err := r.GetUsersByDepartment(dept)
		if err != nil {
			log.Printf("부서 %s 사용자 조회 실패: %v", dept, err)
			continue
		}

		// 2. 각 사용자의 일정 정보 조회
		var usersWithSchedules []model.ScheduleSharingUserWithSchedules
		for _, user := range users {
			schedules, err := r.GetUserSchedules(user.ID)
			if err != nil {
				log.Printf("사용자 %s(%d) 일정 조회 실패: %v", user.Name, user.ID, err)
				schedules = []model.Schedule{} // 빈 일정 배열로 설정
			}

			userWithSchedules := model.ScheduleSharingUserWithSchedules{
				ScheduleSharingUser: user,
				Schedules:           schedules,
			}

			usersWithSchedules = append(usersWithSchedules, userWithSchedules)
		}

		// 3. 부서별 결과 추가
		result = append(result, model.DepartmentUsers{
			Department: dept,
			Users:      usersWithSchedules,
		})

		log.Printf("부서 %s: %d명의 사용자 정보를 로드했습니다", dept, len(usersWithSchedules))
	}

	return result, nil
}

// GetUserByID 특정 ID의 사용자 정보 조회
func (r *ScheduleSharingRepositoryImpl) GetUserByID(userID int) (*model.ScheduleSharingUser, error) {
	query := `
		SELECT id, userid, name, role, department, created_at, hire_date
		FROM users 
		WHERE id = ?
	`

	var user model.ScheduleSharingUser
	err := r.db.QueryRow(query, userID).Scan(
		&user.ID,
		&user.UserID,
		&user.Name,
		&user.Role,
		&user.Department,
		&user.CreatedAt,
		&user.HireDate,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("사용자 ID %d를 찾을 수 없습니다", userID)
		}
		return nil, fmt.Errorf("사용자 조회 실패: %v", err)
	}

	return &user, nil
}

// GetUsersByRole 특정 직급의 사용자 목록 조회 (추가 기능)
func (r *ScheduleSharingRepositoryImpl) GetUsersByRole(role string) ([]model.ScheduleSharingUser, error) {
	query := `
		SELECT id, userid, name, role, department, created_at, hire_date
		FROM users 
		WHERE role = ? 
		ORDER BY department, name
	`

	rows, err := r.db.Query(query, role)
	if err != nil {
		return nil, fmt.Errorf("직급별 사용자 조회 실패: %v", err)
	}
	defer rows.Close()

	var users []model.ScheduleSharingUser
	for rows.Next() {
		var user model.ScheduleSharingUser
		err := rows.Scan(
			&user.ID,
			&user.UserID,
			&user.Name,
			&user.Role,
			&user.Department,
			&user.CreatedAt,
			&user.HireDate,
		)
		if err != nil {
			log.Printf("사용자 데이터 스캔 실패: %v", err)
			continue
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("행 처리 중 오류: %v", err)
	}

	return users, nil
}

// GetUserSchedulesByID 특정 사용자의 일정 목록 조회 (GetUserSchedules와 동일한 기능)
func (r *ScheduleSharingRepositoryImpl) GetUserSchedulesByID(userID int) ([]model.Schedule, error) {
	return r.GetUserSchedules(userID)
}

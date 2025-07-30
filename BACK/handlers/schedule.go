package handlers

import (
	"BACK/config"
	"BACK/models"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// AddSchedule 함수는 일정(스케줄) 추가 요청을 처리합니다.
// POST /api/schedule
// 요청 본문(JSON): user_id, title, content, project_id, start_time, end_time
func AddSchedule(w http.ResponseWriter, r *http.Request) {
	// POST 메서드만 허용
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 요청 본문에서 일정 정보 파싱
	var sched models.Schedule
	err := json.NewDecoder(r.Body).Decode(&sched)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	fmt.Println("[AddSchedule] 받은 데이터:", sched)
	fmt.Println("[AddSchedule] 받은 color 값:", sched.Color)
	fmt.Println("[AddSchedule] 받은 content 값:", sched.Content)

	// DB에 일정 정보 저장
	_, err = config.DB.Exec(
		"INSERT INTO schedules (user_id, title, content, project_id, start_time, end_time, color) VALUES (?, ?, ?, ?, ?, ?, ?)",
		sched.UserID, sched.Title, sched.Content, sched.ProjectID, sched.StartTime, sched.EndTime, sched.Color,
	)
	if err != nil {
		fmt.Println("[AddSchedule] DB 오류:", err)
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}

	// 성공 응답 반환
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
}

// GetSchedulesByUser 함수는 특정 사원의 모든 일정을 조회합니다.
// GET /api/schedules?user_id=...
// user_id: 조회할 사원 ID (쿼리 파라미터)
func GetSchedulesByUser(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "user_id required", http.StatusBadRequest)
		return
	}

	// DB에서 해당 사원의 일정 목록 조회
	rows, err := config.DB.Query("SELECT schedule_id, user_id, title, content, project_id, start_time, end_time, created_at, color FROM schedules WHERE user_id = ?", userID)
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	schedules := []models.Schedule{}
	for rows.Next() {
		var s models.Schedule
		err := rows.Scan(&s.ScheduleID, &s.UserID, &s.Title, &s.Content, &s.ProjectID, &s.StartTime, &s.EndTime, &s.CreatedAt, &s.Color)
		if err != nil {
			continue // 에러 발생 시 해당 row는 건너뜀
		}
		schedules = append(schedules, s)
	}

	// 일정 목록을 JSON으로 반환
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(schedules)
}

// DeleteSchedule 함수는 특정 일정을 삭제합니다.
// DELETE /api/schedule/{schedule_id}
// schedule_id: 삭제할 일정 ID (URL 파라미터)
func DeleteSchedule(w http.ResponseWriter, r *http.Request) {
	// DELETE 메서드만 허용
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// URL에서 schedule_id 추출
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "Invalid schedule ID", http.StatusBadRequest)
		return
	}
	scheduleID := pathParts[3]

	// DB에서 해당 일정 삭제
	result, err := config.DB.Exec("DELETE FROM schedules WHERE schedule_id = ?", scheduleID)
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}

	// 삭제된 행 수 확인
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Schedule not found", http.StatusNotFound)
		return
	}

	// 성공 응답 반환
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "message": "일정이 삭제되었습니다."})
}

// UpdateSchedule 함수는 특정 일정을 수정합니다.
// PUT /api/schedule/{schedule_id}
// schedule_id: 수정할 일정 ID (URL 파라미터)
// 요청 본문(JSON): title, content, project_id, start_time, end_time, color
func UpdateSchedule(w http.ResponseWriter, r *http.Request) {
	// PUT 메서드만 허용
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// URL에서 schedule_id 추출
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "Invalid schedule ID", http.StatusBadRequest)
		return
	}
	scheduleID := pathParts[3]

	// 요청 본문에서 일정 정보 파싱
	var sched models.Schedule
	err := json.NewDecoder(r.Body).Decode(&sched)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// DB에서 해당 일정 수정
	result, err := config.DB.Exec(
		"UPDATE schedules SET title = ?, content = ?, project_id = ?, start_time = ?, end_time = ?, color = ? WHERE schedule_id = ?",
		sched.Title, sched.Content, sched.ProjectID, sched.StartTime, sched.EndTime, sched.Color, scheduleID,
	)
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}

	// 수정된 행 수 확인
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Schedule not found", http.StatusNotFound)
		return
	}

	// 성공 응답 반환
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "message": "일정이 수정되었습니다."})
}

// HandleScheduleOperations 함수는 HTTP 메서드에 따라 삭제 또는 수정을 처리합니다.
// DELETE /api/schedule/{schedule_id} - 일정 삭제
// PUT /api/schedule/{schedule_id} - 일정 수정
func HandleScheduleOperations(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodDelete:
		DeleteSchedule(w, r)
	case http.MethodPut:
		UpdateSchedule(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

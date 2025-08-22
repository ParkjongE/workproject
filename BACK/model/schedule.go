package model

type Schedule struct {
	ID          int    `json:"id"`
	ScheduleID  int    `json:"schedule_id"`
	UserID      int    `json:"user_id"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	Description string `json:"description"`
	ProjectID   int    `json:"project_id"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	Type        string `json:"type"` // 추가
	Color       string `json:"color"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// ==================== Response 구조체들 ====================

// ScheduleAPIResponse 스케줄 API 응답 구조체
type ScheduleAPIResponse struct {
	Success bool      `json:"success"`
	Message string    `json:"message"`
	Data    *Schedule `json:"data,omitempty"`
}

// ScheduleSuccessResponse 스케줄 성공 응답 구조체
type ScheduleSuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

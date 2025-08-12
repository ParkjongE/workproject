package models

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
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	Color       string `json:"color"`
}

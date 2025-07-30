package models

type Schedule struct {
	ScheduleID int    `json:"schedule_id"`
	UserID     int    `json:"user_id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	ProjectID  int    `json:"project_id"`
	StartTime  string `json:"start_time"`
	EndTime    string `json:"end_time"`
	CreatedAt  string `json:"created_at"`
	Color      string `json:"color"`
}

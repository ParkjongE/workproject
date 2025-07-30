package models

// DepartmentStatus는 부서별 업무현황을 나타내는 구조체
type DepartmentStatus struct {
	DepartmentName string `json:"department_name"` // 부서 이름
	TeamName       string `json:"team_name"`       // 팀 이름
	TeamLead       string `json:"team_lead"`       // 팀장 이름
	MemberCount    int    `json:"member_count"`    // 팀원 수
	MainTasks      string `json:"main_tasks"`      // 주요 담당 업무
	TaskName       string `json:"task_name"`       // 업무명
	TaskAssignee   string `json:"task_assignee"`   // 업무 담당자
	Status         string `json:"status"`          // 진행 상태
	Priority       string `json:"priority"`        // 우선순위
	StartDate      string `json:"start_date"`      // 시작일
	Deadline       string `json:"deadline"`        // 마감일
}

// 샘플 데이터 정의
var DepartmentStatuses = []DepartmentStatus{
	{
		DepartmentName: "개발팀",
		TeamName:       "웹팀",
		TeamLead:       "김민수",
		MemberCount:    5,
		MainTasks:      "웹 서비스 개발",
		TaskName:       "메인 페이지 리뉴얼",
		TaskAssignee:   "박지훈",
		Status:         "진행중",
		Priority:       "🔴 긴급",
		StartDate:      "2025-07-01",
		Deadline:       "2025-07-20",
	},
	{
		DepartmentName: "디자인팀",
		TeamName:       "UI팀",
		TeamLead:       "강서연",
		MemberCount:    3,
		MainTasks:      "UI/UX 개선",
		TaskName:       "모바일 UI 수정",
		TaskAssignee:   "한지수",
		Status:         "지연",
		Priority:       "🟠 보통",
		StartDate:      "2025-07-05",
		Deadline:       "2025-07-15",
	},
	{
		DepartmentName: "마케팅팀",
		TeamName:       "콘텐츠팀",
		TeamLead:       "이현우",
		MemberCount:    4,
		MainTasks:      "SNS 운영",
		TaskName:       "캠페인 포스터 검토",
		TaskAssignee:   "최현우",
		Status:         "완료",
		Priority:       "🟢 여유",
		StartDate:      "2025-06-20",
		Deadline:       "2025-07-10",
	},
	{
		DepartmentName: "개발팀",
		TeamName:       "백엔드팀",
		TeamLead:       "정민호",
		MemberCount:    6,
		MainTasks:      "서버 개발 및 유지보수",
		TaskName:       "API 서버 성능 최적화",
		TaskAssignee:   "이준영",
		Status:         "진행중",
		Priority:       "🟠 보통",
		StartDate:      "2025-07-03",
		Deadline:       "2025-07-25",
	},
	{
		DepartmentName: "마케팅팀",
		TeamName:       "광고팀",
		TeamLead:       "김소연",
		MemberCount:    3,
		MainTasks:      "광고 기획 및 집행",
		TaskName:       "여름 시즌 프로모션",
		TaskAssignee:   "박민지",
		Status:         "진행중",
		Priority:       "🔴 긴급",
		StartDate:      "2025-07-10",
		Deadline:       "2025-07-30",
	},
	{
		DepartmentName: "디자인팀",
		TeamName:       "그래픽팀",
		TeamLead:       "황지영",
		MemberCount:    4,
		MainTasks:      "그래픽 디자인",
		TaskName:       "브랜드 로고 리디자인",
		TaskAssignee:   "김태호",
		Status:         "완료",
		Priority:       "🟢 여유",
		StartDate:      "2025-06-15",
		Deadline:       "2025-07-05",
	},
}

// 부서 이름 목록
var Departments = []string{"개발팀", "디자인팀", "마케팅팀"}

//models : 데이터의 구조(모양)를 정의하는 곳

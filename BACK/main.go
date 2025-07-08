package main

import (
    "github.com/gin-gonic/gin" // Gin 웹 프레임워크 임포트
    "net/http"                 // HTTP 상태 코드 정의 사용
)

// DepartmentStatus는 부서별 업무현황을 나타내는 구조체
type DepartmentStatus struct {
    DepartmentName string `json:"department_name"` // 부서 이름
    TeamName       string `json:"team_name"`       // 팀 이름
    TeamLead       string `json:"team_lead"`       // 팀장 이름
    MemberCount    int    `json:"member_count"`    // 팀원 수
    Tasks          string `json:"tasks"`           // 담당 업무
    OngoingProject string `json:"ongoing_project"` // 진행 중인 프로젝트
    Deadline       string `json:"deadline"`        // 마감 기한
}

// departmentStatuses는 부서별 업무현황에 대한 샘플 데이터
var departmentStatuses = []DepartmentStatus{
    {
        DepartmentName: "영업부",
        TeamName:       "영업1팀",
        TeamLead:       "김팀장",
        MemberCount:    5,
        Tasks:          "고객 관리, 매출 보고",
        OngoingProject: "신규 시장 진출",
        Deadline:       "2025-08-15",
    },
    {
        DepartmentName: "개발부",
        TeamName:       "백엔드팀",
        TeamLead:       "이팀장",
        MemberCount:    7,
        Tasks:          "API 개발, 서버 관리",
        OngoingProject: "업무관리 시스템",
        Deadline:       "2025-09-01",
    },
    {
        DepartmentName: "개발부",
        TeamName:       "프론트엔드팀",
        TeamLead:       "박팀장",
        MemberCount:    4,
        Tasks:          "UI/UX 개발, 웹 개발",
        OngoingProject: "사용자 인터페이스 개선",
        Deadline:       "2025-08-30",
    },
}

// departments는 부서 이름 목록에 대한 샘플 데이터
var departments = []string{"영업부", "개발부", "디자인부"}

func main() {
    // Gin 엔진 생성 (기본 미들웨어 포함: 로깅, 복구)
    r := gin.Default()

    // CORS 미들웨어 추가
    r.Use(func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
        
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        
        c.Next()
    })

    // 정적 파일 서빙 설정 - 상위 디렉토리의 FRONT 폴더를 서빙
    r.Static("/static", "../FRONT")
    r.LoadHTMLGlob("../FRONT/*.html")

    // 루트 경로에서 login.html 서빙
    r.GET("/", func(c *gin.Context) {
        c.HTML(http.StatusOK, "login.html", nil)
    })

    // /orgchart 경로에서 orgchart.html 서빙
    r.GET("/orgchart", func(c *gin.Context) {
        c.HTML(http.StatusOK, "orgchart.html", nil)
    })

    // /orgmember 경로에서 orgmember.html 서빙
    r.GET("/orgmember", func(c *gin.Context) {
        c.HTML(http.StatusOK, "orgmember.html", nil)
    })

    // /dep_work 경로에서 dep_work.html 서빙
    r.GET("/dep_work", func(c *gin.Context) {
        c.HTML(http.StatusOK, "dep_work.html", nil)
    })

    // API 1: 모든 부서 이름 목록 조회
    r.GET("/api/departments", func(c *gin.Context) {
        // 부서 이름 목록을 JSON으로 반환
        c.JSON(http.StatusOK, gin.H{
            "departments": departments,
        })
    })

    // API 2: 모든 부서의 업무현황 조회
    r.GET("/api/departments/status", func(c *gin.Context) {
        // 부서별 업무현황 리스트를 JSON으로 반환
        c.JSON(http.StatusOK, departmentStatuses)
    })

    // API 3: 특정 부서의 업무현황 조회
    r.GET("/api/departments/status/:name", func(c *gin.Context) {
        name := c.Param("name")
        for _, dept := range departmentStatuses {
            if dept.DepartmentName == name {
                // 배열 형태로 반환 (프론트가 반복문으로 처리 가능)
                c.JSON(http.StatusOK, []DepartmentStatus{dept})
                return
            }
        }
        c.JSON(http.StatusNotFound, gin.H{"error": "해당 부서를 찾을 수 없습니다"})
    })
    // API 4: 구성원 목록 조회
    r.GET("/api/members", func(c *gin.Context) {
        // 구성원 이름 목록을 JSON으로 반환
        c.JSON(http.StatusOK, gin.H{
            "members": []string{"홍길동", "김영희"},
        })
    })

    // HTTP 서버 실행 (localhost:8080)
    r.Run(":8080")
}

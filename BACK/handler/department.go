package handler

import (
	"BACK/model"
	"encoding/json"
	"net/http"
	"strings"
)

// GetDepartments는 모든 부서 이름 목록을 반환하는 핸들러 함수입니다
func GetDepartments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string][]string{"departments": model.Departments}) // map[키타입]값타입
	// 키타입: string, 값타입: []string
	// 키타입: departments, 값타입: models.Departments
	// 키타입: departments, 값타입: models.Departments
}

// GetAllDepartmentStatuses는 모든 부서의 업무현황을 반환하는 핸들러 함수입니다
func GetAllDepartmentStatuses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(model.DepartmentStatuses)
}

// GetDepartmentStatus는 특정 부서의 업무현황을 반환하는 핸들러 함수입니다
func GetDepartmentStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// URL에서 부서 이름 추출
	path := strings.Split(r.URL.Path, "/")
	if len(path) < 5 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	name := path[4] // /api/departments/status/{name}

	var result []model.DepartmentStatus
	for _, dept := range model.DepartmentStatuses {
		if dept.DepartmentName == name {
			result = append(result, dept)
		}
	}

	if len(result) > 0 {
		json.NewEncoder(w).Encode(result)
	} else {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "해당 부서를 찾을 수 없습니다"})
	}
}

// GetMembers는 구성원 목록을 반환하는 핸들러 함수입니다
func GetMembers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string][]string{"members": {"홍길동", "김영희"}})
}

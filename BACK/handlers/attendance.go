package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"BACK/config"
	"BACK/models"
)

// CheckIn 출근 처리
func CheckIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.AttendanceRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	fmt.Println("[CheckIn] 받은 데이터:", req)

	// 클라이언트가 보낸 시간 사용, 없으면 현재 시간 사용
	var checkInTime time.Time
	if req.CheckInTime != "" {
		parsedTime, err := time.Parse(time.RFC3339, req.CheckInTime)
		if err == nil {
			checkInTime = parsedTime
		} else {
			checkInTime = time.Now()
		}
	} else {
		checkInTime = time.Now()
	}
	today := checkInTime.Format("2006-01-02")

	// 기존 기록 확인
	var recordID int
	var userID int
	var workDate string
	var checkInTimeStr sql.NullString
	var checkOutTimeStr sql.NullString

	err = config.DB.QueryRow(
		"SELECT record_id, user_id, work_date, check_in_time, check_out_time FROM attendance_records WHERE user_id = ? AND work_date = ?",
		req.UserID, today,
	).Scan(&recordID, &userID, &workDate, &checkInTimeStr, &checkOutTimeStr)

	if err == sql.ErrNoRows {
		// 새 기록 생성
		_, err = config.DB.Exec(
			"INSERT INTO attendance_records (user_id, work_date, check_in_time) VALUES (?, ?, ?)",
			req.UserID, today, checkInTime,
		)
		if err != nil {
			fmt.Println("[CheckIn] DB 오류:", err)
			http.Error(w, "DB error", http.StatusInternalServerError)
			return
		}
	} else if err != nil {
		fmt.Println("[CheckIn] DB 조회 오류:", err)
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	} else {
		// 기존 기록이 있고 출근 시간이 없으면 업데이트
		if !checkInTimeStr.Valid {
			_, err = config.DB.Exec(
				"UPDATE attendance_records SET check_in_time = ? WHERE record_id = ?",
				checkInTime, recordID,
			)
			if err != nil {
				fmt.Println("[CheckIn] DB 업데이트 오류:", err)
				http.Error(w, "DB error", http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, "Already checked in today", http.StatusConflict)
			return
		}
	}

	// 성공 응답
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":       true,
		"message":       "출근이 기록되었습니다",
		"check_in_time": checkInTime.Format("15:04:05"),
	})
}

// CheckOut 퇴근 처리
func CheckOut(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.AttendanceRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	fmt.Println("[CheckOut] 받은 데이터:", req)

	// 클라이언트가 보낸 시간 사용, 없으면 현재 시간 사용
	var checkOutTime time.Time
	if req.CheckOutTime != "" {
		parsedTime, err := time.Parse(time.RFC3339, req.CheckOutTime)
		if err == nil {
			checkOutTime = parsedTime
		} else {
			checkOutTime = time.Now()
		}
	} else {
		checkOutTime = time.Now()
	}
	today := checkOutTime.Format("2006-01-02")

	// 기존 기록 확인
	var recordID int
	var userID int
	var workDate string
	var checkInTimeStr sql.NullString
	var checkOutTimeStr sql.NullString

	err = config.DB.QueryRow(
		"SELECT record_id, user_id, work_date, check_in_time, check_out_time FROM attendance_records WHERE user_id = ? AND work_date = ?",
		req.UserID, today,
	).Scan(&recordID, &userID, &workDate, &checkInTimeStr, &checkOutTimeStr)

	if err == sql.ErrNoRows {
		http.Error(w, "No check-in record found for today", http.StatusNotFound)
		return
	} else if err != nil {
		fmt.Println("[CheckOut] DB 조회 오류:", err)
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}

	// 출근 기록이 없으면 오류
	if !checkInTimeStr.Valid {
		http.Error(w, "No check-in record found for today", http.StatusNotFound)
		return
	}

	// 이미 퇴근했으면 오류
	if checkOutTimeStr.Valid {
		http.Error(w, "Already checked out today", http.StatusConflict)
		return
	}

	// 출근 시간 파싱
	checkInTime, err := time.Parse("2006-01-02 15:04:05", checkInTimeStr.String)
	if err != nil {
		fmt.Println("[CheckOut] 출근 시간 파싱 오류:", err)
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}

	// 근무 시간 계산
	workHours := checkOutTime.Sub(checkInTime).Hours()

	// 퇴근 시간 업데이트
	_, err = config.DB.Exec(
		"UPDATE attendance_records SET check_out_time = ? WHERE record_id = ?",
		checkOutTime, recordID,
	)
	if err != nil {
		fmt.Println("[CheckOut] DB 업데이트 오류:", err)
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}

	// 성공 응답
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":        true,
		"message":        "퇴근이 기록되었습니다",
		"check_out_time": checkOutTime.Format("15:04:05"),
		"work_hours":     workHours,
	})
}

// GetAttendanceRecords 사용자의 출퇴근 기록 조회
func GetAttendanceRecords(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// URL에서 user_id 추출
	pathParts := strings.Split(r.URL.Path, "/")
	fmt.Println("[GetAttendanceRecords] URL 경로:", r.URL.Path)
	fmt.Println("[GetAttendanceRecords] 경로 파트:", pathParts)

	if len(pathParts) < 5 {
		fmt.Println("[GetAttendanceRecords] 경로 파트 부족:", len(pathParts))
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	userID, err := strconv.Atoi(pathParts[4])
	if err != nil {
		fmt.Println("[GetAttendanceRecords] user_id 파싱 오류:", err, "값:", pathParts[4])
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	fmt.Println("[GetAttendanceRecords] 추출된 user_id:", userID)

	// 쿼리 파라미터에서 날짜 범위 확인
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	var query string
	var args []interface{}

	if startDate != "" && endDate != "" {
		query = "SELECT record_id, user_id, work_date, check_in_time, check_out_time, status, created_at FROM attendance_records WHERE user_id = ? AND work_date BETWEEN ? AND ? ORDER BY work_date DESC"
		args = []interface{}{userID, startDate, endDate}
	} else {
		query = "SELECT record_id, user_id, work_date, check_in_time, check_out_time, status, created_at FROM attendance_records WHERE user_id = ? ORDER BY work_date DESC LIMIT 30"
		args = []interface{}{userID}
	}

	rows, err := config.DB.Query(query, args...)
	if err != nil {
		fmt.Println("[GetAttendanceRecords] DB 오류:", err)
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var records []models.AttendanceRecord
	for rows.Next() {
		var recordID int
		var userID int
		var workDate string
		var checkInTimeStr sql.NullString
		var checkOutTimeStr sql.NullString
		var status sql.NullString
		var createdAtStr sql.NullString

		err := rows.Scan(
			&recordID,
			&userID,
			&workDate,
			&checkInTimeStr,
			&checkOutTimeStr,
			&status,
			&createdAtStr,
		)
		if err != nil {
			fmt.Println("[GetAttendanceRecords] 스캔 오류:", err)
			continue
		}

		// sql.NullString을 *time.Time으로 변환
		var checkInTime *time.Time
		var checkOutTime *time.Time

		if checkInTimeStr.Valid {
			parsedTime, err := time.Parse("2006-01-02 15:04:05", checkInTimeStr.String)
			if err == nil {
				checkInTime = &parsedTime
			}
		}

		if checkOutTimeStr.Valid {
			parsedTime, err := time.Parse("2006-01-02 15:04:05", checkOutTimeStr.String)
			if err == nil {
				checkOutTime = &parsedTime
			}
		}

		var statusPtr *string
		if status.Valid {
			statusPtr = &status.String
		}

		// createdAt 파싱
		var createdAt time.Time
		if createdAtStr.Valid {
			parsedTime, err := time.Parse("2006-01-02 15:04:05", createdAtStr.String)
			if err == nil {
				createdAt = parsedTime
			} else {
				createdAt = time.Now() // 기본값
			}
		} else {
			createdAt = time.Now() // 기본값
		}

		record := models.AttendanceRecord{
			RecordID:     recordID,
			UserID:       userID,
			WorkDate:     workDate,
			CheckInTime:  checkInTime,
			CheckOutTime: checkOutTime,
			Status:       statusPtr,
			CreatedAt:    createdAt,
		}

		records = append(records, record)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(records)
}

// GetAttendanceByDate 특정 날짜의 출퇴근 기록 조회
func GetAttendanceByDate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// URL에서 user_id와 date 추출
	pathParts := strings.Split(r.URL.Path, "/")
	fmt.Println("[GetAttendanceByDate] URL 경로:", r.URL.Path)
	fmt.Println("[GetAttendanceByDate] 경로 파트:", pathParts)

	if len(pathParts) < 6 {
		fmt.Println("[GetAttendanceByDate] 경로 파트 부족:", len(pathParts))
		http.Error(w, "Invalid user ID or date", http.StatusBadRequest)
		return
	}
	userID, err := strconv.Atoi(pathParts[4])
	if err != nil {
		fmt.Println("[GetAttendanceByDate] user_id 파싱 오류:", err, "값:", pathParts[4])
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	targetDate := pathParts[5]
	fmt.Println("[GetAttendanceByDate] 추출된 user_id:", userID, "date:", targetDate)

	var recordID int
	var dbUserID int
	var workDate string
	var checkInTimeStr sql.NullString
	var checkOutTimeStr sql.NullString
	var status sql.NullString
	var createdAtStr sql.NullString

	err = config.DB.QueryRow(
		"SELECT record_id, user_id, work_date, check_in_time, check_out_time, status, created_at FROM attendance_records WHERE user_id = ? AND work_date = ?",
		userID, targetDate,
	).Scan(
		&recordID,
		&dbUserID,
		&workDate,
		&checkInTimeStr,
		&checkOutTimeStr,
		&status,
		&createdAtStr,
	)

	if err == sql.ErrNoRows {
		// 해당 날짜 기록이 없으면 빈 응답
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"user_id":        userID,
			"work_date":      targetDate,
			"check_in_time":  nil,
			"check_out_time": nil,
		})
		return
	} else if err != nil {
		fmt.Println("[GetAttendanceByDate] DB 오류:", err)
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}

	// sql.NullString을 *time.Time으로 변환
	var checkInTime *time.Time
	var checkOutTime *time.Time

	if checkInTimeStr.Valid {
		parsedTime, err := time.Parse("2006-01-02 15:04:05", checkInTimeStr.String)
		if err == nil {
			checkInTime = &parsedTime
		}
	}

	if checkOutTimeStr.Valid {
		parsedTime, err := time.Parse("2006-01-02 15:04:05", checkOutTimeStr.String)
		if err == nil {
			checkOutTime = &parsedTime
		}
	}

	var statusPtr *string
	if status.Valid {
		statusPtr = &status.String
	}

	// createdAt 파싱
	var parsedCreatedAt time.Time
	if createdAtStr.Valid {
		parsedTime, err := time.Parse("2006-01-02 15:04:05", createdAtStr.String)
		if err == nil {
			parsedCreatedAt = parsedTime
		} else {
			parsedCreatedAt = time.Now() // 기본값
		}
	} else {
		parsedCreatedAt = time.Now() // 기본값
	}

	record := models.AttendanceRecord{
		RecordID:     recordID,
		UserID:       dbUserID,
		WorkDate:     workDate,
		CheckInTime:  checkInTime,
		CheckOutTime: checkOutTime,
		Status:       statusPtr,
		CreatedAt:    parsedCreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(record)
}

// GetTodayAttendance 오늘의 출퇴근 기록 조회
func GetTodayAttendance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// URL에서 user_id 추출
	pathParts := strings.Split(r.URL.Path, "/")
	fmt.Println("[GetTodayAttendance] URL 경로:", r.URL.Path)
	fmt.Println("[GetTodayAttendance] 경로 파트:", pathParts)

	if len(pathParts) < 5 {
		fmt.Println("[GetTodayAttendance] 경로 파트 부족:", len(pathParts))
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	userID, err := strconv.Atoi(pathParts[4])
	if err != nil {
		fmt.Println("[GetTodayAttendance] user_id 파싱 오류:", err, "값:", pathParts[4])
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	fmt.Println("[GetTodayAttendance] 추출된 user_id:", userID)

	today := time.Now().Format("2006-01-02")

	var recordID int
	var dbUserID int
	var workDate string
	var checkInTimeStr sql.NullString
	var checkOutTimeStr sql.NullString
	var status sql.NullString
	var createdAtStr sql.NullString

	err = config.DB.QueryRow(
		"SELECT record_id, user_id, work_date, check_in_time, check_out_time, status, created_at FROM attendance_records WHERE user_id = ? AND work_date = ?",
		userID, today,
	).Scan(
		&recordID,
		&dbUserID,
		&workDate,
		&checkInTimeStr,
		&checkOutTimeStr,
		&status,
		&createdAtStr,
	)

	if err == sql.ErrNoRows {
		// 오늘 기록이 없으면 빈 응답
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"user_id":        userID,
			"work_date":      today,
			"check_in_time":  nil,
			"check_out_time": nil,
		})
		return
	} else if err != nil {
		fmt.Println("[GetTodayAttendance] DB 오류:", err)
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}

	// sql.NullString을 *time.Time으로 변환
	var checkInTime *time.Time
	var checkOutTime *time.Time

	if checkInTimeStr.Valid {
		parsedTime, err := time.Parse("2006-01-02 15:04:05", checkInTimeStr.String)
		if err == nil {
			checkInTime = &parsedTime
		}
	}

	if checkOutTimeStr.Valid {
		parsedTime, err := time.Parse("2006-01-02 15:04:05", checkOutTimeStr.String)
		if err == nil {
			checkOutTime = &parsedTime
		}
	}

	var statusPtr *string
	if status.Valid {
		statusPtr = &status.String
	}

	// createdAt 파싱
	var parsedCreatedAt time.Time
	if createdAtStr.Valid {
		parsedTime, err := time.Parse("2006-01-02 15:04:05", createdAtStr.String)
		if err == nil {
			parsedCreatedAt = parsedTime
		} else {
			parsedCreatedAt = time.Now() // 기본값
		}
	} else {
		parsedCreatedAt = time.Now() // 기본값
	}

	record := models.AttendanceRecord{
		RecordID:     recordID,
		UserID:       dbUserID,
		WorkDate:     workDate,
		CheckInTime:  checkInTime,
		CheckOutTime: checkOutTime,
		Status:       statusPtr,
		CreatedAt:    parsedCreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(record)
}

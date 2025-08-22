package service

import (
	"BACK/model"
	"BACK/repository"
	"fmt"
	"net"
	"strings"
	"time"
)

// AttendanceServiceImpl 출근 서비스 구현체
type AttendanceServiceImpl struct {
	repo repository.AttendanceRepository
}

// NewAttendanceService 출근 서비스 생성자
func NewAttendanceService(repo repository.AttendanceRepository) *AttendanceServiceImpl {
	return &AttendanceServiceImpl{repo: repo}
}

// GetAttendance 출근 기록 조회
func (s *AttendanceServiceImpl) GetAttendance(userID int, startDate, endDate time.Time) ([]model.Attendance, error) {
	return s.repo.FindByUserAndDateRange(userID, startDate, endDate)
}

// PostAttendance 서비스 구현 (IP 검증 포함)
func (s *AttendanceServiceImpl) PostAttendance(userID int, date time.Time, action string, timeStr string, clientIP string) error {
	fmt.Printf("[PostAttendance] %s 처리 - userID: %d, date: %s, IP: %s\n", action, userID, date.Format("2006-01-02"), clientIP)

	// 출근 시에만 IP 체크
	if action == "checkin" {
		if !s.isOfficeNetwork(clientIP) {
			fmt.Printf("[PostAttendance] 회사 네트워크가 아닌 IP에서 출근 시도: %s\n", clientIP)
			return fmt.Errorf("출근은 회사 Wi-Fi(172.30.1.x)에서만 가능합니다")
		}
	}

	switch action {
	case "checkin":
		checkInTime, err := s.parseTime(timeStr, date)
		if err != nil {
			return err
		}

		err = s.repo.CreateOrUpdate(userID, date, "checkin", checkInTime, "")
		if err != nil {
			fmt.Printf("[PostAttendance] 실패: %v\n", err)
			return err
		}
		fmt.Printf("[PostAttendance] 성공\n")
		return nil

	case "checkout":
		checkOutTime, err := s.parseTime(timeStr, date)
		if err != nil {
			return err
		}

		err = s.repo.CreateOrUpdate(userID, date, "checkout", checkOutTime, "")
		if err != nil {
			fmt.Printf("[PostAttendance] 실패: %v\n", err)
			return err
		}
		fmt.Printf("[PostAttendance] 성공\n")
		return nil

	default:
		return fmt.Errorf("지원하지 않는 액션: %s", action)
	}
}

// parseTime 시간 문자열을 time.Time으로 변환
func (s *AttendanceServiceImpl) parseTime(timeStr string, date time.Time) (time.Time, error) {
	if timeStr == "" {
		return time.Now(), nil
	}

	// ISO 형식인지 HH:MM 형식인지 확인
	if len(timeStr) > 10 && (timeStr[4] == '-' || timeStr[10] == 'T') {
		// ISO 형식 (2024-01-15T10:30:00.000Z)
		parsedTime, err := time.Parse(time.RFC3339, timeStr)
		if err != nil {
			return time.Time{}, fmt.Errorf("ISO 시간 형식 오류: %s", timeStr)
		}
		return parsedTime, nil
	} else {
		// HH:MM 형식
		parsedTime, err := time.Parse("15:04", timeStr)
		if err != nil {
			return time.Time{}, fmt.Errorf("시간 형식 오류: %s", timeStr)
		}

		// 날짜와 시간 결합
		result := time.Date(
			date.Year(), date.Month(), date.Day(),
			parsedTime.Hour(), parsedTime.Minute(), 0, 0,
			time.Local,
		)
		return result, nil
	}
}

// isOfficeNetwork 회사 네트워크인지 확인합니다
func (s *AttendanceServiceImpl) isOfficeNetwork(ip string) bool {
	fmt.Printf("IP 체크 시작: %s\n", ip)

	// 회사 Wi-Fi IP 범위 설정
	allowedNetworks := []string{
		"172.30.1.0/24", // 회사 Wi-Fi (192.168.0.1 ~ 192.168.0.254)
		//"127.0.0.1",      // 로컬호스트 (개발/테스트용)
		//"::1",            // IPv6 로컬호스트
	}

	clientIP := net.ParseIP(ip)
	if clientIP == nil {
		fmt.Printf("올바르지 않은 IP 형식: %s\n", ip)
		return false
	}

	for _, network := range allowedNetworks {
		// 정확한 IP 매치 (로컬호스트)
		if network == ip {
			fmt.Printf("정확한 IP 매치: %s\n", ip)
			return true
		}

		// CIDR 범위 체크
		if strings.Contains(network, "/") {
			_, ipNet, err := net.ParseCIDR(network)
			if err == nil && ipNet.Contains(clientIP) {
				fmt.Printf("네트워크 범위 매치: %s in %s\n", ip, network)
				return true
			}
		}
	}

	fmt.Printf("허용되지 않은 네트워크: %s\n", ip)
	return false
}

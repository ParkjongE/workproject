# 개요 (Overview)

이 개발 프로젝트는 조직 내 업무 협업의 효율성 증대와 인력 관리의 체계화를 동시에 달성하기 위해 설계된 통합 협업 관리 플랫폼입니다. 기존의 업무 협업 도구(Flow)와 인력 관리 솔루션(Shiftee)의 장점을 결합하여, 하나의 시스템에서 업무, 일정, 근태, 인사 관리를 모두 처리하는 올인원 솔루션을 제공합니다.

본 시스템의 최종 목표는 다음과 같습니다.

## 주요 목표

- **소통 활성화 및 접근성 향상**: 심리적 장벽 없이 자유롭게 소통하고 협업을 요청할 수 있는 환경을 조성하여, 팀원 간의 소통 부재 문제를 해소하고 업무 효율성을 높입니다.

- **스케줄 및 업무 투명성 확보**: 팀원들의 일정과 업무 진행 상황을 한눈에 파악하여 중복 업무를 방지하고, 자연스러운 협업 제안을 유도함으로써 불필요한 일정 조율 시간을 절감합니다.

- **지식 공유 및 업무 효율성 증대**: 과거 유사 업무 데이터를 업무 기록(History) 형태로 저장하고 관리하여, 불필요한 중복 작업을 방지하고 팀 전체의 지식 자산을 활용함으로써 업무 처리 효율을 극대화합니다.

- **통합적인 조직 및 인력 관리**: 출퇴근, 근태, 휴가 관리부터 회사 문서 및 규칙 관리까지 모든 인력 관리 기능을 통합하여, 관리자와 직원 모두 효율적인 시간 관리와 근무 현황 파악, 그리고 조직 운영 전반에 대한 통합적인 관리를 가능하게 합니다.

- **생산성 극대화**: 모든 기능이 하나의 플랫폼에 통합되어 구성원 간의 업무 연계성과 투명성을 확보하고, 시간 낭비를 최소화하여 궁극적으로 조직의 생산성을 극대화합니다.

# 프로젝트 파일 구조

## BACK/

```
├── main.go                    # 메인 서버 진입점 - 라우터 설정 및 서버 시작
├── go.mod                     # Go 모듈 의존성 관리
├── go.sum                     # Go 모듈 체크섬
├── BACK.exe                   # 컴파일된 실행 파일
├── db_setup.sql              # 데이터베이스 초기 설정 스크립트
├── JWT_README.md             # JWT 인증 시스템 사용법 문서
│
├── config/                    # 설정 관련
│   ├── db.go                 # 데이터베이스 연결 설정
│   └── session.go            # 세션 관리 설정  (사용x)
│
├── model/                     # 데이터 모델 (구조체)
│   ├── user.go               # 사용자 모델
│   ├── auth.go               # 인증 관련 모델
│   ├── attendance.go         # 출근/퇴근 모델
│   ├── leave.go              # 연차/휴가 모델
│   ├── schedule.go           # 일정 모델
│   ├── schedule_sharing.go   # 일정 공유 모델
│   ├── salary.go             # 급여 모델
│   └── department.go         # 부서 모델 
│
├── repository/                # 데이터 접근 계층 (DB 쿼리)
│   ├── register.go           # 사용자 등록 DB 작업
│   ├── attendance.go         # 출근/퇴근 DB 작업
│   ├── leave.go              # 연차/휴가 DB 작업
│   ├── schedule.go           # 일정 DB 작업
│   └── schedule_sharing.go   # 일정 공유 DB 작업
│
├── service/                   # 비즈니스 로직 계층
│   ├── register.go           # 사용자 등록 서비스
│   ├── attendance.go         # 출근/퇴근 서비스
│   ├── leave.go              # 연차/휴가 서비스
│   ├── schedule.go           # 일정 서비스
│   ├── schedule_sharing.go   # 일정 공유 서비스
│   └── salary.go             # 급여 서비스
│
├── handler/                   # HTTP 요청 처리 계층
│   ├── auth.go               # 기본 인증 처리
│   ├── jwt_auth.go           # JWT 인증 처리
│   ├── register.go           # 사용자 등록 처리
│   ├── attendance.go         # 출근/퇴근 처리
│   ├── leave.go              # 연차/휴가 처리
│   ├── schedule.go           # 일정 처리
│   ├── schedule_sharing.go   # 사용자 별 일정 공유 처리
│   ├── salary.go             # 급여명세서 처리 
│   └── department.go         # 부서 처리 (사용x)
│
├── middleware/                # HTTP 요청 전처리
│   ├── cors.go               # CORS 설정
│   ├── jwt.go                # JWT 인증 미들웨어
│   ├── dual_auth.go          # 이중 인증 미들웨어 (사용x)
│   └── session.go            # 세션 관리 미들웨어 (사용x)
│
├── utils/                     # 유틸리티 함수들
│   ├── jwt.go                # JWT 토큰 생성/검증
│   ├── password.go           # 비밀번호 해싱/검증
│   └── leave_calculator.go   # 연차 계산 로직
│
└── files/                     # 정적 파일 저장소
    └── salary/                # 급여명세서 PDF 파일들
        ├── [일반서식] 급여명세서.pdf
        ├── salary_statement_1_2024_1.pdf
        ├── salary_statement_1_2024_2.pdf
        └── salary_statement_1_2024_3.pdf
```

## FRONT/

```
├── vue-project/               # Vue.js 메인 프로젝트
│   ├── public/                # 정적 파일들
│   │   ├── favicon.ico        
│   │   └── index.html         # 메인 HTML 템플릿
│   │
│   ├── src/                   # 소스코드
│   │   ├── main.js            # Vue 앱 진입점
│   │   ├── App.vue            # 루트 컴포넌트
│   │   │
│   │   ├── components/        # 공통 컴포넌트
│   │   │   ├── SideBar.vue    # 사이드바 네비게이션
│   │   │   ├── AuthGuard.vue  # 인증 보호 컴포넌트 (사용x)
│   │   │   
│   │   │
│   │   ├── views/             # 페이지 컴포넌트들
│   │   │   ├── Login/         # 로그인/회원가입
│   │   │   │   ├── Login.vue
│   │   │   │   ├── Register.vue
│   │   │   │   ├── AuthGuard.vue
│   │   │   │   ├── Login.styled.js
│   │   │   │   └── Register.styled.js
│   │   │   │
│   │   │   ├── attendance/    # 출근/퇴근 관리
│   │   │   │   ├── Commute.vue
│   │   │   │   ├── MyAttendance.vue
│   │   │   │   ├── LeaveApply.vue
│   │   │   │   ├── Commute.styled.js
│   │   │   │   ├── MyAttendance.styled.js
│   │   │   │   └── LeaveApply.styled.js
│   │   │   │
│   │   │   ├── myschedule/    # 내 일정 관리
│   │   │   │   ├── Schedule/
│   │   │   │   │   ├── Schedule.vue
│   │   │   │   │   ├── Schedule.styled.js
│   │   │   │   │   ├── components/
│   │   │   │   │   │   ├── CalendarSection.vue
│   │   │   │   │   │   ├── ClockSection.vue
│   │   │   │   │   │   └── InfoCardSection.vue
│   │   │   │   │   ├── composables/
│   │   │   │   │   │   ├── useAttendance.js
│   │   │   │   │   │   ├── useCalendar.js
│   │   │   │   │   │   ├── useClock.js
│   │   │   │   │   │   └── useSchedule.js
│   │   │   │   │   ├── styled.js
│   │   │   │   │   └── utils/
│   │   │   │   │       └── formatters.js
│   │   │   │   │
│   │   │   │   ├── MeetingRequests/ # 협업/미팅 요청 관리
│   │   │   │   │   ├── MeetingRequests.vue
│   │   │   │   │   └── MeetingRequests.styled.js
│   │   │   │   │
│   │   │   │   └── WorkRecords/ # 내 업무 기록
│   │   │   │       ├── WorkRecords.vue
│   │   │   │       └── WorkRecords.styled.js
│   │   │   │
│   │   │   ├── organization/  # 조직도/부서 관리
│   │   │   │   ├── OrgChart.vue #조직도    
│   │   │   │   ├── OrgMember.vue #구성원 관리
│   │   │   │   ├── DepWork.vue #부서별 업무 현황 조회
│   │   │   │   ├── OrgChart.styled.js
│   │   │   │   ├── OrgMember.styled.js
│   │   │   │   └── DepWork.styled.js
│   │   │   │
│   │   │   └── schedulesharing/ # 일정 공유
│   │   │       ├── ScheduleSharing.vue #개인 일정 공유
│   │   │       ├── ScheduleSharing.styled.js
│   │   │       ├── InteractionGraph.vue #상호작용 그래프
│   │   │       ├── RelatedWork.vue
│   │   │       ├── TeamSchedule.vue #팀 일정 공유
│   │   │       ├── InteractionGraph.styled.js
│   │   │       ├── RelatedWork.styled.js
│   │   │       ├── TeamSchedule.styled.js
│   │   │       └── utils/
│   │   │           ├── drag.js
│   │   │           ├── Edit_schedules.js
│   │   │           ├── load.js
│   │   │           └── Modal.js
│   │   │
│   │   ├── router/            # 라우터 설정
│   │   │   ├── index.js       # 메인 라우터
│   │   │   ├── auth.js        # 인증 관련 라우트
│   │   │   ├── guards.js      # 라우트 가드
│   │   │   ├── attendance.js  # 출근 관련 라우트
│   │   │   ├── myschedule.js  # 일정 관련 라우트
│   │   │   ├── organization.js # 조직 관련 라우트
│   │   │   └── schedulesharing.js # 일정 공유 라우트
│   │   │
│   │   ├── stores/            # Pinia 상태 관리
│   │   │   ├── attendance.js  # 출근 상태 관리
│   │   │   └── schedule.js    # 일정 상태 관리 (사용x)
│   │   │
│   │   ├── plugins/           # 플러그인
│   │   │   ├── axios.js       # HTTP 클라이언트 설정
│   │   │   └── styled-components.js # 스타일 컴포넌트
│   │   │
│   │   └── assets/            # 정적 자산
│   │       ├── logo.png       # 로고 이미지
│   │       └── styles/        # 스타일 파일들
│   │
│   ├── package.json           # 프로젝트 의존성 및 스크립트
│   ├── package-lock.json      # 의존성 잠금 파일
│   ├── vue.config.js          # Vue CLI 설정
│   ├── babel.config.js        # Babel 설정
│   ├── jsconfig.json          # JavaScript 설정
│   ├── .gitignore             # Git 무시 파일
│   
│
│
└── img/                       # 이미지 파일들
```

Dump20250821 # mysql dump 파일
## 구현된 기능

### 1. 회원가입 및 인증 시스템

* **회원가입 기능**
    * **사용자 등록:** `handler/register.go`에서 회원가입 요청을 처리하고 사용자 정보를 데이터베이스에 저장
    * **비밀번호 보안:** `utils/password.go`에서 32바이트 랜덤 솔트를 생성하고 SHA256으로 비밀번호를 해싱하여 안전하게 저장
    * **사용자 모델:** `model/user.go`에 사용자 정보 구조체와 회원가입 요청/응답 모델이 정의되어 있음

* **로그인 인증**
    * **JWT 인증:** `handler/jwt_auth.go`에서 사용자 로그인 시 JWT 토큰을 생성하고 발급
    * **토큰 생성:** `utils/jwt.go`에서 사용자 ID, 사용자명, 역할 정보를 포함한 JWT 토큰을 생성
    * **비밀번호 검증:** `utils/password.go`에서 입력된 비밀번호와 저장된 솔트를 조합하여 해시값을 생성하고 데이터베이스의 해시값과 비교하여 인증

* **JWT 토큰 관리**
    * **토큰 발급:** 로그인 성공 시 사용자 정보를 담은 JWT 토큰을 클라이언트에 전달
    * **토큰 검증:** `middleware/jwt.go`에서 모든 보호된 API 요청 시 JWT 토큰의 유효성을 검증

* **미구현 기능 (연습 파일)**
    * **이중 인증 시스템:** `middleware/dual_auth.go`에 JWT와 세션을 모두 검증하는 이중 인증 로직이 구현되어 있지만 실제로는 사용되지 않음
    * **세션 기반 인증:** `middleware/session.go`에 세션 관리 기능이 구현되어 있지만 현재 JWT 인증만 사용하고 있어서 활용되지 않음
    * **세션 설정:** `config/session.go`에 세션 관련 설정이 되어 있지만 실제 세션 인증은 구현되지 않음

---

### 2. 출근/퇴근 관리 시스템

* **출근/퇴근 기록**
    * **체크인/체크아웃:** `handler/attendance.go`에서 출근과 퇴근 시간을 기록하고 데이터베이스에 저장
    * **출근 이력 관리:** `repository/attendance.go`에서 사용자별 출근 기록을 조회하고 관리
    * **비즈니스 로직:** `service/attendance.go`에서 출근/퇴근 관련 업무 규칙과 계산 로직을 처리

* **근무 시간 계산**
    * **자동 계산:** `service/attendance.go`에서 출근 시간과 퇴근 시간을 기반으로 일일 근무 시간을 자동으로 계산
    * **초과 근무 추적:** `service/attendance.go`에서 정상 근무 시간을 초과한 시간을 계산하여 초과 근무 시간을 추적
    * **통계 데이터:** `service/attendance.go`에서 일별, 주별, 월별 근무 시간 통계를 생성하여 사용자와 관리자가 근무 현황을 파악할 수 있도록 지원

* **위치 기반 출근 (WiFi 로직)**
    * **WiFi 기반 인증:** `handler/attendance.go`에서 사용자가 회사 WiFi에 연결되어 있을 때만 출근/퇴근이 가능하도록 위치 기반 인증 시스템 구현
    * **IP 범위 검증:** `service/attendance.go`에서 `isOfficeNetwork()` 함수를 통해 사용자의 IP 주소가 회사 WiFi 네트워크 범위 내에 있는지 검증
    * **네트워크 설정:** `service/attendance.go`에서 허용된 네트워크 IP 범위를 설정하여 회사 WiFi IP 대역(예: 172.30.1.0/24)에서만 출근이 가능하도록 제한

* **데이터 모델**
    * **출근 기록:** `model/attendance.go`에 출근/퇴근 시간, 근무 시간, 위치 정보 등을 포함한 출근 기록 구조체 정의
    * **사용자별 관리:** `repository/attendance.go`에서 각 사용자의 출근 기록을 개별적으로 관리
    * **WiFi 로직 세부 설명:** `service/attendance.go`의 `isOfficeNetwork()` 함수에서 회사 WiFi IP 범위를 설정해야 하며, 현재는 172.30.1.0/24 대역으로 설정되어 있습니다. 실제 운영 환경에서는 회사의 실제 WiFi IP 대역으로 수정해야 합니다.

---

### 3. 본인 일정 수정 CRUD 시스템

* **일정 생성 (Create)**
    * **일정 등록:** `handler/schedule.go`에서 `CreateSchedule()` 함수로 새로운 일정을 생성하고 데이터베이스에 저장
    * **입력값 검증:** `service/schedule.go`에서 일정 제목과 사용자 ID가 필수 입력값인지 검증하는 비즈니스 로직 처리
    * **데이터 저장:** `repository/schedule.go`에서 일정 제목, 내용, 시작/종료 시간, 색상 등의 정보를 데이터베이스에 INSERT

* **일정 조회 (Read)**
    * **사용자별 조회:** `handler/schedule.go`에서 `GetSchedulesByUser()` 함수로 특정 사용자의 모든 일정을 조회
    * **데이터 정렬:** `repository/schedule.go`에서 사용자 ID를 기준으로 일정을 조회하고 시작 시간 순으로 정렬하여 반환
    * **일정 목록:** `model/schedule.go`에 일정 ID, 제목, 내용, 시작/종료 시간, 색상, 생성일 등의 일정 정보 구조체 정의

* **일정 수정 (Update)**
    * **일정 업데이트:** `handler/schedule.go`에서 `UpdateSchedule()` 함수로 기존 일정 정보를 수정
    * **수정 검증:** `service/schedule.go`에서 수정할 일정 ID와 제목이 유효한지 검증하는 비즈니스 로직 처리
    * **데이터 갱신:** `repository/schedule.go`에서 일정 제목, 내용, 시작/종료 시간, 색상 등의 정보를 UPDATE로 갱신

* **일정 삭제 (Delete)**
    * **일정 제거:** `handler/schedule.go`에서 `DeleteSchedule()` 함수로 특정 일정을 삭제
    * **삭제 검증:** `service/schedule.go`에서 삭제할 일정 ID가 유효한지 검증하는 비즈니스 로직 처리
    * **데이터 삭제:** `repository/schedule.go`에서 일정 ID를 기준으로 해당 일정을 데이터베이스에서 DELETE

---

### 4. 연차현황 조회 시스템

* **연차 현황 조회**
    * **개인 연차 현황:** `handler/leave.go`에서 `GetUserLeaveStatus()` 함수로 사용자의 총 연차, 사용한 연차, 잔여 연차를 조회
    * **연차 계산 로직:** `utils/leave_calculator.go`에서 `CalculateUserLeaveStatus()` 함수로 입사일을 기준으로 근무 연수에 따른 연차를 자동 계산
    * **연차 정책:** `utils/leave_calculator.go`에서 기본 15일, 1년 이상 근무 시 2일씩 추가, 최대 25일까지 부여하는 연차 정책 정의

* **연차 계산 및 업데이트**
    * **특정 날짜 연차 계산:** `handler/leave.go`에서 `CalculateLeaveForDate()` 함수로 특정 날짜 기준의 연차 현황을 계산
    * **전체 사용자 연차 업데이트:** `handler/leave.go`에서 `UpdateAllUsersLeaveStatus()` 함수로 모든 사용자의 연차 현황을 일괄 업데이트
    * **근무 연수 계산:** `utils/leave_calculator.go`에서 `CalculateYearsOfService()` 함수로 입사일부터 현재까지의 근무 연수를 계산

* **연차 데이터 모델**
    * **연차 현황 구조체:** `utils/leave_calculator.go`에 사용자 ID, 총 부여 연차, 사용한 연차, 잔여 연차를 포함한 `LeaveStatus` 구조체 정의
    * **연차 정책 구조체:** `utils/leave_calculator.go`에 최대 연차 일수, 최소 근무 연수, 연차 증가량을 포함한 `LeavePolicy` 구조체 정의

---

### 5. 급여명세서 관리 시스템

* **급여명세서 파일 관리**
    * **파일 다운로드:** `handler/salary.go`에서 `GetSalaryStatement()` 함수로 특정 사용자의 연도/월별 급여명세서 PDF 파일을 다운로드
    * **파일 목록 조회:** `handler/salary.go`에서 `GetSalaryStatementList()` 함수로 사용자의 모든 급여명세서 목록을 조회
    * **파일 검증:** `service/salary.go`에서 사용자 ID, 연도, 월의 유효성을 검증하고 파일 존재 여부를 확인

* **급여 파일 처리**
    * **파일 경로 관리:** `service/salary.go`에서 `files/salary/` 디렉토리 내의 급여명세서 파일들을 체계적으로 관리
    * **파일명 규칙:** `service/salary.go`에서 `salary_statement_{userID}_{year}_{month}.pdf` 형식의 파일명 규칙을 적용
    * **파일 정보 구조체:** `model/salary.go`에 파일명, 파일 경로, 파일 크기, 콘텐츠 타입을 포함한 `SalaryFile` 구조체 정의

* **급여 데이터 모델**
    * **급여명세서 정보:** `model/salary.go`에 사용자 ID, 연도, 월, 파일명, 파일 경로를 포함한 `SalaryStatement` 구조체 정의
    * **응답 구조체:** `model/salary.go`에 급여명세서 목록, 파일 다운로드, 에러 응답을 위한 다양한 응답 구조체 정의
    * **파일 저장소:** `files/salary/` 디렉토리에 실제 급여명세서 PDF 파일들을 저장하고 관리

---

### 6. 일정 공유 시스템

* **부서별 사용자 조회**
    * **부서별 사용자 목록:** `handler/schedule_sharing.go`에서 `GetUsersByDepartment()` 함수로 특정 부서(개발부, 디자인부, 마케팅부)의 모든 사용자 정보를 조회
    * **역할별 정렬:** `repository/schedule_sharing.go`에서 부장, 차장, 과장, 대리, 사원 순으로 역할에 따른 계급별 정렬을 적용하여 사용자 목록을 반환
    * **부서 유효성 검증:** `handler/schedule_sharing.go`에서 허용된 부서명만 조회 가능하도록 유효성 검사를 수행

* **전체 부서 사용자 조회**
    * **전체 사용자 통합:** `handler/schedule_sharing.go`에서 `GetAllDepartmentUsers()` 함수로 모든 부서의 사용자 정보를 통합하여 조회
    * **부서별 그룹화:** `service/schedule_sharing.go`에서 각 부서별로 사용자를 그룹화하고 모든 사용자를 하나의 배열로 합쳐서 반환
    * **일정 정보 포함:** `repository/schedule_sharing.go`에서 각 사용자의 일정 정보도 함께 조회하여 사용자와 일정을 연계하여 제공

* **개별 사용자 일정 조회**
    * **사용자별 일정:** `handler/schedule_sharing.go`에서 `GetUserSchedulesByID()` 함수로 특정 사용자의 모든 일정을 조회
    * **일정 상세 정보:** `repository/schedule_sharing.go`에서 일정 ID, 제목, 시작/종료 시간, 색상, 생성일 등의 일정 상세 정보를 제공
    * **빈 일정 처리:** `service/schedule_sharing.go`에서 일정이 없는 사용자의 경우 빈 배열을 반환하여 일관된 응답 구조 유지

* **일정 공유 데이터 모델**
    * **사용자 정보 구조체:** `model/schedule_sharing.go`에 사용자 ID, 이름, 역할, 부서, 입사일을 포함한 `ScheduleSharingUser` 구조체 정의
    * **일정 포함 사용자:** `model/schedule_sharing.go`에 사용자 정보와 해당 사용자의 일정 목록을 포함한 `ScheduleSharingUserWithSchedules` 구조체 정의
    * **부서별 사용자 그룹:** `model/schedule_sharing.go`에 부서명과 해당 부서의 사용자 목록을 포함한 `DepartmentUsers` 구조체 정의
    * **API 응답 구조:** `model/schedule_sharing.go`에 성공/실패 여부, 메시지, 데이터, 에러 정보를 포함한 다양한 응답 구조체 정의

* **일정 공유 비즈니스 로직**
    * **사용자-일정 연계:** `service/schedule_sharing.go`에서 사용자 정보를 조회한 후 각 사용자의 일정을 개별적으로 조회하여 통합된 정보 제공
    * **에러 처리:** `service/schedule_sharing.go`에서 일정 조회에 실패한 경우에도 사용자 정보는 반환하여 부분적인 데이터라도 제공

---

## FRONT 프로젝트 화면 구조 및 기능 설명

### 메인 화면 구조

1.  **메인 화면 (`ScheduleSharing.vue`)**
    * **위치:** `FRONT/vue-project/src/views/schedulesharing/ScheduleSharing.vue`
    * **역할:** 전체 시스템의 메인 대시보드 역할
    * **구성:**
        * 부서별 경계 (개발부, 디자인부, 마케팅부)
        * 본인 노드 (고정 위치)
        * 구성원 노드들 (드래그x)

2.  **본인 일정 모달 (`Schedule.vue`)**
    * **위치:** `FRONT/vue-project/src/views/myschedule/Schedule/Schedule.vue`
    * **호출 방식:** 메인 화면에서 본인 노드 클릭 시 모달로 표시
    * **구조:** 시계, 위젯, 달력으로 기능별 분리

3.  **`Schedule.vue` 컴포넌트 분리 구조**
    * **시계 섹션 (`ClockSection.vue`)**
        * **위치:** `FRONT/vue-project/src/views/myschedule/Schedule/components/ClockSection.vue`
        * **기능:** 시간 기반 일정 표시 및 관리
    * **달력 섹션 (`CalendarSection.vue`)**
        * **위치:** `FRONT/vue-project/src/views/myschedule/Schedule/components/CalendarSection.vue`
        * **기능:** 월별/주별 일정 캘린더 뷰
    * **정보 카드 섹션 (`InfoCardSection.vue`)**
        * **위치:** `FRONT/vue-project/src/views/myschedule/Schedule/components/InfoCardSection.vue`
        * **기능:** 일정 요약 정보, 출근 정보, 공개 설정 위젯

---

### 🔗 백엔드 연동 완료된 페이지들

1.  **출근/퇴근 관리 (`Commute.vue`)**
    * **위치:** `FRONT/vue-project/src/views/attendance/Commute.vue`
    * **연동 상태:** ✅ 완전 연동 완료
    * **기능:**
        * 출근/퇴근 체크
        * 주간 근무 현황 표시
        * 근태 이력 조회
        * 근태 정정요청

2.  **본인 일정 관리 (`Schedule.vue`)**
    * **위치:** `FRONT/vue-project/src/views/myschedule/Schedule/Schedule.vue`
    * **연동 상태:** ✅ 완전 연동 완료
    * **기능:**
        * 개인 일정 CRUD
        * 시계/달력 뷰
        * 출근 정보 연동 (InfoCardSection)
        * 급여명세서 첨부 기능

3.  **일정 공유 메인 (`ScheduleSharing.vue`)**
    * **위치:** `FRONT/vue-project/src/views/schedulesharing/ScheduleSharing.vue`
    * **연동 상태:** ✅ 부분 연동 완료
    * **기능:**
        * 부서별 사용자 조회
        * 일정 공유
        * 미니 시계/달력 표시

---

### 📁 `ScheduleSharing.vue` 파일 분리 현황

* **완료된 파일 분리**
    * **드래그 기능:** `FRONT/vue-project/src/views/schedulesharing/utils/drag.js`
        * `handleDragStart`, `handleDragEnd`, `handleDragOver`, `handleNodeDrop`, `updateNodePosition`
    * **일정 수정:** `FRONT/vue-project/src/views/schedulesharing/utils/Edit_schedules.js`
        * `editScheduleRequest`, `updateScheduleRequest`, `deleteScheduleRequest`, `isMySchedule`
    * **데이터 로딩:** `FRONT/vue-project/src/views/schedulesharing/utils/load.js`
        * `loadMyProfile`, `loadTeamMembers`, `loadUsersByDepartment`, `loadUserSchedules`
    * **모달 관리:** `FRONT/vue-project/src/views/schedulesharing/utils/Modal.js`
        * `showClockView`, `showCalendarView`, `expandMemberSchedule`, `openQuickMeetingModal`

---

### ❌ 미구현 기능들

1.  **일정 요청 시스템**
    * **기능:** 다른 사용자에게 일정 요청 보내기
    * **상태:** 아직 구현되지 않음
    * **영향:** 일정 협업 기능 제한

2.  **사용자별 출퇴근 기록 연동**
    * **기능:** 각 사용자의 출퇴근 기록을 메인 화면에서 조회
    * **상태:** 아직 구현되지 않음
    * **영향:** 팀원 근무 현황 파악 불가

3.  **노드 드래그 기능**
    * **기능:** 사용자 노드들을 자유롭게 드래그하여 위치 이동
    * **상태:** 아직 구현되지 않음
    * **현재 상황:** `utils/drag.js`에 드래그 관련 함수들이 정의되어 있지만 실제 노드 드래그 기능은 작동하지 않음

---

### ✅ 구현 완료된 기능들

1.  **사용자별 일정 확인**
    * **시계 뷰:** 각 사용자의 일정을 시계 형태로 표시
    * **달력 뷰:** 각 사용자의 일정을 달력 형태로 표시
    * **미니 시계:** 축소된 시계에서도 일정 표시 (`MiniClockView`)
    * **미니 달력:** 축소된 달력에서도 일정 표시 (`MiniCalendarView`)

2.  **데이터 연동**
    * **부서별 사용자 조회:** 개발부, 디자인부, 마케팅부 사용자 정보
    * **일정 데이터 연동:** 백엔드 API와 연동하여 실시간 일정 표시
    * **사용자 프로필:** 각 사용자의 기본 정보 및 역할 표시

3.  **시계 드래그 투 셀렉트 기능**
    * **시간 설정:** 시계에서 드래그하여 15분 간격으로 시간 범위 선택 가능
    * **일정 추가:** 선택된 시간 범위에 일정을 추가하면 자동으로 달력과 일정 목록에 표시
    * **날짜 연동:** 해당 날짜를 클릭하면 선택된 날짜의 일정을 조회할 수 있음
    * **실시간 업데이트:** 일정 추가/수정 시 시계, 달력, 일정 목록이 동기화되어 표시

4.  **일정 공유 및 모달 시스템**
    * **노드 클릭 모달:** `ScheduleSharing.vue` 페이지에서 사용자 노드를 클릭하면 해당 사용자의 일정을 모달로 표시
    * **일정 공유:** 각 사용자의 일정을 다른 구성원들과 공유하여 팀 전체의 일정 현황을 파악 가능
    * **모달 통합:** 시계 뷰, 달력 뷰, 정보 카드가 통합된 모달로 사용자별 일정을 상세하게 확인







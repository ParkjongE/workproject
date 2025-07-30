# JWT 인증 시스템

이 프로젝트는 세션 기반 인증과 JWT 기반 인증을 모두 지원합니다.

## 기능

- **JWT 토큰 생성 및 검증**
- **토큰 갱신**
- **역할 기반 접근 제어 (RBAC)**
- **선택적 인증 미들웨어**
- **CORS 지원**

## 설치 및 설정

### 1. 의존성 설치

```bash
go get github.com/golang-jwt/jwt/v5
```

### 2. 환경변수 설정

`.env` 파일에 JWT 시크릿 키를 설정하세요:

```env
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
```

시크릿 키가 설정되지 않으면 기본값이 사용됩니다.

## API 엔드포인트

### 인증 API

#### 1. JWT 로그인
```
POST /api/jwt/login
Content-Type: application/json

{
  "username": "사용자명",
  "password": "비밀번호"
}
```

**응답:**
```json
{
  "success": true,
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "사용자명",
    "name": "사용자 이름",
    "role": "user",
    "department": "부서명"
  },
  "message": "JWT 로그인 성공"
}
```

#### 2. JWT 로그아웃
```
POST /api/jwt/logout
```

**응답:**
```json
{
  "success": true,
  "message": "로그아웃 성공. 클라이언트에서 토큰을 삭제하세요."
}
```

#### 3. 토큰 검증
```
GET /api/jwt/verify
Authorization: Bearer <token>
```

**응답:**
```json
{
  "success": true,
  "valid": true,
  "user": {
    "user_id": 1,
    "username": "사용자명",
    "role": "user"
  },
  "message": "토큰이 유효합니다"
}
```

#### 4. 토큰 갱신
```
POST /api/jwt/refresh
Authorization: Bearer <token>
```

**응답:**
```json
{
  "success": true,
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "message": "토큰 갱신 성공"
}
```

### 보호된 API

#### 1. 사용자 프로필 (JWT 보호)
```
GET /api/jwt/user/profile
Authorization: Bearer <token>
```

#### 2. 관리자 대시보드 (관리자 역할 필요)
```
GET /api/jwt/admin/dashboard
Authorization: Bearer <token>
```

## 미들웨어 사용법

### 1. JWT 인증 미들웨어

```go
import "BACK/middleware"

// JWT 인증이 필요한 핸들러
handler := middleware.JWTAuthMiddleware(yourHandler)
```

### 2. 역할 기반 인증

```go
// 관리자 역할이 필요한 핸들러
adminHandler := middleware.RequireRole("admin")(yourHandler)
```

### 3. 선택적 JWT 인증

```go
// JWT 토큰이 있으면 검증하고, 없어도 통과
optionalHandler := middleware.OptionalJWTAuthMiddleware(yourHandler)
```

### 4. CORS와 JWT 인증 조합

```go
// CORS와 JWT 인증을 함께 적용
combinedHandler := middleware.CORSWithJWTAuth(yourHandler)
```

## 컨텍스트에서 사용자 정보 가져오기

```go
import "BACK/middleware"

func yourHandler(w http.ResponseWriter, r *http.Request) {
    userID, username, role, ok := middleware.GetUserFromContext(r.Context())
    if !ok {
        http.Error(w, "User not authenticated", http.StatusUnauthorized)
        return
    }
    
    // 사용자 정보 사용
    fmt.Printf("User ID: %d, Username: %s, Role: %s\n", userID, username, role)
}
```

## 클라이언트 사용법

### JavaScript/TypeScript

```javascript
// JWT 로그인
const loginResponse = await fetch('/api/jwt/login', {
    method: 'POST',
    headers: {
        'Content-Type': 'application/json',
    },
    body: JSON.stringify({
        username: '사용자명',
        password: '비밀번호'
    })
});

const { token } = await loginResponse.json();

// 보호된 API 호출
const profileResponse = await fetch('/api/jwt/user/profile', {
    headers: {
        'Authorization': `Bearer ${token}`
    }
});
```

### cURL

```bash
# JWT 로그인
curl -X POST http://localhost:8080/api/jwt/login \
  -H "Content-Type: application/json" \
  -d '{"username":"사용자명","password":"비밀번호"}'

# 보호된 API 호출
curl -X GET http://localhost:8080/api/jwt/user/profile \
  -H "Authorization: Bearer <token>"
```

## 보안 고려사항

1. **JWT 시크릿 키**: 프로덕션 환경에서는 강력한 시크릿 키를 사용하세요.
2. **토큰 만료**: 기본적으로 24시간으로 설정되어 있습니다. 필요에 따라 조정하세요.
3. **HTTPS**: 프로덕션 환경에서는 반드시 HTTPS를 사용하세요.
4. **토큰 저장**: 클라이언트에서는 안전한 방법으로 토큰을 저장하세요.

## 테스트

`test_jwt.html` 파일을 브라우저에서 열어 JWT 인증 기능을 테스트할 수 있습니다.

## 기존 세션 인증과의 차이점

| 기능 | 세션 인증 | JWT 인증 |
|------|-----------|----------|
| 서버 상태 | 상태 유지 | 무상태 |
| 확장성 | 제한적 | 높음 |
| 클라이언트 저장 | 쿠키 | 로컬 스토리지/쿠키 |
| 토큰 크기 | 작음 | 큼 |
| 무효화 | 즉시 가능 | 만료까지 대기 |

## 문제 해결

### 1. 토큰 만료 오류
- 토큰 갱신 API를 사용하여 새로운 토큰을 발급받으세요.

### 2. 권한 부족 오류
- 사용자의 역할을 확인하고 필요한 권한이 있는지 확인하세요.

### 3. CORS 오류
- 미들웨어 순서를 확인하세요. CORS 미들웨어가 JWT 미들웨어보다 먼저 적용되어야 합니다. 
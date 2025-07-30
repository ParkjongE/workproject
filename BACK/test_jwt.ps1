# JWT 로그인 테스트 스크립트

Write-Host "JWT 로그인 테스트 시작..." -ForegroundColor Green

# 서버가 실행 중인지 확인
try {
    $response = Invoke-WebRequest -Uri "http://localhost:8080/api/jwt/login" -Method POST -ContentType "application/json" -Body '{"username":"admin","password":"admin123"}' -TimeoutSec 5
    Write-Host "응답 상태: $($response.StatusCode)" -ForegroundColor Yellow
    
    if ($response.StatusCode -eq 200) {
        $data = $response.Content | ConvertFrom-Json
        Write-Host "로그인 성공!" -ForegroundColor Green
        Write-Host "토큰: $($data.token)" -ForegroundColor Cyan
        Write-Host "사용자: $($data.user.name)" -ForegroundColor Cyan
        Write-Host "역할: $($data.user.role)" -ForegroundColor Cyan
        
        # 토큰을 파일에 저장
        $data.token | Out-File -FilePath "jwt_token.txt" -Encoding UTF8
        Write-Host "토큰이 jwt_token.txt 파일에 저장되었습니다." -ForegroundColor Green
        
        # 토큰 검증 테스트
        Write-Host "`n토큰 검증 테스트..." -ForegroundColor Yellow
        $headers = @{
            "Authorization" = "Bearer $($data.token)"
        }
        
        $verifyResponse = Invoke-WebRequest -Uri "http://localhost:8080/api/jwt/verify" -Method GET -Headers $headers -TimeoutSec 5
        if ($verifyResponse.StatusCode -eq 200) {
            $verifyData = $verifyResponse.Content | ConvertFrom-Json
            Write-Host "토큰 검증 성공!" -ForegroundColor Green
            Write-Host "사용자 ID: $($verifyData.user.user_id)" -ForegroundColor Cyan
            Write-Host "사용자명: $($verifyData.user.username)" -ForegroundColor Cyan
        }
        
    } else {
        Write-Host "로그인 실패: $($response.StatusCode)" -ForegroundColor Red
        Write-Host "응답: $($response.Content)" -ForegroundColor Red
    }
} catch {
    Write-Host "오류 발생: $($_.Exception.Message)" -ForegroundColor Red
    Write-Host "서버가 실행 중인지 확인하세요." -ForegroundColor Yellow
}

Write-Host "`n테스트 완료!" -ForegroundColor Green 
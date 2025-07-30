<template>
  <div class="auth-guard-container">
    <div class="auth-guard-box">
      <div class="auth-guard-icon">🔒</div>
      <h2>인증이 필요합니다</h2>
      <p>이 페이지에 접근하려면 로그인이 필요합니다.</p>
      <button @click="goToLogin" class="login-btn">로그인하기</button>
    </div>
  </div>
</template>

<script>
export default {
  name: 'AuthGuard',
  methods: {
    goToLogin() {
      this.$router.push('/');
    }
  }
}
</script>

<style scoped>
.auth-guard-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: #F4F5F5;
  font-family: 'Comfortaa', 'Corben', sans-serif;
}

.auth-guard-box {
  background: #fff;
  border-radius: 10px;
  box-shadow: 0px 4px 4px 0px rgba(0,0,0,0.25);
  padding: 60px 80px;
  text-align: center;
  max-width: 500px;
}

.auth-guard-icon {
  font-size: 64px;
  margin-bottom: 20px;
}

.auth-guard-box h2 {
  font-family: 'Corben', sans-serif;
  font-size: 24px;
  color: #757678;
  margin-bottom: 16px;
}

.auth-guard-box p {
  font-family: 'Comfortaa', sans-serif;
  font-size: 16px;
  color: #757678;
  margin-bottom: 30px;
}

.login-btn {
  background: #218838;
  color: #fff;
  border: none;
  border-radius: 5px;
  padding: 12px 24px;
  font-size: 18px;
  font-family: 'Comfortaa', sans-serif;
  cursor: pointer;
  transition: background 0.2s;
}

.login-btn:hover {
  background: #17692a;
}
</style> 
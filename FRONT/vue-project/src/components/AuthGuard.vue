<template>
  <div>
    <slot v-if="isAuthenticated"></slot>
    <div v-else>
      <p>로그인이 필요합니다. 로그인 페이지로 이동합니다...</p>
    </div>
  </div>
</template>

<script>
import api from '@/plugins/axios'

export default {
  name: 'AuthGuard',
  data() {
    return {
      isAuthenticated: false
    }
  },
  created() {
    this.checkAuth();
  },
  methods: {
    async checkAuth() {
      const user = localStorage.getItem('user');
      if (!user) {
        this.$router.push('/');
        return;
      }
      
      try {
        // 세션 유효성 검증을 위해 서버에 요청
        await api.get('/api/user/profile');
        this.isAuthenticated = true;
      } catch (error) {
        // 세션이 유효하지 않으면 로그인 페이지로 이동
        localStorage.removeItem('user');
        localStorage.removeItem('jwt_token');
        this.$router.push('/');
      }
    }
  }
}
</script> 
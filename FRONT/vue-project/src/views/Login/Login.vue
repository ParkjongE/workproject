<template>
  <LoginContainer>
    <Logo>project</Logo>
    <LoginTitle>로그인</LoginTitle>
    <div v-if="$route.query.reason === 'auth'" style="color: #d32f2f; margin-bottom: 12px; font-weight: bold;">로그인이 필요합니다.</div>
    <form id="loginForm" @submit.prevent="handleLogin">
      <InputLabel for="userid">아이디</InputLabel>
      <LoginInput type="text" id="userid" :value="userid" @input="userid = $event.target.value" placeholder="아이디를 입력하세요" autocomplete="username" />

      <InputLabel for="password">비밀번호</InputLabel>
      <LoginInput type="password" id="password" :value="password" @input="password = $event.target.value" placeholder="비밀번호를 입력하세요" autocomplete="current-password" />

      <LoginOptions>
        <CheckboxLabel>
          <input type="checkbox" v-model="rememberMe" style="width:16px;height:16px;margin:0;" />
          <span class="checkmark"></span>
          로그인 상태 유지
        </CheckboxLabel>
        <ForgotPassword href="#">비밀번호 찾기</ForgotPassword>
      </LoginOptions>

      <ErrorMessage v-if="errorMessage">{{ errorMessage }}</ErrorMessage>
      <LoginBtn type="submit" :disabled="isLoading">{{ isLoading ? '로그인 중...' : '로그인' }}</LoginBtn>
    </form>
    <LoginLinks>
      <span>계정이 없으신가요?</span>
      <router-link to="/register">회원가입</router-link>
    </LoginLinks>
  </LoginContainer>
</template>

<script>
import {
  LoginContainer,
  Logo,
  LoginTitle,
  InputLabel,
  LoginInput,
  LoginOptions,
  CheckboxLabel,
  ForgotPassword,
  LoginBtn,
  ErrorMessage,
  LoginLinks
} from './Login.styled.js';

import api from '@/plugins/axios';

export default {
  name: 'LoginPage',
  components: {
    LoginContainer,
    Logo,
    LoginTitle,
    InputLabel,
    LoginInput,
    LoginOptions,
    CheckboxLabel,
    ForgotPassword,
    LoginBtn,
    ErrorMessage,
    LoginLinks
  },
  data() {
    return {
      userid: '',
      password: '',
      rememberMe: false,
      isLoading: false,
      errorMessage: ''
    }
  },
  methods: {
    async handleLogin() {
      if (!this.userid || !this.password) {
        this.errorMessage = '아이디와 비밀번호를 입력해주세요.';
        return;
      }
      this.isLoading = true;
      this.errorMessage = '';
      
      try {
        const response = await api.post('/api/login', {
          userid: String(this.userid),
          password: String(this.password)
        });
        
        console.log('로그인 응답:', response.data);
        
        if (response.data.success) {
          console.log('로그인 성공, 토큰 저장 시작');
          localStorage.setItem('user', JSON.stringify(response.data.user));
          
          // Access Token 저장
          if (response.data.access_token) {
            console.log('JWT 토큰 저장:', response.data.access_token.substring(0, 20) + '...');
            localStorage.setItem('jwt_token', response.data.access_token);
          } else {
            console.error('access_token이 응답에 없습니다!');
          }
          
          if (this.rememberMe) {
            localStorage.setItem('rememberMe', 'true');
          }
          
          console.log('저장된 토큰 확인:', localStorage.getItem('jwt_token'));
          this.$router.push('/commute');
        } else {
          this.errorMessage = response.data.error || '로그인에 실패했습니다.';
        }
      } catch (error) {
        if (error.response && error.response.data && error.response.data.error) {
          this.errorMessage = error.response.data.error;
        } else {
          this.errorMessage = '로그인에 실패했습니다. 다시 시도해주세요.';
        }
      } finally {
        this.isLoading = false;
      }
    }
  }
}
</script>

<!-- 스타일은 vue-styled-components로 분리됨 --> 
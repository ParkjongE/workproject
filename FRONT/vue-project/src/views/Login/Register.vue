<template>
  <RegisterContainer>
    <Logo>project</Logo>
    <RegisterTitle>회원가입</RegisterTitle>
    <form id="registerForm" @submit.prevent="handleRegister">
      <InputLabel for="userid">아이디</InputLabel>
      <RegisterInput type="text" id="userid" v-model="userid" placeholder="아이디를 입력하세요" autocomplete="username" />

      <InputLabel for="password">비밀번호</InputLabel>
      <RegisterInput type="password" id="password" v-model="password" placeholder="비밀번호를 입력하세요" autocomplete="new-password" />

      <InputLabel for="confirmPassword">비밀번호 확인</InputLabel>
      <RegisterInput type="password" id="confirmPassword" v-model="confirmPassword" placeholder="비밀번호를 다시 입력하세요" autocomplete="new-password" />

      <InputLabel for="name">이름</InputLabel>
      <RegisterInput type="text" id="name" v-model="name" placeholder="이름을 입력하세요" />

      <InputLabel for="department">부서</InputLabel>
      <RegisterSelect id="department" v-model="department">
        <option value="">부서 선택</option>
        <option v-for="dept in departments" :key="dept" :value="dept">{{ dept }}</option>
      </RegisterSelect>

      <InputLabel for="hireDate">입사날짜</InputLabel>
      <RegisterInput 
        type="date" 
        id="hireDate" 
        v-model="hireDate" 
        required 
      />

      <ErrorMessage v-if="errorMessage">{{ errorMessage }}</ErrorMessage>
      <RegisterBtn type="submit" :disabled="isLoading">{{ isLoading ? '처리 중...' : '회원가입' }}</RegisterBtn>
    </form>
    <RegisterLinks>
      <span>이미 계정이 있으신가요?</span>
      <router-link to="/">로그인</router-link>
    </RegisterLinks>
  </RegisterContainer>
</template>

<script>
import {
  RegisterContainer,
  Logo,
  RegisterTitle,
  InputLabel,
  RegisterInput,
  RegisterSelect,
  RegisterBtn,
  ErrorMessage,
  RegisterLinks
} from './Register.styled.js';


export default {
  name: 'RegisterView',
  components: {
    RegisterContainer,
    Logo,
    RegisterTitle,
    InputLabel,
    RegisterInput,
    RegisterSelect,
    RegisterBtn,
    ErrorMessage,
    RegisterLinks
  },
  data() {
    return {
      userid: '',
      password: '',
      confirmPassword: '',
      name: '',
      department: '',
      hireDate: '',
      departments: ['개발팀', '디자인팀', '마케팅팀'],
      errorMessage: '',
      isLoading: false
    }
  },

  methods: {
    async handleRegister() {
      if (!this.userid || !this.password || !this.name || !this.hireDate) {
        this.errorMessage = '아이디, 비밀번호, 이름, 입사날짜는 필수 입력 항목입니다.';
        return;
      }
      if (this.password !== this.confirmPassword) {
        this.errorMessage = '비밀번호가 일치하지 않습니다.';
        return;
      }
      if (this.password.length < 6) {
        this.errorMessage = '비밀번호는 최소 6자 이상이어야 합니다.';
        return;
      }
      this.isLoading = true;
      this.errorMessage = '';
      const userData = {
        userid: this.userid,
        password: this.password,
        name: this.name,
        role: '사원',
        department: this.department || "",
        hire_date: this.hireDate
      };
      try {
        const response = await fetch('/api/register', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify(userData)
        });
        const responseText = await response.text();
        try {
          JSON.parse(responseText);
        } catch (e) {
          // JSON 파싱 실패 시 무시
        }
        if (response.ok) {
          alert('회원가입이 완료되었습니다. 로그인해주세요.');
          this.$router.push('/');
        } else {
          this.errorMessage = responseText || '회원가입에 실패했습니다.';
        }
      } catch (error) {
        if (error.response) {
          this.errorMessage = error.response.data.error || `회원가입 실패 (${error.response.status})`;
        } else if (error.request) {
          this.errorMessage = '서버로부터 응답이 없습니다. 서버가 실행 중인지 확인하세요.';
        } else {
          this.errorMessage = `요청 오류: ${error.message}`;
        }
      } finally {
        this.isLoading = false;
      }
    }
  }
}
</script>

<!-- 스타일은 vue-styled-components로 분리됨 -->
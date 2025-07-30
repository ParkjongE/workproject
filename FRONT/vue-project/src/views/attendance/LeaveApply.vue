<template>
  <Container>
    <SideBar />
    <MainContent>
      <Title>휴가/조퇴/외출 신청</Title>
      <StyledForm @submit.prevent="submitLeaveForm">
        <StyledLabel>신청 종류
          <StyledSelect v-model="leaveType" required>
            <option value="">선택하세요</option>
            <option value="연차">연차</option>
            <option value="반차">반차</option>
            <option value="조퇴">조퇴</option>
            <option value="외출">외출</option>
            <option value="출산휴가">출산휴가</option>
            <option value="배우자출산휴가">배우자출산휴가</option>
            <option value="육아휴직">육아휴직</option>
            <option value="가족돌봄휴가">가족돌봄휴가</option>
            <option value="병가">병가</option>
          </StyledSelect>
          <AnnualInfo v-if="leaveType === '연차' && leaveStatus">
            신청 가능한 연차: {{ leaveStatus.remaining_days }}일
            <br>
            <small>총 발생: {{ leaveStatus.total_earned_days }}일 | 사용: {{ leaveStatus.used_days }}일 | 근무연수: {{ leaveStatus.years_of_service }}년</small>
          </AnnualInfo>
        </StyledLabel>
        <StyledLabel>신청 날짜
          <StyledInput type="date" v-model="leaveDate" required />
        </StyledLabel>
        <StyledLabel>시간(조퇴/외출 시)
          <StyledInput type="time" v-model="leaveTime" />
        </StyledLabel>
        <StyledLabel>사유
          <StyledTextarea v-model="leaveReason" placeholder="사유를 입력하세요" required />
        </StyledLabel>
        <UploadRow>
          <UploadItem>
            <StyledLabel>신청서 업로드</StyledLabel>
            <StyledInput type="file" ref="applicationFile" />
          </UploadItem>
          <UploadItem>
            <StyledLabel>제출서류 업로드</StyledLabel>
            <StyledInput type="file" ref="submitFile" />
          </UploadItem>
        </UploadRow>
        <SubmitButton type="submit">신청하기</SubmitButton>
      </StyledForm>
      <History>
        <HistoryTitle>최근 신청 내역</HistoryTitle>
        <StyledTable>
          <thead>
            <Tr>
              <Th>신청일</Th>
              <Th>종류</Th>
              <Th>날짜</Th>
              <Th>시간</Th>
              <Th>사유</Th>
              <Th>상태</Th>
            </Tr>
          </thead>
          <tbody>
            <Tr v-for="(item, index) in leaveHistory" :key="index">
              <Td>{{ item.applyDate }}</Td>
              <Td>{{ item.type }}</Td>
              <Td>{{ item.date }}</Td>
              <Td>{{ item.time || '-' }}</Td>
              <Td>{{ item.reason }}</Td>
              <Td>{{ item.status }}</Td>
            </Tr>
          </tbody>
        </StyledTable>
      </History>
    </MainContent>
  </Container>
</template>

<script>
import SideBar from '@/components/SideBar.vue';
import {
  Container, MainContent, Title, StyledForm, StyledLabel, StyledSelect, StyledInput, StyledTextarea, SubmitButton, UploadRow, UploadItem, AnnualInfo, History, HistoryTitle, StyledTable, Th, Td, Tr
} from './LeaveApply.styled.js';

export default {
  name: 'LeaveApplication',
  components: {
    SideBar,
    Container, MainContent, Title, StyledForm, StyledLabel, StyledSelect, StyledInput, StyledTextarea, SubmitButton, UploadRow, UploadItem, AnnualInfo, History, HistoryTitle, StyledTable, Th, Td, Tr
  },
  data() {
    return {
      leaveType: '',
      leaveDate: '',
      leaveTime: '',
      leaveReason: '',
      leaveStatus: null,
      leaveHistory: []
    };
  },
  methods: {
    submitLeaveForm() {
      alert('신청이 완료되었습니다!');
      this.leaveType = '';
      this.leaveDate = '';
      this.leaveTime = '';
      this.leaveReason = '';
      if (this.$refs.applicationFile) this.$refs.applicationFile.value = '';
      if (this.$refs.submitFile) this.$refs.submitFile.value = '';
    }
  }
};
</script>

<!-- 스타일은 vue-styled-components로 분리됨 --> 
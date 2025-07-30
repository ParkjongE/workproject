<template>
  <Container>
    <SideBar />
    <MainContent>
      <AttendanceTitle>내 근태내역</AttendanceTitle>
      <FilterArea>
        <label>기간
          <input type="date" v-model="startDate"> ~ <input type="date" v-model="endDate">
        </label>
        <label>유형
          <select v-model="filterType">
            <option value="">전체</option>
            <option value="정상">정상 출근/퇴근</option>
            <option value="연차">연차</option>
            <option value="반차">반차</option>
            <option value="조퇴">조퇴</option>
            <option value="외출">외출</option>
            <option value="지각">지각</option>
            <option value="결근">결근</option>
          </select>
        </label>
      </FilterArea>
      <SummaryBox>
        <div>총 근무일수: <HighlightGreen>{{ totalWorkDays }}일</HighlightGreen></div>
        <div>연차 사용: <HighlightRed>{{ usedAnnualLeave }}일</HighlightRed></div>
        <div>초과근무: <HighlightBlue>{{ overtimeHours }}</HighlightBlue></div>
      </SummaryBox>
      <AnnualLeaveInfo>
        남은 연차: {{ remainingAnnualLeave }}일 <PeriodInfo>(2025.01.01 ~ 2025.12.31)</PeriodInfo>
      </AnnualLeaveInfo>
      <StyledTable>
        <thead>
          <Tr>
            <Th>날짜</Th>
            <Th>유형</Th>
            <Th>출근</Th>
            <Th>퇴근</Th>
            <Th>근무시간</Th>
            <Th>사유</Th>
            <Th>승인상태</Th>
          </Tr>
        </thead>
        <tbody>
          <Tr v-for="(record, index) in filteredRecords" :key="index">
            <Td>{{ record.date }}</Td>
            <Td>{{ record.type }}</Td>
            <Td>{{ record.checkIn || '-' }}</Td>
            <Td>{{ record.checkOut || '-' }}</Td>
            <Td>{{ record.workHours || '-' }}</Td>
            <Td>{{ record.reason || '-' }}</Td>
            <Td>{{ record.approvalStatus || '-' }}</Td>
          </Tr>
        </tbody>
      </StyledTable>
    </MainContent>
  </Container>
</template>

<script>
import SideBar from '@/components/SideBar.vue';
import {
  Container, MainContent, AttendanceTitle, FilterArea, SummaryBox, HighlightGreen, HighlightRed, HighlightBlue, AnnualLeaveInfo, PeriodInfo, StyledTable, Th, Td, Tr
} from './MyAttendance.styled.js';

export default {
  name: 'AttendanceRecord',
  components: {
    SideBar,
    Container, MainContent, AttendanceTitle, FilterArea, SummaryBox, HighlightGreen, HighlightRed, HighlightBlue, AnnualLeaveInfo, PeriodInfo, StyledTable, Th, Td, Tr
  },
  data() {
    return {
      startDate: '2024-07-01',
      endDate: '2024-07-31',
      filterType: '',
      totalWorkDays: 20,
      usedAnnualLeave: 2,
      overtimeHours: '5시간 30분',
      remainingAnnualLeave: 8,
      attendanceRecords: [
        { date: '2024-07-01', type: '정상', checkIn: '09:00', checkOut: '18:00', workHours: '8시간', reason: '', approvalStatus: '' },
        { date: '2024-07-02', type: '연차', checkIn: '', checkOut: '', workHours: '', reason: '휴가', approvalStatus: '승인' },
        { date: '2024-07-03', type: '조퇴', checkIn: '09:00', checkOut: '15:00', workHours: '6시간', reason: '병원진료', approvalStatus: '승인' },
        { date: '2024-07-04', type: '지각', checkIn: '10:30', checkOut: '18:00', workHours: '7.5시간', reason: '지각', approvalStatus: '승인' },
        { date: '2024-07-05', type: '외출', checkIn: '09:00', checkOut: '18:00', workHours: '8시간', reason: '외근', approvalStatus: '승인' },
        { date: '2024-07-06', type: '정상', checkIn: '09:00', checkOut: '18:00', workHours: '8시간', reason: '', approvalStatus: '' }
      ]
    };
  },
  computed: {
    filteredRecords() {
      return this.attendanceRecords.filter(record => {
        const recordDate = new Date(record.date);
        const start = new Date(this.startDate);
        const end = new Date(this.endDate);
        const dateInRange = recordDate >= start && recordDate <= end;
        const typeMatches = this.filterType === '' || record.type === this.filterType;
        return dateInRange && typeMatches;
      });
    }
  }
};
</script>

<!-- 스타일은 vue-styled-components로 분리됨 --> 
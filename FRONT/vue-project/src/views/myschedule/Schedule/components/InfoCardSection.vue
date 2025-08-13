<template>
  <div>
    <!-- 연차 정보 -->
    <InfoCard>
      <SectionHeaderWithPrivacy>
        <SectionTitle>연차 정보</SectionTitle>
        <PrivacyToggle v-if="editingPrivacy">
          <ToggleLabel>
            <input type="checkbox" :checked="privacySettings.leaveInfo" @change="$emit('privacy-changed')">
            <ToggleText>{{ privacySettings.leaveInfo ? '공개' : '비공개' }}</ToggleText>
          </ToggleLabel>
        </PrivacyToggle>
        <PrivacyStatus v-else :public="privacySettings.leaveInfo">
          {{ privacySettings.leaveInfo ? '공개' : '비공개' }}
        </PrivacyStatus>
      </SectionHeaderWithPrivacy>
      <div v-if="leaveLoading" style="text-align: center; padding: 20px;">
        연차 정보를 불러오는 중...
      </div>
      <div v-else-if="leaveError" style="text-align: center; padding: 20px; color: red;">
        {{ leaveError }}
        <button @click="loadLeaveInfo" style="margin-left: 10px; padding: 5px 10px;">다시 시도</button>
      </div>
      <template v-else>
        <LeaveItem>
          <LeaveLabel>부여된 연차:</LeaveLabel>
          <LeaveValue>{{ leaveInfo ? leaveInfo.total_earned_days : totalLeave }}일</LeaveValue>
        </LeaveItem>
        <LeaveItem>
          <LeaveLabel>사용한 연차:</LeaveLabel>
          <LeaveValue>{{ leaveInfo ? leaveInfo.used_days : usedLeave }}일</LeaveValue>
        </LeaveItem>
        <LeaveItem>
          <LeaveLabel>남은 연차:</LeaveLabel>
          <LeaveValue>{{ leaveInfo ? leaveInfo.remaining_days : remainingLeave }}일</LeaveValue>
        </LeaveItem>
      </template>
    </InfoCard>
    
    <!-- 출퇴근 정보 -->
    <InfoCard>
      <SectionHeaderWithPrivacy>
        <h4>{{ isToday ? '오늘의 출퇴근' : `${selectedDate} 출퇴근` }}</h4>
        <PrivacyToggle v-if="editingPrivacy">
          <ToggleLabel>
            <input type="checkbox" :checked="privacySettings.attendance" @change="$emit('privacy-changed')">
            <ToggleText>{{ privacySettings.attendance ? '공개' : '비공개' }}</ToggleText>
          </ToggleLabel>
        </PrivacyToggle>
        <PrivacyStatus v-else :public="privacySettings.attendance">
          {{ privacySettings.attendance ? '공개' : '비공개' }}
        </PrivacyStatus>
      </SectionHeaderWithPrivacy>
      <CurrentTime>{{ currentTime }}</CurrentTime>
      
      <!-- 출퇴근 기록 표시 -->
      <AttendanceStatus>
        <StatusItem>
          <StatusLabel>출근 시간:</StatusLabel>
          <StatusValue :class="{ 'recorded': todayAttendance && todayAttendance.checkInTime }">
            {{ todayAttendance && todayAttendance.checkInTime || '-' }}
          </StatusValue>
        </StatusItem>
        <StatusItem>
          <StatusLabel>퇴근 시간:</StatusLabel>
          <StatusValue :class="{ 'recorded': todayAttendance && todayAttendance.checkOutTime }">
            {{ todayAttendance && todayAttendance.checkOutTime || '-' }}
          </StatusValue>
        </StatusItem>
        <StatusItem>
          <StatusLabel>근무 시간:</StatusLabel>
          <StatusValue :class="{ 'recorded': workingHours !== '-' }">
            {{ workingHours }}
          </StatusValue>
        </StatusItem>
        <StatusItem v-if="todayAttendance && todayAttendance.status">
          <StatusLabel>상태:</StatusLabel>
          <StatusValue :class="{ 'recorded': true }">
            {{ todayAttendance.status }}
          </StatusValue>
        </StatusItem>
      </AttendanceStatus>
      
      <!-- 출퇴근 버튼 -->
      <AttendanceButtons>
        <CheckInBtn 
          @click="handleCheckIn"
          :disabled="!isToday || (todayAttendance && todayAttendance.checkInTime)"
          :class="{ 'disabled': !isToday || (todayAttendance && todayAttendance.checkInTime) }"
        >
          {{ todayAttendance && todayAttendance.checkInTime ? '출근 완료' : '출근' }}
        </CheckInBtn>
        <CheckOutBtn 
          @click="handleCheckOut"
          :disabled="!isToday || !todayAttendance || !todayAttendance.checkInTime || (todayAttendance && todayAttendance.checkOutTime)"
          :class="{ 'disabled': !isToday || !todayAttendance || !todayAttendance.checkInTime || (todayAttendance && todayAttendance.checkOutTime) }"
        >
          {{ todayAttendance && todayAttendance.checkOutTime ? '퇴근 완료' : '퇴근' }}
        </CheckOutBtn>
      </AttendanceButtons>
    </InfoCard>
    
    <!-- 급여명세서 -->
    <InfoCard>
      <SectionHeaderWithPrivacy>
        <h4>급여명세서</h4>
        <PrivacyToggle v-if="editingPrivacy">
          <ToggleLabel>
            <input type="checkbox" :checked="privacySettings.salary" @change="$emit('privacy-changed')">
            <ToggleText>{{ privacySettings.salary ? '공개' : '비공개' }}</ToggleText>
          </ToggleLabel>
        </PrivacyToggle>
        <PrivacyStatus v-else :public="privacySettings.salary">
          {{ privacySettings.salary ? '공개' : '비공개' }}
        </PrivacyStatus>
      </SectionHeaderWithPrivacy>
      <UploadBtn @click="showUploadForm = true">
        <UploadIcon>+</UploadIcon> 파일 첨부
      </UploadBtn>
      <SalaryYears>
        <SalaryYear 
          v-for="year in Object.keys(salaryStatements).sort((a, b) => parseInt(b) - parseInt(a))" 
          :key="year"
          :class="{ 'active': selectedSalaryYear === year }"
        >
          <YearHeader @click="toggleSalaryYear(year)">
            <span>{{ year }}년</span>
            <YearActions>
              <DownloadAllBtn 
                title="모두 다운로드"
                @click.stop="downloadAllStatements(year)"
              >
                <DownloadAllIcon>↓↓</DownloadAllIcon>
              </DownloadAllBtn>
              <DropdownArrow :open="selectedSalaryYear === year"></DropdownArrow>
            </YearActions>
          </YearHeader>
          <SalaryMonths v-if="selectedSalaryYear === year">
            <SalaryMonth 
              v-for="month in getSalaryMonths(year)" 
              :key="`${year}-${month}`"
              @click.stop="downloadSalaryStatement(year, month)"
            >
              <span>{{ month }}월</span>
              <DownloadIcon>↓</DownloadIcon>
            </SalaryMonth>
            <NoStatements v-if="getSalaryMonths(year).length === 0">
              명세서가 없습니다
            </NoStatements>
          </SalaryMonths>
        </SalaryYear>
        <NoStatements v-if="Object.keys(salaryStatements).length === 0">
          등록된 급여명세서가 없습니다
        </NoStatements>
      </SalaryYears>
    </InfoCard>
  </div>
</template>

<script>
import * as styled from '../Schedule.styled.js';
import { watch, computed } from 'vue';

// Stores
import { useAttendanceStore } from '@/stores/attendance';
import api from '@/plugins/axios';

export default {
  name: 'InfoCardSection',
  components: {
    ...styled
  },
  props: {
    editingPrivacy: {
      type: Boolean,
      default: false
    },
    privacySettings: {
      type: Object,
      default: () => ({
        leaveInfo: true,
        attendance: true,
        salary: false
      })
    },
    selectedDate: {
      type: String,
      default: () => new Date().toISOString().split('T')[0]
    }
  },
  emits: ['privacy-changed', 'attendance-updated'],
  setup(props) {

    const attendanceStore = useAttendanceStore();

    // 선택된 날짜가 변경될 때마다 해당 날짜의 출퇴근 기록 로드
    watch(() => props.selectedDate, async (newDate) => {
      if (newDate) {
        await attendanceStore.loadAttendanceByDate(newDate);
      }
    }, { immediate: true });

    // computed 속성들
    const isToday = computed(() => {
      const today = new Date().toISOString().split('T')[0];
      return props.selectedDate === today;
    });

    const isPastDate = computed(() => {
      const today = new Date().toISOString().split('T')[0];
      return props.selectedDate < today;
    });

    return {
      // Attendance Store
      attendanceStore,
      currentTime: computed(() => attendanceStore.currentTime.value),
      todayAttendance: computed(() => attendanceStore.todayAttendance.value),
      workingHours: computed(() => attendanceStore.workingHours.value),
      checkIn: attendanceStore.checkIn,
      checkOut: attendanceStore.checkOut,
      isToday,
      isPastDate,
      loadAttendanceByDate: attendanceStore.loadAttendanceByDate
    };
  },
  mounted() {
    console.log('[InfoCardSection] 컴포넌트 마운트됨, 선택된 날짜:', this.selectedDate);
    console.log('[InfoCardSection] loadLeaveInfo 함수 존재 여부:', typeof this.loadLeaveInfo);
    
    // 컴포넌트가 마운트될 때 선택된 날짜의 출퇴근 기록 로드
    if (this.selectedDate) {
      console.log('[InfoCardSection] mounted에서 출퇴근 기록 로드 시작');
      this.attendanceStore.loadAttendanceByDate(this.selectedDate);
    } else {
      console.log('[InfoCardSection] selectedDate가 없음');
    }
    
    // 연차 정보 로드
    console.log('[InfoCardSection] 연차 정보 로드 시작');
    this.loadLeaveInfo();
    this.loadSalaryStatements();
    console.log('[InfoCardSection] 연차 정보 로드 호출 완료');
  },
  watch: {
    selectedDate: {
      handler(newDate) {
        console.log('[InfoCardSection] 선택된 날짜 변경:', newDate);
        if (newDate) {
          console.log('[InfoCardSection] watch에서 출퇴근 기록 로드 시작');
          this.attendanceStore.loadAttendanceByDate(newDate);
        }
      },
      immediate: true
    }
  },
  data() {
    return {
      // 급여명세서 관련
      selectedSalaryYear: null,
      showUploadForm: false,
      uploadForm: {
        year: new Date().getFullYear(),
        month: new Date().getMonth() + 1,
        file: null,
        fileName: ''
      },
      salaryStatements: {},
      salaryUpdateInterval: null, // 자동 업데이트 인터벌
      
      // 연차 정보
      totalLeave: 0,
      usedLeave: 0,
      leaveInfo: null,
      leaveLoading: false,
      leaveError: null
    };
  },
  computed: {
    remainingLeave() {
      if (this.leaveInfo) {
        return this.leaveInfo.remaining_days;
      }
      return this.totalLeave - this.usedLeave;
    }
  },
  methods: {
    // 연차 정보 관련
    async loadLeaveInfo() {
      try {
        this.leaveLoading = true;
        this.leaveError = null;
        console.log('[loadLeaveInfo] 연차 정보 로드 시작');
        
        // 로그인 상태 확인
        const userStr = localStorage.getItem('user');
        console.log('[loadLeaveInfo] localStorage user:', userStr);
        
        if (!userStr) {
          console.error('[loadLeaveInfo] 로그인 정보가 없습니다');
          this.leaveError = '로그인 정보가 없습니다';
          return;
        }
        
        // 현재 로그인한 사용자 정보 가져오기
        const userResponse = await api.get('/api/user/profile');
        
        console.log('[loadLeaveInfo] 사용자 정보 응답:', userResponse.status);
        
        const userResult = userResponse.data;
        console.log('[loadLeaveInfo] 사용자 정보:', userResult);
        
        if (!userResult.success || !userResult.user) {
          console.error('사용자 정보 없음');
          this.leaveError = '사용자 정보를 가져올 수 없습니다';
          return;
        }
        
        const userId = userResult.user.id;
        
        console.log('[loadLeaveInfo] 연차 정보 요청:', `user_id=${userId}`);
        
        // API 호출하여 연차 정보 가져오기
        const response = await api.get(`/api/leave/status?user_id=${userId}`);
        
        console.log('[loadLeaveInfo] 연차 정보 응답:', response.status);
        
        const result = response.data;
        console.log('[loadLeaveInfo] 연차 정보 결과:', result);
        
        if (result.success && result.data) {
          this.leaveInfo = result.data;
          this.totalLeave = result.data.total_earned_days;
          this.usedLeave = result.data.used_days;
          
          console.log('[loadLeaveInfo] 연차 정보 설정 완료:', {
            totalLeave: this.totalLeave,
            usedLeave: this.usedLeave,
            remainingLeave: this.remainingLeave
          });
        } else {
          this.leaveError = '연차 정보를 가져올 수 없습니다';
        }
      } catch (error) {
        console.error('연차 정보 로드 오류:', error);
        this.leaveError = '연차 정보 로드 중 오류가 발생했습니다';
      } finally {
        this.leaveLoading = false;
      }
    },
    
    // 출근/퇴근 관련
    async handleCheckIn() {
      console.log('[InfoCardSection] 출근 처리 - 날짜:', this.selectedDate);
      
      // 오늘 날짜인 경우에만 실제 출근 처리
      if (this.isToday) {
        await this.attendanceStore.checkIn();
        
        // 출퇴근 기록 업데이트 후 부모 컴포넌트에 알림
        this.$emit('attendance-updated');
      } else {
        alert('오늘 날짜에만 출근 기록을 할 수 있습니다.');
      }
    },
    
    async handleCheckOut() {
      console.log('[InfoCardSection] 퇴근 처리 - 날짜:', this.selectedDate);
      
      // 오늘 날짜인 경우에만 실제 퇴근 처리
      if (this.isToday) {
        await this.attendanceStore.checkOut();
        
        // 출퇴근 기록 업데이트 후 부모 컴포넌트에 알림
        this.$emit('attendance-updated');
      } else {
        alert('오늘 날짜에만 퇴근 기록을 할 수 있습니다.');
      }
    },
    
    // 급여명세서 관련
    async loadSalaryStatements() {
      try {
        console.log('[loadSalaryStatements] 급여명세서 목록 로드 시작');
        
        // 로그인 상태 확인
        const userStr = localStorage.getItem('user');
        if (!userStr) {
          console.error('[loadSalaryStatements] 로그인 정보가 없습니다');
          return;
        }
        
        const user = JSON.parse(userStr);
        const userId = user.id;
        
        console.log('[loadSalaryStatements] 사용자 ID:', userId);
        
        // API 호출하여 급여명세서 목록 가져오기
        const response = await api.get(`/api/salary/list?user_id=${userId}`);
        
        console.log('[loadSalaryStatements] 급여명세서 목록 응답:', response.status);
        
        const result = response.data;
        console.log('[loadSalaryStatements] 급여명세서 목록 결과:', result);
        
        if (result.success && result.data) {
          // 서버 데이터를 프론트엔드 형식으로 변환
          const statements = {};
          
          result.data.forEach((item) => {
            const year = item.year.toString();
            const month = item.month.toString();
            
            if (!statements[year]) {
              statements[year] = {};
            }
            
            statements[year][month] = {
              fileName: item.fileName,
              fileSize: item.fileSize,
              year: item.year,
              month: item.month
            };
          });
          
          console.log('[loadSalaryStatements] 변환된 statements:', statements);
          this.salaryStatements = statements;
          
          // 첫 번째 년도를 자동으로 선택
          const years = Object.keys(statements).sort((a, b) => parseInt(b) - parseInt(a));
          if (years.length > 0 && !this.selectedSalaryYear) {
            this.selectedSalaryYear = years[0];
          }
          
          console.log('[loadSalaryStatements] 급여명세서 로드 완료:', {
            statements: this.salaryStatements,
            selectedYear: this.selectedSalaryYear,
            years: years
          });
        } else {
          console.log('[loadSalaryStatements] 성공하지 않았거나 데이터가 없음');
          this.salaryStatements = {};
        }
      } catch (error) {
        console.error('[loadSalaryStatements] 급여명세서 목록 로드 오류:', error);
        this.salaryStatements = {};
      }
    },
    
    async toggleSalaryYear(year) {
      console.log('[toggleSalaryYear] 토글 요청:', year, '현재 선택된 년도:', this.selectedSalaryYear);
      if (this.selectedSalaryYear === year) {
        this.selectedSalaryYear = null;
      } else {
        this.selectedSalaryYear = year;
        // 년도를 선택할 때 목록 자동 업데이트
        console.log('[toggleSalaryYear] 년도 선택 시 목록 자동 업데이트 시작');
        await this.loadSalaryStatements();
        console.log('[toggleSalaryYear] 목록 자동 업데이트 완료');
      }
      console.log('[toggleSalaryYear] 토글 후 선택된 년도:', this.selectedSalaryYear);
    },
    
    getSalaryMonths(year) {
      if (!this.salaryStatements[year]) return [];
      const months = Object.keys(this.salaryStatements[year]);
      // 월을 숫자 순서대로 정렬 (1, 2, 3, ..., 12)
      return months.sort((a, b) => parseInt(a) - parseInt(b));
    },
    
    handleFileChange(event) {
      const file = event.target.files[0];
      if (file) {
        this.uploadForm.file = file;
        this.uploadForm.fileName = file.name;
      } else {
        this.uploadForm.file = null;
        this.uploadForm.fileName = '';
      }
      this.$emit('privacy-changed');
    },
    
    uploadSalaryStatement() {
      if (!this.uploadForm.file) {
        alert('파일을 선택해주세요.');
        return;
      }
      
      console.log('급여명세서 업로드:', this.uploadForm);
      
      const year = this.uploadForm.year.toString();
      const month = this.uploadForm.month.toString();
      
      if (!this.salaryStatements[year]) {
        this.$set(this.salaryStatements, year, {});
      }
      
      this.$set(this.salaryStatements[year], month, {
        fileName: this.uploadForm.fileName,
        url: URL.createObjectURL(this.uploadForm.file)
      });
      
      this.uploadForm.file = null;
      this.uploadForm.fileName = '';
      this.showUploadForm = false;
      this.selectedSalaryYear = year;
      
      alert('급여명세서가 성공적으로 업로드되었습니다.');
      this.$emit('privacy-changed');
    },
    
    async downloadSalaryStatement(year, month) {
      try {
        console.log('[downloadSalaryStatement] 급여명세서 다운로드 시작:', year, month);
        
        // 로그인 상태 확인
        const userStr = localStorage.getItem('user');
        if (!userStr) {
          console.error('[downloadSalaryStatement] 로그인 정보가 없습니다');
          alert('로그인 정보가 없습니다.');
          return;
        }
        
        const user = JSON.parse(userStr);
        const userId = user.id;
        
        console.log('[downloadSalaryStatement] 급여명세서 다운로드 요청:', `user_id=${userId}, year=${year}, month=${month}`);
        
        // API 호출하여 급여명세서 다운로드
        const response = await api.get(`/api/salary/statement?user_id=${userId}&year=${year}&month=${month}`, {
          responseType: 'blob'
        });
        
        console.log('[downloadSalaryStatement] 급여명세서 다운로드 응답:', response.status);
        
        // 파일 다운로드
        const blob = new Blob([response.data], { type: 'application/pdf' });
        const url = window.URL.createObjectURL(blob);
        const link = document.createElement('a');
        link.href = url;
        link.download = `급여명세서_${year}년_${month}월.pdf`;
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
        window.URL.revokeObjectURL(url);
        
        console.log('[downloadSalaryStatement] 급여명세서 다운로드 완료');
        
        // 다운로드 완료 후 목록 자동 업데이트
        console.log('[downloadSalaryStatement] 목록 자동 업데이트 시작');
        await this.loadSalaryStatements();
        console.log('[downloadSalaryStatement] 목록 자동 업데이트 완료');
        
      } catch (error) {
        console.error('[downloadSalaryStatement] 급여명세서 다운로드 오류:', error);
        if (error.response && error.response.status === 404) {
          alert('해당 급여명세서를 찾을 수 없습니다.');
        } else {
          alert('급여명세서 다운로드 중 오류가 발생했습니다.');
        }
      }
    },
    
    async downloadAllStatements(year) {
      if (!this.salaryStatements[year] || Object.keys(this.salaryStatements[year]).length === 0) {
        alert('해당 년도의 급여명세서가 없습니다.');
        return;
      }
      
      const months = this.getSalaryMonths(year);
      alert(`${year}년 모든 급여명세서(${months.length}개)를 다운로드합니다.`);
      
      // 로그인 상태 확인
      const userStr = localStorage.getItem('user');
      if (!userStr) {
        alert('로그인 정보가 없습니다.');
        return;
      }
      
      const user = JSON.parse(userStr);
      const userId = user.id;
      
      // 각 월별로 순차적으로 다운로드
      for (let i = 0; i < months.length; i++) {
        const month = months[i];
        try {
          console.log(`[downloadAllStatements] ${year}년 ${month}월 다운로드 시작`);
          
          const response = await api.get(`/api/salary/statement?user_id=${userId}&year=${year}&month=${month}`, {
            responseType: 'blob'
          });
          
          const blob = new Blob([response.data], { type: 'application/pdf' });
          const url = window.URL.createObjectURL(blob);
          const link = document.createElement('a');
          link.href = url;
          link.download = `급여명세서_${year}년_${month}월.pdf`;
          document.body.appendChild(link);
          link.click();
          document.body.removeChild(link);
          window.URL.revokeObjectURL(url);
          
          console.log(`[downloadAllStatements] ${year}년 ${month}월 다운로드 완료`);
          
          // 다음 다운로드 전에 잠시 대기
          if (i < months.length - 1) {
            await new Promise(resolve => setTimeout(resolve, 500));
          }
        } catch (error) {
          console.error(`[downloadAllStatements] ${year}년 ${month}월 다운로드 오류:`, error);
          alert(`${year}년 ${month}월 급여명세서 다운로드 중 오류가 발생했습니다.`);
        }
      }
      
      console.log(`[downloadAllStatements] ${year}년 모든 급여명세서 다운로드 완료`);
      
      // 모든 다운로드 완료 후 목록 자동 업데이트
      console.log('[downloadAllStatements] 목록 자동 업데이트 시작');
      await this.loadSalaryStatements();
      console.log('[downloadAllStatements] 목록 자동 업데이트 완료');
    }
  }
};
</script> 
<template>
  <MyScheduleRoot>
    <SideBar />
    <ContentArea>
      <HeaderContainer>
        <Title>나의 일정 수정</Title>
        <EditPrivacyBtn v-if="!editingPrivacy" @click="startEditingPrivacy">
          <i class="edit-icon">✏️</i> 공개 설정 편집
        </EditPrivacyBtn>
        <DonePrivacyBtn v-else @click="finishEditingPrivacy">
          <i class="done-icon">✓</i> 완료
        </DonePrivacyBtn>
      </HeaderContainer>
      <ScheduleManagement>
        <!-- 시계 섹션 -->
        <ClockSection 
          :selectedDate="selectedDate"
          :mySchedules="mySchedules"
          @schedule-selected="handleScheduleSelected"
          @schedule-edited="handleScheduleEdited"
          @schedule-deleted="handleScheduleDeleted"
          @schedule-added="handleScheduleAdded"
        />
        
            <!-- 정보 카드 섹션 -->
    <InfoCardSection 
      :editingPrivacy="editingPrivacy"
      :privacySettings="privacySettings"
      :selectedDate="selectedDate"
      @privacy-changed="handlePrivacyChanged"
      @attendance-updated="handleAttendanceUpdated"
    />
        
        <!-- 달력 섹션 -->
        <CalendarSection 
          :selectedDate="selectedDate"
          :mySchedules="mySchedules"
          @day-selected="onDaySelected"
          @schedule-selected="handleScheduleSelected"
        />
      </ScheduleManagement>
      
      <!-- 파일 업로드 폼 -->
      <div v-if="showUploadForm" class="upload-form">
        <div class="form-header">
          <h4>급여명세서 첨부</h4>
          <button class="close-btn" @click="showUploadForm = false">×</button>
        </div>
        <div class="form-body">
          <div class="form-group">
            <label>년도</label>
            <select v-model="uploadForm.year">
              <option v-for="year in availableYears" :key="year" :value="year">{{ year }}년</option>
            </select>
          </div>
          <div class="form-group">
            <label>월</label>
            <select v-model="uploadForm.month">
              <option v-for="month in 12" :key="month" :value="month">{{ month }}월</option>
            </select>
          </div>
          <div class="form-group">
            <label>파일</label>
            <input type="file" ref="fileInput" @change="handleFileChange">
            <div class="file-name">{{ uploadForm.fileName || '선택된 파일 없음' }}</div>
          </div>
          <div class="form-actions">
            <button class="cancel-btn" @click="showUploadForm = false">취소</button>
            <button class="submit-btn" @click="uploadSalaryStatement" :disabled="!uploadForm.file">첨부</button>
          </div>
        </div>
      </div>
    </ContentArea>
  </MyScheduleRoot>
</template>

<script>
import SideBar from '@/components/SideBar.vue';
import ClockSection from './components/ClockSection.vue';
import InfoCardSection from './components/InfoCardSection.vue';
import CalendarSection from './components/CalendarSection.vue';
// eslint-disable-next-line no-unused-vars
import { MyScheduleRoot, ContentArea, HeaderContainer, Title, EditPrivacyBtn, DonePrivacyBtn, ScheduleManagement } from './Schedule.styled.js';

// Composables
import { useSchedule } from './composables/useSchedule';
import { useCalendar } from './composables/useCalendar';
import { useAttendanceStore } from '@/stores/attendance';

export default {
  name: 'MySchedule',
  components: {
    SideBar,
    ClockSection,
    InfoCardSection,
    CalendarSection,
    MyScheduleRoot,
    ContentArea,
    HeaderContainer,
    Title,
    EditPrivacyBtn,
    DonePrivacyBtn,
    ScheduleManagement
  },
  watch: {
    $route() {
      this.fetchMySchedules();
    },
    mySchedules: {
      handler(newSchedules) {
        console.log('[Schedule] mySchedules 변경됨, 일정 수:', newSchedules.length);
      },
      deep: true
    }
  },
  setup() {
    // Composables 사용
    const {
      mySchedules,
      selectedSchedule,
      fetchMySchedules
    } = useSchedule();

    const {
      availableYears,
      initializeCalendar
    } = useCalendar();

    const attendanceStore = useAttendanceStore();

    return {
      // Schedule
      mySchedules,
      selectedSchedule,
      fetchMySchedules,
      // Calendar
      availableYears,
      initializeCalendar,
      // Attendance
      attendanceStore
    };
  },
  data() {
    const today = new Date();
    return {
      // 공개 설정
      editingPrivacy: false,
      privacySettings: {
        leaveInfo: true,
        attendance: true,
        salary: false
      },
      // 달력 날짜 상태
      selectedDate: `${today.getFullYear()}-${String(today.getMonth() + 1).padStart(2, '0')}-${String(today.getDate()).padStart(2, '0')}`,
      // 급여명세서 관련
      showUploadForm: false,
      uploadForm: {
        year: today.getFullYear(),
        month: today.getMonth() + 1,
        file: null,
        fileName: ''
      },
      // 삭제된 일정 추적
      lastDeletedSchedule: null
    };
  },
  mounted() {
    // 초기화
    this.initializeCalendar();
    this.attendanceStore.initializeAttendance();
    this.fetchMySchedules();
  },
  beforeUpdate() {
    // 컴포넌트 업데이트 전 실행
    console.log('[Schedule] 컴포넌트 업데이트 전, 현재 일정 수:', this.mySchedules.length);
    
    // 삭제된 일정이 있는 경우 처리
    if (this.lastDeletedSchedule) {
      console.log('[Schedule] 삭제된 일정 정보:', this.lastDeletedSchedule.title, 'ID:', this.lastDeletedSchedule.id);
      
      // 삭제된 일정이 현재 목록에서 제거되었는지 확인
      const deletedScheduleExists = this.mySchedules.find(schedule => schedule.id === this.lastDeletedSchedule.id);
      
      if (!deletedScheduleExists) {
        console.log('[Schedule] 일정 삭제 확인됨 - ID:', this.lastDeletedSchedule.id);
        
        // 삭제된 일정의 날짜가 현재 선택된 날짜인지 확인
        const deletedDate = this.lastDeletedSchedule.start.split('T')[0];
        const isCurrentDate = deletedDate === this.selectedDate;
        
        if (isCurrentDate) {
          console.log('[Schedule] 현재 날짜에서 일정 삭제됨 - 화면 재렌더링 시작');
          
          // 현재 날짜의 남은 일정 수 확인
          const currentDateSchedules = this.mySchedules.filter(schedule => {
            const scheduleDate = schedule.start.split('T')[0];
            return scheduleDate === this.selectedDate;
          });
          console.log('[Schedule] 현재 날짜 남은 일정 수:', currentDateSchedules.length);
        } else {
          console.log('[Schedule] 다른 날짜의 일정 삭제됨 - 달력 뷰 업데이트 시작');
        }
      }
      
      // 삭제된 일정 정보 초기화
      this.lastDeletedSchedule = null;
    }
    
    // 전체 일정 상태 확인
    if (this.mySchedules.length === 0) {
      console.log('[Schedule] 모든 일정이 삭제됨, 빈 화면 표시 시작');
    }
  },
  beforeUnmount() {
    // 정리
    this.attendanceStore.stopTimer();
  },
  methods: {
    // 공개 설정 관련
    startEditingPrivacy() {
      this.editingPrivacy = true;
    },
    
    finishEditingPrivacy() {
      this.editingPrivacy = false;
      this.savePrivacySettings();
    },
    
    savePrivacySettings() {
      console.log('공개 설정 저장:', this.privacySettings);
      localStorage.setItem('privacySettings', JSON.stringify(this.privacySettings));
    },
    
    handlePrivacyChanged() {
      this.savePrivacySettings();
    },
    
    // 일정 관련 이벤트 핸들러
    handleScheduleSelected(schedule) {
      console.log('일정 선택됨:', schedule);
    },
    
    handleScheduleEdited(schedule) {
      console.log('일정 편집됨:', schedule);
      this.fetchMySchedules(); // 일정 수정 후 데이터 다시 불러오기
    },
    
    handleScheduleDeleted(schedule) {
      console.log('일정 삭제됨:', schedule);
      // 삭제된 일정 정보 저장 (beforeUpdate 훅에서 활용)
      this.lastDeletedSchedule = schedule;
      
      // 일정 추가와 같은 방식으로 서버에서 데이터 다시 불러오기
      this.fetchMySchedules();
      console.log('[Schedule] 삭제 후 데이터 다시 불러오기 실행');
    },
    
    handleScheduleAdded() {
      console.log('일정 추가됨');
      this.fetchMySchedules(); // 일정 추가 후 데이터 다시 불러오기
    },
    
    // 출퇴근 기록 업데이트 이벤트 핸들러
    handleAttendanceUpdated() {
      console.log('출퇴근 기록 업데이트됨');
      // useAttendance composable의 데이터가 자동으로 업데이트되므로
      // 별도의 전역 이벤트가 필요하지 않음
    },
    
    // 날짜 선택 이벤트 핸들러
    onDaySelected(date) {
      this.selectedDate = date;
    },
    
    // 파일 업로드 관련
    handleFileChange(event) {
      const file = event.target.files[0];
      if (file) {
        this.uploadForm.file = file;
        this.uploadForm.fileName = file.name;
      } else {
        this.uploadForm.file = null;
        this.uploadForm.fileName = '';
      }
    },
    
    uploadSalaryStatement() {
      if (!this.uploadForm.file) {
        alert('파일을 선택해주세요.');
        return;
      }
      
      console.log('급여명세서 업로드:', this.uploadForm);
      alert('급여명세서가 성공적으로 업로드되었습니다.');
      
      this.uploadForm.file = null;
      this.uploadForm.fileName = '';
      this.showUploadForm = false;
    }
  }
};
</script>

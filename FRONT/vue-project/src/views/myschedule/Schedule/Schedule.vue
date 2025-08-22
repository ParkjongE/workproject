<template>
  <MyScheduleRoot>
    <SideBar v-if="!hideSideBar" />
    <ContentArea :class="{ 'full-width': hideSideBar }">
      <HeaderContainer>
        <Title>{{ isOtherMember ? `${memberData?.name || '구성원'}의 일정` : '나의 일정 수정' }}</Title>
        <EditPrivacyBtn v-if="!isOtherMember && !editingPrivacy" @click="startEditingPrivacy">
          <i class="edit-icon">✏️</i> 공개 설정 편집
        </EditPrivacyBtn>
        <DonePrivacyBtn v-if="!isOtherMember && editingPrivacy" @click="finishEditingPrivacy">
          <i class="done-icon">✓</i> 완료
        </DonePrivacyBtn>
      </HeaderContainer>
      <ScheduleManagement>
        <!-- 시계 섹션 -->
        <ClockSection 
          :selectedDate="selectedDate"
          :mySchedules="currentSchedules"
          :isOtherMember="isOtherMember"
          :memberData="memberData"
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
          :isOtherMember="isOtherMember"
          :memberData="memberData"
          @privacy-changed="handlePrivacyChanged"
          @attendance-updated="handleAttendanceUpdated"
        />
        
        <!-- 달력 섹션 -->
        <CalendarSection 
          :selectedDate="selectedDate"
          :mySchedules="currentSchedules"
          :isOtherMember="isOtherMember"
          :memberData="memberData"
          @day-selected="onDaySelected"
          @schedule-selected="handleScheduleSelected"
        />
      </ScheduleManagement>
      
      <!-- 파일 업로드 폼 (본인만 표시) -->
      <div v-if="!isOtherMember && showUploadForm" class="upload-form">
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
  props: {
    hideSideBar: {
      type: Boolean,
      default: false
    },
    // 다른 구성원 데이터를 위한 props
    memberData: {
      type: Object,
      default: null
    },
    isOtherMember: {
      type: Boolean,
      default: false
    }
  },
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
      // 다른 구성원의 일정을 보는 경우에는 라우트 변경 감지를 하지 않음
      if (this.isOtherMember) return;
      this.fetchMySchedules();
    },
    mySchedules: {
      handler(newSchedules) {
        // 다른 구성원의 일정을 보는 경우에는 일정 변경 감지를 하지 않음
        if (this.isOtherMember) return;
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
  computed: {
    // 다른 구성원의 일정을 사용할지 결정
    currentSchedules() {
      if (this.isOtherMember && this.memberData) {
        return this.memberData.schedules || [];
      }
      return this.mySchedules;
    },
    
    // 다른 구성원의 이름
    currentMemberName() {
      if (this.isOtherMember && this.memberData) {
        return this.memberData.name || '구성원';
      }
      return '나';
    }
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
    // 다른 구성원의 일정을 보는 경우에도 달력 초기화는 필요 (드래그 투 셀렉트를 위해)
    this.initializeCalendar();
    
    // 다른 구성원의 일정을 보는 경우에는 데이터 로딩을 하지 않음
    if (this.isOtherMember) {
      console.log('[Schedule] 다른 구성원 일정 모드 - 데이터 로딩 건너뜀, 드래그 투 셀렉트는 활성화');
      return;
    }
    
    // 초기화 (본인 일정 모드에서만)
    this.attendanceStore.initializeAttendance();
    this.fetchMySchedules();
    
    // 초기 선택된 날짜를 localStorage에 저장
    try {
      localStorage.setItem('myScheduleSelectedDate', this.selectedDate);
    } catch (error) {
      console.warn('로컬스토리지에 초기 선택된 날짜 저장 중 오류:', error);
    }
  },
  beforeUpdate() {
    // 다른 구성원의 일정을 보는 경우에는 업데이트 로직을 실행하지 않음
    if (this.isOtherMember) return;
    
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
      console.log('[Schedule] 컴포넌트 업데이트 전, 현재 일정 수:', this.mySchedules.length);
    }
  },
  beforeUnmount() {
    // 다른 구성원의 일정을 보는 경우에는 타이머 정리를 하지 않음
    if (this.isOtherMember) return;
    
    // 정리 (본인 일정 모드에서만)
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
      // 선택된 날짜를 localStorage에 저장하여 다른 페이지에서 참조할 수 있도록 함
      try {
        localStorage.setItem('myScheduleSelectedDate', date);
      } catch (error) {
        console.warn('로컬스토리지에 선택된 날짜 저장 중 오류:', error);
      }
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

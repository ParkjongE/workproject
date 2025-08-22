<template>
  <ScheduleSharingRoot>
    <SideBar />
    <ContentArea>
      <h1>개인 일정 공유</h1>
      
      <DashboardContainer>
        <GraphView 
          ref="graphView" 
          class="graph-view"
          @dragover="handleDragOver"
          @drop="handleDrop"
        >
          <!-- 부서별 경계 -->
          <DepartmentBoundary 
            v-for="department in departments" 
            :key="department.name"
            :style="getDepartmentBoundaryStyle(department)"
          >
            <DepartmentLabel>{{ department.name }}</DepartmentLabel>
          </DepartmentBoundary>
          
          <!-- 본인 노드 (고정 위치) -->
          <MyNode
            :style="getMyNodePosition()"
            @click.stop="selectMember(myProfile.id)"
          >
            <MemberInfo>
              <MemberPosition>{{ myProfile.position }} | {{ myProfile.department }} {{ myProfile.team }}</MemberPosition>
              <MemberName>{{ myProfile.name }}</MemberName>
            </MemberInfo>
            <ScheduleViews>
              <!-- 미니 시계 뷰 -->
              <MiniClockView @click.stop="showClockView(myProfile.id)">
                <MiniClock>
                  <ClockFace>
                    <!-- 시간 숫자 -->
                    <HourNumber v-for="hour in 12" :key="`my-hour-${hour}`" :style="getHourNumberStyle(hour)">
                      {{ hour }}
                    </HourNumber>
                    <!-- 일정 세그먼트 -->
                    <!-- 디버깅: 오늘 일정 개수 {{ getTodaySchedules(myProfile).length }} -->
                    <ScheduleSegment v-for="(schedule, idx) in getTodaySchedules(myProfile)" :key="`my-schedule-${idx}`" 
                      :style="getScheduleSegmentStyle(schedule)"
                      :title="`${schedule.title} (${schedule.color})`">
                    </ScheduleSegment>
                  </ClockFace>
                </MiniClock>
              </MiniClockView>
              
              <!-- 미니 달력 뷰 -->
              <MiniCalendarView @click.stop="showCalendarView(myProfile.id)">
                <MiniCalendar>
                  <MiniCalendarGrid>
                    <!-- 요일 헤더 -->
                    <MiniDayHeader v-for="i in 7" :key="`my-day-${i}`">{{ weekdays[i-1].substring(0, 1) }}</MiniDayHeader>
                    <!-- 날짜 셀 -->
                    <MiniDayCell v-for="i in 35" :key="`my-cell-${i}`">
                      <MiniDayNumber>{{ getMiniCalendarDayNumber(i) }}</MiniDayNumber>
                      <MiniDayEvent v-if="getMiniCalendarDaySchedule(myProfile, i)" 
                        :style="{ backgroundColor: getMiniCalendarDaySchedule(myProfile, i).color }">
                      </MiniDayEvent>
                      <!-- 미니 달력 프로젝트 선 추가 -->
                      <MiniProjectLine v-if="isProjectActive(getMiniCalendarDate(i))" 
                        :style="{ backgroundColor: getProjectForDate(getMiniCalendarDate(i))?.color || '#B4D3F5' }">
                      </MiniProjectLine>
                    </MiniDayCell>
                  </MiniCalendarGrid>
                </MiniCalendar>
              </MiniCalendarView>
            </ScheduleViews>
          </MyNode>
          
          <!-- 구성원 노드 -->
          <MemberNode
            v-for="(member, index) in teamMembers" 
            :key="member.id"
            :data-member-id="member.id"
            :style="getNodePosition(index)"
            :class="{ 'selected': selectedMember === member.id }"
            @click.stop="selectMember(member.id)"
            draggable="true"
            @dragstart="handleDragStart($event, member.id)"
            @dragend="handleDragEnd"
          >
            <MemberInfo>
              <MemberPosition>{{ member.role }} | {{ member.department }}</MemberPosition>
              <MemberName>{{ member.name }}</MemberName>
            </MemberInfo>
            <ScheduleViews>
              <!-- 미니 시계 뷰 -->
              <MiniClockView>
                <MiniClock>
                  <ClockFace>
                    <!-- 시간 숫자 -->
                    <HourNumber v-for="hour in 12" :key="hour" :style="getHourNumberStyle(hour)">
                      {{ hour }}
                    </HourNumber>
                    <!-- 일정 세그먼트 -->
                    <ScheduleSegment v-for="(schedule, idx) in getTodaySchedules(member)" :key="idx" 
                      :style="getScheduleSegmentStyle(schedule)">
                    </ScheduleSegment>
                  </ClockFace>
                </MiniClock>
              </MiniClockView>
              
              <!-- 미니 달력 뷰 -->
              <MiniCalendarView>
                <MiniCalendar>
                  <MiniCalendarGrid>
                    <!-- 요일 헤더 -->
                    <MiniDayHeader v-for="i in 7" :key="`day-${i}`">{{ weekdays[i-1].substring(0, 1) }}</MiniDayHeader>
                    <!-- 날짜 셀 -->
                    <MiniDayCell v-for="i in 35" :key="`cell-${i}`">
                      <MiniDayNumber>{{ getMiniCalendarDayNumber(i) }}</MiniDayNumber>
                      <MiniDayEvent v-if="getMiniCalendarDaySchedule(member, i)" 
                        :style="{ backgroundColor: getMiniCalendarDaySchedule(member, i).color }">
                      </MiniDayEvent>
                      <!-- 미니 달력 프로젝트 선 추가 -->
                      <MiniProjectLine v-if="isProjectActive(getMiniCalendarDate(i))" 
                        :style="{ backgroundColor: getProjectForDate(getMiniCalendarDate(i))?.color || '#B4D3F5' }">
                      </MiniProjectLine>
                    </MiniDayCell>
                  </MiniCalendarGrid>
                </MiniCalendar>
              </MiniCalendarView>
            </ScheduleViews>
            
            <!-- 디버깅용 정보 표시 -->
            <div style="position: absolute; top: -20px; left: 0; font-size: 10px; color: red; background: white; padding: 2px;">
              ID: {{ member.id }}
            </div>
          </MemberNode>
        </GraphView>
      </DashboardContainer>

      <!-- 본인 일정 모달 (MySchedule 컴포넌트) -->
      <ExpandedScheduleModal v-if="showMyScheduleModal">
        <ModalContent style="max-width: 95vw; max-height: 95vh; width: 95vw; height: 95vh;">
          <ModalHeader>
            <h2>{{ myProfile.name }}의 일정</h2>
            <CloseBtn @click="closeMyScheduleModal">×</CloseBtn>
          </ModalHeader>
          <ModalBody style="height: calc(100% - 60px); overflow: auto;">
            <Schedule :hideSideBar="true" />
          </ModalBody>
        </ModalContent>
      </ExpandedScheduleModal>

      <!-- 확대된 일정 뷰 모달 (다른 구성원) -->
      <ExpandedScheduleModal v-if="showExpandedView">
        <ModalContent style="max-width: 95vw; max-height: 95vh; width: 95vw; height: 95vh;">
          <ModalHeader>
            <h2>{{ expandedMember.name }}의 일정</h2>
            <CloseBtn @click="closeExpandedView">×</CloseBtn>
          </ModalHeader>
          <ModalBody style="height: calc(100% - 60px); overflow: auto;">
            <Schedule 
              :hideSideBar="true" 
              :memberData="expandedMember"
              :isOtherMember="true"
            />


                

                



          </ModalBody>
        </ModalContent>
      </ExpandedScheduleModal>
    </ContentArea>
    
    <!-- 빠른 미팅 요청 모달 -->
    <QuickMeetingModal v-if="showQuickModal" @close="closeQuickModal">
      <ModalOverlay @click="closeQuickModal" />
      <QuickModalContent @click.stop>
        <QuickModalHeader>
          <h3>{{ expandedMember ? expandedMember.name : '' }}님께 미팅 요청</h3>
          <CloseButton @click="closeQuickModal">×</CloseButton>
        </QuickModalHeader>
        
        <QuickModalBody>
          <TimeInfo>
            📅 {{ formatDate(getCurrentSelectedDate()) }} {{ selectedTimeSlot ? selectedTimeSlot.hour : 0 }}:00 - {{ selectedTimeSlot ? (selectedTimeSlot.endHour || selectedTimeSlot.hour + 1) : 1 }}:00
          </TimeInfo>
          
          <QuickForm>
            <input 
              v-model="quickMeeting.title" 
              type="text" 
              placeholder="미팅 제목 (예: 프로젝트 논의)"
              @keyup.enter="addQuickMeeting"
            />
            
            <select v-model="quickMeeting.type">
              <option value="meeting">🤝 회의</option>
              <option value="work">💼 업무 협업</option>
              <option value="other">📋 기타</option>
            </select>
          </QuickForm>
        </QuickModalBody>
        
        <QuickModalFooter>
          <CancelButton @click="closeQuickModal">취소</CancelButton>
          <AddButton @click="addQuickMeeting" :disabled="!quickMeeting.title">
            요청 보내기
          </AddButton>
        </QuickModalFooter>
      </QuickModalContent>
    </QuickMeetingModal>
    
    <!-- 일정 수정 모달 -->
    <QuickMeetingModal v-if="showEditModal" @close="closeEditModal">
      <ModalOverlay @click="closeEditModal" />
      <QuickModalContent @click.stop>
        <QuickModalHeader>
          <h3>일정 수정</h3>
          <CloseButton @click="closeEditModal">×</CloseButton>
        </QuickModalHeader>
        
        <QuickModalBody>
          <TimeInfo>
            📅 {{ formatDate(getCurrentSelectedDate()) }} {{ editingSchedule.hour }}:00 - {{ editingSchedule.hour + 1 }}:00
          </TimeInfo>
          
          <QuickForm>
            <input 
              v-model="editingSchedule.title" 
              type="text" 
              placeholder="일정 제목"
              @keyup.enter="updateScheduleRequest"
            />
            
            <select v-model="editingSchedule.type">
              <option value="meeting">🤝 회의</option>
              <option value="work">💼 업무 협업</option>
              <option value="other">📋 기타</option>
            </select>
          </QuickForm>
        </QuickModalBody>
        
        <QuickModalFooter>
          <CancelButton @click="closeEditModal">취소</CancelButton>
          <AddButton @click="updateScheduleRequest" :disabled="!editingSchedule.title">
            수정 완료
          </AddButton>
        </QuickModalFooter>
      </QuickModalContent>
    </QuickMeetingModal>

    <!-- 로딩 상태 표시 -->
    <div v-if="isLoadingTeamMembers" class="loading-overlay">
      <div class="loading-spinner">팀원 데이터 로딩 중...</div>
    </div>

    <!-- 에러 상태 표시 -->
    <div v-if="teamMembersError" class="error-message">
      <p>팀원 데이터 로드 실패: {{ teamMembersError }}</p>
      <button @click="loadTeamMembers" class="retry-button">다시 시도</button>
    </div>
  </ScheduleSharingRoot>
</template>

<script>
import SideBar from '@/components/SideBar.vue';
import Schedule from '@/views/myschedule/Schedule/Schedule.vue';
// import { useScheduleStore } from '@/stores/schedule'
import  * as ScheduleSharing from './ScheduleSharing.styled.js';
import { 
  handleDragStart as dragStartHandler, 
  handleDragEnd as dragEndHandler, 
  handleDragOver as dragOverHandler, 
  handleNodeDrop as nodeDropHandler, 
  updateNodePosition as updateNodePositionHandler 
} from './utils/drag.js';
import {
  editScheduleRequest as editScheduleRequestHandler,
  updateScheduleRequest as updateScheduleRequestHandler,
  deleteScheduleRequest as deleteScheduleRequestHandler,
  isMySchedule as isMyScheduleHandler,
  editSchedule as editScheduleHandler,
  getScheduleColorByType as getScheduleColorByTypeHandler,
  getScheduleType as getScheduleTypeHandler
} from './utils/Edit_schedules.js';
import {
  loadMyProfile as loadMyProfileHandler,
  loadTeamMembers as loadTeamMembersHandler,
  loadFallbackTeamMembers as loadFallbackTeamMembersHandler,
  loadUsersByDepartment as loadUsersByDepartmentHandler,
  loadUserSchedules as loadUserSchedulesHandler,
  getPositionFromRole as getPositionFromRoleHandler
} from './utils/load.js';
import {
  expandMemberSchedule as expandMemberScheduleHandler,
  closeMyScheduleModal as closeMyScheduleModalHandler,
  closeExpandedView as closeExpandedViewHandler,
  showClockView as showClockViewHandler,
  showCalendarView as showCalendarViewHandler,
  openQuickMeetingModal as openQuickMeetingModalHandler,
  closeQuickModal as closeQuickModalHandler,
  closeEditModal as closeEditModalHandler,
  openDragScheduleModal as openDragScheduleModalHandler,
  openUserScheduleModal as openUserScheduleModalHandler
} from './utils/Modal.js';
export default {
  name: 'ScheduleSharing',
  components: {
    SideBar,
    Schedule,
    ...ScheduleSharing
  },
  //vue3
  // setup() {
  //   const scheduleStore = useScheduleStore()
  //   return { scheduleStore }
  // },
  
  //vue2
  data() {
    return {
      // 일정 유형별 색상 정의
      scheduleColors: {
        meeting: '#FF6B6B', // 회의용 빨간색
        work: '#4ECDC4',    // 업무용 청록색
        other: '#45B7D1'    // 기타용 파란색
      },
      // 본인 프로필 (로그인 정보에서 가져옴)
      myProfile: {
        id: 0,
        name: '로딩 중...',
        position: '',
        department: '',
        team: '',
        schedules: []
      },
      showMyScheduleModal: false,
      departments: [
        // 부서 네모 크기 확대 (width와 height를 350에서 450으로 증가)
        // 부서 간 간격도 넓힘 (x 값 조정)
        { name: '개발부', x: 150, y: 300, width: 450, height: 450 },
        { name: '디자인부', x: 650, y: 300, width: 450, height: 450 },
        { name: '마케팅부', x: 1150, y: 300, width: 450, height: 450 }
      ],
      teamMembers: [], // DB에서 부서별로 데이터를 불러올 예정
      selectedMember: null,
      showExpandedView: false,
      expandedMember: null,
      activeTab: 'clock',
      selectedSchedule: null,
      selectedDate: null, // 선택된 날짜 추가
      memberSelectedDates: {}, // 구성원별 선택된 날짜 관리
      currentYear: 2025,
      currentMonth: 6, // 7월 (0-based index)
      weekdays: ['일', '월', '화', '수', '목', '금', '토'],
      draggedMember: null,
      nodePositions: {},
      // 오늘 날짜로 설정 (예시 일정 표시를 위해)
      today: '2023-06-15',
      // 간단한 프로젝트 데이터 - 2025년 7월로 설정
      projects: [
        {
          id: 1,
          title: '웹사이트 리뉴얼',
          startDate: '2025-07-08',
          endDate: '2025-07-15',
          color: '#B4D3F5'
        },
        {
          id: 2,
          title: '모바일 앱 개발',
          startDate: '2025-07-20',
          endDate: '2025-07-25',
          color: '#FF6B6B'
        }
      ],
      // 빠른 미팅 요청 관련
      showQuickModal: false,
      selectedTimeSlot: null,
      quickMeeting: {
        title: '',
        type: 'meeting'
      },
      // 일정 수정 관련
      showEditModal: false,
      editingSchedule: {
        id: null,
        title: '',
        type: 'meeting',
        hour: 0
      },
      // 드래그 투 셀렉트 관련
      isDragging: false,
      dragStartAngle: null,
      dragCurrentAngle: null,
      selectedTimeRange: {
        start: null,
        end: null
      },
      // API 엔드포인트
      API_ENDPOINTS: {
        USERS_ALL: '/api/schedule-sharing/users/all',
        USERS_BY_DEPARTMENT: '/api/schedule-sharing/users/department',
        SCHEDULES_BY_USER: '/api/schedules'
      },
      // 로딩 상태
      isLoadingTeamMembers: false,
      teamMembersError: null,
    };
  },
  methods: {
    // 드래그 관련 메서드들 (drag.js에서 import)
    handleDragStart(event, memberId) {
      dragStartHandler(event, memberId, this);
    },
    handleDragEnd() {
      dragEndHandler(this);
    },
    handleDragOver(event) {
      dragOverHandler(event);
    },
    handleDrop(event) {
      nodeDropHandler(event, this.draggedMember, (memberId, left, top) => {
        updateNodePositionHandler(memberId, left, top, this);
      });
    },

    // 로그인된 사용자 정보 로드 (load.js에서 import)
    async loadMyProfile() {
      await loadMyProfileHandler(this, getScheduleTypeHandler, getScheduleColorByTypeHandler, getPositionFromRoleHandler);
    },
    // 부서별 팀원 데이터 로드 (load.js에서 import)
    async loadTeamMembers() {
      await loadTeamMembersHandler(this, getScheduleColorByTypeHandler, getScheduleTypeHandler);
    },
    // 폴백 데이터 로드 (API 실패 시 사용) (load.js에서 import)
    loadFallbackTeamMembers() {
      loadFallbackTeamMembersHandler(this);
    },
    // 역할(role)에서 직급 추출 (load.js에서 import)
    getPositionFromRole(role) {
      return getPositionFromRoleHandler(role);
    },
    // 부서별 사용자 조회 (필요시) (load.js에서 import)
    async loadUsersByDepartment(department) {
      return await loadUsersByDepartmentHandler(department);
    },
    // 특정 사용자의 일정 조회 (필요시) (load.js에서 import)
    async loadUserSchedules(userId) {
      return await loadUserSchedulesHandler(userId, this, getScheduleColorByTypeHandler, getScheduleTypeHandler);
    },

    // 일정 제목으로 타입 추정 (Edit_schedules.js에서 import)
    getScheduleType(title) {
      return getScheduleTypeHandler(title);
    },
    // 일정 유형에 따라 색상 자동 할당 (Edit_schedules.js에서 import)
    getScheduleColorByType(type) {
      return getScheduleColorByTypeHandler(type, this);
    },
    // 일정 수정 관련 메서드들 (Edit_schedules.js에서 import)
    editScheduleRequest(schedule) {
      editScheduleRequestHandler(schedule, this);
    },
    closeEditModal() {
      closeEditModalHandler(this);
    },
    updateScheduleRequest() {
      updateScheduleRequestHandler(this);
    },
    deleteScheduleRequest(schedule) {
      deleteScheduleRequestHandler(schedule, this);
    },
    isMySchedule(schedule) {
      return isMyScheduleHandler(schedule, this);
    },
    editSchedule() {
      if (!this.selectedSchedule) return;
      editScheduleHandler(this);
    },

    // 본인 노드 위치 계산 메서드 추가
    getMyNodePosition() {
      // 개발부 네모 위쪽에 고정 위치
      const devDepartment = this.departments.find(d => d.name === '개발부');
      const x = devDepartment.x + devDepartment.width / 2 - 50; // 50px 왼쪽으로 이동
      const y = devDepartment.y - 170; // 150px에서 200px로 변경하여 더 위로 올림
      
      return { left: `${x}px`, top: `${y}px` };
    },
    getDepartmentBoundaryStyle(department) {
      return {
        left: `${department.x}px`,
        top: `${department.y}px`,
        width: `${department.width}px`,
        height: `${department.height}px`
      };
    },
    getNodePosition(index) {
      // 항상 초기 위치 계산 (nodePositions 무시)
      // 부서별로 그룹화하여 초기 위치 설정
      const member = this.teamMembers[index];
      const department = member.department;
      
      // 부서별 위치 계산
      let baseX = 0;
      let baseY = 0;
      let offsetX = 0;
      let offsetY = 0;
      
      // 부서별 그룹 위치
      const departmentObj = this.departments.find(d => d.name === department);
      if (departmentObj) {
        // 부서 네모 중앙에 배치하되 아래로 내림
        baseX = departmentObj.x + departmentObj.width / 2 - 120; // 중앙에서 더 왼쪽으로
        baseY = departmentObj.y + departmentObj.height / 2 - 50; // 중앙에서 덜 위로 (120에서 50으로 변경하여 아래로 내림)
      }
      
      // 같은 부서 내에서 위치 조정
      const departmentMembers = this.teamMembers.filter(m => m.department === department);
      const memberIndexInDepartment = departmentMembers.findIndex(m => m.id === member.id);
      
      // 부서 내 행과 열 계산 (2x2 그리드 유지)
      const row = Math.floor(memberIndexInDepartment / 2);
      const col = memberIndexInDepartment % 2;
      
      // 간격을 크게 늘림
      offsetX = col * 240; // 간격을 240px로 크게 늘림
      offsetY = row * 220; // 간격을 220px로 크게 늘림
      
      const x = baseX + offsetX;
      const y = baseY + offsetY;
      
      return { left: `${x}px`, top: `${y}px` };
    },

    // 멤버 선택 및 모달 관련
    selectMember(memberId) {
      console.log('🚀 노드 클릭됨! memberId:', memberId);
      
      // 기존 로직
      this.selectedMember = memberId;
      
      // 모달도 함께 열기 (기본적으로 시계 뷰로)
      this.expandMemberSchedule(memberId);
      
      // API 호출
      loadUserSchedulesHandler(memberId, this, getScheduleColorByTypeHandler, getScheduleTypeHandler)
        .then(schedules => {
          console.log('✅ API 응답 성공:', schedules);
          
          const member = this.teamMembers.find(m => m.id === memberId);
          if (member) {
            member.schedules = schedules;
            console.log(`🎉 ${member.name}의 일정 데이터 로드 완료`);
          }
        })
        .catch(error => {
          console.error('❌ 사용자 일정 로드 오류:', error);
        });
    },
    expandMemberSchedule(memberId) {
      expandMemberScheduleHandler(memberId, this);
    },
    closeMyScheduleModal() {
      closeMyScheduleModalHandler(this);
    },
    closeExpandedView() {
      closeExpandedViewHandler(this);
    },
    showClockView(memberId) {
      showClockViewHandler(memberId, this);
    },
    showCalendarView(memberId) {
      showCalendarViewHandler(memberId, this);
    },
    openUserSchedule(user) {
      openUserScheduleModalHandler(user, this);
    },

    // 일정 관련 메서드
    getTimeSchedulePreview(member) {
      const today = new Date().toISOString().split('T')[0];
      const todaySchedules = member.schedules.filter(s => s.start.startsWith(today));
      if (todaySchedules.length === 0) return '오늘 일정 없음';
      return `오늘 ${todaySchedules.length}개 일정`;
    },
    getCalendarPreview(member) {
      const currentMonth = new Date().getMonth() + 1;
      const currentMonthSchedules = member.schedules.filter(s => {
        const scheduleMonth = new Date(s.start).getMonth() + 1;
        return scheduleMonth === currentMonth;
      });
      return `이번 달 ${currentMonthSchedules.length}개 일정`;
    },
    getHourSchedules(hour) {
      if (!this.expandedMember) return [];
      
      return this.expandedMember.schedules.filter(schedule => {
        const scheduleHour = new Date(schedule.start).getHours();
        return scheduleHour === hour;
      });
    },
    getDaySchedules(dateString) {
      if (!this.expandedMember) return [];
      
      const date = dateString.split('T')[0];
      return this.expandedMember.schedules.filter(schedule => {
        return schedule.start.startsWith(date);
      });
    },
    prevMonth() {
      if (this.currentMonth === 0) {
        this.currentYear--;
        this.currentMonth = 11;
      } else {
        this.currentMonth--;
      }
    },
    nextMonth() {
      if (this.currentMonth === 11) {
        this.currentYear++;
        this.currentMonth = 0;
      } else {
        this.currentMonth++;
      }
    },
    addSchedule() {
      // 실제 구현에서는 일정 추가 모달 표시
      alert('일정 추가 기능은 아직 구현되지 않았습니다.');
      // 일정 추가 시 유형에 따라 색상 자동 설정 예시
      // const newSchedule = {
      //   id: Math.max(...this.expandedMember.schedules.map(s => s.id)) + 1,
      //   title: '새 일정',
      //   start: '2023-06-15T10:00',
      //   end: '2023-06-15T11:00',
      //   type: 'meeting' // 또는 'work', 'other'
      // };
      // newSchedule.color = getScheduleColorByTypeHandler(newSchedule.type, this);
      // this.expandedMember.schedules.push(newSchedule);
    },

    // 시계 관련 메서드
    getHourMarkStyle(hour) {
      const angle = (hour / 12) * 2 * Math.PI - Math.PI / 2;
      const x = Math.cos(angle) * 18;
      const y = Math.sin(angle) * 18;
      return {
        left: `${x + 20}px`,
        top: `${y + 20}px`
      };
    },
    getScheduleSegmentStyle(schedule) {
      try {
        console.log(`[미니시계] 일정 스타일 생성: ${schedule.title}, 색상: ${schedule.color}`);
        
        // 시간 파싱 - Date 객체 사용
        const startDate = new Date(schedule.start);
        const endDate = new Date(schedule.end);
        
        // 시간과 분을 직접 추출
        const startHour = startDate.getHours();
        const startMinute = startDate.getMinutes();
        const endHour = endDate.getHours();
        const endMinute = endDate.getMinutes();
        
        // 시계 위치 계산 (12시간제로 변환)
        const startHour12 = startHour % 12 === 0 ? 12 : startHour % 12;
        const endHour12 = endHour % 12 === 0 ? 12 : endHour % 12;
        
        // 시작 각도 (12시 방향이 0도, 시계 방향으로 회전)
        // 시간은 30도씩, 분은 0.5도씩 (30도/60분 = 0.5도/분)
        // CSS conic-gradient는 3시 방향이 0도이므로 -90도 조정
        // 추가로 90도를 시계방향으로 더 돌림
        const startDegree = (startHour12 * 30) + (startMinute * 0.5);
        const startAngle = startDegree; // -90 + 90 = 0 (상쇄)
        
        // 종료 각도
        const endDegree = (endHour12 * 30) + (endMinute * 0.5);
        let endAngle = endDegree; // -90 + 90 = 0 (상쇄)
        
        // 종료 각도가 시작 각도보다 작으면 360도 더함 (다음 날로 넘어가는 경우)
        if (endAngle < startAngle) {
          endAngle += 360;
        }
        
        // 각도 차이
        const sweepDeg = endAngle - startAngle;
        
        // 디버깅용 로그
        console.log(`미니 시계 - ${schedule.title}: ${startHour}:${startMinute}(${startHour12}) ~ ${endHour}:${endMinute}(${endHour12})`);
        console.log(`시작 시간: ${startHour12}시 ${startMinute}분 => ${startDegree}도 => ${startAngle}도`);
        console.log(`종료 시간: ${endHour12}시 ${endMinute}분 => ${endDegree}도 => ${endAngle}도`);
        console.log(`회전 각도: ${sweepDeg}도`);
        
        // 부채꼴 스타일
        const style = {
          position: 'absolute',
          width: '100%',
          height: '100%',
          left: '0',
          top: '0',
          background: `conic-gradient(from ${startAngle}deg, ${schedule.color} 0deg, ${schedule.color} ${sweepDeg}deg, transparent ${sweepDeg}deg 360deg)`,
          borderRadius: '50%',
          opacity: '0.9',
          zIndex: 10
        };
        
        console.log(`[미니시계] 생성된 스타일:`, style);
        return style;
      } catch (error) {
        console.error('일정 스타일 계산 오류:', error);
        return {
          display: 'none'
        };
      }
    },
    getExpandedScheduleStyle(schedule) {
      try {
        // 시간 파싱 - Date 객체 사용
        const startDate = new Date(schedule.start);
        const endDate = new Date(schedule.end);
        
        // 시간과 분을 직접 추출
        const startHour = startDate.getHours();
        const startMinute = startDate.getMinutes();
        const endHour = endDate.getHours();
        const endMinute = endDate.getMinutes();
        
        // 시계 위치 계산 (12시간제로 변환)
        const startHour12 = startHour % 12 === 0 ? 12 : startHour % 12;
        const endHour12 = endHour % 12 === 0 ? 12 : endHour % 12;
        
        // 시작 각도 (12시 방향이 0도, 시계 방향으로 회전)
        // 시간은 30도씩, 분은 0.5도씩 (30도/60분 = 0.5도/분)
        // CSS conic-gradient는 3시 방향이 0도이므로 -90도 조정
        // 추가로 90도를 시계방향으로 더 돌림
        const startDegree = (startHour12 * 30) + (startMinute * 0.5);
        const startAngle = startDegree; // -90 + 90 = 0 (상쇄)
        
        // 종료 각도
        const endDegree = (endHour12 * 30) + (endMinute * 0.5);
        let endAngle = endDegree; // -90 + 90 = 0 (상쇄)
        
        // 종료 각도가 시작 각도보다 작으면 360도 더함 (다음 날로 넘어가는 경우)
        if (endAngle < startAngle) {
          endAngle += 360;
        }
        
        // 각도 차이
        const sweepDeg = endAngle - startAngle;
        
        // 디버깅용 로그
        console.log(`확대 시계 - ${schedule.title}: ${startHour}:${startMinute}(${startHour12}) ~ ${endHour}:${endMinute}(${endHour12})`);
        console.log(`시작 시간: ${startHour12}시 ${startMinute}분 => ${startDegree}도 => ${startAngle}도`);
        console.log(`종료 시간: ${endHour12}시 ${endMinute}분 => ${endDegree}도 => ${endAngle}도`);
        console.log(`회전 각도: ${sweepDeg}도`);
        
        // 부채꼴 스타일
        return {
          position: 'absolute',
          width: '100%',
          height: '100%',
          left: '0',
          top: '0',
          background: `conic-gradient(from ${startAngle}deg, ${schedule.color} 0deg, ${schedule.color} ${sweepDeg}deg, transparent ${sweepDeg}deg 360deg)`,
          borderRadius: '50%',
          opacity: '0.8',
          zIndex: 1
        };
      } catch (error) {
        console.error('확대 일정 스타일 계산 오류:', error);
        return {
          display: 'none'
        };
      }
    },
    getHourNumberStyle(hour) {
      // 시계 숫자 위치 계산 (12시는 위쪽, 시계 방향으로 회전)
      const angle = ((hour % 12) / 12) * 2 * Math.PI - Math.PI / 2;
      const radius = 34; // 반지름을 키워서 80px 시계에 맞게 조정
      const centerX = 38; // 시계 중심 X 좌표 (80px 시계의 중심)
      const centerY = 38; // 시계 중심 Y 좌표 (80px 시계의 중심)
      
      const x = Math.cos(angle) * radius + centerX;
      const y = Math.sin(angle) * radius + centerY;
      
      return {
        position: 'absolute',
        left: `${x}px`,
        top: `${y}px`,
        transform: 'translate(-50%, -50%)',
        fontSize: '9px',
        fontWeight: 'bold',
        zIndex: 20,
        color: '#333',
        textShadow: '0 0 2px white'
      };
    },
    getExpandedHourNumberStyle(hour) {
      // 확대된 시계 숫자 위치 계산 (500px 시계에 맞게 조정)
      const angle = ((hour % 12) / 12) * 2 * Math.PI - Math.PI / 2;
      const radius = 225; // 반지름을 조정하여 모서리쪽에 배치
      const centerX = 250; // 시계 중심 X 좌표 (500px 시계의 중심은 250)
      const centerY = 250; // 시계 중심 Y 좌표 (500px 시계의 중심은 250)
      
      const x = Math.cos(angle) * radius + centerX;
      const y = Math.sin(angle) * radius + centerY;
      
      return {
        position: 'absolute',
        left: `${x}px`,
        top: `${y}px`,
        transform: 'translate(-50%, -50%)',
        fontSize: '28px',
        fontWeight: 'bold',
        zIndex: 20,
        color: '#333',
        textShadow: '0 0 3px white'
      };
    },
    getTimeSlotStyle(hour) {
      const angle = (hour / 12) * 2 * Math.PI - Math.PI / 2;
      const radius = 225; // 시간 숫자와 동일한 반지름
      const centerX = 250; // 시계 중심 X 좌표
      const centerY = 250; // 시계 중심 Y 좌표
      
      const x = Math.cos(angle) * radius + centerX;
      const y = Math.sin(angle) * radius + centerY;
      
      return {
        position: 'absolute',
        left: `${x}px`,
        top: `${y}px`,
        transform: 'translate(-50%, -50%)'
      };
    },
    getEmptyTimeSlotStyle(slot) {
      const hour = slot.hour > 12 ? slot.hour - 12 : slot.hour;
      const angle = (hour / 12) * 2 * Math.PI - Math.PI / 2;
      const radius = 35; // 일정보다 안쪽에 위치
      const x = Math.cos(angle) * radius;
      const y = Math.sin(angle) * radius;
      
      return {
        position: 'absolute',
        left: `${x + 50}px`,
        top: `${y + 50}px`,
        transform: 'translate(-50%, -50%)'
      };
    },

    // 오늘의 일정 관련
    getTodaySchedules(member) {
      // member나 schedules가 null/undefined인 경우 빈 배열 반환
      if (!member || !member.schedules || !Array.isArray(member.schedules)) {
        return [];
      }
      
      // 본인인 경우 실제 오늘 날짜 사용, 다른 사람은 하드코딩된 날짜 사용
      let today;
      // if (member.id === 0) {
        // 본인 노드인 경우 실제 오늘 날짜
        today = new Date().toISOString().split('T')[0];
      // } else {
      //   // 다른 팀원들은 예시 데이터를 위해 하드코딩된 날짜 사용
      //   today = this.today;
      // }
      
      const schedules = member.schedules.filter(s => {
        // if (member.id === 0) {
          // 시간대 문제를 해결하기 위해 로컬 날짜로 비교
          const scheduleDate = new Date(s.start);
          const localDateString = `${scheduleDate.getFullYear()}-${String(scheduleDate.getMonth() + 1).padStart(2, '0')}-${String(scheduleDate.getDate()).padStart(2, '0')}`;
          return localDateString === today;
        // } else {
        //   return s.start.startsWith(today);
        // }
      });
      
      return schedules;
    },
    getSelectedDateSchedules() {
      if (!this.expandedMember) return [];
      
      const memberId = this.expandedMember.id;
      const targetDate = this.memberSelectedDates[memberId] || this.today;
      
      return this.expandedMember.schedules.filter(s => s.start.startsWith(targetDate));
    },
    getCurrentSelectedDate() {
      if (!this.expandedMember) return this.today;
      return this.memberSelectedDates[this.expandedMember.id] || this.today;
    },

    // 달력 관련 메서드
    getMiniCalendarDaySchedule(member, dayIndex) {
      // 본인 노드인 경우 현재 실제 날짜 기준, 다른 사람은 하드코딩된 달력 기준
      let year, month;
      // if (member.id === 0) {
        // 본인은 현재 실제 날짜 기준
        const now = new Date();
        year = now.getFullYear();
        month = now.getMonth();
      // // } else {
      //   // 다른 팀원들은 하드코딩된 달력 기준
      //   year = this.currentYear;
      //   month = this.currentMonth;
      // }
      
      const currentDate = new Date(year, month, 1);
      const firstDayOfMonth = currentDate.getDay();
      const adjustedIndex = dayIndex - firstDayOfMonth;
      
      if (adjustedIndex < 0 || adjustedIndex > 30) return null;
      
      const day = new Date(year, month, adjustedIndex + 1);
      const dateString = day.toISOString().split('T')[0];
      
      const scheduleForDay = member.schedules.find(s => {
        // if (member.id === 0) {
          // 시간대 문제를 해결하기 위해 로컬 날짜로 비교
          const scheduleDate = new Date(s.start);
          const localDateString = `${scheduleDate.getFullYear()}-${String(scheduleDate.getMonth() + 1).padStart(2, '0')}-${String(scheduleDate.getDate()).padStart(2, '0')}`;
          return localDateString === dateString;
        // } else {
        //   return s.start.startsWith(dateString);
        // }
      });
      return scheduleForDay || null;
    },
    getMiniCalendarDayNumber(index) {
      // 나의일정수정에서 선택된 날짜를 기준으로 하거나, 없으면 현재 월 기준
      const selectedDate = this.getSelectedDateFromSchedule();
      const currentYear = selectedDate.getFullYear();
      const currentMonth = selectedDate.getMonth();
      
      // index는 1부터 시작하므로 0부터 시작하도록 조정
      const zeroBasedIndex = index - 1;
      
      const currentDate = new Date(currentYear, currentMonth, 1);
      const firstDayOfMonth = currentDate.getDay();
      const adjustedIndex = zeroBasedIndex - firstDayOfMonth;
      
      if (adjustedIndex < 0) {
        // 이전 달의 날짜
        const prevMonthLastDate = new Date(currentYear, currentMonth, 0).getDate();
        return prevMonthLastDate + adjustedIndex + 1;
      } else {
        // 현재 달의 날짜 (다음 달로 넘어가는 경우도 현재 달 기준으로 계산)
        const currentMonthLastDate = new Date(currentYear, currentMonth + 1, 0).getDate();
        if (adjustedIndex + 1 <= currentMonthLastDate) {
          return adjustedIndex + 1;
        } else {
          // 다음 달의 날짜
          return adjustedIndex + 1 - currentMonthLastDate;
        }
      }
    },
    getMiniCalendarDate(index) {
      // 나의일정수정에서 선택된 날짜를 기준으로 하거나, 없으면 현재 월 기준
      const selectedDate = this.getSelectedDateFromSchedule();
      const currentYear = selectedDate.getFullYear();
      const currentMonth = selectedDate.getMonth();
      
      // index는 1부터 시작하므로 0부터 시작하도록 조정
      const zeroBasedIndex = index - 1;
      
      const currentDate = new Date(currentYear, currentMonth, 1);
      const firstDayOfMonth = currentDate.getDay();
      const adjustedIndex = zeroBasedIndex - firstDayOfMonth;
      
      if (adjustedIndex < 0) {
        // 이전 달의 날짜
        const prevMonthYear = currentMonth === 0 ? currentYear - 1 : currentYear;
        const prevMonth = currentMonth === 0 ? 11 : currentMonth - 1;
        const prevMonthLastDate = new Date(currentYear, currentMonth, 0).getDate();
        const day = prevMonthLastDate + adjustedIndex + 1;
        return `${prevMonthYear}-${String(prevMonth + 1).padStart(2, '0')}-${String(day).padStart(2, '0')}`;
      } else {
        // 현재 달의 날짜 (다음 달로 넘어가는 경우도 처리)
        const currentMonthLastDate = new Date(currentYear, currentMonth + 1, 0).getDate();
        if (adjustedIndex + 1 <= currentMonthLastDate) {
          const day = adjustedIndex + 1;
          return `${currentYear}-${String(currentMonth + 1).padStart(2, '0')}-${String(day).padStart(2, '0')}`;
        } else {
          // 다음 달의 날짜
          const nextMonthYear = currentMonth === 11 ? currentYear + 1 : currentYear;
          const nextMonth = currentMonth === 11 ? 0 : currentMonth + 1;
          const day = adjustedIndex + 1 - currentMonthLastDate;
          return `${nextMonthYear}-${String(nextMonth + 1).padStart(2, '0')}-${String(day).padStart(2, '0')}`;
        }
      }
    },
    
    formatTime(dateTimeString) {
      const date = new Date(dateTimeString);
      return date.toLocaleTimeString('ko-KR', { hour: '2-digit', minute: '2-digit', hour12: false });
    },
    
    formatDate(dateTimeString) {
      const date = new Date(dateTimeString);
      return date.toLocaleDateString('ko-KR', { month: 'short', day: 'numeric' });
    },
    
    formatFullDate(dateString) {
      const date = new Date(dateString);
      const year = date.getFullYear();
      const month = date.getMonth() + 1;
      const day = date.getDate();
      const weekdays = ['일', '월', '화', '수', '목', '금', '토'];
      const weekday = weekdays[date.getDay()];
      return `${year}년 ${month}월 ${day}일 ${weekday}요일의 일정`;
    },
    
    selectSchedule(schedule) {
      this.selectedSchedule = schedule;
    },
    
    // 날짜 선택 메서드 - 현재 확대된 구성원의 날짜만 변경
    selectDate(dateString) {
      if (this.expandedMember) {
        this.memberSelectedDates[this.expandedMember.id] = dateString;
        // selectedDate도 업데이트
        this.selectedDate = dateString;
        
        // 선택된 날짜의 일정이 있으면 첫 번째 일정을 선택
        const selectedDateSchedules = this.getDaySchedules(dateString);
        if (selectedDateSchedules.length > 0) {
          this.selectedSchedule = selectedDateSchedules[0];
        } else {
          this.selectedSchedule = null;
        }
      }
    },
    
    getWeekSchedules() {
      if (!this.expandedMember) return [];
      
      const weeks = [];
      const daysInMonth = new Date(this.currentYear, this.currentMonth + 1, 0).getDate();
      
      // 주 단위로 일정 그룹화
      for (let weekStart = 1; weekStart <= daysInMonth; weekStart += 7) {
        const weekEnd = Math.min(weekStart + 6, daysInMonth);
        const weekDays = [];
        
        for (let day = weekStart; day <= weekEnd; day++) {
          const dateString = `${this.currentYear}-${String(this.currentMonth + 1).padStart(2, '0')}-${String(day).padStart(2, '0')}`;
          const daySchedules = this.expandedMember.schedules.filter(s => s.start.startsWith(dateString));
          weekDays.push(daySchedules);
        }
        
        // 주의 마지막 날이 월의 마지막 날이 아니면 빈 배열로 채움
        const remainingDays = 7 - weekDays.length;
        for (let i = 0; i < remainingDays; i++) {
          weekDays.push([]);
        }
        
        weeks.push(weekDays);
      }
      
      return weeks;
    },
    
    // 나의일정수정에서 선택된 날짜 가져오기 (localStorage 활용)
    getSelectedDateFromSchedule() {
      try {
        const savedDate = localStorage.getItem('myScheduleSelectedDate');
        if (savedDate) {
          const date = new Date(savedDate);
          // 유효한 날짜인지 확인
          if (!isNaN(date.getTime())) {
            return date;
          }
        }
      } catch (error) {
        console.warn('로컬스토리지에서 선택된 날짜를 가져오는 중 오류:', error);
      }
      // 저장된 날짜가 없거나 유효하지 않으면 현재 날짜 반환
      return new Date();
    },
    
    // 프로젝트가 해당 날짜에 진행 중인지 확인
    isProjectActive(dateString) {
      return this.projects.some(project => {
        const date = new Date(dateString);
        const startDate = new Date(project.startDate);
        const endDate = new Date(project.endDate);
        return date >= startDate && date <= endDate;
      });
    },
    
    // 해당 날짜의 프로젝트 가져오기
    getProjectForDate(dateString) {
      return this.projects.find(project => {
        const date = new Date(dateString);
        const startDate = new Date(project.startDate);
        const endDate = new Date(project.endDate);
        return date >= startDate && date <= endDate;
      });
    },
    
    // 프로젝트 날짜 포맷팅
    formatProjectDate(dateString) {
      const date = new Date(dateString);
      const month = date.getMonth() + 1;
      const day = date.getDate();
      return `${month}월 ${day}일`;
    },

    // 선택된 날짜 표시 포맷 메서드 추가
    formatSelectedDate(dateString) {
      const date = new Date(dateString);
      const year = date.getFullYear();
      const month = date.getMonth() + 1;
      const day = date.getDate();
      const weekdays = ['일', '월', '화', '수', '목', '금', '토'];
      const weekday = weekdays[date.getDay()];
      return `${year}년 ${month}월 ${day}일 ${weekday}요일`;
    },
    
    // 페이지 포커스 리스너 설정
    setupPageFocusListener() {
      this.pageFocusListener = () => {
        // 페이지에 포커스가 돌아왔을 때 미니달력 강제 업데이트
        this.$forceUpdate();
      };
      window.addEventListener('focus', this.pageFocusListener);
    },
    
    // 빈 시간 슬롯 가져오기 (시간 제한 없이 모든 빈 시간)
    getEmptyTimeSlots() {
      if (!this.expandedMember) return [];
      
      const emptySlots = [];
      const existingSchedules = this.getSelectedDateSchedules();
      
      // 0시부터 23시까지 모든 시간 체크
      for (let hour = 0; hour < 24; hour++) {
        const hasSchedule = existingSchedules.some(schedule => {
          const startHour = new Date(schedule.start).getHours();
          const endHour = new Date(schedule.end).getHours();
          return hour >= startHour && hour < endHour;
        });
        
        if (!hasSchedule) {
          emptySlots.push({ hour });
        }
      }
      
      return emptySlots;
    },
    
    // 미팅 요청 추가
    addMeetingRequest(timeSlot) {
      openQuickMeetingModalHandler(timeSlot, this);
    },
    
    // 빠른 모달 닫기
    closeQuickModal() {
      closeQuickModalHandler(this);
    },
    
    // 빠른 미팅 추가
    async addQuickMeeting() {
      try {
        const startHour = this.selectedTimeSlot.hour;
        const endHour = this.selectedTimeSlot.endHour || (this.selectedTimeSlot.hour + 1);
        
        const startTime = `${startHour.toString().padStart(2, '0')}:00`;
        const endTime = `${endHour.toString().padStart(2, '0')}:00`;
        
        const meetingRequest = {
          title: `[요청] ${this.quickMeeting.title}`,
          start: `${this.getCurrentSelectedDate()}T${startTime}`,
          end: `${this.getCurrentSelectedDate()}T${endTime}`,
          type: this.quickMeeting.type,
          color: '#4CAF50', // 요청 일정은 초록색으로 고정
          isRequest: true, // 요청 플래그
          requestedBy: this.myProfile.name // 요청자
        };
        
        // 해당 멤버의 일정에 바로 추가 (실제로는 API 호출)
        this.expandedMember.schedules.push({
          ...meetingRequest,
          id: Date.now() // 임시 ID
        });
        
        alert(`${this.expandedMember.name}님께 미팅 요청을 보냈습니다!`);
        this.closeQuickModal();
        
      } catch (error) {
        console.error('미팅 요청 실패:', error);
        alert('미팅 요청에 실패했습니다.');
      }
    },
    

    
    // 드래그 투 셀렉트 관련 메서드들
    onDragStart(event) {
      if (event.target.closest('.schedule-segment') || event.target.closest('.schedule-tooltip')) return;
      
      this.isDragging = true;
      const angle = this.getAngleFromEvent(event);
      this.dragStartAngle = angle;
      this.dragCurrentAngle = angle;
      this.selectedTimeRange = { start: null, end: null };
    },
    
    onDragMove(event) {
      if (!this.isDragging) return;
      
      const angle = this.getAngleFromEvent(event);
      this.dragCurrentAngle = angle;
    },
    
    onDragEnd() {
      if (!this.isDragging) return;
      
      this.isDragging = false;
      
      // 드래그 범위 계산
      const startHour = this.angleToHour(this.dragStartAngle);
      const endHour = this.angleToHour(this.dragCurrentAngle);
      
      const minHour = Math.min(startHour, endHour);
      const maxHour = Math.max(startHour, endHour);
      
      if (minHour !== maxHour) {
        this.selectedTimeRange = {
          start: minHour,
          end: maxHour
        };
        this.openDragScheduleModal();
      }
      
      this.dragStartAngle = null;
      this.dragCurrentAngle = null;
    },
    
    onClockFaceClick(event) {
      if (this.isDragging) return;
      
      const angle = this.getAngleFromEvent(event);
      const hour = this.angleToHour(angle);
      
      this.selectedTimeRange = {
        start: hour,
        end: hour + 1
      };
      this.openDragScheduleModal();
    },
    
    getAngleFromEvent(event) {
      // ExpandedClockFace 요소를 정확히 찾기
      let clockFace = event.currentTarget;
      
      // 만약 currentTarget이 올바르지 않다면 closest로 찾기
      if (!clockFace.classList || !clockFace.classList.contains('expanded-clock-face')) {
        clockFace = event.target.closest('.expanded-clock-face') || event.currentTarget;
      }
      
      const rect = clockFace.getBoundingClientRect();
      
      // 시계의 중심점 계산 (실제 시계 크기에 맞게)
      const centerX = rect.left + rect.width / 2;
      const centerY = rect.top + rect.height / 2;
      
      const x = event.clientX - centerX;
      const y = event.clientY - centerY;
      
      console.log('좌표 계산:', {
        clientX: event.clientX,
        clientY: event.clientY,
        rectLeft: rect.left,
        rectTop: rect.top,
        centerX: centerX,
        centerY: centerY,
        x: x,
        y: y
      });
      
      let angle = Math.atan2(y, x) * (180 / Math.PI);
      angle = (angle + 90 + 360) % 360; // 12시 방향을 0도로
      
      return angle;
    },
    
    angleToHour(angle) {
      const hour = Math.round(angle / 30); // 30도마다 1시간
      return hour === 0 ? 12 : hour;
    },
    
    getTimeSlotSelected(hour) {
      if (!this.selectedTimeRange.start || !this.selectedTimeRange.end) return false;
      return hour >= this.selectedTimeRange.start && hour <= this.selectedTimeRange.end;
    },
    
    getDragSelectionStyle() {
      if (!this.isDragging || !this.dragStartAngle || !this.dragCurrentAngle) return {};
      
      const startAngle = Math.min(this.dragStartAngle, this.dragCurrentAngle);
      const endAngle = Math.max(this.dragStartAngle, this.dragCurrentAngle);
      const angleDiff = endAngle - startAngle;
      
      // 시간표 크기에 100% 정확히 맞춤 - 테두리까지 완벽하게
      const centerX = 250; // 시계 중심 (500/2)
      const centerY = 250; // 시계 중심 (500/2)
      const radius = 250; // 시간표 전체 반지름 (500/2)
      
      // 디버깅용 로그
      console.log('드래그 원 크기:', {
        centerX,
        centerY,
        radius,
        left: centerX - radius,
        top: centerY - radius,
        width: radius * 2,
        height: radius * 2
      });
      
      return {
        position: 'absolute',
        left: `${centerX - radius}px`,
        top: `${centerY - radius}px`,
        width: `${radius * 2}px`,
        height: `${radius * 2}px`,
        background: `conic-gradient(from ${startAngle}deg, transparent 0deg, rgba(76, 175, 80, 0.8) 0deg, rgba(76, 175, 80, 0.8) ${angleDiff}deg, transparent ${angleDiff}deg)`,
        borderRadius: '50%',
        pointerEvents: 'none',
        zIndex: 5,
        border: 'none',
        boxShadow: '0 0 25px rgba(76, 175, 80, 0.6)'
      };
    },
    
    openDragScheduleModal() {
      openDragScheduleModalHandler(this);
    }
  },
  computed: {
    calendarDays() {
      const days = [];
      const firstDay = new Date(this.currentYear, this.currentMonth, 1);
      const lastDay = new Date(this.currentYear, this.currentMonth + 1, 0);
      
      // 이전 달의 마지막 날짜들
      const prevMonthLastDate = new Date(this.currentYear, this.currentMonth, 0).getDate();
      const firstDayOfWeek = firstDay.getDay();
      
      for (let i = firstDayOfWeek - 1; i >= 0; i--) {
        const prevMonthYear = this.currentMonth === 0 ? this.currentYear - 1 : this.currentYear;
        const prevMonth = this.currentMonth === 0 ? 11 : this.currentMonth - 1;
        const day = prevMonthLastDate - i;
        days.push({
          date: `${prevMonthYear}-${String(prevMonth + 1).padStart(2, '0')}-${String(day).padStart(2, '0')}`,
          dayNumber: day,
          isCurrentMonth: false
        });
      }
      
      // 현재 달의 날짜들
      for (let i = 1; i <= lastDay.getDate(); i++) {
        days.push({
          date: `${this.currentYear}-${String(this.currentMonth + 1).padStart(2, '0')}-${String(i).padStart(2, '0')}`,
          dayNumber: i,
          isCurrentMonth: true
        });
      }
      
      // 다음 달의 시작 날짜들
      const lastDayOfWeek = lastDay.getDay();
      for (let i = 1; i < 7 - lastDayOfWeek; i++) {
        const nextMonthYear = this.currentMonth === 11 ? this.currentYear + 1 : this.currentYear;
        const nextMonth = this.currentMonth === 11 ? 0 : this.currentMonth + 1;
        days.push({
          date: `${nextMonthYear}-${String(nextMonth + 1).padStart(2, '0')}-${String(i).padStart(2, '0')}`,
          dayNumber: i,
          isCurrentMonth: false
        });
      }
      
      return days;
    }
  },
  


  async mounted() {
    console.log('🚀 === ScheduleSharing 컴포넌트 마운트 시작 ===');
    
    // 로그인된 사용자 정보 로드
    await this.loadMyProfile();
    console.log('✅ 본인 프로필 로드 완료:', this.myProfile);
    
    // 팀원 데이터 로드
    await this.loadTeamMembers();
    console.log('✅ 팀원 데이터 로드 완료:', this.teamMembers);
    console.log('📊 팀원 수:', this.teamMembers.length);
    
    // 각 팀원의 ID 확인
    this.teamMembers.forEach((member, index) => {
      console.log(`👤 팀원 ${index + 1}:`, {
        id: member.id,
        name: member.name,
        role: member.role,
        department: member.department
      });
    });
    

    
    // selectedDate 초기화
    this.selectedDate = this.today;
    console.log('📅 selectedDate 초기화:', this.selectedDate);
    
    // MemberNode 요소들 확인
    this.$nextTick(() => {
      const memberNodes = document.querySelectorAll('[data-member-id]');
      console.log('🔍 찾은 MemberNode 요소들:', memberNodes.length);
      memberNodes.forEach((node, index) => {
        console.log(`📍 MemberNode ${index + 1}:`, {
          memberId: node.getAttribute('data-member-id'),
          className: node.className
        });
      });
    });
    
    // localStorage 변경 감지를 위한 이벤트 리스너 (같은 탭에서는 동작하지 않으므로 다른 방법 사용)
    // 대신 페이지 포커스 시 미니달력을 다시 렌더링하도록 설정
    this.setupPageFocusListener();
    
    console.log('🏁 === ScheduleSharing 컴포넌트 마운트 완료 ===');
  },
  
  beforeUnmount() {
    // 이벤트 리스너 정리
    if (this.pageFocusListener) {
      window.removeEventListener('focus', this.pageFocusListener);
    }
  }
};
</script>


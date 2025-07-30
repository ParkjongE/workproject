<template>
  <ScheduleSharingRoot>
    <SideBar />
    <ContentArea>
      <h1>개인 일정 공유</h1>
      <DashboardContainer>
        <GraphView ref="graphView" class="graph-view">
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
            @dblclick="expandMemberSchedule(myProfile.id)"
          >
            <MemberInfo>
              <MemberPosition>{{ myProfile.position }} | {{ myProfile.department }} {{ myProfile.team }}</MemberPosition>
              <MemberName>홍길동</MemberName>
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
                    <ScheduleSegment v-for="(schedule, idx) in myProfile.schedules" :key="`my-schedule-${idx}`" 
                      :style="getScheduleSegmentStyle(schedule)">
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
            :style="getNodePosition(index)"
            :class="{ 'selected': selectedMember === member.id }"
            @click="selectMember(member.id)"
            @dblclick="expandMemberSchedule(member.id)"
            draggable="true"
            @dragstart="dragStart($event, member.id)"
            @dragend="dragEnd"
          >
            <MemberInfo>
              <MemberPosition>{{ member.position }} | {{ member.department }} {{ member.team }}</MemberPosition>
              <MemberName>{{ member.name }}</MemberName>
            </MemberInfo>
            <ScheduleViews>
              <!-- 미니 시계 뷰 -->
              <MiniClockView @click.stop="showClockView(member.id)">
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
              <MiniCalendarView @click.stop="showCalendarView(member.id)">
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
          </MemberNode>
        </GraphView>
      </DashboardContainer>

      <!-- 확대된 일정 뷰 모달 -->
      <ExpandedScheduleModal v-if="showExpandedView">
        <ModalContent>
          <ModalHeader>
            <h2>{{ expandedMember.name }}의 일정</h2>
            <CloseBtn @click="closeExpandedView">×</CloseBtn>
          </ModalHeader>
          <ModalBody>
            <ExpandedViews>
              <!-- 왼쪽: 시계 -->
              <ExpandedClockContainer>
                <!-- expanded-date-label 제거 -->
                <ExpandedClock>
                  <ExpandedClockFace>
                    <!-- 시간 숫자 -->
                    <ExpandedHourNumber v-for="hour in 12" :key="`hour-${hour}`" :style="getExpandedHourNumberStyle(hour)">
                      {{ hour }}
                    </ExpandedHourNumber>
                    
                    <!-- 일정 세그먼트 - 선택된 날짜의 일정 표시 -->
                    <ExpandedScheduleSegment v-for="(schedule, idx) in getSelectedDateSchedules()" :key="`schedule-${idx}`" 
                      :style="getExpandedScheduleStyle(schedule)"
                      @click="selectSchedule(schedule)">
                      <ScheduleTooltip>
                        <div class="schedule-time">
                          {{ formatTime(schedule.start) }} - {{ formatTime(schedule.end) }}
                        </div>
                        <div class="schedule-title">{{ schedule.title }}</div>
                      </ScheduleTooltip>
                    </ExpandedScheduleSegment>
                  </ExpandedClockFace>
                </ExpandedClock>
                
                <!-- view-title도 제거 -->
                <TodaySchedules>
                  <TodayScheduleItem 
                    v-for="(schedule, idx) in getSelectedDateSchedules()" 
                    :key="`today-${idx}`"
                    :class="{ 'selected': selectedSchedule && selectedSchedule.id === schedule.id }"
                    @click="selectSchedule(schedule)"
                  >
                    <ScheduleColorDot :style="{ backgroundColor: schedule.color }"></ScheduleColorDot>
                    <ScheduleInfo>
                      <ScheduleTimeRange>{{ formatTime(schedule.start) }} - {{ formatTime(schedule.end) }}</ScheduleTimeRange>
                      <ScheduleTitle>{{ schedule.title }}</ScheduleTitle>
                    </ScheduleInfo>
                  </TodayScheduleItem>
                  <NoSchedules v-if="getSelectedDateSchedules().length === 0">
                    해당 날짜에 일정이 없습니다.
                  </NoSchedules>
                </TodaySchedules>
              </ExpandedClockContainer>
              <!-- 가운데: info-section(연차/출퇴근/급여명세서) -->
              <AdditionalInfoContainer>
                <!-- info-section(연차/출퇴근/급여명세서) -->
                <InfoSection class="leave-info-compact">
                  <SectionHeaderWithPrivacy>
                    <SectionTitle>연차 정보</SectionTitle>
                    <PrivacyStatus class="public">
                      <span>공개</span>
                    </PrivacyStatus>
                  </SectionHeaderWithPrivacy>
                  <LeaveItem>
                    <LeaveLabel>부여된 연차:</LeaveLabel>
                    <LeaveValue>{{ expandedMember.totalLeave || 15 }}일</LeaveValue>
                  </LeaveItem>
                  <LeaveItem>
                    <LeaveLabel>사용한 연차:</LeaveLabel>
                    <LeaveValue>{{ expandedMember.usedLeave || 5 }}일</LeaveValue>
                  </LeaveItem>
                  <LeaveItem>
                    <LeaveLabel>남은 연차:</LeaveLabel>
                    <LeaveValue>{{ (expandedMember.totalLeave || 15) - (expandedMember.usedLeave || 5) }}일</LeaveValue>
                  </LeaveItem>
                </InfoSection>
                
                <!-- 출퇴근 정보 -->
                <InfoSection class="attendance-info">
                  <SectionHeaderWithPrivacy>
                    <h4>오늘의 출퇴근</h4>
                    <PrivacyStatus class="public">
                      <span>공개</span>
                    </PrivacyStatus>
                  </SectionHeaderWithPrivacy>
                  <AttendanceStatus>
                    <StatusItem>
                      <StatusLabel>출근 시간:</StatusLabel>
                      <StatusValue>{{ expandedMember.todayAttendance?.checkInTime || '09:05' }}</StatusValue>
                    </StatusItem>
                    <StatusItem>
                      <StatusLabel>퇴근 시간:</StatusLabel>
                      <StatusValue>{{ expandedMember.todayAttendance?.checkOutTime || '18:10' }}</StatusValue>
                    </StatusItem>
                    <StatusItem>
                      <StatusLabel>근무 시간:</StatusLabel>
                      <StatusValue>{{ expandedMember.workingHours || '9시간 5분' }}</StatusValue>
                    </StatusItem>
                  </AttendanceStatus>
                </InfoSection>
                
                <!-- 급여명세서 -->
                <InfoSection class="salary-info">
                  <div class="salary-header">
                    <SectionHeaderWithPrivacy>
                      <h4>급여명세서</h4>
                      <PrivacyStatus class="private">
                        <span>비공개</span>
                      </PrivacyStatus>
                    </SectionHeaderWithPrivacy>
                  </div>
                  <SalaryYears>
                    <SalaryYear class="active">
                      <YearHeader>
                        <span>2023년</span>
                        <YearActions>
                          <DownloadAllBtn title="모두 다운로드">
                            <DownloadAllIcon>↓↓</DownloadAllIcon>
                          </DownloadAllBtn>
                        </YearActions>
                      </YearHeader>
                      <SalaryMonths>
                        <SalaryMonth>
                          <span>6월</span>
                          <DownloadIcon>↓</DownloadIcon>
                        </SalaryMonth>
                        <SalaryMonth>
                          <span>5월</span>
                          <DownloadIcon>↓</DownloadIcon>
                        </SalaryMonth>
                      </SalaryMonths>
                    </SalaryYear>
                  </SalaryYears>
                </InfoSection>
              </AdditionalInfoContainer>
              <!-- 오른쪽: 달력 -->
              <ExpandedCalendarContainer>
                <CalendarHeader>
                  <button @click="prevMonth">&lt;</button>
                  <h3>{{ currentYear }}년 {{ currentMonth + 1 }}월</h3>
                  <button @click="nextMonth">&gt;</button>
                </CalendarHeader>
                <CalendarGrid>
                  <WeekdayHeader v-for="day in weekdays" :key="day">{{ day }}</WeekdayHeader>
                  <CalendarDay 
                    v-for="day in calendarDays" 
                    :key="day.date"
                    :class="{ 
                      'other-month': !day.isCurrentMonth,
                      'selected-date': getCurrentSelectedDate() === day.date,
                      'project-active': isProjectActive(day.date)
                    }"
                    @click="selectDate(day.date)"
                  >
                    <DayNumber>{{ day.dayNumber }}</DayNumber>
                    
                    <!-- 프로젝트 선 표시 -->
                    <ProjectLine v-if="isProjectActive(day.date)"
                         :style="{ backgroundColor: getProjectForDate(day.date)?.color || '#B4D3F5' }">
                    </ProjectLine>
                    
                    <DaySchedules>
                      <DayScheduleItem 
                        v-for="(schedule, idx) in getDaySchedules(day.date)" 
                        :key="idx"
                        :style="{ backgroundColor: schedule.color }"
                        @click.stop="selectSchedule(schedule)"
                      >
                        {{ schedule.title }}
                      </DayScheduleItem>
                    </DaySchedules>
                  </CalendarDay>
                </CalendarGrid>
                
                <!-- 주간 일정 표시 -->
                <WeekView>
                  <WeekHeader>
                    <WeekDayHeader v-for="day in weekdays" :key="`week-${day}`">{{ day }}</WeekDayHeader>
                  </WeekHeader>
                  <WeekBody>
                    <WeekRow v-for="(weekSchedules, weekIdx) in getWeekSchedules()" :key="`week-${weekIdx}`">
                      <WeekDay v-for="(daySchedules, dayIdx) in weekSchedules" :key="`day-${dayIdx}`">
                        <WeekScheduleLine 
                          v-for="(schedule, scheduleIdx) in daySchedules" 
                          :key="`schedule-${scheduleIdx}`"
                          :style="{ 
                            backgroundColor: schedule.color,
                            top: `${(scheduleIdx * 4) + 2}px`
                          }"
                          @click="selectSchedule(schedule)"
                        ></WeekScheduleLine>
                      </WeekDay>
                    </WeekRow>
                  </WeekBody>
                </WeekView>
              </ExpandedCalendarContainer>
            </ExpandedViews>
            
            <ScheduleDetails v-if="selectedSchedule">
              <h3>{{ selectedSchedule.title }}</h3>
              <p>{{ formatDate(selectedSchedule.start) }} {{ formatTime(selectedSchedule.start) }} - {{ formatTime(selectedSchedule.end) }}</p>
            </ScheduleDetails>
            
            <!-- 선택된 날짜 정보 표시 -->
            <SelectedDateInfo v-if="selectedDate">
              <h3>{{ formatSelectedDate(selectedDate) }} 일정</h3>
              <DateSchedules>
                <DateScheduleItem 
                  v-for="(schedule, idx) in getDaySchedules(selectedDate)" 
                  :key="idx"
                  :class="{ 'selected': selectedSchedule && selectedSchedule.id === schedule.id }"
                  @click="selectSchedule(schedule)"
                >
                  <ScheduleTime>{{ formatTime(schedule.start) }} - {{ formatTime(schedule.end) }}</ScheduleTime>
                  <ScheduleTitle>{{ schedule.title }}</ScheduleTitle>
                </DateScheduleItem>
                <NoSchedules v-if="getDaySchedules(selectedDate).length === 0">
                  해당 날짜에 일정이 없습니다.
                </NoSchedules>
              </DateSchedules>
            </SelectedDateInfo>
            
            <ScheduleActions>
              <AddBtn @click="addSchedule">일정 추가</AddBtn>
              <EditBtn @click="editSchedule" :disabled="!selectedSchedule">일정 수정</EditBtn>
            </ScheduleActions>
            
            <!-- 모달 하단 프로젝트 정보 -->
            <ModalBottomProjects>
              <h4>진행 중인 프로젝트</h4>
              <ProjectList>
                <ProjectInfo v-for="project in projects" :key="project.id">
                  <ProjectLinePreview :style="{ backgroundColor: project.color }"></ProjectLinePreview>
                  <ProjectDetails>
                    <ProjectName>{{ project.title }}</ProjectName>
                    <ProjectPeriod>{{ formatProjectDate(project.startDate) }} - {{ formatProjectDate(project.endDate) }}</ProjectPeriod>
                  </ProjectDetails>
                </ProjectInfo>
              </ProjectList>
            </ModalBottomProjects>
          </ModalBody>
        </ModalContent>
      </ExpandedScheduleModal>
    </ContentArea>
  </ScheduleSharingRoot>
</template>

<script>
import SideBar from '@/components/SideBar.vue';
import {
  ScheduleSharingRoot, ContentArea, DashboardContainer, GraphView, DepartmentBoundary, DepartmentLabel, MemberNode, MyNode, MemberInfo, MemberPosition, MemberName, ScheduleViews, MiniClockView, MiniClock, ClockFace, HourNumber, ScheduleSegment, MiniCalendarView, MiniCalendar, MiniCalendarGrid, MiniDayHeader, MiniDayCell, MiniDayNumber, MiniDayEvent, MiniProjectLine,
  ExpandedScheduleModal, ModalContent, ModalHeader, CloseBtn, ModalBody, ExpandedViews, ExpandedClockContainer, ExpandedClock, ExpandedClockFace, ExpandedHourNumber, ExpandedScheduleSegment, ScheduleTooltip, TodaySchedules, TodayScheduleItem, ScheduleColorDot, ScheduleInfo, ScheduleTimeRange, ScheduleTitle, NoSchedules, AdditionalInfoContainer, InfoSection, SectionHeaderWithPrivacy, SectionTitle, PrivacyStatus, LeaveItem, LeaveLabel, LeaveValue, AttendanceStatus, StatusItem, StatusLabel, StatusValue,
  ExpandedCalendarContainer, CalendarHeader, CalendarGrid, WeekdayHeader, CalendarDay, DayNumber, ProjectLine, DaySchedules, DayScheduleItem, WeekView, WeekHeader, WeekDayHeader, WeekBody, WeekRow, WeekDay, WeekScheduleLine,
  SalaryYears, SalaryYear, YearHeader, YearActions, DownloadAllBtn, DownloadAllIcon, SalaryMonths, SalaryMonth, DownloadIcon,
  ScheduleDetails, SelectedDateInfo, DateSchedules, DateScheduleItem, ScheduleTime, ScheduleActions, AddBtn, EditBtn,
  ModalBottomProjects, ProjectList, ProjectInfo, ProjectLinePreview, ProjectDetails, ProjectName, ProjectPeriod
} from './ScheduleSharing.styled.js';

export default {
  name: 'ScheduleSharing',
  components: {
    SideBar,
    ScheduleSharingRoot, ContentArea, DashboardContainer, GraphView, DepartmentBoundary, DepartmentLabel, MemberNode, MyNode, MemberInfo, MemberPosition, MemberName, ScheduleViews, MiniClockView, MiniClock, ClockFace, HourNumber, ScheduleSegment, MiniCalendarView, MiniCalendar, MiniCalendarGrid, MiniDayHeader, MiniDayCell, MiniDayNumber, MiniDayEvent, MiniProjectLine,
    ExpandedScheduleModal, ModalContent, ModalHeader, CloseBtn, ModalBody, ExpandedViews, ExpandedClockContainer, ExpandedClock, ExpandedClockFace, ExpandedHourNumber, ExpandedScheduleSegment, ScheduleTooltip, TodaySchedules, TodayScheduleItem, ScheduleColorDot, ScheduleInfo, ScheduleTimeRange, ScheduleTitle, NoSchedules, AdditionalInfoContainer, InfoSection, SectionHeaderWithPrivacy, SectionTitle, PrivacyStatus, LeaveItem, LeaveLabel, LeaveValue, AttendanceStatus, StatusItem, StatusLabel, StatusValue,
    ExpandedCalendarContainer, CalendarHeader, CalendarGrid, WeekdayHeader, CalendarDay, DayNumber, ProjectLine, DaySchedules, DayScheduleItem, WeekView, WeekHeader, WeekDayHeader, WeekBody, WeekRow, WeekDay, WeekScheduleLine,
    SalaryYears, SalaryYear, YearHeader, YearActions, DownloadAllBtn, DownloadAllIcon, SalaryMonths, SalaryMonth, DownloadIcon,
    ScheduleDetails, SelectedDateInfo, DateSchedules, DateScheduleItem, ScheduleTime, ScheduleActions, AddBtn, EditBtn,
    ModalBottomProjects, ProjectList, ProjectInfo, ProjectLinePreview, ProjectDetails, ProjectName, ProjectPeriod
  },
  data() {
    return {
      // 일정 유형별 색상 정의
      scheduleColors: {
        meeting: '#FF6B6B', // 회의용 빨간색
        work: '#4ECDC4',    // 업무용 청록색
        other: '#45B7D1'    // 기타용 파란색
      },
      // 본인 프로필 추가
      myProfile: {
        id: 0,
        name: '박종희',
        position: '사원',
        department: '개발부',
        team: '프론트엔드팀',
        schedules: [
          { id: 101, title: '팀 미팅', start: '2023-06-15T10:00', end: '2023-06-15T11:30', type: 'meeting', color: '#FF6B6B' },
          { id: 102, title: '코드 리뷰', start: '2023-06-15T14:00', end: '2023-06-15T15:30', type: 'work', color: '#4ECDC4' },
          { id: 103, title: '프로젝트 기획', start: '2023-06-15T16:00', end: '2023-06-15T17:30', type: 'other', color: '#45B7D1' }
        ]
      },
      departments: [
        // 부서 네모 크기 확대 (width와 height를 350에서 450으로 증가)
        // 부서 간 간격도 넓힘 (x 값 조정)
        { name: '개발부', x: 150, y: 300, width: 450, height: 450 },
        { name: '디자인부', x: 650, y: 300, width: 450, height: 450 },
        { name: '마케팅부', x: 1150, y: 300, width: 450, height: 450 }
      ],
      teamMembers: [
        // 개발부
        { id: 1, name: '김부장', position: '부장', department: '개발부', team: '팀장', schedules: [
          { id: 1, title: '팀 미팅', start: '2023-06-15T10:00', end: '2023-06-15T11:30', type: 'meeting', color: '#FF6B6B' },
          { id: 2, title: '고객 미팅', start: '2023-06-15T14:00', end: '2023-06-15T15:30', type: 'meeting', color: '#FF6B6B' },
          { id: 3, title: '프로젝트 기획', start: '2023-06-16T09:00', end: '2023-06-16T12:00', type: 'work', color: '#4ECDC4' },
          { id: 12, title: '월간 보고서', start: '2023-06-20T13:00', end: '2023-06-20T15:00', type: 'work', color: '#4ECDC4' },
          { id: 13, title: '신규 프로젝트 검토', start: '2023-06-22T10:00', end: '2023-06-22T12:00', type: 'other', color: '#45B7D1' }
        ]},
        { id: 2, name: '이차장', position: '차장', department: '개발부', team: '백엔드팀', schedules: [
          { id: 4, title: '팀 미팅', start: '2023-06-15T10:00', end: '2023-06-15T11:30', type: 'meeting', color: '#FF6B6B' },
          { id: 5, title: '보고서 작성', start: '2023-06-15T13:00', end: '2023-06-15T17:00', type: 'work', color: '#4ECDC4' },
          { id: 14, title: '클라이언트 미팅', start: '2023-06-19T14:00', end: '2023-06-19T16:00', type: 'meeting', color: '#FF6B6B' },
          { id: 15, title: '예산 검토', start: '2023-06-21T09:00', end: '2023-06-21T11:00', type: 'other', color: '#45B7D1' }
        ]},
        { id: 3, name: '박과장', position: '과장', department: '개발부', team: '프론트엔드팀', schedules: [
          { id: 6, title: '팀 미팅', start: '2023-06-15T10:00', end: '2023-06-15T11:30', type: 'meeting', color: '#FF6B6B' },
          { id: 7, title: '데이터 분석', start: '2023-06-16T13:00', end: '2023-06-16T18:00', type: 'work', color: '#4ECDC4' },
          { id: 16, title: '마케팅 전략 회의', start: '2023-06-18T10:00', end: '2023-06-18T12:00', type: 'meeting', color: '#FF6B6B' },
          { id: 17, title: '분기별 성과 분석', start: '2023-06-23T13:00', end: '2023-06-23T17:00', type: 'work', color: '#4ECDC4' }
        ]},
        { id: 4, name: '최대리', position: '대리', department: '개발부', team: '백엔드팀', schedules: [
          { id: 8, title: '팀 미팅', start: '2023-06-15T10:00', end: '2023-06-15T11:30', type: 'meeting', color: '#FF6B6B' },
          { id: 9, title: '디자인 작업', start: '2023-06-17T09:00', end: '2023-06-17T18:00', type: 'work', color: '#4ECDC4' },
          { id: 18, title: 'UI 디자인 검토', start: '2023-06-19T10:00', end: '2023-06-19T12:00', type: 'other', color: '#45B7D1' },
          { id: 19, title: '사용자 테스트', start: '2023-06-22T14:00', end: '2023-06-22T16:00', type: 'other', color: '#45B7D1' }
        ]},
        
        // 디자인부
        { id: 5, name: '정사원', position: '사원', department: '디자인부', team: 'UI팀', schedules: [
          { id: 10, title: '팀 미팅', start: '2023-06-15T10:00', end: '2023-06-15T11:30', type: 'meeting', color: '#FF6B6B' },
          { id: 11, title: '코딩 작업', start: '2023-06-15T13:00', end: '2023-06-15T18:00', type: 'work', color: '#4ECDC4' },
          { id: 20, title: '버그 수정', start: '2023-06-20T09:00', end: '2023-06-20T12:00', type: 'work', color: '#4ECDC4' },
          { id: 21, title: '코드 리뷰', start: '2023-06-21T14:00', end: '2023-06-21T16:00', type: 'other', color: '#45B7D1' }
        ]},
        { id: 6, name: '한부장', position: '부장', department: '디자인부', team: '팀장', schedules: [
          { id: 22, title: '디자인 리뷰', start: '2023-06-15T09:00', end: '2023-06-15T11:00', type: 'meeting', color: '#FF6B6B' },
          { id: 23, title: '임원 회의', start: '2023-06-15T14:00', end: '2023-06-15T16:00', type: 'meeting', color: '#FF6B6B' },
          { id: 24, title: '신규 프로젝트 기획', start: '2023-06-16T10:00', end: '2023-06-16T12:00', type: 'work', color: '#4ECDC4' }
        ]},
        { id: 7, name: '윤과장', position: '과장', department: '디자인부', team: 'UX팀', schedules: [
          { id: 25, title: '사용자 연구', start: '2023-06-15T11:00', end: '2023-06-15T13:00', type: 'work', color: '#4ECDC4' },
          { id: 26, title: '프로토타입 개발', start: '2023-06-15T15:00', end: '2023-06-15T17:00', type: 'work', color: '#4ECDC4' },
          { id: 27, title: '사용성 테스트', start: '2023-06-17T10:00', end: '2023-06-17T12:00', type: 'other', color: '#45B7D1' }
        ]},
        
        // 마케팅부
        { id: 8, name: '조차장', position: '차장', department: '마케팅부', team: '디지털마케팅팀', schedules: [
          { id: 28, title: '광고 전략 회의', start: '2023-06-15T09:30', end: '2023-06-15T11:00', type: 'meeting', color: '#FF6B6B' },
          { id: 29, title: '소셜 미디어 분석', start: '2023-06-15T13:30', end: '2023-06-15T15:30', type: 'work', color: '#4ECDC4' },
          { id: 30, title: '캠페인 기획', start: '2023-06-16T14:00', end: '2023-06-16T16:00', type: 'work', color: '#4ECDC4' }
        ]},
        { id: 9, name: '임대리', position: '대리', department: '마케팅부', team: '브랜드마케팅팀', schedules: [
          { id: 31, title: '브랜드 전략 회의', start: '2023-06-15T11:00', end: '2023-06-15T12:30', type: 'meeting', color: '#FF6B6B' },
          { id: 32, title: '콘텐츠 기획', start: '2023-06-15T15:00', end: '2023-06-15T17:00', type: 'work', color: '#4ECDC4' },
          { id: 33, title: '마케팅 자료 제작', start: '2023-06-17T09:00', end: '2023-06-17T12:00', type: 'work', color: '#4ECDC4' }
        ]},
        { id: 10, name: '오사원', position: '사원', department: '마케팅부', team: '디지털마케팅팀', schedules: [
          { id: 34, title: '데이터 분석', start: '2023-06-15T09:00', end: '2023-06-15T12:00', type: 'work', color: '#4ECDC4' },
          { id: 35, title: '광고 소재 제작', start: '2023-06-15T14:00', end: '2023-06-15T16:00', type: 'work', color: '#4ECDC4' },
          { id: 36, title: '성과 보고서 작성', start: '2023-06-16T13:00', end: '2023-06-16T15:00', type: 'other', color: '#45B7D1' }
        ]}
      ],
      selectedMember: null,
      showExpandedView: false,
      expandedMember: null,
      activeTab: 'clock',
      selectedSchedule: null,
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
      ]
    };
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
  methods: {
    // 일정 유형에 따라 색상 자동 할당
    getScheduleColorByType(type) {
      switch(type) {
        case 'meeting':
          return this.scheduleColors.meeting;
        case 'work':
          return this.scheduleColors.work;
        case 'other':
          return this.scheduleColors.other;
        default:
          return this.scheduleColors.other;
      }
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
      // 저장된 위치가 있으면 사용
      if (this.nodePositions[this.teamMembers[index].id]) {
        return this.nodePositions[this.teamMembers[index].id];
      }
      
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
    selectMember(memberId) {
      this.selectedMember = memberId;
    },
    expandMemberSchedule(memberId) {
      this.expandedMember = this.teamMembers.find(m => m.id === memberId);
      this.showExpandedView = true;
      this.activeTab = 'clock';
      // 모달 열 때 오늘 날짜로 초기화
      this.memberSelectedDates[memberId] = this.today;
    },


    closeExpandedView() {
      this.showExpandedView = false;
      this.expandedMember = null;
      this.selectedSchedule = null;
      // 모달 닫을 때 선택된 날짜 초기화
      if (this.expandedMember) {
        this.memberSelectedDates[this.expandedMember.id] = null;
      }
    },
    showClockView(memberId) {
      this.expandMemberSchedule(memberId);
      this.activeTab = 'clock';
    },
    showCalendarView(memberId) {
      this.expandMemberSchedule(memberId);
      this.activeTab = 'calendar';
    },
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
      // newSchedule.color = this.getScheduleColorByType(newSchedule.type);
      // this.expandedMember.schedules.push(newSchedule);
    },
    editSchedule() {
      if (!this.selectedSchedule) return;
      // 실제 구현에서는 일정 수정 모달 표시
      alert('일정 수정 기능은 아직 구현되지 않았습니다.');
      // 일정 수정 시 유형에 따라 색상 자동 설정 예시
      // this.selectedSchedule.color = this.getScheduleColorByType(this.selectedSchedule.type);
    },
    dragStart(event, memberId) {
      this.draggedMember = memberId;
      event.dataTransfer.setData('text/plain', memberId);
      event.dataTransfer.effectAllowed = 'move';
    },
    dragEnd() {
      this.draggedMember = null;
    },
    updateNodePosition(memberId, left, top) {
      this.nodePositions[memberId] = { left: `${left}px`, top: `${top}px` };
    },
    getHourMarkStyle(hour) {
      const angle = (hour / 12) * 2 * Math.PI - Math.PI / 2;
      const x = Math.cos(angle) * 18;
      const y = Math.sin(angle) * 18;
      return {
        left: `${x + 20}px`,
        top: `${y + 20}px`
      };
    },
    // 오늘의 일정만 가져오기
    getTodaySchedules(member) {
      // 예시 데이터를 위해 하드코딩된 날짜 사용
      const today = this.today;
      const schedules = member.schedules.filter(s => s.start.startsWith(today));
      
      // 디버깅용 로그
      console.log(`${member.name}의 오늘 일정:`, schedules);
      schedules.forEach(s => {
        const start = new Date(s.start);
        const end = new Date(s.end);
        console.log(`${s.title}: ${start.getHours()}:${start.getMinutes()} ~ ${end.getHours()}:${end.getMinutes()}`);
      });
      
      return schedules;
    },

    getScheduleSegmentStyle(schedule) {
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
        console.log(`미니 시계 - ${schedule.title}: ${startHour}:${startMinute}(${startHour12}) ~ ${endHour}:${endMinute}(${endHour12})`);
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
    
    getMiniCalendarDaySchedule(member, dayIndex) {
      const currentDate = new Date(this.currentYear, this.currentMonth, 1);
      const firstDayOfMonth = currentDate.getDay();
      const adjustedIndex = dayIndex - firstDayOfMonth;
      
      if (adjustedIndex < 0 || adjustedIndex > 30) return null;
      
      const day = new Date(this.currentYear, this.currentMonth, adjustedIndex + 1);
      const dateString = day.toISOString().split('T')[0];
      
      const scheduleForDay = member.schedules.find(s => s.start.startsWith(dateString));
      return scheduleForDay || null;
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
        fontWeight: 'bold'
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
        fontWeight: 'bold'
      };
    },
    
    getMiniCalendarDayNumber(index) {
      const currentDate = new Date(this.currentYear, this.currentMonth, 1);
      const firstDayOfMonth = currentDate.getDay();
      const adjustedIndex = index - firstDayOfMonth;
      
      if (adjustedIndex < 0) {
        // 이전 달의 날짜
        const prevMonthLastDate = new Date(this.currentYear, this.currentMonth, 0).getDate();
        return prevMonthLastDate + adjustedIndex + 1;
      } else if (adjustedIndex >= new Date(this.currentYear, this.currentMonth + 1, 0).getDate()) {
        // 다음 달의 날짜
        return adjustedIndex - new Date(this.currentYear, this.currentMonth + 1, 0).getDate() + 1;
      } else {
        // 현재 달의 날짜
        return adjustedIndex + 1;
      }
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
    
    // 미니 달력의 날짜 인덱스로부터 실제 날짜 계산
    getMiniCalendarDate(index) {
      const currentDate = new Date(this.currentYear, this.currentMonth, 1);
      const firstDayOfMonth = currentDate.getDay();
      const adjustedIndex = index - firstDayOfMonth;
      
      if (adjustedIndex < 0) {
        // 이전 달의 날짜
        const prevMonthYear = this.currentMonth === 0 ? this.currentYear - 1 : this.currentYear;
        const prevMonth = this.currentMonth === 0 ? 11 : this.currentMonth - 1;
        const prevMonthLastDate = new Date(this.currentYear, this.currentMonth, 0).getDate();
        const day = prevMonthLastDate + adjustedIndex + 1;
        return `${prevMonthYear}-${String(prevMonth + 1).padStart(2, '0')}-${String(day).padStart(2, '0')}`;
      } else if (adjustedIndex >= new Date(this.currentYear, this.currentMonth + 1, 0).getDate()) {
        // 다음 달의 날짜
        const nextMonthYear = this.currentMonth === 11 ? this.currentYear + 1 : this.currentYear;
        const nextMonth = this.currentMonth === 11 ? 0 : this.currentMonth + 1;
        const day = adjustedIndex - new Date(this.currentYear, this.currentMonth + 1, 0).getDate() + 1;
        return `${nextMonthYear}-${String(nextMonth + 1).padStart(2, '0')}-${String(day).padStart(2, '0')}`;
      } else {
        // 현재 달의 날짜
        const day = adjustedIndex + 1;
        return `${this.currentYear}-${String(this.currentMonth + 1).padStart(2, '0')}-${String(day).padStart(2, '0')}`;
      }
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

    // 선택된 날짜의 일정만 가져오는 메서드 추가
    getSelectedDateSchedules() {
      if (!this.expandedMember) return [];
      
      const memberId = this.expandedMember.id;
      const targetDate = this.memberSelectedDates[memberId] || this.today;
      
      return this.expandedMember.schedules.filter(s => s.start.startsWith(targetDate));
    },
    
    // 현재 선택된 날짜 가져오기
    getCurrentSelectedDate() {
      if (!this.expandedMember) return this.today;
      return this.memberSelectedDates[this.expandedMember.id] || this.today;
    }
  },
  mounted() {
    // 드래그 앤 드롭 이벤트 리스너 설정
    this.$nextTick(() => {
      // DOM에서 직접 요소를 찾아서 이벤트 리스너 설정
      const graphViewElement = document.querySelector('.graph-view');
      
      if (graphViewElement) {
        graphViewElement.addEventListener('dragover', (e) => {
          e.preventDefault();
        });
        
        graphViewElement.addEventListener('drop', (e) => {
          e.preventDefault();
          if (this.draggedMember) {
            // 드롭된 위치로 노드 이동
            const rect = graphViewElement.getBoundingClientRect();
            const left = e.clientX - rect.left;
            const top = e.clientY - rect.top;
            this.updateNodePosition(this.draggedMember, left, top);
          }
        });
      } else {
        console.warn('GraphView 요소를 찾을 수 없습니다.');
      }
    });
  }
};
</script>


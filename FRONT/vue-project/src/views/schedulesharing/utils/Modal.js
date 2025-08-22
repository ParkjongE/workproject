/**
 * 모달 관련 유틸리티 함수들
 */

/**
 * 멤버 일정 모달 열기
 * @param {string|number} memberId - 멤버 ID
 * @param {Object} context - Vue 컴포넌트 컨텍스트
 */
export function expandMemberSchedule(memberId, context) {
  console.log('🔍 expandMemberSchedule 호출됨, memberId:', memberId);
  
  // 본인 노드(id: 0)인 경우 MySchedule 모달 표시
  if (memberId === 0) {
    console.log('👤 본인 노드 - MySchedule 모달 열기');
    context.showMyScheduleModal = true;
    return;
  }
  
  // 다른 팀원들인 경우 기존 확대 모달 표시
  context.expandedMember = context.teamMembers.find(m => m.id === memberId);
  console.log('👥 팀원 노드 - 확대 모달 열기:', context.expandedMember?.name);
  
  context.showExpandedView = true;
  context.activeTab = 'clock'; // 기본 탭을 시계로 설정
  
  // 모달 열 때 오늘 날짜로 초기화
  context.memberSelectedDates[memberId] = context.today;
  
  console.log('✅ 모달 열기 완료');
}

/**
 * 본인 일정 모달 닫기
 * @param {Object} context - Vue 컴포넌트 컨텍스트
 */
export function closeMyScheduleModal(context) {
  context.showMyScheduleModal = false;
}

/**
 * 확대된 일정 뷰 모달 닫기
 * @param {Object} context - Vue 컴포넌트 컨텍스트
 */
export function closeExpandedView(context) {
  context.showExpandedView = false;
  context.expandedMember = null;
  context.selectedSchedule = null;
  // 모달 닫을 때 선택된 날짜 초기화
  if (context.expandedMember) {
    context.memberSelectedDates[context.expandedMember.id] = null;
  }
}

/**
 * 시계 뷰 모달 열기
 * @param {string|number} memberId - 멤버 ID
 * @param {Object} context - Vue 컴포넌트 컨텍스트
 */
export function showClockView(memberId, context) {
  expandMemberSchedule(memberId, context);
  context.activeTab = 'clock';
}

/**
 * 달력 뷰 모달 열기
 * @param {string|number} memberId - 멤버 ID
 * @param {Object} context - Vue 컴포넌트 컨텍스트
 */
export function showCalendarView(memberId, context) {
  expandMemberSchedule(memberId, context);
  context.activeTab = 'calendar';
}

/**
 * 빠른 미팅 모달 열기
 * @param {Object} timeSlot - 선택된 시간 슬롯
 * @param {Object} context - Vue 컴포넌트 컨텍스트
 */
export function openQuickMeetingModal(timeSlot, context) {
  context.selectedTimeSlot = timeSlot;
  context.quickMeeting = {
    title: '',
    type: 'meeting'
  };
  context.showQuickModal = true;
}

/**
 * 빠른 미팅 모달 닫기
 * @param {Object} context - Vue 컴포넌트 컨텍스트
 */
export function closeQuickModal(context) {
  context.showQuickModal = false;
  context.selectedTimeSlot = null;
}

/**
 * 일정 수정 모달 열기
 * @param {Object} schedule - 수정할 일정
 * @param {Object} context - Vue 컴포넌트 컨텍스트
 */
export function openEditScheduleModal(schedule, context) {
  context.editingSchedule = {
    id: schedule.id,
    title: schedule.title.replace('[요청] ', ''), // [요청] 접두사 제거
    type: schedule.type,
    hour: new Date(schedule.start).getHours()
  };
  context.showEditModal = true;
}

/**
 * 일정 수정 모달 닫기
 * @param {Object} context - Vue 컴포넌트 컨텍스트
 */
export function closeEditModal(context) {
  context.showEditModal = false;
  context.editingSchedule = {
    id: null,
    title: '',
    type: 'meeting',
    hour: 0
  };
}

/**
 * 드래그 스케줄 모달 열기
 * @param {Object} context - Vue 컴포넌트 컨텍스트
 */
export function openDragScheduleModal(context) {
  context.selectedTimeSlot = {
    hour: context.selectedTimeRange.start,
    endHour: context.selectedTimeRange.end
  };
  context.quickMeeting = {
    title: '',
    type: 'meeting'
  };
  context.showQuickModal = true;
}

/**
 * 사용자 스케줄 모달 열기
 * @param {Object} user - 사용자 정보
 * @param {Object} context - Vue 컴포넌트 컨텍스트
 */
export function openUserScheduleModal(user, context) {
  console.log('사용자 스케줄 모달 열기:', user);
  if (context.scheduleStore) {
    context.scheduleStore.openUserScheduleModal(user);
  } else {
    console.warn('scheduleStore가 정의되지 않았습니다.');
  }
}

/**
 * 모든 모달 상태 초기화
 * @param {Object} context - Vue 컴포넌트 컨텍스트
 */
export function resetAllModals(context) {
  context.showMyScheduleModal = false;
  context.showExpandedView = false;
  context.showQuickModal = false;
  context.showEditModal = false;
  context.expandedMember = null;
  context.selectedSchedule = null;
  context.selectedTimeSlot = null;
  context.editingSchedule = {
    id: null,
    title: '',
    type: 'meeting',
    hour: 0
  };
}

/**
 * 모달 상태 확인
 * @param {Object} context - Vue 컴포넌트 컨텍스트
 * @returns {Object} 현재 모달 상태
 */
export function getModalStatus(context) {
  return {
    showMyScheduleModal: context.showMyScheduleModal,
    showExpandedView: context.showExpandedView,
    showQuickModal: context.showQuickModal,
    showEditModal: context.showEditModal,
    hasOpenModal: context.showMyScheduleModal || context.showExpandedView || context.showQuickModal || context.showEditModal
  };
}

/**
 * 모달 열기 전 유효성 검사
 * @param {string} modalType - 모달 타입
 * @param {Object} context - Vue 컴포넌트 컨텍스트
 * @param {Object} data - 모달에 필요한 데이터
 * @returns {boolean} 모달을 열 수 있는지 여부
 */
export function validateModalOpen(modalType, context, data = {}) {
  switch (modalType) {
    case 'mySchedule':
      return context.myProfile && context.myProfile.id !== undefined;
    
    case 'expandedView':
      return data.memberId && context.teamMembers && context.teamMembers.length > 0;
    
    case 'quickMeeting':
      return data.timeSlot && data.timeSlot.hour !== undefined;
    
    case 'editSchedule':
      return data.schedule && data.schedule.id;
    
    default:
      return true;
  }
}

/**
 * 모달 닫기 전 정리 작업
 * @param {string} modalType - 모달 타입
 * @param {Object} context - Vue 컴포넌트 컨텍스트
 */
export function cleanupBeforeClose(modalType, context) {
  switch (modalType) {
    case 'expandedView':
      if (context.expandedMember) {
        context.memberSelectedDates[context.expandedMember.id] = null;
      }
      break;
    
    case 'quickMeeting':
      context.selectedTimeSlot = null;
      context.quickMeeting = { title: '', type: 'meeting' };
      break;
    
    case 'editSchedule':
      context.editingSchedule = { id: null, title: '', type: 'meeting', hour: 0 };
      break;
  }
}

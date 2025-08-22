/**
 * 일정 수정 관련 유틸리티 함수들
 */

/**
 * 일정 수정 모달 열기
 * @param {Object} schedule - 수정할 일정 객체
 * @param {Object} context - Vue 컴포넌트 컨텍스트
 */
export function editScheduleRequest(schedule, context) {
  // 본인이 추가한 일정인지 확인
  if (!isMySchedule(schedule, context)) {
    alert('본인이 추가한 일정만 수정할 수 있습니다.');
    return;
  }
  
  context.editingSchedule = {
    id: schedule.id,
    title: schedule.title.replace('[요청] ', ''), // [요청] 접두사 제거
    type: schedule.type,
    hour: new Date(schedule.start).getHours()
  };
  context.showEditModal = true;
}



/**
 * 일정 수정 완료
 * @param {Object} context - Vue 컴포넌트 컨텍스트
 */
export function updateScheduleRequest(context) {
  try {
    const scheduleIndex = context.expandedMember.schedules.findIndex(s => s.id === context.editingSchedule.id);
    if (scheduleIndex !== -1) {
      // 기존 일정 수정
      context.expandedMember.schedules[scheduleIndex].title = `[요청] ${context.editingSchedule.title}`;
      context.expandedMember.schedules[scheduleIndex].type = context.editingSchedule.type;
      context.expandedMember.schedules[scheduleIndex].color = '#4CAF50'; // 요청 일정은 항상 초록색
      
              alert('일정이 수정되었습니다!');
        // closeEditModal은 Modal.js에서 처리됨
    }
  } catch (error) {
    console.error('일정 수정 실패:', error);
    alert('일정 수정에 실패했습니다.');
  }
}

/**
 * 일정 삭제
 * @param {Object} schedule - 삭제할 일정 객체
 * @param {Object} context - Vue 컴포넌트 컨텍스트
 */
export function deleteScheduleRequest(schedule, context) {
  console.log('🗑️ 삭제 버튼 클릭됨:', schedule.title);
  
  // 본인이 추가한 일정인지 다시 한 번 확인
  if (!isMySchedule(schedule, context)) {
    alert('본인이 추가한 일정만 삭제할 수 있습니다.');
    return;
  }
  
  if (confirm(`"${schedule.title}" 일정을 삭제하시겠습니까?`)) {
    try {
      const scheduleIndex = context.expandedMember.schedules.findIndex(s => s.id === schedule.id);
      if (scheduleIndex !== -1) {
        // 추가 검증: 삭제하려는 일정이 본인 것인지 확인
        const targetSchedule = context.expandedMember.schedules[scheduleIndex];
        if (targetSchedule.requestedBy === context.myProfile.name) {
          context.expandedMember.schedules.splice(scheduleIndex, 1);
          alert('일정이 삭제되었습니다!');
        } else {
          alert('권한이 없습니다.');
        }
      } else {
        alert('삭제할 일정을 찾을 수 없습니다.');
      }
    } catch (error) {
      console.error('일정 삭제 실패:', error);
      alert('일정 삭제에 실패했습니다.');
    }
  }
}

/**
 * 본인이 추가한 일정인지 확인
 * @param {Object} schedule - 확인할 일정 객체
 * @param {Object} context - Vue 컴포넌트 컨텍스트
 * @returns {boolean} 본인 일정 여부
 */
export function isMySchedule(schedule, context) {
  console.log('일정 체크:', {
    scheduleTitle: schedule.title,
    requestedBy: schedule.requestedBy,
    myName: context.myProfile.name,
    isMatch: schedule.requestedBy === context.myProfile.name
  });
  
  return schedule.requestedBy && 
         context.myProfile.name && 
         schedule.requestedBy === context.myProfile.name;
}

/**
 * 일정 수정 (기본 editSchedule 메서드)
 * @param {Object} context - Vue 컴포넌트 컨텍스트
 */
export function editSchedule(context) {
  // 실제 구현에서는 일정 추가 모달 표시
  alert('일정 추가 기능은 아직 구현되지 않았습니다.');
  
  // context를 실제로 사용하여 로깅 (ESLint 오류 방지)
  if (context && context.expandedMember) {
    console.log('일정 추가 준비 중...', {
      memberName: context.expandedMember.name,
      currentSchedulesCount: context.expandedMember.schedules.length
    });
  }
  
  // 일정 추가 시 유형에 따라 색상 자동 설정 예시
  // if (context && context.expandedMember) {
  //   const newSchedule = {
  //     id: Math.max(...context.expandedMember.schedules.map(s => s.id)) + 1,
  //     title: '새 일정',
  //     start: '2023-06-15T10:00',
  //     end: '2023-06-15T11:00',
  //     type: 'meeting' // 또는 'work', 'other'
  //   };
  //   newSchedule.color = getScheduleColorByType(newSchedule.type, context);
  //   context.expandedMember.schedules.push(newSchedule);
  // }
}

/**
 * 일정 유형에 따라 색상 자동 할당
 * @param {string} type - 일정 유형
 * @param {Object} context - Vue 컴포넌트 컨텍스트
 * @returns {string} 색상 코드
 */
export function getScheduleColorByType(type, context) {
  const scheduleColors = context.scheduleColors || {
    meeting: '#FF6B6B', // 회의용 빨간색
    work: '#4ECDC4',    // 업무용 청록색
    other: '#45B7D1'    // 기타용 파란색
  };
  
  switch(type) {
    case 'meeting':
      return scheduleColors.meeting;
    case 'work':
      return scheduleColors.work;
    case 'other':
      return scheduleColors.other;
    default:
      return scheduleColors.other;
  }
}

/**
 * 일정 제목으로 타입 추정
 * @param {string} title - 일정 제목
 * @returns {string} 일정 유형
 */
export function getScheduleType(title) {
  const titleLower = title.toLowerCase();
  if (titleLower.includes('회의') || titleLower.includes('미팅') || titleLower.includes('meeting')) {
    return 'meeting';
  } else if (titleLower.includes('개발') || titleLower.includes('작업') || titleLower.includes('업무') || titleLower.includes('work')) {
    return 'work';
  } else {
    return 'other';
  }
}

/**
 * 데이터 로드 관련 유틸리티 함수들
 */

import api from '@/plugins/axios.js';

/**
 * 로그인된 사용자 정보 로드
 * @param {Object} context - Vue 컴포넌트 컨텍스트
 * @param {Function} getScheduleTypeHandler - 일정 타입 추정 함수
 * @param {Function} getScheduleColorByTypeHandler - 일정 색상 할당 함수
 * @param {Function} getPositionFromRoleHandler - 역할에서 직급 추출 함수
 */
export async function loadMyProfile(context, getScheduleTypeHandler, getScheduleColorByTypeHandler, getPositionFromRoleHandler) {
  try {
    const userStr = localStorage.getItem('user');
    if (!userStr) {
      console.error('로그인 정보가 없습니다');
      return;
    }

    // 현재 로그인한 사용자 정보 가져오기
    const userResponse = await api.get('/api/user/profile');
    
    if (userResponse.data.success && userResponse.data.user) {
      const user = userResponse.data.user;
      
      // 사용자 일정 가져오기
      let schedules = [];
      try {
        const schedulesResponse = await api.get(`/api/schedules?user_id=${user.id}`);
        if (schedulesResponse.data && Array.isArray(schedulesResponse.data)) {
          schedules = schedulesResponse.data.map(schedule => {
            const scheduleType = getScheduleTypeHandler(schedule.title);
            return {
              id: schedule.schedule_id,
              title: schedule.title,
              start: schedule.start_time,
              end: schedule.end_time,
              color: schedule.color || getScheduleColorByTypeHandler(scheduleType, context),
              type: scheduleType
            };
          });
        }
      } catch (scheduleError) {
        console.error('일정 로드 오류:', scheduleError);
      }
      
      context.myProfile = {
        id: 0, // 본인 노드는 항상 id: 0
        name: user.name,
        position: getPositionFromRoleHandler(user.role),
        department: user.department,
        team: '', // 팀 정보가 없으면 빈 문자열
        schedules: schedules
      };
    }
  } catch (error) {
    console.error('사용자 정보 로드 오류:', error);
    context.myProfile.name = '사용자 정보 로드 실패';
  }
}

/**
 * 부서별 팀원 데이터 로드
 * @param {Object} context - Vue 컴포넌트 컨텍스트
 * @param {Function} getScheduleColorByTypeHandler - 일정 색상 할당 함수
 * @param {Function} getScheduleTypeHandler - 일정 타입 추정 함수
 */
export async function loadTeamMembers(context, getScheduleColorByTypeHandler, getScheduleTypeHandler) {
  try {
    console.log('🚀 === 팀원 데이터 로딩 시작 ===');
    
    // 실제 API 호출로 변경
    const response = await api.get('/api/schedule-sharing/users/all');
    
    if (response.data.success) {
      // API 응답 데이터를 teamMembers 형식으로 변환
      context.teamMembers = response.data.data || [];
      console.log('✅ 팀원 데이터 로드 완료:', context.teamMembers);
      
      // 각 팀원의 ID 정보 로깅
      context.teamMembers.forEach((member, index) => {
        console.log(`👤 팀원 ${index + 1}:`, {
          id: member.id,
          name: member.name,
          role: member.role,
          department: member.department,
          idType: typeof member.id
        });
      });
      
      // 각 사용자의 schedules 필드 초기화 및 변환
      context.teamMembers.forEach(member => {
        // schedules가 없으면 빈 배열로 초기화
        if (!member.schedules) {
          member.schedules = [];
        }
        
        // schedules가 배열이 아니면 빈 배열로 초기화
        if (!Array.isArray(member.schedules)) {
          member.schedules = [];
        }
        
        // 일정 데이터 변환
        if (member.schedules.length > 0) {
          member.schedules = member.schedules.map(schedule => ({
            id: schedule.schedule_id || schedule.id,
            title: schedule.title,
            start: schedule.start_time,
            end: schedule.end_time,
            color: schedule.color || getScheduleColorByTypeHandler(schedule.type, context),
            type: schedule.type || getScheduleTypeHandler(schedule.title),
            requestedBy: member.name
          }));
        }
      });
      
    } else {
      throw new Error(response.data.error || '데이터 로드 실패');
    }
    
  } catch (error) {
    console.error('❌ 팀원 데이터 로드 오류:', error);
    console.log('🔄 폴백 데이터로 전환합니다...');
    loadFallbackTeamMembers(context);
  }
}

/**
 * 폴백 데이터 로드 (API 실패 시 사용)
 * @param {Object} context - Vue 컴포넌트 컨텍스트
 */
export function loadFallbackTeamMembers(context) {
  console.log('폴백 데이터 로드 중...');
  context.teamMembers = [
    // 개발부 4명
    {
      id: 1,
      name: '김부장',
      role: '부장',
      department: '개발부',
      schedules: []
    },
    {
      id: 2,
      name: '이차장',
      role: '차장',
      department: '개발부',
      schedules: []
    },
    {
      id: 3,
      name: '박과장',
      role: '과장',
      department: '개발부',
      schedules: []
    },
    {
      id: 4,
      name: '최대리',
      role: '대리',
      department: '개발부',
      schedules: []
    },
    // 디자인부 3명
    {
      id: 5,
      name: '정사원',
      role: '사원',
      department: '디자인부',
      schedules: []
    },
    {
      id: 6,
      name: '한부장',
      role: '부장',
      department: '디자인부',
      schedules: []
    },
    {
      id: 7,
      name: '윤과장',
      role: '과장',
      department: '디자인부',
      schedules: []
    },
    // 마케팅부 3명
    {
      id: 8,
      name: '조차장',
      role: '차장',
      department: '마케팅부',
      schedules: []
    },
    {
      id: 9,
      name: '임대리',
      role: '대리',
      department: '마케팅부',
      schedules: []
    },
    {
      id: 10,
      name: '오사원',
      role: '사원',
      department: '마케팅부',
      schedules: []
    }
  ];
  console.log('폴백 데이터 로드 완료:', context.teamMembers);
}

/**
 * 역할(role)에서 직급 추출
 * @param {string} role - 사용자 역할
 * @returns {string} 직급명
 */
export function getPositionFromRole(role) {
  const roleMap = {
    'admin': '관리자',
    'manager': '부장',
    'team_leader': '팀장',
    'senior': '과장',
    'junior': '대리',
    'intern': '사원'
  };
  return roleMap[role] || '사원';
}

/**
 * 부서별 사용자 조회
 * @param {string} department - 부서명
 * @returns {Array} 부서별 사용자 목록
 */
export async function loadUsersByDepartment(department) {
  try {
    const response = await api.get(`/api/schedule-sharing/users/department/${encodeURIComponent(department)}`);
    
    if (response.data.success) {
      return response.data.data;
    } else {
      throw new Error(response.data.error || '부서별 사용자 조회 실패');
    }
  } catch (error) {
    console.error(`${department} 부서 사용자 조회 오류:`, error);
    return [];
  }
}

/**
 * 특정 사용자의 일정 조회
 * @param {string|number} userId - 사용자 ID
 * @param {Object} context - Vue 컴포넌트 컨텍스트
 * @param {Function} getScheduleColorByTypeHandler - 일정 색상 할당 함수
 * @param {Function} getScheduleTypeHandler - 일정 타입 추정 함수
 * @returns {Array} 사용자 일정 목록
 */
export async function loadUserSchedules(userId, context, getScheduleColorByTypeHandler, getScheduleTypeHandler) {
  console.log('📡 === loadUserSchedules 시작 ===');
  console.log('👤 요청할 userId:', userId);
  console.log('🔗 API URL:', `/api/schedules?user_id=${userId}`);
  
  try {
    console.log('🚀 API 요청 전송 중...');
    const response = await api.get(`/api/schedules?user_id=${userId}`);
    
    console.log('📥 API 응답 받음:', response);
    console.log('📊 응답 상태:', response.status);
    console.log('📋 응답 데이터:', response.data);
    
    if (response.data && Array.isArray(response.data)) {
      const mappedSchedules = response.data.map(schedule => ({
        id: schedule.id,
        title: schedule.title,
        start: schedule.start_time,
        end: schedule.end_time,
        color: schedule.color || getScheduleColorByTypeHandler(schedule.type, context),
        type: schedule.type || getScheduleTypeHandler(schedule.title)
      }));
      
      console.log('🔄 매핑된 일정 데이터:', mappedSchedules);
      console.log('✅ === loadUserSchedules 성공 ===');
      return mappedSchedules;
    } else {
      console.log('⚠️ 응답 데이터가 배열이 아님:', response.data);
      console.log('✅ === loadUserSchedules 완료 (빈 배열 반환) ===');
      return [];
    }
  } catch (error) {
    console.error('❌ === loadUserSchedules 오류 발생 ===');
    console.error('🚨 오류 객체:', error);
    console.error('📝 오류 메시지:', error.message);
    console.error('🔢 HTTP 상태:', error.response?.status);
    console.error('📋 오류 응답 데이터:', error.response?.data);
    console.error('🌐 네트워크 오류:', error.code);
    throw error;
  }
}

/**
 * 사용자 일정 데이터 변환 (공통 로직)
 * @param {Array} schedules - 원본 일정 데이터
 * @param {Object} context - Vue 컴포넌트 컨텍스트
 * @param {Function} getScheduleColorByTypeHandler - 일정 색상 할당 함수
 * @param {Function} getScheduleTypeHandler - 일정 타입 추정 함수
 * @param {string} requestedBy - 요청자 이름 (선택사항)
 * @returns {Array} 변환된 일정 데이터
 */
export function transformScheduleData(schedules, context, getScheduleColorByTypeHandler, getScheduleTypeHandler, requestedBy = null) {
  if (!schedules || !Array.isArray(schedules)) {
    return [];
  }
  
  return schedules.map(schedule => ({
    id: schedule.schedule_id || schedule.id,
    title: schedule.title,
    start: schedule.start_time,
    end: schedule.end_time,
    color: schedule.color || getScheduleColorByTypeHandler(schedule.type, context),
    type: schedule.type || getScheduleTypeHandler(schedule.title),
    ...(requestedBy && { requestedBy })
  }));
}

/**
 * API 응답 데이터 검증
 * @param {Object} response - API 응답 객체
 * @param {string} endpoint - API 엔드포인트명
 * @returns {boolean} 응답 유효성 여부
 */
export function validateApiResponse(response, endpoint) {
  if (!response || !response.data) {
    console.error(`❌ ${endpoint}: 응답 데이터가 없습니다.`);
    return false;
  }
  
  if (response.status !== 200) {
    console.error(`❌ ${endpoint}: HTTP 상태 오류 - ${response.status}`);
    return false;
  }
  
  return true;
}

/**
 * 에러 처리 및 로깅
 * @param {Error} error - 에러 객체
 * @param {string} operation - 수행 중이던 작업명
 * @param {Object} context - Vue 컴포넌트 컨텍스트 (선택사항)
 */
export function handleLoadError(error, operation, context = null) {
  console.error(`❌ ${operation} 오류 발생:`, error);
  
  if (context && context.teamMembersError) {
    context.teamMembersError = error.message;
  }
  
  // 사용자에게 에러 메시지 표시 (필요시)
  if (context && context.$message) {
    context.$message.error(`${operation} 중 오류가 발생했습니다: ${error.message}`);
  }
}

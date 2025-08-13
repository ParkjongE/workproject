import { ref } from 'vue';
import api from '@/plugins/axios';

export function useSchedule() {
  const mySchedules = ref([]);

  const selectedSchedule = ref(null);

  // 일정 목록 불러오기
  const fetchMySchedules = async () => {
    let user = localStorage.getItem('user');
    let user_id = null;
    if (user) {
      try {
        user_id = JSON.parse(user).id;
      } catch (e) { /* ignore */ }
    }
    if (!user_id) return;
    
    try {
      const res = await api.get(`/api/schedules?user_id=${user_id}`);
      console.log('[fetchMySchedules] 일정 개수:', res.data?.length || 0);
      
      // 백엔드 일정 데이터를 프론트 mySchedules 형식에 맞게 변환
      mySchedules.value = (res.data || []).map(s => {
        return {
          id: s.schedule_id,
          title: s.title,
          start: s.start_time.replace(' ', 'T'),
          end: s.end_time.replace(' ', 'T'),
          color: s.color || '#4CAF50',
          description: s.content
        };
      });
    } catch (err) {
      console.error('[fetchMySchedules] API 호출 실패:', err);
      mySchedules.value = [];
    }
  };

  // 일정 추가
  const addSchedule = async (scheduleData) => {
    try {
      await api.post('/api/schedule', scheduleData);
      await fetchMySchedules();
      return true;
    } catch (err) {
      console.error('일정 추가 실패:', err);
      return false;
    }
  };

  // 일정 삭제
  const deleteSchedule = async (scheduleId) => {
    try {
      console.log('[프론트] DELETE 요청 시작 - Schedule ID:', scheduleId);
      const response = await api.delete(`/api/schedule/delete?schedule_id=${scheduleId}`);
      
      if (response.data.success) {
        console.log('[deleteSchedule] 서버 삭제 성공:', response.data.message);
        
        // 삭제 전 일정 정보 로그
        const deletedSchedule = mySchedules.value.find(s => s.id === scheduleId);
        if (deletedSchedule) {
          console.log('[deleteSchedule] 삭제될 일정:', deletedSchedule.title, 'ID:', deletedSchedule.id);
        }
        
        // 삭제 성공 시 로컬 데이터 업데이트는 하지 않음 (서버에서 다시 불러올 예정)
        console.log('[deleteSchedule] 서버 삭제 완료, 로컬 데이터는 fetchMySchedules에서 업데이트됨');
        
        return true;
      } else {
        console.error('[deleteSchedule] 서버 삭제 실패:', response.data);
        return false;
      }
    } catch (err) {
      console.error('[deleteSchedule] 삭제 요청 실패:', err);
      return false;
    }
  };

  // 일정 수정
  const updateSchedule = async (scheduleId, updatedData) => {
    try {
      console.log('[updateSchedule] 수정 요청 시작 - Schedule ID:', scheduleId);
      console.log('[updateSchedule] 수정할 데이터:', updatedData);
      
      const response = await api.put(`/api/schedule/update?schedule_id=${scheduleId}`, updatedData);
      
      if (response.data.success) {
        console.log('[updateSchedule] 서버 수정 성공:', response.data.message);
        // 수정 성공 시 서버에서 최신 데이터 다시 불러오기
        await fetchMySchedules();
        return true;
      } else {
        console.error('[updateSchedule] 서버 수정 실패:', response.data);
        return false;
      }
    } catch (err) {
      console.error('[updateSchedule] 수정 요청 실패:', err);
      return false;
    }
  };

  // 특정 날짜의 일정 가져오기
  const getDaySchedules = (date) => {
    if (!date) return [];
    return mySchedules.value.filter(schedule => {
      const scheduleDate = schedule.start.split('T')[0];
      console.log('[getDaySchedules] 비교:', scheduleDate, date, scheduleDate === date);
      return scheduleDate === date;
    });
  };

  // 일정 선택
  const selectSchedule = (schedule) => {
    selectedSchedule.value = schedule;
  };

  return {
    mySchedules,
    selectedSchedule,
    fetchMySchedules,
    addSchedule,
    deleteSchedule,
    updateSchedule,
    getDaySchedules,
    selectSchedule,
    // updatePreviewSchedule 제거
  };
} 
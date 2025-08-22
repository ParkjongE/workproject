<template>
  <ClockSectionStyle>
    <ClockHeader>
      <ClockHeaderTitle>{{ selectedDate ? formatDate(selectedDate) : '오늘' }}의 일정</ClockHeaderTitle>
    </ClockHeader>
    <ClockView>
      <ClockFace 
        class="clock-face" 
        @click="onClockFaceClick($event)"
        @mousedown="onDragStart($event)"
        @mousemove="onDragMove($event)"
        @mouseup="onDragEnd($event)"
        @mouseleave="onDragEnd($event)"
        style="user-select: none;"
      >
        <!-- 시간 숫자 -->
        <HourNumber v-for="hour in 12" :key="`hour-${hour}`" :style="getHourNumberStyle(hour)">
          {{ hour }}
        </HourNumber>
        <!-- 시계 중앙 원 -->
        <ClockCenter />
        <!-- 시간대 선택 슬롯 -->
        <TimeSlot
          v-for="hour in 12"
          :key="`timeslot-${hour}`"
          :style="getTimeSlotStyle(hour)"
          :class="{ selected: getTimeSlotSelected(hour) }"
        >
        </TimeSlot>
        <!-- 🔥 드래그 선택 영역 -->
        <DragSelection 
          v-show="isDragging"
          :style="getDragSelectionStyle()"
        />
        <!-- 일정 세그먼트 -->
        <ScheduleSegment
          v-for="(schedule, idx) in daySchedules"
          :key="`schedule-${schedule.id || idx}`"
          :style="getScheduleSegmentStyle(schedule)"
          @mouseenter="hoveredIdx = idx"
          @mouseleave="hoveredIdx = null"
          @click.stop
        >
          <ScheduleTooltip v-show="hoveredIdx === idx">
            <ScheduleTime>
              {{ formatTime(schedule.start) }} - {{ formatTime(schedule.end) }}
            </ScheduleTime>
            <ScheduleTitle>{{ schedule.title }}</ScheduleTitle>
            <ScheduleDescription>{{ schedule.description }}</ScheduleDescription>
          </ScheduleTooltip>
        </ScheduleSegment>
        <!-- 일정 추가/요청 모달 -->
        <ModalOverlay v-if="showQuickAddForm" @click="closeQuickAddForm">
          <QuickAddForm :style="getQuickAddFormStyle()" class="quick-add-form" @click.stop>
            <QuickAddHeader>
              <QuickAddTitle>{{ isEditing ? '일정 수정' : (isOtherMember ? '일정 요청' : '일정 추가') }}</QuickAddTitle>
              <QuickAddClose @click.stop="closeQuickAddForm">×</QuickAddClose>
            </QuickAddHeader>
            
            <!-- 요청 유형 선택 (다른 구성원일 때만 표시) -->
            <div v-if="isOtherMember && !isEditing" style="margin-bottom: 16px;">
              <div style="margin-bottom: 8px; font-weight: 600; color: #333;">요청 유형:</div>
              <div style="display: flex; flex-wrap: wrap; gap: 8px;">
                <button
                  v-for="type in requestTypes"
                  :key="type.value"
                  @click="requestType = type.value"
                  :style="{
                    padding: '8px 16px',
                    border: `2px solid ${requestType === type.value ? type.color : '#ddd'}`,
                    borderRadius: '20px',
                    backgroundColor: requestType === type.value ? type.color : 'white',
                    color: requestType === type.value ? 'white' : '#333',
                    cursor: 'pointer',
                    fontSize: '14px',
                    fontWeight: requestType === type.value ? '600' : '400',
                    transition: 'all 0.2s ease'
                  }"
                >
                  {{ type.label }}
                </button>
              </div>
            </div>
            
            <!-- 디버깅용 (나중에 제거) -->
            <div v-if="isOtherMember" style="margin-bottom: 8px; padding: 8px; background-color: #f0f0f0; border-radius: 4px; font-size: 12px; color: #666;">
              Debug: isOtherMember = {{ isOtherMember }}, memberData = {{ memberData?.name || '없음' }}
            </div>
            
            <!-- 요청 유형 선택 테스트 (항상 표시) -->
            <div style="margin-bottom: 8px; padding: 8px; background-color: #e8f5e8; border-radius: 4px; font-size: 12px; color: #333;">
              Test: isOtherMember = {{ isOtherMember }}, memberData = {{ memberData?.name || '없음' }}
            </div>
            
            <input
              type="text"
              v-model="quickAddForm.title"
              :placeholder="isOtherMember ? '요청 제목을 입력하세요' : '일정 제목을 입력하세요'"
              @keyup.enter="saveQuickSchedule"
              style="width:100%;margin-bottom:8px;padding:8px;border:1px solid #ddd;border-radius:4px;font-size:14px;"
            />
            <textarea
              v-model="quickAddForm.description"
              :placeholder="isOtherMember ? '요청 내용을 입력하세요' : '일정 내용을 입력하세요'"
              rows="3"
              style="width:100%;margin-bottom:8px;padding:8px;border:1px solid #ddd;border-radius:4px;font-size:14px;resize:vertical;"
            ></textarea>
            <TimeInputContainer>
              <TimeInputGroup>
                <input
                  type="time"
                  v-model="quickAddForm.startTime"
                  step="900"
                  required
                  style="width:100%;padding:6px;border:1px solid #ddd;border-radius:4px;font-size:14px;"
                />
                <TimeLabel>{{ quickAddForm.ampm }}</TimeLabel>
              </TimeInputGroup>
              <TimeSeparator>~</TimeSeparator>
              <TimeInputGroup>
                <input
                  type="time"
                  v-model="quickAddForm.endTime"
                  step="900"
                  required
                  style="width:100%;padding:6px;border:1px solid #ddd;border-radius:4px;font-size:14px;"
                />
              </TimeInputGroup>
            </TimeInputContainer>
            <ColorSection>
              <ColorLabel>색상 선택:</ColorLabel>
              <ColorOptionsContainer>
                <ColorOptionBtn
                  v-for="color in colorOptions"
                  :key="color"
                  :color="color"
                  :selected="quickAddForm.color === color"
                  :class="{ 
                    selected: quickAddForm.color === color,
                    'color-selected': quickAddForm.color === color && colorAnimationActive
                  }"
                  @click="selectColor(color)"
                >
                  <ColorCheckIcon v-if="quickAddForm.color === color">✔</ColorCheckIcon>
                </ColorOptionBtn>
              </ColorOptionsContainer>
            </ColorSection>
            <QuickAddActions>
              <QuickAddBtn class="save" @click="saveQuickSchedule">
                {{ isEditing ? '수정' : (isOtherMember ? '요청 보내기' : '추가') }}
              </QuickAddBtn>
              <QuickAddBtn class="cancel" @click.stop="closeQuickAddForm">닫기</QuickAddBtn>
            </QuickAddActions>
          </QuickAddForm>
        </ModalOverlay>
      </ClockFace>
    </ClockView>
    <DayScheduleList>
      <h3>일정 목록</h3>
      <ScheduleItem
        v-for="(schedule, idx) in daySchedules"
        :key="idx"
        :selected="selectedSchedule && selectedSchedule.id === schedule.id"
        @click="selectSchedule(schedule)"
      >
        <ScheduleColor :color="schedule.color" />
        <ScheduleInfo>
          <div class="schedule-title">{{ schedule.title }}</div>
          <ScheduleDescriptionLarge class="schedule-desc">{{ schedule.description }}</ScheduleDescriptionLarge>
          <div class="schedule-time">{{ formatTime(schedule.start) }} - {{ formatTime(schedule.end) }}</div>
        </ScheduleInfo>
        <ScheduleActions>
          <EditBtn @click.stop="editSchedule(schedule)">수정</EditBtn>
          <DeleteBtn @click.stop="handleDeleteSchedule(schedule)">삭제</DeleteBtn>
        </ScheduleActions>
      </ScheduleItem>
    </DayScheduleList>
  </ClockSectionStyle>
</template>

<script>
import * as styled from '../Schedule.styled.js';

//
import { useSchedule } from '../composables/useSchedule';
import { useClock } from '../composables/useClock';
import { computed, ref, onMounted } from 'vue';

// Utils
import { formatDate, formatTime } from '../utils/formatters';

export default {
  name: 'ClockSection', 
  components: {
    ...styled
  },
  props: {
    selectedDate: {
      type: String,
      required: true
    },
    mySchedules: {
      type: Array,
      default: () => []
    },
    isOtherMember: {
      type: Boolean,
      default: false
    },
    memberData: {
      type: Object,
      default: null
    }
  },
  emits: ['schedule-selected', 'schedule-edited', 'schedule-deleted', 'schedule-added'],
  setup(props, { emit }) {
    // props로 받은 mySchedules를 사용하여 getDaySchedules 함수 생성
    const getDaySchedules = (date) => {
      if (!date) return [];
      return props.mySchedules.filter(schedule => {
        const scheduleDate = schedule.start.split('T')[0];
        return scheduleDate === date;
      });
    };

    // Composables 사용 (일정 추가/삭제/수정 기능)
    const {
      selectedSchedule,
      selectSchedule,
      addSchedule,
      deleteSchedule,
      updateSchedule,
      updatePreviewSchedule
    } = useSchedule();

    const {
      hoveredIdx,
      showQuickAddForm,
      quickAddForm,
      colorOptions,
      colorAnimationActive,
      // 🔥 드래그 관련 추가
      isDragging,
      selectedTimeRange,
      onDragStart,
      onDragMove,
      onDragEnd,
      getDragSelectionStyle,
      getHourNumberStyle,
      getScheduleSegmentStyle,
      getTimeSlotStyle,
      getQuickAddFormStyle,
      getTimeSlotSelected,
      parseTimeToMinutes,
      handleOutsideClick,
      closeQuickAddForm,
      onClockFaceClick,
      selectColor
    } = useClock();

    // 수정 모드 상태
    const isEditing = ref(false);
    const editingScheduleId = ref(null);
    
    // 요청 유형 (다른 구성원일 때 사용)
    const requestType = ref('meeting');
    const requestTypes = [
      { value: 'meeting', label: '🤝 회의', color: '#FF6B6B' },
      { value: 'collaboration', label: '💼 업무 협업', color: '#4ECDC4' },
      { value: 'discussion', label: '💬 논의', color: '#45B7D1' },
      { value: 'other', label: '📝 기타', color: '#9B59B6' }
    ];

    // getDaySchedules를 computed로 만들어서 props 변경 시 자동 업데이트
    const daySchedules = computed(() => {
      console.log('[ClockSection] daySchedules computed 실행, mySchedules:', props.mySchedules.length);
      return getDaySchedules(props.selectedDate) || [];
    });
    
    // mounted 라이프사이클 추가
    onMounted(() => {
      console.log('[ClockSection] mounted 실행');
      console.log('[ClockSection] isOtherMember:', props.isOtherMember);
      console.log('[ClockSection] memberData:', props.memberData);
    });

    // 빠른 일정 저장 (추가/수정) 또는 요청
    const saveQuickSchedule = async () => {
      if (!quickAddForm.value.title || !quickAddForm.value.title.trim()) {
        alert('일정 제목을 입력해주세요.');
        return;
      }
      
      let user = localStorage.getItem('user');
      let user_id = null;
      if (user) {
        try {
          user_id = JSON.parse(user).id;
        } catch (e) { /* ignore */ }
      }
      if (!user_id) {
        alert('로그인 정보가 없습니다. 다시 로그인 해주세요.');
        return;
      }

      const scheduleData = {
        user_id: user_id,
        title: quickAddForm.value.title,
        content: quickAddForm.value.description,
        project_id: null,
        start_time: `${props.selectedDate}T${quickAddForm.value.startTime}`,
        end_time: `${props.selectedDate}T${quickAddForm.value.endTime}`,
        color: quickAddForm.value.color
      };

      console.log('[saveQuickSchedule] 폼 데이터:', {
        title: quickAddForm.value.title,
        description: quickAddForm.value.description,
        content: scheduleData.content
      });
      
      // description이 비어있는지 확인
      if (!quickAddForm.value.description || quickAddForm.value.description.trim() === '') {
        console.warn('[saveQuickSchedule] description이 비어있음!');
      } else {
        console.log('[saveQuickSchedule] description 길이:', quickAddForm.value.description.length);
      }

      let success = false;
      
      if (isEditing.value) {
        // 수정 모드
        console.log('[일정수정] 수정할 데이터:', scheduleData);
        success = await updateSchedule(editingScheduleId.value, scheduleData);
        
        if (success) {
          closeQuickAddForm();
          emit('schedule-edited', { id: editingScheduleId.value, ...scheduleData });
        } else {
          alert('일정 수정에 실패했습니다.');
        }
      } else if (props.isOtherMember) {
        // 다른 구성원에게 일정 요청 모드
        const requestData = {
          ...scheduleData,
          requestType: requestType.value,
          requestTypeLabel: requestTypes.find(t => t.value === requestType.value)?.label || '기타',
          targetMember: props.memberData?.name || '구성원',
          requester: user_id
        };
        
        console.log('[일정요청] 요청할 데이터:', requestData);
        console.log('[일정요청] 대상 구성원:', props.memberData?.name);
        console.log('[일정요청] 요청 유형:', requestType.value);
        
        // 요청 완료 메시지 표시
        alert(`${props.memberData?.name || '구성원'}님께 ${requestTypes.find(t => t.value === requestType.value)?.label || '일정'} 요청이 전송되었습니다.`);
        closeQuickAddForm();
        emit('schedule-added');
      } else {
        // 본인 일정 추가 모드
        console.log('[일정추가] 추가할 데이터:', scheduleData);
        success = await addSchedule(scheduleData);
        
        if (success) {
          closeQuickAddForm();
          emit('schedule-added');
        } else {
          alert('일정 추가에 실패했습니다.');
        }
      }
    };

    // 일정 편집 모달 열기
    const editSchedule = (schedule) => {
      console.log('[ClockSection] 일정 수정 모달 열기:', schedule);
      console.log('[ClockSection] 일정 내용:', schedule.description);
      
      // 수정 모드로 설정
      isEditing.value = true;
      editingScheduleId.value = schedule.id;
      
      // 폼에 기존 데이터 채우기
      quickAddForm.value.title = schedule.title || '';
      quickAddForm.value.description = schedule.description || '';
      quickAddForm.value.color = schedule.color || '#4CAF50';
      
      // 시간 형식 변환 (ISO 문자열에서 시간만 추출)
      const startTime = schedule.start.split('T')[1]?.substring(0, 5) || '09:00';
      const endTime = schedule.end.split('T')[1]?.substring(0, 5) || '10:00';
      
      quickAddForm.value.startTime = startTime;
      quickAddForm.value.endTime = endTime;
      
      console.log('[ClockSection] 폼에 채워진 데이터:', {
        title: quickAddForm.value.title,
        description: quickAddForm.value.description,
        color: quickAddForm.value.color,
        startTime: quickAddForm.value.startTime,
        endTime: quickAddForm.value.endTime
      });
      
      // 모달 표시
      showQuickAddForm.value = true;
    };

    // 모달 닫기 시 수정 모드 초기화
    const handleCloseQuickAddForm = () => {
      closeQuickAddForm();
      isEditing.value = false;
      editingScheduleId.value = null;
      
      // 폼 데이터 초기화
      quickAddForm.value.title = '';
      quickAddForm.value.description = '';
      quickAddForm.value.color = '#4CAF50';
      quickAddForm.value.startTime = '09:00';
      quickAddForm.value.endTime = '10:00';
      
      console.log('[ClockSection] 수정 모드 초기화 및 폼 데이터 초기화');
    };

    // 일정 삭제
    const handleDeleteSchedule = async (schedule) => {
      if (confirm('이 일정을 삭제하시겠습니까?')) {
        console.log('[ClockSection] 일정 삭제 시작:', schedule.title, 'ID:', schedule.id);
        
        const success = await deleteSchedule(schedule.id);
        if (success) {
          console.log('[ClockSection] 일정 삭제 성공, 이벤트 emit');
          emit('schedule-deleted', schedule);
        } else {
          console.error('[ClockSection] 일정 삭제 실패');
          alert('일정 삭제에 실패했습니다.');
        }
      } else {
        console.log('[ClockSection] 일정 삭제 취소됨');
      }
    };

    // 일정 선택
    const handleSelectSchedule = (schedule) => {
      selectSchedule(schedule);
      emit('schedule-selected', schedule);
    };

    return {
      // Clock
      hoveredIdx,
      showQuickAddForm,
      quickAddForm,
      colorOptions,
      colorAnimationActive,
      // 🔥 드래그 관련 추가
      isDragging,
      selectedTimeRange,
      onDragStart,
      onDragMove,
      onDragEnd,
      getDragSelectionStyle,
      getHourNumberStyle,
      getScheduleSegmentStyle,
      getTimeSlotStyle,
      getQuickAddFormStyle,
      getTimeSlotSelected,
      parseTimeToMinutes,
      handleOutsideClick,
      onClockFaceClick, // 추가
      
      // Schedule
      selectedSchedule,
      daySchedules,
      selectSchedule: handleSelectSchedule,
      addSchedule,
      deleteSchedule,
      updateSchedule,
      updatePreviewSchedule,
      
      // Edit Mode
      isEditing,
      editingScheduleId,
      
      // Methods
      saveQuickSchedule,
      editSchedule,
      closeQuickAddForm: handleCloseQuickAddForm,
      handleDeleteSchedule,
      selectColor,
      
      // Utils
      formatDate,
      formatTime
    };
  }
};
</script> 
import { ref } from 'vue';

export function useClock() {
  const hoveredIdx = ref(null);
  const showQuickAddForm = ref(false);
  const colorAnimationActive = ref(false);

  // 🔥 드래그 상태 관리
  const isDragging = ref(false);
  const dragStartTime = ref(null);
  const dragEndTime = ref(null);
  const selectedTimeRange = ref({ start: null, end: null });

  const quickAddForm = ref({
    title: '',
    description: '',
    startTime: '',
    endTime: '',
    color: '#4CAF50',
    ampm: '',
  });

  const colorOptions = [
    '#4CAF50', // 초록
    '#2196F3', // 파랑
    '#FF9800', // 주황
    '#F44336', // 빨강
    '#9C27B0', // 보라
    '#00BCD4', // 청록
    '#FFEB3B', // 노랑
    '#795548'  // 갈색
  ];

  // 🔥 시계에서 시간 계산 (각도를 시간으로 변환)
  const getTimeFromPosition = (event, element) => {
    const rect = element.getBoundingClientRect();
    const centerX = rect.left + rect.width / 2;
    const centerY = rect.top + rect.height / 2;
    const x = event.clientX - centerX;
    const y = event.clientY - centerY;
    
    // 각도 계산 (0도는 3시 방향)
    let angle = Math.atan2(y, x) * 180 / Math.PI;
    // 12시 방향을 0도로 맞추기 위해 조정
    angle = (angle + 450) % 360;
    
    // 시간 계산 (12시간 단위)
    const hour = Math.floor(angle / 30);
    const minute = Math.floor((angle % 30) * 2); // 30도당 60분이므로 2배
    
    // 🔥 15분 단위 스냅 적용
    const snappedTime = snapToQuarterHour(hour === 0 ? 12 : hour, minute);
    
    return snappedTime;
  };

  // 🔥 15분 단위 스냅 함수 (0, 15, 30, 45분으로 조정)
  const snapToQuarterHour = (hour, minute) => {
    const originalTime = `${hour}:${String(minute).padStart(2, '0')}`;
    
    let snappedMinute;
    if (minute < 7.5) {
      snappedMinute = 0;
    } else if (minute < 22.5) {
      snappedMinute = 15;
    } else if (minute < 37.5) {
      snappedMinute = 30;
    } else if (minute < 52.5) {
      snappedMinute = 45;
    } else {
      // 52.5분 이상은 다음 시간 00분으로
      const nextHour = hour === 12 ? 1 : hour + 1;
      console.log(`[15분 스냅] ${originalTime} → ${nextHour}:00`);
      return { hour: nextHour, minute: 0 };
    }
    
    console.log(`[15분 스냅] ${originalTime} → ${hour}:${String(snappedMinute).padStart(2, '0')}`);
    return { hour, minute: snappedMinute };
  };

  // 🔥 시간을 문자열로 포맷팅
  const formatTimeString = (hour, minute) => {
    return `${String(hour).padStart(2, '0')}:${String(minute).padStart(2, '0')}`;
  };

  // 🔥 12시간 형식을 24시간 형식으로 변환 (AM/PM 기본값 적용)
  const convertTo24Hour = (hour12) => {
    // 1시~7시: 오후(PM) → 13시~19시
    if (hour12 >= 1 && hour12 <= 7) {
      return hour12 + 12;
    }
    // 8시~12시: 오전(AM) → 8시~12시 (12시는 0시로 변환)
    else if (hour12 >= 8 && hour12 <= 11) {
      return hour12;
    }
    // 12시: 오전 → 0시
    else if (hour12 === 12) {
      return 0;
    }
    return hour12;
  };

  // 🔥 드래그 시작
  const onDragStart = (event) => {
    if (showQuickAddForm.value) return;
    
    isDragging.value = true;
    const time = getTimeFromPosition(event, event.currentTarget);
    dragStartTime.value = time;
    dragEndTime.value = time;
    selectedTimeRange.value = { start: time, end: time };
    
    console.log('[Drag Start] 시작 시간:', formatTimeString(time.hour, time.minute));
  };

  // 🔥 드래그 중
  const onDragMove = (event) => {
    if (!isDragging.value) return;
    
    event.preventDefault();
    const time = getTimeFromPosition(event, event.currentTarget);
    dragEndTime.value = time;
    
    // 🔥 시계방향 체크: 시작점에서 시계방향으로만 드래그 허용
    const start = dragStartTime.value;
    const end = time;
    
    // 각도로 변환해서 시계방향 체크
    const startAngle = ((start.hour % 12) * 30 + start.minute * 0.5);
    const endAngle = ((end.hour % 12) * 30 + end.minute * 0.5);
    
    // 시계방향 각도 차이 계산
    let angleDiff = endAngle - startAngle;
    if (angleDiff < 0) angleDiff += 360;
    
    // 🔥 시계방향으로만 허용 (0도~360도)
    if (angleDiff >= 0 && angleDiff <= 360) {
      selectedTimeRange.value = { start, end };
      
      console.log('[Drag Move] 시계방향 드래그:', 
        `${formatTimeString(start.hour, start.minute)} → ${formatTimeString(end.hour, end.minute)}`,
        `(${angleDiff.toFixed(1)}도)`
      );
    } else {
      // 반시계방향은 무시하고 시작점으로 리셋
      selectedTimeRange.value = { start, end: start };
      
      console.log('[Drag Move] 반시계방향 무시됨');
    }
  };

  // 🔥 드래그 종료
  const onDragEnd = () => {
    if (!isDragging.value) return;
    
    isDragging.value = false;
    
    // 선택된 시간대로 폼 설정하고 모달 열기
    if (selectedTimeRange.value.start && selectedTimeRange.value.end) {
      const startHour24 = convertTo24Hour(selectedTimeRange.value.start.hour);
      const endHour24 = convertTo24Hour(selectedTimeRange.value.end.hour);
      
      const startTime = formatTimeString(startHour24, selectedTimeRange.value.start.minute);
      const endTime = formatTimeString(endHour24, selectedTimeRange.value.end.minute);
      
      quickAddForm.value = {
        title: '',
        description: '',
        startTime,
        endTime,
        color: '#4CAF50',
        ampm: '',
      };
      
      showQuickAddForm.value = true;
      
      console.log('[Drag End] 일정 추가 모달 열기:', startTime, '~', endTime);
      console.log('[Drag End] 12시간→24시간 변환:', 
        `${selectedTimeRange.value.start.hour}시→${startHour24}시`, 
        `${selectedTimeRange.value.end.hour}시→${endHour24}시`);
    }
    
    // 드래그 상태 초기화
    dragStartTime.value = null;
    dragEndTime.value = null;
  };

  // 시계 숫자 스타일 계산
  const getHourNumberStyle = (hour) => {
    const clockSize = 520;
    const radius = 240;
    const centerX = clockSize / 2;
    const centerY = clockSize / 2;
    const angle = (hour / 12) * 2 * Math.PI - Math.PI / 2;
    const x = Math.cos(angle) * radius + centerX;
    const y = Math.sin(angle) * radius + centerY;
    return {
      position: 'absolute',
      left: `${x}px`,
      top: `${y}px`,
      transform: 'translate(-50%, -50%)',
      fontSize: '28px',
      fontWeight: 'bold',
      color: '#333',
      zIndex: 2
    };
  };

  // 일정 세그먼트 스타일 계산
  const getScheduleSegmentStyle = (schedule) => {
    const startParts = schedule.start.split('T')[1].split(':');
    const endParts = schedule.end.split('T')[1].split(':');
    let startHour = parseInt(startParts[0]);
    const startMinute = parseInt(startParts[1]);
    let endHour = parseInt(endParts[0]);
    const endMinute = parseInt(endParts[1]);
    
    function to12Hour(hour) {
      if (hour === 0) return 12;
      if (hour > 12) return hour - 12;
      return hour;
    }
    
    const startHour12 = to12Hour(startHour);
    const endHour12 = to12Hour(endHour);
    const startDeg = ((startHour12 % 12) * 30 + startMinute * 0.5);
    let endDeg = ((endHour12 % 12) * 30 + endMinute * 0.5);
    if (endDeg <= startDeg) endDeg += 360;
    const sweepDeg = endDeg - startDeg;
    const clockSize = 520;
    
    return {
      position: 'absolute',
      left: '50%',
      top: '50%',
      width: `${clockSize}px`,
      height: `${clockSize}px`,
      transform: 'translate(-50%, -50%)',
      background: `conic-gradient(from ${startDeg}deg, ${schedule.color} 0deg, ${schedule.color} ${sweepDeg}deg, transparent ${sweepDeg}deg 360deg)`,
      borderRadius: '50%',
      opacity: '0.7',
      zIndex: 0,
      pointerEvents: 'none',
      cursor: 'pointer'
    };
  };

  // 🔥 드래그 선택 영역 스타일 계산
  const getDragSelectionStyle = () => {
    if (!isDragging.value || !selectedTimeRange.value.start || !selectedTimeRange.value.end) {
      return { display: 'none' };
    }

    const start = selectedTimeRange.value.start;
    const end = selectedTimeRange.value.end;
    
    // 🔥 시작점과 끝점이 같으면 색칠하지 않음
    if (start.hour === end.hour && start.minute === end.minute) {
      return { display: 'none' };
    }
    
    // 시작과 끝 각도 계산
    const startDeg = ((start.hour % 12) * 30 + start.minute * 0.5);
    const endDeg = ((end.hour % 12) * 30 + end.minute * 0.5);
    
    // 🔥 무조건 시계방향으로 색칠
    let sweepDeg = endDeg - startDeg;
    if (sweepDeg < 0) sweepDeg += 360;  // 음수면 360도 더하기 (시계방향)
    
    // 🔥 sweepDeg가 0이면 색칠하지 않음
    if (sweepDeg === 0) {
      return { display: 'none' };
    }
    
    const clockSize = 520;
    
    return {
      position: 'absolute',
      left: '50%',
      top: '50%',
      width: `${clockSize}px`,
      height: `${clockSize}px`,
      transform: 'translate(-50%, -50%)',
      background: `conic-gradient(from ${startDeg}deg, rgba(76, 175, 80, 0.3) 0deg, rgba(76, 175, 80, 0.3) ${sweepDeg}deg, transparent ${sweepDeg}deg 360deg)`,
      borderRadius: '50%',
      zIndex: 1,
      pointerEvents: 'none'
    };
  };

  // 시간대 선택 슬롯 스타일
  const getTimeSlotStyle = (hour) => {
    const clockSize = 520;
    const radius = 240;
    const centerX = clockSize / 2;
    const centerY = clockSize / 2;
    const angle = (hour / 12) * 2 * Math.PI - Math.PI / 2;
    const x = Math.cos(angle) * radius + centerX;
    const y = Math.sin(angle) * radius + centerY;
    
    return {
      position: 'absolute',
      left: `${x}px`,
      top: `${y}px`,
      transform: 'translate(-50%, -50%)',
      zIndex: 4
    };
  };

  // 빠른 추가 폼 위치 계산
  const getQuickAddFormStyle = () => {
    return {
      position: 'absolute',
      top: '50%',
      left: '50%',
      transform: 'translate(-50%, -50%)',
      zIndex: 1000
    };
  };

  // 🔥 수정된 시계 클릭 이벤트 처리 (단순 클릭만 처리)
  const onClockFaceClick = (event) => {
    // 드래그 중이거나 이미 모달이 열려있으면 무시
    if (showQuickAddForm.value || isDragging.value) return;
    
    const rect = event.currentTarget.getBoundingClientRect();
    const centerX = rect.left + rect.width / 2;
    const centerY = rect.top + rect.height / 2;
    const x = event.clientX - centerX;
    const y = event.clientY - centerY;
    let angle = Math.atan2(y, x) * 180 / Math.PI;
    angle = (angle + 450) % 360;
    let hour = Math.floor(angle / 30);
    hour = (hour + 12) % 12;
    if (hour === 0) hour = 12;
    
    // 🔥 클릭한 시간을 15분 단위로 스냅
    const clickedTime = snapToQuarterHour(hour, 0); // 정각에서 시작
    const startHour24 = convertTo24Hour(clickedTime.hour);
    
    // 🔥 끝 시간은 1시간 후로 설정 (15분 단위 유지)
    const endHour = clickedTime.hour === 12 ? 1 : clickedTime.hour + 1;
    const endHour24 = convertTo24Hour(endHour);
    
    showQuickAddForm.value = true;
    
    const currentDescription = quickAddForm.value.description || '';
    quickAddForm.value = {
      title: '',
      description: currentDescription,
      startTime: `${String(startHour24).padStart(2, '0')}:${String(clickedTime.minute).padStart(2, '0')}`,
      endTime: `${String(endHour24).padStart(2, '0')}:00`,
      color: '#4CAF50',
      ampm: '',
    };
    
    console.log('[Click] 15분 단위 시간 설정:', 
      `${clickedTime.hour}:${String(clickedTime.minute).padStart(2, '0')} → ${startHour24}:${String(clickedTime.minute).padStart(2, '0')}`,
      `끝: ${endHour24}:00`);
  };

  // 색상 선택
  const selectColor = (color) => {
    quickAddForm.value.color = color;
    colorAnimationActive.value = true;
    setTimeout(() => {
      colorAnimationActive.value = false;
    }, 300);
  };

  // 빠른 일정 추가 폼 닫기
  const closeQuickAddForm = () => {
    showQuickAddForm.value = false;
    // 🔥 드래그 선택도 초기화
    selectedTimeRange.value = { start: null, end: null };
    quickAddForm.value = {
      title: '',
      description: '',
      startTime: '',
      endTime: '',
      color: '#4CAF50',
      ampm: '',
    };
  };

  // 시간대 선택 확인
  const getTimeSlotSelected = (hour) => {
    if (showQuickAddForm.value && quickAddForm.value.startTime && quickAddForm.value.endTime) {
      const start = parseTimeToMinutes(quickAddForm.value.startTime);
      const end = parseTimeToMinutes(quickAddForm.value.endTime);
      return hour >= start && hour <= end;
    }
    return false;
  };

  // 시간을 분으로 변환
  const parseTimeToMinutes = (timeString) => {
    const [hours, minutes] = timeString.split(':').map(Number);
    return hours * 60 + minutes;
  };

  // 시계 외부 클릭 처리
  const handleOutsideClick = (event) => {
    const clockFace = event.target.closest('.clock-face');
    const quickAddFormElement = event.target.closest('.quick-add-form');
    
    if (!clockFace && !quickAddFormElement && showQuickAddForm.value) {
      closeQuickAddForm();
    }
  };

  return {
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
    onClockFaceClick,
    selectColor,
    closeQuickAddForm,
    getTimeSlotSelected,
    parseTimeToMinutes,
    handleOutsideClick
  };
} 
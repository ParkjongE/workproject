import styled from '@vue-styled-components/core';

const colorOptionProps = { color: String, selected: Boolean };

// === 레이아웃 ===
export const MyScheduleRoot = styled.div`
  display: flex;
  height: 100vh;
  width: 100%;
`;

export const ContentArea = styled.div`
  flex: 1;
  padding: 20px;
  overflow-y: auto;
  background-color: #f5f5f5;
  min-height: 100vh;
`;

export const HeaderContainer = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  background-color: white;
  padding: 15px 20px;
  border-radius: 8px;
  box-shadow: 0 2px 5px rgba(0,0,0,0.1);
`;

export const Title = styled.h1`
  margin-bottom: 0;
  color: #333;
`;

export const EditPrivacyBtn = styled.button`
  background-color: #f0f0f0;
  border: none;
  border-radius: 4px;
  padding: 8px 12px;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 14px;
  color: #2196F3;
`;

export const DonePrivacyBtn = styled.button`
  background-color: #4CAF50;
  color: white;
  border: none;
  border-radius: 4px;
  padding: 8px 12px;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 14px;
`;

export const ScheduleManagement = styled.div`
  display: flex;
  gap: 20px;
  height: calc(100vh - 120px);
  margin-top: 20px;
  align-items: flex-start;
`;

// === 시계 섹션 ===
export const ClockSectionStyle = styled.div`
  flex: 0 0 550px;
  display: flex;
  flex-direction: column;
  gap: 15px;
  overflow-y: auto;
`;

export const ClockHeader = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 15px;
  background-color: white;
  padding: 10px 15px;
  border-radius: 6px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.1);
`;

export const ClockHeaderTitle = styled.h2`
  margin: 0;
  color: #333;
  font-size: 16px;
`;

export const ClockView = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  margin: 10px 0;
`;

export const ClockFace = styled.div`
  position: relative;
  width: 520px;
  max-width: 90vw;
  height: 520px;
  max-height: 90vw;
  border-radius: 50%;
  border: 4px solid #bbb;
  background-color: white;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
  margin: 0 auto;
`;

export const HourNumber = styled.div`
  position: absolute;
  font-size: 36px;
  color: #333;
  font-weight: bold;
  z-index: 2;
  /* 위치 계산은 getHourNumberStyle에서 반지름 220~230px 기준으로 조정 */
`;

export const ClockCenter = styled.div`
  position: absolute;
  width: 18px;
  height: 18px;
  border-radius: 50%;
  background-color: #888;
  left: 50%;
  top: 50%;
  transform: translate(-50%, -50%);
  z-index: 3;
  border: 3px solid #fff;
`;

export const ScheduleSegment = styled.div`
  position: absolute;
  left: 50%;
  top: 50%;
  transform: translate(-50%, -50%);
  z-index: 1;
  opacity: 0.7;
  cursor: pointer;
  &:hover {
    opacity: 0.9;
  }
  
  &.preview-schedule {
    opacity: 0.5;
    animation: previewPulse 2s ease-in-out infinite;
    border: 2px dashed rgba(255, 255, 255, 0.8);
  }
  
  @keyframes previewPulse {
    0%, 100% { opacity: 0.5; }
    50% { opacity: 0.8; }
  }
`;

export const ScheduleTooltip = styled.div`
  display: none;
  position: absolute;
  left: 100%;
  top: 0;
  background-color: white;
  border: 1px solid #ccc;
  padding: 8px;
  border-radius: 4px;
  box-shadow: 0 2px 5px rgba(0,0,0,0.2);
  z-index: 10;
  white-space: nowrap;
`;

export const ScheduleTime = styled.div`
  font-size: 12px;
  color: #666;
`;

export const ScheduleTitle = styled.div`
  font-weight: bold;
`;

// === 시간대 선택을 위한 스타일 ===
export const TimeSlot = styled.div`
  position: absolute;
  width: 60px;
  height: 60px;
  border-radius: 50%;
  background-color: transparent;
  cursor: pointer;
  z-index: 4;
  transition: all 0.2s ease;
  
  &:hover {
    background-color: rgba(33, 150, 243, 0.15);
    border: 2px solid rgba(33, 150, 243, 0.5);
    transform: translate(-50%, -50%) scale(1.1);
    box-shadow: 0 2px 8px rgba(33, 150, 243, 0.3);
  }
  
  &.selected {
    background-color: rgba(33, 150, 243, 0.25);
    border: 2px solid #2196F3;
    transform: translate(-50%, -50%) scale(1.05);
    box-shadow: 0 2px 8px rgba(33, 150, 243, 0.4);
  }

  &.dragging {
    box-shadow: 0 0 0 4px rgba(33,150,243,0.15);
    opacity: 0.85;
  }
`;

export const TimeSlotLabel = styled.div`
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  font-size: 12px;
  font-weight: bold;
  color: #333;
  pointer-events: none;
`;

export const QuickAddForm = styled.div`
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  background-color: white;
  border: 1px solid #ccc;
  border-radius: 8px;
  padding: 15px;
  box-shadow: 0 4px 12px rgba(0,0,0,0.15);
  z-index: 1000;
  min-width: 250px;
  animation: fadeIn 0.2s ease-in-out;
  
  @keyframes fadeIn {
    from {
      opacity: 0;
      transform: translate(-50%, -50%) scale(0.9);
    }
    to {
      opacity: 1;
      transform: translate(-50%, -50%) scale(1);
    }
  }
`;

export const QuickAddHeader = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
  padding-bottom: 8px;
  border-bottom: 1px solid #eee;
`;

export const QuickAddTitle = styled.h4`
  margin: 0;
  color: #333;
  font-size: 14px;
`;

export const QuickAddClose = styled.button`
  background: none;
  border: none;
  font-size: 18px;
  cursor: pointer;
  color: #666;
  padding: 0;
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  
  &:hover {
    color: #333;
  }
`;



export const QuickAddActions = styled.div`
  display: flex;
  gap: 8px;
  justify-content: flex-end;
`;

export const QuickAddBtn = styled.button`
  padding: 6px 12px;
  border: none;
  border-radius: 4px;
  font-size: 12px;
  cursor: pointer;
  
  &.cancel {
    background-color: #f5f5f5;
    color: #666;
  }
  
  &.save {
    background-color: #2196F3;
    color: white;
  }
  
  &:hover {
    opacity: 0.8;
  }
`;

export const ColorOptionBtn = styled('button', colorOptionProps)`
  width: 28px;
  height: 28px;
  border-radius: 50%;
  border: ${props => props.selected ? '3px solid #333' : '1px solid #ddd'};
  outline: none;
  cursor: pointer;
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: ${props => props.selected ? '0 0 0 2px #2196F3' : 'none'};
  transition: all 0.2s ease;
  background-color: ${props => props.color};
  
  &:hover {
    transform: scale(1.05);
    box-shadow: 0 2px 8px rgba(0,0,0,0.2);
  }
  
  &.selected {
    transform: scale(1.1);
    box-shadow: 0 0 0 3px #2196F3, 0 2px 8px rgba(33, 150, 243, 0.4);
  }
  
  /* 색상 선택 시 애니메이션 효과 */
  &.color-selected {
    animation: colorPulse 0.3s ease-in-out;
  }
  
  @keyframes colorPulse {
    0% { transform: scale(1); }
    50% { transform: scale(1.2); }
    100% { transform: scale(1.1); }
  }
`;

export const ColorSection = styled.div`
  margin-bottom: 8px;
`;

export const ColorLabel = styled.div`
  font-size: 12px;
  color: #666;
  margin-bottom: 4px;
`;

export const ColorOptionsContainer = styled.div`
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
`;

export const ColorCheckIcon = styled.span`
  color: #333;
  font-size: 18px;
  font-weight: bold;
`;

export const QuickAddInput = styled.input`
  width: 100%;
  margin-bottom: 8px;
  padding: 8px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
  
  &:focus {
    outline: none;
    border-color: #2196F3;
  }
`;

export const QuickAddTextarea = styled.textarea`
  width: 100%;
  margin-bottom: 8px;
  resize: vertical;
  padding: 8px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
  font-family: inherit;
  
  &:focus {
    outline: none;
    border-color: #2196F3;
  }
`;

export const TimeInputContainer = styled.div`
  display: flex;
  gap: 8px;
  margin-bottom: 8px;
  align-items: center;
`;

export const TimeInputGroup = styled.div`
  flex: 1;
  display: flex;
  align-items: center;
`;

export const TimeInput = styled.input`
  flex: 1;
  padding: 6px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
  
  &:focus {
    outline: none;
    border-color: #2196F3;
  }
`;

export const TimeLabel = styled.span`
  margin-left: 4px;
  font-size: 12px;
  color: #666;
  min-width: 2.5em;
`;

export const TimeSeparator = styled.span`
  align-self: center;
  font-weight: bold;
  color: #666;
`;

export const ScheduleDescription = styled.div`
  font-size: 12px;
  color: #666;
  white-space: pre-line;
`;

export const ScheduleDescriptionLarge = styled.div`
  font-size: 14px;
  color: #666;
  white-space: pre-line;
`;

export const ModalOverlay = styled.div`
  position: fixed;
  z-index: 9999;
  left: 0;
  top: 0;
  right: 0;
  bottom: 0;
  background: rgba(0,0,0,0.1);
`;

// === 정보 카드 섹션 ===
export const InfoCardSectionStyled = styled.div`
  flex: 0 0 320px;
  display: flex;
  flex-direction: column;
  gap: 15px;
  overflow-y: auto;
`;

export const InfoCard = styled.div`
  background-color: white;
  border-radius: 8px;
  padding: 15px;
  box-shadow: 0 2px 5px rgba(0,0,0,0.1);
  border-left: 4px solid #2196F3;
`;

// === 공개/비공개 설정 ===
export const SectionHeaderWithPrivacy = styled.div`
  display: flex;
  align-items: center;
  gap: 8px;
`;

export const SectionTitle = styled.span`
  font-weight: bold;
  color: #333;
`;

export const PrivacyToggle = styled.div`
  display: flex;
  align-items: center;
`;

export const ToggleLabel = styled.label`
  display: flex;
  align-items: center;
  cursor: pointer;
  font-size: 11px;
`;

export const ToggleText = styled.span`
  color: #666;
`;

export const PrivacyStatus = styled.div`
  font-size: 11px;
  padding: 2px 6px;
  border-radius: 10px;
  background-color: ${props => props.public ? '#b2ff59' : '#ffebee'};
  color: ${props => props.public ? '#388e3c' : '#F44336'};
`;

// === 연차 정보 ===
export const LeaveItem = styled.div`
  margin-bottom: 8px;
  font-size: 13px;
  display: flex;
  justify-content: space-between;

  &:last-child {
    margin-bottom: 0;
    font-weight: bold;
    color: #4CAF50;
  }
`;

export const LeaveLabel = styled.span`
  color: #555;
`;

export const LeaveValue = styled.span`
  font-weight: bold;
  color: #333;
`;

// === 출퇴근 관리 ===
export const CurrentTime = styled.div`
  font-weight: bold;
  color: #333;
  font-size: 12px;
`;

export const AttendanceStatus = styled.div`
  display: flex;
  flex-direction: column;
  gap: 5px;
`;

export const StatusItem = styled.div`
  display: flex;
  justify-content: space-between;
  padding: 3px 0;
  border-bottom: 1px dotted #ddd;
`;

export const StatusLabel = styled.span`
  color: #666;
`;

export const StatusValue = styled.span`
  color: #333;
  font-weight: bold;
  
  &.recorded {
    color: #555;
    font-weight: 600;
    background-color: #f5f5f5;
    padding: 2px 6px;
    border-radius: 3px;
    border-left: 2px solid #ddd;
  }
`;

export const AttendanceButtons = styled.div`
  display: flex;
  gap: 8px;
  margin-top: 10px;
`;

export const CheckInBtn = styled.button`
  background-color: #4CAF50;
  color: white;
  border: none;
  border-radius: 4px;
  padding: 8px 12px;
  cursor: pointer;
  font-size: 12px;
  font-weight: bold;
  flex: 1;
  transition: all 0.2s;

  &:hover {
    background-color: #45a049;
  }

  &:disabled {
    background-color: #cccccc;
    cursor: not-allowed;
  }

  &.disabled {
    background-color: #cccccc;
    cursor: not-allowed;
  }
`;

export const CheckOutBtn = styled.button`
  background-color: #f44336;
  color: white;
  border: none;
  border-radius: 4px;
  padding: 8px 12px;
  cursor: pointer;
  font-size: 12px;
  font-weight: bold;
  flex: 1;
  transition: all 0.2s;

  &:hover {
    background-color: #da190b;
  }

  &:disabled {
    background-color: #cccccc;
    cursor: not-allowed;
  }

  &.disabled {
    background-color: #cccccc;
    cursor: not-allowed;
  }
`;

// === 급여명세서 ===
export const UploadBtn = styled.button`
  background-color: #FF9800;
  color: white;
  border: none;
  border-radius: 3px;
  padding: 2px 5px;
  font-size: 10px;
  cursor: pointer;
  display: flex;
  align-items: center;
`;

export const UploadIcon = styled.i`
  margin-right: 3px;
  font-weight: bold;
`;

export const SalaryYears = styled.div`
  max-height: 150px;
  overflow-y: auto;
`;

export const SalaryYear = styled.div`
  margin-bottom: 5px;
  border-radius: 4px;
  overflow: hidden;
`;

export const YearHeader = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 4px 6px;
  background-color: #f0f0f0;
  cursor: pointer;
  font-size: 11px;

  ${SalaryYear}.active & {
    background-color: #e3f2fd;
    color: #2196F3;
  }
`;

export const YearActions = styled.div`
  display: flex;
  align-items: center;
  gap: 5px;
`;

export const DownloadAllBtn = styled.button`
  background: none;
  border: none;
  padding: 0;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #2196F3;
  font-size: 11px;
`;

export const DownloadAllIcon = styled.i`
  font-weight: bold;
  font-style: normal;
`;

export const SalaryMonths = styled.div`
  background-color: white;
  border: 1px solid #eee;
  border-top: none;
  max-height: 120px;
  overflow-y: auto;
`;

export const SalaryMonth = styled.div`
  display: flex;
  justify-content: space-between;
  padding: 3px 8px;
  cursor: pointer;
  border-bottom: 1px solid #f0f0f0;
  font-size: 10px;

  &:hover {
    background-color: #f5f5f5;
  }
`;

export const DownloadIcon = styled.i`
  font-weight: bold;
  color: #2196F3;
`;

export const NoStatements = styled.div`
  padding: 8px 10px;
  color: #999;
  font-style: italic;
  text-align: center;
  font-size: 11px;
`;

// === 일정 목록 ===
export const DayScheduleList = styled.div`
  background-color: white;
  border-radius: 8px;
  padding: 15px;
  margin-top: 15px;
  box-shadow: 0 2px 5px rgba(0,0,0,0.1);
`;

export const ScheduleItem = styled.div`
  display: flex;
  align-items: center;
  padding: 10px;
  background-color: ${props => props.selected ? '#e3f2fd' : 'white'};
  border-radius: 4px;
  margin-bottom: 8px;
  cursor: pointer;
  transition: background-color 0.2s;
  border: ${props => props.selected ? '1px solid #2196F3' : 'none'};

  &:hover {
    background-color: #f0f0f0;
  }
`;

const scheduleColorProps = { color: String };
export const ScheduleColor = styled('div', scheduleColorProps)`
  width: 12px;
  height: 12px;
  border-radius: 50%;
  margin-right: 10px;
  background-color: ${props => props.color || '#4CAF50'};
`;

export const ScheduleInfo = styled.div`
  flex: 1;

  .schedule-title {
    font-weight: bold;
    color: #333;
    margin-bottom: 2px;
  }

  .schedule-time {
    font-size: 12px;
    color: #666;
  }
`;

export const ScheduleActions = styled.div`
  display: flex;
  gap: 5px;
`;

export const EditBtn = styled.button`
  padding: 4px 8px;
  border: none;
  border-radius: 3px;
  font-size: 11px;
  cursor: pointer;
  background-color: #2196F3;
  color: white;
`;

export const DeleteBtn = styled.button`
  padding: 4px 8px;
  border: none;
  border-radius: 3px;
  font-size: 11px;
  cursor: pointer;
  background-color: #F44336;
  color: white;
`;

export const AddScheduleBtn = styled.button`
  width: 100%;
  padding: 10px;
  background-color: #4CAF50;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-weight: bold;
  margin-top: 10px;
`;

// === 일정 폼 ===
export const ScheduleFormStyled = styled.div`
  background-color: white;
  border-radius: 8px;
  padding: 20px;
  margin-top: 20px;
  box-shadow: 0 2px 5px rgba(0,0,0,0.1);
`;

export const FormGroup = styled.div`
  margin-bottom: 15px;
`;

export const ColorOptions = styled.div`
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
`;

export const ColorOption = styled.div`
  width: 25px;
  height: 25px;
  border-radius: 50%;
  cursor: pointer;
  border: 2px solid transparent;
  transition: border-color 0.2s;
  border-color: ${props => props.selected ? '#333' : 'transparent'};
`;

export const FormActions = styled.div`
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 20px;
`;

export const CancelBtn = styled.button`
  background-color: #f0f0f0;
  color: #333;
  padding: 8px 16px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
`;

export const SaveBtn = styled.button`
  background-color: #4CAF50;
  color: white;
  padding: 8px 16px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
`;

// === 달력 ===
export const CalendarSectionStyle = styled.div`
  flex: 0 0 780px;
  display: flex;
  flex-direction: column;
  background-color: white;
  border-radius: 26px;
  padding: 24px;
  box-shadow: 0 2px 5px rgba(0,0,0,0.1);
  min-width: 442px;
  max-width: 780px;
  width: 100%;
  height: auto;
  align-items: center;
`;

export const CalendarHeader = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
`;

export const CalendarGrid = styled.div`
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 5px;
  background-color: #f5f5f5;
  padding: 10px;
  border-radius: 16px;
  width: 100%;
  max-width: 100%;
  height: auto;
  box-sizing: border-box;
`;

export const WeekdayHeader = styled.div`
  text-align: center;
  font-weight: bold;
  padding: 8px 6px;
  background-color: #e3f2fd;
  color: #2196F3;
  font-size: 12px;
`;

export const CalendarDay = styled.div`
  border: 2px solid #e0e0e0;
  min-height: 73px;
  aspect-ratio: 1;
  padding: 8px;
  background-color: ${props => props.otherMonth ? '#f9f9f9' : 'white'};
  color: ${props => props.otherMonth ? '#aaa' : '#333'};
  cursor: pointer;
  transition: all 0.2s ease;
  font-size: 21px;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  justify-content: flex-start;
  box-sizing: border-box;

  &:hover {
    background-color: #f0f8ff;
    border-color: #2196F3;
    transform: translateY(-1px);
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  }

  &.selected-day {
    background-color: #e3f2fd;
    border: 3px solid #2196F3;
    box-shadow: 0 2px 8px rgba(33, 150, 243, 0.3);
  }
`;

export const DayNumber = styled.div`
  font-weight: bold;
  margin-bottom: 3px;
  color: #333;
  font-size: 12px;
`;

export const DaySchedules = styled.div`
  display: flex;
  flex-direction: column;
  gap: 2px;
`;

const dayScheduleItemProps = { color: String };
export const DayScheduleItem = styled('div', dayScheduleItemProps)`
  padding: 1px 3px;
  border-radius: 2px;
  font-size: 9px;
  color: white;
  cursor: pointer;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin-bottom: 1px;
  background-color: ${props => props.color || '#4CAF50'};

  &:hover {
    opacity: 0.8;
  }
`;

// === 드롭다운 ===
export const DropdownContainer = styled.div`
  position: relative;
  min-width: 80px;
`;

export const DropdownToggle = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 12px;
  background-color: #f5f5f5;
  border: 1px solid #ddd;
  border-radius: 4px;
  cursor: pointer;
  font-weight: bold;
  font-size: 14px;

  &:hover {
    background-color: #e9e9e9;
  }
`;

export const DropdownArrow = styled.i`
  display: inline-block;
  width: 0;
  height: 0;
  margin-left: 8px;
  border-left: 5px solid transparent;
  border-right: 5px solid transparent;
  border-top: 5px solid #333;
  transition: transform 0.2s;
  transform: ${props => props.open ? 'rotate(180deg)' : 'rotate(0deg)'};
`;

export const DropdownMenu = styled.div`
  position: absolute;
  top: 100%;
  left: 0;
  width: 100%;
  max-height: 200px;
  overflow-y: auto;
  background-color: white;
  border: 1px solid #ddd;
  border-radius: 4px;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
  z-index: 10;
  margin-top: 2px;
`;

export const DropdownItem = styled.div`
  padding: 8px 12px;
  cursor: pointer;
  font-size: 14px;
  background-color: ${props => props.selected ? '#e3f2fd' : 'white'};
  font-weight: ${props => props.selected ? 'bold' : 'normal'};
  color: ${props => props.selected ? '#2196F3' : '#333'};

  &:hover {
    background-color: #f5f5f5;
  }
`;

// === 반응형 디자인 ===
// @media (max-width: 1200px) {
//   .schedule-management {
//     flex-direction: column;
//     gap: 15px;
//   }
//   
//   .clock-section {
//     flex: none;
//     width: 100%;
//   }
//   
//   .info-card-section {
//     flex: none;
//     width: 100%;
//     flex-direction: row;
//     flex-wrap: wrap;
//   }
//   
//   .info-card {
//     flex: 1;
//     min-width: 250px;
//   }
//   
//   .calendar-section {
//     min-width: auto;
//   }
// }

// @media (max-width: 768px) {
//   .content-area {
//     padding: 10px;
//   }
//   
//   .header-container {
//     flex-direction: column;
//     gap: 10px;
//     align-items: flex-start;
//   }
//   
//   .clock-face {
//     width: 220px;
//     height: 220px;
//   }
//   
//   .calendar-grid {
//     gap: 1px;
//     padding: 5px;
//   }
//   
//   .calendar-day {
//     min-height: 80px;
//     padding: 5px;
//   }
// } 
import styled from '@vue-styled-components/core';

export const ScheduleSharingRoot = styled.div`
  display: flex;
  height: 100vh;
  width: 100%;
`;

export const ContentArea = styled.div`
  flex: 1;
  padding: 100px 20px 20px 20px;
  overflow-y: auto;
`;

export const DashboardContainer = styled.div`
  width: 100%;
  height: calc(100vh - 140px);
  position: relative;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  background-color: #f9f9f9;
  margin-top: 30px;
  padding-top: 50px;
  padding-bottom: 50px;
`;

export const GraphView = styled.div`
  width: 100%;
  height: 100%;
  position: relative;
  overflow: hidden;
`;

export const DepartmentBoundary = styled.div`
  position: absolute;
  border: 2px solid #aaa;
  border-radius: 15px;
  background-color: rgba(245, 245, 245, 0.5);
  pointer-events: none;
  z-index: 0;
  box-shadow: 0 3px 10px rgba(0,0,0,0.1);
`;

export const DepartmentLabel = styled.div`
  position: absolute;
  top: 15px;
  left: 25px;
  font-size: 18px;
  font-weight: bold;
  color: #555;
`;

export const MemberNode = styled.div`
  position: absolute;
  width: 220px;
  height: 200px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: flex-start;
  padding: 6px;
  cursor: move;
  transition: all 0.3s ease;
  transform: translate(-50%, -50%);
  z-index: 1;
  &:hover {
    transform: translate(-50%, -50%) scale(1.05);
  }
  &.selected {
    outline: 3px solid #4CAF50;
    outline-offset: 3px;
  }
`;

export const MyNode = styled(MemberNode)`
  border: 2px solid #4682b4;
  background-color: rgba(240, 248, 255, 0.5);
  border-radius: 8px;
  z-index: 2;
`;

export const MemberInfo = styled.div`
  text-align: center;
  margin-bottom: 8px;
`;

export const MemberPosition = styled.div`
  font-size: 13px;
  color: #666;
  margin-bottom: 2px;
  line-height: 1.2;
`;

export const MemberName = styled.div`
  font-size: 14px;
  font-weight: bold;
  margin-bottom: 4px;
  line-height: 1.2;
`;

export const ScheduleViews = styled.div`
  display: flex;
  gap: 10px;
`;

export const MiniClockView = styled.div`
  cursor: pointer;
`;

export const MiniClock = styled.div`
  width: 80px;
  height: 80px;
  position: relative;
  box-sizing: border-box;
`;

export const ClockFace = styled.div`
  width: 100%;
  height: 100%;
  border-radius: 50%;
  border: 2px solid #ccc;
  position: relative;
  background-color: white;
  box-sizing: border-box;
`;

export const HourNumber = styled.div`
  position: absolute;
  font-size: 9px;
  color: #333;
  font-weight: bold;
  z-index: 2;
`;

export const ScheduleSegment = styled.div`
  position: absolute;
  width: 100%;
  height: 100%;
  left: 0;
  top: 0;
  z-index: 1;
  opacity: 0.7;
  border-radius: 50%;
`;

export const MiniCalendarView = styled.div`
  cursor: pointer;
`;

export const MiniCalendar = styled.div`
  width: 80px;
  height: 80px;
  border: 2px solid #ccc;
  background-color: white;
`;

export const MiniCalendarGrid = styled.div`
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  grid-template-rows: repeat(6, 1fr);
  width: 100%;
  height: 100%;
  font-size: 7px;
`;

export const MiniDayHeader = styled.div`
  text-align: center;
  font-size: 7px;
  font-weight: bold;
  color: #666;
  border-bottom: 1px solid #eee;
`;

export const MiniDayCell = styled.div`
  position: relative;
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
  overflow: hidden;
`;

export const MiniDayNumber = styled.div`
  font-size: 7px;
  z-index: 20;
  color: #333;
  position: relative;
  margin-bottom: 3px;
  text-shadow: 0 0 2px white;
  font-weight: bold;
`;

export const MiniDayEvent = styled.div`
  position: absolute;
  width: 100%;
  height: 100%;
  opacity: 0.5;
  z-index: 1;
`;

export const MiniProjectLine = styled.div`
  position: absolute;
  bottom: 0px;
  left: 0;
  right: 0;
  height: 1.3px;
  transform: none;
  opacity: 1;
  z-index: 1;
  border-radius: 0px;
  box-shadow: none;
`;

// 모달 관련 스타일
export const ExpandedScheduleModal = styled.div`
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
`;

export const ModalContent = styled.div`
  background: white;
  border-radius: 8px;
  width: 95%;
  max-width: 1600px;
  height: 90%;
  max-height: 1000px;
  display: flex;
  flex-direction: column;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.3);
`;

export const ModalHeader = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid #e0e0e0;
  
  h2 {
    margin: 0;
    color: #333;
  }
`;

export const CloseBtn = styled.button`
  background: none;
  border: none;
  font-size: 24px;
  cursor: pointer;
  color: #666;
  
  &:hover {
    color: #333;
  }
`;

export const ModalBody = styled.div`
  flex: 1;
  padding: 20px;
  overflow-y: auto;
`;

export const ExpandedViews = styled.div`
  display: flex;
  gap: 20px;
  height: 100%;
`;

export const ExpandedClockContainer = styled.div`
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 20px;
`;

export const ExpandedClock = styled.div`
  width: 500px;
  height: 500px;
  margin: 0 auto;
  position: relative;
`;

export const ExpandedClockFace = styled.div`
  width: 100%;
  height: 100%;
  border-radius: 50%;
  border: 3px solid #ccc;
  position: relative;
  background-color: white;
`;

export const ExpandedHourNumber = styled.div`
  position: absolute;
  font-size: 24px;
  color: #333;
  font-weight: bold;
  z-index: 2;
`;

export const ExpandedScheduleSegment = styled.div`
  position: absolute;
  width: 100%;
  height: 100%;
  left: 0;
  top: 0;
  z-index: 1;
  opacity: 0.8;
  border-radius: 50%;
  cursor: pointer;
  
  &:hover {
    opacity: 1;
  }
`;

export const ScheduleTooltip = styled.div`
  position: absolute;
  background: rgba(0, 0, 0, 0.8);
  color: white;
  padding: 8px;
  border-radius: 4px;
  font-size: 12px;
  z-index: 10;
  pointer-events: none;
  white-space: nowrap;
`;

export const TodaySchedules = styled.div`
  flex: 1;
  background: #f9f9f9;
  border-radius: 8px;
  padding: 15px;
  overflow-y: auto;
`;

export const TodayScheduleItem = styled.div`
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.2s;
  
  &:hover {
    background-color: #e9e9e9;
  }
  
  &.selected {
    background-color: #e3f2fd;
    border: 1px solid #2196f3;
  }
`;

export const ScheduleColorDot = styled.div`
  width: 12px;
  height: 12px;
  border-radius: 50%;
  flex-shrink: 0;
`;

export const ScheduleInfo = styled.div`
  flex: 1;
`;

export const ScheduleTimeRange = styled.div`
  font-size: 12px;
  color: #666;
  margin-bottom: 2px;
`;

export const ScheduleTitle = styled.div`
  font-size: 14px;
  font-weight: 500;
  color: #333;
`;

export const NoSchedules = styled.div`
  text-align: center;
  color: #666;
  font-style: italic;
  padding: 20px;
`;

export const AdditionalInfoContainer = styled.div`
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 15px;
`;

export const InfoSection = styled.div`
  background: white;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  padding: 15px;
`;

export const SectionHeaderWithPrivacy = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
`;

export const SectionTitle = styled.span`
  font-weight: bold;
  color: #333;
`;

export const PrivacyStatus = styled.div`
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 11px;
  font-weight: 500;
  
  &.public {
    background-color: #e8f5e8;
    color: #2e7d32;
  }
  
  &.private {
    background-color: #ffebee;
    color: #c62828;
  }
`;

export const LeaveItem = styled.div`
  display: flex;
  justify-content: space-between;
  margin-bottom: 5px;
`;

export const LeaveLabel = styled.span`
  color: #666;
  font-size: 13px;
`;

export const LeaveValue = styled.span`
  font-weight: 500;
  color: #333;
`;

export const AttendanceStatus = styled.div`
  display: flex;
  flex-direction: column;
  gap: 8px;
`;

export const StatusItem = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
`;

export const StatusLabel = styled.span`
  color: #666;
  font-size: 13px;
`;

export const StatusValue = styled.span`
  font-weight: 500;
  color: #333;
`;

// 달력 관련 스타일
export const ExpandedCalendarContainer = styled.div`
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 25px;
  min-height: 700px;
`;

export const CalendarHeader = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
  
  button {
    background: none;
    border: none;
    font-size: 18px;
    cursor: pointer;
    color: #666;
    
    &:hover {
      color: #333;
    }
  }
  
  h3 {
    margin: 0;
    color: #333;
  }
`;

export const CalendarGrid = styled.div`
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 2px;
  background: white;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  padding: 10px;
`;

export const WeekdayHeader = styled.div`
  text-align: center;
  font-weight: bold;
  color: #666;
  padding: 15px;
  font-size: 16px;
`;

export const CalendarDay = styled.div`
  position: relative;
  min-height: 100px;
  padding: 10px;
  border: 1px solid #f0f0f0;
  cursor: pointer;
  transition: background-color 0.2s;
  
  &:hover {
    background-color: #f9f9f9;
  }
  
  &.other-month {
    background-color: #fafafa;
    color: #ccc;
  }
  
  &.selected-date {
    background-color: #e3f2fd;
    border-color: #2196f3;
  }
  
  &.project-active {
    border-bottom: 3px solid #B4D3F5;
  }
`;

export const DayNumber = styled.div`
  font-size: 16px;
  font-weight: 500;
  margin-bottom: 6px;
`;

export const ProjectLine = styled.div`
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  height: 3px;
`;

export const DaySchedules = styled.div`
  display: flex;
  flex-direction: column;
  gap: 1px;
`;

export const DayScheduleItem = styled.div`
  font-size: 13px;
  padding: 3px 6px;
  border-radius: 4px;
  color: white;
  cursor: pointer;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  
  &:hover {
    opacity: 0.8;
  }
`;

export const WeekView = styled.div`
  background: white;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  padding: 10px;
`;

export const WeekHeader = styled.div`
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 2px;
  margin-bottom: 10px;
`;

export const WeekDayHeader = styled.div`
  text-align: center;
  font-weight: bold;
  color: #666;
  font-size: 15px;
  padding: 8px;
`;

export const WeekBody = styled.div`
  display: flex;
  flex-direction: column;
  gap: 2px;
`;

export const WeekRow = styled.div`
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 2px;
  min-height: 20px;
`;

export const WeekDay = styled.div`
  position: relative;
  border: 1px solid #f0f0f0;
  min-height: 40px;
`;

export const WeekScheduleLine = styled.div`
  position: absolute;
  left: 0;
  right: 0;
  height: 4px;
  cursor: pointer;
  
  &:hover {
    opacity: 0.8;
  }
`;

// 급여명세서 관련 스타일
export const SalaryYears = styled.div`
  display: flex;
  flex-direction: column;
  gap: 10px;
`;

export const SalaryYear = styled.div`
  border: 1px solid #e0e0e0;
  border-radius: 6px;
  padding: 10px;
  
  &.active {
    border-color: #2196f3;
    background-color: #f8f9ff;
  }
`;

export const YearHeader = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
  
  span {
    font-weight: bold;
    color: #333;
  }
`;

export const YearActions = styled.div`
  display: flex;
  gap: 5px;
`;

export const DownloadAllBtn = styled.button`
  background: none;
  border: none;
  cursor: pointer;
  color: #666;
  font-size: 12px;
  
  &:hover {
    color: #333;
  }
`;

export const DownloadAllIcon = styled.i`
  font-style: normal;
`;

export const SalaryMonths = styled.div`
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
`;

export const SalaryMonth = styled.div`
  display: flex;
  align-items: center;
  gap: 5px;
  padding: 5px 8px;
  border: 1px solid #e0e0e0;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.2s;
  
  &:hover {
    background-color: #f5f5f5;
    border-color: #ccc;
  }
  
  span {
    font-size: 12px;
    color: #333;
  }
`;

export const DownloadIcon = styled.i`
  font-style: normal;
  color: #666;
  font-size: 10px;
`;

// 기타 모달 요소들
export const ScheduleDetails = styled.div`
  background: white;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  padding: 15px;
  margin-top: 15px;
  
  h3 {
    margin: 0 0 10px 0;
    color: #333;
  }
  
  p {
    margin: 0;
    color: #666;
  }
`;

export const SelectedDateInfo = styled.div`
  background: white;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  padding: 15px;
  margin-top: 15px;
  
  h3 {
    margin: 0 0 10px 0;
    color: #333;
  }
`;

export const DateSchedules = styled.div`
  display: flex;
  flex-direction: column;
  gap: 8px;
`;

export const DateScheduleItem = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px;
  border: 1px solid #e0e0e0;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.2s;
  
  &:hover {
    background-color: #f9f9f9;
  }
  
  &.selected {
    background-color: #e3f2fd;
    border-color: #2196f3;
  }
`;

export const ScheduleTime = styled.div`
  font-size: 12px;
  color: #666;
`;

export const ScheduleActions = styled.div`
  display: flex;
  gap: 10px;
  margin-top: 15px;
`;

export const AddBtn = styled.button`
  background-color: #4CAF50;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 12px;
  
  &:hover {
    background-color: #45a049;
  }
`;

export const EditBtn = styled.button`
  background-color: #2196F3;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 12px;
  
  &:hover {
    background-color: #1976D2;
  }
  
  &:disabled {
    background-color: #ccc;
    cursor: not-allowed;
  }
`;

export const ModalBottomProjects = styled.div`
  background: white;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  padding: 15px;
  margin-top: 15px;
  
  h4 {
    margin: 0 0 10px 0;
    color: #333;
  }
`;

export const ProjectList = styled.div`
  display: flex;
  flex-direction: column;
  gap: 8px;
`;

export const ProjectInfo = styled.div`
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px;
  border: 1px solid #e0e0e0;
  border-radius: 4px;
`;

export const ProjectLinePreview = styled.div`
  width: 20px;
  height: 3px;
  border-radius: 2px;
`;

export const ProjectDetails = styled.div`
  flex: 1;
`;

export const ProjectName = styled.div`
  font-size: 12px;
  font-weight: 500;
  color: #333;
`;

export const ProjectPeriod = styled.div`
  font-size: 10px;
  color: #666;
`; 
import styled from '@vue-styled-components/core';

export const CommuteRoot = styled.div`
  display: flex;
  min-height: 100vh;
  background: #f8f9fa;
`;

export const MainContent = styled.div`
  flex: 1;
  padding: 30px;
  max-width: 1200px;
  margin: 0 auto;
  @media (max-width: 768px) {
    padding: 20px;
  }
`;

export const Title = styled.h1`
  font-size: 28px;
  font-weight: 600;
  color: #333;
  margin-bottom: 30px;
  @media (max-width: 768px) {
    font-size: 24px;
  }
`;

export const Card = styled.div`
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
  padding: 25px;
  margin-bottom: 25px;
`;

export const CardHeader = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
`;

export const CardTitle = styled.h2`
  font-size: 20px;
  font-weight: 600;
  color: #333;
  margin: 0;
`;

export const CorrectionRequestBtn = styled.button`
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  background: #007bff;
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.2s ease;
  &:hover {
    background: #0056b3;
  }
`;

export const TimeDisplay = styled.div`
  display: flex;
  gap: 30px;
  margin-bottom: 25px;
  width: 100%;
  max-width: 820px;
  margin: 0 auto;
  @media (max-width: 768px) {
    flex-direction: column;
    gap: 15px;
  }
`;

export const TimeItem = styled.div`
  text-align: center;
  flex: 1;
  padding: 0 24px;
  background: #f8f9fa;
  border-radius: 8px;
  border: 1px solid #e9ecef;
  width: 100%;
  max-width: 400px;
  margin: 0 auto;
`;

export const TimeLabel = styled.div`
  font-size: 14px;
  color: #6c757d;
  margin-bottom: 8px;
  font-weight: 500;
`;

export const TimeValue = styled.div`
  font-size: 24px;
  font-weight: 600;
  color: #333;
`;

export const ButtonGroup = styled.div`
  display: flex;
  gap: 20px;
  width: 100%;
  max-width: 820px;
  margin: 0 auto 20px auto;
  @media (max-width: 768px) {
    flex-direction: column;
  }
`;

export const CheckInBtn = styled.button`
  flex: 1;
  padding: 10px 0;
  border: none;
  border-radius: 6px;
  font-size: 15px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
  width: 100%;
  max-width: 400px;
  margin: 0 auto;
  background: #28a745;
  color: white;
  &:hover:not(:disabled) {
    background: #218838;
  }
  &:disabled {
    background: #6c757d;
    color: #fff;
    cursor: not-allowed;
  }
  &.loading {
    background: #ffc107;
    color: #212529;
  }
`;

export const CheckOutBtn = styled(CheckInBtn)`
  background: #dc3545;
  color: white;
  &:hover:not(:disabled) {
    background: #c82333;
  }
`;

export const WeekGrid = styled.div`
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 12px;
  width: 100%;
  @media (max-width: 768px) {
    grid-template-columns: repeat(4, 1fr);
  }
`;

export const WeekDay = styled.div`
  padding: 15px 10px;
  border-radius: 6px;
  text-align: center;
  background: white;
  border: 1px solid #e9ecef;
  transition: border-color 0.2s ease;
  &.today {
    border-color: #007bff;
    background: #f8f9ff;
  }
  &.weekend {
    border-color: #ffc107;
    background: #fffbf0;
  }
`;

export const DayName = styled.div`
  font-size: 14px;
  font-weight: 600;
  color: #333;
  margin-bottom: 5px;
`;

export const DayDate = styled.div`
  font-size: 12px;
  color: #6c757d;
  margin-bottom: 8px;
`;

export const DayStatus = styled.div`
  font-size: 11px;
  font-weight: 500;
  padding: 4px 8px;
  border-radius: 12px;
  display: inline-block;
  &.checked-in {
    background: #d4edda;
    color: #155724;
  }
  &.checked-out {
    background: #d1ecf1;
    color: #0c5460;
  }
  &.absent {
    background: #f8d7da;
    color: #721c24;
  }
  &.weekend {
    background: #fff3cd;
    color: #856404;
  }
`;

export const ModalOverlay = styled.div`
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
`;

export const ModalContent = styled.div`
  background: white;
  border-radius: 8px;
  width: 90%;
  max-width: 500px;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
  @media (max-width: 768px) {
    width: 95%;
    margin: 20px;
  }
`;

export const ModalHeader = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 25px;
  border-bottom: 1px solid #e9ecef;
  @media (max-width: 768px) {
    padding: 15px 20px;
  }
`;

export const ModalBody = styled.div`
  padding: 25px;
  @media (max-width: 768px) {
    padding: 15px 20px;
  }
`;

export const ModalFooter = styled.div`
  display: flex;
  gap: 10px;
  padding: 20px 25px;
  border-top: 1px solid #e9ecef;
  justify-content: flex-end;
  @media (max-width: 768px) {
    padding: 15px 20px;
  }
`;

export const ModalClose = styled.button`
  background: none;
  border: none;
  font-size: 18px;
  color: #6c757d;
  cursor: pointer;
  padding: 4px;
  border-radius: 4px;
  transition: background-color 0.2s ease;
  &:hover {
    background: #f8f9fa;
  }
`;

export const FormGroup = styled.div`
  margin-bottom: 20px;
`;

export const FormInput = styled.input`
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #ced4da;
  border-radius: 4px;
  font-size: 14px;
  transition: border-color 0.2s ease;
  &:focus {
    outline: none;
    border-color: #007bff;
    box-shadow: 0 0 0 2px rgba(0, 123, 255, 0.25);
  }
`;

export const RadioGroup = styled.div`
  display: flex;
  gap: 15px;
  flex-wrap: wrap;
  @media (max-width: 768px) {
    flex-direction: column;
  }
`;

export const RadioItem = styled.label`
  display: flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  padding: 8px 12px;
  border: 1px solid #ced4da;
  border-radius: 4px;
  transition: border-color 0.2s ease;
  background: #f8f9fa;
  &:hover {
    border-color: #007bff;
  }
`;

export const TimeInputs = styled.div`
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
  @media (max-width: 768px) {
    grid-template-columns: 1fr;
  }
`;

export const TimeInputGroup = styled.div`
  display: flex;
  flex-direction: column;
  gap: 4px;
`;

export const FormTextarea = styled.textarea`
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #ced4da;
  border-radius: 4px;
  font-size: 14px;
  min-height: 80px;
  resize: vertical;
  font-family: inherit;
  transition: border-color 0.2s ease;
  &:focus {
    outline: none;
    border-color: #007bff;
    box-shadow: 0 0 0 2px rgba(0, 123, 255, 0.25);
  }
`;

export const BtnSecondary = styled.button`
  padding: 8px 16px;
  border: none;
  border-radius: 4px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.2s ease;
  background: #6c757d;
  color: white;
  &:hover {
    background: #5a6268;
  }
`;

export const BtnPrimary = styled.button`
  padding: 8px 16px;
  border: none;
  border-radius: 4px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.2s ease;
  background: #007bff;
  color: white;
  &:hover:not(:disabled) {
    background: #0056b3;
  }
  &:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }
`;

export const HistoryTable = styled.div`
  overflow-x: auto;
`;

export const Table = styled.table`
  width: 100%;
  border-collapse: collapse;
  background: white;
  border-radius: 6px;
  overflow: hidden;
  border: 1px solid #e9ecef;
`;

export const Th = styled.th`
  background: #f8f9fa;
  font-weight: 600;
  color: #333;
  font-size: 14px;
  padding: 12px;
`;

export const Td = styled.td`
  font-size: 14px;
  color: #333;
  padding: 12px;
`;

export const Tr = styled.tr`
  &:hover {
    background: #f8f9fa;
  }
`;

export const StatusBadge = styled.span`
  padding: 4px 8px;
  border-radius: 12px;
  font-size: 11px;
  font-weight: 500;
  display: inline-block;
  &.normal {
    background: #d4edda;
    color: #155724;
  }
  &.overtime {
    background: #d1ecf1;
    color: #0c5460;
  }
  &.late {
    background: #f8d7da;
    color: #721c24;
  }
  &.absent {
    background: #f8d7da;
    color: #721c24;
  }
`; 
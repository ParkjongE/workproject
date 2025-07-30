import styled from '@vue-styled-components/core';

export const Container = styled.div`
  max-width: 1200px;
  margin: 30px auto;
  background: #fff;
  border: 1px solid #ddd;
  border-radius: 8px;
  padding: 0;
  box-shadow: 0 2px 6px rgba(0,0,0,0.1);
  position: relative;
  display: flex;
  min-height: 700px;
  @media (max-width: 1200px) {
    flex-direction: column;
    max-width: 100vw;
  }
`;

export const MainContent = styled.div`
  flex: 1;
  padding: 40px 40px 40px 40px;
`;

export const AttendanceTitle = styled.div`
  font-size: 24px;
  font-weight: bold;
  margin-bottom: 18px;
  color: #218838;
`;

export const FilterArea = styled.div`
  display: flex;
  gap: 18px;
  align-items: center;
  margin-bottom: 24px;
  & label {
    font-size: 15px;
    font-weight: 500;
  }
  & select, & input[type="date"] {
    font-size: 15px;
    padding: 4px 8px;
    border: 1px solid #bbb;
    border-radius: 6px;
  }
`;

export const SummaryBox = styled.div`
  background: #f5f5f5;
  border-radius: 8px;
  padding: 16px 24px;
  margin-bottom: 24px;
  display: flex;
  gap: 32px;
  font-size: 16px;
  color: #234fa0;
  font-weight: 500;
`;

export const HighlightGreen = styled.span`
  color: #218838;
  font-weight: bold;
`;
export const HighlightRed = styled.span`
  color: #e74c3c;
  font-weight: bold;
`;
export const HighlightBlue = styled.span`
  color: #234fa0;
  font-weight: bold;
`;

export const AnnualLeaveInfo = styled.div`
  margin-bottom: 12px;
  color: #218838;
  font-size: 16px;
  font-weight: 500;
`;
export const PeriodInfo = styled.span`
  color: #888;
  font-size: 15px;
`;

export const StyledTable = styled.table`
  width: 100%;
  border-collapse: collapse;
  font-size: 15px;
  margin-top: 8px;
`;
export const Th = styled.th`
  border: 1px solid #eee;
  padding: 8px 6px;
  text-align: center;
  background: #f5f5f5;
`;
export const Td = styled.td`
  border: 1px solid #eee;
  padding: 8px 6px;
  text-align: center;
`;
export const Tr = styled.tr`
  &:nth-child(even) {
    background: #fafafa;
  }
`; 
import styled from '@vue-styled-components/core';

export const Container = styled.div`
  max-width: 1100px;
  margin: 40px auto;
  background: #fff;
  border-radius: 10px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.08);
  padding: 0;
  display: flex;
  @media (max-width: 900px) {
    flex-direction: column;
    max-width: 100vw;
  }
`;

export const MainContent = styled.div`
  flex: 1;
  padding: 40px 40px 40px 40px;
`;

export const Title = styled.h2`
  font-size: 24px;
  font-weight: bold;
  margin-bottom: 18px;
  color: #218838;
`;

export const StyledForm = styled.form`
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin-bottom: 32px;
`;

export const StyledLabel = styled.label`
  font-size: 16px;
  font-weight: 500;
  margin-bottom: 4px;
`;

export const StyledSelect = styled.select`
  font-size: 15px;
  padding: 6px 10px;
  border: 1px solid #bbb;
  border-radius: 6px;
`;

export const StyledInput = styled.input`
  font-size: 15px;
  padding: 6px 10px;
  border: 1px solid #bbb;
  border-radius: 6px;
`;

export const StyledTextarea = styled.textarea`
  font-size: 15px;
  padding: 6px 10px;
  border: 1px solid #bbb;
  border-radius: 6px;
  resize: vertical;
  min-height: 60px;
`;

export const SubmitButton = styled.button`
  background: #218838;
  color: #fff;
  border: none;
  border-radius: 8px;
  padding: 12px 0;
  font-size: 18px;
  font-weight: bold;
  cursor: pointer;
  margin-top: 8px;
`;

export const UploadRow = styled.div`
  display: flex;
  gap: 32px;
  margin-bottom: 28px;
`;

export const UploadItem = styled.div`
  display: flex;
  flex-direction: column;
  flex: 1;
`;

export const AnnualInfo = styled.span`
  margin-left: 12px;
  color: #218838;
  font-size: 15px;
`;

export const History = styled.div`
  margin-top: 30px;
`;

export const HistoryTitle = styled.div`
  font-size: 18px;
  font-weight: bold;
  margin-bottom: 10px;
  color: #234fa0;
`;

export const StyledTable = styled.table`
  width: 100%;
  border-collapse: collapse;
  font-size: 15px;
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
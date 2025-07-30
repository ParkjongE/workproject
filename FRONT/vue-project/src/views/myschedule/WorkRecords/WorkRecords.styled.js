import styled from '@vue-styled-components/core';

export const WorkRecordsContainer = styled.div`
  display: flex;
  height: 100vh;
  width: 100%;
`;

export const ContentArea = styled.div`
  flex: 1;
  padding: 20px;
  overflow-y: auto;
`;

export const FilterSection = styled.div`
  display: flex;
  flex-wrap: wrap;
  gap: 15px;
  margin-bottom: 20px;
  padding: 15px;
  background-color: #f9f9f9;
  border-radius: 8px;
  align-items: center;
`;

export const DateFilter = styled.div`
  display: flex;
  align-items: center;
  gap: 8px;
`;

export const CategoryFilter = styled.div`
  display: flex;
  align-items: center;
  gap: 8px;
`;

export const SearchBox = styled.div`
  display: flex;
  gap: 5px;
  margin-left: auto;
`;

export const StatsSection = styled.div`
  display: flex;
  gap: 20px;
  margin-bottom: 30px;
`;

export const StatCard = styled.div`
  flex: 1;
  padding: 20px;
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  text-align: center;
`;

export const StatTitle = styled.div`
  font-size: 14px;
  color: #666;
  margin-bottom: 10px;
`;

export const StatValue = styled.div`
  font-size: 24px;
  font-weight: bold;
  color: #333;
`;

export const RecordsSection = styled.div`
  margin-top: 20px;
`;

export const SectionHeader = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
`;

export const AddRecordBtn = styled.button`
  padding: 8px 15px;
  background-color: #4CAF50;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
`;

export const EmptyState = styled.div`
  padding: 40px;
  text-align: center;
  background-color: #f9f9f9;
  border-radius: 8px;
  color: #666;
`;

export const RecordsList = styled.div`
  display: flex;
  flex-direction: column;
  gap: 15px;
`;

export const RecordCard = styled.div`
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  overflow: hidden;
  border-left: ${props => props.completed ? '4px solid #4CAF50' : 'none'};
`;

export const RecordHeader = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px;
  background-color: #f9f9f9;
  border-bottom: 1px solid #e0e0e0;
`;

export const RecordTitle = styled.div`
  font-weight: bold;
  font-size: 16px;
`;

export const RecordStatus = styled.div`
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
  background-color: ${props => props.status === 'completed' ? '#E8F5E9' : '#E3F2FD'};
  color: ${props => props.status === 'completed' ? '#388E3C' : '#1976D2'};
`;

export const RecordDetails = styled.div`
  display: flex;
  flex-wrap: wrap;
  gap: 15px;
  padding: 15px;
  background-color: #fafafa;
`;

export const RecordDate = styled.div`
  font-size: 14px;
  color: #555;
`;

export const RecordTime = styled.div`
  font-size: 14px;
  color: #555;
`;

export const RecordDuration = styled.div`
  font-size: 14px;
  color: #555;
`;

export const RecordCategory = styled.div`
  font-size: 14px;
  color: #555;
  padding: 3px 8px;
  background-color: #f0f0f0;
  border-radius: 4px;
`;

export const RecordDescription = styled.div`
  padding: 15px;
  border-top: 1px solid #f0f0f0;
  font-size: 14px;
  color: #333;
`;

export const RecordActions = styled.div`
  display: flex;
  gap: 10px;
  padding: 15px;
  border-top: 1px solid #f0f0f0;
  background-color: #fafafa;
`;

export const EditBtn = styled.button`
  padding: 6px 12px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 12px;
  background-color: #2196F3;
  color: white;
`;

export const DeleteBtn = styled.button`
  padding: 6px 12px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 12px;
  background-color: #F44336;
  color: white;
`;

export const CompleteBtn = styled.button`
  padding: 6px 12px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 12px;
  background-color: #4CAF50;
  color: white;
`;

export const ModalOverlay = styled.div`
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
`;

export const ModalContent = styled.div`
  width: 600px;
  max-width: 90%;
  background-color: white;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
`;

export const ModalHeader = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px 20px;
  background-color: #f9f9f9;
  border-bottom: 1px solid #e0e0e0;
`;

export const CloseBtn = styled.button`
  background: none;
  border: none;
  font-size: 24px;
  cursor: pointer;
  color: #666;
`;

export const ModalBody = styled.div`
  padding: 20px;
`;

export const FormGroup = styled.div`
  margin-bottom: 15px;
`;

export const FormRow = styled.div`
  display: flex;
  gap: 15px;
`;

export const FormActions = styled.div`
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 20px;
`;

export const CancelBtn = styled.button`
  background-color: #e0e0e0;
  color: #333;
  padding: 8px 15px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-weight: 500;
`;

export const SaveBtn = styled.button`
  background-color: #4CAF50;
  color: white;
  padding: 8px 15px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-weight: 500;
`; 
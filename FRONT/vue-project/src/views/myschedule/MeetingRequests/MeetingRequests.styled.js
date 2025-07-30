import styled from '@vue-styled-components/core';

// 상태별 카드 border-left
export const RequestCard = styled.div`
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  overflow: hidden;
  border-left: 4px solid #eee;
  &.pending { border-left: 4px solid #FFC107; }
  &.accepted { border-left: 4px solid #4CAF50; }
  &.rejected { border-left: 4px solid #F44336; }
`;

// 상태별 RequestStatus
const statusProps = { status: String };
export const RequestStatus = styled('div', statusProps)`
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
  background-color: ${props =>
    props.status === 'pending' ? '#FFF8E1' :
    props.status === 'accepted' ? '#E8F5E9' :
    props.status === 'rejected' ? '#FFEBEE' : '#f9f9f9'};
  color: ${props =>
    props.status === 'pending' ? '#FFA000' :
    props.status === 'accepted' ? '#388E3C' :
    props.status === 'rejected' ? '#D32F2F' : '#666'};
`;

export const Container = styled.div`
  display: flex;
  flex-direction: row;
  height: 100vh;
  width: 100vw;
  background: #f8f9fa;
`;

export const ContentArea = styled.div`
  flex: 1;
  min-width: 0;
  overflow-y: auto;
  padding: 20px;
`;

export const Title = styled.h1`
  margin-bottom: 20px;
  color: #333;
`;

export const Tabs = styled.div`
  display: flex;
  margin-bottom: 20px;
  border-bottom: 1px solid #e0e0e0;
`;

export const Tab = styled.div`
  padding: 10px 20px;
  cursor: pointer;
  border-bottom: 2px solid transparent;
  font-weight: 500;
  &.active {
    border-bottom: 2px solid #4CAF50;
    color: #4CAF50;
  }
`;

export const RequestList = styled.div`
  margin-bottom: 20px;
`;

export const EmptyState = styled.div`
  padding: 40px;
  text-align: center;
  background-color: #f9f9f9;
  border-radius: 8px;
  color: #666;
`;

export const Requests = styled.div`
  display: flex;
  flex-direction: column;
  gap: 15px;
`;

export const RequestHeader = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px;
  background-color: #f9f9f9;
  border-bottom: 1px solid #e0e0e0;
`;

export const RequestTitle = styled.div`
  font-weight: bold;
  font-size: 16px;
`;

export const RequestDetails = styled.div`
  padding: 15px;
`;

export const RequestFrom = styled.div`
  margin-bottom: 5px;
  font-size: 14px;
`;
export const RequestTo = styled.div`
  margin-bottom: 5px;
  font-size: 14px;
`;
export const RequestTime = styled.div`
  margin-bottom: 5px;
  font-size: 14px;
`;
export const RequestDuration = styled.div`
  margin-bottom: 5px;
  font-size: 14px;
`;
export const RequestDescription = styled.div`
  margin-top: 10px;
  padding-top: 10px;
  border-top: 1px solid #f0f0f0;
  font-size: 14px;
  color: #555;
`;

export const RequestActions = styled.div`
  display: flex;
  gap: 10px;
  padding: 15px;
  border-top: 1px solid #f0f0f0;
  background-color: #fafafa;
`;

export const AcceptBtn = styled.button`
  padding: 8px 15px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-weight: 500;
  background-color: #4CAF50;
  color: white;
`;
export const RejectBtn = styled.button`
  padding: 8px 15px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-weight: 500;
  background-color: #F44336;
  color: white;
`;
export const SuggestBtn = styled.button`
  padding: 8px 15px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-weight: 500;
  background-color: #2196F3;
  color: white;
`;
export const CancelBtn = styled.button`
  padding: 8px 15px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-weight: 500;
  background-color: #F44336;
  color: white;
`;
export const EditBtn = styled.button`
  padding: 8px 15px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-weight: 500;
  background-color: #2196F3;
  color: white;
`;

export const RequestAcceptedInfo = styled.div`
  padding: 15px;
  background-color: #fafafa;
  border-top: 1px solid #f0f0f0;
`;
export const RequestRejectedInfo = styled.div`
  padding: 15px;
  background-color: #fafafa;
  border-top: 1px solid #f0f0f0;
`;
export const AcceptedMessage = styled.div`
  color: #388E3C;
`;
export const RejectedMessage = styled.div`
  color: #D32F2F;
`;
export const RejectReason = styled.div`
  margin-top: 5px;
  font-style: italic;
  color: #666;
`;

export const ActionBar = styled.div`
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
`;
export const NewRequestBtn = styled.button`
  padding: 10px 20px;
  background-color: #4CAF50;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-weight: 500;
`;

// 모달 스타일
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
  width: 500px;
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
export const ModalHeaderTitle = styled.h2`
  margin: 0;
  font-size: 18px;
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
export const FormGroupLabel = styled.label`
  display: block;
  margin-bottom: 5px;
  font-weight: 500;
`;
export const FormGroupInput = styled.input`
  width: 100%;
  padding: 8px 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
`;
export const FormGroupSelect = styled.select`
  width: 100%;
  padding: 8px 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
`;
export const FormGroupTextarea = styled.textarea`
  width: 100%;
  padding: 8px 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
`;
export const FormActions = styled.div`
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 20px;
`;
export const SubmitBtn = styled.button`
  padding: 8px 15px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-weight: 500;
  background-color: #4CAF50;
  color: white;
`; 
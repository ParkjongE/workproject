import styled from '@vue-styled-components/core';

export const TeamScheduleContainer = styled.div`
  display: flex;
  height: 100vh;
  width: 100%;
`;

export const ContentArea = styled.div`
  flex: 1;
  padding: 20px;
  overflow-y: auto;
`;

export const DashboardContainer = styled.div`
  width: 100%;
  height: calc(100vh - 100px);
  position: relative;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  background-color: #f9f9f9;
  padding: 20px;
`;

export const TeamGrid = styled.div`
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;
  justify-items: center;
  margin-top: 30px;
`;

export const TeamCard = styled.div`
  background: #fff;
  border-radius: 8px;
  width: 320px;
  height: 180px;
  padding: 18px 22px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.08);
  transition: box-shadow 0.2s, transform 0.2s;
  position: relative;
  display: flex;
  flex-direction: column;
  justify-content: flex-start;
  &:hover {
    box-shadow: 0 6px 18px rgba(33,136,56,0.18);
    transform: translateY(-4px) scale(1.03);
  }
`;

export const TeamTitle = styled.div`
  font-size: 18px;
  font-weight: bold;
  margin-bottom: 10px;
  color: #222;
`;

export const TaskList = styled.ul`
  margin: 0;
  padding: 0;
  list-style: none;
  font-size: 15px;
  color: #234fa0;
`;

export const Deadline = styled.span`
  color: #e74c3c;
  font-size: 13px;
  margin-left: 8px;
`;

export const More = styled.div`
  position: absolute;
  right: 18px;
  bottom: 14px;
  font-size: 13px;
  color: #888;
`;

export const TeamMembers = styled.div`
  margin-top: auto;
  font-size: 14px;
  color: #333;
  padding-top: 10px;
`;

export const TeamLeader = styled.span`
  color: #218838;
  font-weight: bold;
`; 
import styled from '@vue-styled-components/core';

const statusProps = { status: String };
const priorityProps = { priority: String };

export const Container = styled.div`
  max-width: 1500px;
  margin: 20px auto;
  background: #fff;
  border: 1px solid #ddd;
  border-radius: 8px;
  padding: 0;
  box-shadow: 0 2px 6px rgba(0,0,0,0.1);
  position: relative;
  display: flex;
  min-height: 800px;
`;

export const MainContent = styled.div`
  flex: 1;
  padding: 40px 30px 50px 50px;
  overflow-x: auto;
`;

export const HeaderContainer = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
  position: relative;
`;

export const DepartmentSelect = styled.select`
  background-color: #218838;
  border: none;
  border-radius: 12px;
  padding: 8px 20px;
  font-size: 15px;
  color: #fff;
  cursor: pointer;
  outline: none;
  font-weight: bold;
  transition: background 0.2s;
  &:hover {
    background-color: #1e7e34;
  }
`;

export const Table = styled.table`
  width: 100%;
  border-collapse: collapse;
  margin-top: 20px;
  font-size: 15px;
  table-layout: fixed;
`;

export const Th = styled.th`
  border-top: 1px solid #ccc;
  border-right: 1px solid #ccc;
  border-bottom: 1px solid #ccc;
  background-color: #f1f1f1;
  position: sticky;
  top: 0;
  z-index: 10;
  padding: 16px 8px;
  text-align: center;
`;

export const Td = styled.td`
  border-right: 1px solid #ccc;
  border-bottom: 1px solid #ccc;
  padding: 14px 6px;
  text-align: center;
  vertical-align: middle;
  word-wrap: break-word;
  overflow-wrap: break-word;
  line-height: 1.4;
`;

export const BadgeContainer = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100%;
  text-align: center;
`;

export const StatusBadge = styled('span', statusProps)`
  display: inline-block;
  padding: 5px 12px;
  border-radius: 15px;
  font-weight: normal;
  min-width: 90px;
  white-space: nowrap;
  margin: 0 auto;
  box-sizing: border-box;
  background-color: ${props =>
    props.status === '진행중' ? '#e3f2fd' :
    props.status === '완료' ? '#e8f5e9' :
    props.status === '지연' ? '#ffebee' : '#fff'};
  color: ${props =>
    props.status === '진행중' ? '#1976d2' :
    props.status === '완료' ? '#2e7d32' :
    props.status === '지연' ? '#d32f2f' : '#333'};
  border: ${props =>
    props.status === '진행중' ? '1px solid #bbdefb' :
    props.status === '완료' ? '1px solid #c8e6c9' :
    props.status === '지연' ? '1px solid #ffcdd2' : '1px solid #ccc'};
`;

export const PriorityBadge = styled('span', priorityProps)`
  display: inline-block;
  padding: 5px 12px;
  border-radius: 15px;
  font-weight: bold;
  min-width: 100px;
  white-space: nowrap;
  margin: 0 auto;
  box-sizing: border-box;
  background-color: ${props =>
    props.priority === '긴급' || props.priority === '높음' ? '#ffebee' :
    props.priority === '중간' ? '#fff3e0' :
    props.priority === '낮음' ? '#e8f5e9' : '#fff'};
  color: ${props =>
    props.priority === '긴급' || props.priority === '높음' ? '#d32f2f' :
    props.priority === '중간' ? '#ef6c00' :
    props.priority === '낮음' ? '#2e7d32' : '#333'};
  border: ${props =>
    props.priority === '긴급' || props.priority === '높음' ? '1px solid #ffcdd2' :
    props.priority === '중간' ? '1px solid #ffe0b2' :
    props.priority === '낮음' ? '1px solid #c8e6c9' : '1px solid #ccc'};
`; 
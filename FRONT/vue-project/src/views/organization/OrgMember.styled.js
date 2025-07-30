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
`;

export const MainContent = styled.div`
  flex: 1;
  padding: 40px 40px 40px 40px;
`;

export const TopActions = styled.div`
  position: absolute;
  top: 18px;
  right: 12px;
  display: flex;
  gap: 12px;
`;

export const ActionBtn = styled.button`
  background: #218838;
  border: none;
  border-radius: 16px;
  font-size: 20px;
  font-weight: 500;
  padding: 8px 24px;
  cursor: pointer;
  color: #fff;
  font-family: 'Inter', Arial, sans-serif;
  transition: background 0.2s;
`;

export const Title = styled.div`
  font-size: 32px;
  font-weight: 800;
  margin-top: 30px;
  margin-bottom: 24px;
  color: #000;
`;

export const SearchSortRow = styled.div`
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 24px;
`;

export const SearchBox = styled.div`
  display: flex;
  align-items: center;
  background: #f1f1f1;
  border: none;
  border-radius: 8px;
  padding: 0 20px;
  height: 48px;
  min-width: 320px;
  flex: 1;
  max-width: 400px;
`;

export const SearchInput = styled.input`
  border: none;
  background: transparent;
  color: #222;
  font-size: 18px;
  outline: none;
  width: 100%;
  margin-right: 10px;
  padding: 0 4px;
  border-radius: 6px;
  &::placeholder {
    color: #888;
    opacity: 1;
  }
`;

export const SearchBtn = styled.button`
  background: #218838;
  color: #fff;
  border: none;
  border-radius: 6px;
  font-size: 20px;
  cursor: pointer;
  padding: 4px 12px;
  transition: background 0.2s;
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  &:hover {
    background: #17692c;
  }
`;

export const SortGroupBox = styled.div`
  display: flex;
  align-items: center;
  background: #fff;
  border: 1px solid #000;
  border-radius: 16px;
  height: 44px;
  padding: 0 20px;
  gap: 20px;
  span {
    font-family: 'Comfortaa', Arial, sans-serif;
    font-size: 18px;
    font-weight: 400;
    letter-spacing: 0.1em;
    color: #000;
    cursor: pointer;
  }
`;

export const MemberList = styled.div`
  display: flex;
  flex-wrap: wrap;
  gap: 24px 16px;
  justify-content: flex-start;
  margin-top: 24px;
`;

export const MemberCard = styled.div`
  width: 240px;
  min-width: 200px;
  height: 300px;
  background: #fff;
  border: 1px solid #000;
  border-radius: 8px;
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  box-shadow: 0 2px 6px rgba(0,0,0,0.04);
  font-size: 14px;
  padding: 12px 8px 16px 8px;
`;

export const MemberPhoto = styled.div`
  width: 70px;
  height: 70px;
  background: #D9D9D9;
  border-radius: 50%;
  margin-top: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  color: #888;
`;

export const MemberInfo = styled.div`
  width: 100%;
  margin-top: 8px;
  display: flex;
  flex-direction: column;
  gap: 2px;
  font-size: 13px;
`;

export const MemberInfoRow = styled.div`
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
  font-weight: 500;
  color: #000;
  flex-wrap: wrap;
`;

export const Label = styled.span`
  width: 38px;
  color: #888;
  font-size: 12px;
  font-weight: 400;
  flex-shrink: 0;
`;

export const MemberStatus = styled.span`
  background: #218838;
  color: #fff;
  margin-top: 8px;
  position: static;
  display: inline-block;
  margin: 0 4px;
  font-size: 13px;
  font-weight: 500;
  border-radius: 12px;
  padding: 2px 12px;
  border: none;
`;

export const EditBtn = styled.button`
  background: #218838;
  color: #fff;
  margin-top: 8px;
  cursor: pointer;
  border: none;
  position: static;
  display: inline-block;
  margin: 0 4px;
  font-size: 13px;
  font-weight: 500;
  border-radius: 12px;
  padding: 2px 12px;
`; 
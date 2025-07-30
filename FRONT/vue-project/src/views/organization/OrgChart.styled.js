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
  @media (max-width: 1200px) {
    margin-left: 0;
    padding: 24px 8px;
  }
`;

export const MainTitle = styled.div`
  font-size: 36px;
  font-weight: bold;
  margin-bottom: 40px;
  margin-top: 30px;
`;

export const OrgChartArea = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-top: 40px;
  position: relative;
  width: 900px;
  min-height: 500px;
  @media (max-width: 1200px) {
    margin-top: 20px;
    width: 100vw;
    min-width: 0;
  }
`;

export const OrgLevel = styled.div`
  display: flex;
  justify-content: center;
  gap: 40px;
  margin-bottom: 40px;
  position: relative;
  flex-wrap: wrap;
  width: 100%;
  z-index: 1;
  @media (max-width: 1200px) {
    max-width: 100vw;
  }
`;

export const OrgNode = styled.div`
  background: #f1f1f1;
  border: 1px solid #ccc;
  border-radius: 12px;
  min-width: 120px;
  min-height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: 400;
  color: #222;
  font-family: Arial, sans-serif;
  position: relative;
  z-index: 2;
  box-shadow: 0 2px 6px rgba(0,0,0,0.04);
  margin-bottom: 10px;
  margin-top: 10px;
  word-break: keep-all;
`;

export const OrgSvg = styled.svg`
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
  z-index: 0;
`; 
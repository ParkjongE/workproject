import styled from '@vue-styled-components/core';

export const RegisterContainer = styled.div`
  width: 869px;
  margin: 120px auto 0 auto;
  background: #fff;
  box-shadow: 0px 4px 4px 0px rgba(0,0,0,0.25);
  border-radius: 5px;
  padding: 60px 80px 40px 80px;
  box-sizing: border-box;
  position: relative;
  @media (max-width: 1000px) {
    width: 95vw;
    padding: 30px 10px;
  }
`;

export const Logo = styled.div`
  text-align: center;
  font-family: 'Comfortaa', sans-serif;
  font-size: 28px;
  color: #000;
  margin-bottom: 60px;
  letter-spacing: -1.5%;
`;

export const RegisterTitle = styled.div`
  text-align: center;
  font-family: 'Corben', sans-serif;
  font-size: 24px;
  color: #757678;
  margin-bottom: 40px;
`;

export const InputLabel = styled.label`
  font-family: 'Comfortaa', sans-serif;
  font-size: 18px;
  color: #757678;
  margin-bottom: 8px;
  display: block;
  letter-spacing: 0.1em;
`;

export const RegisterInput = styled.input`
  width: 100%;
  height: 48px;
  border: 1px solid #757678;
  border-radius: 5px;
  font-size: 18px;
  padding: 0 16px;
  margin-bottom: 24px;
  box-sizing: border-box;
  font-family: 'Comfortaa', sans-serif;
`;

export const RegisterSelect = styled.select`
  width: 100%;
  height: 48px;
  border: 1px solid #757678;
  border-radius: 5px;
  font-size: 18px;
  padding: 0 16px;
  margin-bottom: 24px;
  box-sizing: border-box;
  font-family: 'Comfortaa', sans-serif;
`;

export const RegisterBtn = styled.button`
  width: 100%;
  height: 56px;
  background: #218838;
  color: #fff;
  border: none;
  border-radius: 5px;
  font-size: 24px;
  font-family: 'Comfortaa', sans-serif;
  font-weight: 400;
  cursor: pointer;
  margin-top: 16px;
  margin-bottom: 24px;
  transition: background 0.2s;
  &:hover {
    background: #17692a;
  }
  &:disabled {
    background: #cccccc;
    cursor: not-allowed;
  }
`;

export const ErrorMessage = styled.div`
  color: #dc3545;
  margin-bottom: 16px;
  font-size: 14px;
  text-align: center;
  font-family: 'Comfortaa', sans-serif;
`;

export const RegisterLinks = styled.div`
  text-align: center;
  font-size: 14px;
  color: #757678;
  font-family: 'Comfortaa', sans-serif;
  margin-top: 20px;
  & a {
    color: #218838;
    text-decoration: none;
    margin-left: 8px;
  }
`; 
import Styled from "styled-components";

const BlankContainer = Styled.div`
  height: 100%;
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  ${({ bg }) => `background: url(${bg}) no-repeat center center;`};
`;

export default BlankContainer;

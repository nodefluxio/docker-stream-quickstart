import Styled from "styled-components";

const Text = Styled.div`
  font-style: normal;
  line-height: ${props => props.lineHeight || `14px`};
  ${({ size }) => size && `font-size: ${size}px;`}
  ${({ color }) => color && `color: ${color};`}
  ${({ weight }) => weight && `font-weight: ${weight};`}
  ${({ lineHeight }) => lineHeight && `line-height: ${lineHeight}px;`}
  ${({ align }) => align && `text-align: ${align};`}
  ${({ textTransform }) => textTransform && `text-transform: ${textTransform};`}
  ${({ marginBottom }) => marginBottom && `margin-bottom: ${marginBottom};`}
`;

export default Text;

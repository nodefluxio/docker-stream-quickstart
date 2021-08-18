import React from "react";
import Styled, { css } from "styled-components";
import PropTypes from "prop-types";

const Tooltip = ({ text, position = "cb", children }) => (
  <Wrapper className="Tooltip" position={position}>
    {children}
    <Label>{text}</Label>
  </Wrapper>
);
export default Tooltip;

Tooltip.propTypes = {
  text: PropTypes.oneOfType([
    PropTypes.object.isRequired,
    PropTypes.string.isRequired
  ]).isRequired,
  position: PropTypes.string,
  children: PropTypes.object.isRequired
};

const Label = Styled.span`
  display: inline-block;
  position: absolute;
  padding: 9px 15px;
  background-color: #4995E9;
  border-radius: 2px;
  color: white;
  width: max-content;
  opacity: 0;
  pointer-events: none;
  font-weight: 500;
  font-size: 12px;
  line-height: 14px;
  z-index: 9;
`;

const Wrapper = Styled.div`
  position: relative;

  ${Label}{
    ${props => {
      switch (props.position) {
        case "rc":
          return css`
            left: calc(100% + 3px);
            top: 50%;
            transform: translateY(-50%);
          `;
        case "lc":
          return css`
            right: calc(100% + 3px);
            top: 50%;
            transform: translateY(-50%);
          `;
        case "ct":
          return css`
            left: 50%;
            bottom: calc(100% + 3px);
            transform: translateX(-50%);
          `;
        default:
          return css`
            left: 50%;
            top: calc(100% + 3px);
            transform: translateX(-50%);
          `;
      }
    }}
  }

  :hover{
    ${Label}{
      opacity: 1;
    }
  }

`;

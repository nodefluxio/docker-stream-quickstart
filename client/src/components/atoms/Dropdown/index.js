import React, { useState } from "react";
import PropTypes from "prop-types";
import styled from "styled-components";
import hexToRgbA from "helpers/hexToRgbA";

export default function Dropdown(props) {
  const [show, setShow] = useState(false);
  const { id, children, overlay, placement, width, className } = props;

  function toggleDropdown() {
    setShow(!show);
  }

  return (
    <BoxDropdown
      id={id}
      onClick={toggleDropdown}
      show={show}
      className={className}
    >
      {children}
      <WrapDropdown show={show} placement={placement} width={width}>
        {overlay}
      </WrapDropdown>
      {show && <Overlay onClick={toggleDropdown} />}
    </BoxDropdown>
  );
}

Dropdown.defaultProps = {
  placement: "",
  className: "",
  id: ""
};

Dropdown.propTypes = {
  overlay: PropTypes.oneOfType([PropTypes.element, PropTypes.array]).isRequired,
  children: PropTypes.element.isRequired,
  placement: PropTypes.string,
  width: PropTypes.string,
  className: PropTypes.string,
  id: PropTypes.string
};

export const Menu = ({ id, children, onClick, className }) => (
  <MenuItemList id={id} onClick={onClick} className={className}>
    {children}
  </MenuItemList>
);

Menu.propTypes = {
  children: PropTypes.any.isRequired,
  onClick: PropTypes.func.isRequired,
  className: PropTypes.string,
  id: PropTypes.string
};

Menu.defaultProps = {
  className: "",
  id: ""
};

const BoxDropdown = styled.div`
  position: relative;
  ${props => props.show && `background-color: #4c12a1`}
`;

const positionPlacement = placement => {
  switch (placement) {
    case "center":
      return `
        margin-left: auto;
        margin-right: auto;
        left: 0;
        right: 0;       
      `;
    case "left":
      return `
        right: unset;
        left: 0; 
      `;
    default:
      return `
        right: 0;
      `;
  }
};

const WrapDropdown = styled.div`
  position: absolute;
  z-index: -1;
  box-sizing: border-box;
  width: ${props => props.width || `max-content`};
  height: auto;
  background-color: ${props => props.theme.blueGem};
  opacity: 0;
  transition: 0.1s ease-in;
  visibility: hidden;
  ${({ show }) =>
    show &&
    `
    z-index: 9;
    opacity: 1;
    visibility: visible;
    transition: 0.2s ease-out;
  `};
  ${({ placement }) => positionPlacement(placement)};
`;

const MenuItemList = styled.div`
  min-width: 100px;
  width: 100%;
  display: flex;
  align-items: center;
  text-align: left;
  color: ${props => props.theme.white};
  font-size: 12px;
  font-weight: 700px;
  box-sizing: border-box;
  font-style: normal;
  font-stretch: normal;
  line-height: 14px;
  padding: 13px;
  outline: none;
  border-bottom: 1px solid;
  border-color: ${props => hexToRgbA(props.theme.white, 0.1)};
  ${({ onClick }) => onClick && `cursor: pointer;`};
  a {
    text-decoration: none;
    color: ${props => props.theme.mulberry};
  }
  :hover {
    background-color: ${props => props.theme.secondary2};
    transition: 0.2s;
  }
`;

const Overlay = styled.div`
  position: fixed;
  content: "";
  height: 100%;
  width: 100%;
  left: 0;
  top: 0;
  z-index: 5;
`;

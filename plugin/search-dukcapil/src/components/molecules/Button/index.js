import React from "react";
import PropTypes from "prop-types";
import styled, { css } from "styled-components";
import LoadingSpinner from "../../atoms/LoadingSpinner";

export default function Button(props) {
  const {
    isLoading,
    onClick,
    className,
    style,
    type,
    children,
    disabled,
    name,
    width,
    id
  } = props;

  return (
    <ButtonStyle
      isLoading={!isLoading}
      disabled={disabled}
      onClick={onClick}
      className={className}
      name={name}
      style={style}
      type={type}
      width={width}
      id={id}
    >
      {isLoading === true ? <LoadingSpinner show={true} /> : children}
    </ButtonStyle>
  );
}

Button.propTypes = {
  onClick: PropTypes.func,
  isLoading: PropTypes.bool,
  name: PropTypes.string,
  className: PropTypes.string,
  style: PropTypes.object,
  type: PropTypes.oneOf(["primary", "icon", "secondary", "blue"]),
  children: PropTypes.oneOfType([
    PropTypes.object.isRequired,
    PropTypes.string.isRequired,
    PropTypes.array.isRequired
  ]),
  disabled: PropTypes.bool,
  width: PropTypes.string,
  id: PropTypes.string
};

Button.defaultProps = {
  onClick: () => {},
  isLoading: false,
  className: "",
  name: "",
  type: "primary",
  disabled: false,
  width: "auto",
  style: {},
  id: ""
};

const ButtonStyle = styled.button`
  font-size: 12px;
  letter-spacing: 1px;
  line-height: 14px;
  text-decoration: none;
  border-radius: 8px;
  outline: none;
  display: flex;
  font-weight: bold;
  justify-content: center;
  align-items: center;
  margin-bottom: 0;
  touch-action: manipulation;
  cursor: pointer;
  background-image: none;
  border: 1px solid transparent;
  white-space: nowrap;
  user-select: none;
  text-transform: uppercase;
  background: ${props => props.theme.blue1};
  min-height: 40px;
  min-width: 147px;
  padding: 0px 16px;
  width: ${props => props.width};

  &.icon-text img {
    margin-right: 10px;
  }
  &:active{
    background: ${props => props.theme.mint1};
    transform: translateY(4px);
  }
  &:hover:enabled {
    background: ${props => props.theme.blue1};
  }
  & > :not(:first-child) {
    margin-left: 8px;
  }
  & > * {
    display: flex !important;
  }
  ${props =>
    props.disabled &&
    `cursor: default;
  opacity: 0.7;`}


  ${props =>
    props.isLoading &&
    css`
      border: none;
      color: ${props.theme.white};
      font-weight: bold;
      color: white;
    `}

  ${props =>
    props.type === "secondary" &&
    css`
      background: ${props.theme.bg};
      border: 2px solid ${props.theme.secondary2};
      &:hover:enabled {
        background: ${props.theme.blueGem};
        border: 2px solid ${props.theme.blueGem};
      }
    `}
  ${props =>
    props.type === "icon" &&
    css`
      border: 1px solid ${props.theme.secondary2};
      border-radius: 8px;
      background-color: transparent;
      height: 40px;
      padding: 0px;
      min-width: 40px;
    `}
`;

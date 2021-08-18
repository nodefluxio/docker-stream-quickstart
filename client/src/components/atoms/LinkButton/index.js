import React from "react";
import Styled from "styled-components";
import PropTypes from "prop-types";

export default function LinkButton({
  children,
  id,
  onClick,
  disabled,
  bordered,
  style,
  className
}) {
  return (
    <LinkButtonStyled
      bordered={bordered}
      id={id}
      onClick={onClick}
      disabled={disabled}
      style={style}
      className={className}
    >
      {children}
    </LinkButtonStyled>
  );
}

LinkButton.propTypes = {
  children: PropTypes.any.isRequired,
  id: PropTypes.string.isRequired,
  onClick: PropTypes.func.isRequired,
  disabled: PropTypes.bool,
  bordered: PropTypes.bool,
  style: PropTypes.object,
  className: PropTypes.string
};

LinkButton.defaultProps = {
  disabled: false,
  bordered: false,
  style: {}
};

const LinkButtonStyled = Styled.button`
  border-top: none;
  border-bottom: none;
  border-right: ${props => (props.bordered ? `1px solid #372463` : `none`)};
  border-left: ${props => (props.bordered ? `1px solid #372463` : `none`)};
  background: none;
  outline: none;
  cursor: pointer;
  display: flex;
  align-items: center;
  margin-right: ${props => (props.bordered ? `0px` : `36px`)};
  height: 100%;

  :last-child{
    margin-right: ${props => (props.bordered ? `0px` : `18px`)};
  }
  
  img,svg{
    margin-right: ${props => (props.bordered ? `0px` : `16px`)};
  }
  span{
    font-family: Barlow;
    font-style: normal;
    font-weight: 500;
    font-size: 12px;
    line-height: 14px;
    text-transform: uppercase;
    color: ${props => props.theme.mercury};
  }

  :hover, &.active {
    background: ${props => props.theme.blueGem};
    border-bottom: 3px solid ${props => props.theme.mint};
  }
`;

import React, { Fragment, useState } from "react";
import Styled, { css } from "styled-components";
import PropTypes from "prop-types";

import Text from "components/atoms/Text";
import Row from "components/atoms/Row";
import hexToRgbA from "helpers/hexToRgbA";

export default function AreaSection(props) {
  const {
    rightIcon,
    leftIcon,
    title,
    titleColor,
    children,
    border,
    accordion,
    secondary
  } = props;
  const [toggle, setToggle] = useState(true);

  function toggleAccordion() {
    setToggle(!toggle);
  }

  return (
    <Fragment>
      <AreaSectionWrap
        onClick={accordion ? toggleAccordion : null}
        border={border}
        accordion={accordion}
        secondary={secondary}
      >
        <Row horizontalGap={10}>
          {leftIcon && <Row>{leftIcon}</Row>}
          <Text size="12px" color={titleColor} weight="bold">
            {title}
          </Text>
        </Row>
        {rightIcon && <RotateIcon toggle={toggle}>{rightIcon}</RotateIcon>}
      </AreaSectionWrap>
      {toggle && children}
    </Fragment>
  );
}

AreaSection.defaultProps = {
  rightIcon: null,
  leftIcon: null,
  children: null,
  accordion: false,
  border: false,
  secondary: false
};

AreaSection.propTypes = {
  rightIcon: PropTypes.element,
  leftIcon: PropTypes.element,
  children: PropTypes.any,
  titleColor: PropTypes.string.isRequired,
  title: PropTypes.string.isRequired,
  accordion: PropTypes.bool,
  border: PropTypes.bool,
  secondary: PropTypes.bool
};

const RotateIcon = Styled.div`
    .rotate{
        transform: rotate(0);
        transition: 0.5s;
        ${({ toggle }) =>
          toggle &&
          `
        transform: rotate(180deg);
        `}
    }
`;

const AreaSectionWrap = Styled.div`
  position: relative;
  width: 100%;
  height: 40px;
  background-color: ${props => props.theme.secondary2} ;
  display: flex;
  justify-content: space-between;
  align-items: center;
  color: ${props => props.theme.mercury};
  ${({ border }) =>
    border &&
    `
  border-bottom: 1px solid ${props => hexToRgbA(props.theme.white, 0.1)};
  `}
  ${({ secondary }) =>
    secondary &&
    css`
      background-color: transparent;
    `}
`;

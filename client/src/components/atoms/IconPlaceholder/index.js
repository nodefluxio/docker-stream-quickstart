import React from "react";
import PropType from "prop-types";
import Styled from "styled-components";

export default function IconPlaceholder(props) {
  const {
    id,
    borderColor,
    children,
    width,
    height,
    onClick,
    className,
    disable
  } = props;
  return (
    <IconWrapper
      id={id}
      onClick={onClick}
      borderColor={borderColor}
      width={width}
      height={height}
      className={className}
      disable={disable}
    >
      {children}
    </IconWrapper>
  );
}

IconPlaceholder.propTypes = {
  borderColor: PropType.string,
  id: PropType.string,
  children: PropType.object.isRequired,
  height: PropType.string,
  onClick: PropType.func,
  width: PropType.string,
  className: PropType.string,
  disable: PropType.bool
};

IconPlaceholder.defaultProps = {
  borderColor: "transparent",
  id: "",
  onClick: null,
  height: "40px",
  width: "40px",
  className: "",
  disable: false
};

const IconWrapper = Styled.div`
  position: relative;
  width: ${props => props.width};
  height: ${props => props.height};
  border-radius: 4px;
  border: 2px solid ${props => props.borderColor};
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  cursor: ${props =>
    typeof props.onClick === "function" && !props.disable
      ? `pointer`
      : `default`}
  ${props => props.disable && `opacity: 0.5;`}
  img { 
    @media (orientation: landscape) {
      max-height: 100%
    }
    @media (orientation: potrait) {
      max-width: 100%;
    }
  }
`;

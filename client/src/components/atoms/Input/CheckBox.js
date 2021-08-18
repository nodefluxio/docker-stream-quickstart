import React from "react";
import PropTypes from "prop-types";
import styled from "styled-components";

export default function Checkbox(props) {
  const { value, text, checked, disable, circle } = props;

  const getTwoFunc = e => {
    const val = e.target.value;
    const { checked: check } = e.target;
    const { onChange } = props;
    onChange(val, check);
  };

  return (
    <Container checked={checked} circle={circle}>
      {text}
      <input
        value={value}
        className="analyticCheckbox"
        onChange={getTwoFunc}
        type="checkbox"
        checked={checked}
        disabled={disable}
      />
      <CustomCheckbox checked={checked} circle={circle}>
        {checked && (
          <div className="wrap-check">
            <div className="custom-check" />
          </div>
        )}
      </CustomCheckbox>
    </Container>
  );
}

Checkbox.propTypes = {
  onChange: PropTypes.func.isRequired,
  value: PropTypes.any.isRequired,
  disable: PropTypes.object,
  circle: PropTypes.bool,
  checked: PropTypes.bool,
  text: PropTypes.string
};

const Container = styled.label`
  display: block;
  position: relative;
  padding-left: 25px;
  margin-bottom: 12px;
  cursor: pointer;
  font-size: 14px;
  user-select: none;
  color: #e5e5e5;
  input {
    position: absolute;
    opacity: 0;
    cursor: pointer;
    height: 0;
    width: 0;
  }
  span {
    .custom-check {
      left: 4px;
      top: 1px;
      width: 3px;
      height: 6px;
      border: solid white;
      border-width: 0 1px 1px 0;
      transform: rotate(45deg);
    }
  }
  :hover input ~ span {
    background: #4995e9;
    border: 2px solid #4995e9;
  }
  ${props =>
    props.circle &&
    `
    input:checked ~ span {
      .wrap-check{
        height: 12)};
        border-radius: 50%;
        width: 12)};
        background-color: #7ED321;
        border: unset;
      }
    }
  `}
`;

const CustomCheckbox = styled.span`
  position: absolute;
  top: 2px;
  z-index: 2px;
  left: 0;
  height: 14px;
  width: 14px;
  border: 2px solid #f7f7f7;
  border-radius: 2px;
  .wrap-check {
    height: 14px;
    width: 14px;
  }
  .custom-check {
    content: "";
    position: absolute;
  }
  ${props =>
    props.circle &&
    `
      border-radius: 50%;
      background-color: unset;
      border:1px solid white;
  `}
  ${props =>
    props.checked &&
    `
      border: 2px solid #4995e9;
      background-color: #4995e9;
    `}
`;

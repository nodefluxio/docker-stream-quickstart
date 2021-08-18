import React, { Component } from "react";
import styled from "styled-components";
import PropTypes from "prop-types";

const RelativeDivCol = styled.div`
  position: relative;
  display: flex;
  flex-direction: column;
  width: 100%;
`;

const RelativeDivRow = styled.div`
  position: relative;
  display: flex;
  flex-direction: row;
  align-items: center;
`;

const InputLabel = styled.label`
  width: 100%;
  font-size: 14px;
  font-weight: 500;
  font-style: normal;
  font-stretch: normal;
  line-height: normal;
  letter-spacing: normal;
  display: block;
  position: relative;
  color: ${props => props.theme.mint};
  margin-top: 10px;
  margin-bottom: 10px;
`;

const InputText = styled.input`
  background-color: transparent;
  width: 100%;
  display: block;
  border: none;
  height: ${props => props.height || `40px`};
  padding-left: ${props => (props.addonBefore ? `30px` : `10px`)};
  padding-right: ${props => (props.addonAfter ? `30px` : `10px`)};
  box-sizing: border-box;
  border-radius: 8px;
  color: ${props => props.theme.mercury};
  border: ${props => {
    if (props.error) {
      return `2px solid ${props.theme.inlineError}`;
    }
    return `2px solid ${props.theme.secondary2}`;
  }};
  font-size: ${props => props.fontSize || `14px`};
  position: relative;
  ${props => props.readOnly && `background-color: ${props.theme.secondary2}`};
  ::placeholder {
    font-size: ${props => props.fontSize || `14px`};
    font-weight: 500;
    font-style: normal;
    font-stretch: normal;
    line-height: revert;
    letter-spacing: normal;
  }
  :focus {
    outline: none;
    border: 2px solid ${props => props.theme.mint};
  }
`;

const InputAfter = styled.div`
  position: absolute;
  right: 0;
`;

const InputBefore = styled.div`
  position: absolute;
  left: 10px;
`;

const InputError = styled(InputLabel)`
  font-size: 12px;
  font-weight: 500;
  color: ${props => props.theme.color8};
  margin-bottom: 5px;
`;
export default class Text extends Component {
  render() {
    const {
      label,
      tooltip,
      children,
      error,
      fontSize,
      addonAfter,
      addonBefore,
      ...defaultProps
    } = this.props;

    if (label || tooltip || children || error || addonAfter || addonBefore) {
      return (
        <RelativeDivCol>
          {label && (
            <InputLabel>
              {label}
              {tooltip && tooltip}
            </InputLabel>
          )}
          <RelativeDivRow>
            {addonBefore && <InputBefore>{addonBefore}</InputBefore>}
            <InputText
              {...defaultProps}
              error={error}
              fontSize={fontSize}
              addonAfter={addonAfter}
              addonBefore={addonBefore}
            />
            {children}
            {addonAfter && <InputAfter>{addonAfter}</InputAfter>}
          </RelativeDivRow>
          {error && <InputError>{error}</InputError>}
        </RelativeDivCol>
      );
    }
    return <InputText {...defaultProps} />;
  }
}

Text.propTypes = {
  label: PropTypes.string,
  tooltip: PropTypes.element,
  error: PropTypes.oneOfType([PropTypes.string, PropTypes.bool]),
  children: PropTypes.element,
  addonAfter: PropTypes.element,
  addonBefore: PropTypes.element,
  fontSize: PropTypes.string
};

import React, { Component } from "react";
import Styled from "styled-components";
import PropTypes from "prop-types";

class RadioButton extends Component {
  render() {
    const RadioGroup = Styled.div`
    position: relative;
    width: 100%;
    display: flex;
    flex: 1 1 auto;
    flex-wrap: wrap;
    justify-content: center;
    margin-top: 10px;
    `;

    const { children } = this.props;
    return <RadioGroup>{children}</RadioGroup>;
  }
}

RadioButton.propTypes = {
  children: PropTypes.element.isRequired
};

export default RadioButton;

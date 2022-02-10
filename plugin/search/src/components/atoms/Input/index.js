import React from "react";
import PropTypes from "prop-types";
import Text from "./Text";
import SliderInput from "./Slider";

export default function Input(props) {
  switch (props.type) {
    case "slider":
      return <SliderInput {...props} />;
    default:
      return <Text {...props} />;
  }
}

Input.propTypes = {
  type: PropTypes.string.isRequired
};

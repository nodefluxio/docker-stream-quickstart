import React from "react";
import PropTypes from "prop-types";
import RadioButton from "./RadioButton";
import TextArea from "./TextArea";
import Text from "./Text";
import Select from "./Selector";
import Checkbox from "./CheckBox";

export default function Input(props) {
  switch (props.type) {
    case "radio":
      return <RadioButton {...props} />;
    case "textarea":
      return <TextArea {...props} />;
    case "select":
      return <Select {...props} />;
    case "checkbox":
      return <Checkbox {...props} />;
    default:
      return <Text {...props} />;
  }
}

Input.propTypes = {
  type: PropTypes.string
};

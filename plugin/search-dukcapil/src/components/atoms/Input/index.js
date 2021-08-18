import React from "react";
import PropTypes from "prop-types";
import Text from "./Text";

export default function Input(props) {
  switch (props.type) {
    default:
      return <Text {...props} />;
  }
}

Input.propTypes = {
  type: PropTypes.string.isRequired
};

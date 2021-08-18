import React from "react";
import PropTypes from "prop-types";
import { COUNTER_LINE } from "constants/roiType";
import CounterLine from "./CounterLine";

export default function Roi(props) {
  switch (props.roiType) {
    case COUNTER_LINE:
      return <CounterLine {...props} />;
    default:
      return <CounterLine {...props} />;
  }
}

Roi.propTypes = {
  roiType: PropTypes.string.isRequired
};

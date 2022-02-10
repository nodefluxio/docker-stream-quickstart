import React from "react";
import PropTypes from "prop-types";
import { COUNTER_LINE, REGION_INTEREST } from "constants/roiType";
import CounterLine from "./CounterLine";
import RegionInterest from "./RegionInterest";

export default function Roi(props) {
  switch (props.roiType) {
    case COUNTER_LINE:
      return <CounterLine {...props} />;
    case REGION_INTEREST:
      return <RegionInterest {...props} />;
    default:
      return <CounterLine {...props} />;
  }
}

Roi.propTypes = {
  roiType: PropTypes.string.isRequired
};

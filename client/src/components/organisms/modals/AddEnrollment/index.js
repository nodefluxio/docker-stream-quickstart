import React from "react";
import PropTypes from "prop-types";
import AddPerson from "./AddPerson";
import AddVehicle from "./AddVehicle";

export default function AddEnrollment(props) {
  const { type } = props;
  switch (type) {
    case "person":
      return <AddPerson {...props} />;
    case "vehicle":
      return <AddVehicle {...props} />;
    default:
      return null;
  }
}

AddEnrollment.propTypes = {
  type: PropTypes.string.isRequired
};

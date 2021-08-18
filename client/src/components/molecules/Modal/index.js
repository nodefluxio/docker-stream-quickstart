import React from "react";
import PropTypes from "prop-types";
import { setAppElement } from "react-modal";
import DefaultModal from "components/atoms/Modal";
import ModalFeedback from "./feedback";
import ConfrimationModal from "./ConfirmationModal";
import SearchModal from "./SearchModal";

setAppElement("body");

export default function Modal(props) {
  switch (props.type) {
    case "feedback": {
      return <ModalFeedback {...props} />;
    }
    case "confirmation": {
      return <ConfrimationModal {...props} />;
    }
    case "search": {
      return <SearchModal {...props} />;
    }
    default:
      return <DefaultModal {...props} />;
  }
}

Modal.propTypes = {
  type: PropTypes.string
};

Modal.defaultProps = {
  type: "default"
};

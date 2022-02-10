import React from "react";
import PropTypes from "prop-types";
import LoadingSpinner from "components/atoms/LoadingSpinner";
import Modal from "components/atoms/Modal";

function BackupEnrollment(props) {
  const { openModal } = props;

  const closeModal = () => {
    // eslint-disable-next-line no-console
    console.log("Backup is still in progress, please wait...");
  };

  return (
    <Modal
      show={openModal}
      className="modal-add-enrollment"
      title="Backup Enrollment"
      close={() => closeModal()}
      padding="20px"
      width="800px"
    >
      <center>
        Backup is in progress, please don't close the window...
        <LoadingSpinner show={openModal} />
      </center>
    </Modal>
  );
}

BackupEnrollment.propTypes = {
  openModal: PropTypes.bool.isRequired,
  onClose: PropTypes.func,
  data: PropTypes.array
};

export default BackupEnrollment;

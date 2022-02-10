import React, { useState, useCallback, useEffect } from "react";
import Styled from "styled-components";
import PropTypes from "prop-types";
import Modal from "components/atoms/Modal";
import { EditableInput } from "components/molecules/Input";
import Button from "components/molecules/Button";
import { register, update } from "api/vehicle";

function AddVehicle(props) {
  const { openModal, onClose, data } = props;
  const formSchema = {
    plate_number: "",
    unique_id: "",
    name: "",
    type: "",
    brand: "",
    color: "",
    status: ""
  };
  const [title, setTitle] = useState("Add");
  const [formData, setFormData] = useState(formSchema);
  const [errorMsg, setErrorMsg] = useState("");
  const [isLoading, setIsLoading] = useState(false);

  const handleChange = useCallback(event => {
    const { name, value } = event.target;
    setFormData(prevState => ({
      ...prevState,
      [name]: value
    }));
  });

  const formatLabel = string =>
    string.replace(/_/g, " ").replace(/\b(\w)/g, str => str.toUpperCase());

  function save() {
    setIsLoading(true);
    if (title === "Add New") {
      register(formData)
        .then(result => {
          if (result.ok) {
            setErrorMsg("");
            onClose(result.enrollment);
            setIsLoading(false);
          }
        })
        .catch(error => {
          setErrorMsg(error.response.data.message);
          setIsLoading(false);
        });
    } else if (title === "Update") {
      update(data.id, formData)
        .then(result => {
          if (result.ok) {
            setIsLoading(false);
            setErrorMsg("");
            onClose(result.enrollment);
          }
        })
        .catch(error => {
          setIsLoading(false);
          setErrorMsg(error.response.data.message);
        });
    }
  }

  useEffect(() => {
    if (openModal === false) {
      setFormData(formSchema);
      setTitle("Add New");
    }
  }, [openModal]);

  useEffect(() => {
    if (Object.keys(data).length !== 0 && data.constructor === Object) {
      setTitle("Update");

      setFormData({
        ...formData,
        plate_number: data.plate_number,
        name: data.name,
        type: data.type,
        brand: data.brand,
        status: data.status,
        unique_id: data.unique_id,
        color: data.color
      });
    }
  }, [data]);

  return (
    <Modal
      show={openModal}
      className="modal-vehicle"
      title={`${title} Vehicle`}
      close={onClose}
      padding="20px"
      width="400px"
    >
      <FormRow>
        {Object.keys(formData).map(
          key =>
            key !== "images" && (
              <EditableInput
                key={key}
                label={formatLabel(key)}
                id={key}
                onChange={handleChange}
                name={key}
                value={formData[key]}
              ></EditableInput>
            )
        )}
      </FormRow>
      {errorMsg !== "" && (
        <div
          style={{
            color: "red",
            textAlign: "left"
          }}
        >
          {errorMsg}
        </div>
      )}
      <Button
        width="100%"
        style={{ marginTop: "40px" }}
        onClick={() => save()}
        isLoading={isLoading}
        disabled={isLoading}
      >
        Save
      </Button>
    </Modal>
  );
}

const FormRow = Styled.div`
  display: flex;
  flex-direction: column;
  .input-group{
    margin-bottom: 20px;
    input {
      margin-bottom: 0px;
    }
  }
`;

AddVehicle.propTypes = {
  openModal: PropTypes.bool.isRequired,
  onClose: PropTypes.func.isRequired,
  data: PropTypes.object.isRequired
};

export default AddVehicle;

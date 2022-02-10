import React, { useState, useEffect, useContext } from "react";
import Styled, { ThemeContext } from "styled-components";
import PropTypes from "prop-types";
import Modal from "components/atoms/Modal";
import Input from "components/atoms/Input";
import { EditableInput } from "components/molecules/Input";
import Text from "components/atoms/Text";
import Button from "components/molecules/Button";

import { createCamera, createStream } from "api";

function AddCamera(props) {
  const themeContext = useContext(ThemeContext);
  const { openModal, onClose } = props;
  const formSchema = {
    stream_name: "",
    stream_address: "",
    stream_longitude: 0,
    stream_latitude: 0,
    stream_site: "",
    stream_brand: ""
  };
  const [formData, setFormData] = useState(formSchema);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const handleChange = event => {
    const { name, value } = event.target;
    setFormData(prevState => ({
      ...prevState,
      [name]: value
    }));
  };
  const formatLabel = string => {
    const splitString = string.split("_");
    const keys =
      splitString[1].charAt(0).toUpperCase() + splitString[1].slice(1);
    return `Camera ${keys}`;
  };

  const submitForm = () => {
    const form = formData;
    form.stream_latitude = parseFloat(formData.stream_latitude);
    form.stream_longitude = parseFloat(formData.stream_longitude);
    setLoading(true);
    createCamera(form).then(result => {
      if (result.code === 200) {
        createStream(result.stream_id).then(results => {
          if (results.code === 200) {
            setLoading(false);
            onClose(true);
          } else {
            setError(results.message);
            setLoading(false);
          }
        });
      } else {
        setError(result.message);
        setLoading(false);
      }
    });
  };

  useEffect(() => {
    setFormData(formSchema);
  }, [openModal]);

  function renderForm(name) {
    switch (name) {
      case "stream_address":
        return (
          <Input
            key={name}
            label="Source"
            type="textarea"
            id="source_field"
            height="232px"
            width="345px"
            value={formData[name]}
            name="stream_address"
            onChange={handleChange}
          ></Input>
        );
      default:
        return (
          <EditableInput
            key={name}
            label={formatLabel(name)}
            id={name}
            onChange={handleChange}
            name={name}
            value={formData[name]}
          ></EditableInput>
        );
    }
  }

  return (
    <Modal
      show={openModal}
      className="modal-confirmation"
      title="Add Camera"
      close={() => onClose(false)}
      padding="20px"
    >
      {error && <Text color={themeContext.theme.inlineError}>{error}</Text>}
      <FormRow>{Object.keys(formData).map(key => renderForm(key))}</FormRow>
      <Button
        width="100%"
        style={{ marginTop: "20px" }}
        onClick={() => submitForm()}
        isLoading={loading}
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

AddCamera.propTypes = {
  openModal: PropTypes.bool.isRequired,
  onClose: PropTypes.func.isRequired
};

export default AddCamera;

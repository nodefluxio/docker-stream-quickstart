import React, { useState, useContext, useEffect } from "react";
import Styled, { ThemeContext } from "styled-components";
import PropTypes from "prop-types";
import { PasswordToggled } from "components/molecules/Input";
import Modal from "components/molecules/Modal";
import Button from "components/molecules/Button";
import Text from "components/atoms/Text";
import { changePassword } from "api";
import { useDispatch } from "react-redux";
import { showGeneralNotification } from "store/actions/notification";

function ChangePasswordModal(props) {
  const { openModal, onClose, data } = props;
  const formSchema = {
    password: "",
    re_password: ""
  };
  const dispatch = useDispatch();
  const themeContext = useContext(ThemeContext);
  const [formData, setFormData] = useState(formSchema);
  const [error, setError] = useState("");
  const [btnDisable, setBtnDisable] = useState(true);
  const [loading, setLoading] = useState(false);

  const formatLabel = string => {
    const keys = string.charAt(0).toUpperCase() + string.slice(1);
    return keys;
  };

  const handleChange = event => {
    const { name, value } = event.target;
    setFormData(prevState => ({
      ...prevState,
      [name]: value
    }));
    setTimeout(() => {
      const comparator = name === "password" ? "re_password" : "password";
      if (value.length < 8) {
        setError("Password must be 8 characters or more");
      } else if (value !== formData[comparator]) {
        setError("Confirmation Password didn't match with Password");
      } else if (value === formData[comparator]) {
        setError("");
      }
    }, 1);
  };

  function callChangePassword() {
    const body = {
      ...formData,
      email: data.email,
      username: data.username
    };
    setBtnDisable(true);
    setLoading(true);
    changePassword(data.id, body)
      .then(result => {
        if (result.ok) {
          dispatch(
            showGeneralNotification({
              type: "success",
              desc: "succesfully change password"
            })
          );
        }
        setBtnDisable(false);
        setLoading(false);
        onClose();
      })
      .catch(() => {
        dispatch(
          showGeneralNotification({
            type: "error",
            desc: "failed change password"
          })
        );
        setBtnDisable(false);
        setLoading(false);
      });
  }

  useEffect(() => {
    const values = Object.values(formData);
    const nullIndex = values.findIndex(value => value === "");
    if (nullIndex === -1 && error === "") {
      setBtnDisable(false);
    } else {
      setBtnDisable(true);
    }
  }, [formData, error]);

  return (
    <Modal
      show={openModal}
      className="modal-change-password"
      title="Change Password"
      close={() => onClose()}
      padding="20px"
    >
      <FormWrapper>
        {Object.keys(formData).map(name => (
          <PasswordToggled
            key={name}
            id={name}
            label={
              name === "re_password" ? "Confirm Password" : formatLabel(name)
            }
            onChange={handleChange}
            name={name}
            value={formData[name]}
          />
        ))}
        {error && (
          <Text color={themeContext.inlineError} style={{ margin: "15px 0px" }}>
            {error}
          </Text>
        )}
        <Button
          width="100%"
          style={{ marginTop: "20px" }}
          onClick={() => callChangePassword()}
          isLoading={loading}
          disabled={btnDisable}
        >
          Save
        </Button>
      </FormWrapper>
    </Modal>
  );
}

export default ChangePasswordModal;

ChangePasswordModal.propTypes = {
  openModal: PropTypes.bool.isRequired,
  onClose: PropTypes.func.isRequired,
  data: PropTypes.object
};

ChangePasswordModal.defaultProps = {
  data: {}
};

const FormWrapper = Styled.div`
    display: flex;
    flex-direction: column;
`;

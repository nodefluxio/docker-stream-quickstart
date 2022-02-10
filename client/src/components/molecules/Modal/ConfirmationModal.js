import React, { useEffect, useState } from "react";
import PropTypes from "prop-types";
import Styled from "styled-components";
import Button from "components/molecules/Button";
import Text from "components/atoms/Text";

import Modal from "components/atoms/Modal";
import Input from "components/atoms/Input";

export default function ConfirmationModal(props) {
  const {
    isShown,
    onClose,
    title,
    buttonTitle,
    onConfirm,
    header,
    headerColor,
    subtitle,
    withCode,
    isButtonConfirmLoading,
    withCancelButton
  } = props;
  const [openModal, setOpenModal] = useState(isShown);
  const [code, setCode] = useState("");
  const [input, setInput] = useState("");
  const [isLoading, setIsLoading] = useState(isButtonConfirmLoading);

  useEffect(() => {
    setOpenModal(isShown);
    setInput("");
    if (withCode) {
      setCode(
        Math.random()
          .toString(36)
          .substr(2, 6)
      );
    } else {
      setCode("");
    }
  }, [isShown]);

  useEffect(() => {
    setCode(
      Math.random()
        .toString(36)
        .substr(2, 6)
    );
  }, [withCode]);

  useEffect(() => {
    setIsLoading(isButtonConfirmLoading);
  }, [isButtonConfirmLoading]);

  return (
    <Modal
      show={openModal}
      className="modal-confirmation"
      title={title}
      close={onClose}
    >
      <ConfirmationModalWrapper>
        <Heading color={headerColor}>{header}</Heading>
        <Text align="center" marginBottom="24px">
          {subtitle}
        </Text>
        {withCode && (
          <Wrapper style={{ marginBottom: "40px" }}>
            <Text align="center" marginBottom="24px">
              Please type <b>{code}</b> below to confirm
            </Text>
            <Input
              placeholder="retype code shown above"
              value={input}
              onChange={e => setInput(e.target.value)}
            />
          </Wrapper>
        )}
        {buttonTitle && (
          <Button
            id="confirmation_button"
            width="100%"
            onClick={onConfirm}
            disabled={withCode ? code !== input : false}
            isLoading={isLoading}
          >
            {buttonTitle}
          </Button>
        )}
        {withCancelButton && (
          <Button
            id="close_confirmation_button"
            type="secondary"
            width="100%"
            onClick={onClose}
          >
            cancel
          </Button>
        )}
      </ConfirmationModalWrapper>
    </Modal>
  );
}

ConfirmationModal.propTypes = {
  isShown: PropTypes.bool.isRequired,
  onConfirm: PropTypes.func.isRequired,
  onClose: PropTypes.func.isRequired,
  title: PropTypes.string.isRequired,
  buttonTitle: PropTypes.string.isRequired,
  header: PropTypes.string.isRequired,
  headerColor: PropTypes.string,
  subtitle: PropTypes.string,
  withCode: PropTypes.bool,
  isButtonConfirmLoading: PropTypes.bool,
  withCancelButton: PropTypes.bool
};

ConfirmationModal.defaultProps = {
  headerColor: "#fff",
  subtitle: "All the data will be affected by the change",
  withCode: false,
  isButtonConfirmLoading: false,
  withCancelButton: true
};
const Wrapper = Styled.div``;

const Heading = Styled.h2`
    text-transform: uppercase;
    text-align: center;
    font-style: normal;
    font-weight: 600;
    font-size: 18px;
    line-height: 22px;
    text-align: center;
    margin-bottom: 24px;
    color: ${props => props.color} ;
`;

const ConfirmationModalWrapper = Styled.div`
  padding: 32px;

  button {
      margin-bottom: 16px;
  }
  input {
    margin-bottom: 24px;
    padding-right: 25px;
  }
`;

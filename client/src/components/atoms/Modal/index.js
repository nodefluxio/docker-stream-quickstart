import React, { useContext } from "react";
import PropTypes from "prop-types";
import styled, { ThemeContext } from "styled-components";
import Modal from "react-modal";

import exitIcon from "assets/icon/visionaire/exit-small.svg";

const ModalHeader = styled.div`
  display: flex;
  flex-direction: row;
  height: 30px;
  background: ${props => props.theme.secondary2};
  padding: 13px;
  align-items: center;
  justify-content: space-between;
  font-size: 14px;
  color: ${props => props.theme.white};
  text-transform: capitalize;
  a {
    cursor: pointer;
  }
`;

const ModalContent = styled.div`
  color: ${props => props.theme.white};
  padding: ${props => props.padding};
  word-break: break-word;
  padding: ${props => (props.fullwidth ? "0px" : "15px")};
`;

function ModalComponent(props) {
  const themeContext = useContext(ThemeContext);
  return (
    <Modal
      isOpen={props.show}
      style={{
        content: {
          backgroundColor: themeContext.bg,
          top: "50%",
          left: "50%",
          flex: "1",
          right: "auto",
          bottom: "auto",
          marginRight: "-50%",
          transform: "translate(-50%, -50%)",
          width: props.width || 400,
          position: "relative",
          boxSizing: "border-box",
          display: "flex",
          justifyContent: "center",
          flexDirection: "column",
          zIndex: "999",
          transition: "opacity 200ms ease-in-out",
          outline: "none"
        },
        overlay: {
          backgroundColor: "rgba(8, 3, 17, 0.5)",
          zIndex: "99",
          transition: "opacity 200ms ease-in-out"
        }
      }}
      onRequestClose={props.close}
      shouldCloseOnOverlayClick={true}
      className={props.className}
    >
      <ModalHeader>
        <span>{props.title}</span>
        <a onClick={props.close}>
          <img src={exitIcon} alt="exit" />
        </a>
      </ModalHeader>
      <ModalContent fullwidth={props.fullwidth}>{props.children}</ModalContent>
    </Modal>
  );
}

ModalComponent.propTypes = {
  show: PropTypes.bool,
  children: PropTypes.any,
  close: PropTypes.func,
  submit: PropTypes.func,
  title: PropTypes.string,
  desc: PropTypes.string,
  className: PropTypes.string,
  padding: PropTypes.string,
  fullwidth: PropTypes.bool,
  width: PropTypes.string
};

ModalComponent.defaultProps = {
  show: false,
  children: <div />,
  close: () => {},
  submit: () => {},
  title: "",
  desc: "",
  className: "",
  padding: "",
  fullwidth: false,
  width: ""
};

export default ModalComponent;

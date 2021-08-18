import React from "react";
import Styled from "styled-components";
import PropTypes from "prop-types";
import Button from "components/molecules/Button";

export default function ModalFeedback(props) {
  const {
    desc,
    subDesc,
    status,
    isOpen,
    duration,
    isClose,
    manualClose,
    ...other
  } = props;

  if (isOpen) {
    setTimeout(() => {
      manualClose();
    }, duration);
  }

  return (
    <Wrapper status={status} isOpen={isOpen} {...other}>
      <Panel>
        <Text>{desc}</Text>
        {subDesc && <SubText>{subDesc}</SubText>}
        <Button type="secondary" width="100%" onClick={() => manualClose()}>
          THANKYOU
        </Button>
      </Panel>
    </Wrapper>
  );
}

ModalFeedback.propTypes = {
  desc: PropTypes.string.isRequired,
  subDesc: PropTypes.string,
  status: PropTypes.oneOf(["success", "error"]),
  isOpen: PropTypes.bool,
  duration: PropTypes.number,
  isClose: PropTypes.func,
  manualClose: PropTypes.func
};

ModalFeedback.defaultProps = {
  type: "success",
  desc: "",
  isOpen: false,
  duration: 3000,
  isClose: () => {},
  manualClose: () => {}
};

const Text = Styled.div`
  font-weight: 600;
  font-size: 18px;
  line-height: 22px;
  margin-bottom: 24px;
  text-transform: uppercase;
`;

const SubText = Styled.div`
  font-weight: 400;
  font-size: 14px;
  line-height: 22px;
  margin-bottom: 24px;
`;

const Panel = Styled.div`
  position: relative;
  width: 100%;
  height: 100%;
  background-color: #21153C;
  border-radius: 8px;
  padding: 32px;
  opacity: 1;
  transform: none;
  pointer-events: auto;
  transition-duration: .3s;

`;

const Wrapper = Styled.div`
  position: fixed;
  left: 30px;
  bottom: 30px;
  width: 280px;
  min-height: 172px;
  z-index: 99;
  visibility: hidden;
  
  ${props =>
    props.status === "success"
      ? `
      ${Panel}{
        border: 2px solid #84C041;
      }
    `
      : `
      ${Panel}{
        border: 2px solid #F36B86;
      }
    `}
  
  ${props =>
    props.isOpen &&
    `
        opacity: 1;
        visibility: visible;
    `}
  
  ${props =>
    !props.isOpen &&
    `
      ${Panel}{
        transform: translateY(70px) scale(.5);
        opacity: 0;
        pointer-events: none;
      }
    `}
`;

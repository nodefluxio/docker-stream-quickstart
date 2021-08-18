import React, { useState } from "react";
import PropTypes from "prop-types";
import styled from "styled-components";
import Input from "components/atoms/Input";

const InputPasswordWrapper = styled.div`
  position: relative;
  display: flex;
  flex-direction: row;
  display: block;
  width: 100%;

  input {
    width: 100%;
  }

  .password-toggler {
    display: flex;
    position: absolute;
    right: 10px;
    top: ${props => (props.error ? `15%` : `30%`)};
  }
`;

export default function InputPassword(props) {
  const [visible, setvisible] = useState(false);
  const { error, id } = props;

  return (
    <InputPasswordWrapper error={error}>
      <Input type={visible ? "text" : "password"} {...props}>
        <div
          id={`toggle_password-${id}`}
          onClick={() => setvisible(!visible)}
          className="password-toggler"
        >
          {!visible ? (
            <svg
              xmlns="http://www.w3.org/2000/svg"
              width="20"
              height="16"
              viewBox="0 0 20 16"
            >
              <path
                fill="#FFFFFF"
                fillRule="nonzero"
                d="M19.92 7.6C17.9 2.91 14.1 0 10 0S2.1 2.91.08 7.6a1 1 0 0 0 0 .8C2.1 13.09 5.9 16 10 16s7.9-2.91 9.92-7.6a1 1 0 0 0 0-.8zM10 14c-3.17 0-6.17-2.29-7.9-6C3.83 4.29 6.83 2 10 2s6.17 2.29 7.9 6c-1.73 3.71-4.73 6-7.9 6zm0-10a4 4 0 1 0 0 8 4 4 0 0 0 0-8zm0 6a2 2 0 1 1 0-4 2 2 0 0 1 0 4z"
              />
            </svg>
          ) : (
            <svg
              xmlns="http://www.w3.org/2000/svg"
              width="20"
              height="20"
              viewBox="0 0 20 20"
            >
              <path
                fill="#FFFFFF"
                fillRule="nonzero"
                d="M8.94 4.08A6.93 6.93 0 0 1 10 4c3.18 0 6.17 2.29 7.91 6a15.23 15.23 0 0 1-.9 1.64 1 1 0 0 0 1.7 1.05 15.77 15.77 0 0 0 1.21-2.3 1 1 0 0 0 0-.79C17.9 4.91 14.1 2 10 2a7.78 7.78 0 0 0-1.4.12 1.014 1.014 0 1 0 .34 2v-.04zM1.71.29A1.004 1.004 0 0 0 .29 1.71l3.1 3.09A14.62 14.62 0 0 0 .08 9.6a1 1 0 0 0 0 .8C2.1 15.09 5.9 18 10 18a9.26 9.26 0 0 0 5.05-1.54l3.24 3.25a1 1 0 0 0 1.42 0 1 1 0 0 0 0-1.42l-18-18zm6.36 9.19l2.45 2.45a2 2 0 0 1-2.45-2.45zM10 16c-3.18 0-6.17-2.29-7.9-6a12.09 12.09 0 0 1 2.7-3.79L6.57 8A4 4 0 0 0 12 13.43L13.59 15A7.24 7.24 0 0 1 10 16z"
              />
            </svg>
          )}
        </div>
      </Input>
    </InputPasswordWrapper>
  );
}

InputPassword.propTypes = {
  id: PropTypes.string,
  error: PropTypes.string
};

InputPassword.defaultProps = {
  id: "",
  error: ""
};

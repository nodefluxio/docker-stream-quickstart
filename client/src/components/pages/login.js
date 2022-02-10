/* eslint-disable no-shadow */
/* eslint-disable react/prop-types */
/* eslint-disable camelcase */
import React, { useState, useEffect, useCallback } from "react";
import { useDispatch, useSelector } from "react-redux";
import { emailLogin, resetAuthError } from "store/actions/auth";
import styled, { ThemeProvider } from "styled-components";

import logo from "assets/icon/visionaire/v-logo.svg";
import globe from "assets/images/globe.svg";
import Input, { PasswordToggled } from "components/molecules/Input";
import Text from "components/atoms/Text";
import Button from "components/molecules/Button";

import theme from "theme";

const AnimatedGlobe = styled.div`
  position: absolute;
  z-index: -1;
`;

const LoginContent = styled.div`
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: center;
  height: 100%;
  width: 100%;
`;

const WrapLogin = styled.div`
  color: white;
  margin: 0 auto;
  min-height: 24em;
  height: 100vh;
  opacity: 0;
  animation: LoginfadeIn 0.5s 0s forwards;
  @keyframes LoginfadeIn {
    from {
      opacity: 0;
    }
    to {
      opacity: 1;
    }
  }
  @keyframes LoginfadeOut {
    from {
      opacity: 1;
    }
    to {
      opacity: 0;
    }
  }

  .background-image {
    width: 36.09vw;
    height: 64.16vh;
    position: absolute;
    z-index: -1;
  }
  .form-container {
    width: 20.83vw;
    padding-bottom: 2.96vh;
    background: ${props => props.theme.bg};
    border: 1px solid ${props => props.theme.secondary2};
    box-sizing: border-box;
    border-radius: 8px;
    display: flex;
    flex-direction: column;
    align-items: center;
  }
  .logo {
    width: 1.875vw;
    height: 2.96vh;
    margin-top: 4.72vh;
  }
  .form-title {
    margin-top: 3.33vh;
    font-style: normal;
    font-weight: 600;
    font-size: 18px;
    line-height: 18px;
    text-align: center;
    .bold-font {
      color: white;
    }
    .fade-font {
      color: ${props => props.theme.mint};
    }
  }
  .form-layout {
    width: 17.5vw;
    margin-top: 3.7vh;
    display: flex;
    flex-direction: column;
  }
  .input-divider {
    margin-top: 1.6vh;
  }
  .button-divider {
    margin-top: 2.22vh;
  }
  .bottom-layout {
    display: flex;
    flex-direction: row;
    width: 100%;
    position: absolute;
    bottom: 0;
    padding-bottom: 20px;
    .footer-left {
      width: 33.33vw;
      padding-left: 20px;
    }
    .footer-center {
      width: 33.33vw;
      display: flex;
      justify-content: center;
    }
    .footer-right {
      width: 33.33vw;
      display: flex;
      justify-content: flex-end;
      padding-right: 20px;
    }
  }
`;

export default function Login(props) {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [errMessage, setErrMessage] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [isError, setIsError] = useState("");

  // redux action dispatcher
  const dispatch = useDispatch();

  // redux state
  const userToken = useSelector(state => state.user.token);
  const errorMessage = useSelector(state => state.user.errorMessage);

  useEffect(() => {
    if (userToken) {
      setIsLoading(false);
      dispatch(resetAuthError());
      props.history.push("/");
    }

    if (errorMessage) {
      setIsLoading(false);
      setIsError(true);
      setErrMessage(errorMessage);
      dispatch(resetAuthError());
    }
  });

  function loginHandler() {
    const userData = {
      user_access: email,
      password
    };

    setIsLoading(true);
    dispatch(emailLogin(userData));
  }

  const handleEnter = useCallback(event => {
    if (event.keyCode === 13) {
      loginHandler();
    }
  });

  return (
    <ThemeProvider theme={theme}>
      <WrapLogin>
        <LoginContent>
          <AnimatedGlobe>
            <img src={globe} />
          </AnimatedGlobe>
          <div className="form-container">
            <img className="logo" alt="visionaire logo" src={logo} />
            <div className="form-title">
              <span className="bold-font">SIGN IN TO </span>
              <span className="fade-font">VISIONAIRE</span>
            </div>
            <div className="form-layout">
              <Input
                id="login_field"
                label="Username / Email Address"
                onChange={e => setEmail(e.target.value)}
                onKeyUp={handleEnter}
                name="login_field"
                value={email}
                placeholder="Not set..."
                height="3.7vh"
                error={isError}
              />
              <PasswordToggled
                id="password_field"
                label="Password"
                onChange={e => setPassword(e.target.value)}
                onKeyUp={handleEnter}
                name="password_field"
                value={password}
                placeholder="Not set..."
                height="3.7vh"
                error={errMessage}
              />
              <div className="button-divider" />
              <Button
                id="login_button"
                type="primary"
                isLoading={isLoading}
                onClick={() => loginHandler()}
              >
                Login Now
              </Button>
              {/* <div className="button-divider" />
            <Button id="forgot_password_button" type="secondary">Forgot Password ?</Button> */}
            </div>
          </div>
        </LoginContent>
        <div className="bottom-layout">
          <div className="footer-left">
            <Text size={12}>PRIVACY POLICY</Text>
          </div>
          <div className="footer-center">
            <Text size={12}>POWERED BY NODEFLUX TEKNOLOGI INDONESIA</Text>
          </div>
          <div className="footer-right">
            <Text size={12}>
              COPYRIGHT {new Date().getFullYear()}.ALL RIGHTS RESERVED
            </Text>
          </div>
        </div>
      </WrapLogin>
    </ThemeProvider>
  );
}

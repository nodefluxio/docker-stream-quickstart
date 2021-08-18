import React from "react";
import PropTypes from "prop-types";
import Styled from "styled-components";
import Loader from "assets/icon/visionaire/loader.svg";

function LoadingSpinner(props) {
  if (props.show) {
    return (
      <SpinnerWrapper>
        <img src={Loader} alt="loading-spinner" />
      </SpinnerWrapper>
    );
  }
  return false;
}

const SpinnerWrapper = Styled.div`
  img {
    animation: rotate 2s linear infinite;
  }
  @keyframes rotate {
    100% {
      transform: rotate(360deg);
    }
  }
`;

LoadingSpinner.propTypes = {
  show: PropTypes.bool.isRequired
};
export default LoadingSpinner;

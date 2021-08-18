import React from "react";
import Styled from "styled-components";
import PropTypes from "prop-types";
import ExitIcon from "assets/icon/visionaire/exit.svg";

export default function FullScreen(Component, title, url, RightButtonGroup) {
  return function FullScreenWrapper(props) {
    FullScreenWrapper.propTypes = {
      history: PropTypes.object.isRequired
    };
    return (
      <Wrapper>
        <HeaderWrapper>
          <TitleWrapper>
            <img
              alt="exit-icon"
              src={ExitIcon}
              onClick={() => props.history.push(url)}
            />
            {title}
          </TitleWrapper>
          <RightButtonGroup {...props} />
        </HeaderWrapper>
        <Component {...props} />
      </Wrapper>
    );
  };
}

const Wrapper = Styled.div`
  min-height: 100%;
  width: 100%;
  overflow-y: auto;
`;

const HeaderWrapper = Styled.div`
  width: 100%;
  height: 48px;
  border-bottom: solid 1px ${props => props.theme.secondary2};
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  color: ${props => props.theme.mint};
  font-weight: 600;
  img{
    cursor: pointer;
  }

`;

const TitleWrapper = Styled.div`
`;

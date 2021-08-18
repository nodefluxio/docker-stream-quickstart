import React from "react";
import PropTypes from "prop-types";
import Styled from "styled-components";
import IconPlaceholder from "components/atoms/IconPlaceholder";

import ExitIcon from "assets/icon/visionaire/exit-small.svg";

export default function PageWrapper(props) {
  const { title, close, buttonGroup, children, history } = props;

  return (
    <Wrapper>
      <ControlContainer>
        <AreaActionWrapper>
          {close && (
            <IconPlaceholder onClick={() => history.push("/")}>
              <img src={ExitIcon} alt="close-icon" />
            </IconPlaceholder>
          )}
          {title}
        </AreaActionWrapper>
        {buttonGroup && (
          <ButtonActionWrapper>{buttonGroup}</ButtonActionWrapper>
        )}
      </ControlContainer>
      {children}
    </Wrapper>
  );
}

PageWrapper.propTypes = {
  title: PropTypes.string.isRequired,
  buttonGroup: PropTypes.oneOfType([
    PropTypes.arrayOf(PropTypes.node),
    PropTypes.node
  ]),
  close: PropTypes.bool,
  history: PropTypes.object.isRequired,
  children: PropTypes.oneOfType([
    PropTypes.arrayOf(PropTypes.node),
    PropTypes.node
  ]).isRequired
};

PageWrapper.defaultProps = {
  buttonGroup: null,
  close: false
};

const Wrapper = Styled.div`
  position: relative;
  height: 100%;
  .mr-16 {
    margin-right: 16px;
  }
  .mr-4 {
    margin-right: 4px;
  }
`;

const ControlContainer = Styled.div`
  display: flex;
  flex-direction: row;
  width: 100%;
  justify-content: space-between;
  align-items: center;
  height: 64px;
  border-bottom: 1px solid #372463;
  .plus-icon{
      margin-right: 10px;
  }
`;

const AreaActionWrapper = Styled.div`
  display: flex;
  flex-direction: row;
  align-items: center;
  text-transform: uppercase;
  color: ${props => props.theme.mint};
  font-weight: 600;
  font-size: 18px;
  padding-left: 20px;
`;

const ButtonActionWrapper = Styled.div`
  display: flex;
  flex-direction: row;
  margin-right: 4px;
  button:not:last-child {
    margin-right: 16px;
  }
`;

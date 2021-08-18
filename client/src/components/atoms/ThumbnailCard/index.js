import React from "react";
import PropTypes from "prop-types";
import Styled from "styled-components";

export default function Thumbnail(props) {
  const { url, title, actions, onClick, ...other } = props;
  return (
    <Wrapper {...other}>
      <Image url={url} onClick={onClick} />
      <Content>
        {title && <Title>{title}</Title>}
        {actions}
      </Content>
    </Wrapper>
  );
}

Thumbnail.propTypes = {
  url: PropTypes.string.isRequired,
  title: PropTypes.string,
  actions: PropTypes.element,
  onClick: PropTypes.func
};

Thumbnail.defaultProps = {
  actions: null,
  onClick: null,
  title: ""
};

const Image = Styled.div`
  width: 100%;
  padding-top: 100%;
  background-image: url(${props => props.url});
  background-size: cover;
  background-repeat: no-repeat;
  background-position: center;
  ${({ onClick }) =>
    onClick &&
    `
    cursor: pointer;
  `}
`;

const Title = Styled.div`
  font-family: Barlow;
  font-style: normal;
  font-weight: 500;
  font-size: 12px;
  line-height: 14px;
  color: #E5E5E5;
`;

const Content = Styled.div`
  display: grid;
  grid-template-columns: auto 40px;
  grid-gap: 10px;
  align-items: center;
  padding-left: 5px;
  height: 40px;
`;

const Wrapper = Styled.div`
  position: relative;
  max-width: 222px;
  ${Image}{
    margin-bottom: 9px;
  }
`;

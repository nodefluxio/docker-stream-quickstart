import React from "react";
import Styled from "styled-components";
import PropTypes from "prop-types";

import SearchDBIcon from "assets/icon/visionaire/search_db.svg";

export default function LinkButton(props) {
  const { onClick } = props;
  const item = {
    label: "Search Dukcapil",
    icon: SearchDBIcon
  };
  return (
    <Button key={item.label} id={item.label} onClick={onClick}>
      <img src={item.icon} alt={item.label} />
      <span>{item.label}</span>
    </Button>
  );
}

LinkButton.propTypes = {
  onClick: PropTypes.func
}

const Button = Styled.button`
  border: none;
  background: none;
  outline: none;
  cursor: pointer;
  display: flex;
  align-items: center;
  margin-right: 36px;

  :last-child{
    margin-right: 18px;
  }
  
  img,svg{
    margin-right: 16px;
    max-width: 15px;
  }
  span{
    font-family: Barlow;
    font-style: normal;
    font-weight: 500;
    font-size: 12px;
    line-height: 14px;
    text-transform: uppercase;
    color: ${props => props.theme.mercury};
  }

  :hover {
    background: ${props => props.theme.blueGem};
    border-bottom: 3px solid ${props => props.theme.mint};
  }
`;

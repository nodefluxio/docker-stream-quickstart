import React from "react";
import PropTypes from "prop-types";
import { ContextMenu, MenuItem } from "react-contextmenu";
import styled from "styled-components";

export function MenuWrapper({ id, children }) {
  return <StyledContextMenu id={id}>{children}</StyledContextMenu>;
}

MenuWrapper.propTypes = {
  id: PropTypes.string.isRequired,
  children: PropTypes.any.isRequired
};

export function ItemWrapper({ onClick, data, children }) {
  return (
    <StyledMenuItem onClick={onClick} data={data}>
      {children}
    </StyledMenuItem>
  );
}

ItemWrapper.propTypes = {
  onClick: PropTypes.func.isRequired,
  data: PropTypes.object,
  children: PropTypes.any.isRequired
};

ItemWrapper.defaultProps = {
  data: {}
};

const StyledContextMenu = styled(ContextMenu)`
  min-width: 216px;
`;

const StyledMenuItem = styled(MenuItem)`
  background-color: #372463;
  cursor: pointer;
  padding: 12px;
  :hover {
    background-color: #4c12a1;
  }
`;

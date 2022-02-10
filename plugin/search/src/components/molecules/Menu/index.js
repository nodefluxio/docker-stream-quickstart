import React, { Fragment } from "react";
import PropTypes from "prop-types";

import Dropdown, { Menu } from "components/atoms/Dropdown";
import LinkButton from "components/atoms/LinkButton";

import SearchDBIcon from "assets/icon/visionaire/folder.svg";

export default function HeaderMenu({ url, history, role }) {
  const dropdownMenu = (
    <Fragment>
      {role === "superadmin" && (
        <Menu onClick={() => history.push(`${url}/person`)}>Search Person</Menu>
      )}
      <Menu onClick={() => history.push(`${url}/vehicle`)}>Search Vehicle</Menu>
    </Fragment>
  );
  return (
    <Dropdown overlay={dropdownMenu} width="150px">
      <LinkButton>
        <img src={SearchDBIcon} alt="user" />
        <span>Search Database</span>
      </LinkButton>
    </Dropdown>
  );
}

HeaderMenu.propTypes = {
  history: PropTypes.object.isRequired,
  url: PropTypes.string.isRequired,
  role: PropTypes.string.isRequired
};

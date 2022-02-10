import React, { Fragment, useEffect, useState } from "react";
import { useDispatch } from "react-redux";
import Styled from "styled-components";
import PropTypes from "prop-types";
import { Link } from "react-router-dom";

import Dropdown, { Menu } from "components/atoms/Dropdown";

import LinkButton from "components/atoms/LinkButton";

import { PLUGIN_HOST, PLUGIN_NAME } from "config";
import { logOut } from "store/actions/auth";

import vLogo from "assets/icon/visionaire/visionaire.svg";
import visionarePlatform from "assets/icon/visionaire/visionare-platform.svg";
import enrollmentList from "assets/icon/visionaire/enrollment.svg";
import SearchDBIcon from "assets/icon/visionaire/search_db.svg";
import ChevronIcon from "assets/icon/visionaire/arrow-up.svg";
import DropdownHover, { MenuHover } from "components/atoms/DropdownHover";
import UserIcon from "assets/icon/visionaire/Union.svg";
import { getCookie } from "helpers/cookies";
import parseJwt from "helpers/parseJWT";
import Plugin from "./PluginContainer";

function generateURL(pluginName) {
  const split = pluginName.split(/(?=[A-Z])/);
  return split.join("-").toLowerCase();
}

function Header({ history }) {
  const dispatch = useDispatch();
  const [role, setRole] = useState("operator");

  useEffect(() => {
    const accessToken = getCookie("access_token");
    if (accessToken) {
      const userCookies = parseJwt(accessToken);
      setRole(userCookies.role);
    }
  }, []);

  const menu = [
    {
      label: "Event History",
      url: "/event-history",
      icon: SearchDBIcon
    },
    {
      label: "enrollment",
      url: "",
      icon: enrollmentList,
      submenu: [
        {
          label: "FACE",
          url: "/enrollment",
          icon: enrollmentList
        },
        {
          label: "VEHICLE",
          url: "/vehicle",
          icon: enrollmentList
        }
      ]
    }
  ];

  function onLogOut() {
    dispatch(logOut);
    history.push(`/login`);
  }

  const userMenu = (
    <Fragment>
      <Menu onClick={() => history.push("/account")}>ACCOUNT</Menu>
      <Menu onClick={() => history.push("/license")}>LICENSE</Menu>
      <Menu onClick={() => onLogOut()}>LOG OUT</Menu>
    </Fragment>
  );

  const system = PLUGIN_HOST
    ? {
        name: generateURL(PLUGIN_NAME),
        url: `${PLUGIN_HOST}/remoteEntry.js`,
        scope: `${PLUGIN_NAME}`,
        module: "./Button"
      }
    : {};
  const submenu = items => (
    <MenuWrapper>
      {items.map(item => (
        <MenuHover
          key={item.label}
          className="menu"
          onClick={() => history.push(item.url)}
        >
          <img src={enrollmentList} alt="check-icon" />
          <span>{item.label}</span>
        </MenuHover>
      ))}
    </MenuWrapper>
  );
  const navBlock = () => (
    <Navigation>
      <Fragment>
        <Plugin
          system={system}
          url={`/plugin/${system.name}`}
          history={history}
          role={role}
        />
        {menu.map(item =>
          item.submenu === undefined ? (
            <LinkButton
              key={item.label}
              id={item.label}
              onClick={() => history.push(item.url)}
            >
              <img src={item.icon} alt={item.label} />
              <span>{item.label}</span>
            </LinkButton>
          ) : (
            <DropdownHover
              key={item.label}
              id="grid_view_icon"
              width="135px"
              overlay={submenu(item.submenu)}
              className="button-control"
            >
              <SubmenuButton>
                <img alt="grid-icon" src={enrollmentList} />
                <span>enrollment</span>
                <img
                  className="arrow-icon"
                  alt="arrow-icon"
                  src={ChevronIcon}
                />
              </SubmenuButton>
            </DropdownHover>
          )
        )}
        <Dropdown overlay={userMenu} width="150px">
          <LinkButton>
            <img src={UserIcon} alt="user" />
            <span>Account</span>
          </LinkButton>
        </Dropdown>
      </Fragment>
    </Navigation>
  );

  return (
    <Wrapper className="Header">
      <Panel>
        <Link to={"/"}>
          <Block>
            <img src={vLogo} alt="VisionAIre" />
            <img src={visionarePlatform} alt="VisionAIre Platform" />
          </Block>
        </Link>
        <Block>{navBlock()}</Block>
      </Panel>
    </Wrapper>
  );
}

Header.propTypes = {
  history: PropTypes.object.isRequired
};

export default Header;

const Navigation = Styled.div`
  display: flex;
`;

const Wrapper = Styled.div`
  position: fixed;
  height: 48px;
  width: 100%;
  z-index: 1;
`;

const Panel = Styled.nav`
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 48px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  z-index: 99;
  background-color: ${props => props.theme.bg};
  border-bottom: 1px solid ${props => props.theme.secondary2};
`;

const Block = Styled.div`
  width: max-content;
  display: flex;
  height: 100%;
`;

const MenuWrapper = Styled.div`
  text-transform: ${props => props.textTransform || `none`} ;
  font-weight: bold;
  font-size: 12px;
  line-height: 14px;
  color: #E5E5E5;

  .menu {
    display: flex;
    flex-direction: row;
    justify-content: left;
    padding-right: 17px;
  }
  .menu img {
    margin-right: 20px;
  }
`;

const SubmenuButton = Styled.div`
  border: none;
  background: none;
  outline: none;
  cursor: pointer;
  display: flex;
  align-items: center;
  
  padding: 15px 15px;
  img:first-child,svg:first-child {
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

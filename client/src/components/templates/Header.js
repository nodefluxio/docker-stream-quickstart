import React, { Fragment } from "react";
import Styled from "styled-components";
import PropTypes from "prop-types";
import { Link } from "react-router-dom";

import LinkButton from "components/atoms/LinkButton";

import { PLUGIN_HOST, PLUGIN_NAME } from "config";

import vLogo from "assets/icon/visionaire/visionaire.svg";
import visionarePlatform from "assets/icon/visionaire/visionare-platform.svg";
import enrollmentList from "assets/icon/visionaire/enrollment.svg";
import QubicIcon from "assets/icon/visionaire/qubic.svg";
import SearchDBIcon from "assets/icon/visionaire/search_db.svg";
import Plugin from "./PluginContainer";

function generateURL(pluginName) {
  const split = pluginName.split(/(?=[A-Z])/);
  return split.join("-").toLowerCase();
}

function Header({ history }) {
  const menu = [
    {
      label: "Event History",
      url: "/event-history",
      icon: SearchDBIcon
    },
    {
      label: "enrollment",
      url: "/enrollment",
      icon: enrollmentList
    },
    {
      label: "Camera",
      url: "/camera",
      icon: QubicIcon
    }
  ];

  const system = PLUGIN_HOST
    ? {
        name: generateURL(PLUGIN_NAME),
        url: `${PLUGIN_HOST}/remoteEntry.js`,
        scope: `${PLUGIN_NAME}`,
        module: "./Button"
      }
    : {};

  return (
    <Wrapper className="Header">
      <Panel>
        <Link to={"/"}>
          <Block>
            <img src={vLogo} alt="VisionAIre" />
            <img src={visionarePlatform} alt="VisionAIre Platform" />
          </Block>
        </Link>
        <Block>
          <Navigation>
            <Fragment>
              <Plugin
                system={system}
                onClick={() => history.push(`/plugin/${system.name}`)}
              />
              {menu.map(item => (
                <LinkButton
                  key={item.label}
                  id={item.label}
                  onClick={() => history.push(item.url)}
                >
                  <img src={item.icon} alt={item.label} />
                  <span>{item.label}</span>
                </LinkButton>
              ))}
            </Fragment>
          </Navigation>
        </Block>
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

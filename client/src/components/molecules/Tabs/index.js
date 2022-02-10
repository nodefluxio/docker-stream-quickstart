import React, { Component, Fragment } from "react";
import PropTypes from "prop-types";
import styled from "styled-components";

const WrapTabs = styled.div`
  width: 100%;
  height: 100%;
`;

const WrapTab = styled.div`
  width: fit-content;
  height: auto;
  color: white;
  display: flex;
  box-sizing: border-box;
`;

const TabButton = styled.div`
  flex: 1;
  height: ${props => props.height || `60px`};
  font-size: 14px;
  font-weight: normal;
  font-style: normal;
  font-stretch: normal;
  line-height: normal;
  letter-spacing: 1.3px;
  cursor: pointer;
  color: rgba(34, 8, 78, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  transition: color 0.3s;
  ${({ active }) => active && `color: ${props => props.theme.primary1}`}
`;

const HR = styled.div`
  width: ${({ length }) => 100 / length}%;
  height: 2px;
  margin-left: ${({ length, id }) =>
    id === "1" ? 0 : (100 / length) * (id - 1)}%;
  background-color: ${props => props.theme.primary3};
  margin-bottom: 20px;
  transition: margin-left 0.3s cubic-bezier(0.645, 0.045, 0.355, 1);
`;

const TabContent = styled.div`
  width: 100%;
  height: 100%;
  visibility: visible;
  opacity: 1;
  ${({ show }) =>
    !show &&
    `
        visibility: hidden;
        opacity: 0;
        height: 0;
        overflow: hidden;
  `};
  transition: visibility 0s, opacity 0.5s linear;
`;

const TabWrapContent = styled.div`
  width: 100%;
  height: 100%;
  box-sizing: border-box;
`;

export default class Tabs extends Component {
  static propTypes = {
    children: PropTypes.any,
    style: PropTypes.object,
    customTab: PropTypes.bool,
    activeTab: PropTypes.number
  };

  constructor(props) {
    super(props);

    this.state = {
      currentIndex: this.props.activeTab || 1
    };
  }

  componentDidUpdate(prevProps) {
    if (this.props.activeTab !== prevProps.activeTab) {
      this.setState({
        currentIndex: this.props.activeTab
      });
    }
  }

  switchTab = idx => () => {
    this.setState({
      currentIndex: idx
    });
  };

  render() {
    const { children, style, customTab } = this.props;
    const { currentIndex } = this.state;
    const childrenWithProps = React.Children.map(children, child =>
      React.cloneElement(child, {
        idb: currentIndex,
        onClick: this.switchTab(child.props.idt),
        customTab
      })
    );
    const childrenTab = children.map((val, idb) => (
      <TabContent key={idb} show={val.props.idt === currentIndex}>
        {val.props.children}
      </TabContent>
    ));
    return (
      <WrapTabs style={style}>
        <WrapTab>{childrenWithProps}</WrapTab>
        {customTab ? null : <HR length={children.length} id={currentIndex} />}
        <TabWrapContent>{childrenTab}</TabWrapContent>
      </WrapTabs>
    );
  }
}

const tabDefaultProps = {
  style: {},
  onClick: () => {},
  idb: 0,
  tab: null,
  customTab: false
};
const tabPropTypes = {
  tab: PropTypes.any,
  idb: PropTypes.number,
  idt: PropTypes.any.isRequired,
  style: PropTypes.object,
  onClick: PropTypes.func,
  customTab: PropTypes.bool
};

export const Tab = ({ tab, idb, idt, onClick, style, customTab }) => {
  // idt = idTab
  const customTabElement = React.Children.map(tab, child =>
    React.cloneElement(child, {
      active: idb === idt,
      onClick,
      id: idt,
      style
    })
  );
  return (
    <Fragment>
      {customTab === false || customTab == null ? (
        <TabButton
          active={idb === idt}
          id={idt}
          onClick={onClick}
          style={style}
        >
          {tab}
        </TabButton>
      ) : (
        customTabElement
      )}
    </Fragment>
  );
};
Tab.defaultProps = tabDefaultProps;
Tab.propTypes = tabPropTypes;

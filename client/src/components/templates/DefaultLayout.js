import React from "react";
import { connect, useDispatch } from "react-redux";
import PropTypes from "prop-types";
import Styled, { ThemeProvider } from "styled-components";
import { Redirect, Route, Switch } from "react-router-dom";
import routes from "router";
import theme from "theme";

import Modal from "components/molecules/Modal";
import { closeNotificationGeneral } from "store/actions/notification";
import Header from "./Header";
import Footer from "./Footer";
import requireAuth from "./AuthRoute";

function DefaultLayout(props) {
  const { history, popupFeedback } = props;
  const dispatch = useDispatch();
  const closeNotifFeedback = () => {
    dispatch(closeNotificationGeneral());
  };
  return (
    <ThemeProvider theme={theme}>
      <Main>
        <Layout>
          <Header history={history} />
          <Switch>
            {routes.map((route, idx) =>
              route.component ? (
                <Route
                  key={idx}
                  path={route.path}
                  exact={route.exact}
                  render={propsComponent => (
                    <Content fullscreen={route.fullscreen}>
                      <route.component {...propsComponent} />
                    </Content>
                  )}
                />
              ) : null
            )}
            <Redirect from="*" to="/404" />
          </Switch>
          <Footer history={history} />
        </Layout>
        <Modal
          desc={popupFeedback.desc || ""}
          subDesc={popupFeedback.subDesc || ""}
          type="feedback"
          status={popupFeedback.type}
          isOpen={popupFeedback.show}
          manualClose={() => closeNotifFeedback()}
          duration={3000}
        />
      </Main>
    </ThemeProvider>
  );
}

DefaultLayout.propTypes = {
  history: PropTypes.object.isRequired,
  popupFeedback: PropTypes.object.isRequired
};

function mapStateToProps(state) {
  return {
    popupFeedback: state.popupFeedback
  };
}

export default requireAuth(connect(mapStateToProps)(DefaultLayout));

const Main = Styled.div`
  position: relative;
  color: white;
  overflow: hidden;
  height: 100%;
  width: 100%;
`;

const Layout = Styled.div`
  display: flex;
  flex-direction: column;
  max-height: 100%;
  height: 100%;
  min-height: 100%;
`;

const Content = Styled.div`
  width: 100%;
  z-index: ${props => (props.fullscreen ? `99` : `0`)};
  padding-top: ${props => (props.fullscreen ? `0px` : `51px`)};
  background: ${props => props.theme.bg};
  max-height: ${props => (props.fullscreen ? `100%` : `calc(100% - 101px)`)};
  min-height: ${props => (props.fullscreen ? `100%` : `calc(100% - 101px)`)}; 
`;

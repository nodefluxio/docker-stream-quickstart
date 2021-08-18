import React from "react";
import PropTypes from "prop-types";
import Styled, { ThemeProvider } from "styled-components";
import { Redirect, Route, Switch } from "react-router-dom";
import routes from "router";
import theme from "theme";

function DefaultLayout() {
  return (
    <ThemeProvider theme={theme}>
      <Main>
        <Layout>
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
        </Layout>
      </Main>
    </ThemeProvider>
  );
}

DefaultLayout.propTypes = {
  history: PropTypes.object.isRequired
};

export default DefaultLayout;

const Main = Styled.div`
  position: relative;
  color: white;
  min-height: 100vh;
  overflow: hidden;
`;

const Layout = Styled.div`
  display: flex;
  height: calc(100vh);
  min-height: calc(100vh);
`;

const Content = Styled.div`
  width: 100%;
  z-index: ${props => (props.fullscreen ? `99` : `0`)};
  padding-top: 0px;
  background: ${props => props.theme.bg};
  min-height: 100%;
`;

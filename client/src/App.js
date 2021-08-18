import React from "react";
import { BrowserRouter, Switch, Route } from "react-router-dom";
import { Provider } from "react-redux";

import { Page404, Login } from "components/pages";

import { DefaultLayout } from "components/templates";
import { initializeStore } from "./store/store";

const store = initializeStore();
function App() {
  return (
    <Provider store={store}>
      <BrowserRouter>
        <Switch>
          <Route path="/404" render={props => <Page404 {...props} />} />
          <Route path="/login" render={props => <Login {...props} />} />
          <Route path="/" render={props => <DefaultLayout {...props} />} />
        </Switch>
      </BrowserRouter>
    </Provider>
  );
}

export default App;

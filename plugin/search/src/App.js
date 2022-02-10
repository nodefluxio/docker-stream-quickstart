import React from "react";
import { BrowserRouter, Switch, Route } from "react-router-dom";

import { DefaultLayout } from "components/templates";

function App() {
  return (
    <BrowserRouter>
      <Switch>
        <Route path="/" render={props => <DefaultLayout {...props} />} />
      </Switch>
    </BrowserRouter>
  );
}

export default App;

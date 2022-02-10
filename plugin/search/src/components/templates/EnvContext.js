/* eslint-disable no-undef */
/* eslint-disable camelcase */
import React, { createContext } from "react";
import useFetchJson from "helpers/useFetchJson";

export const EnvContext = createContext();

export default function EnvWrapper(Component) {
  function requireEnv(props) {
    const { data, loading } = useFetchJson(
      `${__webpack_public_path__}env-config.json`
    );
    return (
      !loading && (
        <EnvContext.Provider value={data.SERVER_API}>
          <Component {...props} />
        </EnvContext.Provider>
      )
    );
  }
  return requireEnv;
}

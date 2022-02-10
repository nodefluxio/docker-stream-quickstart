// page for holding plugin available
import React from "react";
import PropTypes from "prop-types";
import { PLUGIN_HOST } from "config";

import Plugin from "components/templates/PluginContainer";

export default function PluginPage({ history }) {
  function capitalizeFirstLetter(string) {
    return string.charAt(0).toUpperCase() + string.slice(1);
  }

  const paths = history.location.pathname.split("/");
  const system =
    PLUGIN_HOST && paths.length >= 3
      ? {
          url: `${PLUGIN_HOST}/remoteEntry.js`,
          scope: capitalizeFirstLetter(paths[2]),
          module: `./${capitalizeFirstLetter(paths[3])}`
        }
      : {};
  return <Plugin system={system} history={history} />;
}

PluginPage.propTypes = {
  history: PropTypes.object.isRequired
};

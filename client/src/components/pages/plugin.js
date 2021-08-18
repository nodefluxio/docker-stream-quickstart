// page for holding plugin available
import React from "react";
import { PLUGIN_HOST, PLUGIN_NAME } from "config";

import Plugin from "components/templates/PluginContainer";

export default function PluginPage() {
  const system = PLUGIN_HOST
    ? {
        url: `${PLUGIN_HOST}/remoteEntry.js`,
        scope: PLUGIN_NAME,
        module: "./Page"
      }
    : {};
  return <Plugin system={system} />;
}

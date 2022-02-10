/* eslint-disable no-undef */
/* eslint-disable no-console */
/* eslint-disable consistent-return */
import React from "react";
import PropTypes from "prop-types";
import { connect } from "react-redux";

function loadComponent(scope, module) {
  return async () => {
    // Initializes the share scope. This fills it with known provided modules from this build and all remotes
    await __webpack_init_sharing__("default");

    const container = window[scope]; // or get the container somewhere else
    // Initialize the container, it may provide shared modules
    await container.init(__webpack_share_scopes__.default);
    const factory = await window[scope].get(module);
    const Module = factory();
    return Module;
  };
}

const useDynamicScript = args => {
  const [ready, setReady] = React.useState(false);
  const [failed, setFailed] = React.useState(false);

  React.useEffect(() => {
    if (!args.url) {
      return;
    }

    const element = document.createElement("script");

    element.src = args.url;
    element.type = "text/javascript";
    element.async = true;

    setReady(false);
    setFailed(false);

    element.onload = () => {
      console.log(`Dynamic Script Loaded: ${args.url}`);
      setReady(true);
    };

    element.onerror = () => {
      console.error(`Dynamic Script Error: ${args.url}`);
      setReady(false);
      setFailed(true);
    };

    document.head.appendChild(element);

    return () => {
      console.log(`Dynamic Script Removed: ${args.url}`);
      document.head.removeChild(element);
    };
  }, [args.url]);

  return {
    ready,
    failed
  };
};

function Plugin(props) {
  const { ready, failed } = useDynamicScript({
    url: props.system && props.system.url
  });

  if (!props.system) {
    return null;
  }

  if (!ready) {
    return null;
  }

  if (failed) {
    return null;
  }

  const Component = React.lazy(
    loadComponent(props.system.scope, props.system.module)
  );

  return (
    <React.Suspense fallback="Loading System">
      <Component {...props} store={props.store} />
    </React.Suspense>
  );
}

Plugin.propTypes = {
  system: PropTypes.object.isRequired,
  componentProps: PropTypes.object,
  store: PropTypes.object.isRequired
};

function mapStateToProps(state) {
  return {
    store: state
  };
}

export default connect(mapStateToProps)(Plugin);

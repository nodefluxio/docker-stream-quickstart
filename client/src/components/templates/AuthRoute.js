import React, { useEffect } from "react";
import { withRouter } from "react-router-dom";
import { connect, useDispatch } from "react-redux";
import { getCookie } from "helpers/cookies";
import PropTypes from "prop-types";
import { logOut } from "../../store/actions/auth";

export default function requireAuth(Component) {
  function AuthenticatedComponent(props) {
    // redux action dispatcher
    const dispatch = useDispatch();

    useEffect(() => {
      const accessToken = getCookie("access_token");
      // handle when user logout, user clear cache, or network error
      if (!accessToken || accessToken === undefined) {
        dispatch(logOut());
        props.history.push(`/login`);
      }
    });

    return <Component {...props} />;
  }

  AuthenticatedComponent.propTypes = {
    router: PropTypes.object,
    location: PropTypes.object,
    user: PropTypes.object,
    history: PropTypes.object
  };

  function mapStateToProps(state) {
    return {
      user: state.user
    };
  }

  return withRouter(connect(mapStateToProps)(AuthenticatedComponent));
}

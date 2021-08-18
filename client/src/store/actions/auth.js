import Cookies from "universal-cookie";
import { setCookie, getCookie, removeCookie } from "helpers/cookies";
import { userLogin, refreshAccessToken } from "api/auth";
import {
  USER_LOGIN,
  USER_LOGOUT,
  REFRESH_ACCESS_TOKEN,
  USER_FAILED_AUTH,
  USER_UPDATE_DATA,
  USER_RESET_AUTH_ERROR
} from "store/actionType";

const cookie = new Cookies();

export function logOut() {
  removeCookie();
  return dispatch => {
    dispatch({
      type: USER_LOGOUT
    });
  };
}

export function getAccessTokenByRefreshToken() {
  return async dispatch => {
    try {
      const refreshToken = await getCookie("refresh_token");

      if (!refreshToken) {
        return dispatch(logOut());
      }

      const newAccessToken = await refreshAccessToken();

      if (!newAccessToken) {
        return dispatch(logOut());
      }
      setCookie(newAccessToken);
      return dispatch({
        type: REFRESH_ACCESS_TOKEN,
        token: newAccessToken.access_token
      });
    } catch (error) {
      return dispatch({
        type: USER_FAILED_AUTH,
        message: "Service is Unavailable"
      });
    }
  };
}

export function emailLogin(authData) {
  return async dispatch => {
    try {
      const userData = await userLogin(authData);
      setCookie(userData);

      return dispatch({ type: USER_LOGIN, token: userData.access_token });
    } catch (error) {
      if (error.response) {
        if (error.response.status === 422 || error.response.status === 400) {
          return dispatch({
            type: USER_FAILED_AUTH,
            message: "Email or Password is Incorrect"
          });
        }

        if (error.response.status === 405 || error.response.status === 503) {
          return dispatch({
            type: USER_FAILED_AUTH,
            message: "Service is Unavailable"
          });
        }
      }

      // for development only
      return dispatch({ type: USER_FAILED_AUTH, message: "Something Error" });
    }
  };
}

export function updateUserInfo(userData) {
  return dispatch => {
    dispatch({
      type: USER_UPDATE_DATA,
      id: userData.id,
      email: userData.email,
      fullname: userData.fullname,
      role: userData.role,
      username: userData.username,
      versionClient: userData.versionClient,
      versionApi: userData.versionApi
    });
  };
}

export function getAccessToken() {
  return async dispatch => {
    try {
      const accessToken = await getCookie("access_token");

      if (accessToken) {
        return dispatch({
          type: REFRESH_ACCESS_TOKEN,
          token: accessToken
        });
      }

      return dispatch(getAccessTokenByRefreshToken());
    } catch (error) {
      return dispatch({
        type: REFRESH_ACCESS_TOKEN,
        token: null
      });
    }
  };
}

export function getToken() {
  return cookie.get("access_token");
}

export function resetAuthError() {
  return dispatch => {
    dispatch({
      type: USER_RESET_AUTH_ERROR
    });
  };
}

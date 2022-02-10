import Cookies from "universal-cookie";
import { setCookie, removeCookie } from "helpers/cookies";
import { userLogin } from "api/auth";
import {
  USER_LOGIN,
  USER_LOGOUT,
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

/* eslint-disable import/no-cycle */
import fetch from "helpers/fetch";
import { REACT_APP_API_AUTH } from "config";
import { getCookie } from "helpers/cookies";

function userLogin(userData) {
  const url = `${REACT_APP_API_AUTH}/auth/token`;

  return fetch(url, "post", userData)
    .then(result => result)
    .catch(error => {
      throw error;
    });
}

async function refreshAccessToken() {
  try {
    const url = `${REACT_APP_API_AUTH}/refresh-token`;
    const refreshToken = getCookie("refresh_token");
    const newAccessToken = await fetch(url, "post", null, refreshToken);

    return newAccessToken;
  } catch (error) {
    // eslint-disable-next-line no-console
    console.log(error);
    throw error;
  }
}

export { userLogin, refreshAccessToken };

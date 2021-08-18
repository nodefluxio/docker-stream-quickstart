import Cookies from "universal-cookie";
import dayjs from "dayjs";

const cookie = new Cookies();
const hostDomain = window.location.hostname;

function getCookie(cname) {
  return cookie.get(cname);
}

function setCookie(userData) {
  const ACCESS_TOKEN_EXPIRED_TIME = 1;
  const REFRESH_TOKEN_EXPIRED_TIME = 2;

  const oldDateObj = new Date();

  const accessTokenExp = dayjs(oldDateObj)
    .add(ACCESS_TOKEN_EXPIRED_TIME, "d")
    .toDate();
  const refreshTokenExp = dayjs(oldDateObj)
    .add(REFRESH_TOKEN_EXPIRED_TIME, "d")
    .toDate();

  cookie.set("access_token", userData.access_token, {
    path: "/",
    expires: accessTokenExp,
    domain: hostDomain
  });

  cookie.set("refresh_token", userData.refresh_token, {
    path: "/",
    expires: refreshTokenExp,
    domain: hostDomain
  });
}

function removeCookie() {
  cookie.remove("refresh_token", {
    path: "/",
    domain: hostDomain
  });

  cookie.remove("access_token", {
    path: "/",
    domain: hostDomain
  });
}

function getHeaders(isFormData) {
  const token = getCookie("access_token") || null;

  return {
    Authorization: `Bearer ${token}`,
    ...isFormData
  };
}

export { getHeaders, getCookie, setCookie, removeCookie };

import axios from "axios";
import { getCookie } from "helpers/cookies";

const checkStatus = response => {
  if (response.status === 401) {
    window.location.href("/");
  }
  return response.data;
};

const constructHeader = async (isFormData, token) => {
  const accessToken = getCookie("access_token") || null;

  if (token) {
    return {
      Authorization: `Bearer ${token}`,
      ...isFormData
    };
  }

  return {
    Authorization: `Bearer ${accessToken}`,
    ...isFormData
  };
};

const fetch = async (url, method, data, responseType, token) => {
  const isFormData = data instanceof FormData && {
    "Content-Type": "multipart/form-data"
  };

  const Header = {
    headers: await constructHeader(isFormData, token)
  };

  const options = {
    url,
    method,
    data,
    responseType,
    ...Header
  };

  axios.interceptors.response.use(
    response =>
      response.headers["content-type"] !== "text/html; charset=utf-8"
        ? response
        : Promise.reject(response),
    error => Promise.reject(error)
  );

  const reqData = await axios(options);
  await checkStatus(reqData);
  const result = reqData.data;
  return result;
};

export default fetch;

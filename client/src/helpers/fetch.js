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

  const reqData = await axios(options);
  await checkStatus(reqData);
  const result = reqData.data;
  return result;
};

export default fetch;

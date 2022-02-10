import axios from "axios";

const checkStatus = response => {
  if (response.status === 401) {
    window.location.href("/");
  }
  return response.data;
};

const fetch = async (url, method, data, responseType) => {
  const isFormData = data instanceof FormData && {
    "Content-Type": "multipart/form-data"
  };

  const options = {
    url,
    method,
    data,
    responseType,
    ...isFormData
  };

  const reqData = await axios(options);
  await checkStatus(reqData);
  const result = reqData.data;
  return result;
};

export default fetch;

import fetch from "helpers/fetch";

export function getCitizen(SERVER_API, query) {
  const url = `${SERVER_API}/api/search/polri/nik${query}`;
  return fetch(url, "get")
    .then(result => result)
    .catch(error => {
      throw error;
    });
}

export function getToken(SERVER_API, data) {
  const url = `${SERVER_API}/api/search/polri/face-similarity/token`;
  return fetch(url, "post", data)
    .then(result => result)
    .catch(error => {
      throw error;
    });
}

export async function getResultByToken(SERVER_API, token) {
  const url = `${SERVER_API}/api/search/polri/face-similarity/results?token=${token}`;
  const fetchResult = await fetch(url, "get")
    .then(result => result)
    .catch(error => {
      throw error;
    });
  return fetchResult;
}

export function getPlateResult(SERVER_API, query) {
  const url = `${SERVER_API}/api/search/polri/plate${query}`;
  return fetch(url, "get")
    .then(result => result)
    .catch(error => {
      throw error;
    });
}

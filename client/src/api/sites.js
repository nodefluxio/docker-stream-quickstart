import fetch from "helpers/fetch";
import { REACT_APP_API_SITES } from "config";

export function getListCamera(query) {
  const url = `${REACT_APP_API_SITES}/streams${query ? `${query}` : ""}`;
  return fetch(url, "get").then(result => result);
}

export function getCamera(id, query, nodeNumber = 0) {
  const url = `${REACT_APP_API_SITES}/streams/${nodeNumber}/${id}${
    query ? `?${query}` : ""
  }`;
  return fetch(url, "get").then(result => result);
}
export function addSites(data) {
  const url = `${REACT_APP_API_SITES}/sites`;
  return fetch(url, "post", data).then(result => result);
}

export function updateSite(data, siteID) {
  const url = `${REACT_APP_API_SITES}/sites/${siteID}`;
  return fetch(url, "put", data).then(result => result);
}

export function getSites() {
  const url = `${REACT_APP_API_SITES}/sites`;
  return fetch(url, "get").then(result => result);
}

export function deleteSite(siteID) {
  const url = `${REACT_APP_API_SITES}/sites/${siteID}`;
  return fetch(url, "delete").then(result => result);
}

export function assignSite(data, siteID) {
  const url = `${REACT_APP_API_SITES}/sites/${siteID}/assign-stream`;
  return fetch(url, "post", data);
}

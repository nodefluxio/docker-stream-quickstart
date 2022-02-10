import fetch from "helpers/fetch";
import { REACT_APP_API_VEHICLE } from "config";

export function getList(query) {
  const url = `${REACT_APP_API_VEHICLE}/vehicles${query && `?${query}`}`;
  return fetch(url, "get").then(result => result);
}

export function getDetail(id) {
  const url = `${REACT_APP_API_VEHICLE}/vehicles/${id}`;
  return fetch(url, "get").then(result => result);
}

export function register(data) {
  const url = `${REACT_APP_API_VEHICLE}/vehicles`;
  return fetch(url, "post", data).then(result => result);
}

export function update(id, data) {
  const url = `${REACT_APP_API_VEHICLE}/vehicles/${id}`;
  return fetch(url, "put", data).then(result => result);
}

export function deleteSingle(id) {
  const url = `${REACT_APP_API_VEHICLE}/vehicles/${id}`;
  return fetch(url, "delete").then(result => result);
}

export function deleteAll() {
  const url = `${REACT_APP_API_VEHICLE}/vehicles`;
  return fetch(url, "delete").then(result => result);
}

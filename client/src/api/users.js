import fetch from "helpers/fetch";
import { REACT_APP_API_AUTH } from "config";

export function getUsers(query) {
  const url = `${REACT_APP_API_AUTH}/manage-users${query && `?${query}`}`;
  return fetch(url, "get").then(result => result);
}

export function addUser(data) {
  const url = `${REACT_APP_API_AUTH}/manage-users`;
  return fetch(url, "post", data).then(result => result);
}

export function deleteUser(id) {
  const url = `${REACT_APP_API_AUTH}/manage-users/${id}`;
  return fetch(url, "delete").then(result => result);
}

export function editUser(id, data) {
  const url = `${REACT_APP_API_AUTH}/manage-users/${id}`;
  return fetch(url, "put", data).then(result => result);
}

export function changePassword(id, data) {
  const url = `${REACT_APP_API_AUTH}/manage-users/${id}/change-password`;
  return fetch(url, "put", data).then(result => result);
}

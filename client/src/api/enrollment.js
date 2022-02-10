import fetch from "helpers/fetch";
import { REACT_APP_API_ENROLLMENT } from "config";

export function getListEnrollment(query) {
  const url = `${REACT_APP_API_ENROLLMENT}/enrollment${query && `?${query}`}`;
  return fetch(url, "get").then(result => result);
}

export function getEnrollment(id) {
  const url = `${REACT_APP_API_ENROLLMENT}/enrollment/${id}`;
  return fetch(url, "get").then(result => result);
}

export function createEnrollment(data) {
  const url = `${REACT_APP_API_ENROLLMENT}/enrollment`;
  return fetch(url, "post", data).then(result => result);
}

export function backupEnrollment() {
  const url = `${REACT_APP_API_ENROLLMENT}/files/enrollment`;
  return fetch(url, "get", null, "blob").then(result => result);
}

export function updateEnrollment(id, data) {
  const url = `${REACT_APP_API_ENROLLMENT}/enrollment/${id}`;
  return fetch(url, "put", data).then(result => result);
}

export function deleteEnrollment(id) {
  const url = `${REACT_APP_API_ENROLLMENT}/enrollment/${id}`;
  return fetch(url, "delete").then(result => result);
}

export function deleteAllEnrollment() {
  const url = `${REACT_APP_API_ENROLLMENT}/enrollment`;
  return fetch(url, "delete").then(result => result);
}

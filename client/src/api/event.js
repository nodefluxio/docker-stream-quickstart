import fetch from "helpers/fetch";
import { REACT_APP_API_EVENT } from "config";

export function getEvent(query) {
  const url = `${REACT_APP_API_EVENT}/events${query && `?${query}`}`;
  return fetch(url, "get").then(result => result);
}

export function createEventExport(query) {
  const url = `${REACT_APP_API_EVENT}/events/export${query && `?${query}`}`;
  return fetch(url, "get").then(result => result);
}

export function checkExportStatus() {
  const url = `${REACT_APP_API_EVENT}/events/export/status`;
  return fetch(url, "get").then(result => result);
}

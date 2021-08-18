import fetch from "helpers/fetch";
import { REACT_APP_API_CAMERA } from "config";

export function getListAnalytic() {
  const url = `${REACT_APP_API_CAMERA}/analytic_list`;
  return fetch(url, "get").then(result => result);
}

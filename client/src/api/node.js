import { REACT_APP_API_NODE } from "config";
import fetch from "helpers/fetch";

export function getListNode() {
  const url = `${REACT_APP_API_NODE}/node_status`;
  return fetch(url, "get").then(result => result);
}

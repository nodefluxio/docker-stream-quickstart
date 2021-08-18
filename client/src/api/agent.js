import fetch from "helpers/fetch";
import { REACT_APP_HOST } from "config";

export function getAgentStatus() {
  const url = `${REACT_APP_HOST}/v1/agents/status`;
  return fetch(url, "get")
    .then(result => result)
    .catch(error => {
      throw error;
    });
}

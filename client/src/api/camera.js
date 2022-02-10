import fetch from "helpers/fetch";
import { REACT_APP_API_CAMERA } from "config";

export function createCamera(data, nodeNumber = 0) {
  const url = `${REACT_APP_API_CAMERA}/streams/${nodeNumber}`;
  return fetch(url, "post", data).then(result => result);
}

export function createStream(
  id,
  nodeNumber = 0,
  analyticCode = "NFV4-FR",
  pipelineConfig = {}
) {
  const url = `${REACT_APP_API_CAMERA}/pipeline/${nodeNumber}/${id}/${analyticCode}`;
  return fetch(url, "post", pipelineConfig).then(result => result);
}

export function updateCamera(data, id, nodeNumber = 0) {
  const url = `${REACT_APP_API_CAMERA}/streams/${nodeNumber}/${id}`;
  return fetch(url, "put", data).then(result => result);
}

export function deleteCamera(id, nodeNumber = 0) {
  const url = `${REACT_APP_API_CAMERA}/streams/${nodeNumber}/${id}`;
  return fetch(url, "delete").then(result => result);
}

export function getCameraImage(id, nodeNumber = 0) {
  const url = `${REACT_APP_API_CAMERA}/stream_jpeg/${nodeNumber}/${id}`;
  return fetch(url, "get", null, "blob").then(result => result);
}

export function getResource() {
  const url = `${REACT_APP_API_CAMERA}/resource_stats`;
  return fetch(url, "get").then(result => result);
}

export function deleteAnalytic(streamId, analyticId, nodeNumber = 0) {
  const url = `${REACT_APP_API_CAMERA}/pipeline/${nodeNumber}/${streamId}/${analyticId}`;
  return fetch(url, "delete").then(result => result);
}

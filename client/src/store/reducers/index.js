import { combineReducers } from "redux";
import user from "./user";
import popupFeedback from "./notification";
import agent from "./agent";
import exportEvent from "./eventDownloader";
import eventContext from "./eventContext";
import cameraMenu from "./cameraMenu";

export default combineReducers({
  user,
  popupFeedback,
  agent,
  exportEvent,
  eventContext,
  cameraMenu
});

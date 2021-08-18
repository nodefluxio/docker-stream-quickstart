import {
  CREATE_EXPORT,
  EXPORT_STATUS,
  RESET_DOWNLOADER_STATE
} from "store/actionType";

const initialState = {
  createEvent: false,
  status: "",
  disableDownload: false,
  statusMsg: ""
};

export default function exportEvent(state = initialState, action) {
  switch (action.type) {
    case CREATE_EXPORT:
      return {
        ...state,
        createEvent: action.createSuccess,
        disableDownload: action.createSuccess,
        status: action.createSuccess ? "running" : "error"
      };
    case EXPORT_STATUS:
      return {
        ...state,
        status: action.data.status,
        statusMsg: action.data.message,
        disableDownload:
          action.data.status !== "error" ? false : state.disableDownload
      };
    case RESET_DOWNLOADER_STATE:
      return initialState;
    default:
      return state;
  }
}

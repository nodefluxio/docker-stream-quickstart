import {
  SHOW_GENERAL_NOTIFICATION,
  CLOSE_GENERAL_NOTIFICATION
} from "../actionType";

const popupFeedback = (state = { show: false }, action) => {
  switch (action.type) {
    case SHOW_GENERAL_NOTIFICATION:
      return action.payload.result;
    case CLOSE_GENERAL_NOTIFICATION:
      return action.payload.result;
    default:
      return state;
  }
};

export default popupFeedback;

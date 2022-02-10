import {
  SHOW_CONFIRMATION_MODAL,
  CLOSE_CONFIRMATION_MODAL,
  SHOW_MODAL_INFO_CAM,
  CLOSE_MODAL_INFO_CAM
} from "../actionType";

const cameraMenu = (
  state = { showConfirmation: false, showInfo: false },
  action
) => {
  switch (action.type) {
    case SHOW_CONFIRMATION_MODAL:
      return action.payload.result;
    case CLOSE_CONFIRMATION_MODAL:
      return action.payload.result;
    case SHOW_MODAL_INFO_CAM:
      return action.payload.result;
    case CLOSE_MODAL_INFO_CAM:
      return action.payload.result;
    default:
      return state;
  }
};

export default cameraMenu;

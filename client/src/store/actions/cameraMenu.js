import {
  SHOW_CONFIRMATION_MODAL,
  CLOSE_CONFIRMATION_MODAL,
  SHOW_MODAL_INFO_CAM,
  CLOSE_MODAL_INFO_CAM
} from "../actionType";

export const showConfirmationModal = data => ({
  type: SHOW_CONFIRMATION_MODAL,
  payload: {
    result: {
      showConfirmation: true,
      showInfo: false,
      selectedID: data.selectedID,
      deleteFunction: data.deleteFunction
    }
  }
});

export const closeConfirmationModal = () => ({
  type: CLOSE_CONFIRMATION_MODAL,
  payload: {
    result: { showConfirmation: false, showInfo: false }
  }
});

export const showInformationModal = data => ({
  type: SHOW_MODAL_INFO_CAM,
  payload: {
    result: {
      showInfo: true,
      showConfirmation: false,
      selectedID: data.selectedID,
      selectedNode: data.selectedNode
    }
  }
});

export const closeInformationModal = () => ({
  type: CLOSE_MODAL_INFO_CAM,
  payload: {
    result: { showInfo: false, showConfirmation: false }
  }
});

import {
  SHOW_GENERAL_NOTIFICATION,
  CLOSE_GENERAL_NOTIFICATION
} from "../actionType";

export const showGeneralNotification = data => ({
  type: SHOW_GENERAL_NOTIFICATION,
  payload: {
    result: { show: true, type: data.type, desc: data.desc }
  }
});

export const closeNotificationGeneral = () => ({
  type: CLOSE_GENERAL_NOTIFICATION,
  payload: {
    result: { show: false }
  }
});

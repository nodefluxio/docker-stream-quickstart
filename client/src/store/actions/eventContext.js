import {
  CLOSE_ADD_ENROLLMENT,
  SHOW_ADD_ENROLLMENT,
  SET_EVENTCONTEXT_DATA
} from "store/actionType";

export const setDataContext = data => ({
  type: SET_EVENTCONTEXT_DATA,
  data
});

export const openModalEnrollment = () => ({
  type: SHOW_ADD_ENROLLMENT
});

export const closeModalEnrollment = () => ({
  type: CLOSE_ADD_ENROLLMENT
});

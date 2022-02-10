import {
  SHOW_ADD_ENROLLMENT,
  CLOSE_ADD_ENROLLMENT,
  SET_EVENTCONTEXT_DATA
} from "store/actionType";
import { ANALYTIC_METADATA } from "constants/analyticMetadata";

const initialState = {
  openModal: false,
  data: {},
  type: ""
};

export default function eventContext(state = initialState, action) {
  switch (action.type) {
    case SET_EVENTCONTEXT_DATA:
      return {
        ...state,
        data: action.data,
        type: ANALYTIC_METADATA[action.data.analytic_id].enrollment_type || ""
      };
    case SHOW_ADD_ENROLLMENT:
      return {
        ...state,
        openModal: state.type !== ""
      };
    case CLOSE_ADD_ENROLLMENT:
      return initialState;
    default:
      return state;
  }
}

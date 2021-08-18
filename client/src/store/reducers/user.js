import {
  USER_LOGIN,
  USER_UPDATE_DATA,
  USER_LOGOUT,
  USER_FAILED_AUTH,
  USER_RESET_AUTH_ERROR,
  REFRESH_ACCESS_TOKEN
} from "store/actionType";

const initialState = {
  id: null,
  token: null,
  time: null,
  email: null,
  fullname: null,
  role: null,
  username: null,
  versionClient: null,
  versionApi: null,
  errorMessage: null,
  isRetryLogin: null
};

export default function detail(state = initialState, action) {
  switch (action.type) {
    case USER_LOGIN:
      return {
        ...state,
        token: action.token,
        errorMessage: null,
        isRetryLogin: null
      };
    case USER_UPDATE_DATA:
      return {
        ...state,
        id: action.id,
        email: action.email,
        fullname: action.fullname,
        role: action.role,
        username: action.username,
        versionClient: action.versionClient,
        versionApi: action.versionApi
      };
    case USER_LOGOUT:
      return initialState;
    case USER_FAILED_AUTH:
      return {
        ...state,
        token: null,
        errorMessage: action.message,
        isRetryLogin: true
      };
    case USER_RESET_AUTH_ERROR:
      return {
        ...state,
        errorMessage: null
      };
    case REFRESH_ACCESS_TOKEN:
      return {
        ...state,
        token: action.token,
        isRetryLogin: true
      };
    default:
      return state;
  }
}

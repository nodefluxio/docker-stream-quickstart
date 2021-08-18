import { createStore, applyMiddleware, compose } from "redux";
import thunkMiddleware from "redux-thunk";
import rootReducer from "./reducers";

const InitialState = {};
const { NODE_ENV } = process.env;
const enhancers = compose(
  typeof window !== "undefined" && NODE_ENV !== "production"
    ? // eslint-disable-next-line no-underscore-dangle
      window.window.__REDUX_DEVTOOLS_EXTENSION__ &&
        // eslint-disable-next-line no-underscore-dangle
        window.window.__REDUX_DEVTOOLS_EXTENSION__()
    : f => f
);

const createStoreWithMiddleware = applyMiddleware(thunkMiddleware)(createStore);

export function initializeStore(initialState = InitialState) {
  return createStoreWithMiddleware(rootReducer, initialState, enhancers);
}

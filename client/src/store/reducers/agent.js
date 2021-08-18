import { AGENT_STATUS } from "store/actionType";

import { USE_CES } from "config";

function renderStatus(status) {
  switch (status) {
    case "done":
      return "Connected";
    case "starting":
      return "Syncing";
    case "running":
      return "Syncing";
    case "disconnected":
      return "Disconnected";
    default:
      return "Disconnected";
  }
}

export default function agent(state = "-", action) {
  switch (action.type) {
    case AGENT_STATUS:
      return renderStatus(action.payload.status);
    default:
      return !USE_CES ? renderStatus("done") : renderStatus(state);
  }
}

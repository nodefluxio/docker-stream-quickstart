import { getAgentStatus } from "api";
import { AGENT_STATUS } from "store/actionType";

export function getStatus() {
  return async dispatch => {
    try {
      const agentStatus = await getAgentStatus();
      if (agentStatus.ok) {
        return dispatch({
          type: AGENT_STATUS,
          payload: {
            status: agentStatus.result.status
          }
        });
      }
      return dispatch({
        type: AGENT_STATUS,
        payload: {
          status: "disconnected"
        }
      });
    } catch {
      return dispatch({
        type: AGENT_STATUS,
        payload: {
          status: "disconnected"
        }
      });
    }
  };
}

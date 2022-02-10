import QueryString from "qs";
import { getTimeZone } from "helpers/dateTime";
import { checkExportStatus, createEventExport } from "api";
import {
  EXPORT_STATUS,
  CREATE_EXPORT,
  RESET_DOWNLOADER_STATE
} from "../actionType";

export function requestEventExport(query) {
  const queryFromEvent = QueryString.parse(query);
  delete queryFromEvent.page;
  delete queryFromEvent.limit;
  queryFromEvent.timezone = getTimeZone();
  const exportQuery = QueryString.stringify(queryFromEvent);
  return async dispatch => {
    try {
      const statusCreateExport = await createEventExport(exportQuery);
      if (statusCreateExport.ok) {
        return dispatch({
          type: CREATE_EXPORT,
          createSuccess: statusCreateExport.ok
        });
      }
      return dispatch({ type: CREATE_EXPORT, createSuccess: false });
    } catch {
      return dispatch({ type: CREATE_EXPORT, createSuccess: false });
    }
  };
}

export function requestExportStatus() {
  return async dispatch => {
    try {
      const exportStatus = await checkExportStatus();
      if (exportStatus.status === "downloaded") {
        return dispatch({ type: RESET_DOWNLOADER_STATE });
      }
      return dispatch({
        type: EXPORT_STATUS,
        data: { status: exportStatus.status, message: exportStatus.message }
      });
    } catch {
      return dispatch({
        type: EXPORT_STATUS,
        data: { status: "error", message: "Failed to generate export file" }
      });
    }
  };
}

export function resetDownloaderState() {
  return async dispatch => dispatch({ type: RESET_DOWNLOADER_STATE });
}

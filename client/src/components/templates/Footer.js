import React, { useEffect, useState, useRef } from "react";
import PropTypes from "prop-types";
import { connect, useDispatch } from "react-redux";
import { USE_CES, LONG_POLL_INTERVAL, REACT_APP_API_EVENT } from "config";
import styled, { keyframes } from "styled-components";
import theme from "theme";

import { getStatus } from "store/actions/agent";
import LinkButton from "components/atoms/LinkButton";
import {
  requestExportStatus,
  resetDownloaderState
} from "store/actions/eventDownloader";
import LoadingSpinner from "components/atoms/LoadingSpinner";
import IconPlaceholder from "components/atoms/IconPlaceholder";

import DownloadIcon from "assets/icon/visionaire/download.svg";

function Footer({ agent, exportEvent, history }) {
  const [agentStatus, setAgentStatus] = useState(agent);
  const [eventButtonText, setEventButtonText] = useState("");
  const [showButton, setShowButton] = useState(false);
  const [showExportComponent, setShowExportComponent] = useState(false);
  const dispatch = useDispatch();
  const downloadLink = useRef(null);

  function renderColor(status) {
    switch (status) {
      case "Disconnected":
        return theme.inlineError;
      case "Connected":
        return theme.success;
      case "Syncing":
        return theme.yellow;
      default:
        return theme.white;
    }
  }

  const menu = [
    {
      label: "License",
      url: "/license"
    }
  ];

  useEffect(() => {
    if (USE_CES) {
      const interval = setInterval(() => {
        dispatch(getStatus());
      }, LONG_POLL_INTERVAL);
      return () => clearInterval(interval);
    }
    return null;
  }, []);

  function handleEventDownloadClick() {
    if (exportEvent.status !== "error") {
      downloadLink.current.click();
    }
    dispatch(resetDownloaderState());
    setShowExportComponent(false);
  }

  useEffect(() => {
    setAgentStatus(agent);
  }, [agent]);

  useEffect(() => {
    dispatch(requestExportStatus());
  }, []);

  useEffect(() => {
    let interval = null;
    switch (exportEvent.status) {
      case "running": {
        setShowExportComponent(true);
        setShowButton(true);
        setEventButtonText("Prepare Exported Files...");
        interval = setInterval(() => {
          dispatch(requestExportStatus());
        }, LONG_POLL_INTERVAL);
        break;
      }
      case "ready": {
        setShowExportComponent(true);
        setShowButton(true);
        setEventButtonText("Download Exported File");
        if (interval !== null) {
          clearInterval(interval);
        }
        break;
      }
      case "error": {
        setEventButtonText("Failed to Generate Export File");
        if (interval !== null) {
          clearInterval(interval);
        }
        break;
      }
      case "downloaded": {
        dispatch(resetDownloaderState());
        setShowExportComponent(false);
        if (interval !== null) {
          clearInterval(interval);
        }
        break;
      }
      default:
        return null;
    }
    return () => {
      if (interval !== null) {
        clearInterval(interval);
      }
    };
  }, [exportEvent]);

  return (
    <Wrapper>
      <Right>
        <ButtonWrapper>
          {menu.map(item => (
            <LinkButton
              key={item.label}
              id={item.label}
              onClick={() => history.push(item.url)}
            >
              <span>{item.label}</span>
            </LinkButton>
          ))}
        </ButtonWrapper>
      </Right>
      <Left>
        {showExportComponent && (
          <ButtonWrapper show={showButton}>
            <IconPlaceholder className="slide-in-top icon">
              <img src={DownloadIcon} alt="download-icon" />
            </IconPlaceholder>
            <LinkButton
              id="download-event-button"
              bordered
              className={`button ${
                exportEvent.status === "ready" ? "active" : ""
              }`}
              onClick={() => handleEventDownloadClick()}
              disabled={exportEvent.status === "running"}
              style={{ paddingRight: "10px" }}
            >
              <LoadingSpinner show={exportEvent.status === "running"} />
              {exportEvent.status === "ready" && (
                <img src={DownloadIcon} alt="download-icon" />
              )}
              <span
                style={{
                  color:
                    exportEvent.status === "error"
                      ? theme.inlineError
                      : theme.white
                }}
              >
                {eventButtonText}
              </span>
            </LinkButton>
            <a
              style={{ display: "hidden" }}
              ref={downloadLink}
              href={`${REACT_APP_API_EVENT}/events/export/download`}
            />
          </ButtonWrapper>
        )}

        {USE_CES && (
          <CoordinatorStatus spanColor={renderColor(agentStatus)}>
            Coordinator Status: <span>{agentStatus}</span>
          </CoordinatorStatus>
        )}
      </Left>
    </Wrapper>
  );
}

Footer.propTypes = {
  agent: PropTypes.string.isRequired,
  exportEvent: PropTypes.object.isRequired,
  history: PropTypes.object.isRequired
};

function mapStateToProps(state) {
  return {
    agent: state.agent,
    exportEvent: state.exportEvent
  };
}

export default connect(mapStateToProps)(Footer);

const Right = styled.div`
  height: 100%;
  margin-left: 5px;
`;
const Left = styled.div``;

const slideIn = keyframes`
  0% {
    -webkit-transform: translateY(-1000px);
            transform: translateY(-1000px);
    opacity: 0;
  }
  100% {
    -webkit-transform: translateY(0);
            transform: translateY(0);
    opacity: 1;
  }
`;

const Wrapper = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  min-height: 48px;
  width: 100%;
  border-top: 1px solid #372463;
  .slide-in-top {
    -webkit-animation: ${slideIn} 0.3s cubic-bezier(0.25, 0.46, 0.45, 0.94) both;
    animation: ${slideIn} 0.3s cubic-bezier(0.25, 0.46, 0.45, 0.94) both;
  }
`;

const ButtonWrapper = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  height: 100%;
  .icon {
    position: absolute;
    bottom: -50px;
  }
  .button {
    visibility: ${props => (props.show ? `visible` : `hidden`)};
    transition: visibility 0.35s;
    transition-delay: 0.35s;
    z-index: 1;
    padding-right: 15px;
  }
`;

const CoordinatorStatus = styled.div`
  padding: 0 20px;
  span {
    color: ${props => props.spanColor};
  }
`;

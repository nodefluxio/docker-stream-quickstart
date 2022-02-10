import React, { Fragment, useState, useEffect } from "react";
import Styled from "styled-components";
import PropTypes from "prop-types";
import VisualizationWrapper from "components/molecules/VisualizationWrapper";
import ImageWrapper from "components/molecules/ImageWrapper";
import { Menu } from "components/atoms/Dropdown";
import { useDispatch } from "react-redux";

import { getCookie } from "helpers/cookies";
import parseJwt from "helpers/parseJWT";
import { REACT_APP_API_CAMERA } from "config";

import { getCameraImage } from "api";
import ArrowRight from "assets/icon/visionaire/arrow-right.svg";
import {
  showConfirmationModal,
  showInformationModal
} from "store/actions/cameraMenu";

export default function VanillaVisualizationWrapper(props) {
  const {
    streamID,
    streamName,
    isImage,
    width,
    height,
    history,
    deleteData,
    refreshTime,
    selectedAnalytic,
    fps,
    selectedNode
  } = props;

  const [role, setRole] = useState("operator");
  const dispatch = useDispatch();

  const cameraMenu = () => (
    <Fragment>
      <MenuWrapper>
        {history.location.pathname.split("/").length <= 2 && (
          <Menu
            className="menu"
            onClick={() =>
              props.history.push(`/camera/${selectedNode}/${streamID}`)
            }
          >
            View This Camera
            <img src={ArrowRight} alt="arrow-right-icon" />
          </Menu>
        )}
        <Menu
          className="menu"
          onClick={() =>
            dispatch(
              showInformationModal({
                selectedID: streamID,
                selectedNode
              })
            )
          }
        >
          View This Camera Information
          <img src={ArrowRight} alt="arrow-right-icon" />
        </Menu>
        {role === "superadmin" && (
          <Fragment>
            <Menu
              className="menu"
              onClick={() =>
                props.history.push(
                  `/assignment?step=1&selected_camera=${streamID}&selected_node=${selectedNode}`
                )
              }
            >
              Edit This Camera
              <img src={ArrowRight} alt="arrow-right-icon" />
            </Menu>
            <Menu
              className="menu"
              onClick={() =>
                props.history.push(
                  `/assignment?step=2&selected_camera=${streamID}&selected_node=${selectedNode}`
                )
              }
            >
              Configure Analytics
              <img src={ArrowRight} alt="arrow-right-icon" />
            </Menu>
            <Menu
              className="menu"
              onClick={() => {
                dispatch(
                  showConfirmationModal({
                    selectedID: streamID,
                    deleteFunction: deleteData
                  })
                );
              }}
            >
              Remove This Camera
              <img src={ArrowRight} alt="arrow-right-icon" />
            </Menu>
          </Fragment>
        )}
      </MenuWrapper>
    </Fragment>
  );

  useEffect(() => {
    const accessToken = getCookie("access_token");
    const userCookies = parseJwt(accessToken);
    setRole(userCookies.role);
  }, []);

  function generateURLCameraRequest() {
    let url = `${REACT_APP_API_CAMERA}/mjpeg`;
    if (fps > 0) {
      url += `/fps/${fps}`;
    }
    url += `/${selectedNode}/${streamID}`;
    if (selectedAnalytic) {
      url += `/${selectedAnalytic}`;
    }
    return url;
  }

  return isImage ? (
    <ImageWrapper
      id={streamID}
      fetchUrl={() => getCameraImage(streamID, selectedNode)}
      width={width}
      height={height}
      name={streamName}
      refreshInterval={refreshTime}
      onClick={() =>
        history.location.pathname.split("/").length <= 2
          ? history.push(`/camera/${selectedNode}/${streamID}`)
          : {}
      }
      menu={cameraMenu(streamID)}
    />
  ) : (
    <VisualizationWrapper
      url={generateURLCameraRequest()}
      width={width}
      height={height}
      name={streamName}
      onClick={() =>
        history.location.pathname.split("/").length <= 2
          ? history.push(`/camera/${selectedNode}/${streamID}`)
          : {}
      }
      menu={cameraMenu(streamID)}
    />
  );
}

VanillaVisualizationWrapper.propTypes = {
  streamID: PropTypes.string.isRequired,
  streamName: PropTypes.string.isRequired,
  history: PropTypes.object.isRequired,
  width: PropTypes.string,
  height: PropTypes.string,
  deleteData: PropTypes.func,
  refreshTime: PropTypes.number,
  isImage: PropTypes.bool,
  selectedAnalytic: PropTypes.string,
  fps: PropTypes.number,
  selectedNode: PropTypes.number
};

VanillaVisualizationWrapper.defaultProps = {
  isImage: false,
  width: "100%",
  height: "100%",
  deleteData: () => {},
  refreshTime: 3000,
  selectedAnalytic: "",
  fps: 0,
  selectedNode: 0
};

const MenuWrapper = Styled.div`
  text-transform: ${props => props.textTransform || `none`} ;
  font-weight: bold;
  font-size: 12px;
  line-height: 14px;
  color: #E5E5E5;

  .menu {
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    padding-right: 17px;
  }
`;

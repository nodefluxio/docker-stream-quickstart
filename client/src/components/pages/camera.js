/* eslint-disable no-console */
/* eslint-disable no-nested-ternary */
import React, { useState, useEffect, useContext, Fragment } from "react";
import PropTypes from "prop-types";
import Styled, { ThemeContext } from "styled-components";
import { REACT_APP_API_CAMERA, REACT_APP_REFRESH_TIME } from "config";
import Dropdown, { Menu } from "components/atoms/Dropdown";
import BlankContainer from "components/atoms/BlankContainer";
import Button from "components/molecules/Button";
import Text from "components/atoms/Text";
import LoadingSpinner from "components/atoms/LoadingSpinner";
import Modal from "components/molecules/Modal";
import VisualisationWrapper from "components/molecules/VisualisationWrapper";
import ImageWrapper from "components/molecules/ImageWrapper";
import PageWrapper from "components/organisms/pageWrapper";
import { getListCamera, deleteCamera, getCameraImage } from "api";
import { AddCamera } from "components/organisms/modals";

import ArrowRight from "assets/icon/visionaire/arrow-right.svg";
import GlobeBg from "assets/images/globe.svg";
import CheckIcon from "assets/icon/visionaire/check-small.svg";
import AnalyticIcon from "assets/icon/visionaire/analytic.svg";
import GridIcon from "assets/icon/visionaire/grid.svg";
import RefreshIcon from "assets/icon/visionaire/refresh.svg";
import { showVisualisation, stopVisualisation } from "assets/js/visualstreamer";

function Camera(props) {
  const { history } = props;
  const [data, setData] = useState([]);
  const [heightSize, setHeightSize] = useState(0);

  // eslint-disable-next-line no-unused-vars
  const [gridColumn, setGridColumn] = useState(4);
  const [blank, setBlank] = useState(false);
  const [isConfirmOpen, setIsConfirmOpen] = useState(false);
  const [selectedId, setSelectedId] = useState(null);
  const [showAdd, setShowAdd] = useState(false);
  const themeContext = useContext(ThemeContext);
  const [errorFetch, setErrorFetch] = useState("");
  const [refreshTime, setRefreshTime] = useState(3000);
  const [isLoading, setIsLoading] = useState(false);

  const gridMenu = (
    <MenuWrapper textTransform="uppercase">
      <Menu className="menu" onClick={() => setGridColumn(3)}>
        3 Grid View
        {gridColumn === 3 && <img src={CheckIcon} alt="check-icon" />}
      </Menu>
      <Menu className="menu" onClick={() => setGridColumn(4)}>
        4 Grid View
        {gridColumn === 4 && <img src={CheckIcon} alt="check-icon" />}
      </Menu>
      <Menu className="menu" onClick={() => setGridColumn(6)}>
        6 Grid View
        {gridColumn === 6 && <img src={CheckIcon} alt="check-icon" />}
      </Menu>
    </MenuWrapper>
  );

  const timeInterval = [
    {
      label: "1s",
      value: 1000
    },
    {
      label: "3s",
      value: 3000
    },
    {
      label: "5s",
      value: 5000
    },
    {
      label: "25s",
      value: 25000
    },
    {
      label: "1m",
      value: 60000
    },
    {
      label: "5m",
      value: 300000
    }
  ];

  const refreshMenu = (
    <MenuWrapper>
      {timeInterval.map(time => (
        <Menu
          key={time.value}
          className="menu"
          onClick={() => setRefreshTime(time.value)}
        >
          {time.label}
          {refreshTime === time.value && (
            <img src={CheckIcon} alt="check-icon" />
          )}
        </Menu>
      ))}
    </MenuWrapper>
  );

  function callData() {
    setIsLoading(true);
    getListCamera()
      .then(result => {
        const resultdata = result.streams;
        setData(resultdata);
        if (resultdata.length === 0) {
          setBlank(true);
        }
        const calHeightSize = (97.6 / gridColumn / 4) * 3;
        setHeightSize(calHeightSize);
        setIsLoading(false);
      })
      .catch(error => {
        setErrorFetch(error.message);
        setBlank(true);
        setIsLoading(false);
      });
  }

  function deleteData(id) {
    deleteCamera(id).then(result => {
      if (result.code === 200) {
        const index = data.findIndex(x => x.stream_id === result.stream_id);
        const newData = data;
        newData.splice(index, 1);
        setData(newData);
        showVisualisation(true);
      }
    });
  }

  useEffect(() => {
    callData();
    return () => stopVisualisation();
  }, []);

  useEffect(() => {
    const interval = setInterval(() => {
      window.location.reload();
    }, REACT_APP_REFRESH_TIME);
    return () => clearInterval(interval);
  }, [REACT_APP_REFRESH_TIME]);

  useEffect(() => {
    if (data.length > 0 && data.length < 5) {
      showVisualisation();
    }
    return () => stopVisualisation();
  }, [data]);

  useEffect(() => {
    const calHeightSize = (100 / gridColumn / 4) * 3;
    setHeightSize(calHeightSize);
  }, [gridColumn]);

  const closeModal = flag => {
    if (flag) {
      setData([]);
      callData();
    }
    setShowAdd(false);
  };

  return (
    <PageWrapper
      title="All Camera"
      history={history}
      buttonGroup={
        <Fragment>
          {data.length >= 5 && (
            <Dropdown
              id="refresh_interval"
              overlay={refreshMenu}
              width="216px"
              className="button-control"
            >
              <IconPlaceholderButton>
                <img alt="refresh-icon" src={RefreshIcon} />
              </IconPlaceholderButton>
            </Dropdown>
          )}
          <Dropdown
            id="grid_view_icon"
            overlay={gridMenu}
            width="216px"
            className="button-control"
          >
            <IconPlaceholderButton>
              <img alt="grid-icon" src={GridIcon} />
            </IconPlaceholderButton>
          </Dropdown>
          <Button
            id="assign-analytics"
            onClick={() => props.history.push("/assignment")}
            type="secondary"
            className="button-control icon-text"
          >
            <img alt="assign-analytics-icon" src={AnalyticIcon} />
            ASSIGN ANALYTICS
          </Button>
        </Fragment>
      }
    >
      {blank ? (
        <BlankContainer bg={GlobeBg}>
          <BoxedDiv>
            {errorFetch !== "" ? (
              <Text
                size="18"
                color={themeContext.inlineError}
                textTransform="uppercase"
                weight="600"
              >
                {errorFetch}
              </Text>
            ) : (
              <Fragment>
                <Text size="18" color="#fff" weight="600">
                  NO CAMERA FOUND
                </Text>

                <Button onClick={() => setShowAdd(true)} width="336px">
                  CONFIGURE LIVE FEED
                </Button>
              </Fragment>
            )}
          </BoxedDiv>
        </BlankContainer>
      ) : isLoading ? (
        <BlankContainer>
          <LoadingSpinner show={isLoading} />
        </BlankContainer>
      ) : (
        <LiveFeedContainer>
          <ImageContainer>
            {data.length >= 5
              ? data.map(value => (
                  <ImageWrapper
                    key={value.stream_id}
                    id={value.stream_id}
                    fetchUrl={getCameraImage}
                    width={`${100 / gridColumn}%`}
                    height={`${heightSize}vw`}
                    name={value.stream_name}
                    refreshInterval={refreshTime}
                    onClick={() =>
                      props.history.push(`/camera/${value.stream_id}`)
                    }
                    menu={
                      <Fragment>
                        <MenuWrapper>
                          <Menu
                            className="menu"
                            onClick={() =>
                              props.history.push(`/camera/${value.stream_id}`)
                            }
                          >
                            View This Camera
                            <img src={ArrowRight} alt="arrow-right-icon" />
                          </Menu>
                          <Menu
                            className="menu"
                            onClick={() => {
                              setIsConfirmOpen(true);
                              setSelectedId(value.stream_id);
                            }}
                          >
                            Remove This Camera
                            <img src={ArrowRight} alt="arrow-right-icon" />
                          </Menu>
                        </MenuWrapper>
                      </Fragment>
                    }
                  />
                ))
              : data.map(value => (
                  <VisualisationWrapper
                    key={value.stream_id}
                    url={`${REACT_APP_API_CAMERA}/mjpeg/fps/4/0/${value.stream_id}`}
                    width={`${100 / gridColumn}%`}
                    height={`${heightSize}vw`}
                    name={value.stream_name}
                    onClick={() =>
                      props.history.push(`/camera/${value.stream_id}`)
                    }
                    menu={
                      <Fragment>
                        <MenuWrapper>
                          <Menu
                            className="menu"
                            onClick={() =>
                              props.history.push(`/camera/${value.stream_id}`)
                            }
                          >
                            View This Camera
                            <img src={ArrowRight} alt="arrow-right-icon" />
                          </Menu>
                          <Menu
                            className="menu"
                            onClick={() => {
                              setIsConfirmOpen(true);
                              setSelectedId(value.stream_id);
                            }}
                          >
                            Remove This Camera
                            <img src={ArrowRight} alt="arrow-right-icon" />
                          </Menu>
                        </MenuWrapper>
                      </Fragment>
                    }
                  />
                ))}
          </ImageContainer>
        </LiveFeedContainer>
      )}
      <Modal
        type="confirmation"
        isShown={isConfirmOpen}
        onConfirm={() => {
          deleteData(selectedId);
          setIsConfirmOpen(false);
        }}
        onClose={() => setIsConfirmOpen(false)}
        title="Delete Camera"
        buttonTitle="Proceed Delete"
        header="ARE YOU SURE TO PROCEED
        DELETING ALL THE DATA?"
        headerColor={themeContext.color8}
      />
      <AddCamera openModal={showAdd} onClose={flag => closeModal(flag)} />
    </PageWrapper>
  );
}

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

const IconPlaceholderButton = Styled.div`
  width: 40px;
  height: 40px;
  border: 1px solid #372463;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 10px;
  &:hover {
    background: #4995E9;
  }
`;

const ImageContainer = Styled.div`
  display: flex;
  flex-direction: row;
  flex-wrap: wrap;
  height: ${props => props.height};
  background: ${props => props.theme.color19};
  overflow-y: auto;
`;

const BoxedDiv = Styled.div`
  border: 1px solid #372463;
  box-sizing: border-box;
  border-radius: 8px;
  width: 400px;
  height: 152px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: space-evenly;
  background: #21153C;
`;

const LiveFeedContainer = Styled.div`
  height: calc(100vh - 112px);
  width: 100vw;
  overflow-y: auto;
`;

Camera.propTypes = {
  history: PropTypes.object.isRequired,
  location: PropTypes.object,
  value: PropTypes.array
};

export default Camera;

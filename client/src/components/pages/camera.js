/* eslint-disable no-console */
/* eslint-disable no-nested-ternary */
import React, { useState, useEffect, useContext, Fragment } from "react";
import PropTypes from "prop-types";
import Styled, { ThemeContext } from "styled-components";
import { REACT_APP_REFRESH_TIME } from "config";
import Dropdown, { Menu } from "components/atoms/Dropdown";
import BlankContainer from "components/atoms/BlankContainer";
import Button from "components/molecules/Button";
import Text from "components/atoms/Text";
import LoadingSpinner from "components/atoms/LoadingSpinner";
import Modal from "components/molecules/Modal";
import Input from "components/molecules/Input";
import PageWrapper from "components/organisms/pageWrapper";
import { getListCamera, deleteCamera, getSites } from "api";
import { getCookie } from "helpers/cookies";
import parseJwt from "helpers/parseJWT";
import qs from "qs";

import GlobeBg from "assets/images/globe.svg";
import CheckIcon from "assets/icon/visionaire/check-small.svg";
import AnalyticIcon from "assets/icon/visionaire/analytic.svg";
import GridIcon from "assets/icon/visionaire/grid.svg";
import RefreshIcon from "assets/icon/visionaire/refresh.svg";
import AreaIcon from "assets/icon/visionaire/all_area.svg";
import VanillaVisualizationWrapper from "components/organisms/VanillaVisualizationWrapper";
import withModalVisualization from "components/templates/withVisualizationModals";

function Camera(props) {
  const { history } = props;
  const [data, setData] = useState([]);
  const [heightSize, setHeightSize] = useState(0);

  // eslint-disable-next-line no-unused-vars
  const [gridColumn, setGridColumn] = useState(4);
  const [blank, setBlank] = useState(false);
  const themeContext = useContext(ThemeContext);
  const [errorFetch, setErrorFetch] = useState("");
  const [refreshTime, setRefreshTime] = useState(3000);
  const [isLoading, setIsLoading] = useState(false);
  const [role, setRole] = useState("operator");
  const [sites, setSites] = useState([]);
  const [openSiteFilter, setOpenSiteFilter] = useState(false);
  const [selectedSites, setSelectedSites] = useState([]);
  const [filterMode, setFilterMode] = useState(false);

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

  function callData(query) {
    setIsLoading(true);
    getListCamera(query)
      .then(result => {
        if (result.ok) {
          const resultdata = result.data.streams;
          setData([...resultdata]);
          if (resultdata.length === 0) {
            setBlank(true);
          } else {
            setBlank(false);
          }
          const calHeightSize = (97.6 / gridColumn / 4) * 3;
          setHeightSize(calHeightSize);
          setIsLoading(false);
        }
      })
      .catch(error => {
        setErrorFetch(error.message);
        setBlank(true);
        setIsLoading(false);
      });
  }

  function filterBySite() {
    const selectedSiteID = [];
    for (let i = 0; i < selectedSites.length; i += 1) {
      selectedSiteID.push(selectedSites[i].value);
    }
    const query = qs.stringify(
      {
        "filter[site_id]": selectedSiteID
      },
      { arrayFormat: "comma" }
    );
    const editedquery = query.replace(/\+/g, "%2B");
    history.replace(`/camera?${editedquery}`);
    setOpenSiteFilter(false);
  }

  function deleteData(id, node) {
    deleteCamera(id, node).then(result => {
      if (result.code === 200) {
        const index = data.findIndex(x => x.stream_id === result.stream_id);
        const newData = data;
        newData.splice(index, 1);
        setData([...newData]);
        window.showVisualisation(true);
      }
    });
  }

  useEffect(() => {
    const accessToken = getCookie("access_token");
    const userCookies = parseJwt(accessToken);
    if (userCookies) {
      setRole(userCookies.role);
    }
    getSites().then(result => {
      if (result.ok) {
        const newSite = [];
        result.sites.map(site =>
          newSite.push({ value: site.id, label: site.name })
        );
        setSites(newSite);
      }
    });
  }, []);

  useEffect(() => {
    const interval = setInterval(() => {
      window.location.reload();
    }, REACT_APP_REFRESH_TIME);
    return () => clearInterval(interval);
  }, [REACT_APP_REFRESH_TIME]);

  useEffect(() => {
    if (data.length > 0 && data.length < 5) {
      window.showVisualisation();
    }
  }, [data]);

  useEffect(() => {
    const calHeightSize = (100 / gridColumn / 4) * 3;
    setHeightSize(calHeightSize);
  }, [gridColumn]);

  function callCameraWithFilter(query, selectedSite, siteList) {
    const array = selectedSite.split(",");
    const newValue = [];
    for (let i = 0; i < array.length; i += 1) {
      const index = siteList.findIndex(
        x => Number(x.value) === Number(array[i])
      );
      newValue.push(siteList[index]);
    }
    setSelectedSites(newValue);

    callData(query);
    setFilterMode(true);
  }

  useEffect(() => {
    const { search } = history.location;
    if (search !== "") {
      search.replace("?", "");
      const query = new URLSearchParams(search);
      const selectedSite = query.get("filter[site_id]");
      if (sites.length > 0) {
        if (selectedSite) {
          callCameraWithFilter(search, selectedSite, sites);
        }
      } else if (sites.length === 0) {
        let sitesTemp = [];
        getSites().then(result => {
          if (result.ok) {
            const newSite = [];
            result.sites.map(site =>
              newSite.push({ value: site.id, label: site.name })
            );
            sitesTemp = newSite;
            setSites(newSite);
            if (selectedSite) {
              callCameraWithFilter(search, selectedSite, sitesTemp);
            }
          }
        });
      } else {
        setFilterMode(false);
      }
    } else {
      callData();
      setFilterMode(false);
    }
  }, [history.location]);

  return (
    <PageWrapper
      title={
        filterMode ? `Camera in Selected Area(s)` : `All Camera in All Area`
      }
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
          <IconPlaceholderButton
            onClick={() => setOpenSiteFilter(true)}
            active={filterMode}
          >
            <img alt="filter area" src={AreaIcon} />
          </IconPlaceholderButton>
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
          {role === "superadmin" && (
            <Button
              id="assign-analytics"
              onClick={() => props.history.push("/assignment")}
              type="secondary"
              className="button-control icon-text"
            >
              <img alt="assign-analytics-icon" src={AnalyticIcon} />
              ASSIGN ANALYTICS
            </Button>
          )}
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

                <Button
                  onClick={() => history.push("/assignment")}
                  width="336px"
                >
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
            {data.map(value => (
              <VanillaVisualizationWrapper
                key={value.stream_id}
                streamID={value.stream_id}
                width={`${100 / gridColumn}%`}
                height={`${heightSize}vw`}
                streamName={value.stream_name}
                refreshInterval={refreshTime}
                isImage={data.length >= 5}
                deleteData={() =>
                  deleteData(value.stream_id, value.stream_node_num)
                }
                history={history}
                fps={4}
                selectedNode={value.stream_node_num}
              />
            ))}
          </ImageContainer>
        </LiveFeedContainer>
      )}
      <Modal
        show={openSiteFilter}
        className="modal-filter-site"
        title="Filter Camera by Site"
        close={() => setOpenSiteFilter(false)}
        padding="20px"
      >
        <Input
          type="select"
          option={sites}
          placeholder="All"
          onChange={e => setSelectedSites(e)}
          value={selectedSites}
          label="Sites:"
          isMulti={true}
          style={{ marginBottom: "40px" }}
        />
        <Button
          width="100%"
          style={{ marginTop: "20px" }}
          onClick={() => filterBySite()}
        >
          Save
        </Button>
      </Modal>
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
  cursor: pointer;
  ${props => props.active && `background: #4995E9;`}
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
  height: calc(100vh - 167px);
  width: 100vw;
  overflow-y: auto;
`;

Camera.propTypes = {
  history: PropTypes.object.isRequired,
  location: PropTypes.object,
  value: PropTypes.array
};

export default withModalVisualization(Camera);

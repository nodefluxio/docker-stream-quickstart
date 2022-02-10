/* eslint-disable no-console */
/* eslint-disable func-names */
import React, {
  useState,
  Fragment,
  useEffect,
  useContext,
  useRef
} from "react";
import Styled, { ThemeContext } from "styled-components";
import PropTypes from "prop-types";
import { REACT_APP_API_WEB_SOCKET, REACT_APP_REFRESH_TIME } from "config";
import { dateToIso, isoToDate, getTimeZone } from "helpers/dateTime";
import qs from "qs";

import InsightBox from "components/atoms/InsightBox";
import BlankContainer from "components/atoms/BlankContainer";
import Input from "components/molecules/Input";
import withEnrollment from "components/templates/withEnrollmentModal";
import Row from "components/atoms/Row";
import AreaSection from "components/molecules/AreaSection";
import { FilterEvent } from "components/organisms/modals";
import { getCamera, deleteCamera } from "api";

import SearchIcon from "assets/icon/visionaire/search.svg";
import Button from "components/molecules/Button";
import RefreshIcon from "assets/icon/visionaire/refresh.svg";

import PageWrapper from "components/organisms/pageWrapper";
import VanillaVisualizationWrapper from "components/organisms/VanillaVisualizationWrapper";
import EventWrapper from "components/organisms/EventWrapper";
import withModalVisualization from "components/templates/withVisualizationModals";
import { getInsightConstant } from "constants/insightConstant";
import { ANALYTIC_METADATA } from "../../constants/analyticMetadata";

function CameraDetail(props) {
  const { history, location } = props;
  const getParamsUrl = location.search.replace("?", "");
  const defaultValue = qs.parse(getParamsUrl);
  const [data, setData] = useState({});
  const themeContext = useContext(ThemeContext);
  const [newEvent, setNewEvent] = useState({});
  const [analyticList, setAnalyticList] = useState([]);
  const [analyticName, setAnalyticName] = useState("");
  const [selectedAnalytic, setSelectedAnalytic] = useState(() => {
    if (defaultValue.filter) {
      if (defaultValue.filter.analytic_id) {
        setAnalyticName(
          `${
            ANALYTIC_METADATA[defaultValue.filter.analytic_id]
              ? ` - ${
                  ANALYTIC_METADATA[defaultValue.filter.analytic_id]
                    .analytic_name
                }`
              : ""
          }`
        );
        return defaultValue.filter.analytic_id;
      }
      return "";
    }
    return "";
  });
  const [prevAnalytic, setPrevAnalytic] = useState(() => {
    if (defaultValue.filter) {
      if (defaultValue.filter.analytic_id) {
        setAnalyticName(
          `${
            ANALYTIC_METADATA[defaultValue.filter.analytic_id]
              ? ` - ${
                  ANALYTIC_METADATA[defaultValue.filter.analytic_id]
                    .analytic_name
                }`
              : ""
          }`
        );
        return defaultValue.filter.analytic_id;
      }
      return "";
    }
    return "";
  });
  const [openFilterMenu, setOpenFilterMenu] = useState(false);
  const [showResetButton, setShowResetButton] = useState(false);
  const socketRef = useRef(null);
  const [prevQuery, setPrevQuery] = useState("");
  const [prevID, setPrevID] = useState("");
  const [hideUnrecog, setHideUnrecog] = useState(true);

  const now = new Date();
  const dateFrom = dateToIso(
    new Date(now.getFullYear(), now.getMonth(), now.getDate(), 0, 0, 0)
  );
  const dateTo = dateToIso(now);

  const LIMIT_DATA = 20;

  function defaultQuery(reset = true) {
    const paths = location.pathname.split("/");
    const id = paths[paths.length - 1];
    const queryObject = {
      "filter[timestamp_from]": dateFrom,
      "filter[timestamp_to]": dateTo,
      "filter[type]":
        hideUnrecog && selectedAnalytic.includes("FR") ? "recognized" : "",
      "filter[stream_id]": id,
      limit: LIMIT_DATA,
      timezone: getTimeZone()
    };
    if (selectedAnalytic !== "") {
      queryObject["filter[analytic_id]"] = selectedAnalytic;
    }
    if (Object.keys(defaultValue).length > 0) {
      if (!reset) {
        queryObject["filter[timestamp_from]"] =
          defaultValue.filter.timestamp_from;
        queryObject["filter[timestamp_to]"] = defaultValue.filter.timestamp_to;
        if (defaultValue.search) {
          queryObject.search = defaultValue.search;
        }
        return qs.stringify(queryObject, { arrayFormat: "comma" });
      }
    }
    return qs.stringify(queryObject, { arrayFormat: "comma" });
  }

  const onConnect = (id, node) => {
    socketRef.current = new WebSocket(
      `${REACT_APP_API_WEB_SOCKET}/event_channel?stream_id=${id}&node_num=${node}`
    );
    socketRef.current.onopen = function() {
      console.log("Connected");
    };
    socketRef.current.onerror = function() {
      setTimeout(() => {
        onConnect(id, node);
      }, 3000);
    };
    socketRef.current.onmessage = async function(event) {
      const msg = JSON.parse(event.data);
      await setNewEvent(msg);
    };
    socketRef.current.onclose = function(event) {
      console.log("closing socket.... code: ", event.code);
      socketRef.current = null;
      if (event.code !== 1006) {
        setTimeout(() => {
          onConnect(id, node);
        }, 3000);
      }
    };
  };

  useEffect(() => {
    if (socketRef.current) {
      socketRef.current.close();
    }
    const paths = location.pathname.split("/");
    const id = paths[paths.length - 1];
    const node = paths[paths.length - 2];
    if (prevID !== id) {
      onConnect(id, node);
      getCamera(id, null, node)
        .then(result => {
          if (result.ok) {
            const streamData = result.stream;
            setData(streamData);
            if (streamData.pipelines !== undefined) {
              const listPipeline = [];
              const seats = [];
              if (streamData.seats && streamData.seats.length > 0) {
                streamData.seats.forEach(value => {
                  const analyticId = value.analytic_id;
                  seats[analyticId] = value.serial_number;
                  listPipeline.push({
                    value: analyticId,
                    label: ANALYTIC_METADATA[analyticId]
                      ? ANALYTIC_METADATA[analyticId].analytic_name
                      : ""
                  });
                });
              } else {
                streamData.pipelines.map(value =>
                  listPipeline.push({
                    value,
                    label: ANALYTIC_METADATA[value]
                      ? ANALYTIC_METADATA[value].analytic_name
                      : ""
                  })
                );
              }
              setAnalyticList(listPipeline);
              setSelectedAnalytic(listPipeline[0].value);
            }
            window.showVisualisation();
          }
        })
        .catch(() => {});
    }
    setPrevID(id);
    return () => {
      if (socketRef.current) {
        socketRef.current.close();
      }
      window.stopVisualisation();
    };
  }, [location.pathname]);

  useEffect(() => {
    const interval = setInterval(() => {
      window.location.reload();
    }, REACT_APP_REFRESH_TIME);
    return () => clearInterval(interval);
  }, [REACT_APP_REFRESH_TIME]);

  useEffect(() => {
    if (selectedAnalytic !== "" && selectedAnalytic !== prevAnalytic) {
      window.showVisualisation();
      setPrevAnalytic(selectedAnalytic);
      const query = defaultQuery(false);
      const editedquery = query.replace(/\+/g, "%2B");
      history.replace(`${location.pathname}?${editedquery}`);
      setAnalyticName(
        `${
          ANALYTIC_METADATA[selectedAnalytic]
            ? ` - ${ANALYTIC_METADATA[selectedAnalytic].analytic_name}`
            : ""
        }`
      );
    }
  }, [selectedAnalytic]);

  useEffect(() => {
    const query = defaultQuery(false);
    const editedquery = query.replace(/\+/g, "%2B");
    history.replace(`${location.pathname}?${editedquery}`);
  }, [hideUnrecog]);

  useEffect(() => {
    if (location.search !== "" && prevQuery !== location.search) {
      setPrevQuery(location.search);
      const defaultSearch = qs.parse(defaultQuery());
      const end = isoToDate(defaultValue.filter.timestamp_to);
      let flag = false;
      if (
        defaultValue.search !== defaultSearch.search ||
        defaultValue.filter.type !== defaultSearch.filter.type ||
        defaultValue.filter.timestamp_from !== dateFrom ||
        end.toDateString() !== now.toDateString()
      ) {
        flag = true;
      }
      setShowResetButton(flag);
    }
  }, [location.search]);

  useEffect(() => {
    if (socketRef.current.readyState > 1) {
      const paths = location.pathname.split("/");
      const id = paths[paths.length - 1];
      const node = paths[paths.length - 2];
      if (showResetButton) {
        socketRef.current.close();
      } else {
        onConnect(id, node);
      }
    }
  }, [showResetButton]);

  const onSearch = url => {
    const editedquery = url.replace(/\+/g, "%2B");
    history.replace(`${location.pathname}?${editedquery}`);
  };

  const handleChangeAnalytic = event => {
    const { value } = event.target;
    if (value !== selectedAnalytic) {
      window.stopVisualisation();
    }
    setSelectedAnalytic(value);
  };

  const clearFilter = () => {
    const query = defaultQuery();
    const editedquery = query.replace(/\+/g, "%2B");
    const listElm = document.querySelector("#infinite-list");
    listElm.scrollTo({
      top: 0,
      left: 0,
      behavior: "smooth"
    });
    history.replace(`${location.pathname}?${editedquery}`);
  };

  function deleteData(id, node) {
    deleteCamera(id, node).then(result => {
      if (result.code === 200) {
        history.push("/camera");
      }
    });
  }

  const renderStreamer = () => (
    <StreamerWrapper>
      {Object.keys(data).length > 0 ? (
        <VanillaVisualizationWrapper
          streamID={data.stream_id}
          selectedAnalytic={selectedAnalytic}
          streamName={`${data.stream_name}${analyticName}`}
          history={history}
          deleteData={() => deleteData(data.stream_id, data.stream_node_num)}
          selectedNode={data.stream_node_num}
        />
      ) : (
        <BlankContainer>No Camera Data</BlankContainer>
      )}
    </StreamerWrapper>
  );

  return (
    <Fragment>
      <PageWrapper title="CAMERA DETAIL" history={history} close={true}>
        <WrapEventHistory>
          <WrapFilter>
            <ItemContainer fullWidth>{renderStreamer()}</ItemContainer>
            {Object.keys(data).length > 0 && (
              <ItemContainerScroll>
                <AreaSection
                  title="SELECT ANALYTIC RESTREAM"
                  titleColor={themeContext.mint}
                  border
                >
                  <Fragment>
                    <Row
                      justify="space-between"
                      align="center"
                      horizontalPadding={15}
                      height="48px"
                      verticalMargin={20}
                    >
                      <Input
                        key="select-analytic"
                        id="select-analytic"
                        type="select"
                        option={analyticList}
                        placeholder="Select Analytic"
                        onChange={event => handleChangeAnalytic(event)}
                        value={selectedAnalytic}
                      />
                    </Row>
                    {!showResetButton && selectedAnalytic.includes("FR") && (
                      <Row
                        justify="space-between"
                        align="center"
                        horizontalPadding={15}
                        verticalMargin={5}
                      >
                        <Input
                          key="show-all-opt-radio"
                          id="show-all-opt-radio"
                          type="checkbox"
                          text="Only Show Recognized Events"
                          checked={hideUnrecog}
                          value={hideUnrecog}
                          onChange={() => setHideUnrecog(!hideUnrecog)}
                        ></Input>
                      </Row>
                    )}
                  </Fragment>
                </AreaSection>
                {getInsightConstant(selectedAnalytic).length > 0 && (
                  <AreaSection
                    title="INSIGHT"
                    titleColor={themeContext.mint}
                    border
                  >
                    <InsightWrapper>
                      {getInsightConstant(selectedAnalytic).map(insight => (
                        <InsightBox
                          key={insight.label}
                          label={insight.label}
                          keyword={insight.keyword}
                          isRealtime={insight.realtime}
                          analyticID={selectedAnalytic}
                          newEvent={newEvent}
                          streamID={data.stream_id || ""}
                        />
                      ))}
                    </InsightWrapper>
                  </AreaSection>
                )}
              </ItemContainerScroll>
            )}
          </WrapFilter>
          <WrapInfiniteData>
            <SectionAreaTitle>
              <AreaTitle>Camera Event History</AreaTitle>
              <ButtonWrapper>
                {showResetButton && (
                  <Button
                    id="filter-event"
                    onClick={() => clearFilter()}
                    type="secondary"
                    className="button-control icon-text"
                    style={{ marginRight: 10 }}
                  >
                    <img alt="search-icon" src={RefreshIcon} />
                    Reset Filter
                  </Button>
                )}
                <Button
                  id="reset-filter-event"
                  onClick={() => setOpenFilterMenu(true)}
                  type="blue"
                  className="button-control icon-text"
                >
                  <img alt="search-icon" src={SearchIcon} />
                  Filter Event
                </Button>
              </ButtonWrapper>
            </SectionAreaTitle>
            <EventWrapper
              query={history.location.search}
              newDataEvent={newEvent}
              isLive={!showResetButton}
              hideUnrecog={hideUnrecog}
            />
          </WrapInfiniteData>
        </WrapEventHistory>
      </PageWrapper>
      <FilterEvent
        history={history}
        openModal={openFilterMenu}
        onClose={() => setOpenFilterMenu(false)}
        onSubmit={query => onSearch(query)}
      />
    </Fragment>
  );
}

export default withModalVisualization(withEnrollment(CameraDetail));

CameraDetail.propTypes = {
  history: PropTypes.object.isRequired,
  location: PropTypes.object.isRequired
};

const ButtonWrapper = Styled.div`
  display: flex;
  flex-direction: row;
`;

const StreamerWrapper = Styled.div`
  width: 100%;
  height: 300px;
  display: flex;
`;

const ItemContainer = Styled.div`
  padding: 16px 16px 0px;
  video-display #video-container{
    width: ${props => (props.fullWidth ? `100%` : `auto`)};
  }
`;

const ItemContainerScroll = Styled(ItemContainer)`
  height: calc(100% - 400px);
  overflow-y: auto; 
  max-height: 100%;
`;

const WrapEventHistory = Styled.div`
  display:flex;
  width: 100%;
  height: 100%;
  position: relative;
`;

const WrapFilter = Styled.div`
  border-right: solid 1px ${props => props.theme.secondary2};
  background-color: ${props => props.theme.bg};
  position: relative;
  width: 40%;
  min-width: 40%;
`;

const WrapInfiniteData = Styled.div`
  width: 70%;
  height: 94.7%;
  background-color: ${props => props.theme.bg};
  position: relative;
`;

const SectionAreaTitle = Styled.div`
    height: 40px;
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: space-between;
    border-bottom: 1px solid #372463;
    padding-left: 13px;
    padding-right: 13px;
`;

const AreaTitle = Styled.div`
    color: #45E5B7;
    font-weight: bold;
    font-size: 12px;
    line-height: 14px;
    text-transform: uppercase;
`;

const InsightWrapper = Styled.div`
  width: 100%;
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  grid-template-rows: repeat(auto-fill, 96px);
  grid-row-gap: 10px;
  grid-column-gap: 10px;
`;

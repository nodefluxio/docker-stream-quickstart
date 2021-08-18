/* eslint-disable no-unused-vars */
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
import {
  REACT_APP_API_WEB_SOCKET,
  REACT_APP_API_CAMERA,
  REACT_APP_REFRESH_TIME
} from "config";
import { eventDate, dateToIso, isoToDate } from "helpers/dateTime";
import qs from "qs";
import dayjs from "dayjs";

import BlankContainer from "components/atoms/BlankContainer";
import LoadingSpinner from "components/atoms/LoadingSpinner";
import Input from "components/molecules/Input";

import Text from "components/atoms/Text";
import Row from "components/atoms/Row";
import AreaSection from "components/molecules/AreaSection";
import Event from "components/molecules/Event";
import { FilterEvent } from "components/organisms/modals";
import VisualisationWrapper from "components/molecules/VisualisationWrapper";
import Accordion from "components/molecules/Accordion";
import { getCamera, getEvent } from "api";

import NoRecord from "assets/icon/visionaire/no-record.svg";
import SearchIcon from "assets/icon/visionaire/search.svg";
import AlertSound from "assets/audio/alert.ogg";
import Button from "components/molecules/Button";
import RefreshIcon from "assets/icon/visionaire/refresh.svg";

import { showVisualisation, stopVisualisation } from "assets/js/visualstreamer";

import theme from "theme";
import PageWrapper from "components/organisms/pageWrapper";
import { ANALYTIC_METADATA } from "../../constants/analyticMetadata";

export default function CameraDetail(props) {
  const { history, location } = props;
  const getParamsUrl = location.search.replace("?", "");
  const defaultValue = qs.parse(getParamsUrl);

  const [dataEvent, setDataEvent] = useState([]);
  const [data, setData] = useState({});
  const [noData, setNoData] = useState(false);
  const themeContext = useContext(ThemeContext);
  const [newEvent, setNewEvent] = useState(null);
  const [analyticList, setAnalyticList] = useState([]);
  const [seatList, setSeatList] = useState([]);
  const [analyticName, setAnalyticName] = useState("");
  const [selectedAnalytic, setSelectedAnalytic] = useState(() => {
    if (defaultValue.analytic_id !== undefined) {
      setAnalyticName(
        `${
          ANALYTIC_METADATA[defaultValue.analytic_id]
            ? ` - ${ANALYTIC_METADATA[defaultValue.analytic_id].analytic_name}`
            : ""
        }`
      );
      return defaultValue.analytic_id;
    }
    return "";
  });
  const [prevAnalytic, setPrevAnalytic] = useState(() => {
    if (defaultValue.analytic_id !== undefined) {
      setAnalyticName(
        `${
          ANALYTIC_METADATA[defaultValue.analytic_id]
            ? ` - ${ANALYTIC_METADATA[defaultValue.analytic_id].analytic_name}`
            : ""
        }`
      );
      return defaultValue.analytic_id;
    }
    return "";
  });
  const [playAlert, setPlayAlert] = useState("");
  const [openFilterMenu, setOpenFilterMenu] = useState(false);
  const [lastTimestamp, setLastTimestamp] = useState(new Date());
  const [queryFilterOnly, setQueryFilterOnly] = useState("");
  const [queryString, setQueryString] = useState("");
  const [showResetButton, setShowResetButton] = useState(false);
  const socketRef = useRef(null);
  const [page, setPage] = useState(1);
  const [totalPage, setTotalPage] = useState(1);
  const [prevQuery, setPrevQuery] = useState("");
  const [prevID, setPrevID] = useState("");
  const [hideUnrecog, setHideUnrecog] = useState(true);

  const audio = new Audio(AlertSound);
  const now = new Date();
  const dateFrom = dateToIso(
    new Date(now.getFullYear(), now.getMonth(), now.getDate(), 0, 0, 0)
  );
  const dateTo = dateToIso(now);

  const LIMIT_DATA = 20;

  function defaultQuery(reset = true) {
    const paths = location.pathname.split("/");
    const id = paths[paths.length - 1];
    if (Object.keys(defaultValue).length > 0) {
      if (!reset) {
        return qs.stringify(
          {
            page: 1,
            search: defaultValue.search,
            "filter[timestamp_from]": defaultValue.filter.timestamp_from,
            "filter[timestamp_to]": defaultValue.filter.timestamp_to,
            "filter[type]": defaultValue.filter.type,
            "filter[stream_id]": id,
            limit: LIMIT_DATA
          },
          { encode: false }
        );
      }
    }
    return qs.stringify(
      {
        page: 1,
        search: "",
        "filter[timestamp_from]": dateFrom,
        "filter[timestamp_to]": dateTo,
        "filter[type]":
          hideUnrecog && selectedAnalytic.includes("FR") ? "recognized" : "",
        "filter[stream_id]": id,
        limit: LIMIT_DATA
      },
      { encode: false }
    );
  }

  const onConnect = id => {
    socketRef.current = new WebSocket(
      `${REACT_APP_API_WEB_SOCKET}/event_channel?stream_id=${id}`
    );
    socketRef.current.onopen = function() {
      console.log("Connected");
    };
    socketRef.current.onerror = function(event) {
      setTimeout(() => {
        onConnect(id);
      }, 3000);
    };
    socketRef.current.onmessage = async function(event) {
      const msg = JSON.parse(event.data);
      await setNewEvent(msg);
    };
    socketRef.current.onclose = function(event) {
      console.log("closing socket.... code: ", event.code);
      socketRef.current = null;
      if (event.code !== 1006 || event.code !== 1000) {
        setTimeout(() => {
          onConnect(id);
        }, 3000);
      }
    };
  };

  function combineEvent(oldData, newData) {
    Array.prototype.push.apply(
      oldData[oldData.length - 1].data,
      newData[0].data
    );
    newData.shift();
    return oldData.concat(newData);
  }

  function getEventHistoryData(query, isFilterOn) {
    const editedquery = query.replace(/\+/g, "%2B");
    if (query !== queryString) {
      setNoData(false);
      getEvent(editedquery).then(result => {
        setTotalPage(result.results.total_page);
        const { events } = result.results;
        if (events.length > 0) {
          if (queryFilterOnly !== query.replace(/page=\d+/g, "")) {
            if (isFilterOn) {
              setDataEvent(events);
            } else {
              setDataEvent(events[0].data);
              setLastTimestamp(events[0].timestamp);
            }
          } else {
            const firstTimestamp = new Date(events[0].timestamp);
            if (
              lastTimestamp.setHours(0, 0, 0, 0) ===
              firstTimestamp.setHours(0, 0, 0, 0)
            ) {
              setDataEvent(combineEvent(dataEvent, events));
            } else {
              setDataEvent(dataEvent.concat(events));
            }
          }
        } else {
          setDataEvent([]);
        }
      });
    }
    setQueryString(query);
    setQueryFilterOnly(query.replace(/page=\d+/g, ""));
  }

  useEffect(() => {
    if (socketRef.current) {
      socketRef.current.close();
    }
    const paths = location.pathname.split("/");
    const id = paths[paths.length - 1];
    if (prevID !== id) {
      onConnect(id);
      getCamera(id)
        .then(result => {
          setData(result);
          if (result.pipelines !== undefined) {
            const listPipeline = [];
            const seats = [];
            result.seats.forEach(value => {
              const analyticId = value.analytic_id;
              seats[analyticId] = value.serial_number;
              listPipeline.push({
                value: analyticId,
                label: ANALYTIC_METADATA[analyticId]
                  ? ANALYTIC_METADATA[analyticId].analytic_name
                  : ""
              });
            });
            setSeatList(seats);
            setAnalyticList(listPipeline);
            setSelectedAnalytic(listPipeline[0].value);
          }
          showVisualisation();
        })
        .catch(() => {});
    }
    setPrevID(id);
    return () => {
      if (socketRef.current) {
        socketRef.current.close();
      }
      stopVisualisation();
    };
  }, [location.pathname]);

  useEffect(() => {
    if (playAlert !== "" && playAlert !== undefined) {
      audio.play();
    }
  }, [playAlert]);

  useEffect(() => {
    const interval = setInterval(() => {
      window.location.reload();
    }, REACT_APP_REFRESH_TIME);
    return () => clearInterval(interval);
  }, [REACT_APP_REFRESH_TIME]);

  useEffect(() => {
    if (selectedAnalytic !== "" && selectedAnalytic !== prevAnalytic) {
      setPrevAnalytic(selectedAnalytic);
      const query = defaultQuery(false);
      const editedquery = query.replace(/\+/g, "%2B");
      history.replace(
        `${location.pathname}?${editedquery}&analytic_id=${selectedAnalytic}`
      );
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
    if (location.search !== "" && prevQuery !== location.search) {
      setDataEvent([]);
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
      getEventHistoryData(getParamsUrl, flag);
    }
  }, [location.search]);

  useEffect(() => {
    if (newEvent !== null && !showResetButton) {
      const analyticType = newEvent.analytic_id;
      const analytic = analyticType.match(/[^-]*$/g);
      const oldDataEvent = dataEvent;
      if (hideUnrecog && analytic[0] === "FR") {
        if (newEvent.label !== "unrecognized" && dataEvent.length > 50) {
          oldDataEvent.pop();
          setDataEvent([newEvent, ...oldDataEvent]);
        } else if (newEvent.label !== "unrecognized") {
          setDataEvent([newEvent, ...dataEvent]);
        }
      } else if (dataEvent.length > 50) {
        oldDataEvent.pop();
        setDataEvent([newEvent, ...oldDataEvent]);
      } else {
        setDataEvent([newEvent, ...dataEvent]);
      }
      if (newEvent.label !== "unrecognized" && analytic[0] === "FR") {
        setPlayAlert(newEvent.secondary_image);
      }
    }
  }, [newEvent]);

  useEffect(() => {
    if (dataEvent.length > 0 && showResetButton) {
      setNoData(false);
      const lastTime = dataEvent[dataEvent.length - 1].timestamp;
      setLastTimestamp(new Date(lastTime));
    } else if (dataEvent.length === 0) {
      setNoData(true);
    }
  }, [dataEvent]);

  useEffect(() => {
    if (socketRef.current.readyState > 1) {
      const paths = location.pathname.split("/");
      const id = paths[paths.length - 1];
      if (showResetButton) {
        socketRef.current.close();
      } else {
        onConnect(id);
      }
    }
  }, [showResetButton]);

  const formatLabel = string => {
    const splitString = string.split("_");
    const keys =
      splitString[1].charAt(0).toUpperCase() + splitString[1].slice(1);
    return `Camera ${keys}`;
  };

  function handleInfiniteScroll() {
    const listElm = document.querySelector("#infinite-list");
    if (listElm.scrollTop === 0 && showResetButton) {
      setPage(1);
    } else if (
      listElm.scrollTop + listElm.clientHeight >=
      listElm.scrollHeight
    ) {
      if (page < totalPage && showResetButton) {
        setPage(page + 1);
      }
    }
    return null;
  }

  const onSearch = url => {
    const editedquery = url.replace(/\+/g, "%2B");
    history.replace(
      `${location.pathname}?${editedquery}&analytic_id=${selectedAnalytic}`
    );
  };

  const handleChangeAnalytic = event => {
    const { value } = event.target;
    setSelectedAnalytic(value);
    history.push(`${location.pathname}${location.search}&analytic_id=${value}`);
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
    history.replace(
      `${location.pathname}?${editedquery}&analytic_id=${selectedAnalytic}`
    );
  };

  useEffect(() => {
    if (page > 1) {
      const { search } = location;
      const filter = qs.parse(search.replace("?", ""));
      filter.page = page;
      const filterUrl = qs.stringify(filter);
      history.replace(
        `${location.pathname}?${filterUrl}&analytic_id=${selectedAnalytic}`
      );
    }
  }, [page]);

  const renderStreamer = () => (
    <StreamerWrapper>
      {Object.keys(data).length > 0 ? (
        <VisualisationWrapper
          url={`${REACT_APP_API_CAMERA}/mjpeg/0/${data.stream_id}/${selectedAnalytic}`}
          width="100%"
          height="100%"
          name={`${data.stream_name}${analyticName}`}
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
              <ItemContainer>
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
                <AreaSection
                  title="CAMERA INFORMATION"
                  titleColor={themeContext.mint}
                  border
                >
                  <div>
                    {Object.keys(data).map(item => {
                      if (typeof data[item] !== "object") {
                        return (
                          <Fragment key={item}>
                            <Row
                              justify="space-between"
                              align="center"
                              horizontalPadding={15}
                              height="48px"
                              key={item}
                            >
                              <Text
                                size="14px"
                                color={themeContext.white}
                                weight={500}
                              >
                                {formatLabel(item)}
                              </Text>
                              <Text
                                size="14px"
                                color={themeContext.white}
                                weight={500}
                              >
                                {data[item]}
                              </Text>
                            </Row>
                            <Divider color={themeContext.secondary2} />
                          </Fragment>
                        );
                      }
                      return null;
                    })}
                  </div>
                </AreaSection>
                <AreaSection
                  title="LICENSE INFORMATION"
                  titleColor={themeContext.mint}
                  border
                >
                  <Row
                    justify="space-between"
                    align="center"
                    horizontalPadding={15}
                    height="48px"
                    key="serial-number"
                  >
                    <Text size="14px" color={themeContext.white} weight={500}>
                      Serial Number
                    </Text>
                    <Text size="14px" color={themeContext.white} weight={500}>
                      {seatList[selectedAnalytic]}
                    </Text>
                  </Row>
                </AreaSection>
              </ItemContainer>
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
            <EventContainer id="infinite-list" onScroll={handleInfiniteScroll}>
              {!showResetButton && (
                <PurpleRibbon>
                  {dayjs(lastTimestamp).format("dddd D MMMM YYYY")}
                </PurpleRibbon>
              )}
              {dataEvent.length > 0 ? (
                dataEvent.map((val, index) =>
                  showResetButton ? (
                    <Accordion
                      key={index}
                      header={dayjs(val.timestamp).format("dddd D MMMM YYYY")}
                    >
                      {val.data.map(event => (
                        <Event
                          key={`${event.id}-${index}`}
                          primaryImage={
                            event.primary_image !== ""
                              ? `data:image/jpeg;base64,${event.primary_image}`
                              : ""
                          }
                          secondaryImage={
                            event.secondary_image !== ""
                              ? `data:image/jpeg;base64,${event.secondary_image}`
                              : ""
                          }
                          label={event.label}
                          result={event.result}
                          location={event.location}
                          date={eventDate(event.timestamp)}
                          color={
                            event.label === "recognized"
                              ? theme.success
                              : theme.white
                          }
                        />
                      ))}
                    </Accordion>
                  ) : (
                    <Event
                      key={`${val.id}-${index}`}
                      primaryImage={
                        val.primary_image !== ""
                          ? `data:image/jpeg;base64,${val.primary_image}`
                          : ""
                      }
                      secondaryImage={
                        val.secondary_image !== ""
                          ? `data:image/jpeg;base64,${val.secondary_image}`
                          : ""
                      }
                      label={val.label}
                      result={val.result}
                      location={val.location}
                      date={eventDate(val.timestamp)}
                      color={
                        val.label === "recognized" ? theme.success : theme.white
                      }
                    />
                  )
                )
              ) : (
                <BlankContainer>
                  {noData ? (
                    <BlankContainer>
                      <img src={NoRecord} alt="blank-icon" />
                      No Data Found
                    </BlankContainer>
                  ) : (
                    <LoadingSpinner show={!noData} />
                  )}
                </BlankContainer>
              )}
            </EventContainer>
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

const EventContainer = Styled.div`
  width: 100%;
  height: 100%;
  overflow-y: auto;
  overflow-x: hidden;
  max-height: 100%;
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

const Divider = Styled.div`
  width: 100%;
  height: 1px;
  ${({ color }) => color && `background-color: ${color};`}
`;

const PurpleRibbon = Styled.div`
  width: 100%;
  height: 40px;
  background-color: ${props => props.theme.secondary2};
  cursor: pointer;
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  align-items: center;
  font-weight: 700;
  text-transform: capitalize;
  font-size: 14px;
  line-height: 14.4px;
  color: ${props => props.theme.white};
  padding-left: 15px;
`;

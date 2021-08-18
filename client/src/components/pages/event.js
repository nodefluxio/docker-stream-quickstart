import React, { useEffect, useState } from "react";
import { connect, useDispatch } from "react-redux";
import Styled from "styled-components";
import PropTypes from "prop-types";
import qs from "qs";
import dayjs from "dayjs";
import PageWrapper from "components/organisms/pageWrapper";
import Accordion from "components/molecules/Accordion";
import Input from "components/molecules/Input";
import Button from "components/molecules/Button";
import Event from "components/molecules/Event";
import BlankContainer from "components/atoms/BlankContainer";
import Text from "components/atoms/Text";
import LoadingSpinner from "components/atoms/LoadingSpinner";

import SearchIcon from "assets/icon/visionaire/search.svg";
import DateIcon from "assets/icon/visionaire/date.svg";
import DownloadIcon from "assets/icon/visionaire/download.svg";

import { dateToIso, isoToDate, eventDate } from "helpers/dateTime";

import { getEvent } from "api";
import theme from "theme";
import IconPlaceholder from "components/atoms/IconPlaceholder";
import { requestEventExport } from "store/actions/eventDownloader";

const LIMIT_DATA = 20;

function EventHistory(props) {
  const { history, exportEvent } = props;
  const dispatch = useDispatch();
  const now = new Date();
  const defaultDateTo = dateToIso(now);
  const defaultDateFrom = dateToIso(
    new Date(now.getFullYear(), now.getMonth(), now.getDate(), 0, 0, 0)
  );
  const getParamsUrl = history.location.search.replace("?", "");
  const defaultValue = qs.parse(getParamsUrl);
  const defaultType = [
    {
      value: "",
      label: "All"
    },
    {
      value: "unrecognized",
      label: "Unrecognized"
    },
    {
      value: "recognized",
      label: "Recognized"
    }
  ];
  const [disableButton, setDisableButton] = useState(false);
  const [dataEvent, setDataEvent] = useState([]);
  const [selectedType, setselectedType] = useState(
    defaultValue.filter ? defaultValue.filter.type : "recognized"
  );
  const [searchValue, setSearchValue] = useState(defaultValue.search || "");
  const [startDate, setStartDate] = useState(
    defaultValue.filter
      ? defaultValue.filter.timestamp_from &&
          isoToDate(defaultValue.filter.timestamp_from)
      : new Date(now.getFullYear(), now.getMonth(), now.getDate(), 0, 0, 0)
  );
  const [endDate, setEndDate] = useState(
    defaultValue.filter
      ? defaultValue.filter.timestamp_to &&
          isoToDate(defaultValue.filter.timestamp_to)
      : now
  );
  const [page, setPage] = useState(defaultValue.page || 1);
  const [lastTimestamp, setLastTimestamp] = useState(new Date(null));
  const [totalPage, setTotalPage] = useState(1);
  const [queryString, setQueryString] = useState("");
  const [queryFilterOnly, setQueryFilterOnly] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [errMsg, setErrMsg] = useState("");
  const [disableExportButton, setDisableExportButton] = useState(false);

  function getQuery() {
    const dateFrom = startDate ? dateToIso(startDate) : "";
    const dateTo = endDate ? dateToIso(endDate) : "";
    return qs.stringify(
      {
        page,
        search: searchValue,
        "filter[timestamp_from]": dateFrom || defaultDateFrom,
        "filter[timestamp_to]": dateTo || defaultDateTo,
        "filter[type]": selectedType,
        limit: LIMIT_DATA
      },
      { encode: false }
    );
  }

  function applyDate(name, date) {
    if (name === "time_from_field") {
      setStartDate(date);
    } else {
      setEndDate(date);
    }
  }

  function handleInfiniteScroll() {
    const listElm = document.querySelector("#infinite-list");
    if (listElm.scrollTop === 0) {
      setPage(1);
    } else if (
      listElm.scrollTop + listElm.clientHeight >=
      listElm.scrollHeight
    ) {
      if (page < totalPage) {
        setPage(page + 1);
      }
    }
    return null;
  }

  function combineEvent(oldData, newData) {
    Array.prototype.push.apply(
      oldData[oldData.length - 1].data,
      newData[0].data
    );
    newData.shift();
    return oldData.concat(newData);
  }

  function getData() {
    const query = getQuery();
    const editedquery = query.replace(/\+/g, "%2B");
    history.replace(`/event-history?${editedquery}`);
    if (query !== queryString) {
      setIsLoading(true);
      getEvent(editedquery)
        .then(result => {
          setTotalPage(result.results.total_page);
          const newEvent = result.results.events;
          if (newEvent.length > 0) {
            if (queryFilterOnly !== query.replace(/page=\d+/g, "")) {
              setDataEvent(newEvent);
              setIsLoading(false);
            } else {
              const firstTimestamp = new Date(newEvent[0].timestamp);
              if (
                lastTimestamp.setHours(0, 0, 0, 0) ===
                firstTimestamp.setHours(0, 0, 0, 0)
              ) {
                setDataEvent(combineEvent(dataEvent, newEvent));
                setIsLoading(false);
              } else {
                setDataEvent(dataEvent.concat(newEvent));
                setIsLoading(false);
              }
            }
          } else {
            setDataEvent([]);
            setIsLoading(false);
            setErrMsg("");
          }
        })
        .catch(err => {
          setErrMsg(err.message);
          setIsLoading(false);
        });
    }
    setQueryString(query);
    setQueryFilterOnly(query.replace(/page=\d+/g, ""));
  }

  function submitFilter() {
    setDisableButton(true);
    setLastTimestamp(new Date(null));
    const listElm = document.querySelector("#infinite-list");
    listElm.scrollTo({
      top: 0,
      left: 0,
      behavior: "smooth"
    });
    getData();
  }

  function keyHandler(targetKey) {
    if (targetKey.key === "Enter") {
      submitFilter();
    }
  }

  useEffect(() => {
    window.addEventListener("keydown", keyHandler);
    return () => {
      document.removeEventListener("keydown", keyHandler);
    };
  }, []);

  useEffect(() => {
    getData();
  }, [page, history.location.search]);

  useEffect(() => {
    if (dataEvent.length > 0) {
      const lastTime = dataEvent[dataEvent.length - 1].timestamp;
      setLastTimestamp(new Date(lastTime));
    }
  }, [dataEvent]);

  useEffect(() => {
    setDisableExportButton(exportEvent.disableDownload);
  }, [exportEvent]);

  return (
    <PageWrapper
      title="Search Event History"
      history={history}
      buttonGroup={
        dataEvent.length > 0 && (
          <Button
            disabled={disableExportButton}
            style={{ marginRight: "15px" }}
            onClick={() => {
              setDisableExportButton(true);
              dispatch(
                requestEventExport(history.location.search.replace("?", ""))
              );
            }}
          >
            <IconPlaceholder>
              <img alt="download" src={DownloadIcon} />
            </IconPlaceholder>
            EXPORT EVENT
          </Button>
        )
      }
    >
      <Content>
        <RightSide>
          <FilterWrapper>
            <Accordion header="Wild Search" icon={SearchIcon}>
              <Input
                type="text"
                label="Enter Search Key"
                value={searchValue}
                onChange={e => setSearchValue(e.target.value)}
                onKeyPress={e => {
                  if (e.which === 13) {
                    submitFilter();
                  }
                }}
              />
            </Accordion>
            <Accordion header="Find By Time Range" icon={DateIcon}>
              <DateWrapper>
                <Input
                  type="date"
                  label="Start"
                  time={true}
                  style={{ marginRight: "10px" }}
                  name="time_from_field"
                  selectsStart
                  startDate={startDate}
                  endDate={endDate}
                  onSubmit={(name, value) => applyDate(name, value)}
                  popperPlacement="left-center"
                  value={startDate}
                />
                <Input
                  type="date"
                  label="End"
                  time={true}
                  name="time_to_field"
                  selectsEnd
                  startDate={startDate}
                  endDate={endDate}
                  onSubmit={(name, value) => applyDate(name, value)}
                  popperPlacement="right-center"
                  value={endDate}
                />
              </DateWrapper>
            </Accordion>
            <Accordion header="Find By Type">
              <Input
                type="select"
                option={defaultType}
                placeholder="All"
                onChange={e => setselectedType(e.target.value)}
                value={selectedType}
              />
            </Accordion>
          </FilterWrapper>
          <Button
            style={{ margin: "15px" }}
            disable={disableButton}
            onClick={() => submitFilter()}
          >
            FIND NOW
          </Button>
        </RightSide>
        <LeftSide>
          <EventWrapper id="infinite-list" onScroll={handleInfiniteScroll}>
            {dataEvent.length > 0 ? (
              dataEvent.map((event, index) => (
                <Accordion
                  key={index}
                  header={dayjs(event.timestamp).format("dddd D MMMM YYYY")}
                >
                  {event.data.map(data => (
                    <Event
                      key={data.id}
                      primaryImage={`data:image/jpeg;base64,${data.primary_image}`}
                      secondaryImage={
                        data.secondary_image !== ""
                          ? `data:image/jpeg;base64,${data.secondary_image}`
                          : ""
                      }
                      label={data.label}
                      result={data.result}
                      location={data.location}
                      date={eventDate(data.timestamp)}
                      color={
                        data.label === "recognized"
                          ? theme.success
                          : theme.white
                      }
                    />
                  ))}
                </Accordion>
              ))
            ) : (
              <BlankContainer>
                {isLoading ? (
                  <LoadingSpinner show={isLoading} />
                ) : (
                  <Text size="18" color="#fff" weight="600">
                    {errMsg === "" ? "NO DATA FOUND" : errMsg}
                  </Text>
                )}
              </BlankContainer>
            )}
          </EventWrapper>
        </LeftSide>
      </Content>
    </PageWrapper>
  );
}

EventHistory.propTypes = {
  history: PropTypes.object.isRequired,
  exportEvent: PropTypes.object.isRequired
};

function mapStateToProps(state) {
  return {
    exportEvent: state.exportEvent
  };
}

export default connect(mapStateToProps)(EventHistory);

const Content = Styled.div`
  width: 100%;
  height: calc(100% - 64px);
  display: flex;
  flex-direction: row;
`;

const RightSide = Styled.div`
  width: 30%;
  min-height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  border-right: 1px solid ${props => props.theme.secondary2};
`;

const LeftSide = Styled.div`
  width: 70%;
`;

const FilterWrapper = Styled.div`
  display: flex;
  flex-direction: column;
`;

const DateWrapper = Styled.div`
  display: flex;
  flex-direction: row;
`;

const EventWrapper = Styled.div`
  max-height: 100%;
  height: 100%;
  overflow: auto;
`;

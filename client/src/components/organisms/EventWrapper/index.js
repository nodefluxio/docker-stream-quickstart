import React, { useState, useEffect } from "react";
import PropTypes from "prop-types";
import Styled from "styled-components";
import dayjs from "dayjs";
import qs from "qs";
import { useDispatch } from "react-redux";

import Event from "components/molecules/Event";
import BlankContainer from "components/atoms/BlankContainer";
import Text from "components/atoms/Text";
import LoadingSpinner from "components/atoms/LoadingSpinner";
import Accordion from "components/molecules/Accordion";

import AlertSound from "assets/audio/alert.ogg";

import theme from "theme";
import { getEvent } from "api";
import { eventDate } from "helpers/dateTime";
import { setDataContext } from "store/actions/eventContext";

export default function EventWrapper(props) {
  const { query, cbDataEvent, newDataEvent, isLive, hideUnrecog } = props;
  const dispatch = useDispatch();
  const [playAlert, setPlayAlert] = useState("");
  const [lastID, setLastID] = useState(null);
  const [dataEvent, setDataEvent] = useState([]);
  const [isLoading, setIsLoading] = useState(false);
  const [errMsg, setErrMsg] = useState("");
  const [queryString, setQueryString] = useState("");
  const [queryFilterOnly, setQueryFilterOnly] = useState("");
  const [lastTimestamp, setLastTimestamp] = useState(new Date());
  const [lastScrollTop, setLastScrollTop] = useState(0);

  const audio = new Audio(AlertSound);

  function handleInfiniteScroll() {
    const listElm = document.querySelector("#infinite-list");
    const scrollTop = listElm.scrollTop + listElm.clientHeight;
    if (scrollTop >= listElm.scrollHeight && scrollTop > lastScrollTop) {
      const lastDateGroupEvent = dataEvent[dataEvent.length - 1];
      const lastDataEventLength = lastDateGroupEvent.data.length;
      const lastDataEvent = lastDateGroupEvent.data[lastDataEventLength - 1];
      setLastID(lastDataEvent.id);
    } else if (listElm.scrollTop === 0) {
      setLastID(null);
    }
    setLastScrollTop(scrollTop);
    return null;
  }

  function combineEvent(oldData, newData) {
    Array.prototype.push.apply(
      oldData[oldData.length - 1].data,
      newData[0].data
    );
    newData.shift();
    const combinedData = oldData.concat(newData);
    return combinedData;
  }

  useEffect(() => {
    if (playAlert !== "" && playAlert !== undefined) {
      audio.play();
    }
  }, [playAlert]);

  useEffect(() => {
    if (Object.keys(newDataEvent).length > 0 && isLive) {
      const analyticType = newDataEvent.analytic_id;
      const analytic = analyticType.match(/[^-]*$/g);
      const oldDataEvent = dataEvent;
      const timestamp = dayjs(new Date()).format("dddd D MMMM YYYY");
      let oldTimestamp = dayjs(new Date()).format("dddd D MMMM YYYY");
      let selectedAnalytic = "";
      const editedquery = query.replace("?", "");

      const defaultValue = qs.parse(editedquery);
      if (defaultValue.filter) {
        if (defaultValue.filter.analytic_id) {
          selectedAnalytic = defaultValue.filter.analytic_id;
        }
      }
      if (dataEvent.length > 0) {
        if (dataEvent[0].timestamp) {
          oldTimestamp = dayjs(new Date(dataEvent[0].timestamp)).format(
            "dddd D MMMM YYYY"
          );
        }
        if (analyticType === selectedAnalytic) {
          if (dataEvent[0].data.length > 50) {
            oldDataEvent[0].data.pop();
          }
          if (hideUnrecog && analytic[0] === "FR") {
            if (newDataEvent.label !== "unrecognized") {
              if (timestamp === oldTimestamp) {
                setDataEvent([
                  {
                    timestamp,
                    data: [newDataEvent, ...oldDataEvent[0].data]
                  }
                ]);
              } else {
                setDataEvent([
                  {
                    timestamp,
                    data: [newDataEvent]
                  },
                  ...dataEvent[0]
                ]);
              }
            }
          } else if (timestamp === oldTimestamp) {
            setDataEvent([
              {
                timestamp,
                data: [newDataEvent, ...oldDataEvent[0].data]
              }
            ]);
          } else if (timestamp !== oldTimestamp) {
            setDataEvent([
              {
                timestamp,
                data: [newDataEvent]
              },
              ...dataEvent[0]
            ]);
          }
        }
      } else if (analyticType === selectedAnalytic) {
        if (hideUnrecog && analytic[0] === "FR") {
          if (newDataEvent.label !== "unrecognized") {
            setDataEvent([{ timestamp, data: [newDataEvent] }]);
          }
        } else {
          setDataEvent([{ timestamp, data: [newDataEvent] }]);
        }
      }
      if (newDataEvent.label !== "unrecognized" && analytic[0] === "FR") {
        setPlayAlert(newDataEvent.secondary_image);
      }
    }
  }, [newDataEvent]);

  useEffect(() => {
    if (dataEvent.length > 0) {
      const lastTime = dataEvent[dataEvent.length - 1].timestamp;
      setLastTimestamp(new Date(lastTime));
      cbDataEvent(dataEvent);
    }
  }, [dataEvent]);

  useEffect(() => {
    const editedquery = query.replace("?", "");
    let queryToEvent = qs.parse(editedquery);
    const isNewFilter =
      queryFilterOnly !== editedquery.replace(/&last_id=\d+/g, "");
    if (lastID && !isNewFilter) {
      queryToEvent.last_id = lastID;
    }
    queryToEvent = qs.stringify(queryToEvent);
    if (queryToEvent !== queryString) {
      setIsLoading(true);
      getEvent(queryToEvent)
        .then(result => {
          const newData = result.results.events;
          if (newData.length > 0) {
            if (queryFilterOnly !== queryToEvent.replace(/&last_id=\d+/g, "")) {
              setDataEvent(newData);
              setIsLoading(false);
            } else {
              const firstTimestamp = new Date(newData[0].timestamp);
              if (
                lastTimestamp.setHours(0, 0, 0, 0) ===
                firstTimestamp.setHours(0, 0, 0, 0)
              ) {
                setDataEvent(combineEvent(dataEvent, newData));
                setIsLoading(false);
              } else {
                setDataEvent(dataEvent.concat(newData));
                setIsLoading(false);
              }
            }
          } else if (
            newData.length <= 0 &&
            queryFilterOnly !== queryToEvent.replace(/&last_id=\d+/g, "")
          ) {
            setDataEvent([]);
            setIsLoading(false);
            setErrMsg("");
          } else {
            setIsLoading(false);
            setErrMsg("");
          }
        })
        .catch(err => {
          setErrMsg(err.message);
          setIsLoading(false);
        });
    }
    setQueryString(queryToEvent);
    setQueryFilterOnly(queryToEvent.replace(/&last_id=\d+/g, ""));
  }, [query, lastID]);

  return (
    <Wrapper id="infinite-list" onScroll={handleInfiniteScroll}>
      {dataEvent.length > 0 ? (
        dataEvent.map((event, index) => (
          <Accordion
            open={true}
            key={index}
            header={dayjs(event.timestamp).format("dddd D MMMM YYYY")}
          >
            {event.data &&
              event.data.map(data => (
                <Event
                  key={data.id}
                  onContextMenu={() => dispatch(setDataContext(data))}
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
                    data.label === "recognized" ? theme.success : theme.white
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
    </Wrapper>
  );
}

EventWrapper.propTypes = {
  query: PropTypes.string.isRequired,
  cbDataEvent: PropTypes.func,
  newDataEvent: PropTypes.object,
  isLive: PropTypes.bool,
  hideUnrecog: PropTypes.bool
};

EventWrapper.defaultProps = {
  cbDataEvent: () => {},
  newDataEvent: {},
  isLive: false,
  hideUnrecog: false
};

const Wrapper = Styled.div`
  max-height: 100%;
  height: 100%;
  overflow: auto;
`;

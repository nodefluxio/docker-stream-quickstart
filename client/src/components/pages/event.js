/* eslint-disable array-callback-return */
import React, { useEffect, useState } from "react";
import { connect, useDispatch } from "react-redux";
import Styled from "styled-components";
import PropTypes from "prop-types";
import qs from "qs";

import PageWrapper from "components/organisms/pageWrapper";
import Accordion from "components/molecules/Accordion";
import Input from "components/molecules/Input";
import Button from "components/molecules/Button";
import withEnrollment from "components/templates/withEnrollmentModal";

import SearchIcon from "assets/icon/visionaire/search.svg";
import DateIcon from "assets/icon/visionaire/date.svg";
import DownloadIcon from "assets/icon/visionaire/download.svg";
import AnalyticIcon from "assets/icon/visionaire/analytic.svg";

import { dateToIso, isoToDate, getTimeZone } from "helpers/dateTime";

import { getListAnalytic, getListCamera } from "api";
import IconPlaceholder from "components/atoms/IconPlaceholder";
import { ANALYTIC_METADATA } from "constants/analyticMetadata";
import { requestEventExport } from "store/actions/eventDownloader";
import EventWrapper from "components/organisms/EventWrapper";

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
  const [disableExportButton, setDisableExportButton] = useState(false);
  const [analyticList, setAnalyticList] = useState([]);
  const [cameraList, setCameraList] = useState([]);
  const [selectedCamera, setSelectedCamera] = useState([]);

  function getQuery() {
    const dateFrom = startDate ? dateToIso(startDate) : "";
    const dateTo = endDate ? dateToIso(endDate) : "";
    let selectedAnalytic = [];
    let defaultAnalytics = defaultValue.filter
      ? defaultValue.filter.analytic_id || ""
      : "";
    defaultAnalytics = defaultAnalytics.split(",");
    const defaultStream = defaultValue.filter
      ? defaultValue.filter.stream_id || ""
      : "";
    let selectedCams = [];
    if (selectedCamera.length > 0) {
      selectedCamera.map(camera => selectedCams.push(camera.value));
    } else if (defaultStream.length > 0) {
      selectedCams = defaultStream;
    }

    if (analyticList.length > 0) {
      analyticList.map(analytic =>
        analytic.checked ? selectedAnalytic.push(analytic.id) : null
      );
    } else if (defaultAnalytics.length > 0) {
      selectedAnalytic = defaultAnalytics;
    }

    const queryObject = {
      "filter[timestamp_from]": dateFrom || defaultDateFrom,
      "filter[timestamp_to]": dateTo || defaultDateTo,
      "filter[type]": selectedType,
      "filter[analytic_id]": selectedAnalytic,
      "filter[stream_id]": selectedCams,
      limit: LIMIT_DATA,
      timezone: getTimeZone()
    };

    if (searchValue) {
      queryObject.search = searchValue;
    }

    return qs.stringify(queryObject, { arrayFormat: "comma" });
  }

  function applyDate(name, date) {
    if (name === "time_from_field") {
      setStartDate(date);
    } else {
      setEndDate(date);
    }
  }

  function submitFilter() {
    setDisableButton(true);
    const listElm = document.querySelector("#infinite-list");
    listElm.scrollTo({
      top: 0,
      left: 0,
      behavior: "smooth"
    });
    const query = getQuery();
    const editedquery = query.replace(/\+/g, "%2B");
    history.replace(`/event-history?${editedquery}`);
  }

  function keyHandler(targetKey) {
    if (targetKey.key === "Enter") {
      submitFilter();
    }
  }

  function handleCheck(value) {
    setAnalyticList(
      analyticList.map(analytic =>
        analytic.id === value
          ? { ...analytic, checked: !analytic.checked }
          : { ...analytic }
      )
    );
  }

  useEffect(() => {
    window.addEventListener("keydown", keyHandler);
    let defaultAnalytics = defaultValue.filter
      ? defaultValue.filter.analytic_id || ""
      : "";
    defaultAnalytics = defaultAnalytics.split(",");
    let defaultStream = defaultValue.filter
      ? defaultValue.filter.stream_id || ""
      : "";
    defaultStream = defaultStream.split(",");
    const query = getQuery();
    const editedquery = query.replace(/\+/g, "%2B");
    history.replace(`/event-history?${editedquery}`);
    getListAnalytic().then(result => {
      if (result.code === 200) {
        const newAnalyticList = [];
        result.analytics.map(analytic => {
          let codeLength = analytic.id.split("-");
          if (codeLength.length > 2) {
            codeLength.splice(codeLength.length - 1, 1);
          }
          codeLength = codeLength.join("-");
          const index = newAnalyticList.findIndex(
            data => data.id === codeLength
          );
          if (index === -1) {
            newAnalyticList.push({
              id: analytic.id,
              name: analytic.name,
              checked: defaultAnalytics.indexOf(String(analytic.id)) > -1
            });
          }
        });
        setAnalyticList(newAnalyticList);
      }
    });
    getListCamera().then(result => {
      const resultdata = result.data.streams;
      const newCameraList = [];
      const newSelectedCams = [];
      // eslint-disable-next-line array-callback-return
      resultdata.map(camera => {
        if (defaultStream.indexOf(String(camera.stream_id)) > -1) {
          newSelectedCams.push({
            value: camera.stream_id,
            label: camera.stream_name
          });
        }
        newCameraList.push({
          value: camera.stream_id,
          label: camera.stream_name
        });
      });
      setCameraList(newCameraList);
      setSelectedCamera(newSelectedCams);
    });
    return () => {
      document.removeEventListener("keydown", keyHandler);
    };
  }, []);

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
            <Accordion open={true} header="Wild Search" icon={SearchIcon}>
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
            <Accordion open={true} header="Find By Time Range" icon={DateIcon}>
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
            {cameraList.length > 0 && (
              <Accordion
                open={true}
                header="Find By Camera"
                icon={AnalyticIcon}
              >
                <Input
                  type="select"
                  option={cameraList}
                  placeholder="All"
                  onChange={e => setSelectedCamera(e)}
                  value={selectedCamera}
                  label="Camera:"
                  isMulti={true}
                  style={{ marginBottom: "40px" }}
                />
              </Accordion>
            )}
            <Accordion
              open={true}
              header="Find By Analytic Criteria"
              icon={AnalyticIcon}
            >
              <Input
                type="select"
                option={defaultType}
                placeholder="All"
                onChange={e => setselectedType(e.target.value)}
                value={selectedType}
                label="Type:"
                style={{ marginBottom: "40px" }}
              />
              {analyticList &&
                analyticList.map(analytic => (
                  <Input
                    key={analytic.id}
                    type="checkbox"
                    text={analytic.name}
                    value={analytic.id}
                    checked={analytic.checked}
                    onChange={value => handleCheck(value)}
                  />
                ))}
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
          <EventWrapper
            query={history.location.search}
            cbDataEvent={data => setDataEvent(data)}
          />
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

export default withEnrollment(connect(mapStateToProps)(EventHistory));

const Content = Styled.div`
  width: 100%;
  height: calc(100% - 64px);
  display: flex;
  flex-direction: row;
`;

const RightSide = Styled.div`
  width: 30%;
  min-height: 100%;
  height: 100%;
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
  max-height: 100%;
  overflow-y: auto;
  height: 100%;
`;

const DateWrapper = Styled.div`
  display: flex;
  flex-direction: row;
`;

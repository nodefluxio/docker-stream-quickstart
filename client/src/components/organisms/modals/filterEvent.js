import React, { useState, useEffect } from "react";
import qs from "qs";
import PropTypes from "prop-types";
import Modal from "components/atoms/Modal";
import Row from "components/atoms/Row";
import Input from "components/molecules/Input";
import Button from "components/molecules/Button";

import { dateToIso, isoToDate } from "helpers/dateTime";

const LIMIT_DATA = 20;

function FilterEvent(props) {
  const { history, openModal, onClose, onSubmit } = props;
  const now = new Date();
  const defaultDateTo = dateToIso(now);
  const defaultDateFrom = dateToIso(
    new Date(now.getFullYear(), now.getMonth(), now.getDate(), 0, 0, 0)
  );
  const [searchValue, setSearchValue] = useState("");
  const [startDate, setStartDate] = useState(defaultDateFrom);
  const [endDate, setEndDate] = useState(defaultDateTo);
  const [page, setPage] = useState(1);
  const [disableButton] = useState(false);

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
  const [selectedType, setselectedType] = useState("");

  function getQuery(id) {
    const dateFrom = startDate ? dateToIso(startDate) : "";
    const dateTo = endDate ? dateToIso(endDate) : "";
    return qs.stringify(
      {
        page,
        search: searchValue,
        "filter[timestamp_from]": dateFrom || defaultDateFrom,
        "filter[timestamp_to]": dateTo || defaultDateTo,
        "filter[type]": selectedType,
        "filter[stream_id]": id,
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

  function submitFilter() {
    const paths = history.location.pathname.split("/");
    const id = paths[paths.length - 1];
    const query = getQuery(id);
    onSubmit(query);
    setSearchValue("");
    onClose();
  }

  useEffect(() => {
    const getParamsUrl = history.location.search.replace("?", "");
    const defaultValue = qs.parse(getParamsUrl);
    if (defaultValue.search) {
      setSearchValue(defaultValue.search);
    }
    if (defaultValue.filter) {
      if (defaultValue.filter.timestamp_from) {
        setStartDate(isoToDate(defaultValue.filter.timestamp_from));
      }
      if (defaultValue.filter.timestamp_to) {
        setEndDate(isoToDate(defaultValue.filter.timestamp_to));
      }
      if (defaultValue.page) {
        setPage(defaultValue.page);
      }
    }
  }, [history]);

  return (
    <Modal
      show={openModal}
      className="modal-filter-event"
      title="Filter Event"
      close={onClose}
      padding="20px"
      width="400px"
    >
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
      <Row>
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
      </Row>
      <Input
        label="Find by Type"
        type="select"
        option={defaultType}
        placeholder="All"
        onChange={e => setselectedType(e.target.value)}
        value={selectedType}
      />

      <Button
        style={{ marginTop: "25px" }}
        disable={disableButton}
        onClick={() => submitFilter()}
        width="100% "
      >
        FIND NOW
      </Button>
    </Modal>
  );
}

FilterEvent.propTypes = {
  history: PropTypes.object.isRequired,
  openModal: PropTypes.bool.isRequired,
  onClose: PropTypes.func.isRequired,
  onSubmit: PropTypes.func.isRequired
};

export default FilterEvent;

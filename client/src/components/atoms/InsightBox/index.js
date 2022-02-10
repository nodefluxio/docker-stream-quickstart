import React, { useEffect, useState } from "react";
import PropTypes from "prop-types";
import Styled from "styled-components";
import { getInsight } from "api";
import qs from "qs";
import { getTimeZone } from "helpers/dateTime";
import LoadingSpinner from "../LoadingSpinner";

export default function InsightBox(props) {
  const {
    label,
    keyword,
    isRealtime,
    analyticID,
    width,
    newEvent,
    streamID
  } = props;
  const [value, setValue] = useState("");
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    if (streamID !== "") {
      const query = {
        "filter[analytic_id]": analyticID,
        "filter[stream_id]": streamID,
        timezone: getTimeZone()
      };
      if (!keyword.includes("total")) {
        query["filter[type]"] = keyword;
      }
      const queryStringified = qs.stringify(query, { arrayFormat: "comma" });
      setIsLoading(true);
      getInsight(queryStringified)
        .then(result => {
          if (result.ok) {
            if (!keyword.includes("total")) {
              setValue(parseInt(result.insight.total_today, 10));
            } else {
              setValue(parseInt(result.insight[keyword], 10));
            }
            setIsLoading(false);
          } else {
            setValue(0);
            setIsLoading(false);
          }
        })
        .catch(() => {
          setValue(0);
          setIsLoading(false);
        });
    }
  }, [keyword, streamID]);

  useEffect(() => {
    if (Object.keys(newEvent).length > 0 && isRealtime) {
      if (newEvent.analytic_id === analyticID) {
        if (keyword === "total_today") {
          setValue(value + 1);
        } else if (newEvent.label === keyword) {
          setValue(value + 1);
        }
      }
    }
  }, [newEvent]);
  return (
    <Wrapper width={width}>
      <Label>{label}</Label>
      <LoadingSpinner show={isLoading} />
      <Value>{value}</Value>
    </Wrapper>
  );
}

InsightBox.propTypes = {
  label: PropTypes.string.isRequired,
  keyword: PropTypes.string.isRequired,
  analyticID: PropTypes.string.isRequired,
  isRealtime: PropTypes.bool,
  width: PropTypes.string,
  newEvent: PropTypes.object,
  streamID: PropTypes.string
};

const Wrapper = Styled.div`
    width: ${props => props.width || `90%`};
    height:96px;
    border: 2px solid #372463;
    border-radius: 8px;
    font-style: normal;
    font-weight: 600;
    font-size: 14px;
    line-height: 22px;
    display: flex;
    align-items: center;
    text-transform: uppercase;
    padding: 0px 10px;
    margin: 5px;
    justify-content: space-between;
`;

const Label = Styled.div`
`;

const Value = Styled.div`
    min-width: 40px;
    color: ${props => props.theme.mint};
`;

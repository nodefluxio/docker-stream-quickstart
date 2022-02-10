/* eslint-disable no-console */
import React, { useState, useEffect, useContext } from "react";
import Styled from "styled-components";
import PropTypes from "prop-types";

import Input from "components/atoms/Input";
import Button from "components/molecules/Button";
import IconPlaceholder from "components/atoms/IconPlaceholder";
import Accordion from "components/molecules/Accordion";
import EnvWrapper, { EnvContext } from "components/templates/EnvContext";

import { getPlateResult } from "api";

import ErrorIcon from "assets/icon/visionaire/file-error.svg";

function SearchHome({ history }) {
  const API = useContext(EnvContext);
  const [licensePlate, setLicensePlate] = useState("");
  const [detailResult, setDetailResult] = useState({});
  const [errorMsg, setErrorMsg] = useState("");
  const [licenseResult, setLicenseResult] = useState("");

  function onHandleLicensePlate(value) {
    setLicensePlate(value);
  }

  function onSubmitSearch() {
    history.replace({
      search: `nopol=${licensePlate}`
    });
  }

  useEffect(() => {
    const { search } = history.location;
    if (search !== "") {
      const query = new URLSearchParams(search);
      getPlateResult(API, search)
        .then(result => {
          if (result.ok) {
            setDetailResult(result.plate_info);
            setErrorMsg("");
            setLicenseResult(query.get("nopol"));
          }
        })
        .catch(error => {
          setDetailResult({});
          setLicenseResult("");
          setErrorMsg(error.response.data.errors[0]);
        });
    }
  }, [history.location.search]);

  return (
    <Wrapper height="100%">
      <ControlContainer>Database Korlantas Search</ControlContainer>
      <Wrapper flexDirection="row" height="calc(100% - 64px)">
        <Left>
          <Accordion header="Search By License Plate" type="open">
            <Input
              type="text"
              label="Search Vehicle by License Plate"
              placeholder="Vehicle by License Plate..."
              style={{ marginBottom: "20px" }}
              value={licensePlate}
              onChange={e => onHandleLicensePlate(e.target.value)}
            />
            <Button
              onClick={() => onSubmitSearch()}
              width="100%"
              disabled={licensePlate === ""}
            >
              SEARCH
            </Button>
          </Accordion>
        </Left>
        <Right>
          <Ribbon>Search Result</Ribbon>
          <Wrapper flexDirection="row" height="calc(100% - 43px)">
            {Object.keys(detailResult).length > 0 && (
              <DetailWrapper>
                <Ribbon bg="transparent">
                  <span>Vehicle Detail</span>
                </Ribbon>
                <PlateNumber>{licenseResult}</PlateNumber>
                <DetailInfo>
                  {Object.keys(detailResult).map(data => (
                    <Row key={data}>
                      <div>{data.split("_").join(" ")}</div>
                      <div>{detailResult[data] || "-"}</div>
                    </Row>
                  ))}
                </DetailInfo>
              </DetailWrapper>
            )}
            {errorMsg !== "" && (
              <ErrorWrapper>
                <IconPlaceholder>
                  <img src={ErrorIcon} alt="no-file" />
                </IconPlaceholder>
                {errorMsg}
              </ErrorWrapper>
            )}
          </Wrapper>
        </Right>
      </Wrapper>
    </Wrapper>
  );
}

export default EnvWrapper(SearchHome);

SearchHome.propTypes = {
  history: PropTypes.object.isRequired
};

const Wrapper = Styled.div`
    display: flex;
    flex-direction: ${props => props.flexDirection || `column`};
    width: 100%;
    height: ${props => props.height};
    min-height: ${props => props.minHeight};    
`;

const ControlContainer = Styled.div`
  display: flex;
  flex-direction: row;
  width: 100%;
  justify-content: space-between;
  align-items: center;
  height: 64px;
  text-transform: uppercase;
  color: ${props => props.theme.mint};
  font-weight: 600;
  font-size: 18px;
  padding-left: 15px;    
  .plus-icon{
      margin-right: 10px;
  }
`;

const Left = Styled(Wrapper)`
    width: 30%;
    border-right: 1px solid #372463;
    height: 100%;
`;

const Right = Styled(Wrapper)`
    height: 100%;
    width: 70%;
    flex-direction: column;
`;

const Ribbon = Styled.div`
  background-color: ${props => props.bg || `#372463`};
  color: #45E5B7;
  height: 40px;
  border: 1px solid #372463;
  font-family: "Barlow", sans-serif;
  display: flex;
  font-weight: 600;
  text-transform: uppercase;
  align-items: center;
  padding-left: 10px;
  justify-content: space-between;
`;

const Row = Styled.div`
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: space-between;
  height: 48px;
  border-bottom: 1px solid #372463;
  div {
    padding: 0px 16px;
    text-transform: uppercase;
  }
`;

const DetailWrapper = Styled.div`
  width: 100%;
  max-height: 100%;
  overflow-y: auto;
  border-right: 1px solid #372463;
`;
const DetailInfo = Styled.div`
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  column-gap: 40px;
  padding: 16px;  
`;

const PlateNumber = Styled.div`
  height: 80px;
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 48px;
  text-transform: uppercase;
`;

const ErrorWrapper = Styled.div`
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: ${props => props.theme.inlineError};
  div {
    margin-bottom: 15px;
  }
`;

/* eslint-disable no-console */
/* eslint-disable no-nested-ternary */
import React, { useState, useEffect, useContext } from "react";
import PropTypes from "prop-types";
import Styled, { ThemeContext } from "styled-components";
import PageWrapper from "components/organisms/pageWrapper";
import BlankContainer from "components/atoms/BlankContainer";
import Text from "components/atoms/Text";
import LoadingSpinner from "components/atoms/LoadingSpinner";
import GlobeBg from "assets/images/globe.svg";

import { getListAnalytic } from "api";

function License(props) {
  const { history } = props;
  const themeContext = useContext(ThemeContext);
  const [isLoading, setIsLoading] = useState(true);
  const [blank, setBlank] = useState(false);
  const [errorFetch, setErrorFetch] = useState("");
  const [deploymentKey, setDeploymentKey] = useState("");
  const [data, setData] = useState([]);

  function callData() {
    setIsLoading(true);
    getListAnalytic()
      .then(result => {
        const resultdata = result;
        setData(resultdata.analytics);
        setDeploymentKey(resultdata.deployment_key);
        setIsLoading(false);
      })
      .catch(error => {
        setErrorFetch(error.message);
        setBlank(true);
        setIsLoading(false);
      });
  }

  useEffect(() => {
    callData();
  }, []);

  const renderListLicns = data.map((license, i) => {
    const seats = license.seats.map((key, index) => (
      <li key={index}>{key.serial_number}</li>
    ));

    return (
      <LicenseGroup key={i}>
        <div>
          <h2 className="title">{license.name}</h2>
          <h3 className="title2">Total {license.seats.length}</h3>
        </div>
        <LicenseList>{seats}</LicenseList>
      </LicenseGroup>
    );
  });

  return (
    <PageWrapper title="License" history={history}>
      {blank ? (
        <BlankContainer bg={GlobeBg}>
          <BoxedDiv>
            <Text
              size="18"
              color={themeContext.inlineError}
              textTransform="uppercase"
              weight="600"
            >
              {errorFetch !== "" ? errorFetch : "Lincese Not Available"}
            </Text>
          </BoxedDiv>
        </BlankContainer>
      ) : isLoading ? (
        <BlankContainer>
          <LoadingSpinner show={isLoading} />
        </BlankContainer>
      ) : (
        <Wrapper>
          <WrapperLicense>
            <WrapperHeader>
              <h2 className="title">DEPLOYMENT KEY</h2>
              <h2 className="title">{deploymentKey}</h2>
            </WrapperHeader>
            {renderListLicns}
          </WrapperLicense>
        </Wrapper>
      )}
    </PageWrapper>
  );
}

const Wrapper = Styled.div`
  display: flex;
  overflow-y: scroll;
  width: auto;
  max-height: calc(100% - 64px);
`;
const WrapperLicense = Styled.div`
  width: 100%;
  margin: 20px;
`;
const WrapperHeader = Styled.div`
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  padding: 0 50px;
`;

const LicenseList = Styled.ul`
  margin:0;
  padding:0;
  display: grid;
  grid-template-columns: repeat(5,1fr);
  grid-template-rows: auto;
  list-style-type:none;
  gap:20px;
`;

const LicenseGroup = Styled.div`
  display: grid;
  padding:20px 40px 40px 40px;
  gap:20px;
  border: 6px solid #fff;
  margin-bottom: 40px;
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

License.propTypes = {
  history: PropTypes.object.isRequired,
  location: PropTypes.object,
  value: PropTypes.array
};

export default License;

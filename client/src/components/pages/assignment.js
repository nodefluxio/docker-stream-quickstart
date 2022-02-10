/* eslint-disable no-unused-vars */
/* eslint-disable no-console */
/* eslint-disable no-return-assign */
/* eslint-disable no-case-declarations */
import React, { useEffect, useState } from "react";
import Styled from "styled-components";
import PropTypes from "prop-types";
import PageWrapper from "components/organisms/pageWrapper";

import CameraSetting from "components/organisms/Assignment/CameraSetting";
import AnalyticSetting from "components/organisms/Assignment/AnalyticSetting";
import ROISetting from "components/organisms/Assignment/ROISetting";

export default function Assignment(props) {
  const { history } = props;
  const stepSchema = [
    {
      step: "Select Camera",
      button_text: "NEXT"
    },
    {
      step: "Select Analaytic",
      button_text: "NEXT"
    },
    {
      step: "Assignment Configuration",
      button_text: "FINISH"
    }
  ];

  const STEP_CAMERA = 0;
  const STEP_ANALYTIC = 1;
  const STEP_ROI = 2;

  const [selectedCam, setSelectedCam] = useState("");
  const [selectedAnalytic, setSelectedAnalytic] = useState("");
  const [stepCounter, setStepCounter] = useState(STEP_CAMERA);
  const [JSONConfig, setJSONConfig] = useState({});
  const [isEditMode, setIsEditMode] = useState(false);
  const [selectedNode, setSelectedNode] = useState(0);

  useEffect(() => {
    const { search } = history.location;
    if (search === "") {
      history.replace(`${history.location.pathname}?step=1`);
      setStepCounter(STEP_CAMERA);
    } else {
      const query = new URLSearchParams(search);
      const step = query.get("step");
      const selectedCamera = query.get("selected_camera");
      if (selectedCamera !== null && Number(step) > 1) {
        setStepCounter(STEP_ANALYTIC);
        setSelectedCam(selectedCamera);
      } else if (selectedCam !== "" && selectedAnalytic !== "") {
        setStepCounter(STEP_ROI);
      } else if (selectedCam !== "" && selectedAnalytic === "") {
        setStepCounter(STEP_ANALYTIC);
      } else {
        setStepCounter(STEP_CAMERA);
      }
    }
  }, [history.location.search]);

  const handleJSONChange = event => {
    if (!event.error && event.jsObject === undefined) {
      setJSONConfig({});
    } else {
      setJSONConfig(event.jsObject);
    }
  };

  const renderStep = () => {
    switch (stepCounter) {
      case STEP_CAMERA:
        return <CameraSetting history={history} />;

      case STEP_ANALYTIC:
        return (
          <AnalyticSetting
            history={history}
            selectedCam={selectedCam}
            handleJSONChange={e => handleJSONChange(e)}
            handleSelectedAnalytic={value => setSelectedAnalytic(value)}
            editModeCallback={value => setIsEditMode(value)}
          />
        );

      case STEP_ROI:
        return (
          <ROISetting
            history={history}
            selectedAnalytic={selectedAnalytic}
            selectedCam={selectedCam}
            JSONConfig={JSONConfig}
            isEditMode={isEditMode}
          />
        );

      default:
        return null;
    }
  };

  return (
    <PageWrapper title="Analytic Configuration" history={history} close={true}>
      <ContentWrapper>
        {stepCounter < stepSchema.length && (
          <Content step={stepCounter}>{renderStep()}</Content>
        )}
      </ContentWrapper>
    </PageWrapper>
  );
}

Assignment.propTypes = {
  history: PropTypes.object.isRequired
};

const ContentWrapper = Styled.div`
  width: 100%;
  height: calc(100vh - 64px);
  display: flex;
  flex-direction: row;
`;

const Content = Styled.div`
  margin-left: ${props => (props.step === 2 ? "10%" : "25%")};
  width:${props => (props.step === 2 ? "80%" : "50%")};
  max-height: 90%;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  border: 1px solid ${props => props.theme.secondary2};
`;

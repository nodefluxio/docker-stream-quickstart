/* eslint-disable no-console */
/* eslint-disable no-return-assign */
/* eslint-disable prefer-destructuring */
/* eslint-disable no-case-declarations */
import React, { useEffect, useState, useCallback } from "react";
import Styled from "styled-components";
import PropTypes from "prop-types";
import PageWrapper from "components/organisms/pageWrapper";
import Accordion from "components/molecules/Accordion";
import Input, { EditableInput } from "components/molecules/Input";
import Modal from "components/molecules/Modal";
import Button from "components/molecules/Button";
import JSONInput from "react-json-editor-ajrm";

import Roi from "components/organisms/Roi";

import { stopVisualisation } from "assets/js/visualstreamer";

import { ANALYTIC_METADATA } from "constants/analyticMetadata";
import { REACT_APP_API_CAMERA } from "config";

import {
  getListCamera,
  createCamera,
  getListAnalytic,
  createStream
} from "api";

export default function Assignment(props) {
  const { history } = props;
  const formSchemaCamera = {
    stream_name: "",
    stream_address: ""
  };
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
  const directionList = [
    { value: "in", label: "IN" },
    { value: "in-out", label: "IN-OUT" }
  ];

  const STEP_CAMERA = 0;
  const STEP_ANALYTIC = 1;
  const STEP_ROI = 2;

  const [formCamera, setFormCamera] = useState(formSchemaCamera);
  const [cameraList, setCameraList] = useState([]);
  const [selectedCam, setSelectedCam] = useState("");
  const [analyticList, SetAnalyticList] = useState([]);
  const [selectedAnalytic, setSelectedAnalytic] = useState("");
  const [stepCounter, setStepCounter] = useState(STEP_CAMERA);
  const [btnDisable, setBtnDisable] = useState(true);
  const [loading, setLoading] = useState(false);
  const [line, setLine] = useState(0);
  const [lineData, setLinedata] = useState([]);
  const [areaIsEdited, setAreaIsEdited] = useState([]);
  const [reverse, setReverse] = useState({});
  const [reverseVal, setReverseVal] = useState([]);
  const [isConfirmOpen, setIsConfirmOpen] = useState(false);
  const [notifData, setNotifData] = useState({
    type: "",
    message: "",
    subMessage: ""
  });
  const [JSONConfig, setJSONConfig] = useState({});

  useEffect(() => {
    getListCamera()
      .then(result => {
        const resultdata = result.streams;
        const listCam = [];
        resultdata.forEach(val => {
          listCam.push({
            value: val.stream_id,
            label: val.stream_name
          });
        });
        setCameraList(value => [...value, ...listCam]);
      })
      .catch(error => {
        console.log(error.message);
      });

    getListAnalytic()
      .then(result => {
        const resultdata = result.analytics;
        const listAnalytic = [];
        resultdata.forEach(val => {
          listAnalytic.push({
            value: val.id,
            label: val.name
          });
        });
        SetAnalyticList(value => [...value, ...listAnalytic]);
      })
      .catch(error => {
        console.log(error.message);
      });
    return () => stopVisualisation();
  }, []);

  const formatLabel = string => {
    const splitString = string.split("_");
    const keys =
      splitString[1].charAt(0).toUpperCase() + splitString[1].slice(1);
    return `Camera ${keys}`;
  };

  const handleChangeFormCamera = event => {
    const { name, value } = event.target;
    setBtnDisable(true);
    setFormCamera(prevState => ({
      ...prevState,
      [name]: value
    }));
    setSelectedCam("");
  };

  useEffect(() => {
    if (formCamera.stream_address !== "" && formCamera.stream_name !== "") {
      setBtnDisable(false);
    }
  }, [formCamera]);

  const handleChangeSelectExistCam = event => {
    setSelectedCam(event.target.value);
    setFormCamera(formSchemaCamera);
    setBtnDisable(false);
  };

  const resetRoi = () => {
    setLinedata([]);
  };

  const deleteLine = lineNumber => {
    setLinedata(lineData.filter(item => item.lineNumber !== lineNumber));
    setLine(lineNumber);
  };

  const changeDirectionVal = (e, index) => {
    const updatedData = [...lineData];
    updatedData[index].direction = e.target.value;
    setLinedata(updatedData);
  };

  const updateRoi = (point, areaName, color, lineNumber) => {
    if (point !== null) {
      const nData = {
        point,
        areaName,
        direction: "in",
        color,
        lineNumber
      };
      setLinedata(value => [...value, nData]);
      setAreaIsEdited(value => [...value, false]);
      setReverseVal(value => [...value, "end"]);
    }
  };

  const submitCamera = () => {
    const form = formCamera;
    form.stream_latitude = parseFloat(0);
    form.stream_longitude = parseFloat(0);
    setLoading(true);
    createCamera(form)
      .then(result => {
        if (result.code === 200) {
          setLoading(false);
          setSelectedCam(result.stream_id);
        } else {
          setStepCounter(STEP_CAMERA);
          setLoading(false);
        }
      })
      .catch(err => {
        console.log(err);
        setStepCounter(STEP_CAMERA);
        setLoading(false);
      });
  };

  const showNotif = () => (
    <Modal
      type="feedback"
      isOpen={isConfirmOpen}
      duration={5000}
      manualClose={() => {
        setIsConfirmOpen(false);
        if (notifData.type === "success") {
          history.push("/");
        }
      }}
      desc={notifData.message}
      subDesc={notifData.subMessage}
      status={notifData.type}
    />
  );

  const addAssignment = () => {
    let pipelineConf = {};
    if (lineData.length !== 0) {
      const areas = [];
      lineData.forEach(value => {
        areas.push({
          name: value.areaName,
          points: value.point,
          bidirection: value.direction !== "in"
        });
      });
      pipelineConf.areas = areas;
    }
    if (JSONConfig !== {}) {
      pipelineConf = { ...pipelineConf, ...JSONConfig };
    }
    createStream(selectedCam, 0, selectedAnalytic, pipelineConf)
      .then(results => {
        if (results.code === 200) {
          setIsConfirmOpen(true);
          setNotifData({
            type: "success",
            message: "succesfully assign analytic",
            subMessage: ""
          });
          setLoading(false);
        } else {
          setIsConfirmOpen(true);
          setNotifData({
            type: "error",
            message: "failed assign analytic",
            subMessage: ""
          });
          setLoading(false);
        }
      })
      .catch(error => {
        const { response } = error;
        setIsConfirmOpen(true);
        setNotifData({
          type: "error",
          message: "failed assign analytic",
          subMessage: response.data.message
        });
        setLoading(false);
      });
  };

  const handleButtonStep = () => {
    if (selectedCam === "" && stepCounter === STEP_CAMERA) {
      submitCamera();
    }
    if (stepCounter === STEP_ROI) {
      addAssignment();
    }
    if (stepCounter !== STEP_ROI) {
      setStepCounter(value => value + 1);
      if (stepCounter + 1 !== STEP_ROI) {
        setBtnDisable(true);
      }
    }
  };

  const handleAreaIsEdit = index => {
    const updatedData = [...areaIsEdited];
    updatedData[index] = !areaIsEdited[index];
    setAreaIsEdited(updatedData);
  };

  const changeAreaNameVal = useCallback(event => {
    const { name, value } = event.target;
    const updatedData = [...lineData];
    updatedData[name].areaName = value;
    setLinedata(updatedData);
  });

  const reverseLine = (index, lineNumber, direction) => {
    const nDirection = direction === "end" ? "start" : "end";
    const nReverseVal = [...reverseVal];
    const nLineData = [...lineData];
    const tmpPoint = lineData[index].point[0];

    nReverseVal[index] = nDirection;
    nLineData[index].point[0] = lineData[index].point[1];
    nLineData[index].point[1] = tmpPoint;

    setReverse({
      line: lineNumber,
      direction: nDirection
    });
    setReverseVal(nReverseVal);
    setLinedata(nLineData);
  };

  const handleJSONChange = event => {
    if (!event.error && event.jsObject === undefined) {
      setJSONConfig({});
    } else {
      setJSONConfig(event.jsObject);
    }
  };

  useEffect(() => {
    if (
      selectedAnalytic !== "" &&
      (JSONConfig === {} || JSONConfig !== undefined)
    ) {
      setBtnDisable(false);
    } else if (JSONConfig === undefined) {
      setBtnDisable(true);
    }
  }, [selectedAnalytic, JSONConfig]);

  const renderStep = () => {
    switch (stepCounter) {
      case STEP_CAMERA:
        return (
          <Wrapper>
            <Accordion header="Add New Camera">
              {Object.keys(formCamera).map(key => (
                <EditableInput
                  key={key}
                  label={formatLabel(key)}
                  id={key}
                  onChange={handleChangeFormCamera}
                  name={key}
                  value={formCamera[key]}
                ></EditableInput>
              ))}
            </Accordion>
            <Accordion header="Select Available Camera" type={"closed"}>
              <Input
                key="select-camera"
                id="select-camera"
                type="select"
                option={cameraList}
                placeholder="Select Camera"
                onChange={handleChangeSelectExistCam}
                value={selectedCam}
              />
            </Accordion>
          </Wrapper>
        );

      case STEP_ANALYTIC:
        return (
          <Wrapper>
            <Accordion header="Select Analytic Type">
              <Input
                key="select-analytic"
                id="select-analytic"
                label="Analytic Name"
                type="select"
                option={analyticList}
                placeholder="Select Analytic"
                onChange={e => setSelectedAnalytic(e.target.value)}
                value={selectedAnalytic}
              />
            </Accordion>
            <Accordion header="Custom Configuration">
              <JSONInput
                id="custom-configuration"
                height="250px"
                width="100%"
                style={{
                  contentBox: {
                    fontSize: "14px"
                  }
                }}
                onChange={e => handleJSONChange(e)}
              />
            </Accordion>
          </Wrapper>
        );

      case STEP_ROI:
        const config = ANALYTIC_METADATA[selectedAnalytic];
        return (
          <Wrapper>
            <Accordion
              header={`Assignment Configuration ( ${config.analytic_name} )`}
            >
              <WrapperROI>
                <Roi
                  onSave={updateRoi}
                  onReset={resetRoi}
                  roiType={config.roi_type}
                  isRoiNeed={config.roi}
                  image={`${REACT_APP_API_CAMERA}/stream_jpeg/0/${selectedCam}`}
                  delLine={line}
                  currPoints={{}}
                  reverse={reverse}
                />
                {config.roi && (
                  <RoiConfig>
                    <h1>{config.roi_title}</h1>
                    <Rule>
                      {config.roi_rule.map(val => (
                        <li key={val}>{val}</li>
                      ))}
                    </Rule>
                    <h4 style={{ margin: "60px 0 10px 0" }}>Counting Lines</h4>
                    <Row>
                      <Title>Name</Title>
                      <Title>Direction</Title>
                      <Title>Action</Title>
                      {lineData.length > 0 &&
                        lineData.map((value, index) => {
                          if (typeof value === "object") {
                            return (
                              <>
                                <div style={{ display: "flex" }}>
                                  <BulletPoint color={value.color} />
                                  {areaIsEdited[index] ? (
                                    <Input
                                      type="text"
                                      placeholder="Area Name"
                                      onChange={changeAreaNameVal}
                                      value={value.areaName}
                                      name={index}
                                      key={`area-${index}`}
                                      id={`area-${index}`}
                                    />
                                  ) : (
                                    <span>{value.areaName}</span>
                                  )}
                                  <ActionArea
                                    onClick={() => handleAreaIsEdit(index)}
                                  >
                                    {areaIsEdited[index] ? (
                                      "ok"
                                    ) : (
                                      <div className="icon">
                                        <svg
                                          width="9"
                                          height="9"
                                          viewBox="0 0 9 9"
                                          fill="none"
                                        >
                                          <path
                                            d="M8.70711 1.70711C9.09763 1.31658 9.09763 0.683418 8.70711 0.292893C8.31658 -0.0976309 7.68342 -0.097631 7.29289 0.292893L8.70711 1.70711ZM7.29289 0.292893L6.29289 1.29289L7.70711 2.70711L8.70711 1.70711L7.29289 0.292893Z"
                                            fill="#F7F7F7"
                                          />
                                          <path
                                            d="M0.292893 7.29289C-0.0976314 7.68342 -0.097632 8.31658 0.292892 8.70711C0.683417 9.09763 1.31658 9.09763 1.70711 8.70711L0.292893 7.29289ZM5.29289 2.29289L0.292893 7.29289L1.70711 8.70711L6.70711 3.70711L5.29289 2.29289Z"
                                            fill="#F7F7F7"
                                          />
                                        </svg>
                                      </div>
                                    )}
                                  </ActionArea>
                                </div>
                                <div style={{ width: "110px" }}>
                                  <Input
                                    type="select"
                                    option={directionList}
                                    placeholder="Direction"
                                    onChange={e => changeDirectionVal(e, index)}
                                    value={value.direction}
                                    key={`direction-${index}`}
                                    id={`direction-${index}`}
                                  />
                                </div>
                                <div style={{ display: "flex" }}>
                                  <ActionLine
                                    onClick={() =>
                                      reverseLine(
                                        index,
                                        value.lineNumber,
                                        reverseVal[index]
                                      )
                                    }
                                  >
                                    Reverse
                                  </ActionLine>
                                  <span style={{ margin: "0 10px 0 10px" }}>
                                    |
                                  </span>
                                  <ActionLine
                                    onClick={() => deleteLine(value.lineNumber)}
                                  >
                                    Delete
                                  </ActionLine>
                                </div>
                              </>
                            );
                          }
                          return null;
                        })}
                    </Row>
                  </RoiConfig>
                )}
              </WrapperROI>
            </Accordion>
          </Wrapper>
        );

      default:
        return null;
    }
  };

  return (
    <PageWrapper title="Analytic Assignment" history={history} close={true}>
      {isConfirmOpen && showNotif()}
      <ContentWrapper>
        {stepCounter < stepSchema.length && (
          <Content step={stepCounter}>
            {renderStep()}
            <Button
              isLoading={loading}
              disabled={btnDisable}
              style={{ margin: "15px" }}
              onClick={() => handleButtonStep(stepCounter)}
            >
              {stepSchema[stepCounter].button_text}
            </Button>
          </Content>
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

const WrapperROI = Styled.div`
  display: flex;
  flex-direction: row;
`;

const RoiConfig = Styled.div`
  margin-left: 30px
`;

const Wrapper = Styled.div`
  display: flex;
  flex-direction: column;
`;

const Row = Styled.div`
  display: grid;
  grid-template-columns:  1fr 1fr 1fr;
  align-items: center;
  height: 40px;
  justify-items: center;
  row-gap: 5px;
`;

const BulletPoint = Styled.div`
  background-color: ${props => props.color};
  height: 10px;
  width: 10px;
  margin: auto 20px auto 0px;
`;

const Title = Styled.div`
  display: block;
  font-size: 0.83em;
  margin-block-start: 1em;
  margin-block-end: 1em;
  margin-inline-start: 0px;
  margin-inline-end: 0px;
  font-weight: bold;
`;

const ActionLine = Styled.div`
  cursor: pointer;
`;

const ActionArea = Styled.div`
  margin: auto 0px auto 15px;
  cursor: pointer;
`;

const Rule = Styled.ol`
  list-style: none;
  counter-reset: item;
  padding-left: 5px;

  li {
   counter-increment: item;
   margin-bottom: 5px;
 }

  li:before {
    content: counter(item) ". ";
    font-weight: bold;
  }

`;

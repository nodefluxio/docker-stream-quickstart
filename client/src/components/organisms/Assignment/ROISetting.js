/* eslint-disable prefer-destructuring */
import React, { Fragment, useState, useEffect, useCallback } from "react";
import PropTypes from "prop-types";
import Styled from "styled-components";
import { useDispatch } from "react-redux";
import Accordion from "components/molecules/Accordion";
import Input from "components/molecules/Input";
import Button from "components/molecules/Button";
import Roi from "components/organisms/Roi";
import { createStream, deleteAnalytic } from "api";

import { showGeneralNotification } from "store/actions/notification";

import { ANALYTIC_METADATA } from "constants/analyticMetadata";
import { REACT_APP_API_CAMERA } from "config";

export default function ROISetting(props) {
  const {
    history,
    selectedAnalytic,
    selectedCam,
    JSONConfig,
    isEditMode
  } = props;
  const dispatch = useDispatch();
  const config = ANALYTIC_METADATA[selectedAnalytic];

  const directionList = [
    { value: "in", label: "IN" },
    { value: "in-out", label: "IN-OUT" }
  ];

  const [line, setLine] = useState(0);
  const [lineData, setLinedata] = useState([]);
  const [areaIsEdited, setAreaIsEdited] = useState([]);
  const [reverse, setReverse] = useState({});
  const [reverseVal, setReverseVal] = useState([]);
  const [loading, setLoading] = useState(false);
  const [btnDisable, setBtnDisable] = useState(false);
  const [selectedNode, setSelectedNode] = useState("");

  useEffect(() => {
    const { search } = history.location;
    let node = selectedNode;
    if (search !== "") {
      const query = new URLSearchParams(search);
      node = query.get("selected_node");
      setSelectedNode(node);
    }
  }, [history.location]);

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

  const callCreateStream = pipelineConf => {
    createStream(selectedCam, selectedNode, selectedAnalytic, pipelineConf)
      .then(results => {
        if (results.code === 200) {
          dispatch(
            showGeneralNotification({
              type: "success",
              desc: "succesfully assign analytic"
            })
          );
          setBtnDisable(false);
          setLoading(false);
          history.push("/camera");
        } else {
          dispatch(
            showGeneralNotification({
              type: "error",
              desc: "failed assign analytic"
            })
          );
          setBtnDisable(false);
          setLoading(false);
        }
      })
      .catch(error => {
        const { response } = error;
        dispatch(
          showGeneralNotification({
            type: "error",
            desc: "failed assign analytic",
            subDesc: response.data.message
          })
        );
        setBtnDisable(false);
        setLoading(false);
      });
  };

  const addAssignment = () => {
    let pipelineConf = {};
    setBtnDisable(true);
    setLoading(true);
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
    if (isEditMode) {
      deleteAnalytic(selectedCam, selectedAnalytic, selectedNode).then(
        result => {
          if (result.code === 200) {
            callCreateStream(pipelineConf);
          } else {
            setBtnDisable(false);
            setLoading(false);
            dispatch(
              showGeneralNotification({
                type: "error",
                desc: "Error Editing Analytic"
              })
            );
          }
        }
      );
    } else {
      callCreateStream(pipelineConf);
    }
  };

  return (
    <Fragment>
      <Wrapper>
        <Accordion
          open={true}
          header={`Assignment Configuration ( ${config.analytic_name} )`}
        >
          <WrapperROI>
            {selectedNode !== "" && (
              <Roi
                onSave={updateRoi}
                onReset={resetRoi}
                roiType={config.roi_type}
                isRoiNeed={config.roi}
                image={`${REACT_APP_API_CAMERA}/stream_jpeg/${selectedNode}/${selectedCam}`}
                delLine={line}
                currPoints={{}}
                reverse={reverse}
              />
            )}
            {config.roi && (
              <RoiConfig>
                <h1>{config.roi_title}</h1>
                <Rule>
                  {config.roi_rule.map(val => (
                    <li key={val}>{val}</li>
                  ))}
                </Rule>
                <h4 style={{ margin: "60px 0 10px 0" }}>
                  {config.roi_type === "COUNTER_LINE"
                    ? "Counting Lines"
                    : "Area of Interest"}
                </h4>
                <Row column={config.roi_type === "COUNTER_LINE" ? 3 : 2}>
                  <Title>Name</Title>
                  {config.roi_type === "COUNTER_LINE" && (
                    <Title>Direction</Title>
                  )}
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
                            {config.roi_type === "COUNTER_LINE" && (
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
                            )}
                            <div style={{ display: "flex" }}>
                              {config.roi_type === "COUNTER_LINE" && (
                                <Fragment>
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
                                </Fragment>
                              )}
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

      <Button
        isLoading={loading}
        disabled={btnDisable}
        style={{ margin: "15px" }}
        onClick={() => addAssignment()}
      >
        NEXT
      </Button>
    </Fragment>
  );
}

ROISetting.propTypes = {
  history: PropTypes.object.isRequired,
  selectedAnalytic: PropTypes.string,
  selectedCam: PropTypes.string,
  JSONConfig: PropTypes.object,
  isEditMode: PropTypes.bool,
  selectedNode: PropTypes.number
};

const Wrapper = Styled.div`
  display: flex;
  flex-direction: column;
`;

const WrapperROI = Styled.div`
  display: flex;
  flex-direction: row;
`;

const RoiConfig = Styled.div`
  margin-left: 30px
`;

const Row = Styled.div`
  display: grid;
  grid-template-columns: repeat(${props => props.column},  1fr);
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

/* eslint-disable no-unused-vars */
import React, { useState, useEffect, Fragment, useContext } from "react";
import PropTypes from "prop-types";
import { useDispatch } from "react-redux";
import Styled, { ThemeContext } from "styled-components";
import Accordion from "components/molecules/Accordion";
import Input from "components/molecules/Input";
import Button from "components/molecules/Button";
import JSONInput from "react-json-editor-ajrm";
import Modal from "components/molecules/Modal";
import { showGeneralNotification } from "store/actions/notification";

import { deleteAnalytic, getCamera, getListAnalytic } from "api";

export default function AnalyticSetting(props) {
  const {
    history,
    handleJSONChange,
    handleSelectedAnalytic,
    selectedCam,
    editModeCallback
  } = props;
  const dispatch = useDispatch();
  const themeContext = useContext(ThemeContext);
  const [analyticList, setAnalyticList] = useState([]);
  const [selectedAnalytic, setSelectedAnalytic] = useState("");
  const [ownedAnalytic, setOwnedAnalytic] = useState([]);
  const [isConfirmOpen, setIsConfirmOpen] = useState(false);
  const [selectedNode, setSelectedNode] = useState("");

  useEffect(() => {
    const { search } = history.location;
    if (search !== "") {
      const query = new URLSearchParams(search);
      const node = query.get("selected_node");
      setSelectedNode(node);

      getCamera(selectedCam, null, node)
        .then(result => {
          if (result.ok) {
            setOwnedAnalytic(result.stream.pipelines || []);
          }
        })
        .catch(() => setOwnedAnalytic([]));
    }
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
        setAnalyticList(value => [...value, ...listAnalytic]);
      })
      .catch(error => {
        // eslint-disable-next-line no-console
        console.log(error.message);
      });
  }, [history.location]);

  function isEditMode(analyticID) {
    const index = ownedAnalytic.findIndex(analytic => analytic === analyticID);
    const editMode = index > -1;
    editModeCallback(editMode);
    return editMode;
  }

  function handleAnalytic(value) {
    setSelectedAnalytic(value);
    handleSelectedAnalytic(value);
  }

  function callDeleteAnalytic() {
    deleteAnalytic(selectedCam, selectedAnalytic, selectedNode).then(result => {
      if (result.code === 200) {
        dispatch(
          showGeneralNotification({
            type: "success",
            desc: result.message
          })
        );
        history.push("/camera");
      } else {
        dispatch(
          showGeneralNotification({
            type: "error",
            desc: "Error Deleteing Analytic"
          })
        );
      }
    });
  }

  return (
    <Fragment>
      <Wrapper>
        <Accordion header="Select Analytic Type" open={true}>
          <Input
            key="select-analytic"
            id="select-analytic"
            label="Analytic Name"
            type="select"
            option={analyticList}
            placeholder="Select Analytic"
            onChange={e => handleAnalytic(e.target.value)}
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
      <ButtonWrapper>
        {isEditMode(selectedAnalytic) && (
          <Button
            disabled={selectedAnalytic === ""}
            style={{ margin: "15px" }}
            width="50%"
            type="secondary"
            onClick={() => setIsConfirmOpen(true)}
          >
            REMOVE
          </Button>
        )}
        <Button
          width={isEditMode(selectedAnalytic) ? "50%" : "100%"}
          disabled={selectedAnalytic === ""}
          style={{ margin: "15px" }}
          onClick={() =>
            history.replace(
              `${history.location.pathname}?step=3&selected_node=${selectedNode}`
            )
          }
        >
          NEXT
        </Button>
      </ButtonWrapper>
      <Modal
        type="confirmation"
        isShown={isConfirmOpen}
        onConfirm={() => {
          callDeleteAnalytic(selectedAnalytic);
          setIsConfirmOpen(false);
        }}
        onClose={() => setIsConfirmOpen(false)}
        title="Delete Analytic"
        buttonTitle="Proceed Delete"
        header="ARE YOU SURE TO PROCEED
        DELETING THIS ANALYTIC?"
        headerColor={themeContext.inlineError}
      />
    </Fragment>
  );
}

AnalyticSetting.propTypes = {
  history: PropTypes.object.isRequired,
  selectedCam: PropTypes.string,
  handleJSONChange: PropTypes.func.isRequired,
  handleSelectedAnalytic: PropTypes.func.isRequired,
  editModeCallback: PropTypes.func.isRequired
};

AnalyticSetting.defaultProps = {
  selectedCam: ""
};

const Wrapper = Styled.div`
  display: flex;
  flex-direction: column;
`;

const ButtonWrapper = Styled.div`
  display: flex;
  flex-direction: row;
  width: 100%;
`;

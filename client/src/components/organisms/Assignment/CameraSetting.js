import React, { useState, useEffect, Fragment } from "react";
import PropTypes from "prop-types";
import { useDispatch } from "react-redux";
import Styled from "styled-components";
import Accordion from "components/molecules/Accordion";
import Input, { EditableInput } from "components/molecules/Input";
import Button from "components/molecules/Button";
import {
  getListCamera,
  createCamera,
  getSites,
  addSites,
  assignSite,
  getCamera,
  updateCamera,
  getListNode
} from "api";
import { showGeneralNotification } from "store/actions/notification";

export default function CameraSetting(props) {
  const { history, nodeCb } = props;
  const formSchemaCamera = {
    stream_name: "",
    stream_address: "",
    stream_site_id: ""
  };
  const dispatch = useDispatch();
  const [formCamera, setFormCamera] = useState(formSchemaCamera);
  const [previousForm, setPreviousForm] = useState(formSchemaCamera);
  const [cameraList, setCameraList] = useState([]);
  const [mode, setMode] = useState("create");
  const [pageMode, setPageMode] = useState("create");
  const [selectedCam, setSelectedCam] = useState("");
  const [loading, setLoading] = useState(false);
  const [btnDisable, setBtnDisable] = useState(true);
  const [sites, setSites] = useState([]);
  const [isSiteLoading, setIsSiteLoading] = useState(false);
  const [selectedNode, setSelectedNode] = useState("");
  const [nodeList, setNodeList] = useState([]);

  useEffect(() => {
    getSites().then(result => {
      if (result.ok) {
        const newSite = [];
        result.sites.map(site =>
          newSite.push({ value: site.id, label: site.name })
        );
        setSites(newSite);
      }
    });
    getListNode()
      .then(result => {
        const resultdata = result.nodes;
        const newList = resultdata.map(val => ({
          value: val.node_num,
          label: `node ${val.node_num} => ${val.status}`,
          isDisabled: val.status === "offline"
        }));
        setNodeList([...newList]);
      })
      .catch(error => {
        console.log(error.message);
      });
  }, []);

  useEffect(() => {
    if (mode === "create" && pageMode === "edit") {
      if (JSON.stringify(formCamera) !== JSON.stringify(previousForm)) {
        setBtnDisable(false);
      } else {
        setBtnDisable(true);
      }
    } else if (mode === "create" && pageMode === "create") {
      if (
        formCamera.stream_address !== "" &&
        formCamera.stream_name !== "" &&
        formCamera.stream_site_id !== ""
      ) {
        setBtnDisable(false);
      } else {
        setBtnDisable(true);
      }
    } else if (mode === "select") {
      if (selectedCam !== "") {
        setBtnDisable(false);
      } else {
        setBtnDisable(true);
      }
    }
  }, [formCamera, selectedCam, pageMode, mode, previousForm]);

  useEffect(() => {
    const { search } = history.location;
    if (search !== "") {
      const query = new URLSearchParams(search);
      const selectedCamera = query.get("selected_camera");
      const node = query.get("selected_node");
      setSelectedNode(node);
      if (selectedCamera !== null) {
        getCamera(selectedCamera, null, node).then(result => {
          if (result.ok) {
            const { stream } = result;
            setPageMode("edit");
            setFormCamera({
              stream_name: stream.stream_name,
              stream_address: stream.stream_address,
              stream_site_id: stream.stream_site_id || ""
            });
            setPreviousForm({
              stream_name: stream.stream_name,
              stream_address: stream.stream_address,
              stream_site_id: stream.stream_site_id || ""
            });
          }
        });
      } else {
        setPageMode("create");
        getListCamera()
          .then(result => {
            const resultdata = result.data.streams;
            const listCam = [];
            resultdata.forEach(val => {
              listCam.push({
                value: val.stream_id,
                label: val.stream_name,
                node: val.stream_node_num
              });
            });
            setCameraList(value => [...value, ...listCam]);
          })
          .catch(error => {
            // eslint-disable-next-line no-console
            console.log(error.message);
          });
      }
    }
  }, [history.location.search]);

  const handleChangeSelectExistCam = event => {
    const value = event.target.value
    const camera =  cameraList.find(cam => cam.value === value);
    
    setSelectedNode(camera.node)
    setSelectedCam(value);
    setFormCamera(formSchemaCamera);
    setBtnDisable(false);
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

  const formatLabel = string => {
    const splitString = string.split("_");
    const keys =
      splitString[1].charAt(0).toUpperCase() + splitString[1].slice(1);
    return `Camera ${keys}`;
  };

  const submitCamera = () => {
    setLoading(true);
    setBtnDisable(true);
    const body = { ...formCamera };
    delete body.stream_site_id;
    if (pageMode === "create") {
      if (mode === "create") {
        createCamera(
          {
            ...body,
            stream_latitude: parseFloat(0),
            stream_longitude: parseFloat(0)
          },
          selectedNode
        )
          .then(result => {
            if (result.code === 200) {
              assignSite(
                { stream_id: result.stream_id },
                formCamera.stream_site_id
              )
                .then(resultSite => {
                  if (resultSite.ok) {
                    dispatch(
                      showGeneralNotification({
                        type: "success",
                        desc: "successfully assign site"
                      })
                    );
                  }
                })
                .catch(() => {
                  dispatch(
                    showGeneralNotification({
                      type: "error",
                      desc: "failed assign site"
                    })
                  );
                });
              setBtnDisable(false);
              setLoading(false);
              history.push(
                `${history.location.pathname}?step=2&selected_camera=${result.stream_id}&selected_node=${selectedNode}`
              );
            } else {
              setLoading(false);
              setBtnDisable(false);
            }
          })
          .catch(err => {
            // eslint-disable-next-line no-console
            console.log(err);
            setLoading(false);
            setBtnDisable(false);
          });
      } else if (mode === "select" && selectedCam !== "") {
        history.push(
          `${history.location.pathname}?step=2&selected_camera=${selectedCam}&selected_node=${selectedNode}`
        );
      }
    } else if (pageMode === "edit") {
      const query = new URLSearchParams(history.location.search);
      const selectedCamera = query.get("selected_camera");
      if (formCamera.stream_site_id !== previousForm.stream_site_id) {
        assignSite({ stream_id: selectedCamera }, formCamera.stream_site_id)
          .then(result => {
            if (result.ok) {
              dispatch(
                showGeneralNotification({
                  type: "success",
                  desc: "successfully assign site"
                })
              );
              history.push("/camera");
            }
          })
          .catch(() => {
            setLoading(false);
            setBtnDisable(false);
            dispatch(
              showGeneralNotification({
                type: "error",
                desc: "failed assign site"
              })
            );
          });
      }
      if (
        formCamera.stream_name !== previousForm.stream_name ||
        formCamera.stream_address !== previousForm.stream_address
      ) {
        const node = query.get("selected_node");
        updateCamera(body, selectedCamera, node)
          .then(result => {
            if (Object.keys(result).length > 0) {
              dispatch(
                showGeneralNotification({
                  type: "success",
                  desc: "successfully edit camera"
                })
              );
              history.push("/camera");
            }
          })
          .catch(() => {
            setLoading(false);
            setBtnDisable(false);
            dispatch(
              showGeneralNotification({
                type: "error",
                desc: "failed edit camera"
              })
            );
          });
      }
    }
  };

  const handleCreateSite = value => {
    setIsSiteLoading(true);
    addSites({
      name: value
    })
      .then(result => {
        if (result.ok) {
          setSites([
            ...sites,
            {
              value: result.site.id,
              label: result.site.name
            }
          ]);
        }
        setFormCamera({ ...formCamera, stream_site_id: result.site.id });
        setIsSiteLoading(false);
      })
      .catch(() => setIsSiteLoading(false));
  };

  const handleChangeSelectNode = event => {
    setSelectedNode(event.target.value);
    nodeCb(event.target.value);
  };

  return (
    <Fragment>
      <Wrapper>
        <Accordion
          header="Add New Camera"
          open={mode === "create"}
          callback={() => setMode("create")}
        >
          {Object.keys(formCamera).map(key =>
            key === "stream_site_id" ? (
              <Input
                key={key}
                label="Sites:"
                type="select"
                option={sites}
                placeholder="Sites"
                onChange={e =>
                  setFormCamera({
                    ...formCamera,
                    stream_site_id: e.target.value
                  })
                }
                value={formCamera[key]}
                isCreatable
                onCreateOption={value => handleCreateSite(value)}
                isLoading={isSiteLoading}
                isDisabled={isSiteLoading}
              />
            ) : (
              <EditableInput
                key={key}
                label={formatLabel(key)}
                id={key}
                onChange={handleChangeFormCamera}
                name={key}
                value={formCamera[key]}
              ></EditableInput>
            )
          )}
          {pageMode !== "edit" && (
            <Input
              label="Select Node"
              key="select-node"
              id="select-node"
              type="select"
              option={nodeList}
              placeholder="Select Node"
              onChange={handleChangeSelectNode}
              value={selectedNode}
            />
          )}
        </Accordion>
        {pageMode !== "edit" && (
          <Accordion
            header="Select Available Camera"
            open={mode === "select"}
            callback={() => setMode("select")}
          >
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
        )}
      </Wrapper>
      <Button
        isLoading={loading}
        disabled={btnDisable}
        style={{ margin: "15px" }}
        onClick={() => submitCamera()}
      >
        {pageMode === "edit" ? "SAVE" : "NEXT"}
      </Button>
    </Fragment>
  );
}

CameraSetting.propTypes = {
  history: PropTypes.object.isRequired,
  nodeCb: PropTypes.func
};

const Wrapper = Styled.div`
  display: flex;
  flex-direction: column;
`;

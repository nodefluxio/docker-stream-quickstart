import React, { useState, useEffect } from "react";
import Styled from "styled-components";
import PropTypes from "prop-types";
import Dropzone from "react-dropzone";
import csv from "csv";
import Modal from "components/atoms/Modal";
import { EditableInput } from "components/molecules/Input";
import Button from "components/molecules/Button";
import DownloadIcon from "assets/icon/visionaire/download.svg";

import { register } from "api/vehicle";

function BatchAddVehicle(props) {
  const { openModal, onClose } = props;
  const title = "Batch";
  const [valueCSV, setValueCSV] = useState("");
  const [listEnrollment, setListEnrollment] = useState([]);
  const [log, setLog] = useState([]);
  const [showLog, setShowLog] = useState(false);
  const [processedImg, setProcessedImg] = useState(0);
  const [processedImgSuccess, setProcessedImgSuccess] = useState(0);
  const [processStart, setProcessStart] = useState(null);
  const [processButton, setProcessButton] = useState(false);
  const [refreshFlag, setRefreshFlag] = useState(false);
  const [isFinished, setIsFinished] = useState(false);

  function handleOnDropCSV(acceptedFiles) {
    if (acceptedFiles.length > 0) {
      const reader = new FileReader();
      reader.onload = () => {
        // Parse CSV file
        csv.parse(reader.result, (err, dataCSV) => {
          dataCSV.shift();
          setListEnrollment(dataCSV);
        });
      };

      // read file contents
      acceptedFiles.forEach(file => reader.readAsBinaryString(file));
      setValueCSV(acceptedFiles[0].path);
      setProcessButton(false);
    }
  }

  function proccessBatchEnroll() {
    setProcessStart(performance.now());
    setLog([]);
    setShowLog(true);
    setProcessButton(true);
    setProcessedImgSuccess(0);
    setProcessedImg(0);

    const message = {
      message: "Batch enrollment start...",
      type: "start"
    };
    setLog(value => [...value, message]);

    listEnrollment.forEach(x => {
      const newForm = {
        plate_number: x[0],
        unique_id: x[1],
        name: x[2],
        type: x[3],
        brand: x[4],
        color: x[5],
        status: x[6]
      };
      register(newForm)
        .then(result => {
          if (result.ok) {
            const messageSuccess = {
              message: `${newForm.plate_number}@${newForm.name} -- Enrollment is success`,
              type: "success"
            };
            setLog(value => [...value, messageSuccess]);
            setProcessedImgSuccess(value => value + 1);
            setProcessedImg(value => value + 1);
          } else {
            const messageFailed = {
              message: `${newForm.plate_number}@${newForm.name} -- Enrollment is failed with error -- ${result.message}`,
              type: "error"
            };
            setLog(value => [...value, messageFailed]);
            setProcessedImg(value => value + 1);
          }
        })
        .catch(err => {
          const response = err.response.data;
          const errorMsg = response.errors[0] || response.message;
          const messageFailed = {
            message: `${newForm.plate_number}@${newForm.name} -- Enrollment is failed with error -- ${errorMsg}`,
            type: "error"
          };
          setLog(value => [...value, messageFailed]);
          setProcessedImg(value => value + 1);
        });
    });
  }

  function downloadLog() {
    const element = document.createElement("a");
    let logs = "";
    log.forEach(logData => {
      logs += `${logData.message} \n`;
    });
    const file = new Blob([logs], { type: "text/plain" });
    element.href = URL.createObjectURL(file);
    element.download = "batch_enrollmetn_log.log";
    document.body.appendChild(element);
    element.click();
  }

  useEffect(() => {
    if (processedImg === listEnrollment.length && processedImg !== 0) {
      const endTime = performance.now();
      const secondDivider = 1000;
      const timeDiff = (endTime - processStart) / secondDivider; // in second
      // get seconds
      const seconds = Math.round(timeDiff);
      const message = {
        message: `Batch enrollment finished, (${processedImgSuccess}/${processedImg} are success) -- Time Elapsed ${seconds} seconds`,
        type: "end"
      };
      setIsFinished(true);
      setRefreshFlag(true);
      setProcessButton(false);
      setLog(value => [...value, message]);
    }
  }, [processedImg]);

  useEffect(() => {
    if (valueCSV === "") {
      setProcessButton(true);
    }
  }, [valueCSV]);

  function closeModal(flag) {
    setValueCSV("");
    setLog([]);
    setIsFinished(false);
    onClose(flag);
  }

  return (
    <Modal
      show={openModal}
      className="modal-add-enrollment"
      title={`${title} Enrollment`}
      close={() => closeModal(refreshFlag)}
      padding="20px"
      width="800px"
    >
      <Dropzone onDrop={handleOnDropCSV} accept=".csv" multiple={false}>
        {({ getRootProps, getInputProps }) => (
          <section>
            <div {...getRootProps()}>
              <input {...getInputProps()} />
              <EditableInput
                label="CSV File"
                id="csv-file"
                name="csv-file"
                readOnly={true}
                value={valueCSV}
              ></EditableInput>
            </div>
          </section>
        )}
      </Dropzone>
      <LinkButton
        key="need-help"
        id="need-help"
        onClick={() =>
          window.open(
            "https://docs.nodeflux.io/visionaire-docker-stream/features",
            "_blank"
          )
        }
      >
        <span>need help ?</span>
      </LinkButton>{" "}
      <LinkButton
        key="template-csv"
        id="template-csv"
        onClick={() => window.open(`/templates/batch-vehicle.csv`, "_blank")}
      >
        <span>download template</span>
      </LinkButton>
      {showLog && (
        <>
          <ContainerDownload>
            <LinkButton key="download" id="download" onClick={downloadLog}>
              <img src={DownloadIcon} alt="download_log" />
              <span>download log</span>
            </LinkButton>
          </ContainerDownload>
          <LogContainer>
            <ul>
              {log.length > 0 &&
                log.map((logData, index) => (
                  <Logger key={index} type={logData.type}>
                    {logData.message}
                  </Logger>
                ))}
            </ul>
          </LogContainer>
        </>
      )}
      <Button
        disabled={processButton}
        width="100%"
        style={{ marginTop: "40px" }}
        onClick={() =>
          isFinished ? closeModal(refreshFlag) : proccessBatchEnroll()
        }
        type="primary"
      >
        {isFinished ? "Finished" : "Proceed"}
      </Button>
    </Modal>
  );
}

const ContainerDownload = Styled.div`
  float: right;
`;
const LinkButton = Styled.button`
  border: none;
  background: none;
  outline: none;
  cursor: pointer;
  display: flex;
  align-items: center;
  margin-top: 15px;
  
  img,svg{
    max-width: 30px;
  }
  span{
    font-family: Barlow;
    font-style: normal;
    font-weight: 500;
    font-size: 12px;
    line-height: 14px;
    text-transform: uppercase;
    text-decoration: underline;
    color: ${props => props.theme.mercury};
  }
`;

const LogContainer = Styled.div`
  position: relative;
  margin-top: 20px;
  width: 100%;
  min-height: 300px;
  max-height: 500px;
  overflow-y: auto;
  overflow-x: hidden;
  display: flex;
  flex-direction: row;
  flex-wrap: wrap;
  background: #000;
`;

const Logger = Styled.li`
  ${props => {
    let color = "";
    switch (props.type) {
      case "error":
        color = "red";
        break;
      case "success":
        color = "green";
        break;

      default:
        color = "#fff";
        break;
    }
    return `color: ${color}`;
  }}
`;

BatchAddVehicle.propTypes = {
  openModal: PropTypes.bool.isRequired,
  onClose: PropTypes.func.isRequired,
  data: PropTypes.object.isRequired
};

export default BatchAddVehicle;

import React, { useState, useEffect } from "react";
import Styled from "styled-components";
import PropTypes from "prop-types";
import Dropzone from "react-dropzone";
import csv from "csv";
import Modal from "components/atoms/Modal";
import { EditableInput } from "components/molecules/Input";
import Button from "components/molecules/Button";
import DownloadIcon from "assets/icon/visionaire/download.svg";

import { createEnrollment } from "api";

function BatchEnrollment(props) {
  const { openModal, onClose } = props;
  const title = "Batch";
  const [valueCSV, setValueCSV] = useState("");
  const [valueFolder, setValueFolder] = useState("");
  const [listEnrollment, setListEnrollment] = useState([]);
  const [images, setImages] = useState([]);
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

  function handleDropFolder(acceptedFiles) {
    setImages(acceptedFiles);
    const path = acceptedFiles[0].path.replace(/^\//, "");
    const newPath = path.split("/")[0];
    setValueFolder(newPath);
  }

  function proccessBatchEnroll() {
    setProcessStart(performance.now());
    setLog([]);
    setShowLog(true);
    setProcessButton(true);
    setProcessedImgSuccess(0);
    setProcessedImg(0);

    const dataImages = [];
    if (images.length > 0) {
      for (let i = 0; i < images.length; i += 1) {
        const blob = new Blob([images[i]], { type: "image/jpeg" });
        const path = images[i].path.replace(/^\//, "").split("/");
        dataImages.push({
          file_path: images[i].path,
          file_name: images[i].name,
          folder_name: path[1],
          blob
        });
      }
    }

    const message = {
      message: "Batch enrollment start...",
      type: "start"
    };
    setLog(value => [...value, message]);

    listEnrollment.forEach(x => {
      const newForm = new FormData();
      const identityNumber = x[0];
      const name = x[1];
      const status = x[2];
      const folderName = x[3];

      newForm.append("identity_number", identityNumber);
      newForm.append("name", name);
      newForm.append("status", status);

      const imageField = dataImages.filter(j => j.folder_name === folderName);
      imageField.map(image => newForm.append("images", image.blob));

      createEnrollment(newForm)
        .then(result => {
          if (result.ok) {
            const messageSuccess = {
              message: `${identityNumber}@${name} -- Enrollment is success`,
              type: "success"
            };
            setLog(value => [...value, messageSuccess]);
            setProcessedImgSuccess(value => value + 1);
            setProcessedImg(value => value + 1);
          } else {
            const messageFailed = {
              message: `${identityNumber}@${name} -- Enrollment is failed with error -- ${result.message}`,
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
            message: `${identityNumber}@${name} -- Enrollment is failed with error -- ${errorMsg}`,
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
    setValueFolder("");
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
      <Dropzone onDrop={handleDropFolder}>
        {({ getRootProps, getInputProps }) => (
          <section>
            <div {...getRootProps()}>
              <input {...getInputProps()} webkitdirectory="" />
              <EditableInput
                label="Folder"
                id="folder"
                name="folder"
                readOnly={true}
                value={valueFolder}
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

BatchEnrollment.propTypes = {
  openModal: PropTypes.bool.isRequired,
  onClose: PropTypes.func.isRequired,
  data: PropTypes.array.isRequired
};

export default BatchEnrollment;

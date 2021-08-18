import React, { useState, useCallback, useEffect, useRef } from "react";
import Styled from "styled-components";
import PropTypes from "prop-types";
import Dropzone from "react-dropzone";
import Modal from "components/atoms/Modal";
import { EditableInput } from "components/molecules/Input";
import Button from "components/molecules/Button";
import IconPlaceholder from "components/atoms/IconPlaceholder";

import AddIcon from "assets/icon/visionaire/plus-circle.svg";
import ExitIcon from "assets/icon/visionaire/exit.svg";
import { createEnrollment, updateEnrollment } from "api";
import {
  REACT_APP_API_ENROLLMENT,
  REACT_APP_MAX_SIZE_IMG_ENROLLMENT
} from "config";

function AddEnrollment(props) {
  const { openModal, onClose, data } = props;
  const formSchema = {
    identity_number: "",
    name: "",
    status: "",
    images: []
  };
  const [formData, setFormData] = useState(formSchema);
  const [thumbnail, setThumbnail] = useState([]);
  const [title, setTitle] = useState("Add");
  const [errorMsg, setErrorMsg] = useState("");
  const dropzoneRef = useRef(null);
  const [isLoading, setIsLoading] = useState(false);

  const handleChange = useCallback(event => {
    const { name, value } = event.target;
    setFormData(prevState => ({
      ...prevState,
      [name]: value
    }));
  });

  const formatLabel = string =>
    string.replace(/_/g, " ").replace(/\b(\w)/g, str => str.toUpperCase());

  function getFile(acceptedFiles) {
    if (acceptedFiles.length > 0) {
      const newThumbnails = [];
      const newImages = formData.images;
      for (let i = 0; i < acceptedFiles.length; i += 1) {
        const blob = new Blob([acceptedFiles[i]], { type: "image/jpeg" });
        newImages.push(blob);
        const blobURL = URL.createObjectURL(blob);
        newThumbnails.push({ image: blobURL });
      }
      setFormData({ ...formData, images: newImages });
      setThumbnail(thumbnail.concat(newThumbnails));
    }
  }

  function saveEnrollment() {
    setIsLoading(true);
    const newForm = new FormData();
    formData.images.map(image => newForm.append("images", image));
    newForm.append("identity_number", formData.identity_number);
    newForm.append("name", formData.name);
    newForm.append("status", formData.status);
    if (title === "Add") {
      createEnrollment(newForm)
        .then(result => {
          if (result.ok) {
            setErrorMsg("");
            onClose(result.enrollment);
            setIsLoading(false);
          }
        })
        .catch(error => {
          setErrorMsg(error.response.data.message);
          setIsLoading(false);
        });
    } else if (title === "Update") {
      updateEnrollment(data.id, newForm)
        .then(result => {
          if (result.ok) {
            setIsLoading(false);
            setErrorMsg("");
            onClose(result.enrollment);
          }
        })
        .catch(error => {
          setIsLoading(false);
          setErrorMsg(error.response.data.message);
        });
    }
  }

  useEffect(() => {
    if (openModal === false) {
      setFormData(formSchema);
      setTitle("Add");
      setThumbnail([]);
    }
  }, [openModal]);

  useEffect(() => {
    if (Object.keys(data).length !== 0 && data.constructor === Object) {
      setTitle("Update");
      setThumbnail(data.faces);
      setFormData({
        ...formData,
        name: data.name,
        identity_number: data.identity_number,
        status: data.status
      });
    }
  }, [data]);

  function deleteThumbnail(formId, thumbnailId) {
    const newImageArray = formData.images;
    const newThumbnailArray = thumbnail;
    newImageArray.splice(formId, 1);
    newThumbnailArray.splice(thumbnailId, 1);
    setFormData({ ...formData, images: newImageArray });
    setThumbnail(newThumbnailArray);
  }

  return (
    <Modal
      show={openModal}
      className="modal-add-enrollment"
      title={`${title} Enrollment`}
      close={onClose}
      padding="20px"
      width={thumbnail.length === 0 ? "400px" : "800px"}
    >
      <FormRow>
        {Object.keys(formData).map(
          key =>
            key !== "images" && (
              <EditableInput
                key={key}
                label={formatLabel(key)}
                id={key}
                onChange={handleChange}
                name={key}
                value={formData[key]}
              ></EditableInput>
            )
        )}
      </FormRow>
      <ImgContainer length={thumbnail.length}>
        {thumbnail.length > 0 && (
          <ImgWrapper border={true} onClick={() => dropzoneRef.current.open()}>
            <img src={AddIcon} alt="add" />
            Upload Image
          </ImgWrapper>
        )}
        {thumbnail.map((file, index) => (
          <ImgWrapper key={index}>
            {formData.images.length > 0 &&
              index >= thumbnail.length - formData.images.length && (
                <BtnClose>
                  <IconPlaceholder
                    onClick={() =>
                      deleteThumbnail(
                        index - (thumbnail.length - formData.images.length),
                        index
                      )
                    }
                  >
                    <img src={ExitIcon} alt="close-icon" />
                  </IconPlaceholder>
                </BtnClose>
              )}
            <img
              key={index}
              src={
                file.id
                  ? `${REACT_APP_API_ENROLLMENT}/face/image/${file.id}`
                  : `${file.image}`
              }
              alt={`preview-thumbnail-${index}`}
            />
          </ImgWrapper>
        ))}
        <Dropzone
          ref={dropzoneRef}
          multiple={true}
          accept="image/jpeg"
          maxSize={REACT_APP_MAX_SIZE_IMG_ENROLLMENT}
          onDrop={getFile}
        >
          {({ getRootProps, getInputProps, isDragReject, fileRejections }) => (
            <DropFiles
              isDragReject={isDragReject}
              fileRejections={fileRejections}
              {...getRootProps()}
            >
              <input {...getInputProps()} />
              {thumbnail.length === 0 && (
                <BtnAdd>
                  <img src={AddIcon} alt="add" className="add-icon" />
                  Upload Image
                </BtnAdd>
              )}

              {fileRejections.length > 0 && (
                <div
                  style={{
                    color: "red",
                    marginTop: "80px",
                    textAlign: "center"
                  }}
                >
                  Unsupported file type or file size too large. Maximum file
                  allowed is {REACT_APP_MAX_SIZE_IMG_ENROLLMENT / 1000} Kb
                </div>
              )}
            </DropFiles>
          )}
        </Dropzone>
      </ImgContainer>
      {errorMsg !== "" && (
        <div
          style={{
            color: "red",
            textAlign: "left"
          }}
        >
          {errorMsg}
        </div>
      )}
      <Button
        width="100%"
        style={{ marginTop: "40px" }}
        onClick={() => saveEnrollment()}
        isLoading={isLoading}
        disabled={isLoading}
      >
        Save
      </Button>
    </Modal>
  );
}

const FormRow = Styled.div`
  display: flex;
  flex-direction: column;
  .input-group{
    margin-bottom: 20px;
    input {
      margin-bottom: 0px;
    }
  }
`;

const ImgContainer = Styled.div`
  position: relative;
  margin-top: 20px;
  width: 100%;
  ${props => props.length === 0 && `min-height: 300px;`}
  max-height: 500px;
  overflow-y: auto;
  overflow-x: hidden;
  display: flex;
  flex-direction: row;
  flex-wrap: wrap;
`;

const ImgWrapper = Styled.div`
  position: relative;
  width: 17%;
  padding: 10px;
  text-align: center;
  ${props =>
    props.border &&
    `
  border: 1px solid white;
  cursor: pointer;
  `}
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  img {
    max-width: 100%;
  }
  .add-icon {
    max-width: 40px;
  }
`;

const DropFiles = Styled.div`
  width: 100%;
  display: flex;
  justify-content: center;
  flex-direction: column;
  align-items: center;
  cursor: pointer;
  position: absolute;
  outline: none;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
`;

const BtnClose = Styled.div`
  position: absolute;
  right: 0;
  top: 0;
  cursor: pointer;
  width: 40px;
  height: 40px;
  font-size: 24px;
  z-index: 20;
`;

const BtnAdd = Styled.div`
  margin: auto;
  min-width: 150px;
  text-align: center;
  color: white;
  font-size: 12px;
  line-height: 14px;
  font-weight: bold;
  cursor: pointer;
  img{
    width: 20px !important;
    display: block;
    margin: auto;
    margin-bottom: 10px;
  }
`;

AddEnrollment.propTypes = {
  openModal: PropTypes.bool.isRequired,
  onClose: PropTypes.func.isRequired,
  data: PropTypes.array.isRequired
};

export default AddEnrollment;

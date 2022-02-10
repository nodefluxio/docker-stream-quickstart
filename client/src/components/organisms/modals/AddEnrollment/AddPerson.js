import React, { useState, useCallback, useEffect, useRef } from "react";
import Styled from "styled-components";
import PropTypes from "prop-types";
import Dropzone from "react-dropzone";
import SVG from "react-inlinesvg";
import Modal from "components/atoms/Modal";
import Input, { EditableInput } from "components/molecules/Input";
import Button from "components/molecules/Button";
import IconPlaceholder from "components/atoms/IconPlaceholder";

import b64toBlob from "helpers/b64toblob";
import { getDateFormat } from "helpers/dateTime";

import AddIcon from "assets/icon/visionaire/plus-circle.svg";
import DeleteIcon from "assets/icon/visionaire/delete.svg";

import { createEnrollment, updateEnrollment } from "api";
import {
  REACT_APP_API_ENROLLMENT,
  REACT_APP_MAX_SIZE_IMG_ENROLLMENT
} from "config";
import theme from "theme";

function AddPerson(props) {
  const { openModal, onClose, data } = props;
  const formSchema = {
    identity_number: "",
    name: "",
    gender: "",
    birth_place: "",
    birth_date: "",
    status: "",
    images: []
  };
  const genderOption = [
    {
      label: "Pria",
      value: "pria"
    },
    {
      label: "Wanita",
      value: "wanita"
    }
  ];
  const [formData, setFormData] = useState(formSchema);
  const [thumbnail, setThumbnail] = useState([]);
  const [title, setTitle] = useState("Add");
  const [errorMsg, setErrorMsg] = useState("");
  const dropzoneRef = useRef(null);
  const [isLoading, setIsLoading] = useState(false);
  const [dataToDelete, setDataToDelete] = useState([]);

  const handleChange = useCallback(event => {
    const { name, value } = event.target;
    setFormData(prevState => ({
      ...prevState,
      [name]: value
    }));
  });

  const handleChangeDate = (name, value) => {
    setFormData(prevState => ({
      ...prevState,
      [name]: getDateFormat(value)
    }));
  };

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
    newForm.append("gender", formData.gender);
    newForm.append("birth_date", formData.birth_date);
    newForm.append("birth_place", formData.birth_place);
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
          setErrorMsg(error.response.data.errors[0]);
          setIsLoading(false);
        });
    } else if (title === "Update") {
      if (dataToDelete.length > 0) {
        dataToDelete.map(item => newForm.append("deleted_variations", item));
      }
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
      setDataToDelete([]);
    }
  }, [openModal]);

  useEffect(() => {
    if (Object.keys(data).length > 1 && data.constructor === Object) {
      setTitle("Update");
      setThumbnail(data.faces);
      setFormData({
        ...formData,
        name: data.name,
        identity_number: data.identity_number,
        status: data.status,
        gender: data.gender,
        birth_place: data.birth_place,
        birth_date: data.birth_date
      });
    } else if (typeof data === "string") {
      const newImages = formData.images;
      const blob = b64toBlob(data);
      const blobURL = URL.createObjectURL(blob);
      newImages.push(blob);
      setFormData({ ...formData, images: newImages });
      setThumbnail([{ image: blobURL }]);
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

  function deleteImage(thumbnailId) {
    const { variation } = thumbnail[thumbnailId];
    setDataToDelete([...dataToDelete, variation]);
    const newThumbnailArray = thumbnail;
    newThumbnailArray.splice(thumbnailId, 1);
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
        {Object.keys(formData).map(key => {
          if (key === "gender") {
            return (
              <Input
                label={formatLabel(key)}
                key={key}
                id={key}
                type="select"
                option={genderOption}
                placeholder="Select Gender"
                onChange={handleChange}
                name={key}
                value={formData[key]}
              />
            );
          }
          if (key === "birth_date") {
            return (
              <Input
                type="date"
                label="Birth Date"
                key={key}
                time={false}
                name="birth_date"
                onSubmit={(name, value) => handleChangeDate(name, value)}
                value={formData[key]}
                shouldCloseOnSelect={true}
              />
            );
          }
          if (key !== "images") {
            return (
              <EditableInput
                key={key}
                label={formatLabel(key)}
                id={key}
                onChange={handleChange}
                name={key}
                value={formData[key]}
              ></EditableInput>
            );
          }
          return key;
        })}
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
            <BtnClose>
              <IconPlaceholder
                onClick={() =>
                  index >= thumbnail.length - formData.images.length
                    ? deleteThumbnail(
                        index - (thumbnail.length - formData.images.length),
                        index
                      )
                    : deleteImage(index)
                }
              >
                <StyledSVG src={DeleteIcon} fill={theme.inlineError} />
              </IconPlaceholder>
            </BtnClose>
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
        disabled={isLoading || thumbnail.length === 0}
      >
        Save
      </Button>
    </Modal>
  );
}

const StyledSVG = Styled(SVG)`
  margin-right: 5px;
  & path {
    fill: ${({ fill }) => fill};
  }
  & rect {
    fill: ${({ fill }) => fill};
  }
`;

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

AddPerson.propTypes = {
  openModal: PropTypes.bool.isRequired,
  onClose: PropTypes.func.isRequired,
  data: PropTypes.oneOfType([
    PropTypes.string,
    PropTypes.array,
    PropTypes.object
  ])
};

AddPerson.defaultProps = {
  data: []
};

export default AddPerson;

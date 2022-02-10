/* eslint-disable array-callback-return */
/* eslint-disable no-unused-vars */
import React, { useState, useEffect, useContext, useRef } from "react";
import Styled, { ThemeContext } from "styled-components";
import PropTypes from "prop-types";
import Modal from "components/atoms/Modal";
import Input, {
  EditableInput,
  PasswordToggled
} from "components/molecules/Input";
import Text from "components/atoms/Text";
import Button from "components/molecules/Button";
import { getSites, addUser, editUser } from "api";
import validateInput from "helpers/validator";
import { useDispatch } from "react-redux";
import { showGeneralNotification } from "store/actions/notification";
import usePrevious from "helpers/usePrevious";

function AddUser(props) {
  const dispatch = useDispatch();
  const themeContext = useContext(ThemeContext);
  const { openModal, onClose, data } = props;
  const formSchema = {
    fullname: "",
    email: "",
    username: "",
    role: "",
    site_id: ""
  };
  const formObjectText = [
    {
      name: "fullname",
      value: "",
      type: "text"
    },
    {
      name: "email",
      value: "",
      type: "email"
    },
    {
      name: "username",
      value: "",
      type: "text"
    }
  ];
  const [formData, setFormData] = useState(formSchema);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [sites, setSites] = useState([]);
  const [selectedSites, setSelectedSites] = useState([]);
  const [mode, setMode] = useState("add");
  const [formObject, setFormObject] = useState(formObjectText);
  const [btnDisable, setBtnDisable] = useState(true);
  const prevData = usePrevious({ data });

  const role = [
    {
      label: "Superadmin",
      value: "superadmin"
    },
    {
      label: "Operator",
      value: "operator"
    }
  ];

  useEffect(() => {
    const keys = Object.keys(data);
    const hasPassword = formObject.some(val => val.name === "password");
    if (keys.length > 0) {
      let newFormData = formData;
      keys.map(key => {
        const index = Object.keys(formData).findIndex(
          form => String(key) === String(form)
        );
        if (index > -1) {
          newFormData = { ...newFormData, [key]: data[key] };
        }
      });
      if (data.site_id) {
        const newSites = [];
        for (let i = 0; i < sites.length; i += 1) {
          const index = data.site_id.findIndex(
            value => parseInt(value, 10) === parseInt(sites[i].value, 10)
          );
          newSites.push({
            label: sites[i].label,
            value: sites[i].value,
            checked: index !== -1
          });
        }
        setSites([...newSites]);
      }
      setFormData(newFormData);
      setSelectedSites(data.site_id || []);
      setMode("edit");
      setFormObject(formObjectText);
      delete newFormData.password;
      delete newFormData.re_password;
      setFormData(newFormData);
    } else if (keys.length === 0 && prevData !== undefined && !hasPassword) {
      setFormObject([
        ...formObject,
        {
          name: "password",
          value: "",
          type: "password"
        },
        {
          name: "re_password",
          value: "",
          type: "password"
        }
      ]);
      setFormData({ ...formData, password: "", re_password: "" });
      const newSites = [];
      for (let i = 0; i < sites.length; i += 1) {
        newSites.push({
          label: sites[i].label,
          value: sites[i].value,
          checked: false
        });
      }
      setSites([...newSites]);
    }
  }, [data]);

  const formatLabel = string => {
    const keys = string.charAt(0).toUpperCase() + string.slice(1);
    return keys;
  };

  const handleChange = event => {
    const { name, value } = event.target;
    let type = "";
    const index = formObject.findIndex(form => form.name === name);
    if (index > -1) {
      ({ type } = formObject[index]);
    }
    setFormData(prevState => ({
      ...prevState,
      [name]: value
    }));
    setTimeout(() => {
      if (type === "email" || type === "text") {
        if (validateInput(type, value)) {
          setError("");
        } else {
          setError(`Please input correct ${name}`);
        }
      } else if (type === "password") {
        const comparator = name === "password" ? "re_password" : "password";
        if (value.length < 8) {
          setError("Password must be 8 characters or more");
        } else if (value !== formData[comparator]) {
          setError("Confirmation Password didn't match with Password");
        } else if (value === formData[comparator]) {
          setError("");
        }
      }
    }, 1);
  };

  const handleCheck = value => {
    const siteIndex = sites.findIndex(
      site => String(site.value) === String(value)
    );
    // copy value array, not reference
    let newSelectedSite = selectedSites.slice();
    const newSite = sites;
    if (value === "all") {
      const flag = sites.length > 0 && sites.length === selectedSites.length;
      if (!flag) {
        for (let i = 0; i < sites.length; i += 1) {
          newSelectedSite.push(sites[i].value);
        }
      } else {
        newSelectedSite = [];
      }
      for (let i = 0; i < newSite.length; i += 1) {
        newSite[i].checked = !flag;
      }
    } else {
      if (sites[siteIndex].checked === false) {
        newSelectedSite.push(parseInt(value, 10));
      } else {
        newSelectedSite.splice(siteIndex, 1);
      }
      newSite[siteIndex].checked = !newSite[siteIndex].checked;
    }
    setSelectedSites([...newSelectedSite]);
    setSites([...newSite]);
  };

  const submitUser = () => {
    const body = formData;
    setBtnDisable(true);
    setLoading(true);
    if (mode === "add") {
      addUser(body)
        .then(result => {
          if (result.ok) {
            dispatch(
              showGeneralNotification({
                type: "success",
                desc: "succesfully create user"
              })
            );
            onClose(true);
            setBtnDisable(false);
            setLoading(false);
          }
        })
        .catch(() => {
          dispatch(
            showGeneralNotification({
              type: "error",
              desc: "failed create user"
            })
          );
          setBtnDisable(false);
          setLoading(false);
        });
    } else if (mode === "edit") {
      editUser(data.id, body)
        .then(result => {
          if (result.ok) {
            dispatch(
              showGeneralNotification({
                type: "success",
                desc: "succesfully change user"
              })
            );
            onClose(true);
            setBtnDisable(false);
            setLoading(false);
          }
        })
        .catch(() => {
          dispatch(
            showGeneralNotification({
              type: "error",
              desc: "failed change user"
            })
          );
          setBtnDisable(false);
          setLoading(false);
        });
    }
  };

  useEffect(() => {
    if (selectedSites.length > 0) {
      const index = selectedSites.findIndex(value => value === "all");
      const tobesubmitted = selectedSites;
      if (index > -1) {
        tobesubmitted.splice(index, 1);
      }
      setFormData({ ...formData, site_id: tobesubmitted });
    }
  }, [selectedSites]);

  useEffect(() => {
    let sameAsPreviousFlag = false;

    if (formData.role === "superadmin") {
      delete formData.site_id;
    }
    if (Object.keys(data).length > 0) {
      let flag = true;
      // eslint-disable-next-line consistent-return
      Object.keys(data).map(key => {
        if (formData[key] === undefined) {
          return false;
        }
        if (data[key] !== formData[key]) flag = false;
      });
      sameAsPreviousFlag = flag;
    }
    const values = Object.values(formData);
    const nullIndex = values.findIndex(value => value === "");
    if (nullIndex === -1 && error === "" && !sameAsPreviousFlag) {
      setBtnDisable(false);
    } else {
      setBtnDisable(true);
    }
  }, [formData, error]);

  useEffect(() => {
    if (!openModal) {
      setFormData(formSchema);
      setSelectedSites([]);
      setFormObject(formObjectText);
      setMode("add");
      setError("");
    }
  }, [openModal]);

  useEffect(() => {
    getSites().then(result => {
      if (result.ok) {
        if (result.sites.length > 0) {
          const newSites = [];
          for (let i = 0; i < result.sites.length; i += 1) {
            newSites.push({
              label: result.sites[i].name,
              value: result.sites[i].id,
              checked: false
            });
          }
          setSites(newSites);
        }
      }
    });
  }, []);

  const renderForm = (name, type) => {
    switch (type) {
      case "password":
        return (
          <PasswordToggled
            key={name}
            id={name}
            label={
              name === "re_password" ? "Confirm Password" : formatLabel(name)
            }
            onChange={handleChange}
            name={name}
            value={formData[name]}
          />
        );
      default:
        return (
          <EditableInput
            readOnly={
              mode === "edit" && (name === "email" || name === "username")
            }
            key={name}
            label={formatLabel(name)}
            id={name}
            onChange={handleChange}
            name={name}
            value={formData[name]}
          ></EditableInput>
        );
    }
  };

  return (
    <Modal
      show={openModal}
      className="modal-add-user"
      title={`${formatLabel(mode)} User`}
      close={() => onClose(false)}
      padding="20px"
      width="800px"
    >
      <FormWrapper>
        <FormRow>
          {formObject.map(form => renderForm(form.name, form.type))}
        </FormRow>
        <FormRow>
          <Input
            label="Role:"
            name="role"
            type="select"
            option={role}
            placeholder="Role"
            onChange={handleChange}
            value={formData.role}
            style={{ marginBottom: "20px" }}
          />
          {sites.length > 0 && formData.role !== "superadmin" && (
            <CheckboxWrapper mode={mode}>
              <p>Sites: </p>
              <Input
                key="all"
                type="checkbox"
                text="Select All"
                value="all"
                checked={
                  sites.length > 0 && sites.length === selectedSites.length
                }
                onChange={value => handleCheck(value)}
              />
              {sites.map(site => (
                <Input
                  key={site.value}
                  type="checkbox"
                  text={site.label}
                  value={site.value}
                  checked={site.checked}
                  onChange={value => handleCheck(value)}
                />
              ))}
            </CheckboxWrapper>
          )}
        </FormRow>
      </FormWrapper>

      {error && (
        <Text color={themeContext.inlineError} style={{ margin: "15px 0px" }}>
          {error}
        </Text>
      )}
      <Button
        width="100%"
        style={{ marginTop: "20px" }}
        onClick={() => submitUser()}
        isLoading={loading}
        disabled={btnDisable}
      >
        Save
      </Button>
    </Modal>
  );
}

const CheckboxWrapper = Styled.div`
  max-height: ${props => (props.mode === "edit" ? `231px` : `385px`)};
  overflow-y: auto;
  p {
    width: 100%;
    font-size: 14px;
    font-weight: 500;
    font-style: normal;
    font-stretch: normal;
    line-height: normal;
    letter-spacing: normal;
    display: block;
    position: relative;
    color: ${props => props.theme.mint};
    margin-top: 10px;
    margin-bottom: 10px;
  }
  label:first-of-type {
    border-bottom: 1px solid white;
    padding-bottom: 15px;
  }
`;

const FormWrapper = Styled.div`
    width: 100%;
    display: flex;
    flex-direction: row;
    justify-content: space-between;
`;

const FormRow = Styled.div`
  display: flex;
  flex-direction: column;
  width: 48%;
  .input-group{
    margin-bottom: 20px;
    input {
      margin-bottom: 0px;
    }
  }
`;

AddUser.propTypes = {
  openModal: PropTypes.bool.isRequired,
  onClose: PropTypes.func.isRequired,
  data: PropTypes.object
};

AddUser.defaultProps = {
  data: {}
};

export default AddUser;

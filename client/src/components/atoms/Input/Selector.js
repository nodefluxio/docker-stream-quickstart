import React from "react";
import PropTypes from "prop-types";
import styled from "styled-components";
import SelectInput from "react-select";

import theme from "theme";

const setflagColor = color => ({
  alignItems: "center",
  display: "flex",
  width: "100%",
  ":after": {
    backgroundColor: color,
    content: '" "',
    position: "absolute",
    right: 10,
    display: "block",
    marginRight: 8,
    height: 10,
    width: 10
  }
});

const customStyles = {
  menu: provided => ({
    ...provided,
    width: "100%",
    backgroundColor: theme.blueGem,
    marginTop: -5
  }),

  option: (styles, { data }) => {
    const style = styles;
    style.color = theme.mercury;
    style.backgroundColor = theme.blueGem;
    style.height = 40;
    style.paddingLeft = 13;
    style[":hover"] = {
      backgroundColor: theme.secondary1
    };
    if (data.color) {
      return { ...style, ...setflagColor(data.color) };
    }
    return { ...style };
  },

  control: (provided, state) => ({
    ...provided,
    width: "100%",
    backgroundColor: "transparent",
    height: 40,
    borderStyle: "solid",
    borderWidth: 2,
    borderRadius: 8,
    borderColor: state.isFocused ? theme.mint : theme.secondary1,
    boxShadow: "none",
    ":hover": {
      borderColor: theme.mint
    }
  }),

  singleValue: (styles, { data }) => {
    const style = styles;
    style.color = theme.mercury;
    if (data.color) {
      return { ...style, ...setflagColor(data.color) };
    }
    return { ...style };
  },

  input: () => ({
    color: theme.mercury
  })
};

export default function Select({
  option,
  label,
  error,
  name,
  value,
  onChange,
  placeholder
}) {
  const indexInput = option.findIndex(x => x.value === value);
  const handleChange = selectedOption => {
    onChange({
      target: { value: selectedOption.value, name }
    });
  };
  return (
    <WrapSelect>
      {label && <InputLabel>{label}</InputLabel>}
      <SelectInput
        options={option}
        styles={customStyles}
        error={error}
        value={option[indexInput] || ""}
        name={name}
        onChange={handleChange}
        placeholder={placeholder}
      />
      {error && <InputError>{error}</InputError>}
    </WrapSelect>
  );
}

Select.propTypes = {
  option: PropTypes.array,
  label: PropTypes.string,
  style: PropTypes.object,
  error: PropTypes.string,
  onChange: PropTypes.func,
  name: PropTypes.string,
  value: PropTypes.any,
  placeholder: PropTypes.string
};

Select.defaultProps = {
  option: [],
  label: "",
  style: {},
  error: "",
  onChange: () => {},
  name: "",
  value: "",
  placeholder: "Select"
};

const WrapSelect = styled.div`
  display: block;
  width: 100%;
`;

const InputLabel = styled.label`
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
`;

const InputError = styled(InputLabel)`
  font-size: 12px;
  font-weight: 500;
  color: ${props => props.theme.inlineError};
`;

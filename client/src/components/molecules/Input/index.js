import React from "react";
import PropTypes from "prop-types";

import Text from "components/atoms/Input";
import InputPassword from "./InputPassword";
import InputEdit from "./EditableInput";
import DatePicker from "./DatePicker";

class Input extends React.Component {
  static propTypes = {
    onSubmit: PropTypes.func,
    text: PropTypes.string,
    type: PropTypes.string
  };

  static defaultProps = {
    type: "text"
  };

  render() {
    const { type } = this.props;
    switch (type) {
      case "date":
        return <DatePicker {...this.props} />;
      default: {
        return <Text {...this.props} />;
      }
    }
  }
}

export const EditableInput = props => <InputEdit {...props} />;
export const PasswordToggled = props => <InputPassword {...props} />;
export default Input;

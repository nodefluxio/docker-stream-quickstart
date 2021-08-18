import React from "react";
import PropTypes from "prop-types";
import Styled from "styled-components";
import Button from "components/molecules/Button";
import Input from "components/atoms/Input";

import SearchIcon from "assets/icon/visionaire/search.svg";

import Modal from "components/atoms/Modal";

export default function SearchModal(props) {
  const { isShown, onClose, searchInput, searchOnChange, onSearch } = props;
  return (
    <Modal
      show={isShown}
      className="modal-search"
      title="Search Area or Camera"
      close={onClose}
    >
      <SearchModalWrapper>
        <Input
          label="Enter keyword"
          placeholder="Not set..."
          value={searchInput}
          onChange={searchOnChange}
        >
          <img alt="search-icon" src={SearchIcon} />
        </Input>
        <Button width="100%" onClick={onSearch}>
          SEARCH NOW
        </Button>
      </SearchModalWrapper>
    </Modal>
  );
}

SearchModal.propTypes = {
  isShown: PropTypes.bool.isRequired,
  searchInput: PropTypes.string.isRequired,
  searchOnChange: PropTypes.func.isRequired,
  onSearch: PropTypes.func.isRequired,
  onClose: PropTypes.func.isRequired
};

const SearchModalWrapper = Styled.div`
  padding: 16px;
  
  input {
    margin-bottom: 17px;
    padding-right: 25px;
  }

  img {
    display: inline-block;
    margin-left: -20px;
    margin-top: -20px;
  }
`;

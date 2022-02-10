/* eslint-disable no-unused-vars */
import React, { useState, Fragment, useEffect, useContext } from "react";
import Styled, { ThemeContext } from "styled-components";
import PropTypes from "prop-types";
import qs from "qs";
import Button from "components/molecules/Button";
import Dropdown, { Menu as Item } from "components/atoms/Dropdown";
import IconPlaceholder from "components/atoms/IconPlaceholder";
import BlankContainer from "components/atoms/BlankContainer";
import Modal from "components/molecules/Modal";
import Input from "components/atoms/Input";

import PageWrapper from "components/organisms/pageWrapper";
import AddEnrollment from "components/organisms/modals/AddEnrollment";
import BatchAddVehicle from "components/organisms/modals/BatchAddVehicle";
import Pagination from "components/organisms/pagination";
import LoadingSpinner from "components/atoms/LoadingSpinner";
import ListLayout from "components/organisms/ListLayout";
import { getDetail, getList, deleteSingle, deleteAll } from "api/vehicle";

import PlusCircleIcon from "assets/icon/visionaire/plus-circle.svg";
import MoreIcon from "assets/icon/visionaire/more.svg";
import NoRecord from "assets/icon/visionaire/no-record.svg";
import SearchIcon from "assets/icon/visionaire/search.svg";
import ArrowDown from "assets/icon/visionaire/drop-down.svg";
import SinglePerson from "assets/icon/visionaire/Union.svg";
import MultiplePerson from "assets/icon/visionaire/collabolators.svg";
import DeleteIcon from "assets/icon/visionaire/delete.svg";

const LIMIT_DATA = 11;
const dataSchema = [
  { attr: "plate_number", showName: "Plate Number" },
  { attr: "unique_id", showName: "Unique ID" },
  { attr: "name", showName: "Name" },
  { attr: "brand", showName: "Brand" },
  { attr: "type", showName: "Type" },
  { attr: "color", showName: "Color" },
  { attr: "status", showName: "Status" }
];
function Vehicle(props) {
  const { history } = props;
  const [list, setDataList] = useState([]);
  const [openAdd, setOpenAdd] = useState(false);
  const [openBatch, setOpenBatch] = useState(false);
  const [isConfirmOpen, setIsConfirmOpen] = useState(false);
  const [selectedId, setSelectedId] = useState(null);
  const [selectedData, setSelectedData] = useState([]);
  const themeContext = useContext(ThemeContext);
  const [totalPage, setTotalPage] = useState(1);
  const [totalData, setTotalData] = useState(0);
  const [noData, setNoData] = useState(false);
  const [errMsg, setErrMsg] = useState("");
  const [searchValue, setSearchValue] = useState("");
  const [isConfirmAllOpen, setIsConfirmAllOpen] = useState(false);
  const [isLoadingDeleteAll, setIsLoadingDeleteAll] = useState(false);

  const callGetDetail = id => {
    getDetail(id).then(result => {
      if (result.ok) {
        setSelectedData(result.vehicle);
        setOpenAdd(true);
      }
    });
  };

  const callGetList = query => {
    getList(query)
      .then(response => {
        if (response.ok) {
          const { result } = response;
          if (result.total_data > 0) {
            setDataList(result.vehicles);
            setTotalPage(result.total_page);
            setTotalData(result.total_data);
            setNoData(false);
          } else {
            setErrMsg("");
            setNoData(true);
            setDataList([]);
            setTotalData(0);
            setTotalPage(0);
          }
        } else {
          setErrMsg("");
          setNoData(true);
          setDataList([]);
          setTotalData(0);
          setTotalPage(0);
        }
      })
      .catch(err => {
        setErrMsg(err.message);
        setNoData(true);
        setDataList([]);
        setTotalData(0);
        setTotalPage(0);
      });
  };

  const callDelete = id => {
    if (id !== undefined) {
      deleteSingle(id).then(result => {
        if (result.ok) {
          setIsConfirmOpen(false);
          setSelectedId(null);
          const query = history.location.search.replace("?", "");
          callGetList(query);
        }
      });
    } else if (id === undefined) {
      setIsLoadingDeleteAll(true);
      deleteAll()
        .then(result => {
          if (result.ok) {
            setIsLoadingDeleteAll(false);
            setIsConfirmAllOpen(false);
            const query = history.location.search.replace("?", "");
            callGetList(query);
          }
        })
        .catch(() => setIsLoadingDeleteAll(false));
    }
  };

  const elmDropdown = id => {
    const menu = (
      <Fragment>
        <Item
          onClick={() => {
            callGetDetail(id);
          }}
        >
          UPDATE
        </Item>
        <Item
          onClick={() => {
            setSelectedId(id);
            setIsConfirmOpen(true);
          }}
        >
          REMOVE
        </Item>
      </Fragment>
    );
    return (
      <Dropdown overlay={menu} width={"max-content"}>
        <IconPlaceholder borderColor="none">
          <img src={MoreIcon} alt="analytic setting" />
        </IconPlaceholder>
      </Dropdown>
    );
  };

  const closeModal = () => {
    const query = history.location.search.replace("?", "");
    callGetList(query);

    setOpenAdd(false);
  };

  const closeModalBatchAdd = flag => {
    if (flag) {
      const query = history.location.search.replace("?", "");
      callGetList(query);
    }
    setOpenBatch(false);
  };

  function onSubmitSearch() {
    const getQueryUrl = history.location.search.replace("?", "");
    const queryUrl = qs.parse(getQueryUrl);
    queryUrl.search = searchValue;
    const query = qs.stringify(queryUrl);
    history.replace(`${history.location.pathname}?${query}`);
  }

  useEffect(() => {
    const timeout = setTimeout(() => {
      onSubmitSearch();
    }, 700);
    return () => clearTimeout(timeout);
  }, [searchValue]);

  useEffect(() => {
    const query = history.location.search.replace("?", "");
    callGetList(query);
  }, [history.location]);

  const RegisterMenu = (
    <Fragment>
      <Item onClick={() => setOpenAdd(true)}>
        <IconPlaceholder width="24px" height="24px" className="mr-4">
          <img src={SinglePerson} alt="single-person" />
        </IconPlaceholder>
        SINGLE REGISTER
      </Item>
      <Item onClick={() => setOpenBatch(true)}>
        <IconPlaceholder width="24px" height="24px" className="mr-4">
          <img src={MultiplePerson} alt="multiple-person" />
        </IconPlaceholder>
        BATCH REGISTER
      </Item>
    </Fragment>
  );

  return (
    <Fragment>
      <PageWrapper
        title="Vehicles"
        history={history}
        buttonGroup={
          <Fragment>
            <Button
              type="secondary"
              onClick={() => setIsConfirmAllOpen(true)}
              className="mr-16"
            >
              <img src={DeleteIcon} className="plus-icon" alt="Delete" />
              REMOVE ALL VEHICLES
            </Button>
            <Dropdown overlay={RegisterMenu} width="100%" className="mr-16">
              <Button type="blue">
                <img src={PlusCircleIcon} className="plus-icon" alt="Plus" />
                REGISTER VEHICLE
                <img src={ArrowDown} alt="dropdown" />
              </Button>
            </Dropdown>
          </Fragment>
        }
      >
        <SearchBarWrapper>
          <InputWrapper>
            <Input
              addonBefore={<img src={SearchIcon} alt="search-icon" />}
              placeholder="Type to search..."
              value={searchValue}
              onChange={event => setSearchValue(event.target.value)}
              onKeyDown={event => {
                if (event.key === "Enter") {
                  onSubmitSearch();
                }
              }}
            />
          </InputWrapper>
        </SearchBarWrapper>
        <ContentContainer>
          <Pagination
            limit={LIMIT_DATA}
            totalPage={totalPage}
            totalData={totalData}
          >
            <CustomLayout
              data={list}
              columns={dataSchema}
              editAction={""}
              deleteAction={""}
              elmAction={elmDropdown}
            />
          </Pagination>
        </ContentContainer>
      </PageWrapper>
      <Modal
        type="confirmation"
        isShown={isConfirmOpen}
        onConfirm={() => {
          callDelete(selectedId);
          setIsConfirmOpen(false);
        }}
        onClose={() => setIsConfirmOpen(false)}
        title="Delete Vehicle"
        buttonTitle="Proceed Delete"
        header="ARE YOU SURE TO PROCEED
        DELETING THIS DATA?"
        headerColor={themeContext.inlineError}
      />
      <Modal
        type="confirmation"
        isShown={isConfirmAllOpen}
        isButtonConfirmLoading={isLoadingDeleteAll}
        withCode={true}
        onConfirm={() => callDelete()}
        onClose={() => setIsConfirmAllOpen(false)}
        title="Delete All Vehicles"
        buttonTitle="Proceed Delete"
        subtitle=""
        header="ARE YOU SURE TO PROCEED
        DELETE ALL DATA?"
        headerColor={themeContext.inlineError}
      />
      <AddEnrollment
        type="vehicle"
        openModal={openAdd}
        onClose={data => closeModal(data)}
        data={selectedData}
      />
      <BatchAddVehicle
        openModal={openBatch}
        onClose={refreshFlag => closeModalBatchAdd(refreshFlag)}
        data={selectedData}
      />
    </Fragment>
  );
}

Vehicle.propTypes = {
  history: PropTypes.object.isRequired
};

export default Vehicle;

const ContentContainer = Styled.div`
  position: relative;
  overflow: hidden;
  height: calc( 100% - 140px);
`;

const CustomLayout = Styled(ListLayout)`
  padding:0px 15px;
`;
const Layout = Styled.div`
  position: relative;
  display: grid;
  grid-template-columns: repeat(8, calc((100%/8.6)));
  grid-gap: 33px 18px;
  padding: 9px;
  ${({ data }) =>
    data === 0 &&
    `
    grid-template-columns: unset;
    min-height: calc(100vh - 156px);
  `};
`;

const SearchBarWrapper = Styled.div`
  height: 40px
  width: 100%;
  display: flex;
  justify-content: flex-end;
  align-items: center;
  padding-top: 15px;
`;

const InputWrapper = Styled.div`
  width: 250px;
  padding-right: 20px;
`;

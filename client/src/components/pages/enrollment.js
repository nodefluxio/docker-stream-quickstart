/* eslint-disable no-unused-vars */
import React, { useState, Fragment, useEffect, useContext } from "react";
import { connect } from "react-redux";
import dayjs from "dayjs";
import Styled, { ThemeContext } from "styled-components";
import PropTypes from "prop-types";
import qs from "qs";
import { saveAs } from "file-saver";
import Button from "components/molecules/Button";
import Dropdown, { Menu as Item } from "components/atoms/Dropdown";
import IconPlaceholder from "components/atoms/IconPlaceholder";
import BlankContainer from "components/atoms/BlankContainer";
import Modal from "components/molecules/Modal";
import ThumbnailCard from "components/atoms/ThumbnailCard";
import Input from "components/atoms/Input";

import PageWrapper from "components/organisms/pageWrapper";
import {
  AddEnrollment,
  BatchEnrollment,
  BackupEnrollment
} from "components/organisms/modals";
import Pagination from "components/organisms/pagination";
import LoadingSpinner from "components/atoms/LoadingSpinner";

import {
  getListEnrollment,
  deleteEnrollment,
  getEnrollment,
  backupEnrollment,
  deleteAllEnrollment
} from "api";

import PlusCircleIcon from "assets/icon/visionaire/plus-circle.svg";
import MoreIcon from "assets/icon/visionaire/more.svg";
import NoRecord from "assets/icon/visionaire/no-record.svg";
import SearchIcon from "assets/icon/visionaire/search.svg";
import ArrowDown from "assets/icon/visionaire/drop-down.svg";
import SinglePerson from "assets/icon/visionaire/Union.svg";
import MultiplePerson from "assets/icon/visionaire/collabolators.svg";
import Backup from "assets/icon/visionaire/backup.svg";
import DeleteIcon from "assets/icon/visionaire/delete.svg";

import { REACT_APP_API_ENROLLMENT } from "config";

const LIMIT_DATA = 16;

function Enrollment(props) {
  const { history, agent } = props;
  const [dataEnrollment, setDataEnrollment] = useState([]);
  const [openAdd, setOpenAdd] = useState(false);
  const [openBatch, setOpenBatch] = useState(false);
  const [isConfirmOpen, setIsConfirmOpen] = useState(false);
  const [isBackupInProgress, setBackupInProgress] = useState(false);
  const [selectedId, setSelectedId] = useState(null);
  const [selectedData, setSelectedData] = useState([]);
  const themeContext = useContext(ThemeContext);
  const [totalPage, setTotalPage] = useState(1);
  const [totalData, setTotalData] = useState(0);
  const [noData, setNoData] = useState(false);
  const [errMsg, setErrMsg] = useState("");
  const [searchValue, setSearchValue] = useState(undefined);
  const [isConfirmAllOpen, setIsConfirmAllOpen] = useState(false);
  const [isLoadingDeleteAll, setIsLoadingDeleteAll] = useState(false);
  const [isDisconnected, setIsDisconnected] = useState(false);
  const [defaultQuery, setDefaultQuery] = useState("");
  const [dataLoading, setDataLoading] = useState(false);

  const getDetail = id => {
    getEnrollment(id).then(result => {
      if (result.ok) {
        setSelectedData(result.enrollment);
        setOpenAdd(true);
      }
    });
  };

  const backupEnrollmentHandler = () => {
    setBackupInProgress(true);
    backupEnrollment()
      .then(res => {
        const fileName = `data_${dayjs().format("YYYY-MM-DD")}.zip`;
        saveAs(res, fileName);
        setBackupInProgress(false);
      })
      .catch(err => {
        // eslint-disable-next-line no-console
        console.error(err);
        setBackupInProgress(false);
      });
  };

  const callGetListEnrollment = query => {
    setDataLoading(true);
    getListEnrollment(query)
      .then(result => {
        if (result.ok) {
          const { results } = result;
          if (results.total_data > 0) {
            setDataEnrollment(results.enrollments);
            setTotalPage(results.total_page);
            setTotalData(results.total_data);
            setNoData(false);
          } else {
            setErrMsg("");
            setNoData(true);
            setDataEnrollment([]);
            setTotalData(0);
            setTotalPage(0);
          }
        } else {
          setErrMsg("");
          setNoData(true);
          setDataEnrollment([]);
          setTotalData(0);
          setTotalPage(0);
        }
        setDataLoading(false);
      })
      .catch(err => {
        setErrMsg(err.message);
        setNoData(true);
        setDataEnrollment([]);
        setTotalData(0);
        setTotalPage(0);
        setDataLoading(false);
      });
  };

  const callDeleteEnrollment = id => {
    if (id !== undefined) {
      deleteEnrollment(id).then(result => {
        if (result.ok) {
          setIsConfirmOpen(false);
          setSelectedId(null);
          const query = history.location.search.replace("?", "");
          callGetListEnrollment(query);
        }
      });
    } else if (id === undefined) {
      setIsLoadingDeleteAll(true);
      deleteAllEnrollment()
        .then(result => {
          if (result.ok) {
            setIsLoadingDeleteAll(false);
            setIsConfirmAllOpen(false);
            const query = history.location.search.replace("?", "");
            callGetListEnrollment(query);
          }
        })
        .catch(() => setIsLoadingDeleteAll(false));
    }
  };

  const elmDropdown = data => {
    const menu = (
      <Fragment>
        <Item
          onClick={() => {
            getDetail(data.id);
          }}
        >
          UPDATE THIS ENROLLMENT
        </Item>
        <Item
          onClick={() => {
            setSelectedId(data.id);
            setIsConfirmOpen(true);
          }}
        >
          REMOVE THIS ENROLLMENT
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

  const closeModal = data => {
    if (data) {
      const query = history.location.search.replace("?", "");
      callGetListEnrollment(query);
    }
    setOpenAdd(false);
  };

  const closeModalBatchErollment = flag => {
    if (flag) {
      const query = history.location.search.replace("?", "");
      callGetListEnrollment(query);
    }
    setOpenBatch(false);
  };

  function onSubmitSearch() {
    const getQueryUrl = history.location.search.replace("?", "");
    const queryUrl = qs.parse(getQueryUrl);
    queryUrl.search = searchValue;
    queryUrl.page = 1;
    const query = qs.stringify(queryUrl);
    history.replace(`${history.location.pathname}?${query}`);
  }

  function checkAgentStatus(callback) {
    if (agent === "Connected") {
      callback();
    } else {
      setIsDisconnected(true);
    }
  }

  useEffect(() => {
    let timeout = setTimeout(() => {});
    if (searchValue !== undefined) {
      timeout = setTimeout(() => {
        onSubmitSearch();
      }, 700);
    }
    return () => clearTimeout(timeout);
  }, [searchValue]);

  useEffect(() => {
    if (history.location.search) {
      const query = history.location.search.replace("?", "");
      if (query !== defaultQuery) {
        callGetListEnrollment(query);
        setDefaultQuery(query);
      }
    }
  }, [history.location]);

  const AddEnrollmentMenu = (
    <Fragment>
      <Item onClick={() => checkAgentStatus(() => setOpenAdd(true))}>
        <IconPlaceholder width="24px" height="24px" className="mr-4">
          <img src={SinglePerson} alt="single-person" />
        </IconPlaceholder>
        SINGLE ENROLLMENT
      </Item>
      <Item onClick={() => checkAgentStatus(() => setOpenBatch(true))}>
        <IconPlaceholder width="24px" height="24px" className="mr-4">
          <img src={MultiplePerson} alt="multiple-person" />
        </IconPlaceholder>
        BATCH ENROLLMENT
      </Item>
      <Item onClick={() => checkAgentStatus(() => backupEnrollmentHandler())}>
        <IconPlaceholder width="24px" height="24px" className="mr-4">
          <img src={Backup} alt="multiple-person" />
        </IconPlaceholder>
        BACKUP ENROLLMENT
      </Item>
    </Fragment>
  );

  return (
    <Fragment>
      <PageWrapper
        title="Face Enrollment"
        history={history}
        buttonGroup={
          <Fragment>
            <Button
              type="secondary"
              onClick={() => setIsConfirmAllOpen(true)}
              className="mr-16"
            >
              <img src={DeleteIcon} className="plus-icon" alt="Delete" />
              REMOVE ALL ENROLLMENT
            </Button>
            <Dropdown
              overlay={AddEnrollmentMenu}
              width="100%"
              className="mr-16"
            >
              <Button type="blue">
                <img src={PlusCircleIcon} className="plus-icon" alt="Plus" />
                NEW ENROLLMENT
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
            <Layout data={dataEnrollment.length}>
              {dataEnrollment.length > 0 ? (
                dataEnrollment.map(data => (
                  <ThumbnailCard
                    key={data.id}
                    id={data.id}
                    url={
                      data.faces.length > 0
                        ? `${REACT_APP_API_ENROLLMENT}/face/image/${data.faces[0].id}`
                        : ""
                    }
                    actions={elmDropdown(data)}
                    title={data.name}
                  />
                ))
              ) : (
                <BlankContainer>
                  {noData ? (
                    <Fragment>
                      <img src={NoRecord} alt="blank-icon" />
                      {errMsg !== "" ? errMsg : "No Data Found"}
                    </Fragment>
                  ) : (
                    <LoadingSpinner show={!noData} />
                  )}
                </BlankContainer>
              )}
            </Layout>
          </Pagination>
          {dataLoading && dataEnrollment.length > 0 && (
            <BlankContainer
              position="absolute"
              backgroundColor="rgba(0,0,0,0.3)"
              top="0"
            >
              <LoadingSpinner show={dataLoading} />
            </BlankContainer>
          )}
        </ContentContainer>
      </PageWrapper>
      <Modal
        type="confirmation"
        isShown={isConfirmOpen}
        onConfirm={() => {
          callDeleteEnrollment(selectedId);
          setIsConfirmOpen(false);
        }}
        onClose={() => setIsConfirmOpen(false)}
        title="Delete Enrollemnt"
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
        onConfirm={() => callDeleteEnrollment()}
        onClose={() => setIsConfirmAllOpen(false)}
        title="Delete All Enrollemnt"
        buttonTitle="Proceed Delete"
        subtitle=""
        header="ARE YOU SURE TO PROCEED
        DELETE ALL DATA?"
        headerColor={themeContext.inlineError}
      />
      <Modal
        type="confirmation"
        isShown={isDisconnected}
        title=""
        header={
          agent === "Disconnected"
            ? `Enrollment Task is Not Allowed`
            : `Syncronization is on Progress`
        }
        buttonTitle="I Understand"
        subtitle={
          agent === "Disconnected"
            ? `Creating, editing, or deleting enrollment is not allowed while disconnected to coordinator. Please check your connectivity or coordinator status, ensure that the Coordinator Status is connected and try again.`
            : `Syncronization with coordinator is on progress to ensure enrollment database consistency accross multiple deployment. Please wait until sync process is done.`
        }
        headerColor={
          agent === "Disconnected"
            ? themeContext.inlineError
            : themeContext.yellow
        }
        onClose={() => setIsDisconnected(false)}
        onConfirm={() => setIsDisconnected(false)}
        withCancelButton={false}
      />
      <AddEnrollment
        type="person"
        openModal={openAdd}
        onClose={data => closeModal(data)}
        data={selectedData}
      />
      <BatchEnrollment
        openModal={openBatch}
        onClose={refreshFlag => closeModalBatchErollment(refreshFlag)}
        data={selectedData}
      />
      <BackupEnrollment openModal={isBackupInProgress} />
    </Fragment>
  );
}

Enrollment.propTypes = {
  history: PropTypes.object.isRequired,
  agent: PropTypes.string.isRequired
};

function mapStateToProps(state) {
  return {
    agent: state.agent
  };
}

export default connect(mapStateToProps)(Enrollment);

const ContentContainer = Styled.div`
  position: relative;
  margin-top: 20px;
  overflow: hidden;
  height: calc( 100% - 140px);
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

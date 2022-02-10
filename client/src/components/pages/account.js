import React, { useState, useEffect, useContext, Fragment } from "react";
import PropTypes from "prop-types";
import Styled, { ThemeContext } from "styled-components";
import PageWrapper from "components/organisms/pageWrapper";
import Button from "components/molecules/Button";
import Input from "components/molecules/Input";
import { getCookie } from "helpers/cookies";
import Dropdown, { Menu } from "components/atoms/Dropdown";
import IconPlaceholder from "components/atoms/IconPlaceholder";
import Modal from "components/molecules/Modal";
import LoadingSpinner from "components/atoms/LoadingSpinner";
import BlankContainer from "components/atoms/BlankContainer";
import ChangePasswordModal from "components/organisms/modals/ChangePassword";
import { useDispatch } from "react-redux";
import { showGeneralNotification } from "store/actions/notification";
import AddUser from "components/organisms/modals/AddUser";
import { deleteUser, getUsers } from "api/users";
import Table from "components/molecules/Table";
import Text from "components/atoms/Text";
import AreaSection from "components/molecules/AreaSection";
import parseJwt from "helpers/parseJWT";
import PaginationWrapper from "components/organisms/pagination";

import PlusCircle from "assets/icon/visionaire/plus-circle.svg";
import SettingIcon from "assets/icon/visionaire/more.svg";

const LIMIT_DATA = 16;

export default function AccountPage(props) {
  const { history } = props;
  const dispatch = useDispatch();
  const themeContext = useContext(ThemeContext);
  const [openAddModal, setOpenAddModal] = useState(false);
  const [users, setUsers] = useState([]);
  const [userProfile, setUserProfile] = useState({});
  const [selectedUser, setSelectedUser] = useState({});
  const [isConfirmOpen, setIsConfirmOpen] = useState(false);
  const [selectedID, setSelectedID] = useState("");
  const [totalPage, setTotalPage] = useState(1);
  const [totalData, setTotalData] = useState(0);
  const [btnConfirmLoading, setBtnConfirmLoading] = useState(false);
  const [openChangePasswordModal, setOpenChangePasswordModal] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const [errMsg, setErrMsg] = useState("");

  const column = [
    {
      key: "fullname",
      title: "Name"
    },
    {
      key: "role",
      title: "Role"
    },
    {
      key: "email",
      title: "Email"
    },
    {
      key: "action",
      title: ""
    }
  ];

  function deleteMember(id) {
    setSelectedID(id);
    setIsConfirmOpen(true);
  }

  function editMember(data) {
    setSelectedUser(data);
    setOpenAddModal(true);
  }

  function callListUser() {
    setIsLoading(true);
    const query = history.location.search.replace("?", "");
    getUsers(query)
      .then(result => {
        if (result.ok) {
          const { results } = result;
          for (let i = 0; i < results.users.length; i += 1) {
            results.users[i].action = (
              <Dropdown
                id={`user_option_${results.users[i].id}_button`}
                overlay={
                  <Fragment>
                    <Menu
                      id={`edit_member_${results.users[i].id}_button`}
                      onClick={() => editMember(results.users[i])}
                    >
                      EDIT USER
                    </Menu>
                    <Menu
                      id={`change_password_${results.users[i].id}_button`}
                      onClick={() => {
                        setOpenChangePasswordModal(true);
                        setSelectedUser(results.users[i]);
                      }}
                    >
                      CHANGE PASSWORD USER
                    </Menu>
                    <Menu
                      id={`delete_member_${results.users[i].id}_button`}
                      onClick={() => deleteMember(results.users[i].id)}
                    >
                      DELETE USER
                    </Menu>
                  </Fragment>
                }
                width="216px"
              >
                <IconPlaceholder>
                  <img id="action" src={SettingIcon} alt="more" />
                </IconPlaceholder>
              </Dropdown>
            );
          }
          setUsers(results.users);
          setTotalData(results.total_data);
          setTotalPage(results.total_page);
          setIsLoading(false);
        }
      })
      .catch(error => {
        setErrMsg(error.message);
        setIsLoading(false);
      });
  }

  function callDeleteMember(id) {
    // eslint-disable-next-line no-console
    setBtnConfirmLoading(true);
    deleteUser(id)
      .then(result => {
        if (result.ok) {
          dispatch(
            showGeneralNotification({
              type: "success",
              desc: "success delete user"
            })
          );
          callListUser();
          setBtnConfirmLoading(false);
        }
      })
      .catch(() => {
        dispatch(
          showGeneralNotification({
            type: "error",
            desc: "error delete user"
          })
        );
        setBtnConfirmLoading(false);
      });
  }

  useEffect(() => {
    callListUser();
    const accessToken = getCookie("access_token");
    const userCookies = parseJwt(accessToken);
    setUserProfile({
      name: userCookies.name,
      email: userCookies.email,
      role: userCookies.role
    });
  }, []);

  return (
    <PageWrapper
      title="Account Setting"
      history={history}
      buttonGroup={
        userProfile.role === "superadmin" && (
          <Button
            id="assign-analytics"
            onClick={() => {
              setOpenAddModal(true);
              setSelectedUser({});
            }}
            type="blue"
            className="button-control icon-text"
          >
            <img alt="plus-circle-icon" src={PlusCircle} />
            ADD USER
          </Button>
        )
      }
    >
      <ContentWrapper>
        <Left isAdmin={userProfile.role === "superadmin"}>
          <AreaSection
            title="PROFILE INFORMATION"
            titleColor={themeContext.mint}
          >
            <AreaWrapper>
              <InputLikeDiv>
                <span>User Role</span>
                <span>{userProfile.role || ""}</span>
              </InputLikeDiv>
              <Input label="Name" defaultValue={userProfile.name} readOnly />
              <Input
                label="Email"
                defaultValue={userProfile.email || ""}
                readOnly
              />
            </AreaWrapper>
          </AreaSection>
        </Left>
        <Right isAdmin={userProfile.role === "superadmin"}>
          {users.length > 0 ? (
            <PaginationWrapper
              limit={LIMIT_DATA}
              totalPage={totalPage}
              totalData={totalData}
            >
              <Table
                columns={column}
                dataSource={users}
                client={false}
                paddingColumn={["12px", 0]}
                rowHeight="45.2px"
                bodyHeight="100%"
              />
            </PaginationWrapper>
          ) : (
            <BlankContainer>
              {isLoading ? (
                <LoadingSpinner show={isLoading} />
              ) : (
                <Text size="18" color="#fff" weight="600">
                  {errMsg === "" ? "NO DATA FOUND" : errMsg}
                </Text>
              )}
            </BlankContainer>
          )}
        </Right>
      </ContentWrapper>
      <AddUser
        openModal={openAddModal}
        onClose={refreshFlag => {
          setOpenAddModal(false);
          setSelectedUser({});
          if (refreshFlag) {
            callListUser();
          }
        }}
        data={selectedUser}
      />
      <ChangePasswordModal
        openModal={openChangePasswordModal}
        onClose={() => setOpenChangePasswordModal(false)}
        data={selectedUser}
      />
      <Modal
        type="confirmation"
        isShown={isConfirmOpen}
        onConfirm={() => {
          callDeleteMember(selectedID);
          setIsConfirmOpen(false);
        }}
        onClose={() => setIsConfirmOpen(false)}
        title="Delete User"
        buttonTitle="Proceed Delete"
        header="ARE YOU SURE TO PROCEED
        DELETING THIS USER?"
        headerColor={themeContext.inlineError}
        isButtonConfirmLoading={btnConfirmLoading}
      />
    </PageWrapper>
  );
}

AccountPage.propTypes = {
  history: PropTypes.object.isRequired
};

const ContentWrapper = Styled.div`
  display: flex;
  flex-direction: row;
  height: calc(100% - 65px);
  justify-content: center;
`;

const AreaWrapper = Styled.div`
  padding: 15px;
`;

const InputLikeDiv = Styled.div`
  width: 100%;
  display: flex;
  justify-content: space-between;
  align-items: center;
  border: none;
  height: 40px;
  padding-left: 10px;
  padding-right: 10px;
  box-sizing: border-box;
  border-radius: 8px;
  color: ${props => props.theme.mercury};
  border: 2px solid ${props => props.theme.mint};
  font-size: 14px;
  position: relative;
  background-color: ${props => props.theme.secondary2};
  span:last-child {
    font-weight: 800;
    text-transform: uppercase;
  }
`;

const Left = Styled.div`
    width: ${props => (props.isAdmin ? `30%` : `50%`)};
    ${props => !props.isAdmin && `border-left: 1px solid #372463;`}
    border-right: 1px solid #372463;
    height: 100%;
`;

const Right = Styled.div`
    height: 100%;
    width: ${props => (props.isAdmin ? `70%` : `0%`)};
    display: ${props => (props.isAdmin ? `flex` : `none`)};
    flex-direction: column;
`;

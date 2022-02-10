import React, { Fragment, useState, useEffect } from "react";
import { withRouter } from "react-router-dom";
import { connect, useDispatch } from "react-redux";
import PropTypes from "prop-types";
import { AddEnrollment } from "components/organisms/modals";
import {
  closeModalEnrollment,
  openModalEnrollment
} from "store/actions/eventContext";
import { PLUGIN_NAME } from "config";
import { MenuWrapper, ItemWrapper } from "components/molecules/ContextMenu";
import { getCookie } from "helpers/cookies";
import parseJwt from "helpers/parseJWT";

export default function withEnrollment(Component) {
  function AddEnrollmentWrapper(props) {
    const dispatch = useDispatch();
    const { history, eventContext } = props;
    const { openModal, type, data } = eventContext;
    const [open, setOpen] = useState(false);
    const [role, setRole] = useState("operator");

    function closeModal(returnData) {
      if (
        returnData &&
        Object.keys(returnData).length > 1 &&
        returnData.constructor === Object
      ) {
        history.push("enrollment");
      }
      setOpen(false);
      dispatch(closeModalEnrollment());
    }

    useEffect(() => {
      setOpen(openModal);
    }, [openModal]);

    useEffect(() => {
      const accessToken = getCookie("access_token");
      const userCookies = parseJwt(accessToken);
      setRole(userCookies.role);
    }, []);

    function selectDataToEnroll(item, EnrollType) {
      switch (EnrollType) {
        case "person":
          return item.secondary_image;
        case "vehicle":
          return item.label;
        default:
          return "";
      }
    }

    return (
      <Fragment>
        <Component {...props} />
        <MenuWrapper id="event-context">
          {data.analytic_id === "NFV4-FR" && data.label === "unrecognized" && (
            <ItemWrapper onClick={() => dispatch(openModalEnrollment())}>
              Enroll this {eventContext.type}
            </ItemWrapper>
          )}
          {data.analytic_id === "NFV4-FR" &&
            PLUGIN_NAME === "Search" &&
            role === "superadmin" && (
              <ItemWrapper
                onClick={() => history.push("/plugin/search/person")}
              >
                Search this person
              </ItemWrapper>
            )}
        </MenuWrapper>
        <AddEnrollment
          type={type}
          openModal={open}
          onClose={returnData => closeModal(returnData)}
          data={selectDataToEnroll(data, type)}
        />
      </Fragment>
    );
  }

  AddEnrollmentWrapper.propTypes = {
    eventContext: PropTypes.object.isRequired,
    history: PropTypes.object.isRequired
  };

  function mapStateToProps(state) {
    return {
      eventContext: state.eventContext
    };
  }

  return withRouter(connect(mapStateToProps)(AddEnrollmentWrapper));
}

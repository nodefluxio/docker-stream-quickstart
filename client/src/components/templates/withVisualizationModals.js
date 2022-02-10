import React, { Fragment, useContext, useEffect, useState } from "react";
import Styled, { ThemeContext } from "styled-components";
import { connect, useDispatch } from "react-redux";
import PropTypes from "prop-types";
import Modal from "components/molecules/Modal";
import Text from "components/atoms/Text";
import AreaSection from "components/molecules/AreaSection";
import Row from "components/atoms/Row";

import { getCamera } from "api";
import {
  closeConfirmationModal,
  closeInformationModal
} from "store/actions/cameraMenu";

export default function withModalVisualization(Component) {
  function VisualizationModals(props) {
    const dispatch = useDispatch();
    const [seatList, setSeatList] = useState([]);
    const [data, setData] = useState({});
    const { cameraMenu } = props;
    const themeContext = useContext(ThemeContext);

    const formatLabel = string => {
      const splitString = string.split("_");
      let keys = "";
      for (let i = 1; i < splitString.length; i += 1) {
        keys +=
          splitString[i].charAt(0).toUpperCase() + splitString[i].slice(1);
        if (i < splitString.length - 1) {
          keys += " ";
        }
      }
      return `Camera ${keys}`;
    };

    useEffect(() => {
      if (cameraMenu.showInfo) {
        getCamera(cameraMenu.selectedID, null, cameraMenu.selectedNode)
          .then(result => {
            if (result.ok) {
              const streamData = result.stream;
              setData(streamData);
              if (streamData.pipelines !== undefined) {
                const seats = [];
                if (streamData.seats && streamData.seats.length > 0) {
                  streamData.seats.forEach(value => {
                    const analyticId = value.analytic_id;
                    seats[analyticId] = value.serial_number;
                  });
                }
                setSeatList(seats);
              }
            }
          })
          .catch(() => {});
      }
    }, [cameraMenu.showInfo, cameraMenu.selectedID]);

    return (
      <Fragment>
        <Component {...props} />
        <Modal
          type="confirmation"
          isShown={cameraMenu.showConfirmation}
          onConfirm={() => {
            cameraMenu.deleteFunction(cameraMenu.selectedID);
            dispatch(closeConfirmationModal());
          }}
          onClose={() => dispatch(closeConfirmationModal())}
          title="Delete Camera"
          buttonTitle="Proceed Delete"
          header="ARE YOU SURE TO PROCEED
          DELETING ALL THE DATA?"
          headerColor={themeContext.color8}
        />
        <Modal
          show={cameraMenu.showInfo}
          close={() => dispatch(closeInformationModal())}
          title="Camera Information Detail"
          className="modal-camera-info"
          width="600px"
        >
          <AreaSection
            title="CAMERA INFORMATION"
            titleColor={themeContext.mint}
            border
          >
            <div>
              {Object.keys(data).map(item => {
                if (typeof data[item] !== "object") {
                  return (
                    <Fragment key={item}>
                      <Row
                        justify="space-between"
                        align="center"
                        horizontalPadding={15}
                        height="48px"
                        key={item}
                      >
                        <Text
                          size="14px"
                          color={themeContext.white}
                          weight={500}
                        >
                          {formatLabel(item)}
                        </Text>
                        <Text
                          size="14px"
                          color={themeContext.white}
                          weight={500}
                        >
                          {data[item]}
                        </Text>
                      </Row>
                      <Divider color={themeContext.secondary2} />
                    </Fragment>
                  );
                }
                return null;
              })}
            </div>
          </AreaSection>
          <AreaSection
            title="LICENSE INFORMATION"
            titleColor={themeContext.mint}
            border
          >
            {Object.keys(seatList).map(item => (
              <Row
                justify="space-between"
                align="center"
                horizontalPadding={15}
                height="48px"
                key="serial-number"
              >
                <Text size="14px" color={themeContext.white} weight={500}>
                  {item}
                </Text>
                <Text size="14px" color={themeContext.white} weight={500}>
                  {seatList[item]}
                </Text>
              </Row>
            ))}
          </AreaSection>
        </Modal>
      </Fragment>
    );
  }

  const Divider = Styled.div`
  width: 100%;
  height: 1px;
  ${({ color }) => color && `background-color: ${color};`}
`;

  VisualizationModals.propTypes = {
    cameraMenu: PropTypes.object.isRequired
  };

  function mapStateToProps(state) {
    return {
      cameraMenu: state.cameraMenu
    };
  }

  return connect(mapStateToProps)(VisualizationModals);
}

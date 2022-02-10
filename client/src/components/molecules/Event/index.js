import React, { useContext } from "react";
import { ContextMenuTrigger } from "react-contextmenu";
import Styled, { ThemeContext } from "styled-components";
import PropTypes from "prop-types";
import Row from "components/atoms/Row";
import Text from "components/atoms/Text";
import ErrorFileIcon from "assets/icon/visionaire/file-error.svg";

const defaultProps = {
  primaryImage: "",
  secondaryImage: "",
  rightOption: null,
  pinned: false,
  pinfunc: null,
  color: "white",
  result: "",
  location: "",
  onContextMenu: () => {}
};

const propTypes = {
  color: PropTypes.string,
  label: PropTypes.string.isRequired,
  date: PropTypes.string.isRequired,
  location: PropTypes.string,
  result: PropTypes.string,
  primaryImage: PropTypes.string.isRequired,
  secondaryImage: PropTypes.string,
  rightOption: PropTypes.element,
  pinned: PropTypes.bool,
  pinfunc: PropTypes.func,
  onEnlargePrimaryImage: PropTypes.func,
  onContextMenu: PropTypes.func
};

function Event(props) {
  const {
    color,
    label,
    date,
    location,
    result,
    primaryImage,
    secondaryImage,
    rightOption,
    pinned,
    pinfunc,
    onEnlargePrimaryImage,
    onContextMenu
  } = props;
  const themeContext = useContext(ThemeContext);
  return (
    <ContextMenuTrigger id="event-context">
      <CardList highlight={pinfunc && pinned} onContextMenu={onContextMenu}>
        <Row align="center" horizontalGap={5}>
          <Row>
            <BorderContent
              width={64}
              type="solid"
              radius={6}
              color={color}
              size={3}
              onClick={onEnlargePrimaryImage}
            >
              {primaryImage === "" ||
              String(primaryImage) === "data:image/jpeg;base64," ? (
                <img src={ErrorFileIcon} alt="primary" />
              ) : (
                <img
                  src={primaryImage}
                  onError={e => {
                    e.target.onerror = null;
                    e.target.src = ErrorFileIcon;
                  }}
                  alt="primary"
                  loading="lazy"
                />
              )}
            </BorderContent>

            {secondaryImage !== "" && (
              <BorderContent
                width={64}
                type="solid"
                radius={6}
                color={themeContext.white}
                size={1}
              >
                <img
                  src={secondaryImage}
                  alt="secondary"
                  onError={e => {
                    e.target.onerror = null;
                    e.target.src = ErrorFileIcon;
                  }}
                  loading="lazy"
                />
              </BorderContent>
            )}
          </Row>
          <Row direction="column" justify="space-evenly">
            <Text
              size="14"
              marginBottom="5px"
              color={color}
              weight="bold"
              textTransform="uppercase"
            >
              {label}
            </Text>
            <Text size="14" marginBottom="5px" color={themeContext.white}>
              {result}
            </Text>
            {location && (
              <Text size="14" marginBottom="5px" color={themeContext.white}>
                location: {location}
              </Text>
            )}
            <Text size="14" marginBottom="5px" color={themeContext.white}>
              {date}
            </Text>
          </Row>
        </Row>
        {rightOption}
      </CardList>
    </ContextMenuTrigger>
  );
}

Event.defaultProps = defaultProps;

Event.propTypes = propTypes;

export default Event;

const BorderContent = Styled.div`
  ${({ width }) => width && `width:${width}px; height:${width}px;`}
  ${({ type }) => type && `border:${type};`}
  ${({ color }) => color && `border-color:${color};`}
  ${({ size }) => size && `border-width: ${size}px;`}
  ${({ radius }) => radius && `border-radius: ${radius}px;`}
  ${({ onClick }) => onClick && `cursor:pointer;`}
  overflow: hidden;
  display:flex;
  align-items: center;
  justify-content: center;
  img { 
    @media (orientation: landscape) {
      max-height: 100%
    }
    @media (orientation: potrait) {
      max-width: 100%;
    }
  }
`;

const CardList = Styled.div`
  display: flex;
  justify-content: space-between;
  border: 2px solid ${props => props.theme.secondary2};
  width: calc(100% - 20px);
  border-radius: 8px;
  margin: 10px;
  height: 90px;
  ${props => props.highlight && `background: ${props.theme.secondary2};`}
`;

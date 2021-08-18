import React from "react";
import PropTypes from "prop-types";
import Styled from "styled-components";
import Setting from "assets/icon/visionaire/more.svg";
import Dropdown from "components/atoms/Dropdown";

export default function VisualisationWrapper(props) {
  const {
    children,
    width,
    height,
    name,
    menu,
    secondaryMenu,
    secondaryMenuImage,
    maxheight,
    overlayIcon,
    onClick,
    url
  } = props;

  return (
    <ImageWrapper
      height={height}
      width={width}
      maxheight={maxheight}
      data-url={url}
      className="nodeflux-streamer"
    >
      {(name || menu) && (
        <ImageOverlay>
          <ImageTitle>{name}</ImageTitle>
          {secondaryMenu && (
            <Dropdown overlay={secondaryMenu} width="216px">
              <IconSetting>
                <img src={secondaryMenuImage || Setting} alt="setting-icon" />
              </IconSetting>
            </Dropdown>
          )}
          {menu && (
            <Dropdown overlay={menu} width="216px">
              <IconSetting>
                <img src={Setting} alt="setting-icon" />
              </IconSetting>
            </Dropdown>
          )}
        </ImageOverlay>
      )}
      {overlayIcon && (
        <OverlayWrapper>
          <img src={overlayIcon} alt="overlay-icon" />
        </OverlayWrapper>
      )}
      <ClickableArea onClick={onClick} />
      {children}
    </ImageWrapper>
  );
}

VisualisationWrapper.propTypes = {
  url: PropTypes.string.isRequired,
  children: PropTypes.element,
  width: PropTypes.string,
  height: PropTypes.string,
  name: PropTypes.string,
  menu: PropTypes.element,
  maxheight: PropTypes.string,
  secondaryMenu: PropTypes.oneOfType([PropTypes.object, PropTypes.bool]),
  secondaryMenuImage: PropTypes.string,
  overlayIcon: PropTypes.string,
  onClick: PropTypes.func
};

VisualisationWrapper.defaultProps = {
  children: null,
  width: "448px",
  height: "336px",
  name: null,
  menu: null,
  maxheight: "336px",
  secondaryMenu: false,
  secondaryMenuImage: "",
  overlayIcon: "",
  onClick: () => {}
};

const ImageWrapper = Styled.div`
    display: flex;
    width: ${props => props.width};
    height: ${props => props.height};
    background-position: center;
    background-repeat: no-repeat;
    background-size: auto;
    position: relative;
`;

const ImageOverlay = Styled.div`
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 40px;
  background-color: rgba(24,24,24,0.6);
  display: flex;
  flex-direction: row;
  z-index: 5;
`;

const ImageTitle = Styled.div`
  text-transform: uppercase;
  margin: 12px auto 14px 12px;
`;

const IconSetting = Styled.div`
  display: flex;
  width: 40px;
  height: 40px;
  align-items: center;
  justify-content: center;
  ${props => props.active && `background-color: ${props.theme.secondary2};`}
`;

const OverlayWrapper = Styled.div`
  width: 100%;
  height: 100%;
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  img {
    z-index: 1;
    height: 80px;
    width: 80px;
  }
  &:before {
    content: "";
    background: ${props => props.theme.secondary2};
    opacity: 0.4;
    width: 100%;
    height: 100%;
    position: absolute;
    top: 0;
    left: 0;
  }
`;

const ClickableArea = Styled.div`
  width: 100%;
  height: 100%;
  ${props => props.onClick && `cursor: pointer;`}
  position: absolute;
`;

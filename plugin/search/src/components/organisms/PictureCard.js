import React, { useState } from "react";
import PropTypes from "prop-types";
import Styled from "styled-components";
import ErrorFileIcon from "assets/icon/visionaire/broken-image.svg";

export default function PictureCard(props) {
  const { image, name, percentage, onClick, width } = props;
  const [isImageError, setIsImageError] = useState(false);
  return (
    <PictureCardWrapper onClick={onClick} width={width} error={isImageError}>
      <ImageWrapper bg={image} error={isImageError}>
        <img
          src={image}
          alt={name}
          onError={e => {
            e.target.onerror = null;
            e.target.src = ErrorFileIcon;
            setIsImageError(true);
          }}
        />
        <PercentageRibbon>{percentage}</PercentageRibbon>
      </ImageWrapper>
      {name}
    </PictureCardWrapper>
  );
}

PictureCard.propTypes = {
  image: PropTypes.string.isRequired,
  name: PropTypes.string,
  percentage: PropTypes.string,
  onClick: PropTypes.func,
  width: PropTypes.string
};

const PictureCardWrapper = Styled.div`
    height: 250px;
    cursor: pointer;
    text-transform: uppercase;
    font-size: 14px;
    width: ${props => props.width || "100%"};
    margin: 5px;
    display: flex;
    align-items: flex-start;
    flex-direction: column;
`;

const ImageWrapper = Styled.div`
    width: 100%;
    height: 100%;
    position: relative;
    margin-bottom: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
    overflow: hidden;
    ${props => props.error && `border: 1px solid #372463;`}
    img {
        max-width: 100%;
    }
`;
const PercentageRibbon = Styled.div`
    position: absolute;
    height: 30px;
    display: flex;
    align-items: center;
    bottom: 0;
    background: rgba(69, 69, 69, 0.7);
    width: calc(100% - 8px);
    padding-left: 8px;
    font-weight: 600;
`;

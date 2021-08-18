import React from "react";
import PropTypes from "prop-types";
import Styled from "styled-components";

export default function PictureCard(props) {
  const { image, name, percentage, onClick, width } = props;
  return (
    <PictureCardWrapper onClick={onClick} width={width}>
      <ImageWrapper bg={image}>
        <img src={image} alt={name} />
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
    height: min-content;
    width: 100%;
    cursor: pointer;
    text-transform: uppercase;
    font-size: 14px;
    width: ${props => props.width || "220px"};
    margin: 5px;
`;

const ImageWrapper = Styled.div`
    width: 100%;
    position: relative;
    margin-bottom: 8px;
    img {
        max-width: 100%;
        max-height: 100%;
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

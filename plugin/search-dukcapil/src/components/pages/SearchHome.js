/* eslint-disable no-console */
import React, { useState, useRef, Fragment, useEffect } from "react";
import Styled from "styled-components";

import Input from "components/atoms/Input";
import Button from "components/molecules/Button";
import PictureCard from "components/organisms/PictureCard";
import IconPlaceholder from "components/atoms/IconPlaceholder";

import AddIcon from "assets/icon/visionaire/plus-circle.svg";
import CloseIcon from "assets/icon/visionaire/exit.svg";
import Accordion from "components/molecules/Accordion";

function SearchHome() {
  const [mode, setMode] = useState("");
  const [nik, setNik] = useState("");
  const [image, setImage] = useState("");
  const inputRef = useRef(null);
  const [trigger, setTrigger] = useState("");
  const [formData, setFormData] = useState("");
  const [searchResult, setSearchResult] = useState([]);
  const [detailResult, setDetailResult] = useState({});

  function onImageChange(event) {
    if (event.target.files && event.target.files[0]) {
      const blob = new Blob([event.target.files[0]], { type: "image/jpeg" });
      setImage(URL.createObjectURL(event.target.files[0]));
      setFormData({ images: blob });
    }
  }

  function onHandleNIK(value) {
    setNik(value);
  }

  function onSubmitSearchImage() {
    const newResult = [];
    for (let i = 0; i < 10; i += 1) {
      newResult.push({ name: "Ayunda Risu", percentage: "50%", image });
    }
    setSearchResult(newResult);
  }

  function onSelectDetail() {
    setDetailResult({
      image,
      name: "Ayunda Risu"
    });
  }

  useEffect(() => {
    if (mode === "nik") {
      setImage("");
    } else if (mode === "image") {
      setNik("");
    }
  }, [mode]);

  useEffect(() => {
    console.log(formData);
  }, [formData]);

  return (
    <Wrapper minHeight="100vh" height="100%">
      <ControlContainer>Database Citizen Search</ControlContainer>
      <Wrapper flexDirection="row" height="calc(100% - 64px)">
        <Left>
          <Accordion
            header="Search By Image"
            type={mode === "image" ? "open" : "closed"}
            callback={() => setMode("image")}
            trigger={trigger}
          >
            <SearchImageWrapper>
              <ImgWrapper
                border={true}
                onClick={() => inputRef.current.click()}
              >
                {image === "" ? (
                  <Fragment>
                    <img src={AddIcon} alt="add" className="add-icon" />
                    <br />
                    Upload Image
                  </Fragment>
                ) : (
                  <Fragment>
                    <img
                      src={image}
                      alt="image-search"
                      onLoad={() => setTrigger(image)}
                    />
                    <br />
                    Change Image
                  </Fragment>
                )}
              </ImgWrapper>
              <input
                type="file"
                onChange={e => onImageChange(e)}
                ref={inputRef}
              />
              <Button width="100%" onClick={() => onSubmitSearchImage()}>
                SEARCH
              </Button>
            </SearchImageWrapper>
          </Accordion>
          <Accordion
            header="Search By NIK"
            type={mode === "nik" ? "open" : "closed"}
            callback={() => setMode("nik")}
          >
            <Input
              type="text"
              label="Search People by NIK"
              placeholder="People ID Number..."
              style={{ marginBottom: "20px" }}
              value={nik}
              onChange={e => onHandleNIK(e.target.value)}
            />
            <Button width="100%">SEARCH</Button>
          </Accordion>
        </Left>
        <Right>
          <Ribbon>Search Result</Ribbon>
          <Wrapper flexDirection="row" height="calc(100% - 43px)">
            {Object.keys(detailResult).length > 0 && (
              <DetailWrapper>
                <Ribbon bg="transparent">
                  <span>Identity Detail</span>
                  <IconPlaceholder onClick={() => setDetailResult({})}>
                    <img src={CloseIcon} />
                  </IconPlaceholder>
                </Ribbon>
                <DetailImage>
                  <img src={detailResult.image} />
                </DetailImage>
                <DetailInfo>
                  <Row>
                    <div>Full Name</div>
                    <div>{detailResult.name}</div>
                  </Row>
                </DetailInfo>
              </DetailWrapper>
            )}
            {searchResult.length > 0 && (
              <ResultWrapper
                grid={Object.keys(detailResult).length > 0 ? 3 : 6}
              >
                {searchResult.map((result, index) => (
                  <PictureCard
                    key={index}
                    image={result.image}
                    percentage={result.percentage}
                    name={result.name}
                    onClick={onSelectDetail}
                    width={
                      Object.keys(detailResult).length > 0
                        ? `calc(${100 / 2}% - 10px)`
                        : `calc(${100 / 4}% - 10px)`
                    }
                  />
                ))}
              </ResultWrapper>
            )}
          </Wrapper>
        </Right>
      </Wrapper>
    </Wrapper>
  );
}

export default SearchHome;

const Wrapper = Styled.div`
    display: flex;
    flex-direction: ${props => props.flexDirection || `column`};
    width: 100%;
    height: ${props => props.height};
    min-height: ${props => props.minHeight};    
`;

const ControlContainer = Styled.div`
  display: flex;
  flex-direction: row;
  width: 100%;
  justify-content: space-between;
  align-items: center;
  height: 64px;
  text-transform: uppercase;
  color: ${props => props.theme.mint};
  font-weight: 600;
  font-size: 18px;
  padding-left: 15px;    
  .plus-icon{
      margin-right: 10px;
  }
`;

const Left = Styled(Wrapper)`
    width: 30%;
    border-right: 1px solid #372463;
    height: 100%;
`;

const Right = Styled(Wrapper)`
    height: 100%;
    width: 70%;
    flex-direction: column;
`;

const Ribbon = Styled.div`
  background-color: ${props => props.bg || `#372463`};
  color: #45E5B7;
  height: 40px;
  border: 1px solid #372463;
  font-family: "Barlow", sans-serif;
  display: flex;
  font-weight: 600;
  text-transform: uppercase;
  align-items: center;
  padding-left: 10px;
  justify-content: space-between;
`;

const SearchImageWrapper = Styled.div`
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  input {
    display: none;
  }
`;

const ImgWrapper = Styled.div`
  position: relative;
  min-width: 40%;
  max-width: 100%;
  padding: 10px;
  min-height: 200px;
  max-height: 100%
  text-align: center;
  ${props =>
    props.border &&
    `
  border: 1px solid white;
  cursor: pointer;
  `}
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  margin-bottom: 20px;
  img {
    max-width: 100%;
    max-height: 100%;
  }
  .add-icon {
    max-width: 40px;
  }
`;

const Row = Styled.div`
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: space-between;
  height: 48px;
  border-bottom: 1px solid #372463;
`;

const DetailWrapper = Styled.div`
  max-width: 40%;
  border-right: 1px solid #372463;
`;
const DetailImage = Styled.div`
  padding: 16px;
  img {
    max-width: 100%;
  }
`;
const DetailInfo = Styled.div`
  padding: 16px;  
`;
const ResultWrapper = Styled.div`
  width: ${props => (props.grid === 3 ? `60%` : `100%`)};
  display: flex;
  flex-direction: row;
  flex-wrap: wrap;
  padding: 20px 18px;
  overflow-y: auto;
`;

/* eslint-disable no-console */
import React, {
  useState,
  useRef,
  Fragment,
  useEffect,
  useContext
} from "react";
import Styled from "styled-components";
import PropTypes from "prop-types";

import Input from "components/atoms/Input";
import Button from "components/molecules/Button";
import PictureCard from "components/organisms/PictureCard";
import IconPlaceholder from "components/atoms/IconPlaceholder";
import EnvWrapper, { EnvContext } from "components/templates/EnvContext";
import { getCitizen, getToken, getResultByToken } from "api";

import AddIcon from "assets/icon/visionaire/plus-circle.svg";
import CloseIcon from "assets/icon/visionaire/exit.svg";
import ErrorIcon from "assets/icon/visionaire/file-error.svg";
import BrokenImage from "assets/icon/visionaire/broken-image.svg";

import Accordion from "components/molecules/Accordion";

const RETRY_LIMIT = 5;
const DELAY_SECONDS = 3000;
const MAX_SIMILARITY = 99;

function SearchHome({ history, store }) {
  const API = useContext(EnvContext);
  const [mode, setMode] = useState("");
  const [nik, setNik] = useState("");
  const [image, setImage] = useState("");
  const inputRef = useRef(null);
  const [trigger, setTrigger] = useState("");
  const [formData, setFormData] = useState("");
  const [searchResult, setSearchResult] = useState([]);
  const [detailResult, setDetailResult] = useState({});
  const [errorMsg, setErrorMsg] = useState("");
  const [detailImage, setDetailImage] = useState("");
  const [limitResult, setLimitResult] = useState("10");
  const [loading, setLoading] = useState(false);
  const [token, setToken] = useState("");
  const [minProbability, setMinProbability] = useState(50);
  const [fullSearchResult, setFullSearchResult] = useState([]);

  let retry = 0;

  const b64toBlob = (b64Data, contentType = "", sliceSize = 512) => {
    try {
      const byteCharacters = atob(b64Data);
      const byteArrays = [];

      for (
        let offset = 0;
        offset < byteCharacters.length;
        offset += sliceSize
      ) {
        const slice = byteCharacters.slice(offset, offset + sliceSize);

        const byteNumbers = new Array(slice.length);
        for (let i = 0; i < slice.length; i += 1) {
          byteNumbers[i] = slice.charCodeAt(i);
        }

        const byteArray = new Uint8Array(byteNumbers);
        byteArrays.push(byteArray);
      }

      const blob = new Blob(byteArrays, { type: contentType });
      return blob;
    } catch {
      return null;
    }
  };

  function onImageChange(event) {
    console.log(event);
    if (event.target.files && event.target.files[0]) {
      const blob = new Blob([event.target.files[0]], { type: "image/jpeg" });
      setImage(URL.createObjectURL(event.target.files[0]));
      setFormData({ image: blob });
    }
  }

  function onHandleNIK(value) {
    setNik(value);
  }

  function onSubmitSearchImage() {
    const newForm = new FormData();
    newForm.append("limit", parseInt(limitResult, 10));
    newForm.append("image", formData.image);
    setLoading(true);
    getToken(API, newForm)
      .then(result => {
        if (result.ok) {
          setErrorMsg("");
          setToken(result.token);
        }
      })
      .catch(error => {
        console.log(error);
        setLoading(false);
        setErrorMsg("Cannot get Token");
      });
  }

  function onChangeQuery(searchURLMode = false) {
    if (searchURLMode) {
      history.replace({
        search: `nik=${nik}`
      });
    } else history.replace({ search: null });
  }

  function onSelectDetail(data, percentage) {
    setDetailResult({ ...data, percentage });
    const imageData = b64toBlob(data.foto);
    if (imageData !== null) {
      setDetailImage(URL.createObjectURL(imageData));
    } else {
      setDetailImage("");
    }
  }

  useEffect(() => {
    if (mode === "nik") {
      setImage("");
    } else if (mode === "image") {
      setNik("");
      onChangeQuery();
    }
  }, [mode]);

  useEffect(() => {
    setSearchResult(
      fullSearchResult.filter(
        data =>
          data.probability >= minProbability && data.dukcapil_data !== null
      )
    );
  }, [minProbability]);

  useEffect(() => {
    if (store) {
      const { eventContext } = store;
      if (eventContext.data) {
        if (eventContext.data.secondary_image) {
          const blobImage = b64toBlob(eventContext.data.secondary_image);
          const event = {
            target: {
              files: [blobImage]
            }
          };
          onImageChange(event);
        }
      }
    }
  }, [store]);

  function getResult(tokenData) {
    getResultByToken(API, tokenData)
      .then(result => {
        if (result.ok) {
          setErrorMsg("");
          setLoading(false);
          const { results } = result;
          setDetailResult({});
          setFullSearchResult(results);
          setSearchResult(
            results.filter(
              data =>
                data.probability >= minProbability &&
                data.dukcapil_data !== null
            )
          );
          setToken("");
        }
      })
      .catch(error => {
        if (retry < RETRY_LIMIT) {
          setTimeout(() => {
            getResult(tokenData);
            retry += 1;
          }, DELAY_SECONDS);
        } else {
          setLoading(false);
          console.log(error);
          setErrorMsg("Failed to get data");
        }
      });
  }

  useEffect(() => {
    if (token !== "") {
      getResult(token);
    }
  }, [token]);

  useEffect(() => {
    const { search } = history.location;
    if (search !== "") {
      setSearchResult([]);
      const query = new URLSearchParams(search);
      const nikQuery = query.get("nik");
      if (nikQuery) {
        setNik(query.get("nik"));
        setMode("nik");
      }
      getCitizen(API, search)
        .then(result => {
          if (result.ok) {
            setDetailResult(result.citizen_data);
            const imageData = b64toBlob(result.citizen_data.foto);
            setDetailImage(URL.createObjectURL(imageData));
            setErrorMsg("");
          }
        })
        .catch(error => {
          setDetailResult({});
          setErrorMsg(error.response.data.errors[0]);
        });
    }
  }, [history.location.search]);

  function checkMaxSimilarity(probability) {
    if (probability >= MAX_SIMILARITY) {
      const maxProb = "99.99";
      return maxProb;
    }
    return probability;
  }

  return (
    <Wrapper height="100%">
      <ControlContainer>Database Citizen Search</ControlContainer>
      <Wrapper flexDirection="row" height="calc(100% - 64px)">
        <Left>
          <Accordion
            header="Search By Image"
            type={mode !== "nik" ? "open" : "closed"}
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
                style={{ display: "none" }}
                onChange={e => onImageChange(e)}
                ref={inputRef}
              />

              <Input
                type="number"
                pattern="[0-9]*"
                min={10}
                max={99}
                inputmode="numeric"
                label="Limit Result"
                placeholder="Limit result sorted by higher similarity"
                value={limitResult}
                onChange={e => setLimitResult(e.target.value)}
              />
              <Input
                type="slider"
                label="Minimum Confidence Level: "
                minValue={0}
                maxValue={100}
                value={minProbability}
                onChange={value => setMinProbability(value)}
                suffix="%"
              />
              <Button
                width="100%"
                onClick={() => onSubmitSearchImage()}
                isLoading={loading}
                disabled={!formData.image}
              >
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
            <Button
              width="100%"
              disabled={nik === ""}
              onClick={() => onChangeQuery(true)}
            >
              SEARCH
            </Button>
          </Accordion>
        </Left>
        <Right>
          <Ribbon>Search Result</Ribbon>
          <Wrapper flexDirection="row" height="calc(100% - 43px)">
            {Object.keys(detailResult).length > 0 && (
              <DetailWrapper>
                <Ribbon bg="transparent">
                  <span>Identity Detail</span>
                  {mode !== "nik" && (
                    <IconPlaceholder onClick={() => setDetailResult({})}>
                      <img src={CloseIcon} />
                    </IconPlaceholder>
                  )}
                </Ribbon>
                <DetailImageWrapper>
                  <DetailImage img={detailImage} />
                  {detailResult.percentage && (
                    <PercentageRibbon>
                      {detailResult.percentage}
                    </PercentageRibbon>
                  )}
                </DetailImageWrapper>
                <DetailInfo>
                  {Object.keys(detailResult).map(
                    data =>
                      data !== "foto" && (
                        <Row key={data}>
                          <div>{data.split("_").join(" ")}</div>
                          <div>{detailResult[data] || "-"}</div>
                        </Row>
                      )
                  )}
                </DetailInfo>
              </DetailWrapper>
            )}
            {searchResult.length > 0 && (
              <ResultWrapper
                grid={Object.keys(detailResult).length > 0 ? 2 : 5}
              >
                {searchResult.map((result, index) => (
                  <PictureCard
                    key={index}
                    image={`data:image/jpeg;base64,${result.dukcapil_data.foto}`}
                    percentage={`${parseFloat(
                      checkMaxSimilarity(result.probability)
                    ).toFixed(2)}%`}
                    name={result.dukcapil_data.nama_lgkp}
                    onClick={() =>
                      onSelectDetail(
                        result.dukcapil_data,
                        `${parseFloat(
                          checkMaxSimilarity(result.probability)
                        ).toFixed(2)}%`
                      )
                    }
                  />
                ))}
              </ResultWrapper>
            )}
            {errorMsg !== "" && (
              <ErrorWrapper>
                <IconPlaceholder>
                  <img src={ErrorIcon} alt="no-file" />
                </IconPlaceholder>
                {errorMsg}
              </ErrorWrapper>
            )}
          </Wrapper>
        </Right>
      </Wrapper>
    </Wrapper>
  );
}

export default EnvWrapper(SearchHome);

SearchHome.propTypes = {
  history: PropTypes.object.isRequired,
  store: PropTypes.object
};

const Wrapper = Styled.div`
    display: flex;
    flex-direction: ${props => props.flexDirection || `column`};
    width: 100%;
    height: ${props => props.height};
    min-height: ${props => props.minHeight};
    max-height: 100%;
    max-width: 100%;
    overflow-y: auto;
    overflow-x: hidden;    
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
    max-height: 320px;
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
  div {
    padding: 0px 16px;
    text-transform: uppercase;
  }
`;

const DetailImageWrapper = Styled.div`
  padding-top: 32px;
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
`;

const DetailWrapper = Styled.div`
  width: 100%;
  max-height: 100%;
  overflow-y: auto;
  border-right: 1px solid #372463;
`;
const DetailImage = Styled.div`
  padding: 16px;
  width: 100%;
  height: 320px;
  overflow: hidden;
  background-image: url(${props => props.img}), url(${BrokenImage}); 
  background-position: center;
  background-repeat: no-repeat;
  background-size: contain, auto;
`;
const DetailInfo = Styled.div`
  display: grid;
  grid-template-columns: repeat(1, 1fr);
  column-gap: 40px;
  padding: 16px;
  font-size: 13px;  
`;

const ResultWrapper = Styled.div`
  width: ${props => (props.grid === 2 ? `60%` : `100%`)};
  display: grid;
  grid-template-columns: repeat(${props => props.grid}, calc(${props =>
  props.grid === 2 ? `50%` : `20%`} - 15px));
  grid-auto-rows: 250px;
  grid-row-gap: 15px;
  grid-column-gap: 15px;
`;

const ErrorWrapper = Styled.div`
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: ${props => props.theme.inlineError};
  div {
    margin-bottom: 15px;
  }
`;

const PercentageRibbon = Styled.div`
    position: absolute;
    height: 30px;
    display: flex;
    align-items: center;
    justify-content: center;
    bottom: 0;
    background: rgba(69, 69, 69, 0.7);
    width: calc(100% - 8px);
    padding-left: 8px;
    font-weight: 600;
`;

import React, { useState, useEffect, useContext } from "react";
import qs from "qs";
import { useHistory } from "react-router-dom";
import Styled, { ThemeContext } from "styled-components";
import PropTypes from "prop-types";

import IconPlaceHolder from "components/atoms/IconPlaceholder";

import ArrowLeft from "assets/icon/visionaire/arrow-left.svg";
import DoubleLeft from "assets/icon/visionaire/doub-left.svg";

function PaginationWrapper(props) {
  const { children, limit, totalPage, totalData } = props;
  const history = useHistory();
  const themeContext = useContext(ThemeContext);

  const [page, setPage] = useState(1);
  const [queryString, setQueryString] = useState("");

  function getQuery() {
    return qs.stringify({ page, limit }, { encode: false });
  }

  useEffect(() => {
    const getQueryUrl = history.location.search.replace("?", "");
    const queryUrl = qs.parse(getQueryUrl);
    if (Object.keys(queryUrl).length === 0 && queryUrl.constructor === Object) {
      const query = getQuery();
      history.replace(`${history.location.pathname}?${query}`);
    } else {
      const query = getQuery();
      if (query !== queryString) {
        history.replace(`${history.location.pathname}?${query}`);
      }
      setQueryString(query);
    }
  }, [history.location, page]);

  function changePage(pageNum) {
    setPage(pageNum);
  }

  return (
    <Wrapper>
      {children}
      <FooterWrapper>
        <LimitWrapper></LimitWrapper>
        <PageIndicatorWrapper>
          {`${page * limit - limit + 1} - ${
            totalData < page * limit ? totalData : page * limit
          }`}{" "}
          {totalData && `from ${totalData} Data`}
        </PageIndicatorWrapper>
        <ButtonWrapper>
          <IconPlaceHolder
            borderColor={themeContext.secondary2}
            disable={page === 1}
            onClick={() => (page === 1 ? {} : changePage(1))}
          >
            <img src={DoubleLeft} alt="first-page" />
          </IconPlaceHolder>
          <IconPlaceHolder
            borderColor={themeContext.secondary2}
            disable={page === 1}
            onClick={() => (page === 1 ? {} : changePage(page - 1))}
          >
            <img src={ArrowLeft} alt="prev-page" />
          </IconPlaceHolder>
          <IconPlaceHolder
            borderColor={themeContext.secondary2}
            disable={page === totalPage}
            onClick={() => (page === totalPage ? {} : changePage(page + 1))}
            className="reverse"
          >
            <img src={ArrowLeft} alt="next-page" />
          </IconPlaceHolder>
          <IconPlaceHolder
            borderColor={themeContext.secondary2}
            disable={page === totalPage}
            onClick={() => (page === totalPage ? {} : changePage(totalPage))}
            className="reverse"
          >
            <img src={DoubleLeft} alt="last-page" />
          </IconPlaceHolder>
        </ButtonWrapper>
      </FooterWrapper>
    </Wrapper>
  );
}

export default PaginationWrapper;

PaginationWrapper.propTypes = {
  children: PropTypes.element.isRequired,
  limit: PropTypes.number,
  totalPage: PropTypes.number.isRequired,
  totalData: PropTypes.number
};

PaginationWrapper.defaultProps = {
  limit: 10,
  totalData: null
};

const Wrapper = Styled.div`
    position: relative;
    height: 100%;
    width: 100%;
`;

const FooterWrapper = Styled.div`
    position: absolute;
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: space-between;
    bottom: 0;
    left: 0;
    border-top: 1px solid ${props => props.theme.secondary2};
    width: 100%;
    height: 60px;
    color: ${props => props.theme.white};
`;

const LimitWrapper = Styled.div`
    width: 100px;
`;

const PageIndicatorWrapper = Styled.div``;

const ButtonWrapper = Styled.div`
    display: flex;
    flex-direction: row;
    width: 200px;
    justify-content: space-evenly;
    margin-right: 8px;
    .reverse img {
        transform: rotate(180deg);
    }
`;

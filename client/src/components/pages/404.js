import React, { Fragment } from "react";
import PropTypes from "prop-types";
import { createGlobalStyle } from "styled-components";

import Row from "components/atoms/Row";
import Button from "components/molecules/Button";
import Text from "components/atoms/Text";

import GlobeBg from "assets/images/globe.svg";

const GlobalStyle = createGlobalStyle`
  html, body, #root {
    @import url('https://fonts.googleapis.com/css?family=Barlow&display=swap');
    font-family: "Barlow", sans-serif;
    width: 100%;
    height: 100%;
    margin: 0;
    background: #21153C;
  }
  `;

function Page404({ history }) {
  return (
    <Fragment>
      <GlobalStyle />
      <Row
        width="100%"
        height="100vh"
        direction="column"
        align="center"
        justify="center"
        bgUrl={GlobeBg}
      >
        <Row
          direction="column"
          align="center"
          justify="center"
          verticalGap={10}
          width="400px"
          height="168px"
          bgColor="#21153C"
          border="#372463"
          borderRadius={8}
        >
          <Text color="white" size={18} weight={600}>
            404
          </Text>
          <Text color="white" size={18} weight={600}>
            PAGE NOT FOUND
          </Text>
          <Button
            width="80%"
            onClick={() => history.push("/")}
            id="back404"
            style={{ background: "transparent", border: "2px solid #372463" }}
          >
            BACK TO MAIN PAGE
          </Button>
        </Row>
      </Row>
    </Fragment>
  );
}

Page404.propTypes = {
  history: PropTypes.object.isRequired
};

export default Page404;

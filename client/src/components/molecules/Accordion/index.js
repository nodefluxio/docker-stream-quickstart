import React, { useState, useRef, useEffect } from "react";
import Styled from "styled-components";
import PropTypes from "prop-types";

import IconPlaceholder from "components/atoms/IconPlaceholder";
import ChevronIcon from "assets/icon/visionaire/arrow-up.svg";

export default function Accordion(props) {
  const { header, icon, children, type } = props;
  const [isActive, setIsActive] = useState(type !== "closed");
  const [height, setHeight] = useState(0);
  const contentRef = useRef(null);

  function toggleAccordion() {
    setIsActive(!isActive);
  }

  useEffect(() => {
    setHeight(!isActive ? 0 : contentRef.current.scrollHeight);
  }, [isActive, children]);

  return (
    <Wrapper>
      <AccordionButton active={isActive} onClick={() => toggleAccordion()}>
        <Title withIcon={icon}>
          {icon && (
            <IconPlaceholder>
              <img src={icon} alt="accordion-icon" />
            </IconPlaceholder>
          )}
          {header}
        </Title>
        <IconPlaceholder>
          <img
            src={ChevronIcon}
            className="accordion-indicator"
            alt="chevron"
          />
        </IconPlaceholder>
      </AccordionButton>
      <Content ref={contentRef} maxHeight={height}>
        {children}
      </Content>
    </Wrapper>
  );
}

Accordion.propTypes = {
  header: PropTypes.string.isRequired,
  icon: PropTypes.string,
  children: PropTypes.oneOfType([PropTypes.array, PropTypes.element])
    .isRequired,
  type: PropTypes.string
};

Accordion.defaultProps = {
  icon: ""
};

const Wrapper = Styled.div`
    width: 100%;
`;

const AccordionButton = Styled.div`
    width: 100%;
    height: 40px;
    background-color: ${props => props.theme.secondary2};
    cursor: pointer;
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    align-items: center;
    .accordion-indicator {
        transform: ${props =>
          props.active ? `rotate(0deg)` : `rotate(180deg)`};
    }
`;

const Title = Styled.div`
    display: flex;
    flex-direction: row;
    align-items: center;
    font-weight: 700;
    text-transform: ${props => (props.withIcon ? `uppercase` : `capitalize`)} ;
    font-size: 14px;
    line-height: 14.4px;
    color: ${props => (props.withIcon ? props.theme.mint : props.theme.white)};
    ${props => !props.withIcon && `padding-left: 16px;`}
`;

const Content = Styled.div`
    overflow: ${props => (props.maxHeight === 0 ? "hidden" : "visible")};
    max-height: ${props => props.maxHeight}px;
    transition: 0.3s;
    padding: 0px 16px;
    margin: ${props => (props.maxHeight === 0 ? `10px 0px` : `20px 0px`)};
`;

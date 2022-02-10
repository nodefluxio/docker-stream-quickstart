import React from "react";
import Styld from "styled-components";
import PropTypes from "prop-types";

function ListLayout({ columns, className, data, elmAction }) {
  return (
    <Parent className={className}>
      <HeadSection size={columns.length}>
        {columns.map((value, key) => (
          <HeadCell key={`head-column-${key}_${value}`}>
            {value.showName}
          </HeadCell>
        ))}
      </HeadSection>
      {data.map((obj, key) => (
        <RowSection
          key={`${key}-data-section-${obj.name}`}
          size={columns.length}
        >
          {columns.map(
            keyObj =>
              keyObj !== "id" && <RowCell>{obj[keyObj.attr] || "-"}</RowCell>
          )}
          <ActionColumn>{elmAction(obj.id)}</ActionColumn>
        </RowSection>
      ))}
    </Parent>
  );
}
const ActionColumn = Styld.div`
    display:flex;
`;
const HeadCell = Styld.span`
    font-family: Barlow;
    font-style: normal;
    font-weight: 600;
    font-size: 20px;
    line-height: 24px;
    letter-spacing: 0.25px;
    color: ${props => props.theme.mint};
    margin-bottom: 10px;
`;
const RowCell = Styld.div``;
const RowSection = Styld.div`
    display:grid;
    grid-template-columns: repeat(${({ size }) => size + 1}, 1fr);
    justify-items:center;
    align-items: center;
    border-top: 1px solid ${props => props.theme.secondary2};
    padding: 5px 0px;

`;
const HeadSection = Styld.div`
    display: grid;
    grid-template-columns: repeat(${({ size }) => size + 1}, 1fr);
    grid-template-rows: auto;
    padding-bottom: 5px;
    justify-items:center;
`;
const Parent = Styld.div`
    display:flex;
    flex-direction:column;
`;
ListLayout.propTypes = {
  columns: PropTypes.array.isRequired,
  data: PropTypes.array.isRequired,
  className: PropTypes.string,
  elmAction: PropTypes.object
};

export default ListLayout;

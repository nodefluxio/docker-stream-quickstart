import React from "react";
import PropTypes from "prop-types";
import styled from "styled-components";

const WrapTable = styled.div`
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  ${props => {
    if (Number.isInteger(props.padding)) {
      return `padding: ${props.padding};`;
    }
    if (Array.isArray(props.padding)) {
      switch (props.padding.length) {
        case 2:
          return `padding: ${props.padding[0]} ${props.padding[1]};`;
        case 3:
          return `padding: ${props.padding[0]} ${props.padding[1]} ${props.padding[2]};`;
        case 4:
          return `padding: ${props.padding[0]} ${props.padding[1]} ${props.padding[2]} ${props.padding[3]};`;
        default:
          return `padding: auto;`;
      }
    }
    return `padding: auto;`;
  }}
`;

const TableHead = styled.div`
  display: flex;
  width: 100%;
  background-color: ${props => props.theme.secondary2};
`;

const TableBody = styled.div`
  position: relative;
  overflow-x: auto;
  overflow-y: hidden;
  margin-bottom: 10px;
  ${({ bodyHeight }) =>
    bodyHeight &&
    `
    height: ${bodyHeight}; 
    `}
`;

const TableColumn = styled.div`
  flex: 1 1 auto;
  display: flex;
  justify-content: center;
  position: relative;
  align-items: center;
  font-size: 14px;
  font-weight: 600;
  font-style: normal;
  font-stretch: normal;
  line-height: normal;
  letter-spacing: normal;
  text-align: center;
  color: white;
  width: ${props => (props.name === "action" ? "50px" : `100%`)};
  display: flex;
  ${props => {
    if (Number.isInteger(props.padding)) {
      return `padding: ${props.padding};`;
    }
    if (Array.isArray(props.padding)) {
      switch (props.padding.length) {
        case 2:
          return `padding: ${props.padding[0]} ${props.padding[1]};`;
        case 3:
          return `padding: ${props.padding[0]} ${props.padding[1]} ${props.padding[2]};`;
        case 4:
          return `padding: ${props.padding[0]} ${props.padding[1]} ${props.padding[2]} ${props.padding[3]};`;
        default:
          return `padding: 20px 12px;`;
      }
    }
    return `padding: 20px 12px;`;
  }}
`;

const TableWrap = styled.div`
  flex: 1;
  justify-content: center;
  display: flex;
  flex-direction: column;
  background-color: ${props => props.theme.secondary29};
  ${props =>
    props.clickable &&
    `
    cursor: pointer;
    :hover{
      transform: unset;
      background-color: rgba(181, 171, 199, 0.2);
    } 
    `}
`;

const TableWrapContent = styled.div`
  display: flex;
  flex-direction: row;
`;

const TableRow = styled(TableColumn)`
  font-size: 14px;
  font-weight: 500;
  position: relative;
  color: ${props => props.theme.secondary25};
  border-bottom: 1px solid ${props => props.theme.secondary2};
  font-style: normal;
  font-stretch: normal;
  line-height: 21px;
  letter-spacing: normal;
  display: flex;
  flex: 1 1 auto;
  flex-wrap: nowrap;
  -moz-box-pack: justify;
  text-align: center;
  transition-duration: 0.1s;
  margin-bottom: unset;
  margin-top: unset;
  height: ${props => props.height || "47px"};
  ${props => {
    if (Number.isInteger(props.padding)) {
      return `padding: ${props.padding};`;
    }
    if (Array.isArray(props.padding)) {
      switch (props.padding.length) {
        case 2:
          return `padding: ${props.padding[0]} ${props.padding[1]};`;
        case 3:
          return `padding: ${props.padding[0]} ${props.padding[1]} ${props.padding[2]};`;
        case 4:
          return `padding: ${props.padding[0]} ${props.padding[1]} ${props.padding[2]} ${props.padding[3]};`;
        default:
          return `padding: 0;`;
      }
    }
    return `padding: 0;`;
  }};
`;

const TableAlignLeft = styled.div`
  text-align: left;
  width: 100%;
`;

const TableAlignRight = styled.div`
  position: absolute;
  margin: auto 0;
  right: 20px;
`;

const Table = ({
  columns,
  dataSource,
  onRow,
  paddingColumn,
  paddingRow,
  paddingTable,
  bodyHeight,
  rowHeight,
  client
}) => (
  <WrapTable padding={paddingTable}>
    <TableHead>
      {columns.map(val => (
        <TableColumn key={val.key} name={val.key} padding={paddingColumn}>
          {val.title &&
            (typeof val.title === "string" ? val.title : val.title())}
        </TableColumn>
      ))}
    </TableHead>
    <TableBody bodyHeight={bodyHeight}>
      {dataSource &&
        dataSource.map((val, id) => (
          <TableWrap key={val.id} clickable={onRow}>
            <TableWrapContent
              onClick={e => {
                const row = e.target.classList.contains("row");
                if (row) {
                  onRow(e, val, id);
                }
              }}
            >
              {columns.map(column => {
                if (client) {
                  if (val[column.key] == null) {
                    return (
                      <TableRow
                        className="row"
                        padding={paddingRow}
                        height={rowHeight}
                        key={column.key}
                        name={column.key}
                      >
                        {column.render()}
                      </TableRow>
                    );
                  }
                  if (column.render && val[column.key] !== "") {
                    return (
                      <TableRow
                        className="row"
                        padding={paddingRow}
                        height={rowHeight}
                        key={column.key}
                        name={column.key}
                      >
                        {column.render(val[column.key], val)}
                      </TableRow>
                    );
                  }
                }
                return (
                  <TableRow
                    className="row"
                    padding={paddingRow}
                    height={rowHeight}
                    key={column.key}
                    name={column.key}
                  >
                    {val[column.key]}
                  </TableRow>
                );
              })}
            </TableWrapContent>
          </TableWrap>
        ))}
    </TableBody>
  </WrapTable>
);

Table.defaultProps = {
  bodyHeight: "auto",
  rowHeight: "",
  rowState: [],
  client: true,
  onRow: () => {}
};

Table.propTypes = {
  columns: PropTypes.array.isRequired,
  dataSource: PropTypes.array.isRequired,
  client: PropTypes.bool,
  onRow: PropTypes.func,
  rowState: PropTypes.array,
  rowHeight: PropTypes.string,
  bodyHeight: PropTypes.string,
  paddingColumn: PropTypes.oneOfType([PropTypes.number, PropTypes.array]),
  paddingRow: PropTypes.oneOfType([PropTypes.number, PropTypes.array]),
  paddingTable: PropTypes.oneOfType([PropTypes.number, PropTypes.array])
};
export default Table;
export { TableAlignLeft, TableAlignRight };

import Styled from "styled-components";

const Row = Styled.div`
display: flex;
align-content: flex-start;
${({ gap }) =>
  gap &&
  `
    & > div,button{
      margin: ${gap}px ${gap}px;  
    }
  `}
${({ horizontalGap }) =>
  horizontalGap &&
  `
    & > div,button{
      margin: 0 ${horizontalGap}px;  
    }
  `}
${({ verticalGap }) =>
  verticalGap &&
  `
    & > div,button{
      margin: ${verticalGap}px 0;  
    }
  `}
${({ bgUrl }) => bgUrl && `background: url(${bgUrl}) no-repeat center center;`};
${({ onClick }) => onClick && `cursor: pointer;`};
${({ overflow }) => overflow && `overflow: ${overflow};`}
${({ justify }) => justify && `justify-content: ${justify};`}
${({ direction }) => direction && `flex-direction: ${direction};`}
${({ flow }) => flow && `flex-flow: ${flow};`}
${({ align }) => align && `align-items: ${align};`}
${({ wrap }) => wrap && `flex-wrap: ${wrap};`}
${({ width }) => width && `width: ${width};`}
${({ maxWidth }) => maxWidth && `max-width: ${maxWidth};`}
${({ bgColor }) => bgColor && `background-color: ${bgColor};`}
${({ border }) => border && `border: solid 1px ${border};`}
${({ borderWidth }) => borderWidth && `border-width: ${borderWidth}px;`}
${({ borderTop }) => borderTop && `border-top: solid 1px ${borderTop};`}
${({ borderBottom }) =>
  borderBottom && `border-bottom: solid 1px ${borderBottom};`}
${({ borderRight }) => borderRight && `border-right: solid 1px ${borderRight};`}
${({ borderLeft }) => borderLeft && `border-left: solid 1px ${borderLeft};`}
${({ borderRadius }) => borderRadius && `border-radius: ${borderRadius}px;`}
${({ height }) => height && `height: ${height};`}
${({ verticalPadding }) =>
  verticalPadding &&
  `padding-top: ${verticalPadding}px; padding-bottom: ${verticalPadding}px;`}
${({ horizontalPadding }) =>
  horizontalPadding &&
  `padding-right: ${horizontalPadding}px; padding-left: ${horizontalPadding}px;`}
  ${({ verticalMargin }) =>
    verticalMargin &&
    `margin-top: ${verticalMargin}px; margin-bottom: ${verticalMargin}px;`}
${({ horizontalMargin }) =>
  horizontalMargin &&
  `margin-right: ${horizontalMargin}px; margin-left: ${horizontalMargin}px;`}
  ${({ flex }) => flex && `display: ${flex}`}
`;

export default Row;

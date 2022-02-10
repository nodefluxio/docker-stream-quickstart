/* eslint-disable no-console */
import React from "react";
import PropTypes from "prop-types";
import update from "react-addons-update";
import Styled from "styled-components";
import RefreshIcon from "assets/icon/visionaire/refresh.svg";
import PenIcon from "assets/icon/visionaire/pen.svg";
import Remove from "assets/icon/visionaire/remove.svg";

import { colorRoi } from "theme";

const { d3 } = window;
export default class RegionInterest extends React.Component {
  // 1. click the icon draw
  // 2. click the canvas for initial point(x,y)
  // 3. will show the line adjust with mouse position
  // 4. click to creating the line
  constructor(props) {
    super(props);
    this.state = {
      width: props.width,
      height: props.height,
      showUtility: true,
      container: null,
      draw: "no",
      tempLine: null,
      points: [],
      focus: null,
      usedColors: [],
      activeColor: null,
      isFocus: false,
      showRoi: true,
      counterLine: 1,
      arrLine: [],
      areaName: "Area 1",
      image: ""
    };
  }

  componentDidMount() {
    this.initialCanvas();
    this.setState({
      activeColor: this.selectedColor(),
      image: this.props.image
    });
  }

  componentDidUpdate(prevProps) {
    // Check if the suplied props is changed
    if (prevProps.delLine !== this.props.delLine) {
      d3.select(`#line-${this.props.delLine}`).remove();
      d3.select(`#defs-${this.props.delLine}`).remove();
    }

    if (prevProps.reverse !== this.props.reverse) {
      const { line, direction } = this.props.reverse;
      if (direction === "end") {
        d3.select(`#line-${line}`)
          .attr("marker-start", null)
          .attr("marker-end", `url(#arrow-${line})`);
      } else {
        d3.select(`#line-${line}`)
          .attr("marker-start", `url(#arrow-${line})`)
          .attr("marker-end", null);
      }
    }
  }

  generateRandomNumber = (min, max) => {
    const randomNumber = Math.random() * (max - min) + min;
    return Math.floor(randomNumber);
  };

  selectedColor = () => {
    const { usedColors } = this.state;
    const available = colorRoi.filter(color => !usedColors.includes(color));
    const number = this.generateRandomNumber(0, available.length);
    const chosen = available[number];

    this.setState({ usedColors: update(usedColors, { $push: [chosen] }) });
    return chosen;
  };

  initialCanvas = () => {
    const svgContainer = d3
      .select("#canvas")
      .append("svg")
      .attr("width", this.props.width)
      .attr("height", this.props.height);

    //  The SVG Container

    let clicks = 0;
    svgContainer.on("click", () => {
      const coords = d3.mouse(svgContainer.node());
      if (this.state.draw === "initial" || this.state.draw === "drawing") {
        clicks += 1;
        if (clicks === 1) {
          setTimeout(() => {
            if (clicks === 1) {
              this.onLine(coords);
            }
          }, 300);
        }

        if (clicks === 2 && this.state.draw === "drawing") {
          this.onFinish(svgContainer.node(), this.state.counterLine);
          this.setState({ counterLine: this.state.counterLine + 1 });
          clicks = 0;
        }
      }
    });

    svgContainer.on("mousemove", () => {
      if (this.state.draw === "drawing") {
        const coords = d3.mouse(svgContainer.node());
        const line = this.state.tempLine;
        line.attr("x2", coords[0]).attr("y2", coords[1]);
      }
    });

    const { currPoints } = this.props;
    if (currPoints) {
      this.pointsToLine(svgContainer, currPoints);
      this.setState({ showRoi: false });
    }

    this.setState({ container: svgContainer });
  };

  pointsToLine = (container, data) => {
    // eslint-disable-next-line camelcase
    const { points, area_name } = data;
    const color = this.state.activeColor || this.selectedColor();

    if (points) {
      let line = [];

      const pointsToCoords = points.map(d => [
        d.x * this.props.width,
        d.y * this.props.height
      ]);
      pointsToCoords.map((d, idx, elm) => {
        if (elm[idx + 1] && idx !== 1) {
          line = container
            .append("line")
            .attr("x1", d[0])
            .attr("y1", d[1])
            .attr("x2", elm[idx + 1][0])
            .attr("y2", elm[idx + 1][1])
            .attr("stroke-width", 2)
            .attr("stroke", color);
        }
        return line;
      });
      this.setState({ showRoi: true, tempLine: line, areaName: area_name });
    }
  };
  // onLine is function that for starting drawing (onclick the canvas)
  // or add other line connecting

  onLine = coords => {
    const { width, height } = this.props;
    const {
      tempLine,
      container,
      activeColor,
      points,
      counterLine
    } = this.state;
    const svgContainer = container;
    let line = tempLine;
    const nPoints = points;

    if (line) {
      // if this live of multiple line
      line
        .attr("x2", coords[0])
        .attr("y2", coords[1])
        .on("mouseover", () => {
          this.state.focus.style("display", null);
          this.setState({
            isFocus: true
          });
        })
        .on("mouseout", () => {
          this.state.focus.style("display", "none");
          this.setState({
            isFocus: false
          });
        })
        .on("mousemove", () => {
          const moveCoords = d3.mouse(svgContainer.node());
          this.state.focus.attr(
            "transform",
            `translate(${moveCoords[0]},${moveCoords[1]})`
          );
        });
    }

    // this os inserting new line
    svgContainer
      .append("defs")
      .attr("id", `defs-${counterLine}`)
      .append("marker")
      .attr("id", `arrow-${counterLine}`)
      .attr("viewBox", [0, 0, 20, 20])
      .attr("refX", 10)
      .attr("refY", 10)
      .attr("markerWidth", 8)
      .attr("markerHeight", 8)
      .attr("orient", "auto-start-reverse")
      .append("path")
      .attr(
        "d",
        d3.line()([
          [0, 0],
          [0, 20],
          [20, 10]
        ])
      )
      .attr("stroke", activeColor)
      .attr("fill", activeColor);

    line = svgContainer
      .append("line")
      .attr("x1", coords[0])
      .attr("y1", coords[1])
      .attr("x2", coords[0])
      .attr("y2", coords[1])
      .attr("id", `line-${counterLine}`)
      .attr("stroke-width", 2)
      .attr("stroke", activeColor)
      .attr("marker-end", `url(#arrow-${counterLine})`);

    // if first timer
    nPoints.push({ x: coords[0] / width, y: coords[1] / height });
    this.setState({
      tempLine: line,
      points: nPoints,
      draw: "drawing"
    });
  };

  // onFinish is the function when you double click for finish creating line
  onFinish = (container, lineNumber) => {
    const { width, height } = this.props;
    const { points, isFocus, activeColor } = this.state;
    const coords = d3.mouse(container);
    const nPoints = points;
    const isDrawing = this.state.showRoi;

    nPoints.push({ x: coords[0] / width, y: coords[1] / height });

    if (isFocus === false) {
      console.log("lineNumber :", lineNumber);
      this.setState({ activeColor: this.selectedColor() });
      const dataPoint = {};
      nPoints.forEach((value, index) => {
        const pointNumber = index + 1;
        dataPoint[`x${pointNumber}`] = value.x;
        dataPoint[`y${pointNumber}`] = value.y;
      });
      const nArea = `Area ${lineNumber}`;
      this.props.onSave(nPoints, nArea, activeColor, lineNumber);
    }

    this.setState({
      tempLine: null,
      points: [],
      draw: "no",
      showRoi: isDrawing
    });
  };

  onReset = () => {
    d3.selectAll("line").remove();
    d3.selectAll("defs").remove();
    const reset = null;
    this.props.onSave(reset);
    this.setState({ showRoi: true, counterLine: 1 });
    this.props.onReset();
  };

  onRefreshImg = () => {
    const newImg = `${this.props.image}?random=${new Date().getTime()}`;
    this.setState({ image: newImg });
  };

  onDraw = () => {
    this.setState({ showUtility: false, draw: "initial" });
  };

  render() {
    const { isRoiNeed } = this.props;
    const { height, width, draw, image } = this.state;
    return (
      <RoiStyled>
        <Canvas id="canvas" className={`${draw}`} style={{ width, height }}>
          <img src={image} alt="streamer" />
        </Canvas>
        {isRoiNeed && (
          <BtnWrapper>
            <ToolButton onClick={this.onDraw}>
              <img src={PenIcon} alt="icon-draw" />
            </ToolButton>
            <ToolButton onClick={this.onReset} style={{ marginTop: "10px" }}>
              <img src={Remove} alt="refresh" />
            </ToolButton>
            <ToolButton
              onClick={this.onRefreshImg}
              style={{ marginTop: "10px" }}
            >
              <img src={RefreshIcon} alt="refresh" />
            </ToolButton>
          </BtnWrapper>
        )}
      </RoiStyled>
    );
  }
}

RegionInterest.propTypes = {
  width: PropTypes.string,
  height: PropTypes.string,
  onSave: PropTypes.func,
  onReset: PropTypes.func,
  image: PropTypes.string.isRequired,
  isRoiNeed: PropTypes.bool.isRequired,
  currPoints: PropTypes.object,
  roiType: PropTypes.string,
  delLine: PropTypes.string,
  reverse: PropTypes.object
};

RegionInterest.defaultProps = {
  width: "853",
  height: "480",
  currPoints: {},
  roiType: ""
};

const Canvas = Styled.div`
background-color: black;
  position: relative;
  z-index: 10;
  overflow: hidden;
  &.no{
    cursor: default;
  }
  &.initial{
    cursor: cell;
  }
  &.drawing{
    cursor: none;
  }
  &>img{
    position: absolute;
    z-index: -10;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    width: 100%;
  }
  svg{
    border: none;
    user-select: none;
  }
`;

const RoiStyled = Styled.div`
  display: flex;
  position: relative;
`;

const ToolButton = Styled.div`
  width: 35px;
  height: 35px;
  background-color: #439DFF;
  display: flex;
  justify-content: center;
  align-items: center;
  cursor: pointer;
`;

const BtnWrapper = Styled.div`
  /* position: absolute; */
  /* left: 0;
  bottom: 0; */
  display: flex;
  z-index: 99;
  /* grid-template-columns: auto auto;
  grid-gap: 5px; */
  flex-direction: column;
`;

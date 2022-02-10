import React from "react";
import PropTypes from "prop-types";
import update from "react-addons-update";
import styled from "styled-components";

import RefreshIcon from "assets/icon/visionaire/refresh.svg";
import PenIcon from "assets/icon/visionaire/pen.svg";
import Remove from "assets/icon/visionaire/remove.svg";

import { colorRoi } from "theme";

const { d3 } = window;

class RegionInterest extends React.Component {
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
      activeColor: null,
      usedColors: [],
      areaName: "Area 1",
      counterArea: 1
    };
  }

  componentDidMount() {
    this.initialCanvas();
    this.setState({
      activeColor: this.selectedColor()
    });
  }

  componentDidUpdate(prevProps) {
    // Check if the suplied props is changed
    if (prevProps.delLine !== this.props.delLine) {
      d3.selectAll(`.line-${this.props.delLine}`).remove();
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
    // The data for our line
    // This is the accessor function we talked about above
    const svgContainer = d3
      .select("#canvas")
      .append("svg")
      .attr("width", this.props.width)
      .attr("height", this.props.height);
    // The SVG Container
    const $$ = this;
    let clicks = 0;
    svgContainer.on("click", function() {
      const coords = d3.mouse(this);
      clicks += 1;
      if (clicks === 1) {
        setTimeout(() => {
          if (
            clicks === 1 &&
            ($$.state.draw === "initial" || $$.state.draw === "drawing")
          ) {
            $$.onLine(coords);
          }
          clicks = 0;
        }, 300);
      }
      // let coords = d3.mouse(this.state.container[0]);
    });
    svgContainer.on("dblclick", function() {
      if ($$.state.draw === "drawing") {
        $$.onFinish(this);
      }
    });
    svgContainer.on("mousemove", function() {
      if ($$.state.draw === "drawing") {
        const coords = d3.mouse(this);
        const line = $$.state.tempLine;
        line.attr("x2", coords[0]).attr("y2", coords[1]);
      }
      // let coords = d3.mouse(this.state.container[0]);
    });
    const focus = svgContainer
      .append("g")
      .attr("class", "focus")
      .style("display", "none");
    focus.append("circle").attr("r", 6.0);
    this.setState({ container: svgContainer, focus });
  };

  // onLine is function that for starting drawing (onclick the canvas)
  // or add other line connecting

  onLine = coords => {
    const svgContainer = this.state.container;
    let line = this.state.tempLine;
    const { counterArea } = this.state;
    const $$ = this;
    if (line != null) {
      // if this live of multiple line
      line
        .attr("class", `line-${counterArea}`)
        .attr("x2", coords[0])
        .attr("y2", coords[1])
        .on("mouseover", () => {
          $$.state.focus.style("display", null);
        })
        .on("mouseout", () => {
          $$.state.focus.style("display", "none");
        })
        .on("mousemove", function() {
          const coordsMove = d3.mouse(this);
          $$.state.focus.attr(
            "transform",
            `translate(${coordsMove[0]},${coordsMove[1]})`
          );
        });
    }
    // this os inserting new line
    line = svgContainer
      .append("line")
      .attr("class", `line-${counterArea}`)
      .attr("x1", coords[0])
      .attr("y1", coords[1])
      .attr("x2", coords[0])
      .attr("y2", coords[1])
      .attr("stroke-width", 2)
      .attr("stroke", `${this.state.activeColor}`);
    // if first timer
    let { points } = this.state;
    if (this.state.draw === "initial") {
      points = update(this.state.points, {
        $push: [
          [
            {
              x: coords[0] / this.props.width,
              y: coords[1] / this.props.height
            }
          ]
        ]
      });
    } else {
      const temp = update(points[this.state.points.length - 1], {
        $push: [
          { x: coords[0] / this.props.width, y: coords[1] / this.props.height }
        ]
      });
      points[this.state.points.length - 1] = temp;
    }
    this.setState({ tempLine: line, points, draw: "drawing" });
  };

  // onFinish is the function when you double click for finish creating line
  onFinish = container => {
    const coords = d3.mouse(container);
    const { points, activeColor, counterArea } = this.state;
    const temp = update(points[this.state.points.length - 1], {
      $push: [
        { x: coords[0] / this.props.width, y: coords[1] / this.props.height }
      ]
    });
    points[this.state.points.length - 1] = temp;
    const nArea = `Area ${counterArea}`;
    this.props.onSave(
      points[this.state.points.length - 1],
      nArea,
      activeColor,
      counterArea
    );
    this.setState({
      tempLine: null,
      points: [],
      draw: "no",
      activeColor: this.selectedColor(),
      counterArea: counterArea + 1
    });
  };

  onDraw = () => {
    this.setState({ showUtility: false, draw: "initial" });
  };

  onReset = () => {
    d3.selectAll("line").remove();
    d3.selectAll("defs").remove();
    const reset = null;
    this.props.onSave(reset);
    this.setState({ showRoi: true, counterArea: 1 });
    this.props.onReset();
  };

  render() {
    return (
      <RoiStyled>
        <Canvas
          id="canvas"
          className={`${this.state.draw}`}
          style={{ width: this.state.width, height: this.state.height }}
        >
          <img
            alt="camera-preview"
            src={this.props.image}
            style={{
              width: `${this.state.width}px`,
              height: `${this.state.height}px`
            }}
          />
        </Canvas>
        <BtnWrapper>
          <ToolButton onClick={this.onDraw}>
            <img src={PenIcon} alt="icon-draw" />
          </ToolButton>
          <ToolButton onClick={this.onReset} style={{ marginTop: "10px" }}>
            <img src={Remove} alt="refresh" />
          </ToolButton>
          <ToolButton onClick={this.onRefreshImg} style={{ marginTop: "10px" }}>
            <img src={RefreshIcon} alt="refresh" />
          </ToolButton>
        </BtnWrapper>
      </RoiStyled>
    );
  }
}

const Canvas = styled.div`
  display: flex;
  z-index: 10;
  display: flex;
  background-color: black;
  position: relative;
  z-index: 10;
  overflow: hidden;
  &.no {
    cursor: default;
  }
  &.initial {
    cursor: cell;
  }
  &.drawing {
    cursor: none;
  }
  & > img {
    position: absolute;
    z-index: -10;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    width: 100%;
  }
  svg {
    border: none;
    user-select: none;
  }
`;

const BtnWrapper = styled.div`
  display: flex;
  z-index: 99;
  grid-template-columns: auto auto;
  grid-gap: 5px;
  flex-direction: column;
`;

const RoiStyled = styled.div`
  display: flex;
  position: relative;
`;

const ToolButton = styled.div`
  width: 35px;
  height: 35px;
  background-color: #439dff;
  display: flex;
  justify-content: center;
  align-items: center;
  cursor: pointer;
`;

RegionInterest.propTypes = {
  width: PropTypes.number,
  height: PropTypes.number,
  onSave: PropTypes.func,
  onReset: PropTypes.func,
  image: PropTypes.string.isRequired,
  isRoiNeed: PropTypes.bool.isRequired,
  delLine: PropTypes.number
};

RegionInterest.defaultProps = {
  width: 853,
  height: 480
};

export default RegionInterest;

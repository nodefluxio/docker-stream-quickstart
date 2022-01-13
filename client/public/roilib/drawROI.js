/* eslint-disable no-unused-vars */
const { d3 } = window;

const colorRoi = [
  "#e6194b",
  "#3cb44b",
  "#ffe119",
  "#4363d8",
  "#f58231",
  "#911eb4",
  "#46f0f0",
  "#f032e6",
  "#bcf60c",
  "#fabebe",
  "#008080",
  "#e6beff",
  "#9a6324",
  "#fffac8",
  "#800000",
  "#aaffc3",
  "#808000",
  "#ffd8b1",
  "#000075",
  "#808080",
  "#000000"
];

let container = null;
let tempLine = null;
let points = [];
let focus = null;
const usedColors = [];
let activeColor = null;
let isFocus = false;
let counter = 1;
let clicks = 0;

// onLine is function that for starting drawing (onclick the canvas)
// or add other line connecting

const onLine = (coords, width, height) => {
  const svgContainer = container;
  let line = tempLine;
  const nPoints = points;

  if (line) {
    // if this live of multiple line
    line
      .attr("x2", coords[0])
      .attr("y2", coords[1])
      .on("mouseover", () => {
        focus = focus.style("display", null);
        isFocus = true;
      })
      .on("mouseout", () => {
        focus = focus.style("display", "none");
        isFocus = false;
      })
      .on("mousemove", () => {
        const moveCoords = d3.mouse(svgContainer.node());
        focus = focus.attr(
          "transform",
          `translate(${moveCoords[0]},${moveCoords[1]})`
        );
      });
  }

  // this os inserting new line
  svgContainer
    .append("defs")
    .attr("id", `defs-${counter}`)
    .append("marker")
    .attr("id", `arrow-${counter}`)
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
    .attr("id", `line-${counter}`)
    .attr("stroke-width", 2)
    .attr("stroke", activeColor)
    .attr("marker-end", `url(#arrow-${counter})`);

  // if first timer
  nPoints.push({ x: coords[0] / width, y: coords[1] / height });
  tempLine = line;
  points = nPoints;
  document.getElementById("canvas").className = "drawing";
};

const onLineRegion = (coords, width, height) => {
  const svgContainer = container;
  let line = tempLine;

  if (line != null) {
    // if this live of multiple line
    line
      .attr("class", `line-${counter}`)
      .attr("x2", coords[0])
      .attr("y2", coords[1])
      .on("mouseover", () => {
        focus.style("display", null);
      })
      .on("mouseout", () => {
        focus.style("display", "none");
      })
      .on("mousemove", function() {
        const coordsMove = d3.mouse(this);
        focus.attr("transform", `translate(${coordsMove[0]},${coordsMove[1]})`);
      });
  }
  // this os inserting new line
  line = svgContainer
    .append("line")
    .attr("class", `line-${counter}`)
    .attr("x1", coords[0])
    .attr("y1", coords[1])
    .attr("x2", coords[0])
    .attr("y2", coords[1])
    .attr("stroke-width", 2)
    .attr("stroke", `${activeColor}`);
  // if first timer
  points.push({
    x: coords[0] / width,
    y: coords[1] / height
  });

  tempLine = line;
  document.getElementById("canvas").className = "drawing";
};

const generateRandomNumber = (min, max) => {
  const randomNumber = Math.random() * (max - min) + min;
  return Math.floor(randomNumber);
};

const selectedColor = () => {
  const available = colorRoi.filter(color => !usedColors.includes(color));
  const number = generateRandomNumber(0, available.length);
  const chosen = available[number];
  usedColors.push(chosen);
  return chosen;
};

// onFinish is the function when you double click for finish creating line
// eslint-disable-next-line consistent-return
function onFinish(callback, containerParam, lineNumber, width, height) {
  const coords = d3.mouse(containerParam);
  const nPoints = points;

  nPoints.push({ x: coords[0] / width, y: coords[1] / height });

  tempLine = null;
  points = [];
  document.getElementById("canvas").className = "no";

  if (isFocus === false) {
    activeColor = selectedColor();
    const dataPoint = {};
    nPoints.forEach((value, index) => {
      const pointNumber = index + 1;
      dataPoint[`x${pointNumber}`] = value.x;
      dataPoint[`y${pointNumber}`] = value.y;
    });
    const nArea = `Area ${lineNumber}`;
    callback({
      points: nPoints,
      area: nArea,
      color: activeColor,
      lineNumber
    });
  }
}

function initialCanvas(callback, type) {
  const canvas = document.querySelector("#canvas");
  const width = canvas.offsetWidth;
  const height = canvas.offsetHeight;
  counter = 1;
  activeColor = selectedColor();
  let svgContainer = d3.select("#canvas");
  svgContainer.selectAll("svg").remove();
  svgContainer = d3
    .select("#canvas")
    .append("svg")
    .attr("width", width)
    .attr("height", height);

  svgContainer.on("dblclick", () => {
    const draw = document.getElementById("canvas").className;
    if (draw === "drawing") {
      onFinish(callback, svgContainer.node(), counter, width, height);
      counter += 1;
    }
  });

  //  The SVG Container
  svgContainer.on("click", () => {
    const draw = document.getElementById("canvas").className;
    const coords = d3.mouse(svgContainer.node());
    clicks += 1;

    setTimeout(() => {
      if (clicks === 1 && (draw === "initial" || draw === "drawing")) {
        if (type === "line") {
          onLine(coords, width, height);
        } else if (type === "region") {
          onLineRegion(coords, width, height);
        }
      }
      if (type === "region") {
        clicks = 0;
      }
    }, 300);

    if (clicks === 2 && draw === "drawing" && type === "line") {
      onFinish(callback, svgContainer.node(), counter, width, height);
      counter += 1;
      clicks = 0;
    }
  });

  svgContainer.on("mousemove", () => {
    const draw = document.getElementById("canvas").className;
    if (draw === "drawing") {
      const coords = d3.mouse(svgContainer.node());
      const line = tempLine;
      line.attr("x2", coords[0]).attr("y2", coords[1]);
    }
  });

  if (type === "region") {
    focus = svgContainer
      .append("g")
      .attr("class", "focus")
      .style("display", "none");
    focus.append("circle").attr("r", 6.0);
  }

  container = svgContainer;
}

function onReverseLine(id) {
  const direction = d3.select(`#line-${id}`).attr("marker-start");
  if (direction !== null) {
    d3.select(`#line-${id}`)
      .attr("marker-start", null)
      .attr("marker-end", `url(#arrow-${id})`);
  } else {
    d3.select(`#line-${id}`)
      .attr("marker-start", `url(#arrow-${id})`)
      .attr("marker-end", null);
  }
}

function onDeleteRegion(id) {
  d3.selectAll(`.line-${id}`).remove();
}

function onDeleteLine(id) {
  d3.select(`#line-${id}`).remove();
  d3.select(`#defs-${id}`).remove();
}

function onResetROI(callback) {
  d3.selectAll("line").remove();
  d3.selectAll("defs").remove();
  counter = 1;
  callback();
}

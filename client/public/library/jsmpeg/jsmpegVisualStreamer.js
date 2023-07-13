/* eslint-disable no-restricted-globals */
/* eslint-disable func-names */
/* eslint-disable no-console */
/* eslint-disable no-unused-vars */
const GLOBAL_CONTROLLER = [];
const MESSAGE_NETWORK_ERROR = "There is a connection issue for this stream";
const DEFAULT_MAX_RETRY = 3;

let hidden;
let visibilityChange;
if (typeof document.hidden !== "undefined") {
  // Opera 12.10 and Firefox 18 and later support
  hidden = "hidden";
  visibilityChange = "visibilitychange";
} else if (typeof document.msHidden !== "undefined") {
  hidden = "msHidden";
  visibilityChange = "msvisibilitychange";
} else if (typeof document.webkitHidden !== "undefined") {
  hidden = "webkitHidden";
  visibilityChange = "webkitvisibilitychange";
}

// eslint-disable-next-line func-names
const WSSource = function(url, options) {
  this.url = url;
  this.options = options;
  this.socket = null;
  this.streaming = true;

  this.callbacks = { connect: [], data: [] };
  this.destination = null;

  this.reconnectInterval =
    options.reconnectInterval !== undefined ? options.reconnectInterval : 5;
  this.shouldAttemptReconnect = !!this.reconnectInterval;

  this.completed = false;
  this.established = false;
  this.progress = 0;
  this.reconnectTimeoutId = 0;
  this.maxRetries = options.maxRetry || DEFAULT_MAX_RETRY;
  this.counterRetry = 1;

  this.onEstablishedCallback = options.onSourceEstablished;
  this.onCompletedCallback = options.onSourceCompleted; // Never used

  this.svgCanvas = document.createElement("canvas");
  this.svgCanvas.id = `base-cnv-${options.id}`;
  this.svgCanvas.width = this.options.canvas.width;
  this.svgCanvas.height = this.options.canvas.height;

  this.options.canvas.parentNode.insertBefore(
    this.svgCanvas,
    this.options.canvas
  );

  this.svgCanvas.style.zIndex = 2;
  this.svgCanvas.style.position = "absolute";
  this.svgCanvas.style.width = "100%";
  this.svgCanvas.style.height = "100%";
  this.svgCanvas.style.opacity = 0;
  this.svgCanvas.style.backgroundColor = "transparent";

  this.options.canvas.style.zIndex = 0;
  this.options.canvas.style.position = "absolute";
  this.options.canvas.style.width = "100%";
  this.options.canvas.style.height = "100%";
  this.options.canvas.style.opacity = 0;
  this.options.canvas.style.backgroundColor = "transparent";

  this.svgCtx = this.svgCanvas.getContext("2d");
  this.svgImg = new Image();

  this.messageData = "";
  this.messageDataError = false;
  this.additionalData = "";
};

WSSource.prototype.connect = function(destination) {
  this.destination = destination;
};

WSSource.prototype.destroy = function() {
  clearTimeout(this.reconnectTimeoutId);
  this.shouldAttemptReconnect = false;
  this.socket.close();
};

WSSource.prototype.start = function() {
  const { id, loading, protocols } = this.options;

  if (this.counterRetry < this.maxRetries) {
    loading.start(id);
  }
  this.shouldAttemptReconnect = !!this.reconnectInterval;
  this.progress = 0;
  this.established = false;

  if (protocols) {
    this.socket = new WebSocket(this.url, protocols);
  } else {
    this.socket = new WebSocket(this.url);
  }
  GLOBAL_CONTROLLER[id] = this.socket;
  this.socket.binaryType = "arraybuffer";
  this.socket.onmessage = this.onMessage.bind(this);
  this.socket.onopen = this.onOpen.bind(this);
  this.socket.onerror = this.onClose.bind(this);
  this.socket.onclose = this.onClose.bind(this);
  this.socket.destroy = this.destroy.bind(this);
  this.socket.resume = this.resume.bind(this);
};

WSSource.prototype.resume = function() {
  this.svgCtx.clearRect(
    0,
    0,
    this.svgCtx.canvas.width,
    this.svgCtx.canvas.height
  );
  this.shouldAttemptReconnect = true;
  this.start();
};

WSSource.prototype.onOpen = function() {
  this.progress = 1;
};

WSSource.prototype.onClose = function() {
  const { id, message, pastErrorMsg } = this.options;
  if (
    pastErrorMsg[id] !== MESSAGE_NETWORK_ERROR &&
    this.counterRetry >= this.maxRetries
  ) {
    message.error(id, MESSAGE_NETWORK_ERROR);
  }
  this.svgCtx.clearRect(
    0,
    0,
    this.svgCtx.canvas.width,
    this.svgCtx.canvas.height
  );
  if (this.shouldAttemptReconnect) {
    clearTimeout(this.reconnectTimeoutId);
    this.reconnectTimeoutId = setTimeout(() => {
      this.start();
      this.counterRetry += 1;
    }, this.reconnectInterval * 1000);
  }
};

WSSource.prototype.handleVisualData = function(buff) {
  const st = new TextDecoder().decode(buff);

  // additional check to ensure canvas and svgCanvas size sync,
  // we put it here, because this routine will be called periodically
  this.svgImg.src = `data:image/svg+xml;base64,${st}`;
  if (this.svgCanvas.height !== this.options.canvas.height)
    this.svgCanvas.height = this.options.canvas.height;

  const ctx = this.svgCtx;
  const newSvgImg = this.svgImg;
  this.svgImg.onload = function() {
    if (ctx.canvas.width !== newSvgImg.naturalWidth) {
      ctx.canvas.width = newSvgImg.width;
    }
    if (newSvgImg.naturalWidth * newSvgImg.naturalHeight) {
      ctx.clearRect(0, 0, ctx.canvas.width, ctx.canvas.height);
      ctx.drawImage(newSvgImg, 0, 0);
    }
  };
};

WSSource.prototype.handleAdditionalData = function(buff) {
  this.additionalData = new TextDecoder().decode(buff);
  // draw additional data here
};

WSSource.prototype.handleMessageData = function(buff) {
  const { id, message, pastErrorMsg } = this.options;
  this.messageData = new TextDecoder().decode(buff.slice(1, buff.length));
  if (buff[0] === 0x01) {
    this.messageDataError = true;
    if (pastErrorMsg[id] !== this.messageData) {
      message.error(id, this.messageData);
    }
  } else {
    this.messageDataError = false;
    pastErrorMsg[id] = "";
    message.clear();
  }

  // additional task on checking whether canvas already resized
  if (this.svgCanvas.height !== this.options.canvas.height) {
    this.svgCanvas.height = this.options.canvas.height;
    this.svgCtx.clearRect(
      0,
      0,
      this.svgCtx.canvas.width,
      this.svgCtx.canvas.height
    );
    this.svgCtx.drawImage(this.svgImg, 0, 0);
  }
};

WSSource.prototype.onMessage = function(ev) {
  const { id, loading, loadingStatus, message, pastErrorMsg } = this.options;
  const isFirstChunk = !this.established;
  this.established = true;

  if (isFirstChunk && this.onEstablishedCallback) {
    this.onEstablishedCallback(this);
  }

  if (this.destination) {
    if (loadingStatus[id] === "start" && this.established) {
      loading.complete(id);
      this.options.canvas.style.opacity = 1;
      this.svgCanvas.style.opacity = 1;
    }
    const buff = new Uint8Array(ev.data);
    if (
      buff[0] === 0xbe &&
      buff[1] === 0xef &&
      buff[2] === 0xca &&
      buff[3] === 0xfe
    ) {
      pastErrorMsg[id] = "";
      message.clear();
      this.handleVisualData(buff.slice(4, buff.length));
    } else if (
      buff[0] === 0xbe &&
      buff[1] === 0xef &&
      buff[2] === 0xba &&
      buff[3] === 0xbe
    ) {
      pastErrorMsg[id] = "";
      this.handleAdditionalData(buff.slice(4, buff.length));
    } else if (
      buff[0] === 0xbe &&
      buff[1] === 0xef &&
      buff[2] === 0xc0 &&
      buff[3] === 0xde
    )
      this.handleMessageData(buff.slice(4, buff.length));
    else {
      this.destination.write(ev.data);
    }
  }
};

function stopVisualisation() {
  for (let i = 0; i < GLOBAL_CONTROLLER.length; i += 1) {
    GLOBAL_CONTROLLER[i].destroy();
  }
}

function handleVisibilityChange() {
  if (document[hidden]) {
    stopVisualisation();
  } else {
    for (let i = 0; i < GLOBAL_CONTROLLER.length; i += 1) {
      GLOBAL_CONTROLLER[i].resume();
    }
  }
}

function renderVisualisation(container) {
  const pastErrorMsg = new Array(container.length);
  const loadingStatus = new Array(container.length);

  for (let i = 0; i < container.length; i += 1) {
    const { url, maxRetry } = container[i].dataset;
    const visualContainer = document.createElement("div");
    const loadingElem = document.createElement("div");
    const errorElem = document.createElement("div");
    const canvas = document.createElement("canvas");

    visualContainer.setAttribute("class", `visualization-container-${i}`);
    visualContainer.setAttribute(
      "style",
      "position: relative; width: 100%; height: 100%; background-color: #000000;"
    );

    loadingElem.setAttribute("class", `visualization-loading-${i}`);
    loadingElem.setAttribute(
      "style",
      "display: inline-block; width: 80px; height: 80px; position: absolute; top: 50%; left: 50%; z-index: 3; transform: translateX(-50%) translateY(-50%);"
    );
    const styleElem = visualContainer.appendChild(
      document.createElement("style")
    );
    styleElem.innerHTML =
      `.visualization-loading-${i}:after {` +
      `content: ' ';` +
      `display: block;` +
      `width: 64px;` +
      `height: 64px;` +
      `margin: 8px;` +
      `border-radius: 50%;` +
      `border: 6px solid #fff;` +
      `border-color: #fff transparent #fff transparent;` +
      `animation: dual-ring 1.2s linear infinite;` +
      `}` +
      `@-webkit-keyframes dual-ring{` +
      `0% { transform: rotate(0deg); }` +
      `100% { transform: rotate(360deg); }` +
      `}`;

    container[i].appendChild(visualContainer);
    visualContainer.appendChild(canvas);
    canvas.id = `cnv-${i}`;

    const loading = {
      start(id) {
        loadingStatus[id] = "start";
        visualContainer.style.backgroundColor = "rgba(0, 0, 0, 0.5)";
        visualContainer.appendChild(loadingElem);
        errorElem.remove();
      },
      complete(id) {
        loadingStatus[id] = "complete";
        visualContainer.style.backgroundColor = "";
        loadingElem.remove();
        errorElem.remove();
      }
    };

    const message = {
      clear() {
        loadingElem.remove();
        errorElem.remove();
      },
      error(id, msg) {
        const stopped = msg.match(/(Stream is stopped)/g);
        if (stopped !== null) {
          container[id].remove();
        } else {
          visualContainer.style.backgroundColor = "rgba(0, 0, 0, 0.5)";
          errorElem.id = `error-${i}`;
          errorElem.setAttribute(
            "style",
            `display: 
              inline-block; position: absolute; 
              top: 40%; 
              left: 50%; 
              z-index: 3; 
              transform: 
              translateX(-50%) translateY(-50%); 
              color: black; 
              border-style: solid;
              padding: 20px;
              width: 80%;
              background-color: #ffffff9e;`
          );
          if (pastErrorMsg[id] !== msg) {
            errorElem.innerHTML = msg;
            pastErrorMsg[id] = msg;
          }
          loadingElem.remove();
          visualContainer.style.backgroundColor = "rgba(0, 0, 0, 0.5)";
          visualContainer.appendChild(errorElem);
        }
      }
    };

    // eslint-disable-next-line no-undef
    const player = new JSMpeg.Player(url, {
      id: i,
      source: WSSource,
      maxRetry,
      canvas,
      audio: false,
      loading,
      loadingStatus,
      message,
      pastErrorMsg
    });
  }
}

async function showVisualisation(reload = false) {
  const container = await document.getElementsByClassName("nodeflux-streamer");
  if (
    typeof document.addEventListener === "undefined" ||
    hidden === undefined
  ) {
    // eslint-disable-next-line no-console
    console.log(
      "This library requires a browser, such as Google Chrome or Firefox, that supports the Page Visibility API."
    );
  } else {
    // Handle page visibility change
    document.addEventListener(visibilityChange, handleVisibilityChange, false);
  }
  if (container.length === 1) {
    const child = container[0].children;
    for (let i = 0; i < child.length; i += 1) {
      const ToRemove = child[i].className.includes("visualization-container");
      if (ToRemove) {
        child[i].remove();
      }
    }
    renderVisualisation(container);
  } else if (container.length > 1) {
    if (reload) {
      // eslint-disable-next-line no-restricted-globals
      location.reload();
    } else {
      renderVisualisation(container);
    }
  }
}

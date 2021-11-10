const SOI = new Uint8Array(2);
SOI[0] = 0xFF;
SOI[1] = 0xD8;
const CONTENT_LENGTH = 'Content-Length';
const VISUAL_DATA = 'X-NF-Visual-Data';
const VISUAL_DATA_TS = 'X-NF-Visual-Data-Ts';
const TYPE_JPEG = 'image/jpeg';
const MESSAGE_DATA = "X-NF-Message-Data";
const ADDITIONAL_DATA = "X-NF-Additional-Data";
let heatmapInstance;
document.addEventListener("DOMContentLoaded", function (event) {
    let divs = document.getElementsByClassName("nodeflux-streamer");
    for (var i = 0; i < divs.length; i++) {
        // append div and img element + styling
        const container = document.createElement("div");
        container.setAttribute("class", "visualization-container-" + i)

        let video = document.createElement("canvas");
        let ctx = video.getContext("2d");
        sampling_fps = 50;

        let loadingElem = document.createElement("div");
        let errorElem = document.createElement("div");
        container.setAttribute("style", "position: relative; width: 100%; height: 100%");
        video.style.zIndex = 0;
        video.style.position = "absolute";
        loadingElem.setAttribute("class", "visualization-loading-" + i);
        loadingElem.setAttribute("style", "display: inline-block; width: 80px; height: 80px; position: absolute; top: 50%; left: 50%; z-index: 3; transform: translateX(-50%) translateY(-50%);");
        let styleElem = container.appendChild(document.createElement("style"));
        styleElem.innerHTML = ".visualization-loading-" + i + ":after {" +
            "content: ' ';" +
            "display: block;" +
            "width: 64px;" +
            "height: 64px;" +
            "margin: 8px;" +
            "border-radius: 50%;" +
            "border: 6px solid #fff;" +
            "border-color: #fff transparent #fff transparent;" +
            "animation: dual-ring 1.2s linear infinite;" +
            "}" +
            "@-webkit-keyframes dual-ring{" +
            "0% { transform: rotate(0deg); }" +
            "100% { transform: rotate(360deg); }" +
            "}";
        let loading = {
            start: function () {
                console.log("start");
                container.style.backgroundColor = "rgba(0, 0, 0, 0.5)";
                container.appendChild(loadingElem);
            },
            complete: function () {
                container.style.backgroundColor = "";
                loadingElem.remove();
            },
            error: function (msg) {
                container.style.backgroundColor = "rgba(0, 0, 0, 0.5)";
                errorElem.setAttribute("style", "display: inline-block; position: absolute; top: 40%; left: 50%; z-index: 3; transform: translateX(-50%) translateY(-50%); color: white");
                errorElem.appendChild(document.createTextNode(msg));
                loadingElem.remove();
                container.style.backgroundColor = "rgba(0, 0, 0, 0.5)";
                container.appendChild(errorElem);
            }
        };
        divs[i].appendChild(container)
        container.appendChild(video);
        loading.start();

        // is there any atrribute data-url?
        const url = divs[i].dataset.url;
        if (url.includes("CE")) {
            heatmapInstance = window.h337.create({
              // only container is required, the rest will be defaults
              container: container,
              blur: 1
            });
          }
        if (url !== "") {
            fetch(url)
                .then(response => {
                    console.log(response);
                    if (!response.ok) {
                        throw Error(response.status + ' ' + response.statusText)
                    }

                    if (!response.body) {
                        throw Error('ReadableStream not yet supported in this browser.')
                    }

                    const reader = response.body.getReader();

                    let headers = '';
                    let contentLength = -1;
                    let imageBuffer = null;
                    let activeImageBuffer = null;
                    let jpeg_url = '';
                    let active_jpeg_img = null;
                    let active_visual_url = '';
                    let active_visual_img = null;
                    let bytesRead = 0;
                    let pendingLoad = false;


                    // calculating fps. This is pretty lame. Should probably implement a floating window function.
                    let frames = 0;
                    let last_ts = 0;

                    setInterval(() => {
                        console.log("fps : " + frames);
                        frames = 0;
                    }, 1000)

                    setInterval(() => {
                        if (activeImageBuffer != null && !pendingLoad) {
                            URL.revokeObjectURL(jpeg_url);
                            jpeg_url = URL.createObjectURL(new Blob([activeImageBuffer], { type: TYPE_JPEG })); //24fps

                            let jpeg_img = new Image;
                            jpeg_img.src = jpeg_url;
                            pendingLoad = true;

                            jpeg_img.onload = function () {
                                if (jpeg_img.naturalWidth * jpeg_img.naturalHeight) {
                                    active_jpeg_img = jpeg_img;
                                    jpeg_img = null;
                                }
                                pendingLoad = false;
                            };

                            jpeg_img.onerror = function () {
                                pendingLoad = false;
                            };
                            activeImageBuffer = null;

                            if (active_visual_url) {
                                let svg_img = new Image;
                                svg_img.src = active_visual_url;
                                svg_img.onload = function () {
                                    if (svg_img.naturalWidth * svg_img.naturalHeight) {
                                        active_visual_img = svg_img;
                                    }
                                    svg_img = null;
                                }
                                active_visual_url = '';
                            }
                        }
                    }, 1000 / sampling_fps)

                    function animate() {
                        requestAnimationFrame(animate);
                        if (active_jpeg_img != null) {
                            ctx.canvas.width = active_jpeg_img.width;
                            ctx.canvas.height = active_jpeg_img.height;
                            ctx.drawImage(active_jpeg_img, 0, 0);
                            if (active_visual_img != null)
                                ctx.drawImage(active_visual_img, 0, 0);
                            active_jpeg_img = null;
                        }
                    }
                    requestAnimationFrame(animate);


                    const read = () => {

                        reader.read().then(({ done, value }) => {
                            if (done) {
                                controller.close();
                                return;
                            }

                            for (let index = 0; index < value.length; index++) {

                                // we've found start of the frame. Everything we've read till now is the header.
                                if (value[index] === SOI[0] && value[index + 1] === SOI[1]) {
                                    // console.log('header found : ' + newHeader);
                                    contentLength = getLength(headers);
                                    // console.log("Content Length : " + newContentLength);
                                    imageBuffer = new Uint8Array(new ArrayBuffer(contentLength));
                                }
                                // we're still reading the header.
                                if (contentLength <= 0 || contentLength === -1) {
                                    headers += String.fromCharCode(value[index]);
                                }
                                // we're now reading the jpeg. 
                                else if (bytesRead < contentLength) {
                                    imageBuffer[bytesRead++] = value[index];
                                }
                                // we're done reading the jpeg. Time to render it. 
                                else {
                                    const raw = getAdditionalData(headers);
                                    let data = null;
                                    if (typeof raw === "object") {
                                        const rawData = raw.data;
                                        const pointData = [];
                                        for (let y = 0; y < rawData.data.length; y += 1) {
                                            for (let x = 0; x < rawData.data[y].length; x += 1) {
                                                const width =
                                                    heatmapInstance._config.container.clientWidth;
                                                const height =
                                                    heatmapInstance._config.container.clientHeight;
                                                const gridWidth = Number.parseFloat(width / rawData.data[y].length).toFixed(1);
                                                const gridHeight = Number.parseFloat(height / rawData.data.length).toFixed(1);
                                                if (rawData.data[y][x] > 0) {
                                                    pointData.push({
                                                        x: gridWidth * x,
                                                        y: gridHeight * y,
                                                        value: rawData.data[y][x]
                                                    });
                                                }
                                            }
                                        }
                                        data = {
                                            min: 0,
                                            max: 255,
                                            data: pointData
                                        };
                                    }
                                    loading.complete();
                                    activeImageBuffer = imageBuffer;
                                    if (data !== null && heatmapInstance !== null) {
                                        heatmapInstance.setData(data);
                                    }

                                    imageBuffer = null;
                                    if (last_ts < getVisualDataTs(headers) && getVisualData(headers).length) {
                                        active_visual_url = 'data:image/svg+xml;base64,' + getVisualData(headers); //12FPS
                                        last_ts = getVisualDataTs(headers);
                                    }
                                    frames++;
                                    contentLength = 0;
                                    bytesRead = 0;
                                    headers = '';
                                }

                                if (index === value.length - 1 && bytesRead >= contentLength) {
                                    const errMessage = getMessageData(headers);
                                    if (errMessage !== -1) {
                                        loading.error(errMessage);
                                    }
                                }
                            }

                            read();
                        }).catch(error => {
                            loading.error(error);
                        })
                    }

                    read();

                }).catch(error => {
                    loading.error(error);
                })
        }
    }
});


const getLength = (headers) => {
    let contentLength = -1;
    headers.split('\n').forEach((header, _) => {
        const pair = header.split(':');
        if (pair[0] === CONTENT_LENGTH) {
            contentLength = pair[1];
        }
    })
    return contentLength;
};

const getVisualData = (headers) => {
    let visualData = -1;
    headers.split('\n').forEach((header, _) => {
        const pair = header.split(':');
        if (pair[0] === VISUAL_DATA) {
            visualData = pair[1];
        }
    })
    return visualData;
};

const getVisualDataTs = (headers) => {
    let visualDataTs = -1;
    headers.split('\n').forEach((header, _) => {
        const pair = header.split(':');
        if (pair[0] === VISUAL_DATA_TS) {
            visualDataTs = pair[1];
        }
    })
    return visualDataTs;
};

const getMessageData = (headers) => {
    let messageData = -1;
    headers.split("\n").forEach((header, _) => {
        const pair = header.split(":");
        if (pair[0] === MESSAGE_DATA) {
            messageData = pair[1];
        }
    });
    return messageData;
};

const b64_to_utf8 = str => decodeURIComponent(escape(window.atob(str)));

const getAdditionalData = headers => {
    let additionalData = -1;
    headers.split("\n").forEach((header, _) => {
        const pair = header.split(":");
        if (pair[0] === ADDITIONAL_DATA) {
            additionalData = pair[1];
        }
    });
    if (additionalData !== -1 && /\S/.test(additionalData)) {
        const parsed = b64_to_utf8(additionalData);
        return JSON.parse(parsed);
    }
    return additionalData;
};

# Vanilla-Dashboard

## Features

This is list of all feature that is available on vanilla dashboard. Feature listed below have their own page.

API Documentation:
- Stream/Camera related API : https://docs.nodeflux.io/visionaire-docker-stream/integration-point

- Enrollment API: https://gitlab.com/nodefluxio/vanilla-dashboard/-/blob/master/contract.md

- Event API: https://gitlab.com/nodefluxio/vanilla-dashboard/-/blob/master/event-contract.md

- Websocket : https://docs.nodeflux.io/visionaire-docker-stream/developer-guide

- Global Setting API : https://gitlab.com/nodefluxio/vanilla-dashboard/-/blob/master/global-settings.md

### List All CCTV

This page is used to list all CCTV available. Please check [`camera.js`](https://gitlab.com/nodefluxio/vanilla-dashboard/-/blob/master/client/src/components/pages/camera.js) for the source code

API Used on this page:
- Get list of all streams
- Add Streams
- Delete Streams
- Create Pipeline

currently, when camera is succesfully added, pipeline will automatically be created.

To show visualisation, we use visualstreamer.js class that already provided in [this repo](https://gitlab.com/nodefluxio/vanilla-dashboard/-/blob/master/client/public/js/visualstreamer.js). this class, along with the code can be used anywhere you desire. visualisation class is also used in the [`VisualisationWrapper`](https://gitlab.com/nodefluxio/vanilla-dashboard/-/blob/master/client/src/components/molecules/VisualisationWrapper/index.js) component


### CCTV Details and Events

This page is used to look deeper into one cctv camera. it will show visualisation along with event dumps. Please check [`cameradetail.js`](https://gitlab.com/nodefluxio/vanilla-dashboard/-/blob/master/client/src/components/pages/cameradetail.js) for the source code

API used on this page:
- Get Stream
- Websocket

Like on previous page, we also use [`VisualisationWrapper`](https://gitlab.com/nodefluxio/vanilla-dashboard/-/blob/master/client/src/components/molecules/VisualisationWrapper/index.js) to show the camera visualisation.

For the event dumps, we can call the websocket. Then the data will be displayed using component called [`Event`](https://gitlab.com/nodefluxio/vanilla-dashboard/-/blob/master/client/src/components/molecules/Event/index.js). Event component is just a component used to display picture and text data from event dumps.

### Event History

This page is used to get past event data. We can filter and sort the event data based on flag and or date. Please check [`events.js`]((https://gitlab.com/nodefluxio/vanilla-dashboard/-/blob/master/client/src/components/pages/events.js)) for the source code.

API used on this page:
- Get Event

the Get Event API itself already able to filter and sort based

Component used on this page is mainly `Event` component and Input component (text, select and date selector) for filtering and sorting

### Enrollment

This page use to manage all enrolment CRUD. Please check
[`enrollment.js`]((https://gitlab.com/nodefluxio/vanilla-dashboard/-/blob/master/client/src/components/pages/camera.js)) for the source code.

API used on this page:
- Get List Enrollment (get all enrollment data)
- Delete Enrollment
- Add Enrollment
- Update Enrollment
- Get Enrollment (get specific enrollment data)

We currently use `react-dropzone` for uploading enrollment image, and then call Add Enrollment or Update Enrollment to save the data. 

The component used for displaying enrollment image is `ThumbnailCard`. It's just component to display image with title and actions.


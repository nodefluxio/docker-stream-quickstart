# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Add page for list all available license

### Changed

- Add license information serial number at page camera detail

## [1.1.0] - 2021-07-09

### Added

- Add export event async button
- Conditional rendering for CES status
- Add api for export event history
- Add agent CES status on the footer
- Add warning modal when about to do new enrollment but Coordinator status is disconnected or syncing
- fix some inappropriate proptypes
- fix visualstreamer import

### Changed

- add status error at api ces agent get status

## [1.0.0] - 2021-07 - 07

### Added

- Support plugin system with modular federation (need webpack 5, hence use webpack 5 + babel 7 and forget CRA)
- search dukcapil layout

### Added

- add proccess image converter to JPEG, image remove EXIF, and image compressing at API create and update face enrollment

### Changed

- change max size image enrollment from 800 KB to 25 MB

## [0.16.0] - 2021-06-25

### Added

- add new endpoint at CES agent `/v1/agent/status` for get sync status
- add new service CES agent
- add new service CES coordinator
- add API for handle push new event to coordinator
- add GET API for serve data to agent

### Changed

- implement delete face variation when rollback at enrollment API update
- add rollback scenario at enrollment api POST, UPDATE, & DELETE
- integrate api enrollment create, update & delete at vanilla dashboard with CES agent
- create api get enrollment id by face id
- update vanilla create enrollment now can create custom enrollment if face id supplied

## [0.15.3] - 2021-06-25

### Bugfix

- clean refresh token
- add auto refresh

## [0.15.2] - 2021-06-22

### Bugfix

- if no access token generated from refresh token, logout

## [0.15.1] - 2021-06-22

### Bugfix

- validate user_access or password

## [0.15.0] - 2021-06-16

### Added

- add delete all enrollment feature
- login page and feature
- add ping default handler
- add goroutine websocket read data from client vanilla exclude send back data from read data
- change wording on batch enrollment button based on it's completion status
- add delete all enrollment backend
- add `auth/token` POST method
- add validation for `auth/token` POST
- add bcrypt for encrypted password

## [0.14.0] - 2021-06-03

### Added

- add support for analytic vehicler counting high way at analytic assignment

### Changed

- update handler to read analytic by type at event dumping and event websocket

## [0.13.7] - 2021-06-03

### Added

- add flag to filter recognized event only on camera detail

## [0.13.6] - 2021-05-25

- fix image delay
- fix search by url on camera detail

## [0.13.5] - 2021-05-24

### Changed

- change state loading when enroll success
- publish release tag and latest image registry to dockerhub

### Fxed

- button save loading state not change after add enrollment success
- modal batch enrollment not appear

## [0.13.4] - 2021-05-18

### Added

- minimize reload page on camera list changes
- no reload when change analytic
- auto select first analytic on camera detail

## [0.13.3] - 2021-05-11

### Changed

- method for reconnection websocket
- change max similarity persentage to 99.99% at event dumping and websocket for event FR

## [0.13.2] - 2021-05-06

### Added

- loading state on batch and standard add enrollment
- loading state on event search on event page
- loading state on enrollment page
- loading state on camera list page
- search on enrollment page
- reload cameradetail page every 30 mins

## [0.13.1] - 2021-05-05

### Changed

- add log time elapsed when use face enrollment API at fremisn api
- change api image face enrollment `/api/face/image/:faceImgId` from raw image to thumbnail image

## [0.13.0] - 2021-04-29

### Changed

- add LPR support for analytic assignment, dumping and websocket

## [0.12.0] - 2021-04-16

### Changed

- update all config for run apps
- update similarity check
- update dumping image detection decode from base 64 to blob

## [0.11.1] - 2021-03-31

### Changed

- Add Deleted IS NULL when get detail with faceID

### Deleted

- remove function check if face_ids array data is empty

## [v0.11.0]

### Added

- add JSON Config when analytics assignment

### Changed

- limit live event data up to 50 data, otherwise must filter data to look up past data

## [v0.10.0]

## [0.10.0] - 2021-03-30

### Added

- add global config api for add similarity
- analytic assignment feature draw couting line
- alert sound when recognized face
- add event history to camera detail

### Changed

- add face thumbnail to face enrollment
- update data structure for websocket event

### Fixed

- fix dumping service can't save data from analytic at counting event

## [0.9.0] - 2021-03-19

### Added

- add new endpoint `/api/face/image/:id` for generate image
- feature batch face enrollment

### Changed

- list image at page enrollment now consume source image from API `/api/face/image/:id`
- remove face base64 image form response at API `/api/enrollment`
- fix limit 6 cameras by switching to static image with auto refresh interval

## [0.8.0] - 2021-02-10

### Changed

- add column identity_number and status to form update and add enrollment
- add column identity_number and status to api PUT, POST, GET `api/enrollment`

## [v0.7.1] - 2021-01-21

### Fixed

- fix latitude longitude

## [v0.7.0] - 2021-01-21

### Added

- create pipeline after create stream

## [v0.6.0] - 2021-01-19

### Added

- add api for update enrollment by id `/api/enrollment/:id`
- pagination on enrollment page
- upload multiple images
- add api for get detail enrollment by id `/api/enrollment/:id`

### Changed

- api delete enrollment
- update enrollment page
- change payload api get enrollment `/api/enrollment` and add pagination support
- api create enrollment `/api/enrollment` now support multiple image

## [v0.5.1] - 2021-01-12

### Changed

- dumping service now delete face id in FRemis if unsynchronized happen between enrollment database in vanilla dashboard and FRemis
- websocket service `/api/event_channel` now will change `primary_text` to `UNKNOWN` if face id not found in enrollment database vanilla dashboard

## [v0.4.0] - 2021-01-08

### Added

- api for get event history with pagination `/api/events`
- event dumping services
- auto reconnect websocket in event dumping
- cronjob for handling event history partitioning

## [v0.3.1] - 2021-01-08

### Changed

- fix typos on event history page
- fix query string reading

## [v0.3.0] - 2021-01-07

### Added

- event history page

## [v0.2.17] - 2020-12-16

### Added

- CRUD api face enrollment `/api/enrollment`
- websocket endpoint `/api/event_channel` for middleman communication between the visionaire v4 websocket connection and nodeflux vanilla dashbord

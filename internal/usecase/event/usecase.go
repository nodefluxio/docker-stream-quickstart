package event

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/gorilla/websocket"
	"github.com/robfig/cron"
	logutil "gitlab.com/nodefluxio/goutils/pkg/log"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/entity"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/pkg/imageprocessing"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/pkg/util"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/pkg/zipper"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/presenter"
	"gitlab.com/nodefluxio/vanilla-dashboard/internal/repository"
	"golang.org/x/net/context"
)

// ServiceImpl struct to represent quota transaction service
type ServiceImpl struct {
	EnrolledFaceRepo     repository.EnrolledFace
	VehicleRepo          repository.Vehicle
	FRemisRepo           repository.FRemis
	WSHubRepo            repository.WSHub
	URLGridLiteWS        string
	CronjobPartitionSpec string
	EventRepo            repository.Event
	StreamRepo           repository.Stream
	GlobalSettingRepo    repository.GlobalSetting
	SiteRepo             repository.Site
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 1 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512

	dataLayout      = "2006-01-02"
	dateTimelayout  = "2006-01-02 15:04:05"
	dateTimelayout2 = "2006-01-02_15:04:05"
	exportDir       = "./tmp/events"
)

var exportStatus entity.ExportEventStatus

//  readPump to handle data from client
func (s *ServiceImpl) readPump(client *entity.Client) {
	// this function only for handling ping handler.
	// without this function run in background ping handler will not executed
	for {
		_, _, err := client.Conn.NextReader()
		if err != nil {
			return
		}
	}
}

// InitiateDataStream is a middleman between the websocket connection and the hub.
func (s *ServiceImpl) InitiateDataStream(streamID string, nodeNum int, conn *websocket.Conn) {
	client := &entity.Client{
		Conn:     conn,
		Send:     make(chan []byte, 256),
		StreamID: streamID,
	}
	s.WSHubRepo.RegisterClient(client)
	defer func() {
		s.WSHubRepo.UnregisterClient(client)
		client.Conn.Close()
		logutil.LogObj.SetInfoLog(map[string]interface{}{}, "close conn and unregister client")
	}()
	url := fmt.Sprintf("%s/%d/%s", s.URLGridLiteWS, nodeNum, streamID)
	gridLiteWS, _, err := websocket.DefaultDialer.Dial(url, nil)
	logutil.LogObj.SetInfoLog(map[string]interface{}{"url": url}, "trying to established websocket connection...")
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{"err": err}, "failed to established websocket connection")
		return
	}
	logutil.LogObj.SetInfoLog(map[string]interface{}{"url": url}, "websocket connection established")
	defer gridLiteWS.Close()
	done := make(chan struct{})
	defer close(done)
	// ======= set ping handler =======
	// this function to get ping handler and send back pong message if receive (handle by system not userbase ping)
	// this function can test by websocat use --ping-timeout
	client.Conn.SetPingHandler(func(message string) error {
		client.Conn.WriteControl(websocket.PongMessage, []byte(message), time.Time{})
		return nil
	})
	go s.readPump(client)
	// ======= set ping handler =======
	logutil.LogObj.SetInfoLog(map[string]interface{}{}, "before for loop read message")
	for {

		_, message, err := gridLiteWS.ReadMessage()
		if err != nil {
			logutil.LogObj.SetErrorLog(map[string]interface{}{
				"err": err,
				"url": url,
			},
				"error when read message")
			break
		}
		s.sendMessage(message, client)
	}
	logutil.LogObj.SetInfoLog(map[string]interface{}{}, "after for loop read message")
}

func (s *ServiceImpl) sendMessage(message []byte, c *entity.Client) {
	var parsedMessage entity.Message
	err := json.Unmarshal(message, &parsedMessage)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{"err": err}, "failed to parse message")
	}

	c.Conn.SetWriteDeadline(time.Now().Add(writeWait))

	w, err := c.Conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return
	}

	messageOut, _ := s.prepareDataWS(&parsedMessage)
	analyticCode := strings.Split(parsedMessage.AnalyticID, "-")
	if len(analyticCode) < 2 {
		logutil.LogObj.SetErrorLog(map[string]interface{}{"err": err}, "analytic id not valid")
		return
	}
	switch analyticCode[1] {
	case "FR":
		messageOut, _ = s.prepareDataPipelineFR(&parsedMessage, messageOut)
	case "VC":
		messageOut, _ = s.prepareDataPipelineCounting(&parsedMessage, messageOut)
	case "PC":
		messageOut, _ = s.prepareDataPipelineCounting(&parsedMessage, messageOut)
	case "LPR":
		messageOut, _ = s.prepareDataPipelineLPR(&parsedMessage, messageOut)
	case "CE":
		messageOut, _ = s.prepareDataPipelineCrowdEstimation(&parsedMessage, messageOut)
	}

	nMsg, _ := json.Marshal(messageOut)
	w.Write(nMsg)

	if err := w.Close(); err != nil {
		return
	}
}

func (s *ServiceImpl) prepareDataWS(messageIn *entity.Message) (*entity.EventWebSocket, error) {
	// get location name
	stream, err := s.StreamRepo.GetDetail(context.Background(), messageIn.NodeNum, messageIn.StreamID)
	if err != nil {
		return nil, err
	}

	nsi, err := base64.StdEncoding.DecodeString(messageIn.Image)
	if err != nil {
		return nil, err
	}
	messageOut := entity.EventWebSocket{
		AnalyticID:     messageIn.AnalyticID,
		PrimaryImage:   []byte(""),
		SecondaryImage: nsi,
		StreamID:       messageIn.StreamID,
		Timestamp:      time.Unix(int64(messageIn.Timestamp), 0),
		Location:       stream.StreamName,
		Label:          messageIn.PrimaryText,
		Result:         messageIn.SecondaryText,
	}
	return &messageOut, nil
}

func (s *ServiceImpl) prepareDataPipelineFR(messageIn *entity.Message, messageOut *entity.EventWebSocket) (*entity.EventWebSocket, error) {
	pipelineData := messageIn.PipelineData.(map[string]interface{})

	faceID := pipelineData["face_id"].(string)
	messageOut.Result = faceID

	status := pipelineData["status"].(string)

	switch status {
	case "UNKNOWN":
		messageOut.Label = "unrecognized"
		break

	case "KNOWN":
		faceIDuin64, _ := strconv.ParseUint(faceID, 10, 64)
		dataEnrolledFace, err := s.EnrolledFaceRepo.GetDetailwFaceID(context.Background(), faceIDuin64)
		if err != nil && err.Error() != "record not found" {
			logutil.LogObj.SetErrorLog(map[string]interface{}{
				"err":           err,
				"pipeline_data": pipelineData,
			},
				"error on when get to enrolledFaceRepo")
			return messageOut, err
		}

		if dataEnrolledFace != nil {
			similarity := pipelineData["similarity"].(float64)
			similarity = math.Floor(similarity*100) / 100
			globalConf, err := s.GlobalSettingRepo.GetCurrent(context.Background())
			if err != nil && err.Error() != "record not found" {
				logutil.LogObj.SetErrorLog(map[string]interface{}{
					"err": err,
				},
					"error when get global setting data")
				return nil, err
			}
			if globalConf == nil {
				globalConf = &entity.GlobalSetting{}
			}
			logutil.LogObj.SetInfoLog(map[string]interface{}{"detection_confidence": similarity, "global_config_similarity": globalConf.Similarity}, "checking similarity detection with config")
			if similarity >= globalConf.Similarity {
				nSimilarity := fmt.Sprintf("%.0f", similarity*100)
				if nSimilarity == "100" {
					nSimilarity = "99.99"
				}
				messageOut.Result = nSimilarity + "%" + " - " + dataEnrolledFace.Name
				messageOut.Label = "recognized"
				messageOut.PrimaryImage = dataEnrolledFace.Image
			} else {
				messageOut.Label = "unrecognized"
			}
		} else {
			messageOut.Label = "unrecognized"
		}

		break
	}
	return messageOut, nil
}

func (s *ServiceImpl) prepareDataPipelineCounting(messageIn *entity.Message, messageOut *entity.EventWebSocket) (*entity.EventWebSocket, error) {
	pipelineData := messageIn.PipelineData.(map[string]interface{})
	messageOut.Label = pipelineData["label"].(string)
	messageOut.Result = pipelineData["area_name"].(string)

	return messageOut, nil
}

func (s *ServiceImpl) prepareDataPipelineLPR(messageIn *entity.Message, messageOut *entity.EventWebSocket) (*entity.EventWebSocket, error) {
	pipelineData := messageIn.PipelineData.(map[string]interface{})

	pipelineBoundingBox := pipelineData["bounding_box"].(map[string]interface{})
	boundingBox := imageprocessing.BoundingBox{
		Top:    pipelineBoundingBox["top"].(float64),
		Left:   pipelineBoundingBox["left"].(float64),
		Width:  pipelineBoundingBox["width"].(float64),
		Height: pipelineBoundingBox["height"].(float64),
	}

	plateImg, err := imageprocessing.Base64toCroppedJpg(messageIn.Image, &boundingBox)
	if err != nil {
		return messageOut, nil
	}

	messageOut.PrimaryImage = plateImg.Bytes()
	messageOut.Label = pipelineData["plate_number"].(string)

	vehicle, err := s.VehicleRepo.GetByPlateNumber(context.Background(), pipelineData["plate_number"].(string))
	if err != nil {
		logutil.LogObj.SetDebugLog(map[string]interface{}{
			"err":           err,
			"pipeline_data": pipelineData,
		},
			"error when get data with plate number")
		messageOut.Result = pipelineData["label"].(string)
		return messageOut, nil
	}
	messageOut.Result = pipelineData["label"].(string) + "-" + vehicle.Status + "-" + vehicle.Name

	return messageOut, nil
}

func (s *ServiceImpl) prepareDataPipelineCrowdEstimation(messageIn *entity.Message, messageOut *entity.EventWebSocket) (*entity.EventWebSocket, error) {
	pipelineData := messageIn.PipelineData.(map[string]interface{})
	messageOut.Label = fmt.Sprintf("+/- %s", messageIn.PrimaryText)
	messageOut.Result = pipelineData["area"].(string)

	return messageOut, nil
}

// CronjobPartition is function for scheduling running partition
func (s *ServiceImpl) CronjobPartition(ctx context.Context) error {
	logutil.LogObj.SetInfoLog(map[string]interface{}{}, "setup cronjob partition")
	c := cron.New()
	c.AddFunc(s.CronjobPartitionSpec, func() { s.Partition(ctx) })

	// Start cron with one scheduled job
	logutil.LogObj.SetInfoLog(map[string]interface{}{}, "start cron")
	c.Start()
	return nil
}

// Partition is function for running event partition
func (s *ServiceImpl) Partition(ctx context.Context) error {
	logutil.LogObj.SetInfoLog(map[string]interface{}{}, "start running event partition")
	err := s.EventRepo.Partition(ctx, time.Now())
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{"err": err}, "failed running event partition")
		return err
	}

	err = s.EventRepo.Partition(ctx, time.Now().AddDate(0, 0, 1))
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{"err": err}, "failed running event partition")
		return err
	}

	err = s.EventRepo.Partition(ctx, time.Now().AddDate(0, 0, 2))
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{"err": err}, "failed running event partition")
		return err
	}
	logutil.LogObj.SetInfoLog(map[string]interface{}{}, "finish running event partition")
	return nil
}

// Dumping is a middleman between the websocket connection and the hub.
func (s *ServiceImpl) Dumping(ctx context.Context) {
	for {
		done := make(chan struct{})
		doneR := make(chan struct{})

		url := s.URLGridLiteWS
		logutil.LogObj.SetInfoLog(map[string]interface{}{"url": url}, "trying to established websocket connection...")
		ctxDial, _ := context.WithTimeout(ctx, 10*time.Second)
		gridLiteWS, _, err := websocket.DefaultDialer.DialContext(ctxDial, url, nil)
		if err != nil {
			logutil.LogObj.SetErrorLog(map[string]interface{}{"err": err}, "failed to established websocket connection")
			time.Sleep(1 * time.Second)
			continue
		}
		logutil.LogObj.SetInfoLog(map[string]interface{}{"url": url}, "websocket connection established")
		//read
		func() {
			defer close(done)
			defer gridLiteWS.Close()

			defer logutil.LogObj.SetInfoLog(map[string]interface{}{}, "read closed")
			for {
				select {
				case <-doneR:
					return

				default:
					_, message, err := gridLiteWS.ReadMessage()
					if err != nil {
						logutil.LogObj.SetErrorLog(map[string]interface{}{"err": err}, "failed to read message")
						return
					}
					var data entity.Message
					err = json.Unmarshal(message, &data)
					if err != nil {
						logutil.LogObj.SetErrorLog(map[string]interface{}{"err": err}, "failed to parse message")
						continue
					}

					unbased, err := base64.StdEncoding.DecodeString(data.Image)
					if err != nil {
						logutil.LogObj.SetErrorLog(map[string]interface{}{
							"err": err,
						},
							"error on when unbased base64 string ")
					}

					newImage, err := jpeg.Decode(bytes.NewReader(unbased))
					if err != nil {
						logutil.LogObj.SetErrorLog(map[string]interface{}{
							"err": err,
						},
							"error on when decode image")
					}

					buf := new(bytes.Buffer)
					err = jpeg.Encode(buf, newImage, nil)
					if err != nil {
						logutil.LogObj.SetErrorLog(map[string]interface{}{
							"err": err,
						},
							"error on when encode image")
						return
					}

					analyticCode := strings.Split(data.AnalyticID, "-")
					if len(analyticCode) < 2 {
						logutil.LogObj.SetErrorLog(map[string]interface{}{"err": err}, "analytic id not valid")
						return
					}
					dataEvent := &entity.Event{
						EventType:      fmt.Sprintf("%s-%s", analyticCode[0], analyticCode[1]),
						StreamID:       data.StreamID,
						Detection:      data,
						EventTime:      time.Unix(int64(data.Timestamp), 0),
						SecondaryImage: buf.Bytes(),
					}
					switch analyticCode[1] {
					case "FR":
						dataEvent, err = s.generateEventFR(ctx, dataEvent)
						if err != nil {
							logutil.LogObj.SetErrorLog(map[string]interface{}{"err": err}, "fail generate result FR")
							return
						}

					case "VC":
						dataEvent, err = s.generateEventCounter(ctx, dataEvent)
						if err != nil {
							logutil.LogObj.SetErrorLog(map[string]interface{}{"err": err}, "fail generate result Counter")
							return
						}

					case "PC":
						dataEvent, err = s.generateEventCounter(ctx, dataEvent)
						if err != nil {
							logutil.LogObj.SetErrorLog(map[string]interface{}{"err": err}, "fail generate result Counter")
							return
						}

					case "LPR":
						dataEvent, err = s.generateEventLPR(ctx, dataEvent, &data)
						if err != nil {
							logutil.LogObj.SetErrorLog(map[string]interface{}{"err": err}, "fail generate result LPR")
							return
						}

					case "CE":
						dataEvent, err = s.generateEventCrowdEstimation(ctx, dataEvent, &data)
						if err != nil {
							logutil.LogObj.SetErrorLog(map[string]interface{}{"err": err}, "fail generate result Crowd Estimation")
							return
						}

					default:
						// get location name
						stream, err := s.StreamRepo.GetDetail(ctx, data.NodeNum, data.StreamID)
						if err != nil {
							logutil.LogObj.SetErrorLog(map[string]interface{}{
								"err":         err,
								"node_number": data.NodeNum,
								"stream_id":   data.StreamID,
							}, "fail get detail stream")
							return
						}
						result := entity.EventResult{
							Label:     data.PrimaryText,
							Result:    data.SecondaryText,
							Location:  stream.StreamName,
							Timestamp: dataEvent.EventTime,
						}

						nr, err := json.Marshal(result)
						if err != nil {
							logutil.LogObj.SetErrorLog(map[string]interface{}{"err": err}, "fail marshal data result")
							return
						}
						dataEvent.Result = nr
					}

					// remove image in detection because we already save this image in field secondary_image
					if len(dataEvent.SecondaryImage) >= 0 {
						dataEvent.Detection.Image = ""
					}
					err = s.EventRepo.Create(context.Background(), dataEvent)
					if err != nil {
						logutil.LogObj.SetErrorLog(map[string]interface{}{"err": err}, "fail save event")
						return
					}
					logutil.LogObj.SetInfoLog(map[string]interface{}{
						"event_type":   dataEvent.EventType,
						"stream_id":    dataEvent.StreamID,
						"event_time":   dataEvent.EventTime,
						"primary_text": dataEvent.Detection.PrimaryText,
					},
						"dump event to database success")
				}
			}
		}()
	}
}

func (s *ServiceImpl) generateEventFR(ctx context.Context, data *entity.Event) (*entity.Event, error) {
	var result entity.EventResult
	pipelineData := data.Detection.PipelineData.(map[string]interface{})

	status := pipelineData["status"].(string)
	data.Status = status

	// get location name
	stream, err := s.StreamRepo.GetDetail(ctx, data.Detection.NodeNum, data.Detection.StreamID)
	if err != nil {
		return data, err
	}
	result.Location = stream.StreamName

	faceID := pipelineData["face_id"].(string)
	result.Result = faceID
	result.Timestamp = data.EventTime

	switch status {
	case "UNKNOWN":
		result.Label = "unrecognized"
		break

	case "KNOWN":
		faceIDuin64, _ := strconv.ParseUint(faceID, 10, 64)
		dataEnrolledFace, err := s.EnrolledFaceRepo.GetDetailwFaceID(ctx, faceIDuin64)
		if err != nil && err.Error() != "record not found" {
			logutil.LogObj.SetErrorLog(map[string]interface{}{
				"err":           err,
				"pipeline_data": pipelineData,
			},
				"error on when get to enrolledFaceRepo")
			return data, err
		}

		if dataEnrolledFace != nil {
			similarity := pipelineData["similarity"].(float64)
			similarity = math.Floor(similarity*100) / 100
			globalConf, err := s.GlobalSettingRepo.GetCurrent(ctx)
			if err != nil && err.Error() != "record not found" {
				logutil.LogObj.SetErrorLog(map[string]interface{}{
					"err": err,
				},
					"error when get global setting data")
				return nil, err
			}
			if globalConf == nil {
				globalConf = &entity.GlobalSetting{}
			}
			logutil.LogObj.SetInfoLog(map[string]interface{}{"detection_confidence": similarity, "global_config_similarity": globalConf.Similarity}, "checking similarity detection with config")
			if similarity >= globalConf.Similarity {
				nSimilarity := fmt.Sprintf("%.0f", similarity*100)
				if nSimilarity == "100" {
					nSimilarity = "99.99"
				}
				result.Result = nSimilarity + "%" + " - " + dataEnrolledFace.Name
				result.Label = "recognized"
				data.PrimaryImage = dataEnrolledFace.Image
			} else {
				result.Label = "unrecognized"
				data.Status = "UNKNOWN"
			}

		} else {
			result.Label = "unrecognized"
			data.Status = "UNKNOWN"

			// Delete face id in fremis
			duration := 60 // in second
			ctx, cancel := context.WithTimeout(ctx, time.Duration(duration)*time.Second)
			defer cancel()
			err = s.FRemisRepo.FaceDeleteEnrollment(ctx, []string{faceID})
			if err != nil {
				logutil.LogObj.SetErrorLog(map[string]interface{}{
					"err":     err,
					"face_id": faceID,
				},
					"error on when delete face id in FRemis")
			} else {
				logutil.LogObj.SetInfoLog(map[string]interface{}{"face_id": faceID}, "success delete face id in FRemis")
			}
		}

		break
	}
	nr, err := json.Marshal(result)
	if err != nil {
		return data, err
	}
	data.Result = nr
	return data, nil
}

func (s *ServiceImpl) generateEventCounter(ctx context.Context, data *entity.Event) (*entity.Event, error) {
	var result entity.EventResult
	pipelineData := data.Detection.PipelineData.(map[string]interface{})

	// get location name
	stream, err := s.StreamRepo.GetDetail(ctx, data.Detection.NodeNum, data.Detection.StreamID)
	if err != nil {
		return data, err
	}
	result.Location = stream.StreamName
	result.Label = pipelineData["label"].(string)
	result.Result = pipelineData["area_name"].(string)
	result.Timestamp = data.EventTime

	nr, err := json.Marshal(result)
	if err != nil {
		return data, err
	}
	data.Result = nr
	data.Status = pipelineData["label"].(string)
	return data, nil
}

func (s *ServiceImpl) generateEventLPR(ctx context.Context, data *entity.Event, message *entity.Message) (*entity.Event, error) {
	var result entity.EventResult
	pipelineData := data.Detection.PipelineData.(map[string]interface{})

	// get location name
	stream, err := s.StreamRepo.GetDetail(ctx, data.Detection.NodeNum, data.Detection.StreamID)
	if err != nil {
		return data, err
	}
	result.Location = stream.StreamName
	result.Timestamp = data.EventTime
	result.Label = pipelineData["plate_number"].(string)
	result.Result = pipelineData["label"].(string)
	pipelineBoundingBox := pipelineData["bounding_box"].(map[string]interface{})
	boundingBox := imageprocessing.BoundingBox{
		Top:    pipelineBoundingBox["top"].(float64),
		Left:   pipelineBoundingBox["left"].(float64),
		Width:  pipelineBoundingBox["width"].(float64),
		Height: pipelineBoundingBox["height"].(float64),
	}
	plateImg, err := imageprocessing.Base64toCroppedJpg(message.Image, &boundingBox)
	if err != nil {
		return data, nil
	}

	data.PrimaryImage = plateImg.Bytes()

	vehicle, err := s.VehicleRepo.GetByPlateNumber(ctx, pipelineData["plate_number"].(string))
	if err != nil {
		logutil.LogObj.SetDebugLog(map[string]interface{}{
			"err":           err,
			"pipeline_data": pipelineData,
		},
			"error when get data with plate number")
		nr, err := json.Marshal(result)
		if err != nil {
			return data, err
		}
		data.Result = nr
		return data, nil
	}
	result.Result = result.Result + "-" + vehicle.Status + "-" + vehicle.Name
	nr, err := json.Marshal(result)
	if err != nil {
		return data, err
	}
	data.Result = nr
	data.Status = pipelineData["label"].(string)
	return data, nil
}

func (s *ServiceImpl) generateEventCrowdEstimation(ctx context.Context, data *entity.Event, message *entity.Message) (*entity.Event, error) {
	var result entity.EventResult
	pipelineData := data.Detection.PipelineData.(map[string]interface{})

	// get location name
	stream, err := s.StreamRepo.GetDetail(ctx, data.Detection.NodeNum, data.Detection.StreamID)
	if err != nil {
		return data, err
	}

	result.Location = stream.StreamName
	result.Label = fmt.Sprintf("+/- %s", data.Detection.PrimaryText)
	result.Result = pipelineData["area"].(string)
	result.Timestamp = data.EventTime

	nr, err := json.Marshal(result)
	if err != nil {
		return data, err
	}
	data.Result = nr
	return data, nil
}

// GetHistory is function for get event history with pagination, filtering and searching
func (s *ServiceImpl) GetHistory(ctx context.Context, lastID uint64, timezone string, paging *util.Pagination, userInfo *presenter.AuthInfoResponse) (*presenter.EventHistoryPaging, error) {
	var result presenter.EventHistoryPaging
	eventGroup := make([]*presenter.EventGroup, 0)

	tzLocation, err := time.LoadLocation(timezone)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"err": err,
		},
			"error when load time location")
		return nil, err
	}

	// check filter stream_id
	isStreamIDFiltered := false
	if len(paging.Filter["stream_id"]) > 0 {
		isStreamIDFiltered = true
	}

	// apply role base check
	var newStreamID string
	switch userInfo.Role {
	case string(entity.UserRoleOperator):
		if len(userInfo.SiteID) == 0 {
			logutil.LogObj.SetInfoLog(map[string]interface{}{
				"site_id": userInfo.SiteID,
			},
				"this user not assigned to any site, returning empty event history")
			var result presenter.EventHistoryPaging
			result.Limit = paging.Limit
			result.Events = eventGroup
			return &result, nil
		}
		listStream, err := s.SiteRepo.GetSiteWithStream(ctx, userInfo.SiteID)
		if err != nil {
			logutil.LogObj.SetErrorLog(map[string]interface{}{
				"site_id": userInfo.SiteID,
				"err":     err,
			},
				"failed get list available stream id for this site")
			return nil, err
		}
		var streamIDs []string
		for _, allowedStream := range listStream {
			streamIDs = append(streamIDs, allowedStream.StreamID)
		}

		// check if fillter applied
		if isStreamIDFiltered {
			filterStreamID := strings.Split(paging.Filter["stream_id"], ",")
			var newFilterStreamID []string
			for _, streamID := range filterStreamID {
				checkAvailStream := util.ArrayStringAvailability(streamIDs, streamID)
				if checkAvailStream {
					newFilterStreamID = append(newFilterStreamID, streamID)
				}
			}
			streamIDs = newFilterStreamID
		}
		if len(streamIDs) > 0 {
			newStreamID = strings.Join(streamIDs, ",")
		}
	case string(entity.UserRoleSuperAdmin):
		if isStreamIDFiltered {
			newStreamID = paging.Filter["stream_id"]
		}
	}

	// cast new data filter stream id
	if newStreamID != "" {
		paging.Filter["stream_id"] = newStreamID
	}

	event, err := s.EventRepo.GetWithLastID(ctx, lastID, paging)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"err": err,
		},
			"error on get event history from repository")
		return nil, err
	}
	output := make(map[string][]*presenter.EventData)
	for _, v := range event {
		TzEventTime := v.EventTime.In(tzLocation)
		date := TzEventTime.Format(dataLayout)
		var res entity.EventResult
		err := json.Unmarshal(v.Result, &res)
		if err != nil {
			return nil, err
		}
		output[date] = append(output[date], &presenter.EventData{
			ID:             v.ID,
			AnalyticID:     v.EventType,
			PrimaryImage:   v.PrimaryImage,
			SecondaryImage: v.SecondaryImage,
			Label:          res.Label,
			Result:         res.Result,
			Location:       res.Location,
			Timestamp:      TzEventTime,
		})
	}

	for k, v := range output {
		neg := &presenter.EventGroup{
			Timestamp: k,
			Data:      v,
		}
		eventGroup = append(eventGroup, neg)
	}

	result.Limit = paging.Limit
	result.Events = eventGroup

	return &result, nil
}

// ExportEvent is function for export data event history
func (s *ServiceImpl) ExportEvent(ctx context.Context, paging *util.Pagination, timezone string) error {
	var totalProcess int = 10
	var totalQueryPerProccess int = 1000
	tzLocation, err := time.LoadLocation(timezone)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"err": err,
		},
			"error when load time location")
		return err
	}
	exportStatus = entity.ExportEventStatusRunning

	// empty dir
	dir := "./tmp"
	os.RemoveAll(dir)

	totalData, err := s.EventRepo.Count(ctx, paging)
	if err != nil {
		logutil.LogObj.SetErrorLog(map[string]interface{}{
			"err": err,
		},
			"error on count event history from repository")
		return err
	}
	paging.Limit = totalQueryPerProccess
	pgDetail := paging.CreateProperties(totalData)
	totalPage := pgDetail.TotalPage
	chunkSize := (int(totalPage) + totalProcess - 1) / totalProcess

	pageSegment := [][]int{}
	tmpLastPage := 1
	if totalProcess >= int(totalPage) {
		totalProcess = int(totalPage)
	}
	for i := 0; i < totalProcess; i++ {
		pageList := []int{}
		for j := 0; j < chunkSize; j++ {
			pageList = append(pageList, tmpLastPage)
			tmpLastPage++
		}
		pageSegment = append(pageSegment, pageList)
	}

	logutil.LogObj.SetInfoLog(map[string]interface{}{
		"timezone":          timezone,
		"total_data":        totalData,
		"total_page":        totalPage,
		"chunk_size":        chunkSize,
		"page_segment_size": len(pageSegment),
	},
		"information data event to export")

	if totalData == 0 {
		exportStatus = entity.ExportEventStatusDone
		return errors.New("failed export event data, because no data to proccess")
	}

	doneChannel := make(chan bool)
	errorChannel := make(chan bool)
	for i := 0; i < len(pageSegment); i++ {
		segement := pageSegment[i]
		go func(i int) {
			for j := 0; j < chunkSize; j++ {
				logutil.LogObj.SetInfoLog(map[string]interface{}{
					"topic":   "export-event-history",
					"segment": segement[j],
					"section": j,
				}, "starting proccess save image to file")
				offset := ((segement[j] * totalQueryPerProccess) - totalQueryPerProccess)
				newPaging := util.Pagination{
					Limit:  totalQueryPerProccess,
					Sort:   paging.Sort,
					Page:   segement[j],
					Filter: paging.Filter,
					Search: paging.Search,
					Offset: offset,
				}
				event, err := s.EventRepo.Get(ctx, &newPaging)
				if err != nil {
					errorChannel <- true
				}
				for _, v := range event {
					eventTImeWTz := v.EventTime.In(tzLocation)
					date := eventTImeWTz.Format(dataLayout)
					groupDir := fmt.Sprintf("%s/images/%s", exportDir, date)
					if _, err := os.Stat(groupDir); os.IsNotExist(err) {
						os.MkdirAll(groupDir, os.ModeDir|0755)
					}
					if len(v.PrimaryImage) > 0 {
						out, err := os.Create(fmt.Sprintf("%s/image_primary_%s_%d.jpg", groupDir, date, v.ID))
						if err != nil {
							errorChannel <- true
						}
						image, _, err := image.Decode(bytes.NewReader(v.PrimaryImage))
						if err != nil {
							errorChannel <- true
						}
						// write new image to file
						jpeg.Encode(out, image, nil)
						out.Close()
					}
					if len(v.SecondaryImage) > 0 {
						out, err := os.Create(fmt.Sprintf("%s/image_secondary_%s_%d.jpg", groupDir, date, v.ID))
						if err != nil {
							errorChannel <- true
						}
						image, _, err := image.Decode(bytes.NewReader(v.SecondaryImage))
						if err != nil {
							errorChannel <- true
						}
						// write new image to file
						jpeg.Encode(out, image, nil)
						out.Close()
					}
				}
				logutil.LogObj.SetInfoLog(map[string]interface{}{
					"topic":   "export-event-history",
					"segment": segement[j],
					"section": j,
				}, "finished procces save image to file")
			}
			doneChannel <- true
		}(i)
	}

	go func() {
		count := 0
		errCount := 0
		for {
			select {
			case <-doneChannel:
				count++
			case <-errorChannel:
				count++
				errCount++
			}

			if count == len(pageSegment) && errCount != 0 {
				// empty dir
				dir := "./tmp"
				os.RemoveAll(dir)
				exportStatus = entity.ExportEventStatusError
				logutil.LogObj.SetErrorLog(map[string]interface{}{
					"topic":       "export-event-history",
					"error_count": errCount,
				},
					"error happen when export file")
				return
			}

			if count == len(pageSegment) {
				logutil.LogObj.SetInfoLog(map[string]interface{}{"topic": "export-event-history"}, "finished save image to file")
				logutil.LogObj.SetInfoLog(map[string]interface{}{"topic": "export-event-history"}, "starting generate html file")
				newPaging := util.Pagination{
					Sort:   paging.Sort,
					Filter: paging.Filter,
					Search: paging.Search,
				}
				event, _ := s.EventRepo.GetWithoutImage(ctx, &newPaging)
				eventExportItems := make([]*entity.ExportEventItem, 0)
				for _, v := range event {
					eventTImeWTz := v.EventTime.In(tzLocation)
					date := eventTImeWTz.Format(dataLayout)
					var resultEvent entity.EventResult
					err = json.Unmarshal(v.Result, &resultEvent)
					if err != nil {
						exportStatus = entity.ExportEventStatusError
						logutil.LogObj.SetErrorLog(map[string]interface{}{
							"err": err,
						},
							"failed unmarshal result event")
						return
					}
					resultEvent.Timestamp = resultEvent.Timestamp.In(tzLocation)

					newResult, err := json.Marshal(&resultEvent)
					if err != nil {
						exportStatus = entity.ExportEventStatusError
						logutil.LogObj.SetErrorLog(map[string]interface{}{
							"err": err,
						},
							"failed marshal result event")
						return
					}
					eventExportItems = append(eventExportItems, &entity.ExportEventItem{
						ID:             v.ID,
						EventType:      v.EventType,
						StreamID:       v.StreamID,
						PrimaryImage:   fmt.Sprintf("images/%s/image_primary_%s_%d.jpg", date, date, v.ID),
						SecondaryImage: fmt.Sprintf("images/%s/image_secondary_%s_%d.jpg", date, date, v.ID),
						Detection:      string(v.Detection),
						Result:         string(newResult),
						Status:         v.Status,
						EventTime:      eventTImeWTz.Format(dateTimelayout),
						CreatedAt:      v.CreatedAt.In(tzLocation).Format(dateTimelayout),
					})
				}
				templateData := entity.ExportEventTemplate{
					Items: eventExportItems,
				}
				templatePath := "./template/exportEventTemplate.html"
				generatedFile, err := ParseTemplate(templatePath, templateData)
				if err != nil {
					exportStatus = entity.ExportEventStatusError
					logutil.LogObj.SetErrorLog(map[string]interface{}{
						"template_path": templatePath,
					},
						"failed parse data to template")
					return
				}
				fileName := fmt.Sprintf("%s/index_%s.html", exportDir, time.Now().Format(dateTimelayout2))
				err = ioutil.WriteFile(fileName, generatedFile.Bytes(), 0755)
				if err != nil {
					exportStatus = entity.ExportEventStatusError
					logutil.LogObj.SetErrorLog(map[string]interface{}{
						"file_name": fileName,
					},
						"failed write file")
					return
				}
				logutil.LogObj.SetInfoLog(map[string]interface{}{"topic": "export-event-history"}, "finished parse data export event histroy to html file")
				appZipper := zipper.NewZipper()
				_, err = appZipper.Create(ctx, exportDir, "./tmp/exported_event.zip")
				if err != nil {
					exportStatus = entity.ExportEventStatusError
					logutil.LogObj.SetErrorLog(map[string]interface{}{"topic": "export-event-history"}, "failed compress file")
					return
				}
				logutil.LogObj.SetInfoLog(map[string]interface{}{"topic": "export-event-history"}, "finished compress file export event histroy")
				logutil.LogObj.SetInfoLog(map[string]interface{}{"topic": "export-event-history"}, "done export event histroy")
				exportStatus = entity.ExportEventStatusDone
				return
			}
		}
	}()

	return nil
}

func ParseTemplate(templatePath string, data interface{}) (*bytes.Buffer, error) {
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return nil, err
	}
	return buf, nil
}

func (s *ServiceImpl) CheckExportedEvent(ctx context.Context) (string, error) {
	status := "ready"
	dir := "./tmp/exported_event.zip"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		status = string(exportStatus)
		if status == "" {
			status = "not-started"
		}
	}
	if exportStatus == entity.ExportEventStatusDownloaded {
		status = string(exportStatus)
	}
	return status, nil
}

func (s *ServiceImpl) UpdateStatusExportDownload(ctx context.Context) error {
	exportStatus = entity.ExportEventStatusDownloaded
	return nil
}

func (s *ServiceImpl) GetEventInsight(ctx context.Context, data *entity.EventInsight) (*presenter.EventInsightData, error) {

	data.TimeDeffinition = "today"
	dataToday, err := s.EventRepo.GetInsight(ctx, data)
	if err != nil {
		return nil, err
	}
	data.TimeDeffinition = "yesterday"
	dataYesterday, err := s.EventRepo.GetInsight(ctx, data)
	if err != nil {
		return nil, err
	}
	data.TimeDeffinition = "week"
	dataWeek, err := s.EventRepo.GetInsight(ctx, data)
	if err != nil {
		return nil, err
	}
	data.TimeDeffinition = "month"
	dataMonth, err := s.EventRepo.GetInsight(ctx, data)
	if err != nil {
		return nil, err
	}

	return &presenter.EventInsightData{
		TotalToday:     dataToday.Total,
		TotalYesterday: dataYesterday.Total,
		TotalWeek:      dataWeek.Total,
		TotalMonth:     dataMonth.Total,
	}, err

}

import { COUNTER_LINE, FULL_FRAME, COUNTER_LINE_RULES } from "./roiType";
import {
  FACE_RECOGNITION,
  VEHICLE_COUNTING,
  VEHICLE_COUNTING_HIGHWAY,
  PEOPLE_COUNTING,
  LICENSE_PLATE_RECOGNITION
} from "./analyticName";

const AnalyticMetadata = [];
AnalyticMetadata["NFV4-FR"] = {
  analytic_name: FACE_RECOGNITION,
  roi_type: FULL_FRAME,
  roi_rule: [],
  roi_title: "",
  roi: false
};
AnalyticMetadata["NFV4-VC"] = {
  analytic_name: VEHICLE_COUNTING,
  roi_type: COUNTER_LINE,
  roi_rule: COUNTER_LINE_RULES,
  roi_title: "DRAWING COUNTING LINE",
  roi: true
};
AnalyticMetadata["NFV4-VC-HW"] = {
  analytic_name: VEHICLE_COUNTING_HIGHWAY,
  roi_type: COUNTER_LINE,
  roi_rule: COUNTER_LINE_RULES,
  roi_title: "DRAWING COUNTING LINE",
  roi: true
};
AnalyticMetadata["NFV4-PC"] = {
  analytic_name: PEOPLE_COUNTING,
  roi_type: COUNTER_LINE,
  roi_rule: COUNTER_LINE_RULES,
  roi_title: "DRAWING COUNTING LINE",
  roi: true
};
AnalyticMetadata["NFV4-LPR"] = {
  analytic_name: LICENSE_PLATE_RECOGNITION,
  roi_type: COUNTER_LINE,
  roi_rule: COUNTER_LINE_RULES,
  roi_title: "DRAWING DETECTION LINE",
  roi: true
};
export const ANALYTIC_METADATA = AnalyticMetadata;

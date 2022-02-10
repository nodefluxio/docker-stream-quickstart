import {
  Camera,
  Enrollment,
  CameraDetail,
  Event,
  Assignment,
  Vehicle,
  Plugin,
  License,
  Account
} from "components/pages";

const routes = [
  {
    path: "/",
    exact: true,
    fullscreen: false,
    component: Camera
  },
  {
    path: "/camera/:node/:id",
    exact: true,
    fullscreen: true,
    component: CameraDetail
  },
  {
    path: "/camera",
    exact: true,
    fullscreen: false,
    component: Camera
  },
  {
    path: "/enrollment",
    exact: false,
    fullscreen: false,
    component: Enrollment
  },
  {
    path: "/vehicle",
    exact: false,
    fullscreen: false,
    component: Vehicle
  },
  {
    path: "/event-history",
    exact: false,
    fullscreen: false,
    component: Event
  },
  {
    path: "/assignment",
    exact: false,
    fullscreen: true,
    component: Assignment
  },
  {
    path: "/plugin/",
    exact: false,
    fullscreen: false,
    component: Plugin
  },
  {
    path: "/license",
    exact: false,
    fullscreen: false,
    component: License
  },
  {
    path: "/account",
    exact: false,
    fullscreen: false,
    component: Account
  }
];

export default routes;

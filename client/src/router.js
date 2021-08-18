import {
  Camera,
  Enrollment,
  CameraDetail,
  Event,
  Assignment,
  Plugin,
  License
} from "components/pages";

const routes = [
  {
    path: "/",
    exact: true,
    fullscreen: false,
    component: Camera
  },
  {
    path: "/camera/:id",
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
    path: "/event-history",
    exact: false,
    fullscreen: false,
    component: Event
  },
  {
    path: "/assignment",
    exact: false,
    fullscreen: false,
    component: Assignment
  },
  {
    path: "/plugin/*",
    exact: false,
    fullscreen: false,
    component: Plugin
  },
  {
    path: "/license",
    exact: false,
    fullscreen: false,
    component: License
  }
];

export default routes;

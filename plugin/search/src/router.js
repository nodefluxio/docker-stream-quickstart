import { SearchFace, SearchPlate } from "components/pages";

const routes = [
  {
    path: "/",
    exact: true,
    fullscreen: false,
    component: SearchFace
  },
  {
    path: "/vehicle",
    exact: true,
    fullscreen: false,
    component: SearchPlate
  }
];

export default routes;

import { createRouter, createWebHistory } from "vue-router";

import AppHeader from "@/views/AppHeader.vue";
import ArbitraryPayload from "@/views/ArbitraryPayload.vue";
import CoopStatus from "@/views/CoopStatus.vue";
import FirstContact from "@/views/FirstContact.vue";
import GetConfig from "@/views/GetConfig.vue";
import GetPeriodicals from "@/views/GetPeriodicals.vue";

function routeNameProp(route) {
  return {
    route: route.name,
  };
}

const routes = [
  {
    name: "arbitrary_payload",
    path: "/",
    components: {
      default: ArbitraryPayload,
      header: AppHeader,
    },
    props: {
      default: routeNameProp,
      header: routeNameProp,
    },
  },
  {
    name: "first_contact",
    path: "/first_contact/",
    components: {
      default: FirstContact,
      header: AppHeader,
    },
    props: {
      default: routeNameProp,
      header: routeNameProp,
    },
  },
  {
    name: "get_periodicals",
    path: "/get_periodicals/",
    components: {
      default: GetPeriodicals,
      header: AppHeader,
    },
    props: {
      default: routeNameProp,
      header: routeNameProp,
    },
  },
  {
    name: "get_config",
    path: "/get_config/",
    components: {
      default: GetConfig,
      header: AppHeader,
    },
    props: {
      default: routeNameProp,
      header: routeNameProp,
    },
  },
  {
    name: "coop_status",
    path: "/coop_status/",
    components: {
      default: CoopStatus,
      header: AppHeader,
    },
    props: {
      default: routeNameProp,
      header: routeNameProp,
    },
  },
  {
    name: "doc",
    path: "/doc/",
  },
  {
    path: "/:catchAll(.*)",
    redirect: "/",
  },
];

const router = createRouter({
  routes,
  history: createWebHistory("/proto-explorer/"),
});

export default router;

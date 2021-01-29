import { createRouter, createWebHashHistory } from "vue-router";

import Main from "@/views/Main.vue";
import Mission from "@/views/Mission.vue";
import Artifact from "@/views/Artifact.vue";

const routes = [
  {
    name: "home",
    path: "/",
    component: Main,
    props: true,
    children: [
      {
        name: "mission",
        path: "mission/:missionId/",
        components: {
          mission: Mission,
        },
        props: true,
      },
      {
        name: "artifact",
        path: "artifact/:artifactId/",
        components: {
          artifact: Artifact,
        },
        props: true,
      },
    ],
  },
];

const router = createRouter({
  routes,
  history: createWebHashHistory(),
  scrollBehavior(to, from, savedPosition) {
    // always scroll to top
    return { top: 0 };
  },
});

export default router;

<template>
  <div class="max-w-3xl w-full mx-auto px-4 mb-4 text-sm leading-tight space-y-2">
    <div v-if="activeNews.length > 0" class="text-green-500">
      <p class="uppercase">What's New</p>
      <ul v-for="(news, index) in activeNews" :key="index">
        <li>
          {{ news.datetime.toISOString().substring(0, 10) }} &mdash;
          <span v-html="news.content"></span>
        </li>
      </ul>
    </div>
  </div>
  <router-view />
  <font-switcher />
</template>

<script>
import FontSwitcher from "@/components/FontSwitcher.vue";

export default {
  components: {
    FontSwitcher,
  },

  data() {
    return {
      news: [
        {
          datetime: new Date(1614525360000),
          expiry: new Date(1614784563000),
          content: `Boosts support: New <span class="uppercase">active boost effects</span>
          configuration section, new stats “SE gain w/ empty habs start” and “Boost duration”,
          as well as a guide on optimizing SE gain from prestiges.`,
        },
      ],
    };
  },

  computed: {
    activeNews() {
      const now = new Date();
      return this.news.filter(entry => now < entry.expiry);
    },
  },
};
</script>

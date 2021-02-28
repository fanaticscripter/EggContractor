<template>
  <div class="max-w-3xl w-full mx-auto px-4 mb-4 text-sm leading-tight space-y-2">
    <div class="text-red-500">
      You are using an alpha build of this tool. There is no compatibility guarantee for shared
      links, and many possible improvements have not yet been implemented. Future development hinges
      on community response, in terms of traffic and perceived mind share. As such, basic
      information about your visit (e.g. your browser brand and screen size) is collected through
      Google Analytics in order to guide development. Feedback is welcome on
      <a href="https://discord.gg/egginc" target="_blank" class="underline hover:text-dark-70"
        >Discord</a
      >
      and
      <a
        href="https://github.com/fanaticscripter/EggContractor/discussions"
        target="_blank"
        class="underline hover:text-dark-70"
        >GitHub</a
      >.
    </div>
    <div class="text-green-500">
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
</template>

<script>
export default {
  data() {
    return {
      news: [
        {
          datetime: new Date(1614525360000),
          expiry: new Date(1614784563000),
          content: `Boosts support: New <span class="uppercase">active boost effects</span>
          configuration section, new stats “SE gain w/ empty habs start” and “Boost duration”,
          as well as a guide on optimizing SE gain from presitges.`,
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

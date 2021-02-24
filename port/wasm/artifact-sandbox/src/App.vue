<template>
  <div class="max-w-3xl w-full mx-auto px-4 mb-4 text-red-500 text-sm leading-tight">
    You are using an alpha build of this tool. There is no compatibility guarantee for shared links,
    and many possible improvements have not yet been implemented. Future development hinges on
    community response, in terms of traffic and perceived mind share. As such, basic information
    about your visit (e.g. your browser brand and screen size) is collected through Google Analytics
    in order to guide development. Feedback is welcome on
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
  <router-view />
</template>

<script>
export default {};
</script>

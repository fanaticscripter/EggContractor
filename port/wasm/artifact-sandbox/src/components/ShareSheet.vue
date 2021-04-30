<template>
  <div v-if="show" class="fixed z-50 inset-0 overflow-y-auto">
    <div
      class="flex items-end justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0"
    >
      <!-- Background overlay -->
      <div
        class="fixed inset-0 transition-opacity"
        aria-hidden="true"
        @click="$emit('update:show', false)"
      >
        <div class="absolute inset-0 bg-gray-500 opacity-75"></div>
      </div>

      <!-- This element is to trick the browser into centering the modal contents. -->
      <span class="hidden sm:inline-block sm:align-middle sm:h-screen" aria-hidden="true"
        >&#8203;</span
      >

      <!-- Modal panel -->
      <div
        class="inline-block align-bottom bg-dark-25 rounded-lg px-4 pt-5 pb-4 text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg w-full sm:p-6"
        role="dialog"
        aria-modal="true"
        aria-labelledby="modal-headline"
      >
        <div class="absolute top-0 right-0 pt-4 pr-4">
          <button
            type="button"
            class="rounded-md text-dark-70 hover:text-dark-60 focus:outline-none"
            @click="$emit('update:show', false)"
          >
            <span class="sr-only">Close</span>
            <svg
              class="h-5 w-5"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              aria-hidden="true"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M6 18L18 6M6 6l12 12"
              />
            </svg>
          </button>
        </div>

        <div>
          <p class="text-base">Copy the shareable link</p>
          <div class="mt-1 flex rounded-md shadow-sm">
            <div class="relative flex items-stretch flex-grow focus-within:z-10">
              <input
                type="text"
                class="bg-dark-20 focus:ring-blue-500 focus:border-blue-500 block w-full rounded-none rounded-l-md sm:text-sm border truncate"
                :value="link"
                readonly
              />
            </div>
            <button
              class="-ml-px relative inline-flex items-center space-x-2 px-2 py-2 border border-gray-500 text-sm font-medium rounded-r-md bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-1 focus:ring-blue-500 focus:border-blue-500"
              @click="copyLink"
            >
              <svg
                class="h-5 w-5 transition duration-300 ease-in-out"
                :class="linkCopied ? 'text-blue-300' : 'text-gray-50'"
                viewBox="0 0 20 20"
                fill="currentColor"
              >
                <path d="M8 2a1 1 0 000 2h2a1 1 0 100-2H8z" />
                <path
                  d="M3 5a2 2 0 012-2 3 3 0 003 3h2a3 3 0 003-3 2 2 0 012 2v6h-4.586l1.293-1.293a1 1 0 00-1.414-1.414l-3 3a1 1 0 000 1.414l3 3a1 1 0 001.414-1.414L10.414 13H15v3a2 2 0 01-2 2H5a2 2 0 01-2-2V5zM15 11h2a1 1 0 110 2h-2v-2z"
                />
              </svg>
            </button>
          </div>
        </div>

        <div class="relative mt-2 mb-1">
          <div class="absolute inset-0 flex items-center">
            <div class="w-full border-t border-dark-30"></div>
          </div>
          <div class="relative flex justify-center">
            <span class="px-2 bg-dark-25 text-dark-60 text-sm uppercase">Or</span>
          </div>
        </div>

        <div>
          <p class="text-base">Share an image</p>
          <p class="flex items-center">
            <input
              id="show_footnotes"
              name="show_footnotes"
              type="checkbox"
              class="h-4 w-4 bg-dark-20 text-blue-600 focus:ring-blue-500 focus:ring-offset-2 focus:ring-offset-dark-30 rounded"
              v-model="showFootnotes"
              @change="$emit('update:showFootnotes', $event.target.checked)"
            />
            <label for="show_footnotes" class="ml-2 block text-sm">Show effect footnotes</label>
          </p>
          <p class="text-xs text-dark-60">Right click / long press to copy or share</p>
          <div
            class="mt-2 relative w-full border border-dark-30 rounded-lg shadow-lg"
            :class="imageURL === '' ? 'bg-dark-20' : null"
            :style="{ paddingBottom: `${placeholderAspectRatio * 100}%` }"
          >
            <img class="absolute top-0 left-0" v-if="imageURL !== ''" :src="imageURL" />
            <div
              v-else
              class="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 text-lg uppercase"
            >
              <template v-if="generateImageErrored">
                <div class="text-center text-red-500">Error generating image</div>
                <div class="max-w-xs text-xs text-dark-60 normal-case break-all text-justify">
                  {{ encodedError }}
                </div>
                <div class="my-2 text-center">
                  <button
                    type="button"
                    class="inline-flex items-center px-3 py-2 border border-dark-30 shadow-sm text-sm leading-4 font-medium rounded-md bg-red-600 hover:bg-red-700 focus:outline-none"
                    @click="copyErrorCode()"
                  >
                    Copy error code
                  </button>
                </div>
              </template>
              <template v-else>Generating...</template>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import copyTextToClipboard from "copy-text-to-clipboard";
import html2canvas from "html2canvas";
import { Builds } from "@/lib/models";

let runId;

export default {
  props: {
    show: Boolean,
    showFootnotes: Boolean,
    builds: {
      type: Builds,
      required: true,
    },
  },

  data() {
    return {
      linkCopied: false,
      imageURL: "",
      placeholderAspectRatio: 1,
      generateImageErrored: false,
      generateImageError: "",
    };
  },

  emits: ["update:show", "update:showFootnotes"],

  computed: {
    link() {
      return (
        window.location.origin +
        window.location.pathname +
        this.$router.resolve({
          name: "builds",
          params: { serializedBuilds: this.builds.serialize() },
        }).href
      );
    },

    encodedError() {
      return btoa(`${this.generateImageError}` || "");
    },
  },

  watch: {
    async show() {
      if (this.show) {
        this.linkCopied = false;
        await this.generateImage();
      }
    },

    async showFootnotes() {
      // For some reason, await nextTick() doesn't wait until the DOM changes
      // are applied, but a macrotask setTimeout (even at 1ms) fixes it.
      await new Promise(resolve => setTimeout(resolve, 50));
      await this.generateImage();
    },
  },

  methods: {
    copyLink() {
      copyTextToClipboard(this.link);
      this.linkCopied = true;
    },

    async generateImage() {
      // Generate a unique runId for the invocation, so that we don't actually
      // render the generated image if there's a later invocation that overrides
      // this one.
      runId = Math.floor(Math.random() * 65536);
      const currentRunId = runId;

      this.imageURL = "";
      this.generateImageErrored = false;
      this.generateImageError = "";
      const target = document.getElementById("builds");
      if (target.offsetWidth !== 0) {
        this.placeholderAspectRatio = target.offsetHeight / target.offsetWidth;
      }
      try {
        const { x: offsetX, y: offsetY, width, height } = target.getBoundingClientRect();
        const canvas = await html2canvas(target, {
          allowTaint: false,
          useCORS: true,
          backgroundColor: null,
          // Fix possible scrollbar-induced shift.
          scrollX: -window.scrollX,
          scrollY: 0,
          windowWidth: document.documentElement.offsetWidth,
          windowHeight: document.documentElement.offsetHeight,
          onclone: doc => {
            // Fix weird wrong scroll position issue (manifesting as large blank
            // area at the top of the genreated canvas).
            //
            // Kudos:
            // https://github.com/niklasvh/html2canvas/issues/1878#issuecomment-756504779
            //
            // Using html2canvas has always been shots in the dark like this.
            const clone = doc.getElementById("builds");
            clone.style.position = "fixed";
            clone.style.left = `${offsetX}px`;
            clone.style.top = `${offsetY}px`;
            clone.style.width = `${width}px`;
            clone.style.height = `${height}px`;
          },
        });

        if (currentRunId !== runId) {
          return;
        }

        canvas.toBlob(blob => {
          if (currentRunId !== runId) {
            return;
          }
          if (canvas.width !== 0) {
            this.placeholderAspectRatio = canvas.height / canvas.width;
          }
          this.imageURL = window.URL.createObjectURL(blob);
        }, "image/png");
      } catch (e) {
        this.generateImageErrored = true;
        this.generateImageError = e.stack;
        console.error(`error generating image: ${e}`);
      }
    },

    copyErrorCode() {
      copyTextToClipboard(this.encodedError);
    },
  },
};
</script>

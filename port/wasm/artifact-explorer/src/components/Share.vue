<template>
  <div ref="shareWidget" class="relative" data-html2canvas-ignore>
    <div class="cursor-pointer" @click="open = !open">
      <svg
        xmlns="http://www.w3.org/2000/svg"
        viewBox="0 0 20 20"
        fill="currentColor"
        class="h-4 w-4 text-gray-500"
      >
        <path
          d="M15 8a3 3 0 10-2.977-2.63l-4.94 2.47a3 3 0 100 4.319l4.94 2.47a3 3 0 10.895-1.789l-4.94-2.47a3.027 3.027 0 000-.74l4.94-2.47C13.456 7.68 14.19 8 15 8z"
        />
      </svg>
    </div>

    <div
      class="origin-top-right absolute right-0 mt-2 w-48 rounded-md shadow-lg bg-white ring-1 ring-black ring-opacity-5 transform select-none"
      :class="open ? ['opacity-100', 'scale-100'] : ['opacity-0', 'scale-95']"
    >
      <div class="py-1" role="menu" aria-orientation="vertical" aria-labelledby="options-menu">
        <div
          class="block px-4 py-1 text-sm text-gray-700 hover:bg-gray-100 hover:text-gray-900 cursor-pointer"
          role="menuitem"
          @click="copyLink()"
        >
          Copy Link
        </div>
        <div
          class="block px-4 py-1 text-sm text-gray-700 hover:bg-gray-100 hover:text-gray-900 cursor-pointer"
          role="menuitem"
          @click="generateImageCard()"
        >
          Generate Image Card
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import copyTextToClipboard from "copy-text-to-clipboard";
import html2canvas from "html2canvas";

export default {
  props: {
    id: {
      type: String,
      required: true,
    },
    domElementId: {
      type: String,
      required: true,
    },
  },

  data() {
    return {
      open: false,
    };
  },

  mounted() {
    document.addEventListener("click", this.click);
  },

  beforeUnmount() {
    document.removeEventListener("click", this.click);
  },

  methods: {
    click(event) {
      if (!this.$refs.shareWidget.contains(event.target)) {
        this.open = false;
      }
    },

    copyLink() {
      copyTextToClipboard(`https://ei.tcl.sh/${this.id}`);
      this.open = false;
    },

    async generateImageCard() {
      try {
        const canvas = await html2canvas(document.getElementById(this.domElementId), {
          allowTaint: false,
          useCORS: true,
          backgroundColor: null,
          // Fix possible scrollbar-induced shift.
          scrollX: -window.scrollX,
          scrollY: -window.scrollY,
          windowWidth: document.documentElement.offsetWidth,
          windowHeight: document.documentElement.offsetHeight,
        });
        canvas.toBlob(blob => {
          const blobURL = window.URL.createObjectURL(blob);
          // The naive approach of window.open(blobURL, "_blank") doesn't work on iOS Safari.
          // Use the suggestion in https://dev.to/nombrekeff/download-file-from-blob-21ho,
          // create a link and trigger a click.
          const link = document.createElement("a");
          link.href = blobURL;
          link.download = `${this.id}.png`;
          link.style.display = "none";
          document.body.appendChild(link);
          link.dispatchEvent(
            new MouseEvent("click", {
              bubbles: true,
              cancelable: true,
              view: window,
            })
          );
          document.body.removeChild(link);
          this.open = false;
        }, "image/png");
      } catch (e) {
        console.error(e);
      }
    },
  },
};
</script>

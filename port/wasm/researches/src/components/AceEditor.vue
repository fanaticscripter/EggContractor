<template>
  <div ref="editorRef" class="absolute h-full w-full border border-gray-300 rounded-md"></div>
</template>

<script lang="ts">
import {
  defineComponent,
  onBeforeUnmount,
  onMounted,
  PropType,
  ref,
  Ref,
  toRefs,
  watch,
} from "vue";
import { Ace, config, edit } from "ace-builds";
import "ace-builds/src-noconflict/mode-json";
import "ace-builds/src-noconflict/mode-sql";
import "ace-builds/src-noconflict/ext-searchbox";
import "ace-builds/src-noconflict/ext-language_tools";
import AceWorkerJsonInline from "ace-builds/src-noconflict/worker-json.js?url";
import { Emitter } from "mitt";

config.setModuleUrl("ace/mode/json_worker", AceWorkerJsonInline);

// Extend the signature of foldAll.
interface Folding extends Ace.Folding {
  foldAll(
    startRow?: number,
    endRow?: number,
    depth?: number,
    test?: (row: number) => boolean
  ): void;
}

interface FoldingPrivate {
  getFoldWidgetRange(row: number): Ace.Range;
}

export default defineComponent({
  props: {
    modelValue: {
      type: String,
      required: true,
    },
    lang: {
      type: String,
      required: true,
    },
    readonly: {
      type: Boolean,
      default: false,
    },
    foldAtIndentation: {
      type: Number,
      required: false,
    },
    commands: {
      type: Array as PropType<Ace.Command[]>,
      default: [],
    },
    eventBus: {
      // The component listens for getValue events on the event bus, and emits
      // update:modelValue correspondingly.
      type: Object as PropType<Emitter>,
      required: false,
    },
  },
  emits: {
    "update:modelValue": (value: string) => true,
  },
  setup(props, { emit }) {
    const { modelValue, lang, readonly, foldAtIndentation, commands, eventBus } = toRefs(props);
    const editorRef: Ref<HTMLElement | null> = ref(null);

    let editor: Ace.Editor | null = null;

    onMounted(() => {
      editor = edit(editorRef.value!);
      // editor.renderer.setScrollMargin(16, 16, 16, 16);
      editor.setReadOnly(readonly.value);
      editor.setOptions({
        tabSize: 2,
        enableBasicAutocompletion: true,
        enableSnippets: true,
        enableLiveAutocompletion: true,
      });
      const session = editor.session;
      session.setMode(`ace/mode/${lang.value}`);
      session.setUseWrapMode(true);
      session.setValue(modelValue.value);
      commands.value.forEach(command => editor?.commands.addCommand(command));
      fold();

      eventBus?.value?.on("getValue", () => {
        const value = editor?.session.getValue();
        if (value !== undefined) {
          emit("update:modelValue", value);
        }
      });
    });

    const fold = () => {
      if (editor && foldAtIndentation) {
        const session = editor.session;
        const targetIndentation = foldAtIndentation.value;
        // https://github.com/ajaxorg/ace/blob/bfde34b510b263bd5ffa47b20445377cabe85dfb/lib/ace/edit_session/folding.js#L639-L666
        // https://groups.google.com/g/ace-discuss/c/JfMdCm1K8Qc
        (session as Folding).foldAll(undefined, undefined, undefined, row => {
          const range = ((session as unknown) as FoldingPrivate).getFoldWidgetRange(row);
          return range && range.isMultiLine() && range.end.column === targetIndentation;
        });
      }
    };

    watch(modelValue, () => {
      if (editor) {
        const cursor = editor.getCursorPosition();
        editor.session.setValue(modelValue.value);
        editor.moveCursorToPosition(cursor);
        fold();
      }
    });
    watch(lang, () => {
      editor?.session.setMode(`ace/mode/${lang.value}`);
    });
    watch(readonly, () => {
      editor?.setReadOnly(readonly.value);
    });

    onBeforeUnmount(() => {
      editor?.destroy();
    });

    return {
      editorRef,
    };
  },
});
</script>

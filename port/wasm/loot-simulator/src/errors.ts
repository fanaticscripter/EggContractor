export class ModuleWorkerNotSupportedError extends Error {
  constructor() {
    super(
      "Browser does not support module workers. https://caniuse.com/mdn-api_worker_worker_ecmascript_modules"
    );
    this.name = "ModuleWorkerNotSupportedError";
  }
}

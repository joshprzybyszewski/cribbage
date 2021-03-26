const go = new Go();
// memoryBytes is an Uint8Array pointing to the webassembly linear memory.
let memoryBytes;
let mod, inst, bytes;
console.log("initializing gowasm");
WebAssembly.instantiateStreaming(
  fetch("/assets/wasm/wa_output.wasm",
    // change default to 'no-cache' while deving to force getting it
    { cache: 'default' }),
  go.importObject).then((result) => {
    mod = result.module;
    inst = result.instance;
    memoryBytes = new Uint8Array(inst.exports.mem.buffer)
    console.log("initialized gowasm");
    run();
  });
async function run() {
  await go.run(inst);
}
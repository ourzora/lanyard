var __create = Object.create;
var __defProp = Object.defineProperty;
var __getOwnPropDesc = Object.getOwnPropertyDescriptor;
var __getOwnPropNames = Object.getOwnPropertyNames;
var __getProtoOf = Object.getPrototypeOf;
var __hasOwnProp = Object.prototype.hasOwnProperty;
var __export = (target, all) => {
  for (var name in all)
    __defProp(target, name, { get: all[name], enumerable: true });
};
var __copyProps = (to, from, except, desc) => {
  if (from && typeof from === "object" || typeof from === "function") {
    for (let key of __getOwnPropNames(from))
      if (!__hasOwnProp.call(to, key) && key !== except)
        __defProp(to, key, { get: () => from[key], enumerable: !(desc = __getOwnPropDesc(from, key)) || desc.enumerable });
  }
  return to;
};
var __reExport = (target, mod, secondTarget) => (__copyProps(target, mod, "default"), secondTarget && __copyProps(secondTarget, mod, "default"));
var __toESM = (mod, isNodeMode, target) => (target = mod != null ? __create(__getProtoOf(mod)) : {}, __copyProps(
  isNodeMode || !mod || !mod.__esModule ? __defProp(target, "default", { value: mod, enumerable: true }) : target,
  mod
));
var __toCommonJS = (mod) => __copyProps(__defProp({}, "__esModule", { value: true }), mod);
var __async = (__this, __arguments, generator) => {
  return new Promise((resolve, reject) => {
    var fulfilled = (value) => {
      try {
        step(generator.next(value));
      } catch (e) {
        reject(e);
      }
    };
    var rejected = (value) => {
      try {
        step(generator.throw(value));
      } catch (e) {
        reject(e);
      }
    };
    var step = (x) => x.done ? resolve(x.value) : Promise.resolve(x.value).then(fulfilled, rejected);
    step((generator = generator.apply(__this, __arguments)).next());
  });
};
var lib_exports = {};
__export(lib_exports, {
  CreateTree: () => CreateTree,
  GetProof: () => GetProof,
  GetRoots: () => GetRoots,
  GetTree: () => GetTree
});
module.exports = __toCommonJS(lib_exports);
var import_isomorphic_fetch = __toESM(require("isomorphic-fetch"));
__reExport(lib_exports, require("./types"), module.exports);
const baseUrl = "https://lanyard.org/api/v1/";
const client = (path, method, data) => __async(void 0, null, function* () {
  const opts = {
    method
  };
  if (method !== "GET") {
    opts.body = JSON.stringify(data);
    opts.headers = {
      "Content-Type": "application/json"
    };
  }
  const resp = yield (0, import_isomorphic_fetch.default)(baseUrl + path, opts);
  return resp.json();
});
const CreateTree = (req) => {
  return client("tree", "POST", req);
};
const GetTree = (merkleRoot) => {
  return client(`tree?root=${merkleRoot}`, "GET");
};
const GetProof = (root, unhashedLeaf) => {
  return client(`proof?root=${root}&leaf=${unhashedLeaf}`, "GET");
};
const GetRoots = (proof) => {
  return client(`roots?proof=${proof}`, "GET");
};
// Annotate the CommonJS export names for ESM import in node:
0 && (module.exports = {
  CreateTree,
  GetProof,
  GetRoots,
  GetTree
});

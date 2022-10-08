import fetch from "isomorphic-fetch";
import {
  CreateTreeRequest,
  CreateTreeResponse,
  GetProofResponse,
  GetRootsResponse,
  GetTreeResponse,
} from "./types";
const baseUrl = "https://lanyard.org/api/v1/";

const client = async (path: string, method: string, data?: any) => {
  const opts: {
    method: string;
    body?: string;
    headers?: any;
  } = {
    method: method,
  };

  if (method !== "GET") {
    opts.body = JSON.stringify(data);
    opts.headers = {
      "Content-Type": "application/json",
    };
  }

  const resp = await fetch(baseUrl + path, opts);
  return resp.json();
};

export const CreateTree = (
  req: CreateTreeRequest
): Promise<CreateTreeResponse> => {
  return client("tree", "POST", req);
};

export const GetTree = (merkleRoot: string): Promise<GetTreeResponse> => {
  return client(`tree?root=${merkleRoot}`, "GET");
};

export const GetProof = (
  root: string,
  unhashedLeaf: string
): Promise<GetProofResponse> => {
  return client(`proof?root=${root}&leaf=${unhashedLeaf}`, "GET");
};

export const GetRoots = (proof: string): Promise<GetRootsResponse> => {
  return client(`roots?proof=${proof}`, "GET");
};

export * from "./types";

import { CreateTreeRequest, CreateTreeResponse, GetProofResponse, GetRootsResponse, GetTreeResponse } from "./types";
export declare const CreateTree: (req: CreateTreeRequest) => Promise<CreateTreeResponse>;
export declare const GetTree: (merkleRoot: string) => Promise<GetTreeResponse>;
export declare const GetProof: (root: string, unhashedLeaf: string) => Promise<GetProofResponse>;
export declare const GetRoots: (proof: string) => Promise<GetRootsResponse>;
export * from "./types";

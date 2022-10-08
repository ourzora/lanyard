export interface CreateTreeRequest {
    unhashedLeaves: string[];
    leafTypeDescriptor?: string[];
    packedEncoding?: boolean;
}
export interface CreateTreeResponse {
    merkleRoot: string;
}
export interface GetTreeResponse {
    unhashedLeaves: string[];
    leafCount: number;
    leafTypeDescriptor: string[] | null;
    packedEncoding: boolean | null;
}
export interface GetProofResponse {
    proof: string[];
    unhashedLeaf: string;
}
export interface GetRootsResponse {
    roots: string[];
}

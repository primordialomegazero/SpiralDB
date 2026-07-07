declare module '@primordialomegazero/spiraldb-client' {
  interface PutResponse {
    status: string;
    action: string;
    key: string;
    value: number;
    mirrors: number;
    fractal_layers: number;
    fhe_encrypted: boolean;
    double_mirror: boolean;
  }

  interface GetResponse {
    status: string;
    action: string;
    key: string;
    value: number;
    mirrors: number;
    fhe_encrypted: boolean;
    double_mirror: boolean;
  }

  interface FHEComputeResponse {
    status: string;
    action: string;
    op: 'add' | 'mul';
    encrypted_result: string;
    decrypted_result: number;
    expected_plaintext: number;
    correct: boolean;
    computation_blind: boolean;
    self_referential: boolean;
  }

  interface MirrorHealthResponse {
    status: string;
    key: string;
    in_primary: boolean;
    in_cache: boolean;
    in_fractal: boolean;
    mirrors_synced: boolean;
  }

  interface HealthResponse {
    status: string;
    fhe: string;
  }

  class SpiralDBClient {
    constructor(host?: string);
    
    put(key: string, value: number): Promise<PutResponse>;
    get(key: string): Promise<GetResponse>;
    fheAdd(key1: string, key2: string): Promise<FHEComputeResponse>;
    fheMul(key1: string, key2: string): Promise<FHEComputeResponse>;
    mirrorHealth(key: string): Promise<MirrorHealthResponse>;
    health(): Promise<HealthResponse>;
  }

  export { SpiralDBClient, PutResponse, GetResponse, FHEComputeResponse, MirrorHealthResponse, HealthResponse };
}

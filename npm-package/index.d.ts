declare module '@primordialomegazero/spiraldb-client' {
  class SpiralDBClient {
    constructor(host?: string, port?: number);
    put(key: string, value: number): Promise<any>;
    get(key: string): Promise<any>;
    fheAdd(key1: string, key2: string): Promise<any>;
    fheMul(key1: string, key2: string): Promise<any>;
    mirrorHealth(key: string): Promise<any>;
    health(): Promise<any>;
    encrypt(value: number): number;
    decrypt(encrypted: number): number;
  }
  export { SpiralDBClient };
}

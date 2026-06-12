export {};

declare module "*.js";

declare global {
  var Go: {
    new (): {
      importObject: WebAssembly.Imports;
      run(instance: WebAssembly.Instance): Promise<void>;
    };
  };

  // eslint-disable-next-line no-var
  var comicalPolicy:
    | {
        expiryUnix(ttl: string, fallbackSeconds: number): number;
        expiredUnix(timestamp: number): boolean;
        randomSlug(): string;
        validateSlug(slug: string): boolean;
        visitLimitExceeded(maxVisits: number, visitCount: number): boolean;
      }
    | undefined;
}

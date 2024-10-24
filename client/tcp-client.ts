import net from "node:net";
import { Buffer } from "node:buffer";

class MemoryServer {
  constructor(port: number, host: string) {
    this.port = port;
    this.host = host;
  }

  private port: number;
  private host: string;
  private client: undefined | net.Socket;
  private pendingRequests: Map<string, (val: string) => void> = new Map();
  connected = false;

  async connect() {
    return new Promise((resolve, reject) => {
      if (this.connected) return resolve("already connected");
      const socket = net.createConnection({
        port: this.port,
        host: this.host,
      });
      socket.on("connect", () => {
        this.client = socket;
        this.connected = true;
        resolve("connected");
        this.listenConnectionEvents();
      });
      socket.on("error", (err) => {
        if (err.name === "AggregateError") {
          return reject("Memory server is down, unable to connect");
        }
        reject(err);
      });
    });
  }

  async disconnect() {
    return new Promise((resolve, reject) => {
      if (this.connected) {
        this.client?.end(() => {
          this.client = undefined;
          this.connected = false;
          return resolve("disconnected");
        });
      } else {
        return resolve("already disconnected from memory-server");
      }
    });
  }

  private listenConnectionEvents() {
    this.client?.on("data", (data) => {
      const id = data.subarray(0, 36).toString();
      const content = data.subarray(36).toString();

      if (this.pendingRequests.has(id)) {
        const resolve = this.pendingRequests.get(id);
        resolve?.(content);
      } else {
        console.log(`request ${id} not found`);
      }
    });
  }

  private generateRequestId() {
    return crypto.randomUUID();
  }

  private async handleRequest(cmd: string) {
    if (!this.client || !this.connected) {
      throw new Error(`tried to run ${cmd} but memory-serer is disconnected`);
    }

    const requestId = this.generateRequestId();
    const buff = Buffer.concat([
      Buffer.from(requestId),
      Buffer.from(" "),
      Buffer.from(cmd),
    ]);

    this.client.write(buff);

    const responsePromise = new Promise((resolve: (val: string) => void) => {
      this.pendingRequests.set(requestId, resolve);
    });

    const response = await responsePromise;
    this.pendingRequests.delete(requestId);

    return response;
  }

  async set(key: string, value: string, opts?: Record<string, string>) {
    let strOpts = "";

    if (opts) {
      Object.entries(opts).forEach(([k, v]) => {
        if (v) {
          strOpts += ` ${k} ${v}`;
        } else {
          strOpts += ` ${k}`;
        }
      });
    }

    return this.handleRequest(`SET ${key} ${value} ${strOpts}`);
  }

  async get(key: string) {
    const val = await this.handleRequest(`GET ${key}`);
    if (val === "<nil>") return "nil";
    return val;
  }

  async del(key: string | string[]) {
    let keys = key;
    if (Array.isArray(key)) {
      keys = key.map((k) => k.trim()).join(" ");
    }
    return this.handleRequest(`DEL ${keys}`);
  }
}

export default MemoryServer;

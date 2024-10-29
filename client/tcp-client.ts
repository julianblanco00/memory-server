import net from "node:net";
import { Buffer } from "node:buffer";

type RequestFunctions = {
  resolve: (v?: unknown) => void;
  reject: (r?: Error) => void;
};

class MemoryServer {
  constructor(port: number, host: string) {
    this.port = port;
    this.host = host;
  }

  private port: number;
  private host: string;
  private client: undefined | net.Socket;
  private pendingRequests: Map<string, RequestFunctions> = new Map();
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
      const id = data.subarray(0, 16).toString("hex");
      const content = data.subarray(16).toString();

      if (this.pendingRequests.has(id)) {
        const req = this.pendingRequests.get(id);
        if (content === "wrong number of arguments") {
          return req?.reject(new Error(content));
        }
        req?.resolve(content);
      } else {
        console.log(`request ${id} not found`);
      }
    });
  }

  private buildRESPCommand(...vals: string[]): string {
    let str = "";

    vals.forEach((val) => {
      str += `$${val.length}\n${val}\n`;
    });

    return str;
  }

  private generateRequestId() {
    return Buffer.from(crypto.randomUUID().replace(/-/g, ""), "hex");
  }

  private async handleRequest(cmd: string) {
    if (!this.client || !this.connected) {
      throw new Error(`tried to run ${cmd} but memory-serer is disconnected`);
    }

    const requestId = this.generateRequestId();
    const buff = Buffer.concat([requestId, Buffer.from(" "), Buffer.from(cmd)]);

    this.client.write(buff);

    const requestIdStr = requestId.toString("hex");

    const responsePromise = new Promise((resolve, reject) => {
      this.pendingRequests.set(requestIdStr, { resolve, reject });
    });

    const response = await responsePromise;
    this.pendingRequests.delete(requestIdStr);

    return response;
  }

  async set(key: string, value: string, opts?: string[][]) {
    return this.handleRequest(
      this.buildRESPCommand("SET", key, value, ...(opts?.flat() || [])),
    );
  }

  async get(key: string) {
    const val = await this.handleRequest(this.buildRESPCommand("GET", key));
    return val === "<nil>" ? "nil" : val;
  }

  async del(key: string | string[]) {
    const k = typeof key === "string" ? [key] : key;
    return this.handleRequest(this.buildRESPCommand("DEL", ...k));
  }

  async mSet(...params: string[]) {
    if (params.length % 2 === 1) {
      throw new Error("wrong number of arguments");
    }
    return this.handleRequest(this.buildRESPCommand("MSET", ...params));
  }

  async hSet(key: string, ...params: string[] | Record<string, string>[]) {
    let cmdParams: string[] = [key];

    if (typeof params[0] === "object") {
      Object.entries(params[0]).forEach(([k, v]) => {
        cmdParams.push(k, v);
      });
    } else if (typeof params[0] === "string") {
      if (params.length % 2 === 1) {
        throw new Error("wrong number of arguments");
      }
      cmdParams.push(...(params as unknown as string));
    } else {
      throw new Error("invalid input");
    }

    return this.handleRequest(this.buildRESPCommand("HSET", ...cmdParams));
  }

  async hGet(key: string, field: string) {
    if (!key || !field) {
      throw new Error("missing fields for command hGet");
    }
    return this.handleRequest(this.buildRESPCommand("HGET", key, field));
  }

  async hDel(key: string, ...fields: string[]) {
    if (!key || !fields.length) {
      throw new Error("missing fields for command hDel");
    }
    return this.handleRequest(this.buildRESPCommand("HDEL", key, ...fields));
  }
}

export default MemoryServer;

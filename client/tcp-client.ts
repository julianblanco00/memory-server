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
    let strOpts = "";

    if (opts) {
      opts.forEach(([k, v]) => {
        if (v) {
          strOpts += `$${k.length}\n${k}\n$${v.length}\n${v}\n`;
        } else {
          strOpts += `$${k.length}\n${k}\n`;
        }
      });
    }

    return this.handleRequest(
      `$3\nSET\n$${key.length}\n${key}\n$${value.length}\n${value}\n${strOpts}`,
    );
  }

  async mset(...params: string[]) {
    let str = "$4\nMSET\n";

    if (params.length % 2 === 1) {
      throw new Error("wrong number of arguments");
    }

    let i = 0;

    while (i < params.length) {
      const k = params[i];
      i++;
      const v = params[i];
      i++;
      str += `$${k.length}\n${k}\n$${v?.length ?? 0}\n${v}\n`;
    }

    return this.handleRequest(str);
  }

  async hset(key: string, ...params: string[] | Record<string, string>[]) {
    let str = `$4\nHSET\n$${key.length}\n${key}\n`;

    if (typeof params[0] === "object") {
      Object.entries(params[0]).forEach(([k, v]) => {
        str += `$${k.length}\n${k}\n`;
        str += `$${v.length}\n${v}\n`;
      });
    } else if (typeof params[0] === "string") {
      if (params.length % 2 === 1) {
        throw new Error("wrong number of arguments");
      }

      let i = 0;

      while (i < params.length) {
        const k = params[i];
        i++;
        const v = params[i];
        i++;
        str += `$${k.length}\n${k}\n$${v?.length ?? 0}\n${v}\n`;
      }
    } else {
      throw new Error("invalid input");
    }

    return this.handleRequest(str);
  }

  async get(key: string) {
    const val = await this.handleRequest(`$3\nGET\n$${key.length}\n${key}\n`);
    if (val === "<nil>") return "nil";
    return val;
  }

  async del(key: string | string[]) {
    let keys = `$${key.length}\n${key}`;
    if (Array.isArray(key)) {
      keys = key.map((k) => `\n$${k.length}\n${k}`).join(" ");
    }
    return this.handleRequest(`$3\nDEL ${keys}`);
  }
}

export default MemoryServer;

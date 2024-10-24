const net = require("net");
const { Buffer } = require("buffer");

class MemoryServer {
  constructor(port, host) {
    this.client = net.createConnection({ port, host });

    this.client.on("data", (data) => {
      const id = data.slice(0, 36).toString();
      const content = data.slice(36).toString();

      if (this.pendingRequests.has(id)) {
        const resolve = this.pendingRequests.get(id);
        resolve(content);
      } else {
        console.log(`request ${id} not found`);
      }
    });

    this.client.on("end", () => {
      console.log("disconnected from server");
    });
  }

  client;
  pendingRequests = new Map();

  generateRequestId() {
    return crypto.randomUUID();
  }

  async handleRequest(cmd) {
    const requestId = this.generateRequestId();
    const buff = Buffer.concat([
      Buffer.from(requestId),
      Buffer.from(" "),
      Buffer.from(cmd),
    ]);

    this.client.write(buff);

    const responsePromise = new Promise((resolve) => {
      this.pendingRequests.set(requestId, resolve);
    });

    const response = await responsePromise;
    this.pendingRequests.delete(requestId);

    return response;
  }

  async set(key, value, opts) {
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

  async get(key) {
    return this.handleRequest(`GET ${key}`);
  }

  async del(key) {
    let keys = key;
    if (Array.isArray(key)) {
      keys = key.map((k) => k.trim()).join(" ");
    }
    return this.handleRequest(`DEL ${keys}`);
  }
}

async function main() {
  const memoryServer = new MemoryServer(4444, "localhost");

  const resp1 = await memoryServer.set("mykey", "valuefromnode");
  console.log({ resp1 });
  await new Promise((resolve) => setTimeout(resolve, 2000));
  const resp2 = await memoryServer.get("mykey");
  console.log({ resp2 });
  const resp3 = await memoryServer.del(["mykey", "mykey2"]);
  console.log({ resp3 });
  const resp4 = await memoryServer.get("mykey");
  console.log({ resp4 });
}

main();

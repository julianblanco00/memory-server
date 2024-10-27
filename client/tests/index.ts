import { afterEach, beforeEach, describe, it } from "node:test";
import assert from "node:assert";
import MemoryServer from "../tcp-client.js";

const memoryServer = new MemoryServer(4444, "localhost");

beforeEach(async () => {
  await memoryServer.connect();
});

afterEach(async () => {
  await memoryServer.disconnect();
});

describe("memory-server tests", () => {
  it("can set a value in a key", async () => {
    const resp = await memoryServer.set("mykey", "myvalue is the\n best", [
      ["get"],
      ["exat", "120"],
    ]);
    assert.equal(resp, "OK", "error setting key");
  });

  it("can get a key with value", async () => {
    const resp = await memoryServer.get("mykey");
    assert.equal(resp, "myvalue is the\n best", "error getting key");
  });

  it("can delete a key", async () => {
    const resp = await memoryServer.del("mykey");
    assert.equal(resp, "1", "error deleting key");
  });

  it("can get a key with no value", async () => {
    const resp = await memoryServer.get("mykey");
    assert.equal(resp, "nil", "error getting key");
  });

  it("can try to delete non-existing keys", async () => {
    const resp = await memoryServer.del(["mykey", "mykey1"]);
    assert.equal(resp, "0", "error deleting key");
  });

  it("should fail trying to get a key without passing a key", async () => {
    const resp = await memoryServer.get("    v");
    assert.equal(resp, "nil", "error getting key");
  });
});

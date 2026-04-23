import { environment } from "@/core/environment";
import { RedisClient } from "bun";

const { HOST, PORT, PASSWORD, DB } = environment.REDIS;
const url = `redis://${PASSWORD ? `:${PASSWORD}@` : ""}${HOST}:${PORT}/${DB ?? 0}`;

export interface IRedisService {
  get<T>(key: string): Promise<T | null>;
  set<T>(key: string, value: T, expirationSeconds: number): Promise<void>;
  clear(): Promise<void>;
}

class RedisService implements IRedisService {
  constructor(private readonly client: RedisClient) {}

  async get<T>(key: string): Promise<T | null> {
    const val = await this.client.get(key);
    if (val === null) return null;
    return JSON.parse(val) as T;
  }

  async set<T>(
    key: string,
    value: T,
    expirationSeconds: number,
  ): Promise<void> {
    await this.client.set(key, JSON.stringify(value), "EX", expirationSeconds);
  }

  async clear(): Promise<void> {
    await this.client.send("FLUSHDB", []);
  }
}

export const redisClient = new RedisClient(url);
export const redisService = new RedisService(redisClient);

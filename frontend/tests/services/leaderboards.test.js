import createFetchMock from "vitest-fetch-mock";
import { describe, vi } from "vitest";

import { getLeaderboard } from "../../src/services/leaderboards.js";

const BACKEND_URL = import.meta.env.VITE_BACKEND_URL;

describe("leaderboards service", () => {})
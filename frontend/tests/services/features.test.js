import createFetchMock from "vitest-fetch-mock";
import { describe, vi } from "vitest";

import { getAllFeatures, getFeatures } from "../../src/services/features.js";

const BACKEND_URL = import.meta.env.VITE_BACKEND_URL;

describe("feature service", () => {})
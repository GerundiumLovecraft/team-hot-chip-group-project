import createFetchMock from "vitest-fetch-mock";
import { describe, vi } from "vitest";

import { getAllSpots, getSpotById, getSpotsByUser, getSpotByFeature, createSpot} from "../../src/services/spots.js";

const BACKEND_URL = import.meta.env.VITE_BACKEND_URL;

describe("spots service", () => {})
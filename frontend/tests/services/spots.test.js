import createFetchMock from "vitest-fetch-mock";
import { describe, vi } from "vitest";

import { getAllSpots, getSpotById, getSpotsByUser, getSpotByFeature, createSpot } from "../../src/services/spots.js";

const BACKEND_URL = import.meta.env.VITE_BACKEND_URL;

createFetchMock(vi).enableMocks();

const mockSpots = [
    { _id: 1, name: "Spot One", address: "1 Test St" },
    { _id: 2, name: "Spot Two", address: "2 Test St" },
];

const mockSpot = { _id: 1, name: "Spot One", address: "1 Test St" };

const testToken = "testToken123";

describe("spots service", () => {

    describe("getAllSpots", () => {
        test("calls the correct url with the correct method", async () => {
            fetch.mockResponseOnce(JSON.stringify({ spots: mockSpots }), {
                status: 200,
            });

            await getAllSpots();

            const [url, options] = fetch.mock.lastCall;

            expect(url).toEqual(`${BACKEND_URL}/spots`);
            expect(options.method).toEqual("GET");
            expect(options.headers["Content-Type"]).toEqual("application/json");
        });

        test("returns array of spots on success", async () => {
            fetch.mockResponseOnce(JSON.stringify({ spots: mockSpots }), {
                status: 200,
            });

            const data = await getAllSpots();

            expect(data.spots).toHaveLength(2);
            expect(data.spots[0].name).toEqual("Spot One");
        });

        test("throws an error if request fails", async () => {
            fetch.mockResponseOnce(JSON.stringify({ message: "Server error" }), {
                status: 500,
            });

            await expect(getAllSpots()).rejects.toThrow("Unable to fetch all spots!");
        });
    });

    describe("getSpotById", () => {
        test("calls the correct url with the correct method", async () => {
            fetch.mockResponseOnce(JSON.stringify({ Spot: mockSpot }), {
                status: 200,
            });

            await getSpotById(1);

            const [url, options] = fetch.mock.lastCall;

            expect(url).toEqual(`${BACKEND_URL}/spots/1`);
            expect(options.method).toEqual("GET");
            expect(options.headers["Content-Type"]).toEqual("application/json");
        });

        test("returns a spot on success", async () => {
            fetch.mockResponseOnce(JSON.stringify({ Spot: mockSpot }), {
                status: 200,
            });

            const spot = await getSpotById(1);

            expect(spot.name).toEqual("Spot One");
            expect(spot._id).toEqual(1);
        });

        test("throws an error if request fails", async () => {
            fetch.mockResponseOnce(JSON.stringify({ message: "Not found" }), {
                status: 404,
            });

            await expect(getSpotById(9999)).rejects.toThrow("Unable to find that spot!");
        });
    });

    describe("getSpotsByUser", () => {
        test("calls the correct url with the correct method", async () => {
            fetch.mockResponseOnce(JSON.stringify({ spots: mockSpots }), {
                status: 200,
            });

            await getSpotsByUser(testToken);

            const [url, options] = fetch.mock.lastCall;

            expect(url).toEqual(`${BACKEND_URL}/profile/spots`);
            expect(options.method).toEqual("GET");
            expect(options.headers["Content-Type"]).toEqual("application/json");
            expect(options.headers["Authorization"]).toEqual(`Bearer ${testToken}`);
        });

        test("returns an array of spots on success", async () => {
            fetch.mockResponseOnce(JSON.stringify({ spots: mockSpots }), {
                status: 200,
            });

            const spots = await getSpotsByUser(testToken);

            expect(spots).toHaveLength(2);
            expect(spots[0].name).toEqual("Spot One");
        });

        test("throws an error if unauthorised", async () => {
            fetch.mockResponseOnce(JSON.stringify({ message: "Unauthorised" }), {
                status: 401,
            });

            await expect(getSpotsByUser("badtoken")).rejects.toThrow(
                "Unauthorised: Please log in again."
            );
        });

        test("throws an error if request fails", async () => {
            fetch.mockResponseOnce(JSON.stringify({ message: "Server error" }), {
                status: 500,
            });

            await expect(getSpotsByUser(testToken)).rejects.toThrow(
                "Unable to fetch your submitted spots!"
            );
        });
    });

    describe("getSpotByFeature", () => {
        const selectedFeatures = [{ feat_id: 1, value: true }];

        test("calls the correct url with the correct method", async () => {
            fetch.mockResponseOnce(JSON.stringify({ spots: mockSpots }), {
                status: 200,
            });

            await getSpotByFeature(selectedFeatures);

            const [url, options] = fetch.mock.lastCall;

            expect(url).toEqual(`${BACKEND_URL}/spots/filter`);
            expect(options.method).toEqual("POST");
            expect(options.headers["Content-Type"]).toEqual("application/json");
            expect(options.body).toEqual(JSON.stringify(selectedFeatures));
        });

        test("returns an array of spots on success", async () => {
            fetch.mockResponseOnce(JSON.stringify({ spots: mockSpots }), {
                status: 200,
            });

            const data = await getSpotByFeature(selectedFeatures);

            expect(data.spots).toHaveLength(2);
            expect(data.spots[0].name).toEqual("Spot One");
        });

        test("throws an error if request fails", async () => {
            fetch.mockResponseOnce(JSON.stringify({ message: "Server error" }), {
                status: 500,
            });

            await expect(getSpotByFeature(selectedFeatures)).rejects.toThrow(
                "Unable to fetch the filter spots!"
            );
        });
    });

    describe("createSpot", () => {
        const spotData = {
            name: "New Spot",
            address: "3 Test St",
            description: "A test spot",
            open_from: "09:00",
            open_to: "18:00",
            features: [{ feat_id: 1 }],
        };

        test("calls the correct url with the correct method", async () => {
            fetch.mockResponseOnce(JSON.stringify({ spotID: 3 }), { status: 201 });

            await createSpot(spotData, testToken);

            const [url, options] = fetch.mock.lastCall;

            expect(url).toEqual(`${BACKEND_URL}/spots`);
            expect(options.method).toEqual("POST");
            expect(options.headers["Content-Type"]).toEqual("application/json");
            expect(options.headers["Authorization"]).toEqual(`Bearer ${testToken}`);
            expect(options.body).toEqual(JSON.stringify(spotData));
        });

        test("returns response data on success", async () => {
            fetch.mockResponseOnce(JSON.stringify({ spotID: 3 }), { status: 201 });

            const data = await createSpot(spotData, testToken);

            expect(data.spotID).toEqual(3);
        });

        test("throws an error if request fails", async () => {
            fetch.mockResponseOnce(JSON.stringify({ message: "Bad request" }), {
                status: 400,
            });

            await expect(createSpot(spotData, testToken)).rejects.toThrow(
                "Unable to create a spot!"
            );
        });
    });
});
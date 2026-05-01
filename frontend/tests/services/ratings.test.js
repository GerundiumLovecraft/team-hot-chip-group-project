import createFetchMock from "vitest-fetch-mock";
import { describe, vi } from "vitest";

import { SubmitRating } from "../../src/services/ratings.js";
import {isValidLocalDateAndTimeString} from "jsdom/lib/jsdom/living/helpers/dates-and-times.js";

const BACKEND_URL = import.meta.env.VITE_BACKEND_URL;

// Create fetch mock
createFetchMock(vi).enableMocks();

describe("ratings service", () => {
    describe("submitRating",  () => {
        test("calls the correct url with the correct method", async () => {
            fetch.mockResponseOnce();
            const token = "test-token";
            const spotId = 2;
            const value = 5;

            await SubmitRating(token, spotId, value);

            const fetchArguments = fetch.mock.lastCall;
            const url = fetchArguments[0]
            const options = fetchArguments[1]

            expect(url).toEqual(`${BACKEND_URL}/spots/${spotId}/rate`);
            expect(options.method).toEqual("POST");
            expect(options.headers["Content-Type"]).toEqual("application/json");
            expect(options.headers.Authorization).toEqual(`Bearer ${token}`);
            expect(options.body).toEqual(JSON.stringify({"rating": value}));
        });
        test("returns nothing on success", async () => {
            fetch.mockResponseOnce();
            const token = "test-token";
            const spotId = 2;
            const value = 5;

            const responseData = await SubmitRating(token, spotId, value);

            expect(responseData).toBeUndefined();
        });
        test("throws an error if request fails", async () => {
            fetch.mockResponseOnce(JSON.stringify({ message: "Server issue. Try again later." }), { status: 500 });
            const token = "test-token";
            const spotId = 2;
            const value = 5;

            try{
                await SubmitRating(token, spotId, value);
            } catch (e) {
                expect(e.message).toEqual("Couldn't add the rating. Try again later");
            }

        });
    })
})
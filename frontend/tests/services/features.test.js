import createFetchMock from "vitest-fetch-mock";
import { describe, vi } from "vitest";

import { getFeatures } from "../../src/services/features.js";

const BACKEND_URL = import.meta.env.VITE_BACKEND_URL;

// Create fetch mock
createFetchMock(vi).enableMocks();

describe("feature service", () => {
    describe("getFeatures", () => {
        test("calls the correct url with the correct method", async () => {
            fetch.mockResponseOnce(JSON.stringify({
                "message": "List of features",
                "features": [
                    { "feat_id": 1, "feat_name": "Wi-Fi" },
                    { "feat_id": 2, "feat_name": "Power Sockets" }
                ]
            }));

            await getFeatures();

            const fetchArguments = fetch.mock.lastCall;
            const url = fetchArguments[0];
            const options = fetchArguments[1];

            expect(url).toEqual(`${BACKEND_URL}/features`);
            expect(options.method).toEqual("GET");
            expect(options.headers["Content-Type"]).toEqual("application/json");
        });
        test("returns list of features on success", async () => {
            fetch.mockResponseOnce(JSON.stringify(
                {"message": "List of features", "features": [
                        {
                            "feat_id": 1,
                            "feat_name": "Wi-Fi"
                        },
                        {
                            "feat_id": 2,
                            "feat_name": "Power Sockets"
                        }
                    ]
                })
            );

            const response = await getFeatures();

            expect(response.message).toEqual("List of features");
            expect(response.features.length).toEqual(2);
            expect(response.features[1].feat_id).toEqual(2);

        })
        test("throws error if the request fails", async () => {
            fetch.mockResponseOnce(JSON.stringify({ message: "Server issue. Try again later." }), { status: 500 });

            try{
                await getFeatures();
            } catch (e) {
                expect(e.message).toEqual("Server issue. Try again later!");
            }
        })
    })
})
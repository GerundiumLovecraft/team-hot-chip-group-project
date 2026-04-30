import { render, screen, waitFor, fireEvent } from "@testing-library/react";
import { describe, it, expect, vi, beforeEach } from "vitest";

import { MemoryRouter } from "react-router-dom";

import { ProfilePage } from "../../src/pages/Profile/ProfilePage";

// ── Mock service modules ──
vi.mock("../../services/authentication", () => ({
    fetchUserProfile: vi.fn(),
}));

vi.mock("../../services/spots", () => ({
    getSpotsByUser: vi.fn(),
}));

vi.mock("../../helpers/authentication.js", () => ({
    isTokenValid: vi.fn(),
}));

import { fetchUserProfile } from "../../services/authentication";
import { getSpotsByUser } from "../../services/spots";
import { isTokenValid } from "../../helpers/authentication.js";

// ── Shared test data ──
const mockUser = {
    id: "1",
    username: "aluna",
    email: "aluna@gmail.com",
    createdAt: "2026-04-24T00:00:00Z",
    avatar: null,
};

const mockSpots = [
    { _id: 1, name: "Cosy Cove", address: "cosy grove, middlewood, LS3 6AT", image: "https://example.com/cosy.jpg", features: [] },
    { _id: 2, name: "Hot Spot", address: "127 Neo Drive, En City", image: "https://example.com/hot.jpg", features: [] },
];

const renderProfilePage = () =>
    render(
        <MemoryRouter>
            <ProfilePage />
        </MemoryRouter>
    );

beforeEach(() => {
    vi.clearAllMocks();
    localStorage.setItem("token", "mock-token");
    isTokenValid.mockReturnValue(true);
});

// ── Tests ──
describe("ProfilePage", () => {

    describe("authentication", () => {
        it("redirects to /login if token is invalid", async () => {
            isTokenValid.mockReturnValue(false);
            fetchUserProfile.mockResolvedValue({ user: mockUser, token: "new-token" });
            getSpotsByUser.mockResolvedValue([]);

            renderProfilePage();

            await waitFor(() => {
                expect(window.location.pathname).toBe("/login");
            });
        });
    });

    describe("loading states", () => {
        it("shows loading text while profile is fetching", () => {
            isTokenValid.mockReturnValue(true);
            fetchUserProfile.mockReturnValue(new Promise(() => {}));
            getSpotsByUser.mockReturnValue(new Promise(() => {}));

            renderProfilePage();

            expect(screen.getByText("Loading profile…")).toBeInTheDocument();
            expect(screen.getByText("Loading your spots…")).toBeInTheDocument();
        });
    });

    describe("profile data", () => {
        beforeEach(() => {
            fetchUserProfile.mockResolvedValue({ user: mockUser, token: "new-token" });
            getSpotsByUser.mockResolvedValue(mockSpots);
        });

        it("renders username and email after loading", async () => {
            renderProfilePage();

            await waitFor(() => {
                expect(screen.getByText("aluna")).toBeInTheDocument();
                expect(screen.getByText("aluna@gmail.com")).toBeInTheDocument();
            });
        });

        it("renders the member since date", async () => {
            renderProfilePage();

            await waitFor(() => {
                expect(screen.getByText("24 April 2026")).toBeInTheDocument();
            });
        });

        it("shows initials in avatar when no avatar image is set", async () => {
            renderProfilePage();

            await waitFor(() => {
                expect(screen.getByText("AL")).toBeInTheDocument();
            });
        });

        it("shows avatar image when user has one set", async () => {
            fetchUserProfile.mockResolvedValue({
                user: { ...mockUser, avatar: "https://example.com/avatar.jpg" },
                token: "new-token",
            });

            renderProfilePage();

            await waitFor(() => {
                const img = screen.getByAltText("avatar");
                expect(img).toBeInTheDocument();
                expect(img.src).toBe("https://example.com/avatar.jpg");
            });
        });
    });

    describe("error states", () => {
        it("shows an error message if profile fetch fails", async () => {
            fetchUserProfile.mockRejectedValue(new Error("Something went wrong"));
            getSpotsByUser.mockResolvedValue([]);

            renderProfilePage();

            await waitFor(() => {
                expect(screen.getByText("Something went wrong")).toBeInTheDocument();
            });
        });

        it("shows an error message if spots fetch fails", async () => {
            fetchUserProfile.mockResolvedValue({ user: mockUser, token: "new-token" });
            getSpotsByUser.mockRejectedValue(new Error("Unable to fetch your submitted spots!"));

            renderProfilePage();

            await waitFor(() => {
                expect(screen.getByText("Unable to fetch your submitted spots!")).toBeInTheDocument();
            });
        });
    });

    describe("submitted spots", () => {
        beforeEach(() => {
            fetchUserProfile.mockResolvedValue({ user: mockUser, token: "new-token" });
        });

        it("shows empty state message when user has no spots", async () => {
            getSpotsByUser.mockResolvedValue([]);

            renderProfilePage();

            await waitFor(() => {
                expect(screen.getByText("You have not submitted any spots yet.")).toBeInTheDocument();
            });
        });

        it("renders spot cards when user has submitted spots", async () => {
            getSpotsByUser.mockResolvedValue(mockSpots);

            renderProfilePage();

            await waitFor(() => {
                expect(screen.getByText("Cosy Cove")).toBeInTheDocument();
                expect(screen.getByText("Hot Spot")).toBeInTheDocument();
            });
        });

        it("opens the spot modal when a spot card is clicked", async () => {
            getSpotsByUser.mockResolvedValue(mockSpots);

            renderProfilePage();

            await waitFor(() => screen.getByText("Cosy Cove"));
            fireEvent.click(screen.getByText("Cosy Cove"));

            await waitFor(() => {
                expect(screen.getByText("cosy grove, middlewood, LS3 6AT")).toBeInTheDocument();
            });
        });
    });

    describe("avatar editing", () => {
        beforeEach(() => {
            fetchUserProfile.mockResolvedValue({ user: mockUser, token: "new-token" });
            getSpotsByUser.mockResolvedValue([]);
        });

        it("shows the avatar edit input when avatar is clicked", async () => {
            renderProfilePage();

            await waitFor(() => screen.getByText("AL"));
            fireEvent.click(screen.getByText("AL"));

            expect(screen.getByPlaceholderText("Paste image URL…")).toBeInTheDocument();
        });

        it("hides the avatar edit input when cancel is clicked", async () => {
            renderProfilePage();

            await waitFor(() => screen.getByText("AL"));
            fireEvent.click(screen.getByText("AL"));
            fireEvent.click(screen.getByText("Cancel"));

            expect(screen.queryByPlaceholderText("Paste image URL…")).not.toBeInTheDocument();
        });

        it("saves the avatar and updates the UI", async () => {
            const newAvatarUrl = "https://example.com/new-avatar.jpg";

            global.fetch = vi.fn().mockResolvedValue({
                json: () => Promise.resolve({
                    user: { ...mockUser, avatar: newAvatarUrl },
                    token: "new-token",
                }),
            });

            renderProfilePage();

            await waitFor(() => screen.getByText("AL"));
            fireEvent.click(screen.getByText("AL"));

            fireEvent.change(screen.getByPlaceholderText("Paste image URL…"), {
                target: { value: newAvatarUrl },
            });

            fireEvent.click(screen.getByText("Save"));

            await waitFor(() => {
                const img = screen.getByAltText("avatar");
                expect(img.src).toBe(newAvatarUrl);
            });
        });
    });
});
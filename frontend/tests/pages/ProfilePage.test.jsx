import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, waitFor } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import ProfilePage from "../../src/pages/Profile/ProfilePage.jsx";
import { fetchUserProfile } from "../../src/services/authentication";
import { getSpotsByUser } from "../../src/services/spots";
import { isTokenValid } from "../../src/helpers/authentication.js";


vi.mock("../../src/services/authentication", () => ({
  fetchUserProfile: vi.fn(),
  updateAvatar: vi.fn(),
}));
vi.mock("../../src/services/spots", () => ({
  getSpotsByUser: vi.fn(),
}));
vi.mock("../../src/helpers/authentication.js", () => ({
  isTokenValid: vi.fn(),
}));

const mockNavigate = vi.fn();
vi.mock("react-router-dom", async () => {
  const actual = await vi.importActual("react-router-dom");
  return { ...actual, useNavigate: () => mockNavigate };
});

const mockUser = {
  username: "aluna",
  email: "aluna@gmail.com",
  createdAt: "2026-04-24T00:00:00Z",
  avatar: null,
};

const renderProfilePage = () =>
  render(<MemoryRouter><ProfilePage /></MemoryRouter>);

beforeEach(() => {
  vi.clearAllMocks();
  localStorage.setItem("token", "mock-token");
  isTokenValid.mockReturnValue(true);
  fetchUserProfile.mockResolvedValue({ user: mockUser, token: "new-token" });
  getSpotsByUser.mockResolvedValue([]);
});

describe("ProfilePage", () => {

  it("redirects to /login if token is invalid", async () => {
    isTokenValid.mockReturnValue(false);
    renderProfilePage();
    await waitFor(() => {
      expect(mockNavigate).toHaveBeenCalledWith("/login");
    });
  });

  it("shows loading text while fetching", () => {
    fetchUserProfile.mockReturnValue(new Promise(() => {}));
    renderProfilePage();
    expect(screen.getByText("Loading profile…")).toBeTruthy();
  });

  it("renders username and email after loading", async () => {
    renderProfilePage();
    await waitFor(() => {
      expect(screen.getByText("aluna")).toBeTruthy();
      expect(screen.getByText("aluna@gmail.com")).toBeTruthy();
    });
  });

  it("shows error message if profile fetch fails", async () => {
    fetchUserProfile.mockRejectedValue(new Error("Something went wrong"));
    renderProfilePage();
    await waitFor(() => {
      expect(screen.getByText("Something went wrong")).toBeTruthy();
    });
  });

  it("shows message when user has no spots", async () => {
    renderProfilePage();
    await waitFor(() => {
      expect(screen.getByText("You have not submitted any spots yet.")).toBeTruthy();
    });
  });

});
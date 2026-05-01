import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { afterEach, vi } from "vitest";

import { MemoryRouter} from "react-router-dom";
import { signup } from "../../src/services/authentication";

import { SignupPage } from "../../src/pages/Signup/SignupPage";

const navigateMock = vi.fn();

// Mocking React Router's useNavigate function
vi.mock("react-router-dom", async (importOriginal) => {
  const original = await importOriginal();
  return {... original, useNavigate: () => navigateMock };
});

// Mocking localStorage
const localStorageMock = {
  setItem: vi.fn(),
  getItem: vi.fn(),
  removeItem: vi.fn(),
  clear: vi.fn()
}
Object.defineProperty(window, "localStorage", { value: localStorageMock })


// Mocking the signup service
vi.mock("../../src/services/authentication", () => {
  const signupMock = vi.fn();
  return { signup: signupMock };
});

// Reusable function for filling out signup form
const completeSignupForm = async () => {
  const user = userEvent.setup();

  const usernameInputEl = screen.getByPlaceholderText("Username");
  const emailInputEl = screen.getByPlaceholderText("Email address");
  const passwordInputEl = screen.getByPlaceholderText("Create a password");
  const confirmPasswordInputEl = screen.getByPlaceholderText("Confirm password");
  const loginButtonEl = screen.getByRole("button", {name: /sign up/i});

  await user.type(usernameInputEl, "tester")
  await user.type(emailInputEl, "test@email.com");
  await user.type(passwordInputEl, "password1234");
  await user.type(confirmPasswordInputEl, "password1234")
  await user.click(loginButtonEl);
};

describe("Signup Page", () => {
  afterEach(() => {
    vi.resetAllMocks();
  });

  test("allows a user to signup", async () => {
    render(<MemoryRouter><SignupPage /></MemoryRouter>);

    await completeSignupForm();

    expect(signup).toHaveBeenCalledWith("tester", "test@email.com", "password1234", "password1234");
  });

  test("navigates to / on successful signup", async () => {
    vi.spyOn(Storage.prototype, "setItem").mockImplementation(() => {});
    signup.mockResolvedValue({token: "secrettoken123"});
    render(<MemoryRouter><SignupPage /></MemoryRouter>);

    await completeSignupForm();

    await waitFor(() => {
      expect(navigateMock).toHaveBeenCalledWith("/");
    })
  });

  test("shows an error message when the user is unsuccessful in loging in", async () => {
    vi.spyOn(Storage.prototype, "setItem").mockImplementation(() => {});
    signup.mockRejectedValue(new Error("Error signing up"));
    render(<MemoryRouter><SignupPage /></MemoryRouter>);

    await completeSignupForm();

    await waitFor(() => {
      expect(screen.getByText("Error signing up")).toBeTruthy();
    })
  });

  test("store the token to localStorage when the user signup successfully", async () => {
    vi.spyOn(Storage.prototype, "setItem").mockImplementation(() => {});
    signup.mockResolvedValue({token: "secrettoken123"});
    render(<MemoryRouter><SignupPage /></MemoryRouter>);
    await completeSignupForm();

    await waitFor(() => {
      expect(localStorageMock.setItem).toHaveBeenCalledWith("token", "secrettoken123");
    })
  });
});

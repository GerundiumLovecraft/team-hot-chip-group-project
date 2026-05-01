import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { vi } from "vitest";

import { MemoryRouter} from "react-router-dom";
import { login } from "../../src/services/authentication";

import { LoginPage } from "../../src/pages/Login/LoginPage";

const navigateMock = vi.fn();

// Mocking React Router's useNavigate function
vi.mock("react-router-dom", async(importOriginal) => {
  const original = await importOriginal()
  return { ...original, useNavigate: () => navigateMock };
});

// Mocking localStorage
const localStorageMock = {
  setItem: vi.fn(),
  getItem: vi.fn(),
  removeItem: vi.fn(),
  clear: vi.fn()
};
Object.defineProperty(window, "localStorage", { value: localStorageMock })

// Mocking the login service
vi.mock("../../src/services/authentication", () => {
  const loginMock = vi.fn();
  return { login: loginMock };
});

// Reusable function for filling out login form
const completeLoginForm = async () => {
  const user = userEvent.setup();

  const emailOrUsernameInputEl = screen.getByPlaceholderText("Username or email");
  const passwordInputEl = screen.getByPlaceholderText("Password");
  const loginButtonEl = screen.getByRole("button", {name: /login/i});

  await user.type(emailOrUsernameInputEl, "test@email.com");
  await user.type(passwordInputEl, "password1234");
  await user.click(loginButtonEl);
};

describe("Login Page", () => {
  afterEach(() => {
    vi.resetAllMocks();
  })
  
  test("allows a user to login", async () => {
    render(<MemoryRouter><LoginPage/></MemoryRouter>);
    await completeLoginForm();

    expect(login).toHaveBeenCalledWith("test@email.com", "password1234");
  });

  test("navigates to / on successful login", async () => {
    vi.spyOn(Storage.prototype, "setItem").mockImplementation(() => {});
    login.mockResolvedValue("secrettoken123");
    render(<MemoryRouter><LoginPage/></MemoryRouter>);
    await completeLoginForm();
    await waitFor(() => {
      expect(navigateMock).toHaveBeenCalledWith("/");
    })
  });

  test("shows an error message when the user is unsuccessful in loging in", async () => {
    vi.spyOn(Storage.prototype, "setItem").mockImplementation(() => {});
    login.mockRejectedValue(new Error("Error logging in"));
    render(<MemoryRouter><LoginPage/></MemoryRouter>);

    await completeLoginForm();

    await waitFor(() => {
      expect(screen.getByText("Invalid username/email or password.")).toBeTruthy();
    })
  });

  test("store the token to localStorage when the user login successfully", async () => {
    vi.spyOn(Storage.prototype, "setItem").mockImplementation(() => {});
    login.mockResolvedValue("secrettoken123");
    render(<MemoryRouter><LoginPage/></MemoryRouter>);
    await completeLoginForm();
    await waitFor(() => {
      expect(localStorageMock.setItem).toHaveBeenCalledWith("token", "secrettoken123");
    })
  })
});

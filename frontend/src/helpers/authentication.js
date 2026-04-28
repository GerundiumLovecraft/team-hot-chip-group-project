/*
There are two cases where authentication will fail:
    - No token is present in the local storage
    - Token is expired
Checks will be done in this helper function to reduce the amount of written code on separate pages
 */
export function isTokenValid() {
    // Get the token from the local storage
    const token = localStorage.getItem("token")

    // Check if token is present
    if (!token) return false;

    /*
    Token consists of three parts separated by ".": header.payload.signature
    We need to extract the payload from the token and decode from base64 into the readable format
    Once that will be done, we will have { "sub": ..., "iat":..., "exp":... } object
     */
    try{
        const payload = JSON.parse(atob(token.split(".")[1]));
        // Times in payload are in seconds, Date.now() returns the value in ms
        const isValid = payload.exp * 1000 > Date.now();

        if (!isValid) {
            localStorage.removeItem("token");
            return isValid;
        }

        return isValid;
    } catch {
        localStorage.removeItem("token");
        return false;
    }

}
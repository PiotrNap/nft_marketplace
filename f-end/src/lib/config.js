const baseUrl = process.env.NEXT_PUBLIC_API_URL

module.exports = {
  API: {
    CHECK_USERNAME: baseUrl + "/users/username",
    GET_CHALLENGE: baseUrl + "/auth/challenge",
    SIGN_UP: baseUrl + "/auth/sign-up",
  },
}

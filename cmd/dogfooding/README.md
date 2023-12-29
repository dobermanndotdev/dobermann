# Dogfooding

Goals

- Load a dataset of many well-known urls
- Create monitors based on those urls
- Observe how the system behaves (ideally I should have email notifications turned off to preserve credits)

## How to run

1. Create a `.env.secrets` in the root of this repo
2. Log in to Dobermann and copy the JWT from the cookies
3. Define the variable `DOGFOODING_JWT=<the-jwt>` in the `.env.secrets` file
4. Run it.


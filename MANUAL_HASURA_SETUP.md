# Manual Hasura Actions Setup

Since automated scripts may not work perfectly, here's how to configure Hasura actions manually:

## Step 1: Open Hasura Console

Open: http://localhost:8080

Secret: `myhasurasecret`

---

## Step 2: Configure Login Action

1. Go to **Actions** tab
2. Click **Create**
3. Enter this configuration:

### Action Definition:
```graphql
type Mutation {
  login(arg: LoginInput!): LoginOutput!
}
```

### New types definition:
```graphql
input LoginInput {
  email: String!
  password: String!
}

type LoginOutput {
  token: String!
  user_id: Int!
  name: String!
  email: String!
}
```

### Handler:
```
http://host.docker.internal:8081/hasura/login
```

4. Click **Create Action**

---

## Step 3: Configure Signup Action

1. Still in **Actions** tab
2. Click **Create** again
3. Enter this configuration:

### Action Definition:
```graphql
type Mutation {
  signup(arg: SignupInput!): SignupOutput!
}
```

### New types definition:
```graphql
input SignupInput {
  name: String!
  email: String!
  password: String!
}

type SignupOutput {
  id: Int!
  name: String!
  email: String!
}
```

### Handler:
```
http://host.docker.internal:8081/hasura/signup
```

4. Click **Create Action**

---

## Step 4: Test

Once both actions are created:

1. Refresh your frontend: http://localhost:3000
2. Try registering an account
3. Try logging in

The GraphQL mutations should now work!


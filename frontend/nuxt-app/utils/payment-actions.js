import gql from 'graphql-tag';

const normalizeInitializeArg = (arg) => {
  return {
    amount: String(arg?.amount ?? ''),
    email: String(arg?.email ?? ''),
    user_name: String(arg?.user_name ?? '').trim() || 'Customer',
    recipe_id: Number.isInteger(arg?.recipe_id) ? arg.recipe_id : Number(arg?.recipe_id || 0)
  };
};

const INITIALIZE_PAYMENT_MUTATION = gql`
  mutation InitializePayment(
    $amount: String!
    $email: String!
    $user_name: String!
    $recipe_id: Int!
  ) {
    result: initializePayment(
      amount: $amount
      email: $email
      user_name: $user_name
      recipe_id: $recipe_id
    ) {
      status
      checkout_url
      tx_ref
      message
    }
  }
`;

const VERIFY_PAYMENT_MUTATION = gql`
  mutation VerifyPayment($tx_ref: String, $recipe_id: Int) {
    result: verifyPayment(tx_ref: $tx_ref, recipe_id: $recipe_id) {
      status
      message
    }
  }
`;

const toError = (errors, fallbackMessage) => {
  if (!errors || errors.length === 0) {
    return new Error(fallbackMessage);
  }
  const message = errors
    .map((e) => String(e?.message || '').trim())
    .filter(Boolean)
    .join('; ');
  return new Error(message || fallbackMessage);
};

const extractActionResult = (data) => {
  if (!data || typeof data !== 'object') {
    return undefined;
  }

  if (Object.prototype.hasOwnProperty.call(data, 'result')) {
    return data.result;
  }

  const keys = Object.keys(data);
  if (keys.length === 1) {
    return data[keys[0]];
  }

  return undefined;
};

const runActionMutation = async (client, mutation, variables) => {
  const res = await client.mutate({
    mutation,
    variables,
    fetchPolicy: 'no-cache'
  });

  const graphQLErrors = Array.isArray(res?.errors) ? res.errors : [];
  if (graphQLErrors.length > 0) {
    throw toError(graphQLErrors, 'Payment action mutation failed');
  }

  const result = extractActionResult(res?.data);
  if (result === undefined || result === null) {
    throw new Error('Missing result in action response');
  }
  return result;
};

export const initializePaymentAction = (client, arg) => {
  const payload = normalizeInitializeArg(arg);
  return runActionMutation(client, INITIALIZE_PAYMENT_MUTATION, {
    amount: payload.amount,
    email: payload.email,
    user_name: payload.user_name,
    recipe_id: payload.recipe_id
  });
};

export const verifyPaymentAction = (client, arg) => {
  return runActionMutation(client, VERIFY_PAYMENT_MUTATION, {
    tx_ref: arg?.tx_ref ?? null,
    recipe_id: Number.isInteger(arg?.recipe_id) ? arg.recipe_id : null
  });
};

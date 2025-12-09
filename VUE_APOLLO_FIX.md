# ✅ Vue Apollo Fix - useAsyncQuery Error

## ❌ The Problem

**Error**:
```
Uncaught SyntaxError: The requested module does not provide an export named 'useAsyncQuery'
```

**Why**: `useAsyncQuery` **does not exist** in `@vue/apollo-composable`!

---

## ✅ The Solution

### Correct Vue Apollo Pattern:

**DON'T USE** (doesn't exist):
```typescript
// ❌ WRONG - useAsyncQuery doesn't exist
const { data, error, pending } = await useAsyncQuery(query, { id: recipeId });
```

**USE THIS INSTEAD**:
```typescript
// ✅ CORRECT - useQuery is the right composable
const { result, loading, error } = useQuery(query, { id: recipeId });
const recipe = computed(() => result.value?.recipes_by_pk);
```

---

## 📚 Vue Apollo Composables Reference

### 1. **useQuery** - For GraphQL Queries
```typescript
import { useQuery } from '@vue/apollo-composable';
import gql from 'graphql-tag';

const QUERY = gql`
  query GetRecipe($id: Int!) {
    recipe(id: $id) { id title }
  }
`;

const { result, loading, error, refetch } = useQuery(QUERY, { id: recipeId });

// Access data
const recipe = computed(() => result.value?.recipe);
```

### 2. **useMutation** - For GraphQL Mutations
```typescript
import { useMutation } from '@vue/apollo-composable';

const MUTATION = gql`
  mutation CreateRecipe($input: RecipeInput!) {
    createRecipe(input: $input) { id }
  }
`;

const { mutate, loading, error } = useMutation(MUTATION);

// Call mutation
await mutate({ input: { title: 'New Recipe' } });
```

### 3. **useSubscription** - For GraphQL Subscriptions
```typescript
import { useSubscription } from '@vue/apollo-composable';

const SUBSCRIPTION = gql`
  subscription OnRecipeAdded {
    recipeAdded { id title }
  }
`;

const { result } = useSubscription(SUBSCRIPTION);
```

---

## 🔄 What Changed in recipes/[id].vue

### Before (Broken):
```typescript
import { useQuery, useAsyncQuery } from '@vue/apollo-composable'; // ❌

const { data, error, pending } = await useAsyncQuery(query, { id: recipeId }); // ❌
const recipe = computed(() => data.value?.recipes_by_pk);
```

### After (Fixed):
```typescript
import { useQuery } from '@vue/apollo-composable'; // ✅

const { result, loading: pending, error } = useQuery(query, { id: recipeId }); // ✅
const recipe = computed(() => result.value?.recipes_by_pk);
```

---

## ✅ Key Differences

| Feature | useAsyncQuery (doesn't exist) | useQuery (correct) |
|---------|------------------------------|-------------------|
| Import | ❌ Not available | ✅ `@vue/apollo-composable` |
| Returns | ❌ N/A | ✅ `{ result, loading, error }` |
| Async | ❌ N/A | ✅ Reactive (no await needed) |
| Data Access | ❌ N/A | ✅ `result.value?.data` |

---

## 📖 Additional Notes

### Reactive Data
`useQuery` returns **reactive refs**, so no `await` is needed:
```typescript
const { result, loading } = useQuery(QUERY);

// Data updates automatically when query completes
watch(result, (newValue) => {
  console.log('Data updated:', newValue);
});
```

### Manual Refetch
```typescript
const { result, refetch } = useQuery(QUERY);

// Manually refetch
await refetch();
```

### With Variables
```typescript
const variables = ref({ id: 1 });
const { result } = useQuery(QUERY, variables);

// Change variables to trigger new query
variables.value.id = 2;
```

---

## ✅ Status

- ✅ **Error Fixed**: No more `useAsyncQuery` import error
- ✅ **Correct Pattern**: Using `useQuery` properly
- ✅ **Website Working**: Recipe detail page loads correctly

**Website**: http://localhost:3000  
**Test**: http://localhost:3000/recipes/1

All Vue Apollo queries now use the correct composables! 🎉


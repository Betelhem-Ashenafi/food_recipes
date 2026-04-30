import gql from 'graphql-tag';

export const RECIPE_QUERY = gql`
  query GetRecipe($id: Int!) {
    recipes_by_pk(id: $id) {
      id
      title
      description
      price
      created_at
      preparation_time
      user_id
      category_id
      user {
        name
      }
    }
  }
`;

export const GET_ACCESSIBLE_RECIPE_QUERY = gql`
  query GetAccessibleRecipe($userId: Int!, $recipeId: Int!) {
    get_accessible_recipe(args: { p_user_id: $userId, p_recipe_id: $recipeId }) {
      id
      title
      description
      price
    }
  }
`;

export const CATEGORY_QUERY = gql`
  query GetCategory($id: Int!) {
    categories_by_pk(id: $id) {
      id
      name
      image_url
    }
  }
`;

export const CATEGORIES_QUERY = gql`
  query GetCategories {
    categories(order_by: { name: asc }) {
      id
      name
      image_url
    }
  }
`;

export const RECIPE_IMAGES_QUERY = gql`
  query GetRecipeImages($id: Int!) {
    recipe_images(where: { recipe_id: { _eq: $id } }, order_by: { id: asc }) {
      id
      url
      is_featured
    }
  }
`;

export const RECIPE_INGREDIENTS_QUERY = gql`
  query GetRecipeIngredients($id: Int!) {
    recipe_ingredients(where: { recipe_id: { _eq: $id } }, order_by: { id: asc }) {
      id
      name
      quantity
      unit_id
      unit {
        id
        name
      }
    }
  }
`;

export const RECIPE_STEPS_QUERY = gql`
  query GetRecipeSteps($id: Int!) {
    recipe_steps(where: { recipe_id: { _eq: $id } }, order_by: { step_number: asc }) {
      id
      step_number
      instruction
    }
  }
`;

export const UNITS_QUERY = gql`
  query GetUnits {
    units(order_by: { id: asc }) {
      id
      name
    }
  }
`;

export const RECIPE_COMMENTS_QUERY = gql`
  query GetRecipeComments($id: Int!) {
    comments(where: { recipe_id: { _eq: $id } }, order_by: { created_at: desc }) {
      id
      user_id
      content
      created_at
      user {
        name
      }
    }
  }
`;

export const RECIPE_COMMENTS_QUERY_CAMEL = gql`
  query GetRecipeCommentsCamel($id: Int!) {
    comments(where: { recipeId: { _eq: $id } }, order_by: { createdAt: desc }) {
      id
      userId
      content
      createdAt
      user {
        name
      }
    }
  }
`;

export const RECIPE_RATING_QUERY = gql`
  query GetRecipeRating($id: Int!) {
    recipes_by_pk(id: $id) {
      recipe_average_rating
      recipe_likes_count
    }
  }
`;

export const CHECK_PURCHASES_RATINGS_QUERY = gql`
  query CheckPurchasesRatings($recipeId: Int!, $userId: Int!) {
    purchases(where: { recipe_id: { _eq: $recipeId }, user_id: { _eq: $userId }, status: { _eq: "success" } }) {
      id
      status
    }
    ratings(where: { recipe_id: { _eq: $recipeId }, user_id: { _eq: $userId } }) {
      rating
    }
  }
`;

export const CHECK_SOCIAL_QUERY = gql`
  query CheckSocial($recipeId: Int!, $userId: Int!) {
    likes(where: { recipe_id: { _eq: $recipeId }, user_id: { _eq: $userId } }) {
      recipe_id
    }
    bookmarks(where: { recipe_id: { _eq: $recipeId }, user_id: { _eq: $userId } }) {
      recipe_id
    }
  }
`;

export const LIKE_RECIPE_MUTATION = gql`
  mutation LikeRecipe($recipeId: Int!) {
    insert_likes_one(object: { recipe_id: $recipeId }) {
      recipe_id
    }
  }
`;

export const UNLIKE_RECIPE_MUTATION = gql`
  mutation UnlikeRecipe($recipeId: Int!, $userId: Int!) {
    delete_likes(where: { recipe_id: { _eq: $recipeId }, user_id: { _eq: $userId } }) {
      affected_rows
    }
  }
`;

export const BOOKMARK_RECIPE_MUTATION = gql`
  mutation BookmarkRecipe($recipeId: Int!) {
    insert_bookmarks_one(object: { recipe_id: $recipeId }) {
      recipe_id
    }
  }
`;

export const UNBOOKMARK_RECIPE_MUTATION = gql`
  mutation UnbookmarkRecipe($recipeId: Int!, $userId: Int!) {
    delete_bookmarks(where: { recipe_id: { _eq: $recipeId }, user_id: { _eq: $userId } }) {
      affected_rows
    }
  }
`;

export const RATE_RECIPE_MUTATION = gql`
  mutation RateRecipe($recipeId: Int!, $rating: Int!) {
    insert_ratings_one(
      object: { recipe_id: $recipeId, rating: $rating }
      on_conflict: { constraint: ratings_recipe_user_unique, update_columns: [rating] }
    ) {
      rating
    }
  }
`;

export const POST_COMMENT_MUTATION = gql`
  mutation PostComment($recipeId: Int!, $content: String!) {
    insert_comments_one(object: { recipe_id: $recipeId, content: $content }) {
      id
    }
  }
`;

export const POST_COMMENT_MUTATION_CAMEL = gql`
  mutation PostCommentCamel($recipeId: Int!, $userId: Int!, $content: String!) {
    insert_comments_one(object: { recipeId: $recipeId, userId: $userId, content: $content }) {
      id
    }
  }
`;

export const UPDATE_COMMENT_MUTATION = gql`
  mutation UpdateComment($id: Int!, $content: String!) {
    update_comments(where: { id: { _eq: $id } }, _set: { content: $content }) {
      affected_rows
    }
  }
`;

export const DELETE_COMMENT_MUTATION = gql`
  mutation DeleteComment($id: Int!) {
    delete_comments(where: { id: { _eq: $id } }) {
      affected_rows
    }
  }
`;

export const UPDATE_RECIPE_INLINE_MUTATION = gql`
  mutation UpdateRecipeInline(
    $recipeId: Int!
    $recipe: recipes_set_input!
    $ingredients: [recipe_ingredients_insert_input!]!
    $steps: [recipe_steps_insert_input!]!
  ) {
    update_recipes_by_pk(pk_columns: { id: $recipeId }, _set: $recipe) {
      id
    }
    delete_recipe_ingredients(where: { recipe_id: { _eq: $recipeId } }) {
      affected_rows
    }
    insert_recipe_ingredients(objects: $ingredients) {
      affected_rows
    }
    delete_recipe_steps(where: { recipe_id: { _eq: $recipeId } }) {
      affected_rows
    }
    insert_recipe_steps(objects: $steps) {
      affected_rows
    }
  }
`;

// server/api/seed-recipes.js
export default defineEventHandler(async (event) => {
    const { data, errors } = await $fetch('http://localhost:8080/v1/graphql', {
      method: 'POST',
      headers: {
        'x-hasura-admin-secret': 'myadminsecretkey'
      },
      body: JSON.stringify({
        query: `
          mutation InsertRecipes {
            insert_recipes(objects: [
              {
                title: "Spaghetti Carbonara",
                description: "Classic Italian pasta dish",
                cooking_time: 20,
                difficulty: "Medium",
                image_url: "https://example.com/carbonara.jpg"
              },
              {
                title: "Chicken Curry",
                description: "Spicy Indian curry",
                cooking_time: 45,
                difficulty: "Easy",
                image_url: "https://example.com/curry.jpg"
              }
            ]) {
              returning {
                id
                title
              }
            }
          }
        `
      })
    })
  
    return data
  })
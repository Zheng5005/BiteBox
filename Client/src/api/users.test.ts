import { describe, it, expect } from 'vitest';
import { getUserRecipes } from './users';
import { server } from '../test/setup';
import { http, HttpResponse } from 'msw';

describe('getUserRecipes', () => {
  it('should correctly parse id and rating from strings to numbers', async () => {
    const mockRecipes = [
      {
        id: '1',
        name_recipe: 'Test Recipe 1',
        description: 'Description 1',
        meal_type_id: '10',
        img_url: 'http://example.com/image1.jpg',
        rating: '4.5',
      },
      {
        id: '2',
        name_recipe: 'Test Recipe 2',
        description: 'Description 2',
        meal_type_id: '20',
        img_url: 'http://example.com/image2.jpg',
        rating: '5.0',
      },
    ];

    server.use(
      http.get('http://localhost:8080/api/users', () => {
        return HttpResponse.json(mockRecipes);
      })
    );

    const recipes = await getUserRecipes();

    expect(recipes).toHaveLength(2);

    expect(typeof recipes[0].id).toBe('number');
    expect(recipes[0].id).toBe(1);
    expect(typeof recipes[0].rating).toBe('number');
    expect(recipes[0].rating).toBe(4.5);
    expect(typeof recipes[0].meal_type_id).toBe('number');
    expect(recipes[0].meal_type_id).toBe(10);

    expect(typeof recipes[1].id).toBe('number');
    expect(recipes[1].id).toBe(2);
    expect(typeof recipes[1].rating).toBe('number');
    expect(recipes[1].rating).toBe(5.0);
    expect(typeof recipes[1].meal_type_id).toBe('number');
    expect(recipes[1].meal_type_id).toBe(20);
  });
});

// This file contains E2E tests for the DungeonPlayground component
// It can be used with Cypress or similar E2E testing frameworks

describe('DungeonPlayground E2E Tests', () => {
  beforeEach(() => {
    // Visit the component playground page
    cy.visit('/playground');
    
    // Select the DungeonPlayground component
    cy.get('select').first().select('DungeonPlayground');
  });

  it('should load dungeons from the API', () => {
    // Intercept the API call to get dungeons
    cy.intercept('GET', '**/dungeons', { 
      statusCode: 200, 
      body: [
        { id: '1', name: 'Test Dungeon 1', floors: 3, difficulty: 'easy', createdAt: '2023-01-01', playerCount: 0 },
        { id: '2', name: 'Test Dungeon 2', floors: 5, difficulty: 'hard', createdAt: '2023-01-02', playerCount: 2 },
      ]
    }).as('getDungeons');
    
    // Wait for the API call to complete
    cy.wait('@getDungeons');
    
    // Verify the dungeons are displayed
    cy.contains('Test Dungeon 1').should('be.visible');
    cy.contains('easy').should('be.visible');
    cy.contains('3 floors').should('be.visible');
  });

  it('should create a new dungeon', () => {
    // Intercept the initial API call to get dungeons
    cy.intercept('GET', '**/dungeons', { 
      statusCode: 200, 
      body: [] 
    }).as('getDungeons');
    
    // Intercept the API call to create a dungeon
    cy.intercept('POST', '**/dungeons', {
      statusCode: 201,
      body: { 
        id: '3', 
        name: 'New E2E Test Dungeon', 
        floors: 4, 
        difficulty: 'medium', 
        createdAt: new Date().toISOString(), 
        playerCount: 0 
      }
    }).as('createDungeon');
    
    // Wait for the initial API call to complete
    cy.wait('@getDungeons');
    
    // Fill in the form
    cy.get('input[placeholder="Enter dungeon name"]').clear().type('New E2E Test Dungeon');
    
    // Set the number of floors to 4
    cy.get('input[aria-label="Number of Floors"]').clear().type('4');
    
    // Select medium difficulty
    cy.get('select').eq(1).select('medium');
    
    // Click the create button
    cy.contains('button', 'Create Dungeon').click();
    
    // Wait for the API call to complete
    cy.wait('@createDungeon');
    
    // Verify the new dungeon is displayed
    cy.contains('New E2E Test Dungeon').should('be.visible');
    cy.contains('medium').should('be.visible');
    cy.contains('4 floors').should('be.visible');
  });

  it('should allow switching between floors', () => {
    // Intercept the API call to get dungeons
    cy.intercept('GET', '**/dungeons', { 
      statusCode: 200, 
      body: [
        { id: '1', name: 'Test Dungeon', floors: 3, difficulty: 'easy', createdAt: '2023-01-01', playerCount: 0 }
      ]
    }).as('getDungeons');
    
    // Wait for the API call to complete
    cy.wait('@getDungeons');
    
    // Verify the floor tabs are displayed
    cy.contains('Floor 1').should('be.visible');
    cy.contains('Floor 2').should('be.visible');
    cy.contains('Floor 3').should('be.visible');
    
    // Click on Floor 2
    cy.contains('Floor 2').click();
    
    // Verify Floor 2 is active
    cy.contains('Floor 2').should('have.attr', 'aria-selected', 'true');
    
    // Click on Floor 3
    cy.contains('Floor 3').click();
    
    // Verify Floor 3 is active
    cy.contains('Floor 3').should('have.attr', 'aria-selected', 'true');
  });

  it('should handle API errors gracefully', () => {
    // Intercept the API call to get dungeons and return an error
    cy.intercept('GET', '**/dungeons', { 
      statusCode: 500, 
      body: { error: 'Internal Server Error' }
    }).as('getDungeonsError');
    
    // Wait for the API call to complete
    cy.wait('@getDungeonsError');
    
    // Verify the error message is displayed
    cy.contains('Failed to load dungeons. Please try again.').should('be.visible');
    
    // Intercept the retry API call and return success
    cy.intercept('GET', '**/dungeons', { 
      statusCode: 200, 
      body: [
        { id: '1', name: 'Test Dungeon', floors: 3, difficulty: 'easy', createdAt: '2023-01-01', playerCount: 0 }
      ]
    }).as('getDungeonsRetry');
    
    // Click the retry button
    cy.contains('button', 'Retry').click();
    
    // Wait for the API call to complete
    cy.wait('@getDungeonsRetry');
    
    // Verify the dungeon is displayed
    cy.contains('Test Dungeon').should('be.visible');
  });
}); 